package apis

import (
    "fmt"
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

type GoodsVip struct {
	api.Api
}

// GetPage 获取GoodsVip列表
// @Summary 获取GoodsVip列表
// @Description 获取GoodsVip列表
// @Tags GoodsVip
// @Param layer query string false "排序"
// @Param enable query string false "开关"
// @Param cId query string false "大BID"
// @Param goodsId query string false "商品ID"
// @Param specsId query string false "规格ID"
// @Param gradeId query string false "VipId"
// @Param customPrice query string false "自定义价格"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.GoodsVip}} "{"code": 200, "data": [...]}"
// @Router /api/v1/goods-vip [get]
// @Security Bearer
func (e GoodsVip) GetPage(c *gin.Context) {
    req := dto.GoodsVipGetPageReq{}
    s := service.GoodsVip{}
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
	list := make([]models.GoodsVip, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取GoodsVip失败，\r\n失败信息 %s", err.Error()))
        return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取GoodsVip
// @Summary 获取GoodsVip
// @Description 获取GoodsVip
// @Tags GoodsVip
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.GoodsVip} "{"code": 200, "data": [...]}"
// @Router /api/v1/goods-vip/{id} [get]
// @Security Bearer
func (e GoodsVip) Get(c *gin.Context) {
	req := dto.GoodsVipGetReq{}
	s := service.GoodsVip{}
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
	var object models.GoodsVip

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取GoodsVip失败，\r\n失败信息 %s", err.Error()))
        return
	}

	e.OK( object, "查询成功")
}

// Insert 创建GoodsVip
// @Summary 创建GoodsVip
// @Description 创建GoodsVip
// @Tags GoodsVip
// @Accept application/json
// @Product application/json
// @Param data body dto.GoodsVipInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/goods-vip [post]
// @Security Bearer
func (e GoodsVip) Insert(c *gin.Context) {
    req := dto.GoodsVipInsertReq{}
    s := service.GoodsVip{}
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
		e.Error(500, err, fmt.Sprintf("创建GoodsVip失败，\r\n失败信息 %s", err.Error()))
        return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改GoodsVip
// @Summary 修改GoodsVip
// @Description 修改GoodsVip
// @Tags GoodsVip
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.GoodsVipUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/goods-vip/{id} [put]
// @Security Bearer
func (e GoodsVip) Update(c *gin.Context) {
    req := dto.GoodsVipUpdateReq{}
    s := service.GoodsVip{}
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
		e.Error(500, err, fmt.Sprintf("修改GoodsVip失败，\r\n失败信息 %s", err.Error()))
        return
	}
	e.OK( req.GetId(), "修改成功")
}

// Delete 删除GoodsVip
// @Summary 删除GoodsVip
// @Description 删除GoodsVip
// @Tags GoodsVip
// @Param data body dto.GoodsVipDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/goods-vip [delete]
// @Security Bearer
func (e GoodsVip) Delete(c *gin.Context) {
    s := service.GoodsVip{}
    req := dto.GoodsVipDeleteReq{}
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

	userDto, err := customUser.GetUserDto(e.Orm, c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}
	err = s.Remove(&req, userDto.CId)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("删除GoodsVip失败，\r\n失败信息 %s", err.Error()))
        return
	}
	e.OK( req.GetId(), "删除成功")
}
