package models

import (


	"go-admin/common/models"

)

type ShopOrderRecord struct {
    models.Model
    
    ShopId string `json:"shopId" gorm:"type:bigint(20);comment:关联的小B客户"` 
    ShopName string `json:"shopName" gorm:"type:varchar(30);comment:客户名称"` 
    Money string `json:"money" gorm:"type:double;comment:订单金额"` 
    Number string `json:"number" gorm:"type:bigint(20);comment:订单量"` 
    models.ModelTime
    models.ControlBy
}

func (ShopOrderRecord) TableName() string {
    return "shop_order_record"
}

func (e *ShopOrderRecord) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *ShopOrderRecord) GetId() interface{} {
	return e.Id
}