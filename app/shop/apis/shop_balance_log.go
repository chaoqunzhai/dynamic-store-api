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

type ShopBalanceLog struct {
	api.Api
}

// GetPage 获取ShopBalanceLog列表
// @Summary 获取ShopBalanceLog列表
// @Description 获取ShopBalanceLog列表
// @Tags ShopBalanceLog
// @Param shopId query string false "小BID"
// @Param money query string false "变动金额"
// @Param scene query string false "变动场景"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.ShopBalanceLog}} "{"code": 200, "data": [...]}"
// @Router /api/v1/shop-balance-log [get]
// @Security Bearer
func (e ShopBalanceLog) GetPage(c *gin.Context) {
    req := dto.ShopBalanceLogGetPageReq{}
    s := service.ShopBalanceLog{}
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

	p := actions.GetPermissionFromContext(c)
	list := make([]models.ShopBalanceLog, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取ShopBalanceLog失败，\r\n失败信息 %s", err.Error()))
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

	var shopListObject []models.Shop
	var shopObject models.Shop
	e.Orm.Model(&models.Shop{}).Scopes(actions.PermissionSysUser(shopObject.TableName(),userDto)).Select("name,id").Where("id in ? ",shopList).Find(&shopListObject)
	shopRowMap:=make(map[int]string,0)
	for _,row:=range shopListObject{
		shopRowMap[row.Id] = row.Name
	}

	result :=make([]interface{},0)
	//聚合下shop,防止多次查询
	for _,row:=range list{

		r:=map[string]interface{}{
			"id":row.Id,
			"money":row.Money,
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

// Get 获取ShopBalanceLog
// @Summary 获取ShopBalanceLog
// @Description 获取ShopBalanceLog
// @Tags ShopBalanceLog
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.ShopBalanceLog} "{"code": 200, "data": [...]}"
// @Router /api/v1/shop-balance-log/{id} [get]
// @Security Bearer
func (e ShopBalanceLog) Get(c *gin.Context) {
	req := dto.ShopBalanceLogGetReq{}
	s := service.ShopBalanceLog{}
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
	var object models.ShopBalanceLog

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取ShopBalanceLog失败，\r\n失败信息 %s", err.Error()))
        return
	}

	e.OK( object, "查询成功")
}

// Insert 创建ShopBalanceLog
// @Summary 创建ShopBalanceLog
// @Description 创建ShopBalanceLog
// @Tags ShopBalanceLog
// @Accept application/json
// @Product application/json
// @Param data body dto.ShopBalanceLogInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/shop-balance-log [post]
// @Security Bearer
func (e ShopBalanceLog) Insert(c *gin.Context) {
    req := dto.ShopBalanceLogInsertReq{}
    s := service.ShopBalanceLog{}
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
		e.Error(500, err, fmt.Sprintf("创建ShopBalanceLog失败，\r\n失败信息 %s", err.Error()))
        return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改ShopBalanceLog
// @Summary 修改ShopBalanceLog
// @Description 修改ShopBalanceLog
// @Tags ShopBalanceLog
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.ShopBalanceLogUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/shop-balance-log/{id} [put]
// @Security Bearer
func (e ShopBalanceLog) Update(c *gin.Context) {
    req := dto.ShopBalanceLogUpdateReq{}
    s := service.ShopBalanceLog{}
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
		e.Error(500, err, fmt.Sprintf("修改ShopBalanceLog失败，\r\n失败信息 %s", err.Error()))
        return
	}
	e.OK( req.GetId(), "修改成功")
}

// Delete 删除ShopBalanceLog
// @Summary 删除ShopBalanceLog
// @Description 删除ShopBalanceLog
// @Tags ShopBalanceLog
// @Param data body dto.ShopBalanceLogDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/shop-balance-log [delete]
// @Security Bearer
func (e ShopBalanceLog) Delete(c *gin.Context) {
    s := service.ShopBalanceLog{}
    req := dto.ShopBalanceLogDeleteReq{}
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
		e.Error(500, err, fmt.Sprintf("删除ShopBalanceLog失败，\r\n失败信息 %s", err.Error()))
        return
	}
	e.OK( req.GetId(), "删除成功")
}
