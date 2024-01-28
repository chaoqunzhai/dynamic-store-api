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
	"go.uber.org/zap"
	"gorm.io/gorm"
	"math"
	"time"

	"go-admin/global"

	"strconv"
	"strings"


	"go-admin/app/company/models"
	"go-admin/app/company/service"
	"go-admin/app/company/service/dto"
	"go-admin/common/actions"
)

type Orders struct {
	api.Api
}


func (e Orders) OrderAction(c *gin.Context) {
	req := dto.OrdersActionReq{}
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

	actionCN:=""
	//更新
	updateMap:=map[string]interface{}{
		"status":req.Action,
	}

	switch req.Action {

	case global.OrderWaitConfirm:
		actionCN = "配送中"
	case global.OrderStatusOver:
		actionCN = "收货完成"
		//审核也通过
		updateMap["approve_status"] = global.OrderApproveOk
	default:
		e.Error(500, nil,"不可识别的操作")
		return
	}
	splitTableRes := business.GetTableName(userDto.CId, e.Orm)

	e.Orm.Table(splitTableRes.OrderTable).Where("c_id = ? and order_id in ?",userDto.CId,req.OrderList).Updates(updateMap)

	e.Orm.Table(splitTableRes.OrderSpecs).Where("c_id = ? and order_id in ?",userDto.CId,req.OrderList).Updates(map[string]interface{}{
		"status":req.Action,
	})
	zap.S().Infof("用户 %v,操作订单 %v 进行 %v,备注:%v",userDto.Username,strings.Join(req.OrderList,","),actionCN,req.Msg)
	

	e.OK("","successful")
	return

}
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

	openApprove,hasApprove:=service.IsHasOpenApprove(userDto,e.Orm)

	//配送周期传入的值是:14_2023-09-23
	//配送周期查询

	//下单周期查询


	list := make([]models.Orders, 0)
	var count int64
	req.CId = userDto.CId
	err = s.GetPage(openApprove,splitTableRes, &req, p, &list, &count)
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
			"pay_money":utils.StringDecimal(row.PayMoney),
			"deduction_money":utils.StringDecimal(row.DeductionMoney),
			"order_money":         utils.StringDecimal(row.OrderMoney),
			"coupon_money":utils.StringDecimal(row.CouponMoney),
			"goods_money":utils.StringDecimal(row.GoodsMoney),
			"delivery_money":utils.StringDecimal(row.DeliveryMoney),
			"line":row.Line,
			"delivery_type":global.GetExpressCn(row.DeliveryType), //配送类型
			"pay_type":global.GetPayType(row.PayType),//支付类型
			"pay_int":row.PayType,
			"source_type":global.GetOrderSource(row.SourceType),//订单来源
			"status_int":row.Status,
			"status":         global.OrderStatus(row.Status), //成为DB的订单都是支付成功的订单
			"pay_status":     global.GetOrderPayStatus(row.PayStatus),
			"created_at":     row.CreatedAt,
			"delivery":row.DeliveryType,
			"approve_status":row.ApproveStatus,
		}
		if openApprove {
			//开启了审核
			if row.ApproveStatus == 0 { //还没审核,这个订单就是审核中
				r["status"] = "待审核"
			}
		}
		if row.ApproveStatus == global.OrderApproveReject {
			r["status"] = "已驳回"
		}
		if row.DeliveryType == global.ExpressStore{

			if row.Status == global.OrderWaitConfirm{
				r["status"] = "待取"
			}
		}
		if row.Edit{
			//订单调整了
			r["goods_money"] = row.OrderMoney
			//调整的时候 又把优惠卷加回去了
			r["money"] = utils.StringDecimal(row.OrderMoney - row.CouponMoney)
		}
		switch row.DeliveryType {
		case global.ExpressStore:
			r["address"] = cacheStoreAddressMap[row.AddressId]
		default:
			r["address"] = cacheAddressMap[row.AddressId]

		}
		result = append(result, r)
	}


	resultData:=map[string]interface{}{
		"list":result,
		"pageIndex": req.GetPageIndex(),
		"pageSize": req.GetPageSize(),
		"count":int(count),
		"hasApprove":hasApprove,
		"openApprove":openApprove,
	}
	e.OK(resultData,"successful")
	return
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
	result,err :=s.DetailOrder(orderId,userDto)
	if err!=nil{
		e.Error(500, err, err.Error())
		return
	}
	e.OK(result, "查询成功")
	return
}


