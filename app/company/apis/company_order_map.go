package apis

import (
    "fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"

	"go-admin/app/company/models"
	"go-admin/app/company/service"
	"go-admin/app/company/service/dto"
	"go-admin/common/actions"
)

type CompanyOrderMap struct {
	api.Api
}

// GetPage 获取CompanyOrderMap列表
// @Summary 获取CompanyOrderMap列表
// @Description 获取CompanyOrderMap列表
// @Tags CompanyOrderMap
// @Param layer query string false "排序"
// @Param enable query string false "开关"
// @Param cId query string false "公司ID"
// @Param type query string false "映射表的类型"
// @Param orderTable query string false "对应表的名称"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.CompanyOrderMap}} "{"code": 200, "data": [...]}"
// @Router /api/v1/company-order-map [get]
// @Security Bearer
func (e CompanyOrderMap) GetPage(c *gin.Context) {
    req := dto.CompanyOrderMapGetPageReq{}
    s := service.CompanyOrderMap{}
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
	list := make([]models.CompanyOrderMap, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取CompanyOrderMap失败，\r\n失败信息 %s", err.Error()))
        return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取CompanyOrderMap
// @Summary 获取CompanyOrderMap
// @Description 获取CompanyOrderMap
// @Tags CompanyOrderMap
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.CompanyOrderMap} "{"code": 200, "data": [...]}"
// @Router /api/v1/company-order-map/{id} [get]
// @Security Bearer
func (e CompanyOrderMap) Get(c *gin.Context) {
	req := dto.CompanyOrderMapGetReq{}
	s := service.CompanyOrderMap{}
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
	var object models.CompanyOrderMap

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取CompanyOrderMap失败，\r\n失败信息 %s", err.Error()))
        return
	}

	e.OK( object, "查询成功")
}

// Insert 创建CompanyOrderMap
// @Summary 创建CompanyOrderMap
// @Description 创建CompanyOrderMap
// @Tags CompanyOrderMap
// @Accept application/json
// @Product application/json
// @Param data body dto.CompanyOrderMapInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/company-order-map [post]
// @Security Bearer
func (e CompanyOrderMap) Insert(c *gin.Context) {
    req := dto.CompanyOrderMapInsertReq{}
    s := service.CompanyOrderMap{}
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

	err = s.Insert(&req)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("创建CompanyOrderMap失败，\r\n失败信息 %s", err.Error()))
        return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改CompanyOrderMap
// @Summary 修改CompanyOrderMap
// @Description 修改CompanyOrderMap
// @Tags CompanyOrderMap
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.CompanyOrderMapUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/company-order-map/{id} [put]
// @Security Bearer
func (e CompanyOrderMap) Update(c *gin.Context) {
    req := dto.CompanyOrderMapUpdateReq{}
    s := service.CompanyOrderMap{}
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

	err = s.Update(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("修改CompanyOrderMap失败，\r\n失败信息 %s", err.Error()))
        return
	}
	e.OK( req.GetId(), "修改成功")
}

// Delete 删除CompanyOrderMap
// @Summary 删除CompanyOrderMap
// @Description 删除CompanyOrderMap
// @Tags CompanyOrderMap
// @Param data body dto.CompanyOrderMapDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/company-order-map [delete]
// @Security Bearer
func (e CompanyOrderMap) Delete(c *gin.Context) {
    s := service.CompanyOrderMap{}
    req := dto.CompanyOrderMapDeleteReq{}
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

	err = s.Remove(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("删除CompanyOrderMap失败，\r\n失败信息 %s", err.Error()))
        return
	}
	e.OK( req.GetId(), "删除成功")
}
