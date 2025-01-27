package models

import (
	"go-admin/common/models"
)

type GoodsTag struct {
	models.Model

	Layer  int    `json:"layer" gorm:"type:tinyint(4);comment:排序"`
	Enable bool   `json:"enable" gorm:"type:tinyint(1);comment:开关"`
	Desc   string `json:"desc" gorm:"type:varchar(25);comment:描述信息"`
	CId    int    `json:"-" gorm:"type:bigint(20);comment:大BID"`
	Name   string `json:"name" gorm:"type:varchar(15);comment:商品标签名称"`
	Color  string `json:"color" gorm:"type:varchar(10);comment:标签颜色"`
	//只是做数据组装
	GoodsCount int64 `json:"goods_count" gorm:"-"`
	models.ModelTime
	models.ControlBy
}

func (GoodsTag) TableName() string {
	return "goods_tag"
}

func (e *GoodsTag) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *GoodsTag) GetId() interface{} {
	return e.Id
}
