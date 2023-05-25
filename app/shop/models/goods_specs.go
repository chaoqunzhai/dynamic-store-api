package models

import (
	"go-admin/common/models"
)

type GoodsSpecs struct {
	models.Model

	Layer     string `json:"layer" gorm:"type:tinyint(4);comment:排序"`
	Enable    string `json:"enable" gorm:"type:tinyint(1);comment:开关"`
	CId       string `json:"cId" gorm:"type:bigint(20);comment:大BID"`
	GoodsId   string `json:"goodsId" gorm:"type:bigint(20);comment:商品ID"`
	Name      string `json:"name" gorm:"type:varchar(30);comment:规格名称"`
	Price     string `json:"price" gorm:"type:float;comment:售价"`
	Original  string `json:"original" gorm:"type:float;comment:原价"`
	Inventory string `json:"inventory" gorm:"type:bigint(20);comment:库存"`
	Unit      string `json:"unit" gorm:"type:varchar(8);comment:单位"`
	Limit     string `json:"limit" gorm:"type:bigint(20);comment:起售量"`
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