package models

import (
     
     
     "go-admin/common/models"

)

type GoodsBrand struct {
    models.Model
    
    Layer int `json:"layer" gorm:"type:tinyint(4);comment:排序"`
    CId int `json:"cId" gorm:"type:bigint(20);comment:大BID"`
    Name string `json:"name" gorm:"type:varchar(12);comment:品牌名称"` 
    models.ModelTime
    models.ControlBy
}

func (GoodsBrand) TableName() string {
    return "goods_brand"
}

func (e *GoodsBrand) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *GoodsBrand) GetId() interface{} {
	return e.Id
}