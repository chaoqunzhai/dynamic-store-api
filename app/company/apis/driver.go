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

type Driver struct {
	api.Api
}

// GetPage 获取Driver列表
// @Summary 获取Driver列表
// @Description 获取Driver列表
// @Tags Driver
// @Param layer query string false "排序"
// @Param enable query string false "开关"
// @Param cId query string false "大BID"
// @Param userId query string false "关联的用户ID"
// @Param name query string false "司机名称"
// @Param phone query string false "联系手机号"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.Driver}} "{"code": 200, "data": [...]}"
// @Router /api/v1/driver [get]
// @Security Bearer
func (e Driver) GetPage(c *gin.Context) {
    req := dto.DriverGetPageReq{}
    s := service.Driver{}
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
	list := make([]models.Driver, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取Driver失败，\r\n失败信息 %s", err.Error()))
        return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取Driver
// @Summary 获取Driver
// @Description 获取Driver
// @Tags Driver
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.Driver} "{"code": 200, "data": [...]}"
// @Router /api/v1/driver/{id} [get]
// @Security Bearer
func (e Driver) Get(c *gin.Context) {
	req := dto.DriverGetReq{}
	s := service.Driver{}
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
	var object models.Driver

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取Driver失败，\r\n失败信息 %s", err.Error()))
        return
	}

	e.OK( object, "查询成功")
}

// Insert 创建Driver
// @Summary 创建Driver
// @Description 创建Driver
// @Tags Driver
// @Accept application/json
// @Product application/json
// @Param data body dto.DriverInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/driver [post]
// @Security Bearer
func (e Driver) Insert(c *gin.Context) {
    req := dto.DriverInsertReq{}
    s := service.Driver{}
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
		e.Error(500, err, fmt.Sprintf("创建Driver失败，\r\n失败信息 %s", err.Error()))
        return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改Driver
// @Summary 修改Driver
// @Description 修改Driver
// @Tags Driver
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.DriverUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/driver/{id} [put]
// @Security Bearer
func (e Driver) Update(c *gin.Context) {
    req := dto.DriverUpdateReq{}
    s := service.Driver{}
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
		e.Error(500, err, fmt.Sprintf("修改Driver失败，\r\n失败信息 %s", err.Error()))
        return
	}
	e.OK( req.GetId(), "修改成功")
}

// Delete 删除Driver
// @Summary 删除Driver
// @Description 删除Driver
// @Tags Driver
// @Param data body dto.DriverDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/driver [delete]
// @Security Bearer
func (e Driver) Delete(c *gin.Context) {
    s := service.Driver{}
    req := dto.DriverDeleteReq{}
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
		e.Error(500, err, fmt.Sprintf("删除Driver失败，\r\n失败信息 %s", err.Error()))
        return
	}
	e.OK( req.GetId(), "删除成功")
}