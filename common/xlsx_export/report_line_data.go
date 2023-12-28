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

	UpCloud bool //是否保存云端
}
type DriverRow struct {
	Value string
	Line string
}

func (e *ReportLineObj)ReadLineDetail() (dat map[int]*SheetRow,err error ){

	//根据路线获取这个路线的商品
	linSheetMap:=make(map[int]*SheetRow,0)
	var orderList []models.Orders

	splitTableRes := business.GetTableName(e.Dat.CId, e.Orm)
	e.Orm.Table(splitTableRes.OrderTable).Select(
		"order_id,line_id,created_at").Where("line_id in ?",e.Dat.LineId).Find(&orderList)


	//线路和司机信息
	lineList:=make([]models.Line,0)
	linCnfMap:=make(map[int]*DriverRow,0)
	e.Orm.Model(&models.Line{}).Where("id in ?",e.Dat.LineId).Find(&lineList)
	for _,row:=range lineList{
		var DriverObj models.Driver
		e.Orm.Model(&models.Driver{}).Select("id,name,phone").Where("id = ?",row.DriverId).Limit(1).Find(&DriverObj)
		dd :=&DriverRow{
			Line: row.Name,
		}
		if DriverObj.Id > 0 {
			dd.Value = fmt.Sprintf("%v/%v",DriverObj.Name,DriverObj.Phone)
		}else {
			dd.Value = "暂无司机"
		}
		linCnfMap[row.Id] = dd

	}

	for _,orderRow:=range orderList{
		//放到一个线路里面
		sheetRow,ok:=linSheetMap[orderRow.LineId]
		lineDbData,lineDbOk:=linCnfMap[orderRow.LineId]
		if !lineDbOk{
			continue
		}
		if !ok{

			sheetRow =&SheetRow{
				SheetName: lineDbData.Line,
				OrderCreateTime: orderRow.CreatedAt.Format("2006-01-02 15:04"),
				TitleVal: lineDbData.Value, //放司机信息
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

		linSheetMap[orderRow.LineId] = sheetRow

	}
	//基于上面汇聚的值 在做一次汇总统计

	for s :=range linSheetMap{
		sheetRowObject :=linSheetMap[s]
		//对table的数据进行汇总
		for index,v :=range linSheetMap[s].Table{
			v.Id = index + 1
			sheetRowObject.AllNumber+=v.Number
			sheetRowObject.AllMoney = utils.RoundDecimalFlot64(sheetRowObject.AllMoney) + v.TotalMoney
		}
		sheetRowObject.MoneyCn = utils.ConvertNumToCny(sheetRowObject.AllMoney)
		linSheetMap[s] = sheetRowObject
	}



	return linSheetMap, err
}


func (e *ReportLineObj)SaveExportXlsx(redisRow global.ExportRedisInfo,sheetData map[int]*SheetRow)  {
	export :=XlsxBaseExport{
		ExportUser: redisRow.ExportUser,
		ExportTime: redisRow.ExportTime,

	}

	FileName := export.SetLineXlsxRun(redisRow.CId,sheetData)

	//上传云端
	if e.UpCloud {
		//上传云端
		buckClient :=qiniu.QinUi{CId: redisRow.CId}
		buckClient.InitClient()
		fileName,_:=buckClient.PostFile(FileName)
		e.FileName = fileName
		defer func() {


			_=os.Remove(FileName)
		}()
	}else {
		e.FileName = FileName
	}



}