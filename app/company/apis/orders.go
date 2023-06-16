package apis

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"
	models2 "go-admin/cmd/migrate/migration/models"
	"go-admin/common/business"
	customUser "go-admin/common/jwt/user"
	"go-admin/common/utils"
	"go-admin/global"
	"gorm.io/gorm"
	"strings"
	"time"

	"go-admin/app/company/models"
	"go-admin/app/company/service"
	"go-admin/app/company/service/dto"
	"go-admin/common/actions"
)

type Orders struct {
	api.Api
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
	err = s.GetPage(business.GetTableName(userDto.CId, e.Orm), &req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取订单失败,%s", err.Error()))
		return
	}
	//统一查询商家shop_id
	cacheShopId := make([]int, 0)
	for _, row := range list {
		cacheShopId = append(cacheShopId, row.ShopId)
	}
	cacheShopId = utils.RemoveRepeatInt(cacheShopId)
	//查询到对象
	cacheShopObject := make([]models2.Shop, 0)
	e.Orm.Model(&models2.Shop{}).Select("name,id,phone").Where("c_id = ? and id in ?",
		userDto.CId, cacheShopId).Find(&cacheShopObject)
	//保存为map
	cacheShopMap := make(map[int]map[string]interface{}, 0)
	for _, k := range cacheShopObject {
		cacheShopMap[k.Id] = map[string]interface{}{
			"name":  k.Name,
			"phone": k.Phone,
			"id":    k.Id,
		}
	}

	result := make([]map[string]interface{}, 0)
	for _, row := range list {
		r := map[string]interface{}{
			"id":   fmt.Sprintf("%v", row.Id),
			"shop": cacheShopMap[row.ShopId],
			//下单周期
			"cycle_place": row.CreatedAt.Format("2006-01-02"),
			//配送周期
			"cycle_give":     row.CycleTime.Format("2006-01-02"),
			"cycle_give_str": row.CycleStr,
			"count":          row.Number,
			"money":          row.Money,
			"s":              row.Status,
			"p":              row.PayStatus,
			"status":         global.OrderStatus(row.Status),
			"pay_status":     global.OrderPayStatus(row.PayStatus),
			"created_at":     row.CreatedAt,
		}
		result = append(result, r)
	}
	e.PageOK(result, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
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

	orderTableName := business.GetTableName(userDto.CId, e.Orm)
	orderErr := e.Orm.Table(orderTableName).First(&object, req.Id).Error
	if orderErr != nil && errors.Is(orderErr, gorm.ErrRecordNotFound) {

		e.Error(500, orderErr, "订单不存在")
		return
	}
	if orderErr != nil {
		e.Error(500, orderErr, "订单不存在")
		return
	}

	var shopRow models2.Shop
	e.Orm.Model(&models2.Shop{}).Where("id = ? and c_id = ?", object.ShopId, userDto.CId).Limit(1).Find(&shopRow)

	result := map[string]interface{}{
		"order_id":   object.Id,
		"created_at": object.CreatedAt,
		"cycle_time": object.CycleTime,
		//"cycle_str":     object.CycleStr,
		"pay":           global.GetPayStr(object.Pay),
		"shop_name":     shopRow.Name,
		"shop_username": shopRow.UserName,
		"shop_phone":    shopRow.Phone,
		"shop_address":  shopRow.Address,
	}
	var orderSpecs []models.OrderSpecs
	//是一个分表的名称
	specsTable := business.OrderSpecsTableName(orderTableName)

	e.Orm.Table(specsTable).Where("order_id = ?", object.Id).Find(&orderSpecs)

	specsList := make([]map[string]interface{}, 0)
	for _, row := range orderSpecs {
		ss := map[string]interface{}{
			"name":   row.SpecsName,
			"spec":   row.Unit,
			"status": row.Status,
			"money":  row.Money,
		}
		specsList = append(specsList, ss)
	}
	result["specs_list"] = specsList
	e.OK(result, "查询成功")
}

// 代客下单就没有时间的限制了
func (e Orders) ValetOrder(c *gin.Context) {
	req := dto.ValetOrderReq{}
	s := service.Orders{}
	err := e.MakeContext(c).
		Bind(&req).
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
	orderTableName := business.GetTableName(userDto.CId, e.Orm)
	specsTable := business.OrderSpecsTableName(orderTableName)
	orderExtend := business.OrderExtendTableName(orderTableName)
	var shopObject models2.Shop
	e.Orm.Model(&models2.Shop{}).Where("id = ? and enable =? and c_id = ?", req.Shop, true, userDto.CId).Limit(1).Find(&shopObject)
	if shopObject.Id == 0 {
		e.Error(500, errors.New("商户不存在"), "商户不存在")
		return
	}
	if shopObject.LineId == 0 {
		e.Error(500, errors.New("商家暂无路线"), "商家暂无路线")
		return
	}
	var DeliveryObject models.CycleTimeConf
	e.Orm.Model(&models2.CycleTimeConf{}).Where("id = ? and enable =? and c_id = ?", req.Cycle, true, userDto.CId).Limit(1).Find(&DeliveryObject)
	if DeliveryObject.Id == 0 {
		e.Error(500, nil, "时间区间不存在")
		return
	}
	var lineObject models2.Line
	e.Orm.Model(&models2.Line{}).Where("id = ? and c_id = ? and enable = ?", shopObject.LineId, userDto.CId, true).Limit(1).Find(&lineObject)

	lineName := lineObject.Name
	var DriverObject models2.Driver
	e.Orm.Model(&models2.Driver{}).Where("id = ? and c_id = ? and enable = ?", lineObject.DriverId, userDto.CId, true).Limit(1).Find(&DriverObject)
	driverName := DriverObject.Name
	for classId, goodsList := range req.Goods {

		orderRow := &models.Orders{
			Enable:  true,
			ShopId:  req.Shop,
			Line:    lineName,
			LineId:  lineObject.Id,
			CId:     userDto.CId,
			ClassId: classId,
		}
		//
		orderId := utils.GenUUID()
		orderRow.Id = orderId
		//代客下单,需要把配送周期保存，方便周期配送
		orderRow.CycleTime = s.CalculateTime(DeliveryObject.GiveDay)
		orderRow.CycleStr = DeliveryObject.GiveTime
		orderRow.CycleUid = DeliveryObject.Uid
		orderRow.CreateBy = userDto.UserId
		orderRow.GoodsId = goodsList[0].GoodsId

		e.Orm.Table(orderExtend).Create(&models.OrderExtend{
			OrderId: orderRow.Id,
			Desc:    req.Desc,
			Driver:  driverName,
			Source:  1,
		})

		var orderMoney float64
		var goodsNumber int
		var goodsName string
		specsOrderId := make([]int, 0)
		for _, spec := range goodsList {
			//如果下单的次数>库存的值，那就是非法数据 直接跳出
			var goodsSpecs models.GoodsSpecs
			e.Orm.Model(&models.GoodsSpecs{}).Where("id = ? and c_id = ?", spec.Id, userDto.CId).Limit(1).Find(&goodsSpecs)
			if goodsSpecs.Id == 0 {
				continue
			}
			if spec.Number > goodsSpecs.Inventory {
				continue
			}

			orderMoney += spec.Price
			goodsNumber += spec.Number

			e.Orm.Model(&models.GoodsSpecs{}).Where("id = ? and c_id =?", spec.Id, userDto.CId).Updates(map[string]interface{}{
				"inventory": goodsSpecs.Inventory - spec.Number,
			})

			specRow := &models.OrderSpecs{
				Number:    spec.Number,
				Money:     spec.Price,
				Unit:      spec.Unit,
				SpecsName: goodsSpecs.Name,
			}
			var goodsObject models2.Goods
			e.Orm.Model(&models2.Goods{}).Where("id = ? and c_id = ? and enable = ?", spec.GoodsId, userDto.CId, true).Limit(1).Find(&goodsObject)
			if goodsObject.Id == 0 {
				continue
			}
			goodsName = goodsObject.Name
			e.Orm.Table(specsTable).Create(specRow)
			specsOrderId = append(specsOrderId, specRow.Id)
		}
		orderRow.Number = goodsNumber
		orderRow.Money = orderMoney
		orderRow.GoodsName = goodsName
		e.Orm.Table(orderTableName).Create(orderRow)
		e.Orm.Table(specsTable).Where("id in ?", specsOrderId).Updates(map[string]interface{}{
			"order_id": orderRow.Id,
		})

	}
	e.OK("", "successful")
	return
}
func (e Orders) ToolsOrders(c *gin.Context) {
	req := dto.ToolsOrdersUpdateReq{}
	err := e.MakeContext(c).
		Bind(&req).
		MakeOrm().
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
	orderTableName := business.GetTableName(userDto.CId, e.Orm)
	specsTable := business.OrderSpecsTableName(orderTableName)
	switch req.Type {
	case global.OrderToolsActionStatus: //状态更新
		switch req.Status {
		case global.OrderStatusWait:
		case global.OrderStatusOk:
		case global.OrderStatusReturn:
		case global.OrderStatusRefund:
		case global.OrderStatusLoading:

		default:
			e.Error(500, nil, "状态非法")
			return

		}
		e.Orm.Table(orderTableName).Where("id = ? and enable = ?", req.Id, true).Updates(map[string]interface{}{
			"status":    req.Status,
			"desc":      req.Desc,
			"update_by": userDto.UserId,
		})
		e.Orm.Table(specsTable).Where("order_id = ?", req.Id).Updates(map[string]interface{}{
			"status": req.Status,
		})
	case global.OrderToolsActionDelivery: //周期更改
		if req.Delivery > 0 {
			e.Orm.Table(orderTableName).Where("id = ? and enable = ?", req.Id, true).Updates(map[string]interface{}{
				"delivery":  req.Delivery,
				"update_by": userDto.UserId,
			})
		} else {
			e.Error(500, nil, "状态非法")
			return
		}
	}

	e.OK("", "successful")
	return
}

func (e Orders) OrderCycleList(c *gin.Context) {
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
	//获取所有的周期列表
	datalist := make([]models.OrderCycleList, 0)
	e.Orm.Model(&models.OrderCycleList{}).Select("name,uid,id,start_time,end_time,cycle_time,cycle_str").Where(
		"c_id = ?", userDto.CId).Order(global.OrderTimeKey).Find(&datalist)

	//下单周期
	createTime := make([]map[string]interface{}, 0)
	//配送周期
	giveTime := make([]map[string]interface{}, 0)

	//下单周期和配送周期是成对出现的,
	for _, row := range datalist {

		create := row.StartTime.Format("2006-01-02")
		t1 := map[string]interface{}{
			"id": row.Id,
			"t":  fmt.Sprintf("%v_%v", create, row.Uid),
			"value": fmt.Sprintf("%v 开始:%v至结束:%v",
				row.Name, create, row.EndTime.Format("15:04")),
		}
		createTime = append(createTime, t1)
		give := row.CycleTime.Format("2006-01-02")
		t2 := map[string]interface{}{
			"id":    row.Id,
			"t":     fmt.Sprintf("%v_%v", give, row.Uid),
			"value": fmt.Sprintf("%v %v", give, row.CycleStr),
		}
		giveTime = append(giveTime, t2)
	}
	//放到不同的列表中
	result := map[string]interface{}{
		"cycle_create": createTime,
		"cycle_give":   giveTime,
	}
	e.OK(result, "successful")

	return
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
	e.Orm.Model(&models.CycleTimeConf{}).Where("c_id = ? and enable = ?", userDto.CId, true).Order(global.OrderLayerKey).Find(&lists)

	result := make([]map[string]interface{}, 0)
	for _, row := range lists {
		m := map[string]interface{}{
			"id":      row.Id,
			"placing": "",
			"give":    "",
		}
		placing := ""
		give := ""
		if row.GiveDay == 0 {
			give = fmt.Sprintf("下单当天,%v", row.GiveTime)
		} else {
			give = fmt.Sprintf("下单后 第%v天,%v", row.GiveDay, row.GiveTime)
		}
		switch row.Type {
		case global.CyCleTimeDay:
			placing = fmt.Sprintf("每天%v-%v", row.StartTime, row.EndTime)
		case global.CyCleTimeWeek:
			placing = fmt.Sprintf("每周%v %v-每周%v %v", global.WeekIntToMsg(row.StartWeek), row.StartTime,
				global.WeekIntToMsg(row.EndWeek), row.EndTime,
			)
		default:
			continue
		}
		m["placing"] = placing
		m["give"] = give
		result = append(result, m)
	}
	e.PageOK(result, len(result), 1, -1, "successful")
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
	result := s.ValidTimeConf(userDto.CId)
	if !result.Valid {
		e.Error(500, errors.New("非下单时间段"), "非下单时间段")
		return
	}
	e.OK(map[string]interface{}{
		"time": result.CycleTime,
		"str":  result.CycleStr,
	}, "successful")
	return
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
	var goodsObject models2.Goods
	e.Orm.Model(&models2.Goods{}).Where("id = ? and c_id = ? and enable = ?", req.GoodsId, userDto.CId, true).Limit(1).Find(&goodsObject)
	if goodsObject.Id == 0 {
		e.Error(500, errors.New("无此商品"), "无此商品")
		return
	}

	//根据下单的时间区间，来匹配
	//todo:配送周期
	//根据下单的时间区间来自动匹配,
	//1.查询这个时间段内是否配置了cycle_time_conf得值,如果创建了进行关联即可
	//2.如果没有这个时间
	timeConfResult := s.ValidTimeConf(userDto.CId)
	if !timeConfResult.Valid {
		e.Error(500, errors.New("非下单时间段"), "非下单时间段")
		return
	}

	//下单库存校验
	//存放有用的规格
	ofSpecsUseList := make([]dto.OrderGoodsSpecs, 0)
	errorSpecs := make([]string, 0)
	for _, goods := range req.GoodsSpecs {
		var goodsSpecs models.GoodsSpecs
		e.Orm.Model(&models.GoodsSpecs{}).Where("id = ? and c_id = ?", goods.SpecsId, userDto.CId).Limit(1).Find(&goodsSpecs)
		if goodsSpecs.Id == 0 {
			errorSpecs = append(errorSpecs, fmt.Sprintf("暂无此,%v规格", goods.Name))
			continue
		}
		//库存校验,库存是否 > 下单库存
		if goods.Number > goodsSpecs.Inventory {
			errorSpecs = append(errorSpecs, fmt.Sprintf("%v,库存不足", goods.Name))
			continue
		}
		goods.Inventory = goodsSpecs.Inventory
		ofSpecsUseList = append(ofSpecsUseList, goods)
	}
	if len(errorSpecs) > 0 {
		e.Error(500, errors.New(strings.Join(errorSpecs, ",")), strings.Join(errorSpecs, ","))
		return
	}
	userId := user.GetUserId(c)
	//todo:获取订单表名
	orderTableName := business.GetTableName(userDto.CId, e.Orm)

	var data models.Orders
	data.Id = utils.GenUUID()
	data.CId = userDto.CId
	data.Enable = true
	data.Status = global.OrderStatusWait
	var shopObject models2.Shop
	e.Orm.Model(&models2.Shop{}).Where("id = ? and c_id = ? and enable = ?", req.ShopId, userDto.CId, true).Limit(1).Find(&shopObject)
	if shopObject.Id == 0 {
		e.Error(500, errors.New("暂无商家"), "暂无商家")
		return
	}
	if shopObject.LineId == 0 {
		e.Error(500, errors.New("商家暂无路线"), "商家暂无路线")
		return
	}
	//选择了商家,获取商家关联的路线
	data.ShopId = req.ShopId
	//线路
	var lineObject models2.Line
	e.Orm.Model(&models2.Line{}).Where("id = ? and c_id = ? and enable = ?", shopObject.LineId, userDto.CId, true).Limit(1).Find(&lineObject)

	if lineObject.Id == 0 {
		e.Error(500, errors.New("商家暂无路线"), "商家暂无路线")
		return
	}
	data.Line = lineObject.Name
	data.LineId = lineObject.Id
	//商品
	data.GoodsId = goodsObject.Id
	data.GoodsName = goodsObject.Name
	//todo:配送周期
	data.CycleTime = timeConfResult.CycleTime
	data.CycleStr = timeConfResult.CycleStr
	data.CycleUid = timeConfResult.RandUid
	data.CreateBy = userId
	createErr := e.Orm.Table(orderTableName).Create(&data).Error
	if createErr != nil {
		e.Error(500, createErr, "订单创建失败")
		return
	}
	var DriverObject models2.Driver
	e.Orm.Model(&models2.Driver{}).Where("id = ? and c_id = ? and enable = ?", lineObject.DriverId, userDto.CId, true).Limit(1).Find(&DriverObject)

	//扩展表
	orderExtend := business.OrderExtendTableName(orderTableName)
	e.Orm.Table(orderExtend).Create(&models.OrderExtend{
		OrderId: data.Id,
		Desc:    req.Desc,
		Source:  0,
		Driver:  DriverObject.Name,
		Phone:   DriverObject.Phone,
	})

	//分表检测
	specsTable := business.OrderSpecsTableName(orderTableName)

	//收益总价
	var orderMoney float64
	//售出量
	var goodsSoldNumber int
	//规格下单
	//下单后:1.对商品减库存 2.对商品增加出售量
	for _, goods := range ofSpecsUseList {
		//库存校验,库存是否 = 规格库存
		if goods.Number > goods.Inventory {
			continue
		}
		//对商品规格的库存减法

		nowIN := goods.Inventory - goods.Number
		e.Orm.Model(&models.GoodsSpecs{}).Where("id = ? and c_id =?", goods.SpecsId, userDto.CId).Updates(map[string]interface{}{
			"inventory": nowIN,
		})
		orderMoney += goods.Money
		goodsSoldNumber += goods.Number

		e.Orm.Table(specsTable).Create(&models.OrderSpecs{
			OrderId:   data.Id,
			SpecsName: goods.Name,
			Unit:      goods.Unit,
			Status:    global.OrderStatusWait,
			Money:     goods.Money,
			Number:    goods.Number,
		})
	}
	inventory := goodsObject.Inventory
	soldNumber := goodsObject.Sold
	//销量+,库存-
	soldNumber += goodsSoldNumber
	inventory -= goodsSoldNumber
	e.Orm.Model(&models.Goods{}).Where("c_id = ? and id = ?", userDto.CId, goodsObject.Id).Updates(map[string]interface{}{
		"sold":      soldNumber,
		"inventory": inventory,
	})
	e.Orm.Model(&models.Orders{}).Table(orderTableName).Where("id = ?", data.Id).Updates(map[string]interface{}{
		"number": goodsSoldNumber,
		"money":  orderMoney,
	})
	//订单创建成功了,同时做一个周期列表数据得保存
	orderCreateName := time.Now().Format("2006-01-02")
	var cycleObject models.OrderCycleList
	e.Orm.Model(&models.OrderCycleList{}).Where("c_id = ? and name = ? and uid = ?",
		userDto.CId, orderCreateName, timeConfResult.RandUid).Limit(1).Find(&cycleObject)

	//周期订单数据更新更新
	SoldMoney := cycleObject.SoldMoney
	SoldMoney += orderMoney
	GoodsAll := cycleObject.GoodsAll
	GoodsAll += goodsSoldNumber

	//如果没有这个周期,那就进行创建
	if cycleObject.Id == 0 {
		cycleMode := models.OrderCycleList{
			CId:       userDto.CId,
			Name:      orderCreateName,
			Uid:       timeConfResult.RandUid,
			StartTime: timeConfResult.StartTime,
			EndTime:   timeConfResult.EndTime,
			CycleTime: timeConfResult.CycleTime,
			CycleStr:  timeConfResult.CycleStr,
			SoldMoney: SoldMoney,
			GoodsAll:  GoodsAll,
		}
		e.Orm.Create(&cycleMode)
	} else {
		e.Orm.Model(&models.OrderCycleList{}).Where("id = ?", cycleObject.Id).Updates(map[string]interface{}{
			"sold_money": SoldMoney,
			"goods_all":  GoodsAll,
		})
	}

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
	err = s.Update(business.GetTableName(userDto.CId, e.Orm), &req, p)
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

	err = s.Remove(business.GetTableName(userDto.CId, e.Orm), &req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("删除Orders失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "删除成功")
}
