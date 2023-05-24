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

type GoodsTag struct {
	api.Api
}

// GetPage 获取GoodsTag列表
// @Summary 获取GoodsTag列表
// @Description 获取GoodsTag列表
// @Tags GoodsTag
// @Param layer query string false "排序"
// @Param enable query string false "开关"
// @Param cId query string false "大BID"
// @Param name query string false "商品标签名称"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.GoodsTag}} "{"code": 200, "data": [...]}"
// @Router /api/v1/goods-tag [get]
// @Security Bearer
func (e GoodsTag) GetPage(c *gin.Context) {
    req := dto.GoodsTagGetPageReq{}
    s := service.GoodsTag{}
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
	list := make([]models.GoodsTag, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取GoodsTag失败，\r\n失败信息 %s", err.Error()))
        return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取GoodsTag
// @Summary 获取GoodsTag
// @Description 获取GoodsTag
// @Tags GoodsTag
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.GoodsTag} "{"code": 200, "data": [...]}"
// @Router /api/v1/goods-tag/{id} [get]
// @Security Bearer
func (e GoodsTag) Get(c *gin.Context) {
	req := dto.GoodsTagGetReq{}
	s := service.GoodsTag{}
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
	var object models.GoodsTag

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取GoodsTag失败，\r\n失败信息 %s", err.Error()))
        return
	}

	e.OK( object, "查询成功")
}

// Insert 创建GoodsTag
// @Summary 创建GoodsTag
// @Description 创建GoodsTag
// @Tags GoodsTag
// @Accept application/json
// @Product application/json
// @Param data body dto.GoodsTagInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/goods-tag [post]
// @Security Bearer
func (e GoodsTag) Insert(c *gin.Context) {
    req := dto.GoodsTagInsertReq{}
    s := service.GoodsTag{}
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
		e.Error(500, err, fmt.Sprintf("创建GoodsTag失败，\r\n失败信息 %s", err.Error()))
        return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改GoodsTag
// @Summary 修改GoodsTag
// @Description 修改GoodsTag
// @Tags GoodsTag
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.GoodsTagUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/goods-tag/{id} [put]
// @Security Bearer
func (e GoodsTag) Update(c *gin.Context) {
    req := dto.GoodsTagUpdateReq{}
    s := service.GoodsTag{}
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
		e.Error(500, err, fmt.Sprintf("修改GoodsTag失败，\r\n失败信息 %s", err.Error()))
        return
	}
	e.OK( req.GetId(), "修改成功")
}

// Delete 删除GoodsTag
// @Summary 删除GoodsTag
// @Description 删除GoodsTag
// @Tags GoodsTag
// @Param data body dto.GoodsTagDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/goods-tag [delete]
// @Security Bearer
func (e GoodsTag) Delete(c *gin.Context) {
    s := service.GoodsTag{}
    req := dto.GoodsTagDeleteReq{}
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
		e.Error(500, err, fmt.Sprintf("删除GoodsTag失败，\r\n失败信息 %s", err.Error()))
        return
	}
	e.OK( req.GetId(), "删除成功")
}
