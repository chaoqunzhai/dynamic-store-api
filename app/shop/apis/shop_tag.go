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

type ShopTag struct {
	api.Api
}

// GetPage 获取ShopTag列表
// @Summary 获取ShopTag列表
// @Description 获取ShopTag列表
// @Tags ShopTag
// @Param layer query string false "排序"
// @Param enable query string false "开关"
// @Param desc query string false "描述信息"
// @Param cId query string false "大BID"
// @Param name query string false "客户标签名称"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.ShopTag}} "{"code": 200, "data": [...]}"
// @Router /api/v1/shop-tag [get]
// @Security Bearer
func (e ShopTag) GetPage(c *gin.Context) {
    req := dto.ShopTagGetPageReq{}
    s := service.ShopTag{}
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
	list := make([]models.ShopTag, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取ShopTag失败，\r\n失败信息 %s", err.Error()))
        return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取ShopTag
// @Summary 获取ShopTag
// @Description 获取ShopTag
// @Tags ShopTag
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.ShopTag} "{"code": 200, "data": [...]}"
// @Router /api/v1/shop-tag/{id} [get]
// @Security Bearer
func (e ShopTag) Get(c *gin.Context) {
	req := dto.ShopTagGetReq{}
	s := service.ShopTag{}
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
	var object models.ShopTag

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取ShopTag失败，\r\n失败信息 %s", err.Error()))
        return
	}

	e.OK( object, "查询成功")
}

// Insert 创建ShopTag
// @Summary 创建ShopTag
// @Description 创建ShopTag
// @Tags ShopTag
// @Accept application/json
// @Product application/json
// @Param data body dto.ShopTagInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/shop-tag [post]
// @Security Bearer
func (e ShopTag) Insert(c *gin.Context) {
    req := dto.ShopTagInsertReq{}
    s := service.ShopTag{}
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
		e.Error(500, err, fmt.Sprintf("创建ShopTag失败，\r\n失败信息 %s", err.Error()))
        return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改ShopTag
// @Summary 修改ShopTag
// @Description 修改ShopTag
// @Tags ShopTag
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.ShopTagUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/shop-tag/{id} [put]
// @Security Bearer
func (e ShopTag) Update(c *gin.Context) {
    req := dto.ShopTagUpdateReq{}
    s := service.ShopTag{}
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
		e.Error(500, err, fmt.Sprintf("修改ShopTag失败，\r\n失败信息 %s", err.Error()))
        return
	}
	e.OK( req.GetId(), "修改成功")
}

// Delete 删除ShopTag
// @Summary 删除ShopTag
// @Description 删除ShopTag
// @Tags ShopTag
// @Param data body dto.ShopTagDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/shop-tag [delete]
// @Security Bearer
func (e ShopTag) Delete(c *gin.Context) {
    s := service.ShopTag{}
    req := dto.ShopTagDeleteReq{}
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
		e.Error(500, err, fmt.Sprintf("删除ShopTag失败，\r\n失败信息 %s", err.Error()))
        return
	}
	e.OK( req.GetId(), "删除成功")
}