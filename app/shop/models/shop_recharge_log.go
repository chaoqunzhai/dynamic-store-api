package models

import (
     
     
     
     
     
     "time"

	"go-admin/common/models"

)

type ShopRechargeLog struct {
    models.Model
    
    ShopId string `json:"shopId" gorm:"type:bigint(20);comment:小BID"` 
    Uuid string `json:"uuid" gorm:"type:varchar(10);comment:订单号"` 
    Source string `json:"source" gorm:"type:varchar(16);comment:充值方式"` 
    Money string `json:"money" gorm:"type:double;comment:支付金额"` 
    GiveMoney string `json:"giveMoney" gorm:"type:double;comment:赠送金额"` 
    PayStatus string `json:"payStatus" gorm:"type:tinyint(1);comment:支付状态"` 
    PayTime time.Time `json:"payTime" gorm:"type:datetime(3);comment:支付时间"` 
    models.ModelTime
    models.ControlBy
}

func (ShopRechargeLog) TableName() string {
    return "shop_recharge_log"
}

func (e *ShopRechargeLog) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *ShopRechargeLog) GetId() interface{} {
	return e.Id
}