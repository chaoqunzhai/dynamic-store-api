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
	"go.uber.org/zap"
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
	LineData map[int]*SheetRow //xlsx数据  | key:线路 value:xlsx数据、
	DeliveryData  map[int]map[int]*SheetRow //xlsx数据 key:小B value:xlsx数据
}

func (e ReportDeliveryLineObj)ReadLineDeliveryDetail() (ResultData map[int]*LineMapping,err error )   {

	//多个路线 进行查询
	ResultData = make(map[int]*LineMapping,0)
	nowTimeObj :=time.Now()
	//根据路线获取这个路线的商品
	var orderList []models.Orders

	splitTableRes := business.GetTableName(e.Dat.CId, e.Orm)
	//查这个配送周期下 路线的 有效数据
	e.Orm.Table(splitTableRes.OrderTable).Select(
		"order_id,line_id,created_at,shop_id").Where("line_id in ? and uid = ? and status in ? and order_money > 0",
			e.Dat.LineId,e.CycleUid,global.OrderEffEct()).Find(&orderList)

	//线路和司机信息
	lineList:=make([]models.Line,0)
	e.Orm.Model(&models.Line{}).Where("id in ?",e.Dat.LineId).Find(&lineList)
	for _,row:=range lineList{
		var DriverObj models.Driver
		e.Orm.Model(&models.Driver{}).Select("id,name,phone").Where("id = ?",row.DriverId).Limit(1).Find(&DriverObj)
		dd :=&LineMapping{
			LineName: row.Name,
			DeliveryData: make(map[int]map[int]*SheetRow,0),
		}
		if DriverObj.Id > 0 {
			dd.DriverVal = fmt.Sprintf("%v/%v",DriverObj.Name,DriverObj.Phone)
		}else {
			dd.DriverVal = "暂无司机"
		}
		ResultData[row.Id] = dd
	}
	//循环所有的订单 ,不同路线的订单应该是

	//应该是路线下的大B
	siteMap:=make(map[string]map[int]*SheetRow,0)

	for index,orderRow:=range orderList{
		//放到一个线路里面
		lineRowsData,lineDbOk:=ResultData[orderRow.LineId]
		if !lineDbOk{
			continue
		}
		//保持 路线-商家是一一对应的
		KEY :=fmt.Sprintf("%v-%v",orderRow.LineId,orderRow.ShopId)
		//订单关联了商家信息
		sheetMapRow,ok:=siteMap[KEY]
		if !ok{
			//新的小B在查一次,防止查多次
			var shopRow models2.Shop
			e.Orm.Model(&models2.Shop{}).Where("id = ? ", orderRow.ShopId).Limit(1).Find(&shopRow)
			if shopRow.Id == 0 {
				continue
			}
			sheetRow :=&SheetRow{
				LineId: orderRow.LineId,
				OrderA2: fmt.Sprintf("DCY.%v.%v",nowTimeObj.Format("20060102"),index + 1),
				SheetName: shopRow.Name,
				ShopPhone: shopRow.Phone,
				ShopAddress: shopRow.Address,
				ShopUserValue: shopRow.UserName,
				OrderCreateTime: orderRow.CreatedAt.Format("2006-01-02 15:04"),
				DriverVal: lineRowsData.DriverVal, //放司机信息
			}
			sheetMapRow =make(map[int]*SheetRow,0)
			sheetMapRow[orderRow.ShopId] = sheetRow
			siteMap[KEY] = sheetMapRow
		}

		//获取订单关联的具体规格
		var orderSpecs []models.OrderSpecs


		e.Orm.Table(splitTableRes.OrderSpecs).Where("order_id = ?", orderRow.OrderId).Find(&orderSpecs)

		specsList := make([]*XlsxTableRow, 0)
		for _, row := range orderSpecs {
			xlsx :=&XlsxTableRow{
				Key: fmt.Sprintf("%v_%v",row.GoodsId,row.SpecId),
				GoodsName: row.GoodsName,
				GoodsSpecs: row.SpecsName,
				Unit: row.Unit,
				Number: row.Number,
				Price: row.Money,
				TotalMoney: utils.RoundDecimalFlot64(row.Money  * float64(row.Number)),
			}
			specsList = append(specsList, xlsx)
		}
		sheetMapRow[orderRow.ShopId].Table = append(sheetMapRow[orderRow.ShopId].Table ,specsList...)

		//保存商家信息

		siteMap[KEY] = sheetMapRow

		//上面把数据已经放到小B中了
		//需要把多个小B数据 放到指定的路线中

		//获取路线下的 小BMAP
		sieShopMap,ok:=lineRowsData.DeliveryData[orderRow.LineId]
		if !ok{
			sieShopMap = siteMap[KEY]
		}
		//设置订单的商家  = 缓存查询的商家
		sieShopMap[orderRow.ShopId] = sheetMapRow[orderRow.ShopId]
		lineRowsData.DeliveryData[orderRow.LineId] = sieShopMap
		//fmt.Println("线路!!!",orderRow.LineId,"名字",orderRow.Line,"商家",KEY)
		ResultData[orderRow.LineId] = lineRowsData
	}

	for l :=range ResultData{
		sheetRowObject :=ResultData[l]

		SheetDataRowMap,ok := sheetRowObject.DeliveryData[l]
		if !ok{
			zap.S().Errorf("导出配送表时,不在数据Map中,ResultData 和 sheetRowObject.Data 线路数据不匹配")
			continue
		}
		//循环每一个小B
		for sheetShopIndex,sheetShop:=range SheetDataRowMap{
			//对table的数据进行汇总

			mergeMap:=make(map[string]*XlsxTableRow,0)
			for _,xlsxRow :=range sheetShop.Table{

				//对数据进行去重
				mergeDat,mergeOk:=mergeMap[xlsxRow.Key]
				if !mergeOk{
					mergeMap[xlsxRow.Key] = xlsxRow
					continue
				}else {
					mergeDat.Number +=xlsxRow.Number
					mergeDat.TotalMoney += xlsxRow.TotalMoney

				}
				//合并起来,统一存放到xlsx中
				mergeMap[xlsxRow.Key] = mergeDat
			}
			xlsxIndex:=0
			newTable:=make([]*XlsxTableRow,0)
			//进行数据合并
			for _,mergeXlsx:=range mergeMap{
				xlsxIndex+=1
				mergeXlsx.Id = xlsxIndex
				newTable = append(newTable,mergeXlsx)

				sheetShop.AllNumber +=mergeXlsx.Number
				sheetShop.AllMoney = utils.RoundDecimalFlot64(sheetShop.AllMoney) + mergeXlsx.TotalMoney
			}
			sheetShop.Table = newTable
			sheetShop.MoneyCn = utils.ConvertNumToCny(sheetShop.AllMoney)
			//设置回去
			SheetDataRowMap[sheetShopIndex] = sheetShop
		}

		//fmt.Println("小B",SheetRowVal.SheetName,SheetRowVal.AllMoney,SheetRowVal.AllNumber,sheetRowObject.LineName,sheetRowObject)
		//回传设置到上层
		sheetRowObject.DeliveryData[l] = SheetDataRowMap

		//回传设置到上传
		ResultData[l] = sheetRowObject

	}



	return ResultData, err
}