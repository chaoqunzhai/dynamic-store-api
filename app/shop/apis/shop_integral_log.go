package apis

import (
    "fmt"
	customUser "go-admin/common/jwt/user"
	"go-admin/global"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"

	"go-admin/app/shop/models"
	"go-admin/app/shop/service"
	"go-admin/app/shop/service/dto"
	"go-admin/common/actions"
)

type ShopIntegralLog struct {
	api.Api
}

// GetPage 获取ShopIntegralLog列表
// @Summary 获取ShopIntegralLog列表
// @Description 获取ShopIntegralLog列表
// @Tags ShopIntegralLog
// @Param shopId query string false "小BID"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.ShopIntegralLog}} "{"code": 200, "data": [...]}"
// @Router /api/v1/shop-integral-log [get]
// @Security Bearer
func (e ShopIntegralLog) GetPage(c *gin.Context) {
    req := dto.ShopIntegralLogGetPageReq{}
    s := service.ShopIntegralLog{}
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
	list := make([]models.ShopIntegralLog, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取ShopIntegralLog失败，\r\n失败信息 %s", err.Error()))
        return
	}
	userDto, err := customUser.GetUserDto(e.Orm, c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}
	shopMapCache:=make(map[int]bool,0)
	shopList:=make([]int,0)
	for _,row:=range list{
		ok:=shopMapCache[row.ShopId]
		if !ok {
			shopList = append(shopList,row.ShopId)
			shopMapCache[row.ShopId] = true
		}
	}

	var shopObject []models.Shop
	var object models.Shop
	e.Orm.Model(&models.Shop{}).Scopes(actions.PermissionSysUser(object.TableName(), userDto)).Select("name,id").Where("id in ? ",shopList).Find(&shopObject)
	shopRowMap:=make(map[int]string,0)
	for _,row:=range shopObject{
		shopRowMap[row.Id] = row.Name
	}

	result :=make([]interface{},0)
	//聚合下shop,防止多次查询
	for _,row:=range list{

		r:=map[string]interface{}{
			"id":row.Id,
			"number":row.Number,
			"desc":row.Desc,
			"scene":row.Scene,
			"created_at":row.CreatedAt,
			"type":global.GetScanStr(row.Type),
		}
		if row.ShopId > 0 {
			if shopName,ok:=shopRowMap[row.ShopId];ok{
				r["shop_name"] = shopName
			}
		}
		result = append(result,r)
	}
	e.PageOK(result, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取ShopIntegralLog
// @Summary 获取ShopIntegralLog
// @Description 获取ShopIntegralLog
// @Tags ShopIntegralLog
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.ShopIntegralLog} "{"code": 200, "data": [...]}"
// @Router /api/v1/shop-integral-log/{id} [get]
// @Security Bearer
func (e ShopIntegralLog) Get(c *gin.Context) {
	req := dto.ShopIntegralLogGetReq{}
	s := service.ShopIntegralLog{}
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
	var object models.ShopIntegralLog

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取ShopIntegralLog失败，\r\n失败信息 %s", err.Error()))
        return
	}

	e.OK( object, "查询成功")
}

// Insert 创建ShopIntegralLog
// @Summary 创建ShopIntegralLog
// @Description 创建ShopIntegralLog
// @Tags ShopIntegralLog
// @Accept application/json
// @Product application/json
// @Param data body dto.ShopIntegralLogInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/shop-integral-log [post]
// @Security Bearer
func (e ShopIntegralLog) Insert(c *gin.Context) {
    req := dto.ShopIntegralLogInsertReq{}
    s := service.ShopIntegralLog{}
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
		e.Error(500, err, fmt.Sprintf("创建ShopIntegralLog失败，\r\n失败信息 %s", err.Error()))
        return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改ShopIntegralLog
// @Summary 修改ShopIntegralLog
// @Description 修改ShopIntegralLog
// @Tags ShopIntegralLog
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.ShopIntegralLogUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/shop-integral-log/{id} [put]
// @Security Bearer
func (e ShopIntegralLog) Update(c *gin.Context) {
    req := dto.ShopIntegralLogUpdateReq{}
    s := service.ShopIntegralLog{}
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
		e.Error(500, err, fmt.Sprintf("修改ShopIntegralLog失败，\r\n失败信息 %s", err.Error()))
        return
	}
	e.OK( req.GetId(), "修改成功")
}

// Delete 删除ShopIntegralLog
// @Summary 删除ShopIntegralLog
// @Description 删除ShopIntegralLog
// @Tags ShopIntegralLog
// @Param data body dto.ShopIntegralLogDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/shop-integral-log [delete]
// @Security Bearer
func (e ShopIntegralLog) Delete(c *gin.Context) {
    s := service.ShopIntegralLog{}
    req := dto.ShopIntegralLogDeleteReq{}
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
		e.Error(500, err, fmt.Sprintf("删除ShopIntegralLog失败，\r\n失败信息 %s", err.Error()))
        return
	}
	e.OK( req.GetId(), "删除成功")
}
