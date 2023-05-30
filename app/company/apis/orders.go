package apis

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"
	models2 "go-admin/cmd/migrate/migration/models"
	customUser "go-admin/common/jwt/user"
	"go-admin/global"
	"gorm.io/gorm"
	"strings"

	"go-admin/app/company/models"
	"go-admin/app/company/service"
	"go-admin/app/company/service/dto"
	"go-admin/common/actions"
)

type Orders struct {
	api.Api
}

func (e Orders) getTableName(cid int) string {
	//先在split分表中查询

	var splitRow models2.SplitTableMap
	e.Orm.Model(&models2.SplitTableMap{}).Where("c_id = ? and enable = ? and type = ?", cid, true, global.SplitOrder).Limit(1).Find(&splitRow)

	tableName := ""
	if splitRow.Id > 0 {
		tableName = splitRow.Name
	} else {
		tableName = global.SplitOrderDefaultTableName
	}

	return tableName
}

// GetPage 获取Orders列表
// @Summary 获取Orders列表
// @Description 获取Orders列表
// @Tags Orders
// @Param layer query string false "排序"
// @Param enable query string false "开关"
// @Param cId query string false "大BID"
// @Param shopId query string false "关联客户"
// @Param status query string false "配送状态"
// @Param number query string false "下单数量"
// @Param delivery query string false "配送周期"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.Orders}} "{"code": 200, "data": [...]}"
// @Router /api/v1/orders [get]
// @Security Bearer
func (e Orders) GetPage(c *gin.Context) {
	req := dto.OrdersGetPageReq{}
	s := service.Orders{}
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
	list := make([]models.Orders, 0)
	var count int64
	req.CId = userDto.CId
	err = s.GetPage(e.getTableName(userDto.CId), &req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取订单失败,%s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取Orders
// @Summary 获取Orders
// @Description 获取Orders
// @Tags Orders
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.Orders} "{"code": 200, "data": [...]}"
// @Router /api/v1/orders/{id} [get]
// @Security Bearer
func (e Orders) Get(c *gin.Context) {
	req := dto.OrdersGetReq{}
	s := service.Orders{}
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

	var object models.Orders

	orderTableName :=e.getTableName(userDto.CId)
	orderErr := e.Orm.Table(orderTableName).First(&object,req.Id).Error
	if orderErr != nil && errors.Is(orderErr, gorm.ErrRecordNotFound) {

		e.Error(500, orderErr, "订单不存在")
		return
	}
	if orderErr != nil {
		e.Error(500, orderErr, "订单不存在")
		return
	}

	var shopRow models2.Shop
	e.Orm.Model(&models2.Shop{}).Where("id = ? and c_id = ?",object.ShopId,userDto.CId).Limit(1).Find(&shopRow)

	var timeCnf models.CycleTimeConf
	e.Orm.Model(&models.CycleTimeConf{}).Where("id = ? and c_id = ?",object.Delivery,userDto.CId).Limit(1).Find(&timeCnf)
	result :=map[string]interface{}{
		"order_id":object.Id,
		"created_at":object.CreatedAt,
		"delivery":timeCnf.Id,
		"delivery_give":timeCnf.GiveTime,
		"pay":global.GetPayStr(object.Pay),
		"shop_name":shopRow.Name,
		"shop_username":shopRow.UserName,
		"shop_phone":shopRow.Phone,
		"shop_address":shopRow.Address,
	}
	var orderSpecs []models.OrderSpecs
	//是一个分表的名称
	specsTable:=e.OrderSpecsTableName(orderTableName)

	e.Orm.Table(specsTable).Where("order_id = ?",object.Id).Find(&orderSpecs)

	specsList :=make([]map[string]interface{},0)
	for _,row:=range orderSpecs{
		var specRow models.GoodsSpecs
		e.Orm.Model(&models.GoodsSpecs{}).Where("id = ? and c_id = ?",row.SpecsId,userDto.CId).Limit(1).Find(&specRow)
		ss :=map[string]interface{}{
			"name":specRow.Name,
			"spec":fmt.Sprintf("%v%v",row.Number,specRow.Unit),
			"status":row.Status,
			"money":row.Money,
		}
		specsList  = append(specsList,ss)
	}
	result["specs_list"]  = specsList
	e.OK(result, "查询成功")
}


func (e Orders) ValetOrder(c *gin.Context) {
	err := e.MakeContext(c).
		MakeOrm().
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}



}
func (e Orders) Times(c *gin.Context) {
	s := service.Orders{}
	err := e.MakeContext(c).
		MakeOrm().
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
	var lists []models.CycleTimeConf
	e.Orm.Model(&models.CycleTimeConf{}).Where("c_id = ? and enable = ?",userDto.CId,true).Order(global.OrderLayerKey).Find(&lists)

	e.PageOK(lists,len(lists),1,-1,"successful")
	return
}
func (e Orders) ValidTimeConf(c *gin.Context) {
	s := service.Orders{}
	err := e.MakeContext(c).
		MakeOrm().
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
	timeConf,_ :=s.ValidTimeConf(userDto.CId)
	if !timeConf{
		e.Error(500, errors.New("非下单时间段"), "非下单时间段")
		return
	}
	e.OK("当前时间可下单","successful")
	return
}

