/**
@Author: chaoqun
* @Date: 2023/12/28 15:37
*/
package xlsx_export

import (
	"fmt"
	"go-admin/app/company/models"
	"go-admin/common/business"
	"go-admin/common/qiniu"
	"go-admin/common/utils"
	"go-admin/global"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"os"
)

//配送报表中 路线的导出
//获取路线下商品规格的详细列表


type ReportLineObj struct {
	Orm *gorm.DB
	Dat global.ExportRedisInfo
	RedisKey string
	FileName string
	CycleUid string //配送周期的UID
	UpCloud bool //是否保存云端
}

func (e *ReportLineObj)ReadLineDetail() (ResultData map[int]*LineMapping,err error ){

	//根据路线获取这个路线的商品
	//多个路线 进行查询
	ResultData = make(map[int]*LineMapping,0)

	var orderList []models.Orders

	splitTableRes := business.GetTableName(e.Dat.CId, e.Orm)
	//查这个配送周期下 路线的
	e.Orm.Table(splitTableRes.OrderTable).Select(
		"order_id,line_id,created_at").Where("line_id in ? and uid = ?",e.Dat.LineId,e.CycleUid).Find(&orderList)


	//线路和司机信息
	lineList:=make([]models.Line,0)
	e.Orm.Model(&models.Line{}).Where("id in ?",e.Dat.LineId).Find(&lineList)
	for _,row:=range lineList{
		var DriverObj models.Driver
		e.Orm.Model(&models.Driver{}).Select("id,name,phone").Where("id = ?",row.DriverId).Limit(1).Find(&DriverObj)
		dd :=&LineMapping{
			LineName: row.Name,
			LineData: make(map[int]*SheetRow,0),
		}
		if DriverObj.Id > 0 {
			dd.DriverVal = fmt.Sprintf("%v/%v",DriverObj.Name,DriverObj.Phone)
		}else {
			dd.DriverVal = "暂无司机"
		}
		ResultData[row.Id] = dd
	}

	linSheetMap:=make(map[int]*SheetRow,0)
	for _,orderRow:=range orderList{

		//放到一个线路里面
		lineRowsData,lineDbOk:=ResultData[orderRow.LineId]
		if !lineDbOk{
			continue
		}
		sheetRow,ok:=linSheetMap[orderRow.LineId]
		if !ok{
			sheetRow =&SheetRow{
				SheetName: lineRowsData.LineName,
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
		//重新赋值
		linSheetMap[orderRow.LineId] = sheetRow

		//只把自己的路线 放到里面
		lineRowsData.LineData = map[int]*SheetRow{
			orderRow.LineId:linSheetMap[orderRow.LineId],
		}

		ResultData[orderRow.LineId] = lineRowsData

	}
	//基于上面汇聚的值 在做一次汇总统计

	for l :=range ResultData{
		sheetRowObject :=ResultData[l]

		SheetRowVal,ok := sheetRowObject.LineData[l]

		if !ok{
			zap.S().Errorf("导出配送路线时,不在数据Map中,ResultData 和 sheetRowObject.Data 线路数据不匹配")
			continue
		}
		//对table的数据进行汇总
		for index,v :=range SheetRowVal.Table{
			v.Id = index + 1
			SheetRowVal.AllNumber+=v.Number
			SheetRowVal.AllMoney = utils.RoundDecimalFlot64(SheetRowVal.AllMoney) + v.TotalMoney
		}
		SheetRowVal.MoneyCn = utils.ConvertNumToCny(SheetRowVal.AllMoney)
		//fmt.Println("线路",SheetRowVal.SheetName,SheetRowVal.AllMoney,SheetRowVal.AllNumber,sheetRowObject.LineName,sheetRowObject)
		//回传设置到上层
		sheetRowObject.LineData[l] = SheetRowVal

		//回传设置到上传
		ResultData[l] = sheetRowObject

	}



	return ResultData, err
}

//当一条路线的时候,就保留时间文件名

//如果多个路线时候,压缩包是文件名, 里面的excel是路线名称
//sheetData map[int]*SheetRow

//多个线路时 就需要做一个压缩包


func SaveLineExportXlsx(xlsxType,zipFile string,UpCloud bool,redisRow global.ExportRedisInfo,lineSheetData  map[int]*LineMapping) string {
	export :=XlsxBaseExport{
		ExportUser: redisRow.ExportUser,
		ExportTime: redisRow.ExportTime,


	}
	zipList:=make([]string,0)
	for lineId:=range lineSheetData {

		lineName :=lineSheetData[lineId].LineName

		var FileName string

		switch xlsxType {

		case "line":
			sheetMapData,ok:=lineSheetData[lineId]
			if !ok{
				zap.S().Error("警告！路线数据保存失败,线路ID不存在")
				continue
			}
			FileName = export.SetLineXlsxRun(redisRow.CId,lineName,sheetMapData.LineData)
		case "delivery":
			sheetShopMapData,ok:=lineSheetData[lineId]
			if !ok{
				zap.S().Error("警告！路线配送数据保存失败,线路ID不存在")
				continue
			}
			sheetShop,ok:=sheetShopMapData.DeliveryData[lineId]
			if !ok{
				zap.S().Error("警告！路线配送数据保存失败,小B数据线路ID不存在")
				continue
			}
			FileName = export.SetLineDeliveryXlsxRun(redisRow.CId,lineName,sheetShop)
		default:

			continue
		}

		zipList = append(zipList,FileName)

	}

	var FileName string
	//就一个文件 不做zip压缩包
	if len(zipList) == 1{
		FileName = zipList[0]
	}else {
		//多个文件 压缩包
		FileName,_ = utils.ZipFile(zipFile,zipList)
	}
	var err error
	//上传云端
	if UpCloud {
		//上传云端
		buckClient :=qiniu.QinUi{CId: redisRow.CId}
		buckClient.InitClient()

		if FileName,err =buckClient.PostFile(FileName);err!=nil{
			zap.S().Errorf("SaveLineExportXlsx 文件:%v 保存云端失败 %v",FileName,err)
		}

		defer func() {

			_=os.Remove(FileName)
			if len(zipList) >  1{ //是压缩包 那清理下文件
				for _,f:=range zipList{
					os.Remove(f)
				}
			}
		}()
	}
	return FileName



}