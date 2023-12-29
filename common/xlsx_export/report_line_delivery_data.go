/**
@Author: chaoqun
* @Date: 2023/12/28 15:37
*/
package xlsx_export

import (
	"fmt"
	"go-admin/app/company/models"
	models2 "go-admin/cmd/migrate/migration/models"
	"go-admin/common/business"
	"go-admin/common/utils"
	"go-admin/global"
	"gorm.io/gorm"
	"time"
)

type ReportDeliveryLineObj struct {
	Orm *gorm.DB
	Dat global.ExportRedisInfo
	RedisKey string

	CycleUid string //配送周期的UID
	UpCloud bool //是否保存云端
}

type LineMapping struct {
	DriverVal string `json:"driver_val"` //司机信息
	LineName string `json:"line_name"` //线路名称
	Data map[int]*SheetRow //xlsx数据 key:小B value:xlsx数据
}
func (e ReportDeliveryLineObj)ReadLineDeliveryDetail() (ResultData map[int]*LineMapping,err error )   {

	//多个路线 进行查询
	ResultData = make(map[int]*LineMapping,0)
	nowTimeObj :=time.Now()
	//根据路线获取这个路线的商品
	var orderList []models.Orders

	splitTableRes := business.GetTableName(e.Dat.CId, e.Orm)
	//查这个配送周期下 路线的
	e.Orm.Table(splitTableRes.OrderTable).Select(
		"order_id,line_id,created_at,shop_id").Where("line_id in ? and uid = ?",e.Dat.LineId,e.CycleUid).Find(&orderList)

	//线路和司机信息
	lineList:=make([]models.Line,0)
	e.Orm.Model(&models.Line{}).Where("id in ?",e.Dat.LineId).Find(&lineList)
	for _,row:=range lineList{
		var DriverObj models.Driver
		e.Orm.Model(&models.Driver{}).Select("id,name,phone").Where("id = ?",row.DriverId).Limit(1).Find(&DriverObj)
		dd :=&LineMapping{
			LineName: row.Name,
			Data: make(map[int]*SheetRow,0),
		}
		if DriverObj.Id > 0 {
			dd.DriverVal = fmt.Sprintf("%v/%v",DriverObj.Name,DriverObj.Phone)
		}else {
			dd.DriverVal = "暂无司机"
		}
		ResultData[row.Id] = dd
	}
	//循环所有的订单 ,不同路线的订单应该是

	siteMap:=make(map[int]*SheetRow,0)

	for index,orderRow:=range orderList{
		//放到一个线路里面
		lineRowsData,lineDbOk:=ResultData[orderRow.LineId]
		if !lineDbOk{
			continue
		}
		//订单选择了多个,存在订单是同一个小B发起的
		sheetRow,ok:=siteMap[orderRow.ShopId]
		if !ok{
			//新的小B在查一次,防止查多次
			var shopRow models2.Shop
			e.Orm.Model(&models2.Shop{}).Where("id = ? ", orderRow.ShopId).Limit(1).Find(&shopRow)
			if shopRow.Id == 0 {
				continue
			}
			sheetRow =&SheetRow{
				OrderA2: fmt.Sprintf("DCY.%v.%v",nowTimeObj.Format("20060102"),index + 1),
				SheetName: shopRow.Name,
				ShopPhone: shopRow.Phone,
				ShopAddress: shopRow.Address,
				ShopUserValue: shopRow.UserName,
				OrderCreateTime: orderRow.CreatedAt.Format("2006-01-02 15:04"),
				DriverVal: lineRowsData.DriverVal, //放司机信息
			}
		}

		//获取订单关联的具体规格
		var orderSpecs []models.OrderSpecs


		e.Orm.Table(splitTableRes.OrderSpecs).Where("order_id = ?", orderRow.OrderId).Find(&orderSpecs)

		specsList := make([]*XlsxTableRow, 0)
		for _, row := range orderSpecs {
			xlsx :=&XlsxTableRow{
				GoodsName: row.GoodsName,
				GoodsSpecs: row.SpecsName,
				Unit: row.Unit,
				Number: row.Number,
				Price: row.Money,
				TotalMoney: utils.RoundDecimalFlot64(row.Money  * float64(row.Number)),
			}
			specsList = append(specsList, xlsx)
		}
		sheetRow.Table = append(sheetRow.Table ,specsList...)

		siteMap[orderRow.ShopId] = sheetRow
		//上面把数据已经放到小B中了
		//需要把小B数据 放到指定的路线中

		lineRowsData.Data = siteMap
		ResultData[orderRow.LineId] = lineRowsData
	}

	for l :=range ResultData{
		sheetRowObject :=ResultData[l]

		for s :=range sheetRowObject.Data {
			//对table的数据进行汇总
			SheetRowVal :=sheetRowObject.Data[s]
			for index,v :=range SheetRowVal.Table{
				v.Id = index + 1
				SheetRowVal.AllNumber+=v.Number
				SheetRowVal.AllMoney = utils.RoundDecimalFlot64(SheetRowVal.AllMoney) + v.TotalMoney
			}
			SheetRowVal.MoneyCn = utils.ConvertNumToCny(SheetRowVal.AllMoney)
			//回传设置到上层
			sheetRowObject.Data[s] = SheetRowVal
		}
		//回传设置到上传
		ResultData[l] = sheetRowObject

	}



	return ResultData, err
}