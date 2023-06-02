package models

import (
	"go-admin/common/models"
)

type ShopTag struct {
	models.Model

	Layer  int    `json:"layer" gorm:"type:tinyint(4);comment:排序"`
	Enable bool   `json:"enable" gorm:"type:tinyint(1);comment:开关"`
	Desc   string `json:"desc" gorm:"type:varchar(25);comment:描述信息"`
	CId    int    `json:"-" gorm:"type:bigint(20);comment:大BID"`
	Name   string `json:"name" gorm:"type:varchar(35);comment:客户标签名称"`
	//只是做数据组装
	ShopCount int64 `json:"shop_count" gorm:"-"`
	models.ModelTime
	models.ControlBy
}

func (ShopTag) TableName() string {
	return "shop_tag"
}

func (e *ShopTag) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *ShopTag) GetId() interface{} {
	return e.Id
}
