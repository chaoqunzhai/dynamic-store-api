/*
*
@Author: chaoqun
* @Date: 2023/6/3 23:14
*/
package business

import (
	"encoding/json"
	"fmt"
	models2 "go-admin/cmd/migrate/migration/models"
	"go-admin/common/redis_db"
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
	OrderEdit string `json:"order_edit"` //订单修改表
	OrderReturn string `json:"order_return"` //订单退换货表
}

func (t *GetSplitTable)GetDbTableMapCnf() (res TableRow)  {
	var splitRow models2.SplitTableMap
	res = TableRow{
		OrderTable: global.SplitOrderDefaultTableName,
		OrderSpecs: global.SplitOrderDefaultSubTableName,
		OrderCycle: global.SplitOrderCycleSubTableName,
		OrderEdit:global.SplitOrderEdit,
		OrderReturn: global.SplitOrderReturn,
	}
	t.Orm.Model(&models2.SplitTableMap{}).Where("c_id = ? and enable = ? ", t.CId, true).Limit(1).Find(&splitRow)

	if splitRow.Id == 0 {
		return  res
	}
	//增加无自定义表 默认读取原表
	return TableRow{
		OrderTable:  func()  string {
			if splitRow.OrderTable == ""{
				return global.SplitOrderDefaultTableName
			}
			return splitRow.OrderTable
		}(),
		OrderSpecs:  func()  string {
			if splitRow.OrderSpecs == ""{
				return global.SplitOrderDefaultSubTableName
			}
			return splitRow.OrderSpecs
		}(),
		OrderCycle:  func()  string {
			if splitRow.OrderCycle == ""{
				return global.SplitOrderCycleSubTableName
			}
			return splitRow.OrderCycle
		}(),
		OrderEdit:   func()  string {
			if splitRow.OrderEdit == ""{
				return global.SplitOrderEdit
			}
			return splitRow.OrderEdit
		}(),
		OrderReturn: func()  string {
			if splitRow.OrderReturn == ""{
				return global.SplitOrderReturn
			}
			return splitRow.OrderReturn
		}(),
	}
}
//请求频率比较高，需要缓存到redis中
func (t *GetSplitTable)GetTableMap() (res TableRow)  {


	//从redis中读取
	//如果redis中没有,那就读取DB数据,并把数据load到redis中

	var redisErr error
	redisData,redisErr:=redis_db.GetSplitTableCnf(t.CId)

	if redisErr !=nil{
		//读取db 并返回
		//go 协程写入redis中
		fmt.Println("redis暂无数据,返回DB数据, 数据开始写入redis中")
		dbSplitCnf :=t.GetDbTableMapCnf()

		go func() {
			redis_db.SetCompanyTableSplitCnf(t.CId,dbSplitCnf)
		}()
		return dbSplitCnf
	}else {
		//序列化成TableRow配置
		fmt.Println("读取redis分表配置成功！")
		var tableRow TableRow
		unbarErr := json.Unmarshal([]byte(redisData),&tableRow)
		//json 失败 返回
		if unbarErr !=nil{
			fmt.Println("splitTableCnf配置反序列化失败,读取DB中配置！")
			return t.GetDbTableMapCnf()
		}
		//如果没有问题 就返回配置
		return tableRow
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
