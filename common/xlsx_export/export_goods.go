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
	SpecName string `json:"spec_name"`
	Unit string `json:"unit"`
	Original float64 `json:"original"`
	Price float64 `json:"price"`
	Stock int `json:"stock"`
	State string `json:"state"`
	SerialNumber string `json:"serial_number"`
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
	x.File.SetCellValue(SheetName,fmt.Sprintf("A%v",x.XlsxRowIndex),"商品名称")
	x.File.SetCellValue(SheetName,fmt.Sprintf("B%v",x.XlsxRowIndex),"商品规格")
	x.File.SetCellValue(SheetName,fmt.Sprintf("C%v",x.XlsxRowIndex),"商品库存")
	x.File.SetCellValue(SheetName,fmt.Sprintf("D%v",x.XlsxRowIndex),"商品原价")
	x.File.SetCellValue(SheetName,fmt.Sprintf("E%v",x.XlsxRowIndex),"商品售价")
	x.File.SetCellValue(SheetName,fmt.Sprintf("F%v",x.XlsxRowIndex),"商品状态")
	x.File.SetCellValue(SheetName,fmt.Sprintf("G%v",x.XlsxRowIndex),"商品编号")

	x.File.SetCellStyle(SheetName,"A1","G1",x.StyleRowSubtitleId)
	for index,datRow:=range table{

		//因为上面有4个是标题
		x.XlsxRowIndex = index + 2


		sliceList :=[]interface{}{
			datRow.GoodsName,datRow.SpecName,datRow.Stock,
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