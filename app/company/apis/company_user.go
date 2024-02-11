package apis

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"
	sys "go-admin/app/admin/models"
	"go-admin/app/company/models"
	"go-admin/common/business"
	"go-admin/common/dto"
	"go-admin/common/utils"
	"go-admin/config"
	"go-admin/global"
	"golang.org/x/crypto/bcrypt"
	"strings"

	customUser "go-admin/common/jwt/user"
)

type CompanyUserGetPage struct {
	dto.Pagination `search:"-"`
	UserName           string `form:"username"  search:"type:exact;column:username;table:sys_user" comment:""`
	Phone          string `form:"phone"  search:"type:exact;column:enable;table:sys_user" comment:""`
	Role int `form:"role"  search:"type:exact;column:role_id;table:sys_user" comment:""`
	All bool `json:"all" search:"-"`
	BeginTime      string `form:"beginTime" search:"type:gte;column:created_at;table:sys_user" comment:"创建时间"`
	EndTime        string `form:"endTime" search:"type:lte;column:created_at;table:sys_user" comment:"创建时间"`

}

func (m *CompanyUserGetPage) GetNeedSearch() interface{} {
	return *m
}

type UpdateReq struct {
	Id       int    `uri:"id" comment:"主键编码"` // 主键编码
	Layer    int    `json:"layer" comment:"排序"`
	Roles   []int    `json:"roles"` //给用户分配的角色ID
	ThisRole int `json:"this_role"` //系统用户必须是82 roleid
	Status   string `json:"status" comment:"用户状态"`
	UserName string `json:"username" comment:"用户名称" binding:"required"`
	Phone    string `json:"phone" comment:"手机号"`
	PassWord string `json:"password" comment:"密码" binding:"required"`
	AuthExamine bool `json:"auth_examine"`
}
type RenewPass struct {
	PasswordConfirm    string    `json:"password_confirm" gorm:"column:password_confirm"`
	UserName    string    `json:"user_name" gorm:"column:user_name"`
	RealName    string    `json:"real_name" gorm:"column:real_name"`
}
type CategoryReq struct {
	Type int `json:"type" binding:"required"`
}
type MakeCodeUser struct {
	Ids []int `json:"ids" binding:"required"`
}
type OfflineReq struct {
	Ids []int `json:"ids" binding:"required"`
}
type RoleBindUser struct {
	RoleId int `json:"role_id"`
	UserId int `json:"user_id"`
}
type RoleUserRow struct {
	RoleId   int    `json:"role_id"`
	RoleName string `json:"role_name"`
}

// todo:创建推广码
func (e Company) MakeCode(c *gin.Context) {
	req := MakeCodeUser{}
	err := e.MakeContext(c).
		Bind(&req, binding.JSON, nil).
		MakeOrm().
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	userDto, err := customUser.GetUserDto(e.Orm, c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}
	InvitationCode := utils.GenValidateCode(6)

	e.Orm.Model(&sys.SysUser{}).Where("c_id = ? and enable = ? and user_id in ?",
		userDto.CId, true, req.Ids).Updates(map[string]interface{}{
		"invitation_code": InvitationCode,
	})

	e.OK("", "successful")
	return
}


// 查询业务员的信息
func (e Company) PromotionCode(c *gin.Context) {
	err := e.MakeContext(c).
		MakeOrm().
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	userDto, err := customUser.GetUserDto(e.Orm, c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}

	url:=fmt.Sprintf("%v/%v",config.ExtConfig.PromotionCode,userDto.CId)
	e.OK(url, "successful")
	return
}

