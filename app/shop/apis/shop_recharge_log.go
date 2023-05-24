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

type ShopRechargeLog struct {
	api.Api
}

// GetPage 获取ShopRechargeLog列表
// @Summary 获取ShopRechargeLog列表
// @Description 获取ShopRechargeLog列表
// @Tags ShopRechargeLog
// @Param shopId query string false "小BID"
// @Param uuid query string false "订单号"
// @Param source query string false "充值方式"
// @Param money query string false "支付金额"
// @Param payStatus query string false "支付状态"
// @Param payTime query time.Time false "支付时间"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.ShopRechargeLog}} "{"code": 200, "data": [...]}"
// @Router /api/v1/shop-recharge-log [get]
// @Security Bearer
func (e ShopRechargeLog) GetPage(c *gin.Context) {
    req := dto.ShopRechargeLogGetPageReq{}
    s := service.ShopRechargeLog{}
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
	list := make([]models.ShopRechargeLog, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取ShopRechargeLog失败，\r\n失败信息 %s", err.Error()))
        return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取ShopRechargeLog
// @Summary 获取ShopRechargeLog
// @Description 获取ShopRechargeLog
// @Tags ShopRechargeLog
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.ShopRechargeLog} "{"code": 200, "data": [...]}"
// @Router /api/v1/shop-recharge-log/{id} [get]
// @Security Bearer
func (e ShopRechargeLog) Get(c *gin.Context) {
	req := dto.ShopRechargeLogGetReq{}
	s := service.ShopRechargeLog{}
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
	var object models.ShopRechargeLog

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取ShopRechargeLog失败，\r\n失败信息 %s", err.Error()))
        return
	}

	e.OK( object, "查询成功")
}

// Insert 创建ShopRechargeLog
// @Summary 创建ShopRechargeLog
// @Description 创建ShopRechargeLog
// @Tags ShopRechargeLog
// @Accept application/json
// @Product application/json
// @Param data body dto.ShopRechargeLogInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/shop-recharge-log [post]
// @Security Bearer
func (e ShopRechargeLog) Insert(c *gin.Context) {
    req := dto.ShopRechargeLogInsertReq{}
    s := service.ShopRechargeLog{}
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
		e.Error(500, err, fmt.Sprintf("创建ShopRechargeLog失败，\r\n失败信息 %s", err.Error()))
        return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改ShopRechargeLog
// @Summary 修改ShopRechargeLog
// @Description 修改ShopRechargeLog
// @Tags ShopRechargeLog
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.ShopRechargeLogUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/shop-recharge-log/{id} [put]
// @Security Bearer
func (e ShopRechargeLog) Update(c *gin.Context) {
    req := dto.ShopRechargeLogUpdateReq{}
    s := service.ShopRechargeLog{}
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
		e.Error(500, err, fmt.Sprintf("修改ShopRechargeLog失败，\r\n失败信息 %s", err.Error()))
        return
	}
	e.OK( req.GetId(), "修改成功")
}

// Delete 删除ShopRechargeLog
// @Summary 删除ShopRechargeLog
// @Description 删除ShopRechargeLog
// @Tags ShopRechargeLog
// @Param data body dto.ShopRechargeLogDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/shop-recharge-log [delete]
// @Security Bearer
func (e ShopRechargeLog) Delete(c *gin.Context) {
    s := service.ShopRechargeLog{}
    req := dto.ShopRechargeLogDeleteReq{}
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
		e.Error(500, err, fmt.Sprintf("删除ShopRechargeLog失败，\r\n失败信息 %s", err.Error()))
        return
	}
	e.OK( req.GetId(), "删除成功")
}
