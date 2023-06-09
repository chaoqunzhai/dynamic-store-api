package apis

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"
	sys "go-admin/app/admin/models"
	"go-admin/app/company/models"
	"go-admin/common/dto"
	"go-admin/global"
	"golang.org/x/crypto/bcrypt"
	"strings"

	customUser "go-admin/common/jwt/user"
)

type CompanyUserGetPage struct {
	dto.Pagination `search:"-"`
	Name           string `form:"name"  search:"type:exact;column:name;table:sys_user" comment:""`
	Phone          string `form:"phone"  search:"type:exact;column:enable;table:sys_user" comment:""`
}

func (m *CompanyUserGetPage) GetNeedSearch() interface{} {
	return *m
}

type UpdateReq struct {
	Id       int    `uri:"id" comment:"主键编码"` // 主键编码
	Layer    int    `json:"layer" comment:"排序"`
	RoleId   int    `json:"role_id"`
	Status   string `json:"status" comment:"用户状态"`
	UserName string `json:"username" comment:"用户名称" binding:"required"`
	Phone    string `json:"phone" comment:"手机号"`
	PassWord string `json:"password" comment:"密码" binding:"required"`
}
type CategoryReq struct {
	Type int `json:"type" binding:"required"`
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
	UserId   int    `json:"user_id"`
	RoleName string `json:"role_name"`
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
	e.Orm.Model(&sys.SysUser{}).Select("user_id,username,password,avatar,email,sex,phone,"+
		"status,created_at,enable,layer,invitation_code").Where("c_id = ? and enable = ?", userDto.CId, true).Scopes(
		dto.MakeCondition(req.GetNeedSearch()),
		dto.Paginate(req.GetPageSize(), req.GetPageIndex()),
	).Order(global.OrderLayerKey).Find(&userLists).Count(&count)

	//角色查询
	cacheUserIds := make([]string, 0)
	for _, row := range userLists {
		cacheUserIds = append(cacheUserIds, fmt.Sprintf("%v", row.UserId))
	}
	userBindRoleMap := make(map[int]RoleUserRow, 0)
	if len(cacheUserIds) > 0 {
		roleMap := make([]RoleBindUser, 0)
		//关联的用户查询到角色ID
		sql := fmt.Sprintf("select * from company_role_user where user_id in (%v)",
			strings.Join(cacheUserIds, ","))
		e.Orm.Raw(sql).Scan(&roleMap)
		//角色ID+大BID查询到角色具体名称
		roleIds := make([]int, 0)
		roleBindUser := make(map[int]RoleUserRow, 0)
		for _, row := range roleMap {
			roleIds = append(roleIds, row.RoleId)
			roleBindUser[row.RoleId] = RoleUserRow{
				UserId: row.UserId,
			}
		}
		roleRows := make([]models.CompanyRole, 0)
		e.Orm.Model(&models.CompanyRole{}).Where("c_id = ? and enable = ? and id in ?",
			userDto.CId, true, roleIds).Find(&roleRows)
		//查询到的role 和 user做一个map关联
		for _, role := range roleRows {
			bindData, ok := roleBindUser[role.Id]
			if ok {
				userBindRoleMap[bindData.UserId] = RoleUserRow{
					RoleName: role.Name,
					RoleId:   role.Id,
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
			"disable": func() bool {
				if row.UserId == userDto.UserId {
					return true
				}
				return false
			}(),
		}
		if roleData, ok := userBindRoleMap[row.UserId]; ok {
			userRow["role"] = roleData.RoleName
			userRow["role_id"] = roleData.RoleId
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
	}
	var runSql string
	//更新第三张表角色ID
	if req.RoleId > 0 {
		var roleId int
		sql := fmt.Sprintf("select count(*) from  company_role_user where user_id = %v", req.Id)
		e.Orm.Raw(sql).Scan(&roleId)
		if roleId > 0 {
			runSql = fmt.Sprintf("update company_role_user set role_id = %v where user_id = %v", req.RoleId, req.Id)
		} else {
			runSql = fmt.Sprintf("INSERT INTO  company_role_user VALUES (%v,%v)", req.RoleId, req.Id)
		}

	} else {
		runSql = fmt.Sprintf("delete from company_role_user where user_id = %v", req.Id)
	}
	e.Orm.Exec(runSql)

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
		RoleId:   global.RoleCompanyUser,
		Layer:    req.Layer,
	}
	userObject.CreateBy = userDto.UserId
	e.Orm.Create(&userObject)
	if req.RoleId > 0 {
		runSql := fmt.Sprintf("INSERT INTO  company_role_user VALUES (%v,%v)", req.RoleId, userObject.UserId)
		e.Orm.Exec(runSql)
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
	e.OK("", "successful")
	return
}
