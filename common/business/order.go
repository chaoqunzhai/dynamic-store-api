/*
*
@Author: chaoqun
* @Date: 2023/6/3 23:14
*/
package business

import (
	models2 "go-admin/cmd/migrate/migration/models"
	"go-admin/global"
	"gorm.io/gorm"
)

type GetSplitTable struct {
	CId int
	Orm  *gorm.DB
	SplitRow models2.SplitTableMap
}
type TableRow struct {
	OrderTable string `json:"order_table"`                   //订单表
	OrderSpecs string `json:"order_specs"` //订单规格表
	OrderCycle string `json:"order_cycle"` //周期配送下单索引表
}
func (t *GetSplitTable)GetTableMap() (res TableRow)  {
	var splitRow models2.SplitTableMap
	res = TableRow{
		OrderTable: global.SplitOrderDefaultTableName,
		OrderSpecs: global.SplitOrderDefaultSubTableName,
		OrderCycle: global.SplitOrderCycleSubTableName,
	}
	t.Orm.Model(&models2.SplitTableMap{}).Select("order_table,order_specs,order_cycle,id").Where("c_id = ? and enable = ? ", t.CId, true).Limit(1).Find(&splitRow)

	if splitRow.Id == 0 {
		return  res
	}


	return TableRow{
		OrderTable: splitRow.OrderTable,
		OrderSpecs: splitRow.OrderSpecs,
		OrderCycle: splitRow.OrderCycle,
	}
}

func GetTableName(cid int, orm *gorm.DB) (res TableRow)  {
	//先在split分表中查询

	split:=GetSplitTable{
		CId: cid,
		Orm: orm,
	}
	return split.GetTableMap()

}
//func OrderExtendTableName(orderTable string) string {
//	//子表默认名称
//	specsTable := global.SplitOrderExtendSubTableName
//	//判断是否分表了
//	//默认是 orders 表名，如果分表后就是 orders_大BID_时间戳后6位
//
//	if orderTable != global.SplitOrderDefaultTableName {
//		//拼接位 order_specs_大BID_时间戳后6位
//		specsTable = fmt.Sprintf("%v%v", specsTable, strings.Replace(orderTable, global.SplitOrderDefaultTableName, "", -1))
//
//	}
//	return specsTable
//}
//func OrderSpecsTableName(orderTable string) string {
//	//子表默认名称
//	table := global.SplitOrderDefaultSubTableName
//	//判断是否分表了
//	//默认是 orders 表名，如果分表后就是 orders_大BID_时间戳后6位
//
//	if orderTable != global.SplitOrderDefaultTableName {
//		//拼接位 order_specs_大BID_时间戳后6位
//		table = fmt.Sprintf("%v%v", table, strings.Replace(orderTable, global.SplitOrderDefaultTableName, "", -1))
//	}
//	return table
//}
//
//func OrderCycleTableName(orderTable string) string {
//	//子表默认名称
//	table := global.SplitOrderCycleSubTableName
//	//判断是否分表了
//	//默认是 orders 表名，如果分表后就是 orders_大BID_时间戳后6位
//
//	if orderTable != global.SplitOrderDefaultTableName {
//		//拼接位 order_specs_大BID_时间戳后6位
//		table = fmt.Sprintf("%v%v", table, strings.Replace(orderTable, global.SplitOrderDefaultTableName, "", -1))
//	}
//	return table
//}