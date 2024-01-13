/**
@Author: chaoqun
* @Date: 2024/1/7 10:42
*/
package apis

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	sys "go-admin/app/admin/models"
	"go-admin/app/company/service"
	"go-admin/app/company/service/dto"
	"go-admin/cmd/migrate/migration/models"
	"go-admin/common/actions"
	"go-admin/common/business"
	cDto "go-admin/common/dto"
	customUser "go-admin/common/jwt/user"
	"go-admin/common/utils"
	"go-admin/global"
	"go.uber.org/zap"
	"sort"
	"strings"
	"time"
)
type OrdersRefund struct {
	api.Api
}

func (e OrdersRefund)GetPage(c *gin.Context) {
	req := dto.OrdersRefundPageReq{}
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

	if req.Status == -10 {
		req.Status = 0
	}


	refundList := make([]models.OrderReturn, 0)
	var count int64
	err = e.Orm.Table(splitTableRes.OrderReturn).
		Scopes(
			cDto.MakeSplitTableCondition(req.GetNeedSearch(),splitTableRes.OrderReturn),
			cDto.Paginate(req.GetPageSize(), req.GetPageIndex()),
			actions.Permission(splitTableRes.OrderReturn,p)).Order(global.OrderTimeKey).
		Find(&refundList).Limit(-1).Offset(-1).
		Count(&count).Error
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}
	//统一查询优化
	shopIds:=make([]int,0)
	addressIds:=make([]int,0)
	lineIds:=make([]int,0)
	driverIds:=make([]int,0)
	userIds:=make([]int,0)
	orderIds:=make([]string,0)

	for _,row:=range refundList{
		shopIds = append(shopIds,row.ShopId)
		addressIds = append(addressIds,row.AddressId)
		lineIds = append(lineIds,row.LineId)
		driverIds = append(driverIds,row.DriverId)
		orderIds = append(orderIds,row.OrderId)
		if row.AuditBy > 0 {
			userIds = append(userIds,row.AuditBy)
		}
	}
	orderIds = utils.RemoveRepeatStr(orderIds)
	shopIds = utils.RemoveRepeatInt(shopIds)
	addressIds = utils.RemoveRepeatInt(addressIds)

	lineIds = utils.RemoveRepeatInt(lineIds)
	driverIds = utils.RemoveRepeatInt(driverIds)

	userIds = utils.RemoveRepeatInt(userIds)
	//订单
	var orderList []models.Orders
	e.Orm.Table(splitTableRes.OrderTable).Select("order_id,coupon_money").Where("order_id in ? and c_id = ?",orderIds,userDto.CId).Find(&orderList)

	ordersMap:=make(map[string]models.Orders,0)
	for _,order:=range orderList{
		ordersMap[order.OrderId] = order
	}
	//地址
	var addressList []models.DynamicUserAddress
	e.Orm.Model(&models.DynamicUserAddress{}).Where("id in ? and c_id = ?",addressIds,userDto.CId).Find(&addressList)
	addressMap:=make(map[int]models.DynamicUserAddress,0)
	for _,address:=range addressList{
		addressMap[address.Id] = address
	}
	//商家
	var shopList []models.Shop
	e.Orm.Model(&models.Shop{}).Select("id,name").Where("id in ? and c_id = ?",shopIds,userDto.CId).Find(&shopList)
	shopMap:=make(map[int]models.Shop,0)
	for _,shop:=range shopList{
		shopMap[shop.Id] = shop
	}

	//路线信息
	var lineList []models.Line
	e.Orm.Model(&models.Line{}).Select("id,name").Where("id in ? and c_id = ?",lineIds,userDto.CId).Find(&lineList)
	lineMap:=make(map[int]models.Line,0)
	for _,line:=range lineList{
		lineMap[line.Id] = line
	}
	//司机信息
	var driverList []models.Driver
	e.Orm.Model(&models.Driver{}).Select("id,name").Where("id in ? and c_id = ?",driverIds,userDto.CId).Find(&driverList)
	driverMap:=make(map[int]models.Driver,0)
	for _,d:=range driverList{
		driverMap[d.Id] = d
	}
	//用户信息

	var userList []sys.SysUser
	e.Orm.Model(&sys.SysUser{}).Select("user_id,username").Where("user_id in ? and c_id = ?",userIds,userDto.CId).Find(&userList)
	userMap:=make(map[int]sys.SysUser,0)
	for _,d:=range userList{
		userMap[d.UserId] = d
	}
	result:=make([]dto.RefundDto,0)

	//因为是批量订货软件,需要处理 同一个订单order_id的进行合并
	cacheParentMap:=make(map[string]dto.RefundDto,0)
	sortKey:=make([]string,0)

	for _,row:=range refundList{

		shopObj,ok:=shopMap[row.ShopId]
		if !ok{continue}

		addressObj,addressOk:=addressMap[row.AddressId]

		if !addressOk{continue}


		RefundRow,ok:=cacheParentMap[row.ReturnId]

		RefundGoodsRow:=dto.RefundOrderRow{
			RefundId: row.Id,
			GoodsName: row.GoodsName,
			SpecName: row.SpecsName,
			Price: utils.StringDecimal(row.Price),
			Number: row.Number,
			Edit: row.Number,
			Image: 	business.GetGoodsPathFirst(row.CId,row.Image,global.GoodsPath),
			Unit: row.Unit,
			InNumber: row.Number,
			SourceNumber: row.Source,
		}
		if !ok{
			sortKey = append(sortKey,row.ReturnId)
			RefundRow = dto.RefundDto{
				OrderID: row.OrderId,
				ReturnID: row.ReturnId,
				Reason: row.Reason,
				ShopName: shopObj.Name,
				RefundDeliveryMoney: utils.StringDecimal(row.RefundDeliveryMoney),
				RefundMoney: utils.StringDecimal(row.RefundApplyMoney),
				RefundMoneyType: global.RefundMoneyTypeStr(row.RefundMoneyType),
				Status: row.Status,
				StatusCn: global.GetRefundStatus(row.Status),
				SDesc: row.SDesc,
				CDesc: row.CDesc,
				CreatedAt: row.CreatedAt.Format("2006-01-02 15:04:05"),

			}
			lineObj,lineOk:=lineMap[row.LineId]
			if lineOk {
				RefundRow.Line = lineObj.Name
			}

			driverObj,driverOk:=driverMap[row.DriverId]
			if driverOk {
				RefundRow.Driver = 	driverObj.Name
			}

			if row.AuditBy > 0 {
				userObj,AuditOk:=userMap[row.AuditBy]
				if AuditOk {
					RefundRow.AuditName = userObj.Username
				}
			}
			ordersObj,orderOk:=ordersMap[row.OrderId]
			if orderOk {
				RefundRow.CouponMoney = ordersObj.CouponMoney

			}
			RefundRow.Address = dto.RefundAddress{
				Name: addressObj.Name,
				Address:addressObj.Address,
				Mobile: addressObj.Mobile,
			}

			if !row.RefundTime.IsZero() {
				RefundRow.RefundTime = row.RefundTime.Format("2006-01-02 15:04:05")
			}

			RefundTypeAction:=make([]dto.RefundTypeAction,0)
			switch row.PayType {
			case global.PayTypeBalance:		//余额支付 只能退余额
				RefundTypeAction = append(RefundTypeAction,dto.RefundTypeAction{
					Name: "退款到余额",
					Value: global.RefundMoneyBalance,
				})
			case global.PayTypeCredit:		//如果是授信额支付 只能退授信额
				RefundTypeAction = append(RefundTypeAction,dto.RefundTypeAction{
					Name: "退款到授信额",
					Value: global.RefundMoneyCredit,
				})
			case global.PayTypeOffline: // 只能线下退款
				RefundTypeAction = append(RefundTypeAction,dto.RefundTypeAction{
					Name: "线下退款",
					Value: global.RefundMoneyOffline,
				})
			case global.PayTypeOnlineWechat,global.PayTypeOnlineAli: ////线上支付 可以线下退款 和退款余额
				RefundTypeAction = append(RefundTypeAction,dto.RefundTypeAction{
					Name: "线下退款",
					Value: global.RefundMoneyOffline,
				})
				RefundTypeAction = append(RefundTypeAction,dto.RefundTypeAction{
					Name: "退款到余额",
					Value: global.RefundMoneyBalance,
				})
			}
			RefundRow.RefundTypeAction = RefundTypeAction

		}
		RefundRow.RefundGoods = append(RefundRow.RefundGoods,RefundGoodsRow)
		//只需要进行数量 和商品的叠加
		RefundRow.Number +=row.Number
		RefundRow.RefundTodoMoney += utils.RoundDecimalFlot64(float64(row.Number) * row.Price)
		cacheParentMap[row.ReturnId] = RefundRow
	}
	sort.Slice(sortKey, func(i, j int) bool {
		return sortKey[i] > sortKey[j]
	})
	for _, key := range sortKey {
		result = append(result,cacheParentMap[key])
	}



	e.PageOK(result, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
	return
}

