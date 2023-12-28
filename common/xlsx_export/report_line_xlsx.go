/**
@Author: chaoqun
* @Date: 2023/12/28 15:37
*/
package xlsx_export

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"
)

func (x *XlsxBaseExport)SetLineXlsxRun(cid int,data map[int]*SheetRow) string {

	//实例化对象
	x.NewFile()
	for _,row:=range data {
		var err error
		//创建sheet页面
		// 创建一个工作表
		_, err = x.File.NewSheet(row.SheetName)
		if err != nil {
			zap.S().Errorf("SetXlsxRun error %v",err)
			continue
		}

		//第一行 合并表头
		x.File.MergeCell(row.SheetName,"A1","H1")
		x.File.SetCellStyle(row.SheetName,"A1","H1",x.StyleTitleId)
		x.File.SetCellValue(row.SheetName,"A1",fmt.Sprintf("%v 订货单",row.SheetName))

		//第二行 司机信息即可
		x.File.MergeCell(row.SheetName,"A2","H2")
		x.File.SetCellStyle(row.SheetName,"A2","H2",x.StyleSubtitleId)
		x.File.SetCellValue(row.SheetName,"A2",fmt.Sprintf("配送司机: %v",row.TitleVal))
		x.File.SetRowHeight(row.SheetName,2,21.95)

		//第三行: 标头
		x.SetRowBackGroundCellValue(3,row.SheetName)

		//第四行开始: 进行商品保存
		for index,datRow:=range row.Table{

			x.XlsxRowIndex = index + 4
			start:=fmt.Sprintf("A%v",x.XlsxRowIndex)
			end:=fmt.Sprintf("H%v",x.XlsxRowIndex)

			//datRow 转数组
			sliceList :=[]interface{}{
				datRow.Id,datRow.GoodsName,datRow.GoodsSpecs,datRow.Unit,datRow.Number,datRow.Price,datRow.TotalMoney,
			}
			for cellIndex,tableValue:=range sliceList{

				cellValue :=XlsxIndexRowMap[cellIndex]
				startCell,_:=excelize.JoinCellName(fmt.Sprintf("%v",cellValue),x.XlsxRowIndex)
				x.File.SetCellValue(row.SheetName,startCell,tableValue)

			}
			x.File.SetRowHeight(row.SheetName,x.XlsxRowIndex,21.95)
			x.File.SetCellStyle(row.SheetName,start,end,x.StyleRowInfoId)
		}

		x.SetTotal(false,row)
	}

	_=x.File.DeleteSheet("Sheet1")


	xlsxName:=fmt.Sprintf("%v-路线导出.xlsx",x.ExportTime)
	if err := x.File.SaveAs(xlsxName); err != nil {
		zap.S().Errorf("路线数据导出 大B%v,错误err%v",cid,err.Error())
		return ""
	}
	//释放文件
	defer func() {
		if err := x.File.Close(); err != nil {
			fmt.Println(err)
			zap.S().Errorf("路线数据导出 大B%v,错误关闭文件句柄失败:%v",cid,err.Error())
		}
	}()

	return xlsxName
}