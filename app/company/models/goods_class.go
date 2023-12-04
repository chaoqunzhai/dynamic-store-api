package models

import (
	"go-admin/common/models"
)

type GoodsClass struct {
	models.Model

	Layer  int    `json:"layer" gorm:"type:tinyint(4);comment:排序"`
	Enable bool   `json:"enable" gorm:"type:tinyint(1);comment:开关"`
	Desc   string `json:"desc" gorm:"type:varchar(25);comment:描述信息"`
	CId    int    `json:"cId" gorm:"type:bigint(20);comment:大BID"`
	Recommend  bool `json:"recommend"`
	Name   string `json:"name" gorm:"type:varchar(35);comment:商品分类名称"`
	Image string `gorm:"size:60;comment:商品分类图片路径"`
	//只是做数据组装
	GoodsCount int64 `json:"goods_count" gorm:"-"`
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