func (e OrdersRefund)Audit(c *gin.Context)  {
	req := dto.RefundAuditReq{}
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


	updateMap:= map[string]interface{}{
		"c_desc":req.CDesc,
		"status":req.Status,
		"refund_time": time.Now(),
		"audit_by":userDto.UserId,
	}
	var refundObjectList []models.OrderReturn

	e.Orm.Table(splitTableRes.OrderReturn).Where("c_id = ? and return_id = ?",userDto.CId,req.RefundOrderId).Find(&refundObjectList)
	//只取第一个即可
	refundFirstObject :=refundObjectList[0]
	//todo:驳回的操作
	//售后订单也修改驳回
	if req.Status == global.RefundOkOverReject {
		e.Orm.Table(splitTableRes.OrderReturn).Where("c_id = ? and return_id = ?",userDto.CId,req.RefundOrderId).Updates(&updateMap)
		//大订单也需要驳回

		e.Orm.Table(splitTableRes.OrderTable).Where("c_id = ? and order_id = ?",userDto.CId,refundFirstObject.OrderId).Updates(map[string]interface{}{
			"after_status":global.RefundOkOverReject,
		})
		e.OK("","审批成功")
		return
	}
	//todo:审核通过的操作
	refundMoney :=req.RefundMoney //退款的金额

	if refundFirstObject.Status != global.RefundDefault{
		e.OK("","售后订单状态,非审批中")
		return
	}

	var orderObject models.Orders
	e.Orm.Table(splitTableRes.OrderTable).Where("c_id = ? and order_id = ?",userDto.CId,refundFirstObject.OrderId).Limit(1).Find(&orderObject)

	if orderObject.Id == 0 {
		e.Error(500,nil,"售后原始订单不存在")
		return
	}
	var count int64
	e.Orm.Table(splitTableRes.OrderSpecs).Where("c_id = ? and order_id = ? and spec_id = ?",userDto.CId,refundFirstObject.OrderId,refundFirstObject.SpecId).Count(&count)

	if count == 0 {
		e.Error(500,nil,"售后原始订单规格配置不存在")
		return
	}
	var openInventory bool
	var Inventory models.InventoryCnf
	e.Orm.Model(&models.InventoryCnf{}).Select("id,enable").Where("c_id = ?",userDto.CId).Limit(1).Find(&Inventory)
	if Inventory.Id == 0 {
		openInventory = false
	}else {
		openInventory = Inventory.Enable
	}

	//获取提交的客户信息
	var shopUserObject sys.SysShopUser
	e.Orm.Model(&shopUserObject).Where("c_id = ? and user_id = ?",userDto.CId,refundFirstObject.CreateBy).Limit(1).Find(&shopUserObject)


	//如果通过
	//金额退回选择的方式中

	//查询客户的小B积分配置
	var shop models.Shop
	e.Orm.Model(&shop).Where("c_id = ? and id = ?",userDto.CId,refundFirstObject.ShopId).Limit(1).Find(&shop)
	if shop.Id == 0 {
		e.Error(500,nil,"售后订单客户不存在")
		return
	}

	//如果有运费也是需要抵扣的

	if refundMoney < 0 {
		e.Error(500,nil,"退货费不足以抵扣,请核对售后单")
		return
	}


	updateShopMap:=make(map[string]interface{},0)
	switch req.RefundMoneyType {
	case global.RefundMoneyBalance:		//退余额
		updateShopMap["balance"] = utils.RoundDecimalFlot64(shop.Balance + refundMoney)


	case global.RefundMoneyCredit:		//退授信额
		updateShopMap["credit"] = utils.RoundDecimalFlot64(shop.Credit + refundMoney)

	case global.RefundMoneyOffline: // 线下退款
		//线下退款就是暂无了
	}

	if len(updateShopMap) > 0 {
		shopRes :=e.Orm.Model(&models.Shop{}).Where("c_id = ? and id = ?",userDto.CId,refundFirstObject.ShopId).Updates(&updateShopMap)
		if shopRes.Error !=nil{
			zap.S().Errorf("售后订单客户积分更新失败,updateShopMap:%v err:%v ",updateShopMap, shopRes.Error)
			e.Error(500,nil,"客户积分增加失败")
			return
		}
	}
	if req.RefundMoney > 0{
		//退成功后在进行记录
		switch req.RefundMoneyType {
		case global.RefundMoneyBalance:		//退余额

			row:=models.ShopBalanceLog{
				CId: userDto.CId,
				ShopId: refundFirstObject.ShopId,
				Money: refundMoney,
				Scene:fmt.Sprintf("用户[%v] 提交售后单,%v审批通过,退回余额:%v",shopUserObject.Username, userDto.Username,refundMoney),
				Action: global.UserNumberAdd, //增加
				Type: global.ScanAdmin,
			}
			row.CreateBy = userDto.UserId
			e.Orm.Create(&row)
		case global.RefundMoneyCredit:		//退授信额
			row:=models.ShopCreditLog{
				CId: userDto.CId,
				ShopId: refundFirstObject.ShopId,
				Number: refundMoney,
				Scene:fmt.Sprintf("用户[%v] 提交售后单,%v审批通过,退回授信额:%v",shopUserObject.Username, userDto.Username,refundMoney),
				Action: global.UserNumberAdd, //增加
				Type: global.ScanAdmin,
			}
			row.CreateBy = userDto.UserId
			e.Orm.Create(&row)

		}
	}

	refundAllNumber :=0 //需要退货的总数
	for _,row:=range refundObjectList {

		var InNumber int //入库数
		var LossNumber int //损耗
		//前段选择的退货数
		refundData,refundOk:=req.RefundData[row.Id]
		if refundOk{  //有这个配置 读取前段配置即可
			InNumber = refundData.InNumber
			LossNumber = refundData.LossNumber
		}else {//如果前段没有这个配置
			InNumber = row.Number
			LossNumber = row.LossNumber
		}

		if InNumber < 0 {
			InNumber = 0
		}
		if LossNumber < 0 {
			LossNumber = 0
		}
		//订单必须存在
		var orderSpecsObject models.OrderSpecs
		e.Orm.Table(splitTableRes.OrderSpecs).Where("c_id = ? and order_id = ? and spec_id = ?",userDto.CId,refundFirstObject.OrderId,row.SpecId).Limit(1).Find(&orderSpecsObject)
		if orderSpecsObject.Id == 0 {continue}


		var goodsObject models.Goods
		e.Orm.Model(&goodsObject).Where("id = ? and c_id = ?",row.GoodsId,userDto.CId).Limit(1).Find(&goodsObject)

		if goodsObject.Id == 0 {
			continue
		}
		//规格库存增加
		var goodsSpecs models.GoodsSpecs
		e.Orm.Model(&goodsSpecs).Where("id = ? and c_id = ?",row.SpecId,userDto.CId).Limit(1).Find(&goodsSpecs)
		if goodsSpecs.Id == 0 {
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
		//如果开启了仓库 数据应该回退到仓库中
		if openInventory{
			var InventoryObj models.Inventory
			e.Orm.Model(&InventoryObj).Where("c_id = ? and goods_id =? and spec_id = ?",userDto.CId,row.GoodsId,row.SpecId).Limit(1).Find(&InventoryObj)
			var thisObj models.Inventory
			var SourceNumber int //原数据

			if InventoryObj.Id == 0 {

				//退回的时候 如果没有 给创建一次数据
				createObj := models.Inventory{
					SpecId: row.SpecId,
					GoodsId: row.GoodsId,
					Stock: InNumber,
					OriginalPrice:float64(goodsSpecs.Original),
				}
				createObj.CId = userDto.CId
				e.Orm.Create(&createObj)
				thisObj = createObj
				SourceNumber = 0
			}else {

				SourceNumber = InventoryObj.Stock
				InventoryObj.Stock +=InNumber
				e.Orm.Model(&InventoryObj).Where("id = ?",InventoryObj.Id).Updates(map[string]interface{}{
					"stock":InventoryObj.Stock,
				})
				//叠加后的数据 最后赋值
				thisObj = InventoryObj
			}
			//fmt.Println("库存ID",thisObj.Id,"原库存",SourceNumber,"现库存",thisObj.Stock,"当前入库价",thisObj.OriginalPrice)
			//增加一条入库记录
			RecordLog:=models.InventoryRecord{
				CId: userDto.CId,
				CreateBy:userDto.Username,
				OrderId: fmt.Sprintf("%v",utils.GenUUID()),
				Action: global.InventoryRefundIn, //入库
				Image: imageVal,
				GoodsId: thisObj.GoodsId,
				GoodsName: goodsObject.Name,
				GoodsSpecName: goodsSpecs.Name,
				SpecId: thisObj.SpecId,
				SourceNumber:SourceNumber, //原库存
				ActionNumber:InNumber, //操作的库存
				CurrentNumber:thisObj.Stock, //那现库存 就是 原库存 + 操作的库存
				OriginalPrice:thisObj.OriginalPrice,
				SourcePrice: thisObj.OriginalPrice,
				Unit:goodsSpecs.Unit,
			}
			e.Orm.Table(splitTableRes.InventoryRecordLog).Create(&RecordLog)
		}else {
			//没有开启仓库,那就操作商品规格即可

			//商品库存增加
			//获取到商品id 和 规格  + 退货的商品

			e.Orm.Model(&goodsObject).Where("id = ? and c_id = ?",row.GoodsId,userDto.CId).Updates(map[string]interface{}{
				"inventory":goodsObject.Inventory + InNumber,
			})


			e.Orm.Model(&goodsSpecs).Where("id = ? and c_id = ?",row.SpecId,userDto.CId).Updates(map[string]interface{}{
				"inventory":goodsSpecs.Inventory + InNumber,
			})

		}

		//订单规格:orderSpecsObject 订单规格数量 - 售后数量
		specsNumber := orderSpecsObject.Number - InNumber

		if specsNumber <=0 {
			specsNumber = 0
		}

		//3、已经退货的规格订单也需要标记
		e.Orm.Table(splitTableRes.OrderSpecs).Where("id = ?",orderSpecsObject.Id).Updates(map[string]interface{}{
			"after_status":global.RefundOk,
			"number":specsNumber,
		})


		//总退货数相加, 需要更新到商品订单中
		refundAllNumber += InNumber

		//4.更新下入库 退货的数量 损耗的数量
		e.Orm.Table(splitTableRes.OrderReturn).Where("id = ?",row.Id).Updates(map[string]interface{}{
			"in_number":InNumber,
			"loss_number":LossNumber,
		})
	}


	//把售后单完结掉
	updateMap["refund_apply_money"] = refundMoney
	updateMap["refund_money_type"] = req.RefundMoneyType

	e.Orm.Table(splitTableRes.OrderReturn).Where("c_id = ? and return_id = ?",userDto.CId,req.RefundOrderId).Updates(&updateMap)


	//订单order:orderObject 订单总数量 - 总售后数量
	orderAllNumber := orderObject.Number -  refundAllNumber
	if orderAllNumber <= 0 {
		orderAllNumber = 0
	}

	//把客户的订单的规格 after_status状态改为已经退货,因为订单可能还要查看详情,需要考虑到整个单子都退情况
	//修改订单状态, 修改订单的数量(原始数量 - 退货数量)

	//1、当一个订单下的规格都减完了 那这个订单就是一个退货的状态,个人中心需要保留这个订单,订单状态为 已退货
	updateOrderMap:=map[string]interface{}{
		"after_status":global.RefundOk,
		"number":orderAllNumber,
	}
	if orderAllNumber == 0 {
		updateOrderMap["status"] = global.OrderStatusReturn //售后处理完毕
	}
	//把订单商品的总金额也减少
	GoodsMoney :=orderObject.GoodsMoney -  refundMoney
	if GoodsMoney <= 0{
		GoodsMoney = 0
	}
	updateOrderMap["goods_money"] = GoodsMoney
	//折扣的价格不管了

	//把商品的订单金额也减少
	OrderMoney :=orderObject.OrderMoney -  refundMoney
	if OrderMoney <= 0{
		OrderMoney = 0
	}
	updateOrderMap["order_money"] = OrderMoney
	//更新订单
	e.Orm.Table(splitTableRes.OrderTable).Where("id = ?",orderObject.Id).Updates(updateOrderMap)



	e.OK("","successful")
	return

}

func (e OrdersRefund)Edit(c *gin.Context)  {
	req := dto.RefundEditReq{}
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
	var refundOrder models.OrderReturn

	e.Orm.Table(splitTableRes.OrderReturn).Where("c_id = ? and return_id = ?",userDto.CId,req.RefundOrderId).Limit(1).Find(&refundOrder)
	if refundOrder.Id == 0 {
		e.Error(500,nil,"无售后订单")
		return
	}

	for _,row:=range req.EditList{
		var orderRowRefund models.OrderReturn
		e.Orm.Table(splitTableRes.OrderReturn).Where("c_id = ? and id = ?",userDto.CId,row.RefundId).Limit(1).Find(&orderRowRefund)
		if orderRowRefund.Id == 0 {
			continue
		}
		//没有编辑的时候 才会更新source记录
		updateMap:=map[string]interface{}{
			"number":row.EditNumber,
			"edit":true,
		}
		if !orderRowRefund.Edit {
			updateMap["source"] = row.SourceNumber
		}
		e.Orm.Table(splitTableRes.OrderReturn).Where("c_id = ? and id = ?",userDto.CId,row.RefundId).Updates(updateMap)
	}


	e.OK("","successful")
	return
}