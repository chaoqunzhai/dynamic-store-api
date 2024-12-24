/**
@Author: chaoqun
* @Date: 2024/5/12 10:22
*/
package xlsx_export

import (
	"fmt"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/utils"
	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"
	"time"
)

type GoodsExport struct {
	GoodsName string `json:"goods_name"`
	GoodsId int `json:"goods_id"`
	GoodsLayer int `json:"goods_layer"`
	SpecName string `json:"spec_name"`
	Unit string `json:"unit"`
	Original float64 `json:"original"`
	Price float64 `json:"price"`
	Stock int `json:"stock"`
	State string `json:"state"`
	SerialNumber string `json:"serial_number"`

	Class string `json:"class"`
	Brand string `json:"brand"`
}
func (x *XlsxBaseExport)GoodsExport(cid int,table []GoodsExport) string{
	//导出
	x.NewFile()

	var err error

	SheetName := "Sheet1"
	if err =x.File.SetColWidth(SheetName,"A","A",52.23);err!=nil{
		fmt.Println("err1",err)
	}
	if err =x.File.SetColWidth(SheetName,"B","B",20.23);err!=nil{
		fmt.Println("err1",err)
	}
	if err =x.File.SetColWidth(SheetName,"C","C",16.23);err!=nil{
		fmt.Println("err1",err)
	}
	if err =x.File.SetColWidth(SheetName,"D","D",15.98);err!=nil{
		fmt.Println("err1",err)
	}
	if err =x.File.SetColWidth(SheetName,"E","E",12);err!=nil{
		fmt.Println("err1",err)
	}
	if err =x.File.SetColWidth(SheetName,"F","E",15.98);err!=nil{
		fmt.Println("err1",err)
	}
	if err =x.File.SetColWidth(SheetName,"G","E",15.98);err!=nil{
		fmt.Println("err1",err)
	}
	x.XlsxRowIndex = 1
	x.File.SetRowHeight(SheetName,x.XlsxRowIndex,22.95)
	x.File.SetCellValue(SheetName,fmt.Sprintf("A%v",x.XlsxRowIndex),"名称")
	x.File.SetCellValue(SheetName,fmt.Sprintf("B%v",x.XlsxRowIndex),"规格")
	x.File.SetCellValue(SheetName,fmt.Sprintf("C%v",x.XlsxRowIndex),"单位")
	x.File.SetCellValue(SheetName,fmt.Sprintf("D%v",x.XlsxRowIndex),"分类")
	x.File.SetCellValue(SheetName,fmt.Sprintf("E%v",x.XlsxRowIndex),"品牌")
	x.File.SetCellValue(SheetName,fmt.Sprintf("F%v",x.XlsxRowIndex),"库存")
	x.File.SetCellValue(SheetName,fmt.Sprintf("G%v",x.XlsxRowIndex),"原价")
	x.File.SetCellValue(SheetName,fmt.Sprintf("H%v",x.XlsxRowIndex),"售价")
	x.File.SetCellValue(SheetName,fmt.Sprintf("I%v",x.XlsxRowIndex),"状态")
	x.File.SetCellValue(SheetName,fmt.Sprintf("J%v",x.XlsxRowIndex),"编号")

	x.File.SetCellStyle(SheetName,"A1","J1",x.StyleRowSubtitleId)

	newTable:=x.SortGoodsLayer(table)
	for index,datRow:=range newTable{

		//因为上面有4个是标题
		x.XlsxRowIndex = index + 2


		sliceList :=[]interface{}{
			datRow.GoodsName,datRow.SpecName,datRow.Unit,datRow.Class,datRow.Brand,datRow.Stock,
			datRow.Original,datRow.Price,datRow.State,datRow.SerialNumber,
		}
		for cellIndex,tableValue:=range sliceList{

			cellValue :=XlsxIndexRowMap[cellIndex]
			startCell,_:=excelize.JoinCellName(fmt.Sprintf("%v",cellValue),x.XlsxRowIndex)
			x.File.SetCellValue("Sheet1",startCell,tableValue)

		}

	}


	x.File.SetRowHeight("Sheet1",x.XlsxRowIndex,21.95)


	utils.IsNotExistMkDir("cache_export")
	xlsxName:=fmt.Sprintf("cache_export/商品导出-%v.xlsx",time.Now().Format("2006-01-02 15:04"))
	if saveErr := x.File.SaveAs(xlsxName); saveErr != nil {
		zap.S().Errorf("商品数据 大B%v,导出err%v",cid,saveErr)
		return ""
	}
	return xlsxName

}

