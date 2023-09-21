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
	models3 "go-admin/common/models"
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
	//查询是否进行了分表
	orderTableName := business.GetTableName(userDto.CId, e.Orm)

	p := actions.GetPermissionFromContext(c)
	list := make([]models.Orders, 0)
	var count int64
	req.CId = userDto.CId
	err = s.GetPage(orderTableName, &req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取订单失败,%s", err.Error()))
		return
	}
	//统一查询商家shop_id
	cacheShopId := make([]int, 0)
	//统一查询订单关联的规格
	cacheOrderId := make([]int, 0)
	for _, row := range list {
		cacheShopId = append(cacheShopId, row.ShopId)
		cacheOrderId = append(cacheOrderId, row.Id)
	}
	cacheShopId = utils.RemoveRepeatInt(cacheShopId)

	//查询客户信息
	cacheShopObject := make([]models2.Shop, 0)
	e.Orm.Model(&models2.Shop{}).Select("name,id,phone,address,user_name").Where("c_id = ? and id in ?",
		userDto.CId, cacheShopId).Find(&cacheShopObject)
	//保存为map
	//应该是获取的用户下单的地址
	cacheShopMap := make(map[int]map[string]interface{}, 0)
	for _, k := range cacheShopObject {
		cacheShopMap[k.Id] = map[string]interface{}{
			"name":  k.Name,
			"phone": k.Phone,
			"id":    k.Id,
			"username":k.UserName,
			"address":k.Address,
		}
	}
	//查询订单关联的规格信息
	var orderSpecs []models.OrderSpecs
	//是一个分表的名称
	specsTable := business.OrderSpecsTableName(orderTableName)

	e.Orm.Table(specsTable).Select("order_id,goods_name,number").Where("order_id in ?", cacheOrderId).Find(&orderSpecs)

	//缓存订单关联的商品
	orderSpecsMap := make(map[string]map[string]dto.OrderSpecsRow, 0)
	for _, k := range orderSpecs {
		rows, ok := orderSpecsMap[k.OrderId]
		fmt.Println("ok", ok, rows)
		if !ok {
			specsMap := make(map[string]dto.OrderSpecsRow, 0)

			specsMap[k.SpecsName] = dto.OrderSpecsRow{
				Number: k.Number,
			}
			orderSpecsMap[k.OrderId] = specsMap

		} else {
			parentDat := rows[k.SpecsName]
			parentDat.Number += k.Number
			rows[k.SpecsName] = parentDat

			orderSpecsMap[k.OrderId] = rows
		}
	}

	result := make([]map[string]interface{}, 0)
	for _, row := range list {

		specsDataMap := orderSpecsMap[row.OrderId]


		//如果支付金额为0
		PayMoney:=row.PayMoney
		if row.PayMoney == 0  && row.PayType < global.PayTypeOnlineWechat{
			if row.DeductionMoney > 0 {
				PayMoney = row.DeductionMoney
			}
		}
		r := map[string]interface{}{
			"order_id":   row.OrderId,
			"order_no_id":row.OrderNoId,
			"shop": cacheShopMap[row.ShopId],
			"cycle_place": row.CreatedAt.Format("2006-01-02"), 			//下单周期
			"delivery_time":     row.DeliveryTime.Format("2006-01-02"), 			//配送周期
			"delivery_str": row.DeliveryStr,
			"count":          row.Number,
			"money":         PayMoney,
			"delivery_type":global.GetExpressCn(row.DeliveryType), //配送类型
			"pay_type":global.GetPayType(row.PayType),//支付类型
			"source_type":global.GetOrderSource(row.SourceType),//订单来源
			"status":         global.OrderStatus(row.Status),
			"pay_status":     global.GetOrderPayStatus(row.PayStatus),
			"created_at":     row.CreatedAt,
			"goods":          specsDataMap,
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

	orderId:=c.Param("orderId")
	var object models.Orders

	orderTableName := business.GetTableName(userDto.CId, e.Orm)
	orderErr := e.Orm.Table(orderTableName).Where("order_id = ?",orderId).First(&object).Error
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
		"order_id":       object.Id,
		"created_at":     object.CreatedAt.Format("2006-01-02 15:04:05"),
		"delivery_time":     object.DeliveryTime.Format("2006-01-02"),
		"delivery_str":      object.DeliveryStr,
		"pay":            global.GetPayType(object.PayType),
		"pay_status_str": global.GetOrderPayStatus(object.PayStatus),
		"pay_status":     object.PayStatus,
		"shop_name":      shopRow.Name,
		"shop_username":  shopRow.UserName,
		"shop_phone":     shopRow.Phone,
		"shop_address":   shopRow.Address,
	}
	var orderSpecs []models.OrderSpecs
	//是一个分表的名称
	specsTable := business.OrderSpecsTableName(orderTableName)

	e.Orm.Table(specsTable).Where("order_id = ?", orderId).Find(&orderSpecs)

	specsList := make([]map[string]interface{}, 0)
	for _, row := range orderSpecs {
		ss := map[string]interface{}{
			"id":         row.Id,
			"name":       row.SpecsName,
			"goods_name":row.GoodsName,
			"created_at": row.CreatedAt.Format("2006-01-02 15:04:05"),
			"specs":      fmt.Sprintf("%v%v", row.Number, row.Unit),
			"status":     global.OrderStatus(row.Status),
			"money":      row.Money,
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
	//抵扣类型
	//查询选择的用户额度是否够
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
	var DeductionAllMoney float64
	switch req.DeductionType {
	case global.DeductionBalance:
		DeductionAllMoney = shopObject.Balance
	case global.DeductionCredit:
		DeductionAllMoney = shopObject.Credit
	default:
		e.OK(-1, "未知抵扣类型")
		return
		
	}
	fmt.Println("DeductionAllMoney",DeductionAllMoney)
	//因为折扣的算法和下单是在一起,所以就下单的时候检测下总价是否可以达到减扣的，防止出现负数 或者钱不够还下单
	var AllOrderMoney float64
	for _, goodsList := range req.Goods {
		for _, spec := range goodsList {
			//当前用户是vip哪个等级,使用客户的VIP的价格
			var goodsVip models.GoodsVip
			e.Orm.Model(&goodsVip).Select("id,custom_price").Where("goods_id = ? and specs_id = ? and grade_id = ?",spec.GoodsId,spec.Id,shopObject.GradeId).Limit(1).Find(&goodsVip)

			var specPrice float64
			if goodsVip.Id > 0 {
				//使用vip的价格
				specPrice = goodsVip.CustomPrice
			}else {
				specPrice = spec.Price
			}
			AllOrderMoney += utils.RoundDecimalFlot64(specPrice  * float64(spec.Number))
		}
	}

	if DeductionAllMoney < AllOrderMoney {
		e.OK(-1, "抵扣额不足！！！")
		return
	}

	orderTableName := business.GetTableName(userDto.CId, e.Orm)
	specsTable := business.OrderSpecsTableName(orderTableName)
	orderExtend := business.OrderExtendTableName(orderTableName)
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


	//保存商品和规格的一些映射
	goodsCacheList:=make(map[int]service.ValetOrderGoodsRow,0)

	//一个订单下了很多商品
	orderRow := &models.Orders{
		Enable:  true,
		ShopId:  req.Shop,
		Line:    lineName,
		LineId:  lineObject.Id,
		CId:     userDto.CId,
		SourceType: global.OrderSourceValet, //代客下单
	}
	orderId := utils.GenUUID()
	orderRow.OrderId = fmt.Sprintf("%v",orderId)
	//orderRow.Id = orderId
	//代客下单,需要把配送周期保存，方便周期配送

	orderRow.PayType = req.DeductionType //抵扣的方式
	orderRow.Status = global.OrderStatusWaitSend
	orderRow.PayStatus = global.OrderStatusPaySuccess
	orderRow.DeliveryCode = service.DeliveryCode()
	orderRow.PayTime = models3.XTime{
		Time:time.Now(),
	}
	orderRow.Phone = shopObject.Phone
	orderRow.DeliveryTime = s.CalculateTime(DeliveryObject.GiveDay)
	orderRow.DeliveryStr = DeliveryObject.GiveTime
	orderRow.DeliveryUid = DeliveryObject.Uid
	orderRow.DeliveryID = DeliveryObject.Id
	orderRow.DeliveryType = global.ExpressLocal

	orderRow.CreateBy = userDto.UserId
	e.Orm.Table(orderExtend).Create(&models.OrderExtend{
		CId: userDto.CId,
		OrderId: orderRow.OrderId,
		Driver:  DriverObject.Name,
		Phone: DriverObject.Phone,
		UserName: shopObject.UserName,
		UserAddress: shopObject.Address,
		Buyer: req.Desc,
	})

	var orderMoney float64
	var goodsNumber int

	for _, goodsList := range req.Goods {
		//fmt.Println("classId",classId)
		specsOrderId := make([]int, 0)
		for _, spec := range goodsList {
			//如果商品不存在
			var goodsObject models2.Goods
			e.Orm.Model(&models2.Goods{}).Select("id,sale,inventory,name,image").Where("id = ? and c_id = ? and enable = ?", spec.GoodsId, userDto.CId, true).Limit(1).Find(&goodsObject)
			if goodsObject.Id == 0 {
				continue
			}

			//如果下单的次数>库存的值，那就是非法数据 直接跳出
			var goodsSpecs models.GoodsSpecs
			e.Orm.Model(&models.GoodsSpecs{}).Where("id = ? and c_id = ?", spec.Id, userDto.CId).Limit(1).Find(&goodsSpecs)
			if goodsSpecs.Id == 0 {
				continue
			}
			if spec.Number > goodsSpecs.Inventory {
				continue
			}
			//当前用户是vip哪个等级,使用客户的VIP的价格
			var goodsVip models.GoodsVip
			e.Orm.Model(&goodsVip).Select("id,custom_price").Where("goods_id = ? and specs_id = ? and grade_id = ?",spec.GoodsId,spec.Id,shopObject.GradeId).Limit(1).Find(&goodsVip)

			var specPrice float64
			if goodsVip.Id > 0 && goodsVip.CustomPrice > 0 {
				//使用vip的价格
				specPrice = goodsVip.CustomPrice
			}else {
				//不是VIP就使用售价
				specPrice = spec.Price
			}

			specPrice = utils.RoundDecimalFlot64(specPrice)

			specRow := &models.OrderSpecs{
				OrderId: orderRow.OrderId,
				Number:    spec.Number,
				Money:    specPrice,
				Unit:      spec.Unit,
				GoodsName: goodsObject.Name,
				GoodsId: goodsObject.Id,
				SpecsName: goodsSpecs.Name,
				SpecId:   spec.Id,
				CId: userDto.CId,
			}
			if goodsSpecs.Image == ""{
				//拿商品的首图吧
				if goodsObject.Image != ""{
					specRow.Image = strings.Split(goodsObject.Image,",")[0]
				}
			}else {
				specRow.Image = goodsSpecs.Image
			}
			txRes:=e.Orm.Table(specsTable).Create(&specRow)
			if txRes.Error !=nil{
				continue
			}
			//规格减库存 + 销量
			e.Orm.Model(&models.GoodsSpecs{}).Where("id = ? and c_id =?", spec.Id, userDto.CId).Updates(map[string]interface{}{
				"inventory": goodsSpecs.Inventory - spec.Number,
				"sale":	goodsSpecs.Sale + spec.Number,
			})

			orderMoney += utils.RoundDecimalFlot64(specPrice  * float64(spec.Number))
			goodsNumber += spec.Number

			//缓存订单ID
			specsOrderId = append(specsOrderId, specRow.Id)


			cacheValetOrderRow,ok:=goodsCacheList[spec.GoodsId]
			if !ok{
				goodsCacheList[spec.GoodsId] = service.ValetOrderGoodsRow{
					Number: spec.Number,
				}
			}else{
				cacheValetOrderRow.Number +=spec.Number
				goodsCacheList[spec.GoodsId] = cacheValetOrderRow
			}

		}

	}
	orderRow.Number = goodsNumber

	orderMoney = utils.RoundDecimalFlot64(orderMoney)
	orderRow.PayMoney = orderMoney
	orderRow.OrderMoney = orderMoney
	orderRow.GoodsMoney = orderMoney
	orderRow.DeductionMoney = orderMoney
	e.Orm.Table(orderTableName).Create(&orderRow)

	for goodsId,goodsRow:=range goodsCacheList{
		//商品减库存 + 销量
		var goodsObject models.Goods
		e.Orm.Model(&models.Goods{}).Select("sale,inventory,id").Where("id = ? and c_id = ? and enable = ?", goodsId, userDto.CId, true).Limit(1).Find(&goodsObject)
		if goodsObject.Id == 0 {
			continue
		}
		e.Orm.Model(&models.Goods{}).Where("c_id = ? and id = ?", userDto.CId, goodsId).Updates(map[string]interface{}{
			"sale":    goodsObject.Sale + goodsRow.Number  ,
			"inventory": goodsObject.Inventory - goodsRow.Number,
		})
	}
	//授信额减免

	switch req.DeductionType {
	case global.DeductionBalance:

		Balance:= shopObject.Balance - orderMoney
		if Balance < 0 {
			Balance = 0
		}
		Balance = utils.RoundDecimalFlot64(Balance)
		e.Orm.Model(&models2.Shop{}).Where("id = ?",shopObject.Id).Updates(map[string]interface{}{
			"amount":Balance,
		})
		//余额变动记录
		row:=models2.ShopBalanceLog{
			CId: userDto.CId,
			ShopId: shopObject.Id,
			Money: orderMoney,
			Scene:fmt.Sprintf("用户[%v] 代客下单,使用余额抵扣费:%v",userDto.Username,orderMoney),
			Action: global.UserNumberReduce, //抵扣
			Type: global.ScanShopUse,
		}
		e.Orm.Create(&row)
	case global.DeductionCredit:
		Credit:=  shopObject.Credit - orderMoney
		if Credit < 0 {
			Credit = 0
		}
		Credit = utils.RoundDecimalFlot64(Credit)
		e.Orm.Model(&models2.Shop{}).Where("id = ?",shopObject.Id).Updates(map[string]interface{}{
			"credit":Credit,
		})
		//授信变动记录
		row:=models2.ShopBalanceLog{
			CId: userDto.CId,
			ShopId: shopObject.Id,
			Money: orderMoney,
			Scene:fmt.Sprintf("用户[%v] 代客下单,使用授信额抵扣费:%v",userDto.Username,orderMoney),
			Action: global.UserNumberReduce, //抵扣
			Type: global.ScanShopUse,
		}
		e.Orm.Create(&row)
	default:
		e.OK(-1, "未知抵扣类型")
		return

	}

	e.OK(1, "successful")
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

func (e Orders) ShopOrderList(c *gin.Context) {
	s := service.Orders{}
	req := dto.OrdersShopGetPageReq{}
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
	shopId := c.Param("id")
	fmt.Println("商家ID", shopId)
	result := make([]map[string]interface{}, 0)
	var count int64
	e.PageOK(result, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
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
	data.Status = global.OrderStatusWaitSend
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
	//todo:配送周期
	data.DeliveryTime = timeConfResult.CycleTime
	data.DeliveryStr = timeConfResult.CycleStr
	data.DeliveryUid = timeConfResult.RandUid
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
		OrderId:   fmt.Sprintf("%v",data.Id),
		Driver:  DriverObject.Name,
		Phone:   DriverObject.Phone,
	})

	//分表检测
	specsTable := business.OrderSpecsTableName(orderTableName)

	//收益总价
	var orderMoney float64
	//售出量
	var goodsSaleNumber int
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
		goodsSaleNumber += goods.Number

		e.Orm.Table(specsTable).Create(&models.OrderSpecs{
			OrderId:   fmt.Sprintf("%v",data.Id),
			SpecsName: goods.Name,
			Unit:      goods.Unit,
			Status:    global.OrderStatusWaitSend,
			Money:     goods.Money,
			Number:    goods.Number,
		})
	}
	inventory := goodsObject.Inventory
	SaleNumber := goodsObject.Sale
	//销量+,库存-
	SaleNumber += goodsSaleNumber
	inventory -= goodsSaleNumber
	e.Orm.Model(&models.Goods{}).Where("c_id = ? and id = ?", userDto.CId, goodsObject.Id).Updates(map[string]interface{}{
		"sale":      SaleNumber,
		"inventory": inventory,
	})
	e.Orm.Model(&models.Orders{}).Table(orderTableName).Where("id = ?", data.Id).Updates(map[string]interface{}{
		"number": goodsSaleNumber,
		"money":  orderMoney,
	})
	//订单创建成功了,同时做一个周期列表数据得保存
	orderCreateName := time.Now().Format("2006-01-02")
	var cycleObject models.OrderCycleList
	e.Orm.Model(&models.OrderCycleList{}).Where("c_id = ? and name = ? and uid = ?",
		userDto.CId, orderCreateName, timeConfResult.RandUid).Limit(1).Find(&cycleObject)

	//周期订单数据更新更新
	SaleMoney := cycleObject.SaleMoney
	SaleMoney += orderMoney
	GoodsAll := cycleObject.GoodsAll
	GoodsAll += goodsSaleNumber

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
			SaleMoney: SaleMoney,
			GoodsAll:  GoodsAll,
		}
		e.Orm.Create(&cycleMode)
	} else {
		e.Orm.Model(&models.OrderCycleList{}).Where("id = ?", cycleObject.Id).Updates(map[string]interface{}{
			"sale_money": SaleMoney,
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
