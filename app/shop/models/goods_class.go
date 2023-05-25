package models

import (
	"go-admin/common/models"
)

type GoodsClass struct {
	models.Model

	Layer  string `json:"layer" gorm:"type:tinyint(4);comment:排序"`
	Enable string `json:"enable" gorm:"type:tinyint(1);comment:开关"`
	Desc   string `json:"desc" gorm:"type:varchar(25);comment:描述信息"`
	CId    string `json:"cId" gorm:"type:bigint(20);comment:大BID"`
	Name   string `json:"name" gorm:"type:varchar(35);comment:商品分类名称"`
	models.ModelTime
	models.ControlBy
}

func (GoodsClass) TableName() string {
	return "goods_class"
}

func (e *GoodsClass) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *GoodsClass) GetId() interface{} {
	return e.Id
}
