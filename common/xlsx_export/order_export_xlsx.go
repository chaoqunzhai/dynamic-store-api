/**
@Author: chaoqun
* @Date: 2023/12/26 10:40
*/
package xlsx_export

import (
	"fmt"
	"github.com/go-admin-team/go-admin-core/sdk"
	"github.com/xuri/excelize/v2"
	"go-admin/app/company/models"
	"go-admin/global"
	"go.uber.org/zap"
	"sort"

	"strings"
)

//订单导出的模板

type XlsxBaseExport struct {
	File *excelize.File
	DriverVal string //司机信息
	XlsxRowIndex int //行的索引 只要插入了新的数据 就增1
	ExportUser string //操作人
	ExportTime string //操作时间 也就是录入redis的时间
	Buyer string
	StyleTitleId int `json:"style_title_id"` //标题
	StyleSubtitleId int `json:"style_subtitle_id"` //副标题
	StyleRowSubtitleId int `json:"style_row_subtitle_id"` //行标头样式
	StyleRowInfoId int `json:"style_row_info"` //内容

}
//处理带有特殊字符 导致save xlsx失败
func (x *XlsxBaseExport)ReplaceAllString(originalString string) string  {
	for _,r:=range []string{"/","*","?","[","]",".",","}{

		originalString = strings.Replace(originalString, r, "_", -1)

	}
	return originalString
}
func (x *XlsxBaseExport)SetRowBackGroundCellValue(index int,row string)  {
	var err error
	//设置:长度
	if err =x.File.SetColWidth(row,"A","A",18.2);err!=nil{
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
	x.File.SetRowHeight(row,index,21.95)
	x.File.SetCellValue(row,fmt.Sprintf("A%v",index),"行号")
	x.File.SetCellValue(row,fmt.Sprintf("B%v",index),"商品名称")
	x.File.SetCellValue(row,fmt.Sprintf("C%v",index),"商品规格")
	x.File.SetCellValue(row,fmt.Sprintf("D%v",index),"单位")
	x.File.SetCellValue(row,fmt.Sprintf("E%v",index),"数量")
	x.File.SetCellValue(row,fmt.Sprintf("F%v",index),"单价")
	x.File.SetCellValue(row,fmt.Sprintf("G%v",index),"小计(元)")
	x.File.SetCellValue(row,fmt.Sprintf("H%v",index),"备注")
	x.File.SetCellStyle(row,fmt.Sprintf("A%v",index),fmt.Sprintf("H%v",index),x.StyleRowSubtitleId)

}
func (x *XlsxBaseExport)NewStyle()  {


	StyleTitle, _ := x.File.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:   true,
			Italic: false,
			Family: "微软雅黑",
			Size:   22,
			Color:  "",
		},
		Alignment:&excelize.Alignment{
			Horizontal: "center",
		},

	})
	x.StyleTitleId = StyleTitle


	styleId2, _ := x.File.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:   false,
			Italic: false,
			Family: "微软雅黑",
			Size:   12,
			Color:  "",
		},
		Alignment:&excelize.Alignment{
			Vertical: "center",
		},

	})
	x.StyleSubtitleId = styleId2

	styleId3, _ := x.File.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:   false,
			Italic: false,
			Family: "微软雅黑",
			Size:   16,
			Color:  "",

		},

		Fill: excelize.Fill{
			Type: "pattern",
			Color: []string{"#C0C0C0"},
			Pattern: 1,

		},
		Border:[]excelize.Border{

			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		//白底,背景,深色 15%
	})
	x.StyleRowSubtitleId = styleId3


	styleId4, _ := x.File.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:   false,
			Italic: false,
			Family: "微软雅黑",
			Size:   12,
			Color:  "",
		},
		Border:[]excelize.Border{

			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Alignment:&excelize.Alignment{
			Vertical: "center",
		},

	})
	x.StyleRowInfoId = styleId4

}
func (x *XlsxBaseExport)NewFile()  {
	//实例化对象
	x.File = excelize.NewFile()
	//设置样式
	x.NewStyle()
}
func (x *XlsxBaseExport)SetXlsxRun(cid int,data map[int]*SheetRow) string  {
	//实例化对象
	x.NewFile()
	for _,row:=range data {
		var err error
		//创建sheet页面
		// 创建一个工作表
		row.SheetName = x.ReplaceAllString(row.SheetName)
		_, err = x.File.NewSheet(row.SheetName)
		if err != nil {
			zap.S().Errorf("SetXlsxRun error %v",err)
			continue
		}
		x.Buyer = row.Buyer

		x.SetCell(row.SheetName)

		x.SetSubtitleValue(row)

		x.SetCellRow(row.SheetName,row.Table)

		x.SetTotal(true,row)


	}

	_=x.File.DeleteSheet("Sheet1")


	xlsxName:=fmt.Sprintf("%v-订单导出.xlsx",x.ExportTime)
	if err := x.File.SaveAs(xlsxName); err != nil {
		zap.S().Errorf("配送订单 大B:%v选中数据导出错误,err%v",cid,err.Error())
		return""
	}
	//释放文件
	defer func() {
		if err := x.File.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	return xlsxName

}
//进行合并长度等设置
func (x *XlsxBaseExport)SetCell(row string)  {

	//第一行 合并表头
	x.File.MergeCell(row,"A1","H1")


	//第二行 订货信息 收件人
	x.File.MergeCell(row,"B2","C2")
	x.File.MergeCell(row,"D2","H2")

	//第三行 合并收货地址
	x.File.MergeCell(row,"A3","H3")


	//设置副标题
	x.SetRowBackGroundCellValue(4,row)

	//设置高度
	x.File.SetRowHeight(row,2,21.95)
	x.File.SetRowHeight(row,3,21.95)
}

//设置副标题
func (x *XlsxBaseExport)SetSubtitleValue(sheetRow *SheetRow)  {

	x.File.SetCellValue(sheetRow.SheetName,"A2",sheetRow.OrderA2)

	//x.File.SetCellValue(sheetRow.SheetName,"B2",fmt.Sprintf("下单日期：2023-12-23 16:26"))

	if x.DriverVal !=""{
		x.File.SetCellValue(sheetRow.SheetName,"B2",fmt.Sprintf("配送司机：%v",x.DriverVal))
	}

	userInfo :=fmt.Sprintf("客户信息: %v  %v / %v",sheetRow.SheetName,sheetRow.ShopUserValue,sheetRow.ShopPhone)
	x.File.SetCellValue(sheetRow.SheetName,"D2",userInfo)
	x.File.SetCellValue(sheetRow.SheetName,"A3",fmt.Sprintf("收货地址：%v",sheetRow.ShopAddress))

	x.File.SetCellStyle(sheetRow.SheetName,"A2","D3",x.StyleSubtitleId)
}

//内容合并

func (x *XlsxBaseExport)SetCellRow(row string,table []*XlsxTableRow)  {
	x.File.SetCellValue(row,"A1",fmt.Sprintf("%v 订货单",row))

	x.File.SetCellStyle(row,"A1","H1",x.StyleTitleId)

	sortTable :=x.SortLayer(table)
	for index,datRow:=range sortTable{

		//因为上面有4个是标题
		x.XlsxRowIndex = index + 5
		start:=fmt.Sprintf("A%v",x.XlsxRowIndex)
		end:=fmt.Sprintf("H%v",x.XlsxRowIndex)

		//datRow 转数组
		sliceList :=[]interface{}{
			index + 1,datRow.GoodsName,datRow.GoodsSpecs,datRow.Unit,datRow.Number,datRow.Price,datRow.TotalMoney,
		}
		for cellIndex,tableValue:=range sliceList{

			cellValue :=XlsxIndexRowMap[cellIndex]
			startCell,_:=excelize.JoinCellName(fmt.Sprintf("%v",cellValue),x.XlsxRowIndex)
			x.File.SetCellValue(row,startCell,tableValue)

		}
		x.File.SetRowHeight(row,x.XlsxRowIndex,21.95)
		x.File.SetCellStyle(row,start,end,x.StyleRowInfoId)
	}


}

//设置总计
func (x *XlsxBaseExport)SetTotal(freight bool,sheetRow *SheetRow)  {
	//fmt.Println("设置total.SheetName",sheetRow.SheetName,"sheetRow",sheetRow.ExportTime,"freight",freight)
	//最后开始
	x.XlsxRowIndex += 1
	start:=x.XlsxRowIndex
	x.File.SetCellValue(sheetRow.SheetName,fmt.Sprintf("A%v",x.XlsxRowIndex),"本页小计:")
	x.XlsxRowIndex += 1
	x.File.SetCellValue(sheetRow.SheetName,fmt.Sprintf("A%v",x.XlsxRowIndex),"合计:")

	x.File.SetCellValue(sheetRow.SheetName,fmt.Sprintf("B%v",x.XlsxRowIndex),fmt.Sprintf("大写: %v",sheetRow.MoneyCn))
	x.File.SetCellValue(sheetRow.SheetName,fmt.Sprintf("E%v",x.XlsxRowIndex),sheetRow.AllNumber)
	x.File.SetCellValue(sheetRow.SheetName,fmt.Sprintf("G%v",x.XlsxRowIndex),fmt.Sprintf("￥%v",sheetRow.AllMoney))
	x.XlsxRowIndex += 1
	//插入空行
	x.File.InsertRows(sheetRow.SheetName,x.XlsxRowIndex,1)
	//并且合并空行
	x.File.MergeCell(sheetRow.SheetName,fmt.Sprintf("A%v",x.XlsxRowIndex),fmt.Sprintf("H%v",x.XlsxRowIndex))
	x.XlsxRowIndex += 1


	//合并列 给备注放更多信息
	x.File.MergeCell(sheetRow.SheetName,fmt.Sprintf("B%v",x.XlsxRowIndex),fmt.Sprintf("F%v",x.XlsxRowIndex))
	x.File.SetCellValue(sheetRow.SheetName,fmt.Sprintf("A%v",x.XlsxRowIndex),fmt.Sprintf("备注:%v",x.Buyer))
	if freight {
		x.File.SetCellValue(sheetRow.SheetName,fmt.Sprintf("G%v",x.XlsxRowIndex),"运费:")
	}
	//结束
	x.XlsxRowIndex+=1
	//最后2行进行合并
	x.File.MergeCell(sheetRow.SheetName,fmt.Sprintf("A%v",x.XlsxRowIndex),fmt.Sprintf("C%v",x.XlsxRowIndex))
	x.File.MergeCell(sheetRow.SheetName,fmt.Sprintf("D%v",x.XlsxRowIndex),fmt.Sprintf("H%v",x.XlsxRowIndex))
	x.File.SetCellValue(sheetRow.SheetName,fmt.Sprintf("A%v",x.XlsxRowIndex),fmt.Sprintf("操作时间:%v",x.ExportTime))
	x.File.SetCellValue(sheetRow.SheetName,fmt.Sprintf("E%v",x.XlsxRowIndex),fmt.Sprintf("操作员:%v",x.ExportUser))
	end:=x.XlsxRowIndex
	x.File.SetCellStyle(sheetRow.SheetName,fmt.Sprintf("A%v",start),
		fmt.Sprintf("H%v",end),x.StyleRowInfoId)
}

func (x *XlsxBaseExport)SortLayer( Table []*XlsxTableRow) []*XlsxTableRow  {

	Orm := sdk.Runtime.GetDbByKey("*")

	if Orm == nil{
		return  Table
	}
	var goodsId []int
	for _,datRow:=range Table{
		goodsId = append(goodsId,datRow.GoodsId)
	}
	var goodsList []models.Goods
	Orm.Model(&models.Goods{}).Select("id,layer").Where("id in ?",goodsId).Order(global.OrderLayerKey).Find(&goodsList)

	layerMap:=make(map[int]int,0)
	for _,row:=range goodsList{
		layerMap[row.Id] = row.Layer
	}
	for _,datRow :=range Table{
		layerNum := layerMap[datRow.GoodsId]
		datRow.GoodsLayer = layerNum
	}
	sort.Sort(XlsxTableRowOrderList(Table))

	return  Table
}

func (x *XlsxBaseExport)SortGoodsLayer( Table []GoodsExport) []GoodsExport  {

	Orm := sdk.Runtime.GetDbByKey("*")

	if Orm == nil{
		return  Table
	}
	var goodsId []int
	for _,datRow:=range Table{
		goodsId = append(goodsId,datRow.GoodsId)
	}
	var goodsList []models.Goods
	Orm.Model(&models.Goods{}).Select("id,layer").Where("id in ?",goodsId).Order(global.OrderLayerKey).Find(&goodsList)

	layerMap:=make(map[int]int,0)
	for _,row:=range goodsList{
		layerMap[row.Id] = row.Layer
	}
	for _,datRow :=range Table{
		layerNum := layerMap[datRow.GoodsId]
		datRow.GoodsLayer = layerNum
	}
	sort.Sort(XlsxTableGoodsRowOrderList(Table))

	return  Table
}