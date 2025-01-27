package xlsx_export

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"go-admin/common/utils"
	"go.uber.org/zap"
	"sort"
	"strconv"
	"strings"
	"time"
)

func (x *XlsxBaseExport)CustomerBindUser(Cid int,deliveryTime string,allGoodsMap map[string][]string) (xlsxPath,xlsxName string)  {
	x.NewFile()
	var err error
	styleId, _ := x.File.NewStyle(&excelize.Style{
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
			WrapText: true,
		},

	})
	SheetName := "Sheet1"
	x.File.SetCellValue(SheetName,"A1","行号")
	x.File.SetCellValue(SheetName,"B1","客户名称")

	goodsIndex:=2 //商品x轴的位置

	xlsxName =fmt.Sprintf("%v-%v-客户明细单.xlsx",deliveryTime,time.Now().Format("2006-01-02 15-04-05"))



	sortKey:= make([]string,0)
	for goodsNameKey :=range allGoodsMap{
		sortKey = append(sortKey,goodsNameKey)
	}

	sort.Strings(sortKey)
	var lastRow string
	var userIndex int


	for _,goodsVal:=range sortKey{
		goodsIndex+=1
		//cell是x轴的字母
		cell,_:=excelize.ColumnNumberToName(goodsIndex)

		xValue := fmt.Sprintf("%v1",cell)


		goodsKeyList := strings.Split(goodsVal,"DEVOPS")
		goodsName :=goodsKeyList[1]


		x.File.SetCellValue(SheetName,xValue,goodsName)

		lastRow = cell //设置最后的一个元素
		userList:=allGoodsMap[goodsVal]
		userIndex=1 //客户x轴

		//设置行号
		for i:=0;i<len(userList);i++{
			indexRow:=fmt.Sprintf("A%v",i + 2)
			x.File.SetCellValue(SheetName,indexRow,i+1)
		}
		for _,userValue:=range userList{
			userNameDat :=strings.Split(userValue,"DEVOPS")

			if len(userNameDat) != 2{continue}
			userIndex +=1
			userName :=userNameDat[0]
			userNumber :=userNameDat[1]
			userRow:=fmt.Sprintf("B%v",userIndex)
			x.File.SetCellValue(SheetName,userRow,userName)
			//设置Y轴的数字
			yValue:=fmt.Sprintf("%v%v",cell,userIndex)
			if userNumber == "0" {
				x.File.SetCellValue(SheetName,yValue,"")
			}else {
				userNumberInt,_:=strconv.Atoi(userNumber)
				x.File.SetCellValue(SheetName,yValue,userNumberInt)
			}
		}

	}


	if err = x.File.SetCellStyle(SheetName,"A1",fmt.Sprintf("%v%v",lastRow,userIndex),styleId);err!=nil{
		fmt.Println("set error!!!",err)
	}
	if err =x.File.SetColWidth(SheetName,"A",fmt.Sprintf("%v",lastRow),14.73);err!=nil{
		fmt.Println("err1",err)
	}
	x.File.SetColWidth(SheetName,"A","A",7)
	x.File.SetColWidth(SheetName,"B","B",30)


	utils.DirNotCreate(fmt.Sprintf("%v",Cid))
	xlsxPath =fmt.Sprintf("%v/%v",Cid,xlsxName)

	if saveErr := x.File.SaveAs(xlsxPath); saveErr != nil {
		zap.S().Errorf("用户关联的商品汇总 大B:%v err%v",Cid,saveErr.Error())
		return
	}
	defer func() {
		if err := x.File.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	return
}
