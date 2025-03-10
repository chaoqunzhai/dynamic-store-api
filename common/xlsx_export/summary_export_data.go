/**
@Author: chaoqun
* @Date: 2023/12/27 15:08
*/
package xlsx_export

import (
	"errors"
	"fmt"
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
	e.Orm.Table(splitTableRes.OrderTable).Select("order_id").Where("uid = ? and status in ? and order_money > 0", data.Uid,global.OrderEffEct()).Find(&orderList)
	orderIds:=make([]string,0)
	for _,k:=range orderList{
		orderIds = append(orderIds,k.OrderId)
	}
	orderSpecs:=make([]models.OrderSpecs,0)
	//查下数据 获取规格 在做一次统计
	//orderId 是一一对应的
	e.Orm.Table(splitTableRes.OrderSpecs).Where("order_id in ?",orderIds).Find(&orderSpecs)

	tableRow := make([]*XlsxTableRow,0)
	//进行数据合并。同等的商品规格:总价
	//汇总表 只需要数量进行合并汇总即可
	mergeMap:=make(map[string]*XlsxTableRow,0)

	for _,row:=range orderSpecs{
		key:=fmt.Sprintf("%v_%v",row.GoodsId,row.SpecId)
		TotalMoney:=utils.RoundDecimalFlot64(row.Money  * float64(row.Number))

		xlsx :=&XlsxTableRow{
			Key: key,
			GoodsName: row.GoodsName,
			GoodsId: row.GoodsId,
			GoodsSpecs: row.SpecsName,
			Unit: row.Unit,
			Number: row.Number,
			Price: row.Money,
			TotalMoney: TotalMoney,
		}
		mergeDat,ok:=mergeMap[key]
		if !ok{
			mergeMap[key] = xlsx
			continue
		}else {
			mergeDat.Number +=row.Number
			mergeDat.TotalMoney += TotalMoney

		}
		//合并起来,统一存放到xlsx中
		mergeMap[key] = mergeDat
	}

	xlsxIndex:=0
	for _,xlsx:=range mergeMap{
		xlsxIndex+=1
		xlsx.Id =xlsxIndex
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