// 查询业务员的信息
func (e Company) MiniList(c *gin.Context) {
	req :=CompanyUserGetPage{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	userDto, err := customUser.GetUserDto(e.Orm, c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}
	datalist := make([]sys.SysUser, 0)

	//and invitation_code IS NOT NULL
	orm := e.Orm.Model(&sys.SysUser{}).Select("user_id,username")
	if req.All {
		orm = orm.Where("c_id = ?",userDto.CId)
	}else {
		orm = orm.Where("c_id = ? and enable = ? ",
			userDto.CId, true)
	}
	orm.Scopes(dto.MakeCondition(req.GetNeedSearch())).Find(&datalist)
	result := make([]map[string]interface{}, 0)
	for _, row := range datalist {
		result = append(result, map[string]interface{}{
			"id":   row.UserId,
			"name": row.Username,
		})
	}
	e.OK(result, "successful")
	return
}

func (e Company) List(c *gin.Context) {
	req := CompanyUserGetPage{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	userDto, err := customUser.GetUserDto(e.Orm, c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}

	var userLists []sys.SysUser
	var count int64
	//必须只能更新 大B下的用户,防止随意根据用户ID更改信息
	e.Orm.Model(&sys.SysUser{}).Where("c_id = ? and enable = ?", userDto.CId, true).Scopes(
		dto.MakeCondition(req.GetNeedSearch()),
		dto.Paginate(req.GetPageSize(), req.GetPageIndex()),
	).Order(global.OrderLayerKey).Find(&userLists).Count(&count)

	//角色查询
	cacheUserIds := make([]string, 0)
	for _, row := range userLists {
		cacheUserIds = append(cacheUserIds, fmt.Sprintf("%v", row.UserId))
	}
	RoleMap := make(map[int]RoleUserRow, 0)
	UserRoleMap:=make(map[int][]RoleUserRow,0)
	if len(cacheUserIds) > 0 {
		roleMap := make([]RoleBindUser, 0)
		//关联的用户查询到角色ID
		sql := fmt.Sprintf("select * from company_role_user where user_id in (%v)",
			strings.Join(cacheUserIds, ","))
		e.Orm.Raw(sql).Scan(&roleMap)

		roleIds := make([]int, 0)
		userBindRole:=make(map[int][]int,0)
		for _, row := range roleMap {
			roleIds = append(roleIds, row.RoleId)
			UserRoleMap[row.UserId] = make([]RoleUserRow,0)
			userBindRole[row.UserId] = append(userBindRole[row.UserId],row.RoleId)
		}
		if len(roleIds) > 0 {//统一查询角色
			roleRows := make([]models.CompanyRole, 0)
			e.Orm.Model(&models.CompanyRole{}).Where("c_id = ? and enable = ? and id in ?",
				userDto.CId, true, roleIds).Find(&roleRows)
			for _, role := range roleRows {
				RoleMap[role.Id] = RoleUserRow{
					RoleName: role.Name,
					RoleId:   role.Id,
				}
			}
		}
		for userIdKey,userRoles:=range userBindRole{
			userRoleList,ok:=UserRoleMap[userIdKey]
			if !ok{
				continue
			}
			for _,roleId:=range userRoles{
				getRow,rowOk:=RoleMap[roleId]
				if rowOk{
					userRoleList = append(userRoleList,getRow)
					UserRoleMap[userIdKey] = userRoleList
				}
			}
		}
	}
	result := make([]map[string]interface{}, 0)
	for _, row := range userLists {

		userRow := map[string]interface{}{
			"phone":           row.Phone,
			"username":        row.Username,
			"user_id":         row.UserId,
			"sex":             row.Sex,
			"password":        row.Password,
			"email":           row.Email,
			"invitation_code": row.InvitationCode,
			"avatar":          row.Avatar,
			"status":          row.Status,
			"layer":           row.Layer,
			"created_at":      row.CreatedAt,
			"auth_examine":row.AuthExamine,
			"disable": func() bool {
				if row.UserId == userDto.UserId {
					return true
				}
				return false
			}(),
		}
		if roleData, ok := UserRoleMap[row.UserId]; ok {
			userRow["roles"] = roleData
		}
		result = append(result, userRow)
	}
	e.PageOK(result, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")

	return
}
func (e Company) UpdateUser(c *gin.Context) {
	req := UpdateReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req, binding.JSON, nil).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	userDto, err := customUser.GetUserDto(e.Orm, c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}
	var userObject sys.SysUser
	//必须只能更新 大B下的用户,防止随意根据用户ID更改信息
	e.Orm.Model(&sys.SysUser{}).Where("c_id = ? and user_id = ? ", userDto.CId, req.Id).Limit(1).Find(&userObject)
	if userObject.UserId == 0 {
		e.Error(500, nil, "用户不存在")
		return
	}
	var validUser sys.SysUser
	e.Orm.Model(&sys.SysUser{}).Where("username = ? and c_id = ? and enable = ?",
		req.UserName, userDto.CId, true).Limit(1).Find(&validUser)
	if validUser.UserId != req.Id {
		e.Error(500, errors.New("用户名已经存在"), "用户名已经存在")
		return
	}
	e.Orm.Model(&sys.SysUser{}).Where("phone = ? and c_id = ? and enable = ?",
		req.Phone, userDto.CId, true).Limit(1).Find(&validUser)
	if validUser.UserId != req.Id {
		e.Error(500, errors.New("手机号已经存在"), "手机号已经存在")
		return
	}

	updateMap := map[string]interface{}{
		"username": req.UserName,
		"phone":    req.Phone,
		"layer":    req.Layer,
		"status":   req.Status,
		"enable":true,
		"role_id":req.ThisRole,
		"auth_examine": req.AuthExamine,
	}

	//先情况
	e.Orm.Exec(fmt.Sprintf("delete from company_role_user where user_id = %v", req.Id))
	//更新第三张表角色ID
	if len(req.Roles) > 0 {
		for _,roleId:=range req.Roles {
			e.Orm.Exec(fmt.Sprintf("INSERT INTO  company_role_user VALUES (%v,%v)", roleId, req.Id))
		}

	}

	//密码更新
	if req.PassWord != userObject.Password {
		var hash []byte
		var hasErr error
		if hash, hasErr = bcrypt.GenerateFromPassword([]byte(req.PassWord), bcrypt.DefaultCost); hasErr != nil {
			e.Error(500, err, "密码生成失败")
			return
		} else {
			updateMap["password"] = string(hash)
		}
	}
	e.Orm.Model(&sys.SysUser{}).Where("c_id = ? and user_id = ?", userDto.CId, req.Id).Updates(&updateMap)
	e.OK("", "successful")
	return
}


func (e Company) CreateSalesManUser(c *gin.Context) {
	req := UpdateReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req, binding.JSON, nil).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	userDto, err := customUser.GetUserDto(e.Orm, c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}
	if req.ThisRole == 0 {
		e.Error(500, errors.New("请选择角色"), "请选择角色")
		return
	}
	if req.ThisRole != global.RoleSaleMan{
		e.Error(500, errors.New("非法角色"), "非法角色")
		return
	}
	//检测业务员数量配置
	CompanyCnf := business.GetCompanyCnf(userDto.CId, "salesman_number", e.Orm)
	MaxNumber := CompanyCnf["salesman_number"]

	var thisCount int64
	e.Orm.Model(&sys.SysUser{}).Where("role_id = ? and c_id = ?",global.RoleSaleMan,userDto.CId).Count(&thisCount)

	if thisCount >= int64(MaxNumber) {
		msg:=fmt.Sprintf("最多只可创建%v个业务员",MaxNumber)
		e.Error(500, errors.New(msg), msg)
		return
	}
	//大B下的用户名是唯一的
	var count int64
	e.Orm.Model(&sys.SysUser{}).Where("username = ? and c_id = ? and enable = ?", req.UserName, userDto.CId, true).Count(&count)
	if count > 0 {
		e.Error(500, errors.New("用户名已经存在"), "用户名已经存在")
		return
	}
	var phoneCount int64
	e.Orm.Model(&sys.SysUser{}).Where("phone = ? and c_id = ? and enable = ?", req.Phone, userDto.CId, true).Count(&phoneCount)
	if phoneCount > 0 {
		e.Error(500, errors.New("手机号已经存在"), "手机号已经存在")
		return
	}
	userObject := sys.SysUser{
		Username: req.UserName,
		Phone:    req.Phone,
		Enable:   true,
		Status:   req.Status,
		Password: req.PassWord,
		CId:      userDto.CId,
		RoleId:   global.RoleSaleMan,
		Layer:    req.Layer,
		AuthExamine: req.AuthExamine,
	}
	userObject.CreateBy = userDto.UserId
	e.Orm.Create(&userObject)
	e.OK("successful", "创建成功")
	return
}

func (e Company) CreateUser(c *gin.Context) {
	req := UpdateReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req, binding.JSON, nil).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	userDto, err := customUser.GetUserDto(e.Orm, c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}
	if req.ThisRole == 0 {
		e.Error(500, errors.New("请选择角色"), "请选择角色")
		return
	}
	if req.ThisRole != global.RoleCompanyUser{
		e.Error(500, errors.New("非法角色"), "非法角色")
		return
	}

	//大B下的用户名是唯一的
	var count int64
	e.Orm.Model(&sys.SysUser{}).Where("username = ? and c_id = ? and enable = ?", req.UserName, userDto.CId, true).Count(&count)
	if count > 0 {
		e.Error(500, errors.New("用户名已经存在"), "用户名已经存在")
		return
	}
	var phoneCount int64
	e.Orm.Model(&sys.SysUser{}).Where("phone = ? and enable = ?", req.Phone,true).Count(&phoneCount)
	if phoneCount > 0 {
		e.Error(500, errors.New("手机号已经存在"), "手机号已经存在")
		return
	}
	userObject := sys.SysUser{
		Username: req.UserName,
		Phone:    req.Phone,
		Enable:   true,
		Status:   req.Status,
		Password: req.PassWord,
		CId:      userDto.CId,
		RoleId:   global.RoleCompanyUser,
		Layer:    req.Layer,
		AuthExamine: req.AuthExamine,
	}
	userObject.CreateBy = userDto.UserId
	e.Orm.Create(&userObject)
	if len(req.Roles) > 0 {
		for _,roleId:=range req.Roles{
			runSql := fmt.Sprintf("INSERT INTO  company_role_user VALUES (%v,%v)", roleId, userObject.UserId)
			e.Orm.Exec(runSql)
		}
	}
	e.OK("successful", "创建成功")
	return
}
func (e Company) Offline(c *gin.Context) {
	req := OfflineReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req, binding.JSON, nil).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	userDto, err := customUser.GetUserDto(e.Orm, c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}
	//必须只能更新 大B下的用户,防止随意根据用户ID更改信息
	for _, row := range req.Ids {
		var user sys.SysUser
		//如果是自己的用户,不可操作
		if row == userDto.UserId {
			continue
		}
		e.Orm.Model(&sys.SysUser{}).Where("c_id = ? and user_id = ?", userDto.CId, row).Limit(1).Find(&user)
		if user.UserId == 0 {
			continue
		}
		e.Orm.Model(&sys.SysUser{}).Where("c_id = ? and user_id = ?", userDto.CId, row).Updates(map[string]interface{}{
			"enable": false,
			"status": global.SysUserDisable,
		})
		//删除角色和用户的关联
		runSql := fmt.Sprintf("delete from company_role_user where user_id = %v", row)
		e.Orm.Exec(runSql)

	}
	e.OK("", "操作成功")
	return
}