func (e Orders)OrderSpecsTableName(orderTable string)  string {
	//子表默认名称
	specsTable:=global.SplitOrderDefaultSubTableName
	//判断是否分表了
	//默认是 orders 表名，如果分表后就是 orders_大BID_时间戳后6位

	if orderTable != global.SplitOrderDefaultTableName {
		//拼接位 order_specs_大BID_时间戳后6位
		specsTable =fmt.Sprintf("%v%v",specsTable,strings.Replace(orderTable,global.SplitOrderDefaultTableName,"",-1))
	}
	return specsTable
}
// Insert 创建Orders
// @Summary 创建Orders
// @Description 创建Orders
// @Tags Orders
// @Accept application/json
// @Product application/json
// @Param data body dto.OrdersInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/orders [post]
// @Security Bearer
func (e Orders) Insert(c *gin.Context) {
	req := dto.OrdersInsertReq{}
	s := service.Orders{}
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
	userDto, err := customUser.GetUserDto(e.Orm, c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}
	//根据下单的时间区间，来匹配
	//todo:配送周期
	//根据下单的时间区间来自动匹配,
	//1.查询这个时间段内是否配置了cycle_time_conf得值,如果创建了进行关联即可
	//2.如果没有这个时间
	timeConf,DeliveryId :=s.ValidTimeConf(userDto.CId)
	if !timeConf{
		e.Error(500, errors.New("非下单时间段"), "非下单时间段")
		return
	}
	userId := user.GetUserId(c)
	//todo:获取表名
	orderTableName:=e.getTableName(userDto.CId)

	var data models.Orders

	data.CId = userDto.CId
	data.Enable = true
	data.Layer = req.Layer
	data.Status = global.OrderStatusWait
	data.Desc = req.Desc
	data.ShopId = req.ShopId

	//todo:配送周期
	data.Delivery = DeliveryId

	data.CreateBy = userId
	createErr := e.Orm.Table(orderTableName).Create(&data).Error
	if createErr != nil {
		e.Error(500, createErr, "订单创建失败")
		return
	}
	//分表检测
	specsTable:=e.OrderSpecsTableName(orderTableName)

	var orderMoney float64
	var goodsNumber int
	for _,good:=range req.Goods{
		var count int64
		e.Orm.Model(&models.GoodsSpecs{}).Where("id = ?",good.SpecsId).Count(&count)
		if count == 0 {continue}
		orderMoney+=good.Money
		goodsNumber++

		e.Orm.Table(specsTable).Create(&models.OrderSpecs{
			OrderId: data.Id,
			SpecsId: good.SpecsId,
			Status: global.OrderStatusWait,
			Money: good.Money,
		})
	}

	e.Orm.Model(&models.Orders{}).Table(orderTableName).Where("id = ?",data.Id).Updates(map[string]interface{}{
		"number":goodsNumber,
		"money":orderMoney,
	})
	e.OK(req.GetId(), "创建成功")
}

// Update 修改Orders
// @Summary 修改Orders
// @Description 修改Orders
// @Tags Orders
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.OrdersUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/orders/{id} [put]
// @Security Bearer
func (e Orders) Update(c *gin.Context) {
	req := dto.OrdersUpdateReq{}
	s := service.Orders{}
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
	userDto, err := customUser.GetUserDto(e.Orm, c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}
	err = s.Update(e.getTableName(userDto.CId), &req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("修改Orders失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "修改成功")
}

// Delete 删除Orders
// @Summary 删除Orders
// @Description 删除Orders
// @Tags Orders
// @Param data body dto.OrdersDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/orders [delete]
// @Security Bearer
func (e Orders) Delete(c *gin.Context) {
	s := service.Orders{}
	req := dto.OrdersDeleteReq{}
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

	err = s.Remove(e.getTableName(userDto.CId), &req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("删除Orders失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "删除成功")
}
