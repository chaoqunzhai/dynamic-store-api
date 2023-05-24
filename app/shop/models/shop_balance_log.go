package models

import (
     
     
     "time"

	"go-admin/common/models"

)

type ShopBalanceLog struct {
    models.Model
    
    ShopId string `json:"shopId" gorm:"type:bigint(20);comment:小BID"` 
    Money string `json:"money" gorm:"type:double;comment:变动金额"` 
    Scene string `json:"scene" gorm:"type:varchar(30);comment:变动场景"` 
    Desc string `json:"desc" gorm:"type:varchar(50);comment:描述/说明"` 
    models.ModelTime
    models.ControlBy
}

func (ShopBalanceLog) TableName() string {
    return "shop_balance_log"
}

func (e *ShopBalanceLog) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *ShopBalanceLog) GetId() interface{} {
	return e.Id
}