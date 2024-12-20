package models

import (
	"go-admin/common/models"
	"gorm.io/gorm"
	"time"
)

type GoodsSpecs struct {
	models.Model

	Layer     int            `json:"layer" gorm:"type:tinyint(4);comment:排序"`
	Enable    bool           `json:"enable" gorm:"type:tinyint(1);comment:开关"`
	CId       int            `json:"-" gorm:"type:bigint(20);comment:大BID"`
	GoodsId   int            `json:"goods_id" gorm:"type:bigint(20);comment:商品ID"`
	Name      string         `json:"name" gorm:"type:varchar(20);comment:规格名称"`
	Price     float64        `json:"price" gorm:"type:float;comment:售价"`
	Market float64        `json:"market" gorm:"type:float;comment:市场价"`
	Original  float64        `json:"original" gorm:"type:float;comment:成本价"`
	Inventory int            `json:"inventory" gorm:"type:bigint(20);comment:库存"`
	Sale int   `gorm:"comment:销售量"`
	UnitId      int         `json:"unit_id" gorm:"type:varchar(8);comment:单位"`
	UnitName      string         `json:"unit" gorm:"-"`
	Limit     int            `json:"limit" gorm:"type:bigint(20);comment:起售量"`
	Max int  `json:"max" gorm:"type:bigint(20);comment:起售量"`
	Image     string         `json:"image" gorm:"size:100;comment:商品图片路径"`
	Code      string         `json:"code" gorm:"type:varchar(30);comment:规格名称"`
	VirtuallySale int `json:"virtually_sale" gorm:"comment:虚拟库存"`
	SerialNumber string `json:"serial_number" gorm:"size:20;comment:编号"`
	CreateBy  int            `json:"-" gorm:"index;comment:创建者"`
	CreatedAt time.Time      `json:"created_at" gorm:"comment:创建时间"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间"`
}

func (GoodsSpecs) TableName() string {
	return "goods_specs"
}

func (e *GoodsSpecs) GetId() interface{} {
	return e.Id
}
