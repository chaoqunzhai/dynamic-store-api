package models

import (
	"go-admin/common/models"
)

type GoodsSales struct {
	models.Model

	Layer       string `json:"layer" gorm:"type:tinyint(4);comment:排序"`
	Enable      string `json:"enable" gorm:"type:tinyint(1);comment:开关"`
	CId         string `json:"-" gorm:"type:bigint(20);comment:大BID"`
	ProductId   string `json:"productId" gorm:"type:bigint(20);comment:产品ID"`
	ProductName string `json:"productName" gorm:"type:varchar(30);comment:产品名称"`
	Sales       string `json:"sales" gorm:"type:bigint(20);comment:当时销量"`
	Inventory   string `json:"inventory" gorm:"type:bigint(20);comment:当时剩余库存"`
	models.ModelTime
	models.ControlBy
}

func (GoodsSales) TableName() string {
	return "goods_sales"
}

func (e *GoodsSales) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *GoodsSales) GetId() interface{} {
	return e.Id
}
