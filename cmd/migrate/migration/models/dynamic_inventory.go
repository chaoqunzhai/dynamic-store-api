/**
@Author: chaoqun
* @Date: 2024/1/9 14:38
*/
package models

import (
	models2 "go-admin/common/models"
	"gorm.io/gorm"
	"time"
)

type InventoryCnf struct {
	MiniGlobal
	CId            int       `json:"c_id" gorm:"index;comment:公司(大B)ID"`
}
func (InventoryCnf) TableName() string {
	return "inventory_cnf"
}

//商品仓库数据
type Inventory struct {
	BigBRichGlobal
	GoodsName string `json:"goods_name" gorm:"size:50;comment:入库商品"`
	GoodsSpecName string `json:"goods_spec_name" gorm:"size:50;comment:入库商品规格"`
	GoodsId int `json:"goods_id" gorm:"index;comment:商品ID"`
	SpecId int `json:"spec_id" gorm:"index;comment:规格ID"`
	Stock int `json:"stock" gorm:"comment:仓库数量"`
	Image     string  `gorm:"size:15;comment:商品图片路径"`
	OriginalPrice float64 `json:"original_price" gorm:"comment:当前入库价/成本价"`
	Status     int          `json:"status" gorm:"type:tinyint(1);default:1;index;comment:销售状态  1:销售中 0:下线"`
}
func (Inventory) TableName() string {
	return "inventory"
}


//入库单 每次操作录入的记录
type InventoryOrder struct {
	Model
	CId            int       `json:"c_id" gorm:"index;comment:公司(大B)ID"`
	CreateBy  string            `json:"create_by" gorm:"size:20;comment:创建者"`
	CreatedAt models2.XTime      `json:"created_at" gorm:"comment:创建时间"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间"`
	Desc   string `json:"desc" gorm:"size:50;comment:描述信息"` //描述
	OrderId string `json:"order_id" gorm:"index;size:20;comment:出入库ID"`
	Action int `json:"type" gorm:"type:tinyint(1);default:0;index;comment:入库类型 1:入库 2:出库"`
	Number int `json:"number" gorm:"comment:商品数量"`
	DocumentMoney float64 `json:"document_money" gorm:"comment:单据金额"`
}
func (InventoryOrder) TableName() string {
	return "inventory_order"
}
//出入库单流水 流水做了分表策略
type InventoryRecord struct {
	Model
	CreatedAt time.Time      `json:"createdAt" gorm:"comment:操作时间"`
	CId int `gorm:"index;comment:大BID"`
	CreateBy  string            `json:"createBy" gorm:"size:20;comment:操作人"`
	OrderId string `json:"order_id" gorm:"index;size:20;comment:出入库单号ID"`
	Action int  `json:"action" gorm:"type:tinyint(1);default:1;index;comment:操作, 1:入库 0:出库"`

	ArtNo string `json:"art_no" gorm:"size:20;comment:货架编号"`
	Code      string  `gorm:"size:20;comment:条形码"`
	Image     string  `gorm:"size:15;comment:商品图片路径"`
	GoodsName string `json:"goods_name" gorm:"size:50;comment:入库商品"`
	GoodsSpecName string `json:"goods_spec_name" gorm:"size:50;comment:入库商品规格"`
	GoodsId int `json:"goods_id" gorm:"index;comment:商品ID"`
	SpecId int `json:"spec_id" gorm:"index;comment:规格ID"`
	Unit string `json:"unit" gorm:"size:8;comment:入库商品单位"`
	SourceNumber int `json:"source_number" gorm:"comment:原库存"`
	ActionNumber int `json:"action_number"  gorm:"comment:出入库操作数量"`
	CurrentNumber int `json:"current_number" gorm:"comment:现库存"`
	OriginalPrice float64 `json:"original_price" gorm:"comment:入库价/成本价"`

}
func (InventoryRecord) TableName() string {
	return "inventory_record"
}
