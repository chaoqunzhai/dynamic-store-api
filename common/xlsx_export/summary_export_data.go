/**
@Author: chaoqun
* @Date: 2023/12/27 15:08
*/
package xlsx_export

import (
	"errors"
	"go-admin/app/company/models"
	"go-admin/common/business"
	"go-admin/common/qiniu"
	"go-admin/common/utils"
	"go-admin/global"
	"gorm.io/gorm"
	"os"
)

type SummaryExportObj struct {
	Orm *gorm.DB
	Dat global.ExportRedisInfo
	RedisKey string
	FileName string

	UpCloud bool //是否保存云端
}

func (e *SummaryExportObj)ReadSummaryDetail() (dat *SheetRow,err error )  {
	dat = &SheetRow{
		Table: make([]*XlsxTableRow,0),
		SheetName:"Sheet1",//汇总表就一个分页

	}

	splitTableRes := business.GetTableName(e.Dat.CId, e.Orm)

	var data models.OrderCycleCnf
	e.Orm.Table(splitTableRes.OrderCycle).Select("uid,id,delivery_str").Model(
		&models.OrderCycleCnf{}).Where("id = ?",e.Dat.Cycle).Limit(1).Find(&data)
	if data.Id == 0 {
		return nil,errors.New("暂无周期")
	}
	//
	dat.TitleVal = data.DeliveryStr

	orderList:=make([]models.Orders,0)
	//根据配送UID 统一查一下
	e.Orm.Table(splitTableRes.OrderTable).Select("id").Where("uid = ? and status in ? and order_money > 0", data.Uid,global.OrderEffEct()).Find(&orderList)
	orderIds:=make([]int,0)
	for _,k:=range orderList{
		orderIds = append(orderIds,k.Id)
	}
	orderSpecs:=make([]models.OrderSpecs,0)
	//查下数据 获取规格 在做一次统计
	e.Orm.Table(splitTableRes.OrderSpecs).Where("id in ?",orderIds).Find(&orderSpecs)

	tableRow := make([]*XlsxTableRow,0)
	for index,row:=range orderSpecs{
		index +=1
		xlsx :=&XlsxTableRow{
			GoodsName: row.GoodsName,
			GoodsSpecs: row.SpecsName,
			Unit: row.Unit,
			Number: row.Number,
			Price: row.Money,
			Id: index,
			TotalMoney: utils.RoundDecimalFlot64(row.Money  * float64(row.Number)),
		}
		dat.AllNumber +=xlsx.Number
		dat.AllMoney = utils.RoundDecimalFlot64(dat.AllMoney) + xlsx.TotalMoney
		tableRow = append(tableRow, xlsx)
	}
	dat.MoneyCn = utils.ConvertNumToCny(dat.AllMoney)
	dat.Table = tableRow

	return dat, err


}
func (e *SummaryExportObj)SaveExportXlsx(redisRow global.ExportRedisInfo,sheetData *SheetRow)  {
	//fmt.Println("保存到本地zip文件",e.Dat)

	export :=XlsxBaseExport{
		ExportUser: redisRow.ExportUser,
		ExportTime: redisRow.ExportTime,

	}
	// 大B/order_export/2023年12月26日15:46:06.xlsx
	FileName := export.SetSummaryXlsxRun(redisRow.CId,sheetData)

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


	return
}
