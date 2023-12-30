/**
@Author: chaoqun
* @Date: 2023/12/25 18:56
*/
package xlsx_export

//订单选中导出
import (
	"fmt"
	"go-admin/cmd/migrate/migration/models"
	"go-admin/common/business"
	"go-admin/common/qiniu"
	"go-admin/common/utils"
	"go-admin/global"
	"gorm.io/gorm"
	"os"
	"time"
)

type OrderExportObj struct {
	Orm *gorm.DB
	Dat global.ExportRedisInfo
	RedisKey string
	FileName string

	UpCloud bool //是否保存云端
}


type SheetRow struct {
	LineId int
	OrderA2 string //索引标记而已
	SheetName string //分页名称 小B名称 或者路线名称
	TitleVal string //标题的内容
	DriverVal string //司机信息
	OrderCreateTime string //订单创建时间
	ShopAddress string //小B地址
	ShopPhone string //联系电话
	ShopUserValue string //联系人信息
	DeliveryMoney string //运费
	AllNumber int `json:"all_number"`//总数
	AllMoney float64 `json:"all_money"` //总价
	MoneyCn string `json:"money_cn"` //总价的中文
	Table []*XlsxTableRow `json:"table"`

}
type XlsxTableRow struct {
	Id int `json:"id"`
	GoodsName string `json:"goods_name"`
	GoodsSpecs string `json:"goods_specs"`
	Unit string `json:"unit"`
	Number int `json:"number"`
	Price float64 `json:"price"`
	TotalMoney float64  `json:"total_money"`
}

//开始执行数据导出
//多个订单save为一个excel文件
//多个订单 同一个小B放到一个sheet中

func (e *OrderExportObj)ReadOrderDetail() (dat map[int]*SheetRow,err error )  {
	//多个订单 同一个小B的数据放在一起

	siteMap:=make(map[int]*SheetRow,0)
	//
	nowTimeObj :=time.Now()
	var orderList []models.Orders

	splitTableRes := business.GetTableName(e.Dat.CId, e.Orm)
	e.Orm.Table(splitTableRes.OrderTable).Select(
		"order_id,shop_id,created_at").Where("order_id in ?",e.Dat.Order).Find(&orderList)


	//基于订单做一次聚会查询
	for index,orderRow:=range orderList{

		//订单选择了多个,存在订单是同一个小B发起的
		sheetRow,ok:=siteMap[orderRow.ShopId]
		if !ok{
			//新的小B在查一次,防止查多次
			var shopRow models.Shop
			e.Orm.Model(&models.Shop{}).Where("id = ? ", orderRow.ShopId).Limit(1).Find(&shopRow)
			if shopRow.Id == 0 {
				continue
			}
			sheetRow =&SheetRow{
				OrderA2: fmt.Sprintf("DCY.%v.%v",nowTimeObj.Format("20060102"),index + 1),
				SheetName: shopRow.Name,
				ShopPhone: shopRow.Phone,
				ShopAddress: shopRow.Address,
				ShopUserValue: shopRow.UserName,
				OrderCreateTime: orderRow.CreatedAt.Format("2006-01-02 15:04"),
			}
		}

		//获取订单关联的具体规格
		var orderSpecs []models.OrderSpecs


		e.Orm.Table(splitTableRes.OrderSpecs).Where("order_id = ?", orderRow.OrderId).Find(&orderSpecs)

		specsList := make([]*XlsxTableRow, 0)
		for _, row := range orderSpecs {
			xlsx :=&XlsxTableRow{
				GoodsName: row.GoodsName,
				GoodsSpecs: row.SpecsName,
				Unit: row.Unit,
				Number: row.Number,
				Price: row.Money,
				TotalMoney: utils.RoundDecimalFlot64(row.Money  * float64(row.Number)),
			}
			specsList = append(specsList, xlsx)
		}
		sheetRow.Table = append(sheetRow.Table ,specsList...)

		siteMap[orderRow.ShopId] = sheetRow

	}
	//基于上面汇聚的值 在做一次汇总统计

	for s :=range siteMap{
		sheetRowObject :=siteMap[s]
		//对table的数据进行汇总
		for index,v :=range siteMap[s].Table{
			v.Id = index + 1
			sheetRowObject.AllNumber+=v.Number
			sheetRowObject.AllMoney = utils.RoundDecimalFlot64(sheetRowObject.AllMoney) + v.TotalMoney
		}
		sheetRowObject.MoneyCn = utils.ConvertNumToCny(sheetRowObject.AllMoney)
		siteMap[s] = sheetRowObject

	}
	return siteMap,nil
}


//导出数据保存在云端

func (e *OrderExportObj)SaveExportXlsx(redisRow global.ExportRedisInfo,sheetData map[int]*SheetRow)  {
	//fmt.Println("保存到本地zip文件",e.Dat)

	export :=XlsxBaseExport{
		ExportUser: redisRow.ExportUser,
		ExportTime: redisRow.ExportTime,
	}
	// 大B/order_export/2023年12月26日15:46:06.xlsx
	FileName := export.SetXlsxRun(redisRow.CId,sheetData)

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
