/**
@Author: chaoqun
* @Date: 2024/1/9 14:38
*/
package models

type Inventory struct {
	BigBRichGlobal

}
func (Inventory) TableName() string {
	return "company_inventory"
}
//出入库记录

type InventoryRecord struct {
	MiniLog
	GoodsId int `json:"goods_id" gorm:"comment:关联商品"`
	GoodsSpec int `json:"goods_spec" gorm:"comment:关联商品规格"`
	Action int  `json:"action" gorm:"type:tinyint(1);default:1;index;comment:操作, 1:入库 0:出库"`

}
func (InventoryRecord) TableName() string {
	return "company_inventory_record"
}