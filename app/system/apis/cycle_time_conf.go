package apis

import (
    "fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"

	"go-admin/app/system/models"
	"go-admin/app/system/service"
	"go-admin/app/system/service/dto"
	"go-admin/common/actions"
)

type CycleTimeConf struct {
	api.Api
}

// GetPage 获取CycleTimeConf列表
// @Summary 获取CycleTimeConf列表
// @Description 获取CycleTimeConf列表
// @Tags CycleTimeConf
// @Param layer query string false "排序"
// @Param enable query string false "开关"
// @Param cId query string false "大BID"
// @Param type query string false "类型,每天,每周"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.CycleTimeConf}} "{"code": 200, "data": [...]}"
// @Router /api/v1/cycle-time-conf [get]
// @Security Bearer
func (e CycleTimeConf) GetPage(c *gin.Context) {
    req := dto.CycleTimeConfGetPageReq{}
    s := service.CycleTimeConf{}
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
	list := make([]models.CycleTimeConf, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取CycleTimeConf失败，\r\n失败信息 %s", err.Error()))
        return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取CycleTimeConf
// @Summary 获取CycleTimeConf
// @Description 获取CycleTimeConf
// @Tags CycleTimeConf
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.CycleTimeConf} "{"code": 200, "data": [...]}"
// @Router /api/v1/cycle-time-conf/{id} [get]
// @Security Bearer
func (e CycleTimeConf) Get(c *gin.Context) {
	req := dto.CycleTimeConfGetReq{}
	s := service.CycleTimeConf{}
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
	var object models.CycleTimeConf

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取CycleTimeConf失败，\r\n失败信息 %s", err.Error()))
        return
	}

	e.OK( object, "查询成功")
}

// Insert 创建CycleTimeConf
// @Summary 创建CycleTimeConf
// @Description 创建CycleTimeConf
// @Tags CycleTimeConf
// @Accept application/json
// @Product application/json
// @Param data body dto.CycleTimeConfInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/cycle-time-conf [post]
// @Security Bearer
func (e CycleTimeConf) Insert(c *gin.Context) {
    req := dto.CycleTimeConfInsertReq{}
    s := service.CycleTimeConf{}
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
		e.Error(500, err, fmt.Sprintf("创建CycleTimeConf失败，\r\n失败信息 %s", err.Error()))
        return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改CycleTimeConf
// @Summary 修改CycleTimeConf
// @Description 修改CycleTimeConf
// @Tags CycleTimeConf
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.CycleTimeConfUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/cycle-time-conf/{id} [put]
// @Security Bearer
func (e CycleTimeConf) Update(c *gin.Context) {
    req := dto.CycleTimeConfUpdateReq{}
    s := service.CycleTimeConf{}
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
		e.Error(500, err, fmt.Sprintf("修改CycleTimeConf失败，\r\n失败信息 %s", err.Error()))
        return
	}
	e.OK( req.GetId(), "修改成功")
}

// Delete 删除CycleTimeConf
// @Summary 删除CycleTimeConf
// @Description 删除CycleTimeConf
// @Tags CycleTimeConf
// @Param data body dto.CycleTimeConfDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/cycle-time-conf [delete]
// @Security Bearer
func (e CycleTimeConf) Delete(c *gin.Context) {
    s := service.CycleTimeConf{}
    req := dto.CycleTimeConfDeleteReq{}
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
		e.Error(500, err, fmt.Sprintf("删除CycleTimeConf失败，\r\n失败信息 %s", err.Error()))
        return
	}
	e.OK( req.GetId(), "删除成功")
}
