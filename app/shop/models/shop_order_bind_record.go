package models

import (
     


	"go-admin/common/models"

)

type ShopOrderBindRecord struct {
    models.Model
    
    ShopId string `json:"shopId" gorm:"type:bigint(20);comment:关联的小B客户"` 
    RecordId string `json:"recordId" gorm:"type:bigint(20);comment:每次记录的总ID"` 
    OrderId string `json:"orderId" gorm:"type:bigint(20);comment:订单ID"` 
    models.ModelTime
    models.ControlBy
}

func (ShopOrderBindRecord) TableName() string {
    return "shop_order_bind_record"
}

func (e *ShopOrderBindRecord) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *ShopOrderBindRecord) GetId() interface{} {
	return e.Id
}