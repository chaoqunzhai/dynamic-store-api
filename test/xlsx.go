package main

import (
	"fmt"
	"time"

	"github.com/xuri/excelize/v2"
)

func main() {
	fileMap :=map[int]interface{}{
		0:"A",
		1:"B",
		2:"C",
		3:"D",
		4:"E",
		5:"F",
		6:"G",
		7:"H",
	}
	file := excelize.NewFile()
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	styleId1, _ := file.NewStyle(&excelize.Style{
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
	styleId2, _ := file.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:   false,
			Italic: false,
			Family: "微软雅黑",
			Size:   12,
			Color:  "",
		},
		Border:[]excelize.Border{
		},
		Alignment:&excelize.Alignment{
			Vertical: "center",
		},

	})
	styleId3, _ := file.NewStyle(&excelize.Style{
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
	styleId4, _ := file.NewStyle(&excelize.Style{
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


	xlsxName:=fmt.Sprintf("%v-订单数据.xlsx",time.Now().Format("2006-01-02 15-04-05"))
	ListSite:=[]string{"小B1","小B2"}
	//defaultHeight:=100.0
	//{"动创云订货测试超群 订货单"},
	//{"DH.20231223.006 ","下单日期：2023-12-23 16:26","客户名称：韩测试"},
	//{"联系人：韩先生","联系电话：18710876788","收货地址：陕西省西安市新城区西部大道"},
	data :=[][]interface{}{
		{1,"香蕉","斤","箱",100,30,0,"1"},
		{1,"香蕉","斤","箱",100,30,0,"2"},
		{1,"香蕉","斤","箱",100,30,0,"3"},
		{1,"香蕉","斤","箱",100,30,0,"4"},
		{33,"苹果","三十","个",100,30,0,"4"},
	}

	for _,row:=range ListSite {
		var err error
		//创建sheet页面
		// 创建一个工作表
		_, err = file.NewSheet(row)
		if err != nil {
			fmt.Println(err)
			return
		}

		//第一行 合并表头

		file.MergeCell(row,"A1","H1")
		file.SetCellValue(row,"A1","动创云订货测试超群 订货单")
		file.SetCellStyle(row,"A1","H1",styleId1)
		//第二行 订货信息

		file.MergeCell(row,"B2","C2")
		file.MergeCell(row,"D2","H2")


		////第三行 联系信息
		file.MergeCell(row,"B3","C3")
		file.MergeCell(row,"A3","H3")


		//第四行 设置标题
		file.SetCellValue(row,"A4","行号")
		if err =file.SetColWidth(row,"A","A",18.19);err!=nil{
			fmt.Println("err1",err)
		}
		if err =file.SetColWidth(row,"B","B",40.23);err!=nil{
			fmt.Println("err1",err)
		}
		if err =file.SetColWidth(row,"C","C",16.23);err!=nil{
			fmt.Println("err1",err)
		}
		if err =file.SetColWidth(row,"D","D",8.98);err!=nil{
			fmt.Println("err1",err)
		}
		if err =file.SetColWidth(row,"E","E",14.09);err!=nil{
			fmt.Println("err1",err)
		}
		if err =file.SetColWidth(row,"F","F",11.36);err!=nil{
			fmt.Println("err1",err)
		}
		if err =file.SetColWidth(row,"G","G",16.16);err!=nil{
			fmt.Println("err1",err)
		}
		if err =file.SetColWidth(row,"G","G",15.73);err!=nil{
			fmt.Println("err1",err)
		}
		if err =file.SetColWidth(row,"H","H",16.73);err!=nil{
			fmt.Println("err1",err)
		}

		file.SetCellValue(row,"B4","商品名称")
		file.SetCellValue(row,"C4","商品规格")
		file.SetCellValue(row,"D4","单位")
		file.SetCellValue(row,"E4","数量")
		file.SetCellValue(row,"F4","单价")
		file.SetCellValue(row,"G4","小计(元)")
		file.SetCellValue(row,"H4","备注")

		file.SetCellValue(row,"A2","DH.20231223.006")

		file.SetCellValue(row,"B2","下单日期：2023-12-23 16:26")
		file.SetCellValue(row,"D2","客户信息：韩测试  韩先生/18710876788")
		file.SetCellValue(row,"A3","收货地址：陕西省西安市新城区西部大道")
		if err = file.SetCellStyle(row,"A2","D3",styleId2);err!=nil{
			fmt.Println("set error",err)
		}
		if err = file.SetCellStyle(row,"A4","H4",styleId3);err!=nil{
			fmt.Println("set error",err)
		}
		file.SetRowHeight(row,2,21.95)
		file.SetRowHeight(row,3,21.95)
		//每条
		endIndex :=1
		for i,datRow:=range data{

			//因为上面有4个是标题
			index:= i + 5
			start:=fmt.Sprintf("A%v",index)
			end:=fmt.Sprintf("H%v",index)
			for k,v:=range datRow{

				enValue :=fileMap[k]
				startCell,_:=excelize.JoinCellName(fmt.Sprintf("%v",enValue),index)
				file.SetCellValue(row,startCell,v)
				endIndex = index
			}
			file.SetRowHeight(row,index,21.95)
			file.SetCellStyle(row,start,end,styleId4)
		}
		//最后开始
		lastIndexStart :=endIndex + 1
		moneyIndex :=endIndex + 2
		file.SetCellValue(row,fmt.Sprintf("A%v",lastIndexStart),"本页小计:")
		file.SetCellValue(row,fmt.Sprintf("A%v",moneyIndex),"合计:")

		file.SetCellValue(row,fmt.Sprintf("B%v",moneyIndex),"大写: 壹佰零柒圆零分整")
		file.SetCellValue(row,fmt.Sprintf("E%v",moneyIndex),"500")
		file.SetCellValue(row,fmt.Sprintf("G%v",moneyIndex),"￥107.00 ")
		emptyIndex:=moneyIndex + 1
		//插入空行
		file.InsertRows(row,emptyIndex,1)
		//并且合并空行
		file.MergeCell(row,fmt.Sprintf("A%v",emptyIndex),fmt.Sprintf("H%v",emptyIndex))

		file.SetCellValue(row,fmt.Sprintf("A%v",endIndex + 4),"备注:")
		file.SetCellValue(row,fmt.Sprintf("G%v",endIndex + 4),"运费: ¥ 0.00 ")
		//结束
		lastEndIndex :=endIndex + 5
		//最后2行进行合并
		file.MergeCell(row,fmt.Sprintf("A%v",lastEndIndex),fmt.Sprintf("C%v",lastEndIndex))
		file.MergeCell(row,fmt.Sprintf("D%v",lastEndIndex),fmt.Sprintf("H%v",lastEndIndex))
		file.SetCellValue(row,fmt.Sprintf("A%v",lastEndIndex),"操作时间:2023年12月25日13:50:50 ")
		file.SetCellValue(row,fmt.Sprintf("E%v",lastEndIndex),"操作员: 超群 ")
		file.SetCellStyle(row,fmt.Sprintf("A%v",lastIndexStart),
			fmt.Sprintf("H%v",lastEndIndex),styleId4)
	}
	// 根据指定路径保存文件
	file.DeleteSheet("Sheet1")

	if err := file.SaveAs(xlsxName); err != nil {
		fmt.Println(err)
	}
}