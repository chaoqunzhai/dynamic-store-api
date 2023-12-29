/**
@Author: chaoqun
* @Date: 2023/12/28 15:37
*/
package xlsx_export

import (
	"fmt"
	"go.uber.org/zap"
)

func (x *XlsxBaseExport)SetLineDeliveryXlsxRun(cid int,lineName string,data map[int]*SheetRow) string {

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
		x.DriverVal = row.DriverVal
		x.SetCell(row.SheetName)

		x.SetSubtitleValue(row)

		x.SetCellRow(row.SheetName,row.Table)

		x.SetTotal(true,row)

	}

	_=x.File.DeleteSheet("Sheet1")

	xlsxName:=fmt.Sprintf("%v-%v配送表.xlsx",x.ExportTime,lineName)
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