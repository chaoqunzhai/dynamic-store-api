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

type GoodsSales struct {
	api.Api
}

// GetPage 获取GoodsSales列表
// @Summary 获取GoodsSales列表
// @Description 获取GoodsSales列表
// @Tags GoodsSales
// @Param layer query string false "排序"
// @Param enable query string false "开关"
// @Param cId query string false "大BID"
// @Param productId query string false "产品ID"
// @Param productName query string false "产品名称"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.GoodsSales}} "{"code": 200, "data": [...]}"
// @Router /api/v1/goods-sales [get]
// @Security Bearer
func (e GoodsSales) GetPage(c *gin.Context) {
	req := dto.GoodsSalesGetPageReq{}
	s := service.GoodsSales{}
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
	list := make([]models.GoodsSales, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取GoodsSales失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取GoodsSales
// @Summary 获取GoodsSales
// @Description 获取GoodsSales
// @Tags GoodsSales
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.GoodsSales} "{"code": 200, "data": [...]}"
// @Router /api/v1/goods-sales/{id} [get]
// @Security Bearer
func (e GoodsSales) Get(c *gin.Context) {
	req := dto.GoodsSalesGetReq{}
	s := service.GoodsSales{}
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
	var object models.GoodsSales

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取GoodsSales失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}

// Insert 创建GoodsSales
// @Summary 创建GoodsSales
// @Description 创建GoodsSales
// @Tags GoodsSales
// @Accept application/json
// @Product application/json
// @Param data body dto.GoodsSalesInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/goods-sales [post]
// @Security Bearer
func (e GoodsSales) Insert(c *gin.Context) {
	req := dto.GoodsSalesInsertReq{}
	s := service.GoodsSales{}
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
		e.Error(500, err, fmt.Sprintf("创建GoodsSales失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改GoodsSales
// @Summary 修改GoodsSales
// @Description 修改GoodsSales
// @Tags GoodsSales
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.GoodsSalesUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/goods-sales/{id} [put]
// @Security Bearer
func (e GoodsSales) Update(c *gin.Context) {
	req := dto.GoodsSalesUpdateReq{}
	s := service.GoodsSales{}
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
		e.Error(500, err, fmt.Sprintf("修改GoodsSales失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "修改成功")
}

// Delete 删除GoodsSales
// @Summary 删除GoodsSales
// @Description 删除GoodsSales
// @Tags GoodsSales
// @Param data body dto.GoodsSalesDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/goods-sales [delete]
// @Security Bearer
func (e GoodsSales) Delete(c *gin.Context) {
	s := service.GoodsSales{}
	req := dto.GoodsSalesDeleteReq{}
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
		e.Error(500, err, fmt.Sprintf("删除GoodsSales失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "删除成功")
}
