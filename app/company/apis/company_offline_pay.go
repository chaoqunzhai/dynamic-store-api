package apis

import (
	"errors"
	"fmt"
	sys2 "go-admin/app/admin/models"
	"go-admin/common/business"
	customUser "go-admin/common/jwt/user"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"

	"go-admin/app/company/models"
	"go-admin/app/company/service"
	"go-admin/app/company/service/dto"
	"go-admin/common/actions"
)

type CompanyOfflinePay struct {
	api.Api
}

// GetPage 获取CompanyOfflinePay列表
// @Summary 获取CompanyOfflinePay列表
// @Description 获取CompanyOfflinePay列表
// @Tags CompanyOfflinePay
// @Param cId query string false "大BID"
// @Param name query string false "线下支付名称"
// @Param layer query string false "排序"
// @Param enable query string false "开关"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.CompanyOfflinePay}} "{"code": 200, "data": [...]}"
// @Router /api/v1/company-offline-pay [get]
// @Security Bearer
func (e CompanyOfflinePay) GetPage(c *gin.Context) {
    req := dto.CompanyOfflinePayGetPageReq{}
    s := service.CompanyOfflinePay{}
    err := e.MakeContext(c).
        MakeOrm().
        Bind(&req).
        MakeService(&s.Service).
        Errors
   	if err != nil {
   		e.Logger.Error(err)
   		e.Error(500, err, err.Error())
   		return
   	}

	p := actions.GetPermissionFromContext(c)
	list := make([]models.CompanyOfflinePay, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取CompanyOfflinePay失败，\r\n失败信息 %s", err.Error()))
        return
	}
	result:=make([]map[string]interface{},0)

	for _,row:=range list{
		var sys sys2.SysUser
		e.Orm.Model(&sys).Select("nick_name,user_id").Where("user_id = ? and c_id = ?",row.CreateBy,row.CId).Limit(1).Find(&sys)
		createUser :=""
		if sys.UserId > 0 {
			createUser = sys.NickName
		}
		result = append(result, map[string]interface{}{
			"id":row.Id,
			"name":        row.Name,
			"create_time": row.CreatedAt.Format("2006-01-02 15:04:05"),
			"create_user": createUser,
		})
	}
	e.PageOK(result, len(result), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取CompanyOfflinePay
// @Summary 获取CompanyOfflinePay
// @Description 获取CompanyOfflinePay
// @Tags CompanyOfflinePay
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.CompanyOfflinePay} "{"code": 200, "data": [...]}"
// @Router /api/v1/company-offline-pay/{id} [get]
// @Security Bearer
func (e CompanyOfflinePay) Get(c *gin.Context) {
	req := dto.CompanyOfflinePayGetReq{}
	s := service.CompanyOfflinePay{}
    err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	var object models.CompanyOfflinePay

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取CompanyOfflinePay失败，\r\n失败信息 %s", err.Error()))
        return
	}

	e.OK( object, "查询成功")
}

// Insert 创建CompanyOfflinePay
// @Summary 创建CompanyOfflinePay
// @Description 创建CompanyOfflinePay
// @Tags CompanyOfflinePay
// @Accept application/json
// @Product application/json
// @Param data body dto.CompanyOfflinePayInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/company-offline-pay [post]
// @Security Bearer
func (e CompanyOfflinePay) Insert(c *gin.Context) {
    req := dto.CompanyOfflinePayInsertReq{}
    s := service.CompanyOfflinePay{}
    err := e.MakeContext(c).
        MakeOrm().
        Bind(&req).
        MakeService(&s.Service).
        Errors
    if err != nil {
        e.Logger.Error(err)
        e.Error(500, err, err.Error())
        return
    }

	// 设置创建人
	req.SetCreateBy(user.GetUserId(c))
	userDto, err := customUser.GetUserDto(e.Orm, c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}
	var countAll int64
	e.Orm.Model(&models.CompanyOfflinePay{}).Where("c_id = ?", userDto.CId).Count(&countAll)

	CompanyCnf := business.GetCompanyCnf(userDto.CId, "offline_pay", e.Orm)
	MaxNumber := CompanyCnf["offline_pay"]

	if countAll >= int64(MaxNumber) {
		e.Error(500, errors.New(fmt.Sprintf("线下支付最多只可创建%v个", MaxNumber)), fmt.Sprintf("线下支付最多只可创建%v个", MaxNumber))
		return
	}

	req.CId = userDto.CId
	err = s.Insert(&req)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("创建CompanyOfflinePay失败，\r\n失败信息 %s", err.Error()))
        return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改CompanyOfflinePay
// @Summary 修改CompanyOfflinePay
// @Description 修改CompanyOfflinePay
// @Tags CompanyOfflinePay
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.CompanyOfflinePayUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/company-offline-pay/{id} [put]
// @Security Bearer
func (e CompanyOfflinePay) Update(c *gin.Context) {
    req := dto.CompanyOfflinePayUpdateReq{}
    s := service.CompanyOfflinePay{}
    err := e.MakeContext(c).
        MakeOrm().
        Bind(&req).
        MakeService(&s.Service).
        Errors
    if err != nil {
        e.Logger.Error(err)
        e.Error(500, err, err.Error())
        return
    }
	req.SetUpdateBy(user.GetUserId(c))
	p := actions.GetPermissionFromContext(c)
	userDto, err := customUser.GetUserDto(e.Orm, c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}

	req.CId = userDto.CId
	err = s.Update(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("修改CompanyOfflinePay失败，\r\n失败信息 %s", err.Error()))
        return
	}
	e.OK( req.GetId(), "修改成功")
}

// Delete 删除CompanyOfflinePay
// @Summary 删除CompanyOfflinePay
// @Description 删除CompanyOfflinePay
// @Tags CompanyOfflinePay
// @Param data body dto.CompanyOfflinePayDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/company-offline-pay [delete]
// @Security Bearer
func (e CompanyOfflinePay) Delete(c *gin.Context) {
    s := service.CompanyOfflinePay{}
    req := dto.CompanyOfflinePayDeleteReq{}
    err := e.MakeContext(c).
        MakeOrm().
        Bind(&req).
        MakeService(&s.Service).
        Errors
    if err != nil {
        e.Logger.Error(err)
        e.Error(500, err, err.Error())
        return
    }

	// req.SetUpdateBy(user.GetUserId(c))
	p := actions.GetPermissionFromContext(c)
	userDto, err := customUser.GetUserDto(e.Orm, c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}

	req.CId = userDto.CId
	err = s.Remove(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("删除CompanyOfflinePay失败，\r\n失败信息 %s", err.Error()))
        return
	}
	e.OK( req.GetId(), "删除成功")
}
