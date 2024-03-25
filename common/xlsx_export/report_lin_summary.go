/**
@Author: chaoqun
* @Date: 2024/3/25 14:32
*/
package xlsx_export

import (
	"fmt"
	"go-admin/common/utils"
	"go.uber.org/zap"
	"time"
)

//导出路线的汇总表,直接导出返回,不做异步任务中心
type LineSummaryRow struct {

	Layer int `json:"layer"`
	ShopId int `json:"shop_id"`
	ShopName string `json:"shop_name"`
	ShopAddress string `json:"shop_address"`
	ShopPhone string `json:"phone"`
	OrderCount int `json:"order_count"`
	OrderMoney float64 `json:"order_money"`
}
func (x *XlsxBaseExport)ExportLineSummary(Cid int,lineName string,data map[int]LineSummaryRow) (xlsxPath,xlsxName string)  {

	x.NewFile()
	var err error

	SheetName := "Sheet1"
	if err =x.File.SetColWidth(SheetName,"A","A",20.2);err!=nil{
		fmt.Println("err1",err)
	}
	if err =x.File.SetColWidth(SheetName,"B","B",40.23);err!=nil{
		fmt.Println("err1",err)
	}
	if err =x.File.SetColWidth(SheetName,"C","C",16.23);err!=nil{
		fmt.Println("err1",err)
	}
	if err =x.File.SetColWidth(SheetName,"D","D",15.98);err!=nil{
		fmt.Println("err1",err)
	}
	if err =x.File.SetColWidth(SheetName,"E","E",15.98);err!=nil{
		fmt.Println("err1",err)
	}
	x.XlsxRowIndex = 1
	x.File.SetRowHeight(SheetName,x.XlsxRowIndex,22.95)
	x.File.SetCellValue(SheetName,fmt.Sprintf("A%v",x.XlsxRowIndex),"客户名称")
	x.File.SetCellValue(SheetName,fmt.Sprintf("B%v",x.XlsxRowIndex),"客户地址")
	x.File.SetCellValue(SheetName,fmt.Sprintf("C%v",x.XlsxRowIndex),"客户电话")
	x.File.SetCellValue(SheetName,fmt.Sprintf("D%v",x.XlsxRowIndex),"订货总数")
	x.File.SetCellValue(SheetName,fmt.Sprintf("E%v",x.XlsxRowIndex),"订货总金额")
	//设置标题
	x.File.SetCellStyle(SheetName,"A1","E1",x.StyleRowSubtitleId)
	start:=fmt.Sprintf("A%v",x.XlsxRowIndex + 1)

	for _,row:=range data {
		x.XlsxRowIndex += 1

		x.File.SetCellValue(SheetName,fmt.Sprintf("A%v",x.XlsxRowIndex),row.ShopName)

		x.File.SetCellValue(SheetName,fmt.Sprintf("B%v",x.XlsxRowIndex),row.ShopAddress)

		x.File.SetCellValue(SheetName,fmt.Sprintf("C%v",x.XlsxRowIndex),row.ShopPhone)

		x.File.SetCellValue(SheetName,fmt.Sprintf("D%v",x.XlsxRowIndex),row.OrderCount)

		x.File.SetCellValue(SheetName,fmt.Sprintf("E%v",x.XlsxRowIndex),row.OrderMoney)

		x.File.SetRowHeight(SheetName,x.XlsxRowIndex,21.95)

	}
	end:=fmt.Sprintf("E%v",x.XlsxRowIndex)
	//内容开始到结尾设置一个样式
	x.File.SetCellStyle(SheetName,start,end,x.StyleRowInfoId)

	var saveErr error

	utils.DirNotCreate(fmt.Sprintf("%v",Cid))

	xlsxName =fmt.Sprintf("%v-%v_汇总表.xlsx",time.Now().Format("2006-01-02 15-04-05"),lineName)

	xlsxPath =fmt.Sprintf("%v/%v",Cid,xlsxName)

	if saveErr = x.File.SaveAs(xlsxPath); saveErr != nil {
		zap.S().Errorf("路线汇总表导出 大B:%v err%v",Cid,saveErr.Error())
		return
	}
	defer func() {
		if err := x.File.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	return
}