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

type ShopOrderRecord struct {
	api.Api
}

// GetPage 获取ShopOrderRecord列表
// @Summary 获取ShopOrderRecord列表
// @Description 获取ShopOrderRecord列表
// @Tags ShopOrderRecord
// @Param shopId query string false "关联的小B客户"
// @Param shopName query string false "客户名称"
// @Param number query string false "订单量"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.ShopOrderRecord}} "{"code": 200, "data": [...]}"
// @Router /api/v1/shop-order-record [get]
// @Security Bearer
func (e ShopOrderRecord) GetPage(c *gin.Context) {
    req := dto.ShopOrderRecordGetPageReq{}
    s := service.ShopOrderRecord{}
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
	list := make([]models.ShopOrderRecord, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取ShopOrderRecord失败，\r\n失败信息 %s", err.Error()))
        return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取ShopOrderRecord
// @Summary 获取ShopOrderRecord
// @Description 获取ShopOrderRecord
// @Tags ShopOrderRecord
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.ShopOrderRecord} "{"code": 200, "data": [...]}"
// @Router /api/v1/shop-order-record/{id} [get]
// @Security Bearer
func (e ShopOrderRecord) Get(c *gin.Context) {
	req := dto.ShopOrderRecordGetReq{}
	s := service.ShopOrderRecord{}
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
	var object models.ShopOrderRecord

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取ShopOrderRecord失败，\r\n失败信息 %s", err.Error()))
        return
	}

	e.OK( object, "查询成功")
}

// Insert 创建ShopOrderRecord
// @Summary 创建ShopOrderRecord
// @Description 创建ShopOrderRecord
// @Tags ShopOrderRecord
// @Accept application/json
// @Product application/json
// @Param data body dto.ShopOrderRecordInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/shop-order-record [post]
// @Security Bearer
func (e ShopOrderRecord) Insert(c *gin.Context) {
    req := dto.ShopOrderRecordInsertReq{}
    s := service.ShopOrderRecord{}
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
		e.Error(500, err, fmt.Sprintf("创建ShopOrderRecord失败，\r\n失败信息 %s", err.Error()))
        return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改ShopOrderRecord
// @Summary 修改ShopOrderRecord
// @Description 修改ShopOrderRecord
// @Tags ShopOrderRecord
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.ShopOrderRecordUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/shop-order-record/{id} [put]
// @Security Bearer
func (e ShopOrderRecord) Update(c *gin.Context) {
    req := dto.ShopOrderRecordUpdateReq{}
    s := service.ShopOrderRecord{}
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
		e.Error(500, err, fmt.Sprintf("修改ShopOrderRecord失败，\r\n失败信息 %s", err.Error()))
        return
	}
	e.OK( req.GetId(), "修改成功")
}

// Delete 删除ShopOrderRecord
// @Summary 删除ShopOrderRecord
// @Description 删除ShopOrderRecord
// @Tags ShopOrderRecord
// @Param data body dto.ShopOrderRecordDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/shop-order-record [delete]
// @Security Bearer
func (e ShopOrderRecord) Delete(c *gin.Context) {
    s := service.ShopOrderRecord{}
    req := dto.ShopOrderRecordDeleteReq{}
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
		e.Error(500, err, fmt.Sprintf("删除ShopOrderRecord失败，\r\n失败信息 %s", err.Error()))
        return
	}
	e.OK( req.GetId(), "删除成功")
}
