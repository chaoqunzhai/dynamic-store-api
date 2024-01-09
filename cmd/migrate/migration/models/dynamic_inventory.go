/**
@Author: chaoqun
* @Date: 2024/1/9 14:38
*/
package models


type InventoryCnf struct {
	MiniGlobal
	CId            int       `json:"c_id" gorm:"index;comment:公司(大B)ID"`
}
func (InventoryCnf) TableName() string {
	return "company_inventory_cnf"
}
type Inventory struct {
	BigBRichGlobal
	GoodsId int `json:"goods_id" gorm:"index;comment:商品ID"`
	SpecId int `json:"spec_id" gorm:"index;comment:规格ID"`
	Image     string  `gorm:"size:15;comment:商品图片路径"`
	Stock int64 `json:"stock" gorm:"comment:仓库数量"`
	OriginalPrice float64 `json:"original_price" gorm:"comment:当前入库价"`
}
func (Inventory) TableName() string {
	return "company_inventory"
}

//出入库记录

type InventoryRecord struct {
	MiniLog
	ArtNo string `json:"art_no" gorm:"size:20;comment:货架编号"`
	Code      string  `gorm:"size:20;comment:条形码"`
	Image     string  `gorm:"size:15;comment:商品图片路径"`
	GoodsName string `json:"goods_name" gorm:"size:50;comment:入库商品"`
	GoodsSpecName string `json:"goods_spec_name" gorm:"size:50;comment:入库商品规格"`
	Unit string `json:"unit" gorm:"size:20;comment:入库商品单位"`
	Number int `json:"number"  gorm:"comment:数量"`
	OriginalPrice float64 `json:"original_price" gorm:"comment:入库价"`
	Action int  `json:"action" gorm:"type:tinyint(1);default:1;index;comment:操作, 1:入库 0:出库"`

}
func (InventoryRecord) TableName() string {
	return "company_inventory_record"
}