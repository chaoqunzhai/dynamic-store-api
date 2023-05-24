package apis

import (
    "fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"

	"go-admin/app/shop/models"
	"go-admin/app/shop/service"
	"go-admin/app/shop/service/dto"
	"go-admin/common/actions"
)

type ShopOrderBindRecord struct {
	api.Api
}

// GetPage 获取ShopOrderBindRecord列表
// @Summary 获取ShopOrderBindRecord列表
// @Description 获取ShopOrderBindRecord列表
// @Tags ShopOrderBindRecord
// @Param shopId query string false "关联的小B客户"
// @Param recordId query string false "每次记录的总ID"
// @Param orderId query string false "订单ID"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.ShopOrderBindRecord}} "{"code": 200, "data": [...]}"
// @Router /api/v1/shop-order-bind-record [get]
// @Security Bearer
func (e ShopOrderBindRecord) GetPage(c *gin.Context) {
    req := dto.ShopOrderBindRecordGetPageReq{}
    s := service.ShopOrderBindRecord{}
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
	list := make([]models.ShopOrderBindRecord, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取ShopOrderBindRecord失败，\r\n失败信息 %s", err.Error()))
        return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取ShopOrderBindRecord
// @Summary 获取ShopOrderBindRecord
// @Description 获取ShopOrderBindRecord
// @Tags ShopOrderBindRecord
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.ShopOrderBindRecord} "{"code": 200, "data": [...]}"
// @Router /api/v1/shop-order-bind-record/{id} [get]
// @Security Bearer
func (e ShopOrderBindRecord) Get(c *gin.Context) {
	req := dto.ShopOrderBindRecordGetReq{}
	s := service.ShopOrderBindRecord{}
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
	var object models.ShopOrderBindRecord

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取ShopOrderBindRecord失败，\r\n失败信息 %s", err.Error()))
        return
	}

	e.OK( object, "查询成功")
}

// Insert 创建ShopOrderBindRecord
// @Summary 创建ShopOrderBindRecord
// @Description 创建ShopOrderBindRecord
// @Tags ShopOrderBindRecord
// @Accept application/json
// @Product application/json
// @Param data body dto.ShopOrderBindRecordInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/shop-order-bind-record [post]
// @Security Bearer
func (e ShopOrderBindRecord) Insert(c *gin.Context) {
    req := dto.ShopOrderBindRecordInsertReq{}
    s := service.ShopOrderBindRecord{}
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
		e.Error(500, err, fmt.Sprintf("创建ShopOrderBindRecord失败，\r\n失败信息 %s", err.Error()))
        return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改ShopOrderBindRecord
// @Summary 修改ShopOrderBindRecord
// @Description 修改ShopOrderBindRecord
// @Tags ShopOrderBindRecord
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.ShopOrderBindRecordUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/shop-order-bind-record/{id} [put]
// @Security Bearer
func (e ShopOrderBindRecord) Update(c *gin.Context) {
    req := dto.ShopOrderBindRecordUpdateReq{}
    s := service.ShopOrderBindRecord{}
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
		e.Error(500, err, fmt.Sprintf("修改ShopOrderBindRecord失败，\r\n失败信息 %s", err.Error()))
        return
	}
	e.OK( req.GetId(), "修改成功")
}

// Delete 删除ShopOrderBindRecord
// @Summary 删除ShopOrderBindRecord
// @Description 删除ShopOrderBindRecord
// @Tags ShopOrderBindRecord
// @Param data body dto.ShopOrderBindRecordDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/shop-order-bind-record [delete]
// @Security Bearer
func (e ShopOrderBindRecord) Delete(c *gin.Context) {
    s := service.ShopOrderBindRecord{}
    req := dto.ShopOrderBindRecordDeleteReq{}
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
		e.Error(500, err, fmt.Sprintf("删除ShopOrderBindRecord失败，\r\n失败信息 %s", err.Error()))
        return
	}
	e.OK( req.GetId(), "删除成功")
}