func (e Orders)RichData(c *gin.Context) {
	req := dto.RichOrderDataReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		Errors
	if err != nil {
		e.Logger.Error(err)
		fmt.Println("err!", err)
		e.Error(500, err, err.Error())
		return
	}
	userDto, err := customUser.GetUserDto(e.Orm, c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}
	fmt.Println("userDto",userDto.CId)
	return
}


func (e Orders)Cycle(c *gin.Context)  {
	req := dto.OrderCyCleReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
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

	result:=make([]map[string]interface{},0)

	splitTableRes := business.GetTableName(userDto.CId, e.Orm)

	//默认展示最近10条的配送周期
	datalist := make([]models2.OrderCycleCnf, 0)
	var count int64


	e.Orm.Table(splitTableRes.OrderCycle).Scopes(
		cDto.MakeSplitTableCondition(req.GetNeedSearch(),splitTableRes.OrderCycle),
		actions.PermissionSysUser(splitTableRes.OrderCycle,userDto)).Select("delivery_str,create_str,uid,id,delivery_time").Order(global.OrderTimeKey).Find(&datalist).Limit(-1).Offset(-1).
		Count(&count)
	
	for _,row:=range datalist {
		var value string
		switch req.CyCle {
		case 1:
			value = row.DeliveryStr
		case 2:
			value = fmt.Sprintf("%v %v", row.DeliveryTime.Format("2006-01-02"), row.CreateStr)
		default:
			continue
		}
		dd :=map[string]interface{}{
			"value":value,
			"uid":row.Uid,
			"id":row.Id,
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
	var lineObject models2.Line
	e.Orm.Model(&models2.Line{}).Scopes(actions.PermissionSysUser(lineObject.TableName(),userDto)).Where("id = ?  and enable = ?", shopObject.LineId, true).Limit(1).Find(&lineObject)
	if lineObject.Id == 0 {
		e.Error(500, errors.New("商家路线路线未开启"), "商家路线路线未开启")
		return
	}
	if msg,ExpiredOrNot :=service.CheckLineExpire(userDto.CId,shopObject.LineId,e.Orm);!ExpiredOrNot{
		e.Error(500, errors.New(msg), msg)
		return
	}
	var DeductionAllMoney float64
	switch req.DeductionType {
	case global.DeductionBalance:
		DeductionAllMoney = shopObject.Balance
	case global.DeductionCredit:
		DeductionAllMoney = shopObject.Credit
	case global.DeductionOffline:
		DeductionAllMoney = 0
	default:
		e.OK(-1, "未知抵扣类型")
		return
	}
	//前段的金额
	PayOrderMoney:=utils.RoundDecimalFlot64(req.PayMoney)

	if req.DeductionType != global.DeductionOffline { //非线下支付的时候 才会去校验是否可以抵扣
		if DeductionAllMoney < PayOrderMoney {
			e.OK(-1, fmt.Sprintf("%v 不足！！！",global.GetPayCn(req.DeductionType)))
			return
		}
	}
	//一个订单下了很多商品
	orderRow := &models.Orders{
		Enable:  true,
		ShopId:  req.Shop,
		CId:     userDto.CId,
		SourceType: global.OrderSourceValet, //代客下单
	}
	//扣款方式
	if req.OfflinePayId > 0 {
		orderRow.PayType = global.PayTypeOffline
		orderRow.OfflinePayId = req.OfflinePayId
	}else {
		orderRow.PayType = req.DeductionType
	}
	//是否开启了库存管理
	openInventory:=service.IsOpenInventory(userDto.CId,e.Orm)

	splitTableRes := business.GetTableName(userDto.CId, e.Orm)


	if req.DeliveryType == 2{ //配送
		var DeliveryObject models.CycleTimeConf
		e.Orm.Model(&models2.CycleTimeConf{}).Where("c_id = ? and id = ? and enable =? ",userDto.CId, req.Cycle, true).Limit(1).Find(&DeliveryObject)
		if DeliveryObject.Id == 0 {
			e.Error(500, nil, "时间区间不存在")
			return
		}
		//获取到统一配送的配送UUID
		orderRow.Uid = service.CheckOrderCyCleCnfIsDb(userDto.CId,splitTableRes.OrderCycle,DeliveryObject,e.Orm)
		orderRow.DeliveryTime = service.CalculateTime(DeliveryObject.GiveDay)
		orderRow.DeliveryStr = DeliveryObject.GiveTime
		orderRow.DeliveryID = DeliveryObject.Id
		orderRow.Status = global.OrderStatusWaitSend

		//代客下单,地址就是客户的默认地址
		var defaultAddress models2.DynamicUserAddress
		e.Orm.Model(&defaultAddress).Scopes(actions.PermissionSysUser(defaultAddress.TableName(),userDto)).Select("id").Where(" is_default = 1 and user_id = ?",shopObject.UserId).Limit(1).Find(&defaultAddress)
		//用户是一定有一个默认地址的
		orderRow.AddressId = defaultAddress.Id
	}else { //自提
		orderRow.Status = global.OrderWaitConfirm
		orderRow.DeliveryRunAt = models3.XTime{
			Time:time.Now(),
		}
		orderRow.AddressId = req.StoreAddressId
		orderRow.DeliveryStr = req.DeliveryStr
	}


	orderRow.Line = lineObject.Name
	orderRow.LineId = lineObject.Id

	var DriverObject models2.Driver
	e.Orm.Model(&models2.Driver{}).Scopes(actions.PermissionSysUser(DriverObject.TableName(),userDto)).Where("id = ? and enable = ?", lineObject.DriverId, true).Limit(1).Find(&DriverObject)

	if DriverObject.Id == 0 {
		e.Error(500, errors.New("路线暂无司机"), "路线暂无司机")
		return
	}

	//保存商品和规格的一些映射
	goodsCacheList:=make(map[int]service.ValetOrderGoodsRow,0)


	orderRow.OrderId = fmt.Sprintf("%v",utils.GenUUID())
	//代客下单,需要把配送周期保存，方便周期配送

	orderRow.PayStatus = global.OrderStatusPaySuccess
	orderRow.DeliveryCode = service.DeliveryCode()
	orderRow.PayTime = models3.XTime{
		Time:time.Now(),
	}
	//代客下单，是不需要审批的,
	orderRow.ApproveStatus =  global.OrderApproveOk
	orderRow.DeliveryType = req.DeliveryType
	orderRow.Phone = shopObject.Phone
	//代客下单时的用户是管理员用户！！！
	orderRow.CreateBy = userDto.UserId
	orderRow.Buyer = req.Desc
	orderRow.DriverId = DriverObject.Id

	var goodsNumber int

	for key, selectSpec := range req.Goods {
		cardKey:=strings.Split(key,"_")
		if len(cardKey) != 2{
			continue
		}
		GoodsId,_:=strconv.Atoi(string(cardKey[0]))
		SpecId,_:=strconv.Atoi(string(cardKey[1]))

		specsOrderId := make([]int, 0)

		//如果商品不存在
		var goodsObject models2.Goods
		e.Orm.Model(&models2.Goods{}).Scopes(
			actions.PermissionSysUser(
				goodsObject.TableName(),userDto)).Select(
					"id,sale,inventory,name,image").Where(
						"id = ?  and enable = ?", GoodsId,  true).Limit(1).Find(&goodsObject)
		if goodsObject.Id == 0 {
			continue
		}

		//如果下单的次数>库存的值，那就是非法数据 直接跳出
		var goodsSpecs models.GoodsSpecs
		e.Orm.Model(&models.GoodsSpecs{}).Scopes(
			actions.PermissionSysUser(
				goodsSpecs.TableName(),userDto)).Where(
					"id = ? ", SpecId).Limit(1).Find(&goodsSpecs)
		if goodsSpecs.Id == 0 {
			continue
		}
		var goodsSpecsStock int
		var InventoryObject models2.Inventory
		if openInventory{ //开启了库存管理
			e.Orm.Model(&models2.Inventory{}).Where(
				"c_id = ? and goods_id = ? and spec_id = ?",
				userDto.CId,GoodsId,SpecId).Limit(1).Find(&InventoryObject)

			goodsSpecsStock = InventoryObject.Stock
		}else {
			goodsSpecsStock = goodsSpecs.Inventory
		}
		if selectSpec.Number > goodsSpecsStock {
			continue
		}

		Money:=utils.RoundDecimalFlot64(selectSpec.CachePrice)
		specRow := &models.OrderSpecs{
			OrderId: orderRow.OrderId,
			Number:    selectSpec.Number,
			Money:    Money,
			Unit:      selectSpec.UnitName,
			GoodsName: goodsObject.Name,
			GoodsId: goodsObject.Id,
			SpecsName: goodsSpecs.Name,
			SpecId:   SpecId,
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
		if openInventory {
			if InventoryObject.Id > 0 { //有库存
				CurrentNumber :=InventoryObject.Stock - selectSpec.Number
				e.Orm.Model(&models2.Inventory{}).Where(
					"c_id = ? and goods_id = ? and spec_id = ?",userDto.CId,GoodsId,SpecId).Updates(map[string]interface{}{
					"stock":CurrentNumber,
				})

				//同时需要做出入库记录
				RecordLog:=models2.InventoryRecord{
					CId: userDto.CId,
					CreateBy:userDto.Username,
					OrderId: orderRow.OrderId,
					Action: global.InventoryHelpOut,
					Source: 2, //大B发起
					GoodsId: GoodsId,
					GoodsName: goodsObject.Name,
					GoodsSpecName: goodsSpecs.Name,
					SpecId: SpecId,
					SourceNumber:InventoryObject.Stock, //原库存
					ActionNumber:selectSpec.Number, //操作的库存
					CurrentNumber:CurrentNumber, //那现库存
					OriginalPrice:InventoryObject.OriginalPrice,
					SourcePrice: InventoryObject.OriginalPrice,
					Unit:selectSpec.UnitName,
				}
				e.Orm.Table(splitTableRes.InventoryRecordLog).Create(&RecordLog)
			}

		}else {
			//规格减库存 + 销量
			e.Orm.Model(&models.GoodsSpecs{}).Where("id = ? and c_id = ?", SpecId, userDto.CId).Updates(map[string]interface{}{
				"inventory": goodsSpecs.Inventory - selectSpec.Number,
				"sale":	goodsSpecs.Sale + selectSpec.Number,
			})
		}

		goodsNumber += selectSpec.Number

		//缓存订单ID
		specsOrderId = append(specsOrderId, specRow.Id)


		cacheValetOrderRow,ok:=goodsCacheList[GoodsId]
		if !ok{
			goodsCacheList[GoodsId] = service.ValetOrderGoodsRow{
				Number: selectSpec.Number,
			}
		}else{
			cacheValetOrderRow.Number +=selectSpec.Number
			goodsCacheList[GoodsId] = cacheValetOrderRow
		}

	}
	orderRow.Number = goodsNumber

	//设置价格

	DiscountMoney:=utils.RoundDecimalFlot64(req.DiscountMoney)

	var PayOkMoney float64 //实扣金额

	PayOkMoney = utils.RoundDecimalFlot64(PayOrderMoney - DiscountMoney)

	if DiscountMoney > PayOrderMoney{ //优惠金额大于实扣金额时
		PayOkMoney = 0
		DiscountMoney = PayOrderMoney
	}

	orderRow.PayMoney = PayOkMoney  //支付金额
	orderRow.OrderMoney = PayOrderMoney //订单金额
	orderRow.GoodsMoney = utils.RoundDecimalFlot64(req.GoodsMoney) //商品金额
	orderRow.DeductionMoney = PayOkMoney //抵扣金额 因为不是实际的付款,也是要存抵扣金额的
	orderRow.CouponMoney = DiscountMoney //优惠的金额 在一个优惠卷字段来存储
	e.Orm.Table(splitTableRes.OrderTable).Create(&orderRow)


	if !openInventory {//没有开启库存,那就是最后直接更新商品的总库存即可
		for goodsId,goodsRow:=range goodsCacheList{
			//商品减库存 + 销量
			var goodsObject models.Goods
			e.Orm.Model(&models.Goods{}).Scopes(actions.PermissionSysUser(goodsObject.TableName(),userDto)).Select("sale,inventory,id").Where("id = ?  and enable = ?", goodsId, userDto.CId, true).Limit(1).Find(&goodsObject)
			if goodsObject.Id == 0 {
				continue
			}
			e.Orm.Model(&models.Goods{}).Scopes(actions.PermissionSysUser(goodsObject.TableName(),userDto)).Where(" id = ?", goodsId).Updates(map[string]interface{}{
				"sale":    goodsObject.Sale + goodsRow.Number,
				"inventory": goodsObject.Inventory - goodsRow.Number,
			})
		}
	}
	//授信额减免

	switch req.DeductionType {
	case global.DeductionBalance:

		Balance:= shopObject.Balance - PayOkMoney
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
			Money: PayOkMoney,
			Scene:fmt.Sprintf("管理员[%v] 代客下单,抵扣:%v",userDto.Username,PayOkMoney),
			Action: global.UserNumberReduce, //抵扣
			Type: global.ScanShopUse,
		}
		e.Orm.Create(&row)
	case global.DeductionCredit:
		Credit:=  shopObject.Credit - PayOkMoney
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
			Number: PayOkMoney,
			Scene:fmt.Sprintf("管理员[%v] 代客下单,抵扣:%v",userDto.Username,PayOkMoney),
			Action: global.UserNumberReduce, //抵扣
			Type: global.ScanShopUse,
		}
		e.Orm.Create(&row)
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
		actions.PermissionSysUser(splitTableRes.OrderCycle,userDto)).Order(global.OrderTimeKey).Find(&datalist)

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
			"value":fmt.Sprintf("%v %v", row.DeliveryTime.Format("2006-01-02"), row.CreateStr),
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

func (e Orders) EditOrder(c *gin.Context) {
	req := dto.OrdersEditReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	orderId:=c.Param("orderId")

	userDto, err := customUser.GetUserDto(e.Orm, c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}
	isOpenInventory:=service.IsOpenInventory(userDto.CId,e.Orm)
	//订单编辑: 库存检测  是否需要进行库存检测

	//分表配置
	splitTableRes := business.GetTableName(userDto.CId, e.Orm)

	var orderObject models.Orders

	orderErr := e.Orm.Table(splitTableRes.OrderTable).Where("order_id = ?",orderId).Limit(1).Find(&orderObject).Error
	if orderErr != nil && errors.Is(orderErr, gorm.ErrRecordNotFound) {
		e.Error(500, nil,"订单不存在")
		return
	}
	if orderErr != nil {
		e.Error(500, nil,"订单不存在")
		return
	}

	var shopRow models2.Shop
	e.Orm.Model(&models2.Shop{}).Where("id = ? ", orderObject.ShopId).Limit(1).Find(&shopRow)

	var CompareAmount float64
	switch req.Deduction {
	case global.PayTypeBalance:

		CompareAmount = shopRow.Balance
	case global.PayTypeCredit:

		CompareAmount = shopRow.Credit
	default:
		CompareAmount = -1
	}
	if CompareAmount > 0 {
		if CompareAmount < req.Money{
			e.Error(500, nil,"不足以进行费用抵扣")

			return
		}
	}


	//如果反复的进行对订单操作,
	//1.只记录一条记录


	sourceOrderNumber := orderObject.Number
	sourceOrderMoney := orderObject.OrderMoney //订单金额进行操作

	//fmt.Printf("原订单 数量:%v 金额:%v\n",sourceOrderNumber,sourceOrderMoney)

	//增加映射map
	editGoodsMap:=make(map[string]int,0)
	RecordOrderMap:=make(map[string]string,0)
	for _,order:=range req.EditList{

		var orderSpecs models2.OrderSpecs
		//是否已经创建记录
		e.Orm.Table(splitTableRes.OrderSpecs).Where("order_id = ? and id = ? and c_id = ?",orderId,order.Id,userDto.CId).Limit(1).Find(&orderSpecs)
		if orderSpecs.Id == 0 {
			continue
		}
		if orderSpecs.AllMoney == 0 {
			orderSpecs.AllMoney = utils.RoundDecimalFlot64(orderSpecs.Money  * float64(orderSpecs.Number))
		}
		//只要一直操作就会一直记录
		editRow:=&models2.OrderEdit{
			CreateBy: userDto.UserId,
			OrderId: orderId,
			SpecId: orderSpecs.Id,
			SourerMoney: orderSpecs.AllMoney,
			SourerNumber: orderSpecs.Number,
			Number: order.NewAllNumber,
			Money: order.NewAllMoney,
			Desc: req.Desc,
		}
		editRow.CId = userDto.CId
		e.Orm.Table(splitTableRes.OrderEdit).Create(&editRow)

		//同时修改规格的订单
		e.Orm.Table(splitTableRes.OrderSpecs).Where("id = ?",orderSpecs.Id).Updates(map[string]interface{}{
			"edit":true, //修改了
			"all_money":order.NewAllMoney, //规格用自己新的价格
			"number":order.NewAllNumber, //变更后的数量
		})
		//总的大订单Order也是需要进行重新计价
		//总价只需要记录多余 还是少于即可,因为总价里面有 优惠卷各种抵扣
		//新数据 - 原数据
		sourceOrderNumber +=  order.NewAllNumber -  orderSpecs.Number //订单商品 进行修改

		mapKey :=fmt.Sprintf("%v_%v",orderSpecs.GoodsId,orderSpecs.SpecId)
		editGoodsMap[mapKey] = orderSpecs.Number - order.NewAllNumber //对库存 进行修改
		RecordOrderMap[mapKey] = orderSpecs.OrderId
		sourceOrderMoney  +=  order.NewAllMoney - orderSpecs.AllMoney

		//fmt.Printf("新的订单 数量:%v 金额:%v\n",sourceOrderNumber,sourceOrderMoney)
		if sourceOrderNumber < 0 {
			sourceOrderNumber = 0
		}
		if sourceOrderMoney < 0 {
			sourceOrderMoney = 0
		}
	}
	//把优惠卷的价格也加上, 因为原来的价格是 抛去优惠卷算回来的
	sourceOrderMoney +=orderObject.CouponMoney

	var ActionMode string
	var Scene string
	var EditAction string
	if req.Reduce { //减少数量 那就是返钱
		ActionMode = global.UserNumberAdd
		Scene = fmt.Sprintf("订单编辑退回 %v",math.Abs(req.Money))
		EditAction = "退回"
	}
	if req.Increase { //新增数量 那就是需要额外扣钱
		ActionMode = global.UserNumberReduce
		Scene = fmt.Sprintf("订单编辑抵扣 %v",math.Abs(req.Money))
		EditAction = "抵扣"
	}

	updateMap:=make(map[string]interface{},0)
	//操作余额的时候 也是需要进行记录
	switch req.Deduction {
	case global.PayTypeBalance:
		EditAction = "余额" + EditAction
		shopRow.Balance -=req.Money
		updateMap["balance"] = shopRow.Balance
		row:=models2.ShopBalanceLog{
			CId: userDto.CId,
			ShopId: shopRow.Id,
			Desc: req.Desc,
			Money: math.Abs(req.Money),
			Scene:fmt.Sprintf("管理员[%v] %v",userDto.Username,Scene),
			Action: ActionMode,
			Type: global.ScanAdmin,
		}
		row.CreateBy = user.GetUserId(c)
		e.Orm.Create(&row)
	case global.PayTypeCredit:
		EditAction = "授信额" + EditAction
		shopRow.Credit -=req.Money
		updateMap["credit"] = shopRow.Credit
		row:=models2.ShopCreditLog{
			CId: userDto.CId,
			ShopId: shopRow.Id,
			Desc: req.Desc,
			Number: math.Abs(req.Money),
			Scene:fmt.Sprintf("管理员[%v] %v",userDto.Username,Scene),
			Action: ActionMode,
			Type: global.ScanAdmin,
		}
		row.CreateBy = user.GetUserId(c)
		e.Orm.Create(&row)
	}
	//fmt.Println("小Bid",shopRow.Id,"订单ID",orderObject.Id,"价格",sourceOrderMoney,sourceOrderNumber)
	e.Orm.Table(splitTableRes.OrderTable).Where("id = ?",orderObject.Id).Updates(map[string]interface{}{
		"order_money":sourceOrderMoney,
		"number":sourceOrderNumber,
		"edit":true,
		"edit_action":EditAction,
	})
	e.Orm.Model(&models2.Shop{}).Where("id = ?",shopRow.Id).Updates(&updateMap)

	//操作库存

	for key,ActionNumber:=range editGoodsMap{
		keyData:=strings.Split(key,"_")
		if len(keyData) != 2 {
			continue
		}
		goodsId := keyData[0]
		specId:=keyData[1]
		//更新规格数量
		var goodsSpecs models2.GoodsSpecs
		e.Orm.Model(&models2.GoodsSpecs{}).Where("c_id = ? and goods_id = ? and id = ?",userDto.CId,goodsId,specId).Limit(1).Find(&goodsSpecs)
		if goodsSpecs.Id == 0 {
			continue
		}

		var goodsObject models2.Goods
		e.Orm.Model(&models2.Goods{}).Where("c_id = ? and id = ?",userDto.CId,goodsId).Limit(1).Find(&goodsObject)
		if goodsObject.Id == 0 {
			continue
		}
		imageVal := goodsSpecs.Image
		if goodsSpecs.Image == ""{
			//商品如果有图片,那获取第一张图片即可
			if goodsObject.Image != ""{
				imageVal = strings.Split( goodsObject.Image,",")[0]
			}else {
				imageVal = ""
			}

		}
		if isOpenInventory{
			var Inventory models2.Inventory
			e.Orm.Model(&models2.Inventory{}).Select("id,stock,original_price").Where("c_id = ? and goods_id = ? and spec_id = ?",userDto.CId,goodsId,specId).Limit(1).Find(&Inventory)
			if Inventory.Id == 0 {
				continue
			}
			SourceNumber := Inventory.Stock
			Inventory.Stock +=ActionNumber
			e.Orm.Model(&models2.Inventory{}).Where("c_id = ? and goods_id = ? and spec_id = ?",userDto.CId,goodsId,specId).Updates(map[string]interface{}{
				"stock":Inventory.Stock,
			})

			RecordLog:=models2.InventoryRecord{
				CId: userDto.CId,
				CreateBy:userDto.Username,
				OrderId: RecordOrderMap[key],
				Action: global.InventoryOut, //入库
				Image: imageVal,
				GoodsId: goodsObject.Id,
				GoodsName: goodsObject.Name,
				GoodsSpecName: goodsSpecs.Name,
				SpecId: goodsSpecs.Id,
				SourceNumber:SourceNumber, //原库存
				CurrentNumber:SourceNumber + ActionNumber, //那现库存 就是 原库存 + 操作的库存
				SourcePrice:Inventory.OriginalPrice,
				OriginalPrice:Inventory.OriginalPrice,
				Unit:service.GetUnitName(userDto.CId,goodsSpecs.UnitId,e.Orm),

			}
			//流水创建
			if ActionNumber > 0 {
				RecordLog.Action = global.InventoryEditIn
				RecordLog.ActionNumber = ActionNumber
			}else { //出库
				RecordLog.Action = global.InventoryEditOut
				RecordLog.ActionNumber = -ActionNumber
			}
			e.Orm.Table(splitTableRes.InventoryRecordLog).Create(&RecordLog)



		}else {
			//规格总数
			goodsSpecs.Inventory +=ActionNumber
			e.Orm.Model(&models2.GoodsSpecs{}).Where("c_id = ? and goods_id = ? and id = ?",userDto.CId,goodsId,specId).Updates(map[string]interface{}{
				"inventory":goodsSpecs.Inventory,
			})
			//更新商品总数
			goodsObject.Inventory +=ActionNumber
			e.Orm.Model(&models2.Goods{}).Where("c_id = ? and id = ?",userDto.CId,goodsId).Updates(map[string]interface{}{
				"inventory":goodsObject.Inventory,
			})
		}
	}


	e.OK("","更新成功")
	return

}

func (e Orders)BatchCancelOrder(c *gin.Context) {
	req := dto.ApproveReq{}
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
	splitTableRes := business.GetTableName(userDto.CId, e.Orm)

	OrderSpecId:=make([]int,0)
	//进行作废了 需要退库 退钱
	for _,orderId:=range req.OrderList{
		if cancelErr :=s.CancelOrder(global.InventoryCancelIn,true,orderId,OrderSpecId,req.Desc,splitTableRes,userDto);cancelErr!=nil{
			continue
		}
	}

	e.OK("","successful")
	return

}
func (e Orders)CancelOrder(c *gin.Context)  {
	req := dto.OrdersRefundReq{}
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

	splitTableRes := business.GetTableName(userDto.CId, e.Orm)
	//作废
	if cancelErr :=s.CancelOrder(global.InventoryCancelIn,req.All,req.OrderId,req.OrderSpecId,req.Desc,splitTableRes,userDto);cancelErr!=nil{
		e.Error(500, cancelErr, cancelErr.Error())
		return
	}


	e.OK("","操作成功")
	return

}