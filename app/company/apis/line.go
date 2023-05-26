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

type Line struct {
	api.Api
}

// GetPage 获取Line列表
// @Summary 获取Line列表
// @Description 获取Line列表
// @Tags Line
// @Param layer query string false "排序"
// @Param enable query string false "开关"
// @Param cId query string false "大BID"
// @Param name query string false "路线名称"
// @Param driverId query string false "关联司机"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.Line}} "{"code": 200, "data": [...]}"
// @Router /api/v1/line [get]
// @Security Bearer
func (e Line) GetPage(c *gin.Context) {
    req := dto.LineGetPageReq{}
    s := service.Line{}
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
	list := make([]models.Line, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取Line失败，\r\n失败信息 %s", err.Error()))
        return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取Line
// @Summary 获取Line
// @Description 获取Line
// @Tags Line
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.Line} "{"code": 200, "data": [...]}"
// @Router /api/v1/line/{id} [get]
// @Security Bearer
func (e Line) Get(c *gin.Context) {
	req := dto.LineGetReq{}
	s := service.Line{}
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
	var object models.Line

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取Line失败，\r\n失败信息 %s", err.Error()))
        return
	}

	e.OK( object, "查询成功")
}

// Insert 创建Line
// @Summary 创建Line
// @Description 创建Line
// @Tags Line
// @Accept application/json
// @Product application/json
// @Param data body dto.LineInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/line [post]
// @Security Bearer
func (e Line) Insert(c *gin.Context) {
    req := dto.LineInsertReq{}
    s := service.Line{}
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
		e.Error(500, err, fmt.Sprintf("创建Line失败，\r\n失败信息 %s", err.Error()))
        return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改Line
// @Summary 修改Line
// @Description 修改Line
// @Tags Line
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.LineUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/line/{id} [put]
// @Security Bearer
func (e Line) Update(c *gin.Context) {
    req := dto.LineUpdateReq{}
    s := service.Line{}
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
		e.Error(500, err, fmt.Sprintf("修改Line失败，\r\n失败信息 %s", err.Error()))
        return
	}
	e.OK( req.GetId(), "修改成功")
}

// Delete 删除Line
// @Summary 删除Line
// @Description 删除Line
// @Tags Line
// @Param data body dto.LineDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/line [delete]
// @Security Bearer
func (e Line) Delete(c *gin.Context) {
    s := service.Line{}
    req := dto.LineDeleteReq{}
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
		e.Error(500, err, fmt.Sprintf("删除Line失败，\r\n失败信息 %s", err.Error()))
        return
	}
	e.OK( req.GetId(), "删除成功")
}