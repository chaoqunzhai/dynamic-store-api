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
	cDto "go-admin/common/dto"
	customUser "go-admin/common/jwt/user"
	models3 "go-admin/common/models"
	"go-admin/common/utils"
	"go-admin/config"
	"go-admin/global"
	"gorm.io/gorm"
	"strconv"
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
	splitTableRes := business.GetTableName(userDto.CId, e.Orm)

	p := actions.GetPermissionFromContext(c)

	//处理不是数字字符串的问题
	if req.Uid != ""{
		uid,uidErr:=strconv.Atoi(req.Uid)
		if uidErr==nil{
			req.Uid = fmt.Sprintf("%v",uid)
		}else {
			req.Uid = ""
		}
	}

	//配送周期传入的值是:14_2023-09-23
	//配送周期查询

	//下单周期查询


	list := make([]models.Orders, 0)
	var count int64
	req.CId = userDto.CId
	err = s.GetPage(splitTableRes.OrderTable, &req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取订单失败,%s", err.Error()))
		return
	}
	//统一查询商家shop_id
	cacheShopId := make([]int, 0)
	//统一查询用户地址
	cacheAddressId:=make([]int,0)
	cacheStoreAddressId:=make([]int,0)
	for _, row := range list {
		if row.ShopId > 0 {
			cacheShopId = append(cacheShopId,row.ShopId)
		}

		if row.AddressId > 0 {
			switch row.DeliveryType {
			case global.ExpressStore:
				cacheStoreAddressId = append(cacheStoreAddressId,row.AddressId)
			case global.ExpressLocal:
				cacheAddressId = append(cacheAddressId,row.AddressId)
			}

		}

	}
	cacheShopId = utils.RemoveRepeatInt(cacheShopId)
	cacheAddressId = utils.RemoveRepeatInt(cacheAddressId)
	//查询客户信息
	cacheShopMap := make(map[int]map[string]interface{}, 0)


	if len(cacheShopId) > 0 {
		cacheShopObject := make([]models2.Shop, 0)
		e.Orm.Model(&models2.Shop{}).Select("name,id,phone,address,user_name,grade_id").Where("c_id = ? and id in ?",
			userDto.CId, cacheShopId).Find(&cacheShopObject)
		//保存为map
		//应该是获取的用户下单的地址
		for _, k := range cacheShopObject {
			GradeName:=""
			if k.GradeId > 0 {
				var gradeRow models2.GradeVip
				e.Orm.Model(&models2.GradeVip{}).Select("name,id").Where("id = ? and enable = ?",k.GradeId,true).Limit(1).Find(&gradeRow)
				if gradeRow.Id > 0 {
					GradeName = gradeRow.Name
				}
			}
			cacheShopMap[k.Id] = map[string]interface{}{
				"name":  k.Name,
				"phone": k.Phone,
				"id":    k.Id,
				"username":k.UserName,
				"grade_name":GradeName,
			}
		}
	}


	cacheAddressMap := make(map[int]map[string]interface{}, 0)
	if len(cacheAddressId) > 0 {
		cacheAddressObject := make([]models2.DynamicUserAddress, 0)
		e.Orm.Model(&models2.DynamicUserAddress{}).Select("id,address").Where("c_id = ? and id in ?",
			userDto.CId, cacheAddressId).Find(&cacheAddressObject)
		//保存为map
		//应该是获取的用户下单的地址
		for _, k := range cacheAddressObject {
			cacheAddressMap[k.Id] = map[string]interface{}{
				"value":k.Address,
			}
		}
	}
	cacheStoreAddressMap:=make(map[int]map[string]interface{}, 0)
	if len(cacheStoreAddressId) > 0 {
		storeAddressObject := make([]models2.CompanyExpressStore, 0)
		e.Orm.Model(&models2.CompanyExpressStore{}).Select("id,address").Where("c_id = ? and id in ?",
			userDto.CId, cacheStoreAddressId).Find(&storeAddressObject)
		//保存为map
		//应该是获取的用户下单的地址
		for _, k := range storeAddressObject {
			cacheStoreAddressMap[k.Id] = map[string]interface{}{
				"value":k.Address,
			}
		}
	}


	result := make([]map[string]interface{}, 0)
	for _, row := range list {

		//如果支付金额为0
		PayMoney:=row.PayMoney
		if row.PayMoney == 0  && row.PayType < global.PayTypeOnlineWechat{
			if row.DeductionMoney > 0 {
				PayMoney = row.DeductionMoney
			}
		}
		var specCount int64
		e.Orm.Table(splitTableRes.OrderSpecs).Where("order_id = ?", row.OrderId).Count(&specCount)


		r := map[string]interface{}{
			"order_id":   row.OrderId,
			"order_no_id":row.OrderNoId,
			"shop": cacheShopMap[row.ShopId],
			"cycle_place": row.CreatedAt.Format("2006-01-02"), 			//下单周期
			"delivery_time":     row.DeliveryTime.Format("2006-01-02"), 			//配送周期
			"delivery_str": row.DeliveryStr,
			"count":          row.Number,
			"specs_count":specCount,
			"money":         PayMoney,
			"line":row.Line,
			"delivery_type":global.GetExpressCn(row.DeliveryType), //配送类型
			"pay_type":global.GetPayType(row.PayType),//支付类型
			"source_type":global.GetOrderSource(row.SourceType),//订单来源
			"status":         global.OrderStatus(row.Status), //成为DB的订单都是支付成功的订单
			"pay_status":     global.GetOrderPayStatus(row.PayStatus),
			"created_at":     row.CreatedAt,
			"delivery":row.DeliveryType,
		}
		switch row.DeliveryType {
		case global.ExpressStore:
			r["address"] = cacheStoreAddressMap[row.AddressId]
		default:
			r["address"] = cacheAddressMap[row.AddressId]

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

	splitTableRes := business.GetTableName(userDto.CId, e.Orm)
	orderErr := e.Orm.Table(splitTableRes.OrderTable).Where("order_id = ?",orderId).First(&object).Error
	if orderErr != nil && errors.Is(orderErr, gorm.ErrRecordNotFound) {

		e.Error(500, orderErr, "订单不存在")
		return
	}
	if orderErr != nil {
		e.Error(500, orderErr, "订单不存在")
		return
	}
	var shopRow models2.Shop
	e.Orm.Model(&models2.Shop{}).Scopes(actions.PermissionSysUser(shopRow.TableName(),userDto)).Where("id = ? ", object.ShopId).Limit(1).Find(&shopRow)

	result := map[string]interface{}{
		"order_id":       object.Id,
		"created_at":     object.CreatedAt.Format("2006-01-02 15:04:05"),
		"cycle_time":object.CreatedAt.Format("2006-01-02"),
		"delivery_time":     object.DeliveryTime.Format("2006-01-02"),
		"delivery_str":      object.DeliveryStr,
		"pay":            global.GetPayType(object.PayType),
		"pay_status_str": global.GetOrderPayStatus(object.PayStatus),
		"pay_status":     object.PayStatus,
		"shop_name":      shopRow.Name,
		"shop_username":  shopRow.UserName,
		"shop_phone":     shopRow.Phone,
		"shop_address":   shopRow.Address,
		"delivery_type":object.DeliveryType,
		"day":time.Now().Format("2006-01-02"),
		"now":time.Now().Format("2006-01-02 15:04:05"),
		"this_user":userDto.Username,
		//https://weapp.dongchuangyun.com/d1#/'
		"url":fmt.Sprintf("%vd%v#/",config.ExtConfig.H5Url,userDto.CId),
		"desc":object.Desc,
	}
	//如果是同城配送那就获取
	switch object.DeliveryType {
	case global.ExpressLocal:
		var userAddress models2.DynamicUserAddress
		e.Orm.Model(&models2.DynamicUserAddress{}).Scopes(actions.PermissionSysUser(userAddress.TableName(),userDto)).Select("id,address").Where(" id = ?",
			 object.AddressId).Limit(1).Find(&userAddress)
		if userAddress.Id > 0{
			result["address"] = map[string]interface{}{
				"address":userAddress.Address,
			}
		}

	case global.ExpressStore:
		var expressStore models2.CompanyExpressStore
		e.Orm.Model(&models2.CompanyExpressStore{}).Scopes(actions.PermissionSysUser(expressStore.TableName(),userDto)).Select("id,address,name").Where(" id = ?",
			object.AddressId).Limit(1).Find(&expressStore)
		if expressStore.Id > 0{
			result["address"] = map[string]interface{}{
				"name":expressStore.Name,
				"address":expressStore.Address,
			}
		}

	}


	//var orderExtend models.OrderExtend
	//orderExtendTable:=business.OrderExtendTableName(orderTableName)
	//e.Orm.Table(orderExtendTable).Where("order_id = ?",orderId).Limit(1).Find(&orderExtend)

	var driverCnf models.Driver
	e.Orm.Model(&driverCnf).Scopes(actions.PermissionSysUser(driverCnf.TableName(),userDto)).Where("id = ? ",object.DriverId).Limit(1).Find(&driverCnf)
	if driverCnf.Id > 0 {
		result["driver_name"] = driverCnf.Name
		result["driver_phone"] = driverCnf.Phone
	}

	var orderSpecs []models.OrderSpecs


	e.Orm.Table(splitTableRes.OrderSpecs).Where("order_id = ?", orderId).Find(&orderSpecs)

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

func (e Orders)Cycle(c *gin.Context)  {
	req := dto.OrderCyCleReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		Errors
	if err != nil {
		e.Logger.Error(err)
		fmt.Println("err!",err)
		e.Error(500, err, err.Error())
		return
	}
	userDto, err := customUser.GetUserDto(e.Orm, c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}

	result:=make([]map[string]interface{},0)

	//fmt.Println("req!!",req)
	splitTableRes := business.GetTableName(userDto.CId, e.Orm)

	//默认展示最近10条的配送周期
	datalist := make([]models2.OrderCycleCnf, 0)
	var count int64
	e.Orm.Model(&models2.OrderCycleCnf{}).Scopes(
		actions.PermissionSysUser(splitTableRes.OrderCycle,userDto)).Select("delivery_str,create_str,uid").Order(global.OrderTimeKey).Find(&datalist).Limit(-1).Offset(-1).
		Count(&count)
	
	for _,row:=range datalist {
		var value string
		switch req.CyCle {
		case 1:
			value = row.DeliveryStr
		case 2:
			value = row.CreateStr
		default:
			continue
		}
		dd :=map[string]interface{}{
			"value":value,
			"uid":row.Uid,
			//"count":"1",
		}
		result = append(result,dd)
	}

	e.PageOK(result, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
	return
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
	e.Orm.Model(&models2.Shop{}).Scopes(actions.PermissionSysUser(shopObject.TableName(),userDto)).Where("id = ? and enable =? ", req.Shop, true).Limit(1).Find(&shopObject)
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
	//fmt.Println("DeductionAllMoney",DeductionAllMoney)
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





	splitTableRes := business.GetTableName(userDto.CId, e.Orm)
	var DeliveryObject models.CycleTimeConf
	e.Orm.Model(&models2.CycleTimeConf{}).Scopes(actions.PermissionSysUser(DeliveryObject.TableName(),userDto)).Where("id = ? and enable =? ", req.Cycle, true).Limit(1).Find(&DeliveryObject)
	if DeliveryObject.Id == 0 {
		e.Error(500, nil, "时间区间不存在")
		return
	}


	//获取到统一配送的配送UUID
	uid:=service.CheckOrderCyCleCnfIsDb(userDto.CId,splitTableRes.OrderCycle,DeliveryObject,e.Orm)

	var lineObject models2.Line
	e.Orm.Model(&models2.Line{}).Scopes(actions.PermissionSysUser(lineObject.TableName(),userDto)).Where("id = ?  and enable = ?", shopObject.LineId, true).Limit(1).Find(&lineObject)

	lineName := lineObject.Name
	var DriverObject models2.Driver
	e.Orm.Model(&models2.Driver{}).Scopes(actions.PermissionSysUser(DriverObject.TableName(),userDto)).Where("id = ? and enable = ?", lineObject.DriverId, true).Limit(1).Find(&DriverObject)


	//保存商品和规格的一些映射
	goodsCacheList:=make(map[int]service.ValetOrderGoodsRow,0)

	//一个订单下了很多商品
	orderRow := &models.Orders{
		Uid: uid,
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
	orderRow.DeliveryTime = service.CalculateTime(DeliveryObject.GiveDay)
	orderRow.DeliveryStr = DeliveryObject.GiveTime
	orderRow.DeliveryID = DeliveryObject.Id
	orderRow.DeliveryType = global.ExpressLocal

	orderRow.CreateBy = userDto.UserId
	orderRow.Buyer = req.Desc
	orderRow.DriverId = DriverObject.Id

	//代客下单,地址就是客户的默认地址
	var defaultAddress models2.DynamicUserAddress
	e.Orm.Model(&defaultAddress).Scopes(actions.PermissionSysUser(defaultAddress.TableName(),userDto)).Select("id").Where(" is_default = 1 and user_id = ?",shopObject.UserId).Limit(1).Find(&defaultAddress)
	//用户是一定有一个默认地址的
	orderRow.AddressId = defaultAddress.Id

	var orderMoney float64
	var goodsNumber int

	for _, goodsList := range req.Goods {
		//fmt.Println("classId",classId)
		specsOrderId := make([]int, 0)
		for _, spec := range goodsList {
			//如果商品不存在
			var goodsObject models2.Goods
			e.Orm.Model(&models2.Goods{}).Scopes(actions.PermissionSysUser(goodsObject.TableName(),userDto)).Select("id,sale,inventory,name,image").Where("id = ?  and enable = ?", spec.GoodsId,  true).Limit(1).Find(&goodsObject)
			if goodsObject.Id == 0 {
				continue
			}

			//如果下单的次数>库存的值，那就是非法数据 直接跳出
			var goodsSpecs models.GoodsSpecs
			e.Orm.Model(&models.GoodsSpecs{}).Scopes(actions.PermissionSysUser(goodsSpecs.TableName(),userDto)).Where("id = ? ", spec.Id).Limit(1).Find(&goodsSpecs)
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
			txRes:=e.Orm.Table(splitTableRes.OrderSpecs).Create(&specRow)
			if txRes.Error !=nil{
				continue
			}
			//规格减库存 + 销量
			e.Orm.Model(&models.GoodsSpecs{}).Where("id = ? and c_id = s?", spec.Id, userDto.CId).Updates(map[string]interface{}{
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

	//设置价格
	orderMoney = utils.RoundDecimalFlot64(orderMoney)
	orderRow.PayMoney = orderMoney
	orderRow.OrderMoney = orderMoney
	orderRow.GoodsMoney = orderMoney
	orderRow.DeductionMoney = orderMoney
	e.Orm.Table(splitTableRes.OrderTable).Create(&orderRow)

	for goodsId,goodsRow:=range goodsCacheList{
		//商品减库存 + 销量
		var goodsObject models.Goods
		e.Orm.Model(&models.Goods{}).Scopes(actions.PermissionSysUser(goodsObject.TableName(),userDto)).Select("sale,inventory,id").Where("id = ?  and enable = ?", goodsId, userDto.CId, true).Limit(1).Find(&goodsObject)
		if goodsObject.Id == 0 {
			continue
		}
		e.Orm.Model(&models.Goods{}).Scopes(actions.PermissionSysUser(goodsObject.TableName(),userDto)).Where(" id = ?", goodsId).Updates(map[string]interface{}{
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
			"balance":Balance,
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
		row:=models2.ShopCreditLog{
			CId: userDto.CId,
			ShopId: shopObject.Id,
			Number: orderMoney,
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
	splitTableRes := business.GetTableName(userDto.CId, e.Orm)

	switch req.Type {
	case global.OrderToolsActionStatus: //状态更新

		e.Orm.Table(splitTableRes.OrderTable).Where("id = ? and enable = ?", req.Id, true).Updates(map[string]interface{}{
			"status":    req.Status,
			"desc":      req.Desc,
			"update_by": userDto.UserId,
		})
		e.Orm.Table(splitTableRes.OrderSpecs).Where("order_id = ?", req.Id).Updates(map[string]interface{}{
			"status": req.Status,
		})
	case global.OrderToolsActionDelivery: //周期更改
		if req.Delivery > 0 {
			e.Orm.Table(splitTableRes.OrderTable).Where("id = ? and enable = ?", req.Id, true).Updates(map[string]interface{}{
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
	req := dto.CyClePageReq{}
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
	splitTableRes := business.GetTableName(userDto.CId, e.Orm)

	//默认展示最近10条的配送周期
	datalist := make([]models2.OrderCycleCnf, 0)

	e.Orm.Table(splitTableRes.OrderCycle).Scopes(
		cDto.MakeSplitTableCondition(req.GetNeedSearch(),splitTableRes.OrderCycle),
		cDto.Paginate(req.GetPageSize(), req.GetPageIndex()),
		actions.PermissionSysUser(splitTableRes.OrderCycle,userDto)).Model(&models2.OrderCycleCnf{}).Order(global.OrderTimeKey).Find(&datalist)

	//下单周期
	createTime := make([]map[string]interface{}, 0)
	//配送周期
	giveTime := make([]map[string]interface{}, 0)

	//下单周期和配送周期是成对出现的,
	for _, row := range datalist {
		t1 := map[string]interface{}{
			"id": row.Id,
			"color":"",
			"t":  row.Uid,
			"value": row.CreateStr,
		}
		createTime = append(createTime, t1)

		t2 := map[string]interface{}{
			"id":    row.Id,
			"color":"",
			"t":     row.Uid,
			"value": row.DeliveryStr,
		}
		giveTime = append(giveTime, t2)
	}
	if len(createTime) > 1 {
		createTime = append(createTime, map[string]interface{}{
			"color":"#1890ff",
			"t":     "create",
			"value": "查看更多周期列表",
		})
		giveTime = append(giveTime, map[string]interface{}{
			"color":"#1890ff",
			"t":     "give",
			"value": "查看更多周期列表",
		})
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
		//
		_,giveStr:=service.GetOrderCyClyCnf(row)
		m["placing"] = service.GetOrderCreateStr(row)
		m["give"] = giveStr
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
	splitTableRes := business.GetTableName(userDto.CId, e.Orm)
	err = s.Update(splitTableRes.OrderTable, &req, p)
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
	splitTableRes := business.GetTableName(userDto.CId, e.Orm)
	err = s.Remove(splitTableRes.OrderTable, &req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("删除Orders失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "删除成功")
}
