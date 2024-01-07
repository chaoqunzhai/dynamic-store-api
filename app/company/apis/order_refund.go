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
	"time"
)
type OrdersRefund struct {
	api.Api
}

type RefundType struct {
	Name string `json:"name"`
	Value int `json:"value"`
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


	list := make([]models.OrderReturn, 0)
	var count int64
	err = e.Orm.Table(splitTableRes.OrderReturn).
		Scopes(
			cDto.MakeSplitTableCondition(req.GetNeedSearch(),splitTableRes.OrderReturn),
			cDto.Paginate(req.GetPageSize(), req.GetPageIndex()),
			actions.Permission(splitTableRes.OrderReturn,p)).Order(global.OrderTimeKey).
		Find(&list).Limit(-1).Offset(-1).
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
	for _,row:=range list{
		shopIds = append(shopIds,row.ShopId)
		addressIds = append(addressIds,row.AddressId)
		lineIds = append(lineIds,row.LineId)
		driverIds = append(driverIds,row.DriverId)
		if row.AuditBy > 0 {
			userIds = append(userIds,row.AuditBy)
		}
	}
	shopIds = utils.RemoveRepeatInt(shopIds)
	addressIds = utils.RemoveRepeatInt(addressIds)

	lineIds = utils.RemoveRepeatInt(lineIds)
	driverIds = utils.RemoveRepeatInt(driverIds)

	userIds = utils.RemoveRepeatInt(userIds)
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
	result:=make([]interface{},0)


	for _,row:=range list{
		rowVal := utils.StructToMap(row)

		shopObj,ok:=shopMap[row.ShopId]
		if !ok{continue}

		addressObj,addressOk:=addressMap[row.AddressId]

		if !addressOk{
			continue}

		lineObj,lineOk:=lineMap[row.LineId]
		if lineOk {
			rowVal["line"] = lineObj.Name
		}

		driverObj,driverOk:=driverMap[row.DriverId]
		if driverOk {
			rowVal["driver"] = 	driverObj.Name
		}


		if row.AuditBy > 0 {
			userObj,AuditOk:=userMap[row.AuditBy]
			if AuditOk {
				rowVal["audit_name"] = 	userObj.Username
			}
		}
		rowVal["address"] = addressObj
		rowVal["shop"] = shopObj

		rowVal["refund_money_cn"] = global.RefundMoneyTypeStr(row.RefundMoneyType)
		rowVal["status_cn"] = global.GetRefundStatus(row.Status)
		rowVal["id"] = row.Id

		if row.RefundTime.IsZero() {
			rowVal["refund_time"] = nil
		}
		//计算的价格
		rowVal["refund_todo_money"] = utils.RoundDecimalFlot64(float64(row.Number) * row.Price)

		RefundTypeAction:=make([]RefundType,0)
		switch row.PayType {
		case global.PayTypeBalance:		//余额支付 只能退余额
			RefundTypeAction = append(RefundTypeAction,RefundType{
				Name: "退款到余额",
				Value: global.RefundMoneyBalance,
			})
		case global.PayTypeCredit:		//如果是授信额支付 只能退授信额
			RefundTypeAction = append(RefundTypeAction,RefundType{
				Name: "退款到授信额",
				Value: global.RefundMoneyCredit,
			})
		case global.PayTypeOffline: // 只能线下退款
			RefundTypeAction = append(RefundTypeAction,RefundType{
				Name: "线下退款",
				Value: global.RefundMoneyOffline,
			})
		case global.PayTypeOnlineWechat,global.PayTypeOnlineAli: ////线上支付 可以线下退款 和退款余额
			RefundTypeAction = append(RefundTypeAction,RefundType{
				Name: "线下退款",
				Value: global.RefundMoneyOffline,
			})
			RefundTypeAction = append(RefundTypeAction,RefundType{
				Name: "退款到余额",
				Value: global.RefundMoneyBalance,
			})
		}
		rowVal["refund_type_action"] = RefundTypeAction
		result = append(result,rowVal)
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
	var refundObject models.OrderReturn
	e.Orm.Table(splitTableRes.OrderReturn).Where("c_id = ? and return_id = ?",userDto.CId,req.RefundId).Limit(1).Find(&refundObject)

	//todo:驳回的操作
	//售后订单也修改驳回
	if req.Status == global.RefundOkOverReject {
		e.Orm.Table(splitTableRes.OrderReturn).Where("c_id = ? and return_id = ?",userDto.CId,req.RefundId).Updates(&updateMap)
		//大订单也需要驳回

		e.Orm.Table(splitTableRes.OrderTable).Where("c_id = ? and order_id = ?",userDto.CId,refundObject.OrderId).Updates(map[string]interface{}{
			"after_status":global.RefundOkOverReject,
		})
		e.OK("","审批成功")
		return
	}
	//todo:审核通过的操作
	refundMoney :=req.RefundMoney //退款的金额

	if refundObject.Status != global.RefundDefault{
		e.OK("","售后订单状态,非审批中")
		return
	}

	var orderObject models.Orders
	e.Orm.Table(splitTableRes.OrderTable).Where("c_id = ? and order_id = ?",userDto.CId,refundObject.OrderId).Limit(1).Find(&orderObject)

	if orderObject.Id == 0 {
		e.Error(500,nil,"售后原始订单不存在")
		return
	}
	//如果产生了,优惠卷的折扣,需要扣除
	if orderObject.CouponMoney > 0 {
		refundMoney -= orderObject.CouponMoney
	}


	//获取提交的客户信息
	var shopUserObject sys.SysShopUser
	e.Orm.Model(&shopUserObject).Where("c_id = ? and user_id = ?",userDto.CId,refundObject.CreateBy).Limit(1).Find(&shopUserObject)


	//如果通过
	//金额退回选择的方式中

	//查询客户的小B积分配置
	var shop models.Shop
	e.Orm.Model(&shop).Where("c_id = ? and id = ?",userDto.CId,refundObject.ShopId).Limit(1).Find(&shop)
	if shop.Id == 0 {
		e.Error(500,nil,"售后订单客户不存在")
		return
	}

	//如果有运费也是需要抵扣的
	if refundObject.RefundDeliveryMoney > 0{
		refundMoney -= refundObject.RefundDeliveryMoney
	}
	//如果有优惠卷 也需要抵扣




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
		shopRes :=e.Orm.Model(&models.Shop{}).Where("c_id = ? and id = ?",userDto.CId,refundObject.ShopId).Updates(&updateShopMap)
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
				ShopId: refundObject.ShopId,
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
				ShopId: refundObject.ShopId,
				Number: refundMoney,
				Scene:fmt.Sprintf("用户[%v] 提交售后单,%v审批通过,退回授信额:%v",shopUserObject.Username, userDto.Username,refundMoney),
				Action: global.UserNumberAdd, //增加
				Type: global.ScanAdmin,
			}
			row.CreateBy = userDto.UserId
			e.Orm.Create(&row)

		}
	}

	//商品库存增加
	//获取到商品id 和 规格  + 退货的商品
	var goodsObject models.Goods
	e.Orm.Model(&goodsObject).Select("id,inventory").Where("id = ? and c_id = ?",refundObject.GoodsId,userDto.CId).Limit(1).Find(&goodsObject)

	if goodsObject.Id > 0{
		e.Orm.Model(&goodsObject).Where("id = ? and c_id = ?",refundObject.GoodsId,userDto.CId).Updates(map[string]interface{}{
			"inventory":goodsObject.Inventory + refundObject.Number,
		})
	}

	//规格库存增加
	var goodsSpecs models.GoodsSpecs
	e.Orm.Model(&goodsSpecs).Select("id,inventory").Where("id = ? and c_id = ?",refundObject.SpecId,userDto.CId).Limit(1).Find(&goodsSpecs)

	if goodsSpecs.Id > 0{
		e.Orm.Model(&goodsSpecs).Where("id = ? and c_id = ?",refundObject.SpecId,userDto.CId).Updates(map[string]interface{}{
			"inventory":goodsSpecs.Inventory + refundObject.Number,
		})
	}

	//把售后单完结掉
	e.Orm.Table(splitTableRes.OrderReturn).Where("c_id = ? and return_id = ?",userDto.CId,req.RefundId).Updates(&updateMap)
	//把客户的订单的规格 after_status状态改为已经退货,因为订单可能还要查看详情,如果整个单子都退了。
	//修改订单状态, 修改订单的数量(原始数量 - 退货数量)
	//1、当一个订单下的规格都减完了 那这个订单就是一个退货的状态,
	//2、如果没有减完 那只需增加一个有退货的标记,
	//3、已经退货的商品也需要标记
	//4、个人中心需要保留这个订单,订单状态为 已退货


	//e.Orm.Table(splitTableRes.OrderTable).Where("c_id = ? and order_id = ?",userDto.CId,req.RefundId)
	e.OK("","successful")
	return

}