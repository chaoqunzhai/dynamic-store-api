package models

import (
     


	"go-admin/common/models"

)

type ShopTag struct {
    models.Model
    
    Layer string `json:"layer" gorm:"type:tinyint(4);comment:排序"` 
    Enable string `json:"enable" gorm:"type:tinyint(1);comment:开关"` 
    Desc string `json:"desc" gorm:"type:varchar(25);comment:描述信息"` 
    CId string `json:"cId" gorm:"type:bigint(20);comment:大BID"` 
    Name string `json:"name" gorm:"type:varchar(35);comment:客户标签名称"` 
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