/**
@Author: chaoqun
* @Date: 2023/12/27 15:08
*/
package xlsx_export

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"
)

func (x *XlsxBaseExport)SetSummaryXlsxRun(cid int,data *SheetRow) string {
	var err error
	row:="Sheet1"
	//实例化对象
	x.NewFile()

	//第一行 合并表头
	x.File.MergeCell(row,"A1","H1")
	x.File.SetCellStyle(row,"A1","H1",x.StyleTitleId)
	x.File.SetCellValue(row,"A1",fmt.Sprintf("%v 汇总表",data.TitleVal))

	//第二行 表头
	x.File.SetCellValue(row,"A2","行号")
	x.File.SetCellValue(row,"B2","商品名称")
	x.File.SetCellValue(row,"C2","商品规格")
	x.File.SetCellValue(row,"D2","单位")
	x.File.SetCellValue(row,"E2","数量")
	x.File.SetCellValue(row,"F2","单价")
	x.File.SetCellValue(row,"G2","小计(元)")
	x.File.SetCellValue(row,"H2","备注")
	x.File.SetRowHeight(row,2,21.95)
	x.File.SetCellStyle(row,"A2","H2",x.StyleRowSubtitleId)
	if err =x.File.SetColWidth(row,"A","A",10.2);err!=nil{
		fmt.Println("err1",err)
	}
	if err =x.File.SetColWidth(row,"B","B",40.23);err!=nil{
		fmt.Println("err1",err)
	}
	if err =x.File.SetColWidth(row,"C","C",16.23);err!=nil{
		fmt.Println("err1",err)
	}
	if err =x.File.SetColWidth(row,"D","D",8.98);err!=nil{
		fmt.Println("err1",err)
	}
	if err =x.File.SetColWidth(row,"E","E",14.09);err!=nil{
		fmt.Println("err1",err)
	}
	if err =x.File.SetColWidth(row,"F","F",11.36);err!=nil{
		fmt.Println("err1",err)
	}
	if err =x.File.SetColWidth(row,"G","G",16.16);err!=nil{
		fmt.Println("err1",err)
	}
	if err =x.File.SetColWidth(row,"G","G",15.73);err!=nil{
		fmt.Println("err1",err)
	}
	if err =x.File.SetColWidth(row,"H","H",16.73);err!=nil{
		fmt.Println("err1",err)
	}
	//第三行 数据开始

	for index,datRow:=range data.Table{

		//因为上面有只有2个属性配置也就需要 从第三行开始
		x.XlsxRowIndex = index + 3
		start:=fmt.Sprintf("A%v",x.XlsxRowIndex)
		end:=fmt.Sprintf("H%v",x.XlsxRowIndex)

		//datRow 转数组
		sliceList :=[]interface{}{
			datRow.Id,datRow.GoodsName,datRow.GoodsSpecs,datRow.Unit,datRow.Number,datRow.Price,datRow.TotalMoney,
		}
		for cellIndex,tableValue:=range sliceList{

			cellValue :=XlsxIndexRowMap[cellIndex]
			startCell,_:=excelize.JoinCellName(fmt.Sprintf("%v",cellValue),x.XlsxRowIndex)
			x.File.SetCellValue(row,startCell,tableValue)

		}
		x.File.SetRowHeight(row,x.XlsxRowIndex,21.95)
		x.File.SetCellStyle(row,start,end,x.StyleRowInfoId)
	}

	x.SetTotal(false,data)

	xlsxName:=fmt.Sprintf("%v-配送周期导出.xlsx",data.ExportTime)
	if err = x.File.SaveAs(xlsxName); err != nil {
		zap.S().Errorf("配送周期导出 大B:%v选中数据导出错误,err%v",cid,err.Error())
		return ""
	}
	//释放文件
	defer func() {
		if err = x.File.Close(); err != nil {
			fmt.Println(err)
			zap.S().Errorf("配送周期导出 大B:%v选中数据导出错误,关闭文件句柄失败:%v",cid,err.Error())
		}
	}()

	return xlsxName
}