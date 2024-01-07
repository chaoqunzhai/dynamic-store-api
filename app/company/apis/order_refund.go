/**
@Author: chaoqun
* @Date: 2024/1/7 10:42
*/
package apis

import (
	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"go-admin/app/company/service"
	"go-admin/app/company/service/dto"
	"go-admin/cmd/migrate/migration/models"
	"go-admin/common/actions"
	"go-admin/common/business"
	cDto "go-admin/common/dto"
	customUser "go-admin/common/jwt/user"
	"go-admin/common/utils"
	"go-admin/global"
)
type OrdersRefund struct {
	api.Api
}


func (e OrdersRefund) GetPage(c *gin.Context) {
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
	if req.OverStatus == -10 {
		req.OverStatus = 0
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
	for _,row:=range list{
		shopIds = append(shopIds,row.ShopId)
		addressIds = append(addressIds,row.AddressId)
		lineIds = append(lineIds,row.LineId)
		driverIds = append(driverIds,row.DriverId)
	}
	shopIds = utils.RemoveRepeatInt(shopIds)
	addressIds = utils.RemoveRepeatInt(addressIds)

	lineIds = utils.RemoveRepeatInt(lineIds)
	driverIds = utils.RemoveRepeatInt(driverIds)

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

		rowVal["address"] = addressObj
		rowVal["shop"] = shopObj

		rowVal["refund_money_cn"] = global.RefundMoneyTypeStr(row.RefundMoneyType)
		rowVal["status_cn"] = global.GetRefundStatus(row.Status)
		rowVal["over_status_cn"] = global.GetRefundStatus(row.OverStatus)
		rowVal["id"] = row.Id

		if row.RefundTime.IsZero() {
			rowVal["refund_time"] = nil
		}
		result = append(result,rowVal)
	}

	e.PageOK(result, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
	return
}