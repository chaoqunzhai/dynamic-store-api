package apis

import (

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"
	sys "go-admin/app/admin/models"
	"go-admin/app/shop/models"
	"go-admin/common/actions"
	"go-admin/common/business"
	cDto "go-admin/common/dto"
	customUser "go-admin/common/jwt/user"

	"go-admin/global"

)

type ShopRegisterList struct {
	api.Api
}
type ShopRegisterListGetPageReq struct {
	cDto.Pagination `search:"-"`
	Phone          string `form:"phone"  search:"type:contains;column:phone;table:company_register_user_verify" `
	Status           string `form:"status"  search:"type:exact;column:status;table:company_register_user_verify" `
	BeginTime      string `form:"beginTime" search:"type:gte;column:created_at;table:company_register_user_verify" `
	EndTime        string `form:"endTime" search:"type:lte;column:created_at;table:company_register_user_verify"`
}

func (m *ShopRegisterListGetPageReq) GetNeedSearch() interface{} {
	return *m
}
type DeleteReq struct {
	Ids []int `json:"ids"`
}
type UpdateReq struct {
	Id             int      `uri:"id" comment:"主键编码"` // 主键编码
	Status int `json:"status"`
	Info string `json:"info"`
}
func (s *UpdateReq) GetId() interface{} {
	return s.Id
}

func (e ShopRegisterList) GetPage(c *gin.Context) {
	req := ShopRegisterListGetPageReq{}
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
	p := actions.GetPermissionFromContext(c)
	list := make([]models.CompanyRegisterUserVerify, 0)
	var count int64

	var data models.CompanyRegisterUserVerify

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(req.GetNeedSearch()),
			cDto.Paginate(req.GetPageSize(), req.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).Order("created_at desc").
		Find(&list).Limit(-1).Offset(-1).
		Count(&count).Error

	result:=make([]interface{},0)

	for _,row:=range list{
		if row.Salesman > 0 {

			var user sys.SysUser
			e.Orm.Model(&user).Select("user_id,username").Where("c_id = ? and user_id = ?",userDto.CId,row.Salesman).Limit(1).Find(&user)
			if user.UserId > 0 {
				row.SalesmanUser = user.Username
			}
		}
		if AppTypeName :=global.GetAppTypeName(row.AppTypeName);AppTypeName!=""{
			row.AppTypeName = AppTypeName
		}
		result = append(result,row)
	}
	e.PageOK(result, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

func (e ShopRegisterList) Detail(c *gin.Context) {
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
	approveId :=c.Param("id")
	var data models.CompanyRegisterUserVerify
	e.Orm.Model(&models.CompanyRegisterUserVerify{}).Where("c_id = ? and id = ?",userDto.CId, approveId).Limit(1).Find(&data)

	if data.Id == 0 {
		e.OK(business.Response{Code:-1,Msg: "无此用户"},"")
		return
	}

	if data.Status == -1 {
		e.OK(business.Response{Code:-1,Msg: "已驳回"},"")
		return
	}
	if data.Status == 2 {
		e.OK(business.Response{Code:-1,Msg: "门店已创建"},"")
		return
	}
	e.OK(business.Response{Code:0,Data: data},"")
	return
}

func (e ShopRegisterList) Update(c *gin.Context) {

	req := UpdateReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	_, err = customUser.GetUserDto(e.Orm, c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}
	p := actions.GetPermissionFromContext(c)
	var data models.CompanyRegisterUserVerify

	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).First(&data, req.GetId())

	//必须是没有创建门店 才可以把状态扭转
	if data.Status == 2  {
		e.Error(500, nil,"门店已经创建")
		return
	}
	//更新描述和状态
	data.Status  = req.Status
	data.Info = req.Info

	if data.Status == 1 {
		//data.AdoptUser = userDto.Username
		//data.AdoptTime = models3.XTime{
		//	Time:time.Now(),
		//}
		////检测手机号是否已经存在
		//var phoneCount int64
		//e.Orm.Model(&sys.SysShopUser{}).Select("user_id").Where("c_id = ? and phone = ? ",data.CId,data.Phone).Count(&phoneCount)
		//if phoneCount > 0 {
		//	e.Error(500, errors.New("手机号已经存在"), "手机号已经存在")
		//	return
		//}
		////创建用户
		//shopUserDto:=sys.SysShopUser{
		//	Username: data.Value,
		//	NickName:data.Value,
		//	Phone: data.Phone,
		//	Password:data.Password,
		//	Enable: true,
		//	CId: data.CId,
		//	Status:global.SysUserSuccess,
		//	RoleId:global.RoleShop,
		//}
		////设置创建用户为大B的名字
		//shopUserDto.CreateBy = userDto.UserId
		//
		//e.Orm.Create(&shopUserDto)
		////创建的ID保存进去
		//data.ShopUserId = shopUserDto.UserId
		//


	}else {
		//搜索这个手机号是否是有门店的
		var shopCount int64
		e.Orm.Model(&models.Shop{}).Where("c_id = ? and phone = ?",data.CId,data.Phone).Count(&shopCount)
		if shopCount > 0 {
			e.Error(500, nil,"手机号以创建门店,不可驳回")
			return
		}

		//驳回后,删除用户,也就是软删除
		var sysUser sys.SysShopUser
		e.Orm.Model(&sys.SysShopUser{}).Where("c_id = ? and phone = ?",data.CId,data.Phone).Delete(&sysUser)
	}
	saveErr2:=e.Orm.Save(&data).Error
	if saveErr2 != nil{
		e.Error(500, saveErr2,"审核更新失败,请联系管理员")
		return
	}
	e.OK("", "审核成功")
}

func (e ShopRegisterList) Delete(c *gin.Context) {

	req := DeleteReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	p := actions.GetPermissionFromContext(c)
	var data models.CompanyRegisterUserVerify
	e.Orm.Model(&data).Scopes(
		actions.Permission(data.TableName(), p),
	).Unscoped().Delete(&data,req.Ids)

	e.OK("", "删除成功")
}