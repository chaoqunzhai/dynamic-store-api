/**
@Author: chaoqun
* @Date: 2024/1/16 14:57
*/
package service

import (
	"errors"
	"fmt"
	sys "go-admin/app/admin/models"
	"go-admin/app/company/models"
	"go-admin/app/company/service/dto"
	models2 "go-admin/cmd/migrate/migration/models"
	"go-admin/common/business"
	"go-admin/common/utils"
	"go-admin/global"
	"strings"
)

func (e *Orders)CancelOrder(RecordAction int,reqAll bool,reqOrderId string,reqOrderSpecId []int,reqDesc string,splitTableRes business.TableRow,userDto *sys.SysUser) error  {

	var orderObject models.Orders
	//获取订单对象
	e.Orm.Table(splitTableRes.OrderTable).Where("c_id = ?  and order_id = ?",userDto.CId,reqOrderId).Limit(1).Find(&orderObject)

	if orderObject.Id == 0  {

		return errors.New("订单不存在")
	}
	if orderObject.AfterStatus == global.OrderStatusCancel {

		return errors.New("订单已作废")
	}

	if orderObject.Status == global.OrderStatusOver {
		return errors.New("订单已验收")
	}
	var shopRow models2.Shop
	e.Orm.Model(&models2.Shop{}).Where("c_id = ? and id = ? ",userDto.CId,orderObject.ShopId).Limit(1).Find(&shopRow)

	if shopRow.Id == 0  {
		return errors.New("客户不存在")
	}

	//查看是否开启了库存
	openInventory := IsOpenInventory(userDto.CId,e.Orm)


	//1.如果整个订单退回, 把支付的order_money 都退回原路

	var returnOrderMoney float64 //退还金额


	refundMap:=make([]dto.OrderRefund,0) //退回的map映射

	//获取这个订单下 所有的规格数据
	var orderSpecsList []models.OrderSpecs

	orm :=e.Orm.Table(splitTableRes.OrderSpecs).Select("order_id,id,goods_id,spec_id,number,after_status,all_money")
	if len(reqOrderSpecId) > 0 {
		orm = orm.Where("c_id = ? and order_id = ? and id in ?",userDto.CId,reqOrderId,reqOrderSpecId)
	}else {
		orm = orm.Where("c_id = ? and order_id = ?",userDto.CId,reqOrderId)
	}
	orm.Find(&orderSpecsList)
	if len(orderSpecsList) == 0 {

		return errors.New("订单无规格数据")
	}
	isAllAfterStatus := true //是否全部已经退回
	specsId:=make([]int,0)
	var returnOrderSpecMoney float64
	for _,row:=range orderSpecsList{
		specsId = append(specsId,row.Id)

		returnOrderSpecMoney +=row.AllMoney
		refundMap =append(refundMap,dto.OrderRefund{
			GoodsId: row.GoodsId,
			Specs: dto.OrderRefundSpec{
				SpecId: row.SpecId,
				Number: row.Number,
				OrderId: row.OrderId,
			},
		})
	}

	if reqAll { //如果是整个订单退回,那就更新整个订单
		//获取整个订单的金额, 用于退款
		returnOrderMoney += utils.RoundDecimalFlot64(orderObject.OrderMoney)

		isAllAfterStatus = true //全部退回订单

	}else {//不是全部退还 那就是用查询到的规格的价格 叠加退还
		returnOrderMoney = returnOrderSpecMoney

		var allOrderSpecs []models.OrderSpecs
		e.Orm.Table(splitTableRes.OrderSpecs).Select("id,after_status").Where("c_id = ? and order_id = ?",userDto.CId,reqOrderId).Find(&allOrderSpecs)
		for _,row:=range allOrderSpecs {
			//查询到规格ID 不是当前要操作的ID, 其他的ID都是为作废了,那这个订单也就是作废了
			if utils.IsArrayInt(row.Id,reqOrderSpecId) { //不检测当前规格ID，因为当前规格ID在进行操作
				continue
			}
			if row.AfterStatus != global.RefundCompanyCancelCType { //只要有一个规格订单 不是大B退回操作 那就不需要修改整个订单
				isAllAfterStatus = false
			}
		}

	}

	if openInventory{ //开启库存,更新库存 + 规格的商品, 如果不存在 那就不管了？

		for _,dat:=range refundMap{


			var goodsObject models.Goods
			e.Orm.Model(&goodsObject).Where("id = ? and c_id = ?",dat.GoodsId,userDto.CId).Limit(1).Find(&goodsObject)

			if goodsObject.Id == 0 {
				continue
			}
			//规格库存增加
			var goodsSpecs models.GoodsSpecs
			e.Orm.Model(&goodsSpecs).Where("goods_id = ? and id = ? and c_id = ?",dat.GoodsId,dat.Specs.SpecId,userDto.CId).Limit(1).Find(&goodsSpecs)
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

			var Inventory models2.Inventory

			e.Orm.Model(&models2.Inventory{}).Select("id,stock,original_price").Where(
				"c_id = ? and goods_id = ? and spec_id = ?",userDto.CId,dat.GoodsId,dat.Specs.SpecId).Limit(1).Find(&Inventory)
			if Inventory.Id == 0 {
				continue //商品没了,那就不操作仓库?
			}
			SourceNumber:=Inventory.Stock
			//只更新 指定的商品和规格的库存
			e.Orm.Model(&models2.Inventory{}).Where("c_id = ? and goods_id = ? and spec_id = ?",userDto.CId,dat.GoodsId,dat.Specs.SpecId).Updates(map[string]interface{}{
				"stock":dat.Specs.Number + Inventory.Stock,
			})
			//并且增加入库记录
			RecordLog:=models2.InventoryRecord{
				CId: userDto.CId,
				CreateBy:userDto.Username,
				OrderId: dat.Specs.OrderId,
				Action: RecordAction,
				Image: imageVal,
				GoodsId: dat.GoodsId,
				GoodsName: goodsObject.Name,
				GoodsSpecName: goodsSpecs.Name,
				Source: 2,//大B发起的
				SpecId: dat.Specs.SpecId,
				SourceNumber:SourceNumber, //原库存
				ActionNumber:dat.Specs.Number, //操作的库存
				CurrentNumber:dat.Specs.Number + Inventory.Stock, //那现库存 就是 原库存 + 操作的库存
				OriginalPrice:Inventory.OriginalPrice,
				SourcePrice: Inventory.OriginalPrice,
				Unit:goodsSpecs.Unit,
			}
			e.Orm.Table(splitTableRes.InventoryRecordLog).Create(&RecordLog)


		}

	}else {//没有开启库存,直接操作商品和规格
		for _,dat:=range refundMap{

			allGoodsNumber :=0


			allGoodsNumber += dat.Specs.Number //订单的数量

			var goodsSpecsObject models.GoodsSpecs
			e.Orm.Model(&models.GoodsSpecs{}).Select("id,inventory").Where("c_id = ? and goods_id = ? and id = ?",userDto.CId,dat.GoodsId,dat.Specs.SpecId).Limit(1).Find(&goodsSpecsObject)
			if goodsSpecsObject.Id == 0 {
				continue //商品没了,那就不操作,并不能影响客户退费
			}
			//只更新商品规格的库存
			e.Orm.Model(&models.GoodsSpecs{}).Where("c_id = ? and goods_id = ? and id = ?",userDto.CId,dat.GoodsId,dat.Specs.SpecId).Updates(map[string]interface{}{
				"inventory":dat.Specs.Number + goodsSpecsObject.Inventory,
			})


			var goodsObject models.Goods
			e.Orm.Model(&models.Goods{}).Select("id,inventory").Where("c_id = ? and id = ?",userDto.CId,dat.GoodsId).Limit(1).Find(&goodsObject)
			if goodsObject.Id == 0 {
				continue //商品没了,那就不操作,并不能影响客户退费
			}
			//只需更新商品的库存
			e.Orm.Model(&models.Goods{}).Where("c_id = ? and id = ?",userDto.CId,dat.GoodsId).Updates(map[string]interface{}{
				"inventory":allGoodsNumber + goodsObject.Inventory,
			})
		}
	}
	var SceneText string
	if RecordAction == global.InventoryCancelIn {
		SceneText  = "作废订单"
	}else  if RecordAction == global.InventoryApproveIn {
		SceneText  = "审批驳回订单"
	}

	updateMoneyMap:=make(map[string]interface{})
	//支付方式
	switch orderObject.PayType {

	case global.PayTypeCredit:
		updateMoneyMap["credit"] = shopRow.Credit + returnOrderMoney
		row:=models2.ShopCreditLog{
			CId: userDto.CId,
			ShopId: shopRow.Id,
			Number: returnOrderMoney,
			Scene:fmt.Sprintf("管理员[%v] %v 退回:%v",userDto.Username,SceneText,returnOrderMoney),
			Action: global.UserNumberAdd, //增加
			Type: global.ScanAdmin,
		}
		e.Orm.Create(&row)
	default:
		//其他支付都默认退回余额
		updateMoneyMap["balance"] = shopRow.Balance + returnOrderMoney
		row:=models2.ShopBalanceLog{
			CId: userDto.CId,
			ShopId: shopRow.Id,
			Money: returnOrderMoney,
			Scene:fmt.Sprintf("管理员[%v] %v 退回:%v",userDto.Username,SceneText,returnOrderMoney),
			Action: global.UserNumberAdd, //增加
			Type: global.ScanShopUse,
		}
		e.Orm.Create(&row)
	}

	updateMasterMap:=make(map[string]interface{},0)
	if isAllAfterStatus { //增加更改主订单的状态
		//订单直接都改为作废
		updateMasterMap["status"] = global.OrderStatusCancel
	}
	updateMasterMap["edit"] = true
	updateMasterMap["approve_msg"] = reqDesc
	updateMasterMap["after_status"] = global.RefundCompanyCancelCType
	updateMasterMap["approve_status"] = global.OrderApproveReject
	//更新主订单的信息
	e.Orm.Table(splitTableRes.OrderTable).Where("c_id = ? and order_id = ?",userDto.CId,reqOrderId).Updates(updateMasterMap) //更新主订单
	//循环查询到规格ID
	for _,row:=range specsId{
		e.Orm.Table(splitTableRes.OrderSpecs).Where("order_id = ? and id = ?",reqOrderId,row).Updates(map[string]interface{}{
			"after_status":global.RefundCompanyCancelCType,
			"status":global.OrderStatusCancel,
			"edit_action":reqDesc,
			"edit":true,
		}) //更新子订单
	}

	e.Orm.Model(&models2.Shop{}).Where("id = ?",shopRow.Id).Updates(updateMoneyMap)

	return nil
}