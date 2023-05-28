package models

import (
	"go-admin/common/models"
)

type GoodsSpecs struct {
	models.Model

	Layer     int `json:"layer" gorm:"type:tinyint(4);comment:排序"`
	Enable    bool `json:"enable" gorm:"type:tinyint(1);comment:开关"`
	CId       int `json:"cId" gorm:"type:bigint(20);comment:大BID"`
	GoodsId   int `json:"goodsId" gorm:"type:bigint(20);comment:商品ID"`
	Name      string `json:"name" gorm:"type:varchar(30);comment:规格名称"`
	Price     float64 `json:"price" gorm:"type:float;comment:售价"`
	Original  float64 `json:"original" gorm:"type:float;comment:原价"`
	Inventory int `json:"inventory" gorm:"type:bigint(20);comment:库存"`
	Unit      string `json:"unit" gorm:"type:varchar(8);comment:单位"`
	Limit     int `json:"limit" gorm:"type:bigint(20);comment:起售量"`
	models.ModelTime
	models.ControlBy
}

func (GoodsSpecs) TableName() string {
	return "goods_specs"
}

func (e *GoodsSpecs) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *GoodsSpecs) GetId() interface{} {
	return e.Id
}