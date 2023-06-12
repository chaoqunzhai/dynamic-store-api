/*
*
@Author: chaoqun
* @Date: 2023/6/3 23:14
*/
package business

import (
	"fmt"
	models2 "go-admin/cmd/migrate/migration/models"
	"go-admin/global"
	"gorm.io/gorm"
	"strings"
)

func GetTableName(cid int, orm *gorm.DB) string {
	//先在split分表中查询

	var splitRow models2.SplitTableMap
	orm.Model(&models2.SplitTableMap{}).Where("c_id = ? and enable = ? and type = ?", cid, true, global.SplitOrder).Limit(1).Find(&splitRow)

	tableName := ""
	if splitRow.Id > 0 {
		tableName = splitRow.Name
	} else {
		tableName = global.SplitOrderDefaultTableName
	}

	return tableName
}
func OrderExtendTableName(orderTable string) string {
	//子表默认名称
	specsTable := global.SplitOrderExtendSubTableName
	//判断是否分表了
	//默认是 orders 表名，如果分表后就是 orders_大BID_时间戳后6位

	if orderTable != global.SplitOrderDefaultTableName {
		//拼接位 order_specs_大BID_时间戳后6位
		specsTable = fmt.Sprintf("%v%v", specsTable, strings.Replace(orderTable, global.SplitOrderExtendSubTableName, "", -1))

	}
	return specsTable
}
func OrderSpecsTableName(orderTable string) string {
	//子表默认名称
	specsTable := global.SplitOrderDefaultSubTableName
	//判断是否分表了
	//默认是 orders 表名，如果分表后就是 orders_大BID_时间戳后6位

	if orderTable != global.SplitOrderDefaultTableName {
		//拼接位 order_specs_大BID_时间戳后6位
		specsTable = fmt.Sprintf("%v%v", specsTable, strings.Replace(orderTable, global.SplitOrderDefaultTableName, "", -1))
	}
	return specsTable
}
