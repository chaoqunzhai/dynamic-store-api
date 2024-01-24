package models

import (
     
     "go-admin/common/models"

)

type GoodsUnit struct {
    models.Model
    
    Layer int `json:"layer" gorm:"type:tinyint(4);comment:排序"`
    CId int `json:"cId" gorm:"type:bigint(20);comment:大BID"`
    Name string `json:"name" gorm:"type:varchar(8);comment:单位"` 
    models.ModelTime
    models.ControlBy
}

func (GoodsUnit) TableName() string {
    return "goods_unit"
}

func (e *GoodsUnit) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *GoodsUnit) GetId() interface{} {
	return e.Id
}