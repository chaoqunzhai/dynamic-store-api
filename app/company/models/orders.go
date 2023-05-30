package models

import (
	"go-admin/common/models"
)

type Orders struct {
	models.Model

	Layer    int `json:"layer" gorm:"type:tinyint;comment:排序"`
	Enable   bool `json:"enable" gorm:"type:tinyint(1);comment:开关"`
	Desc     string `json:"desc" gorm:"type:varchar(35);comment:描述信息"`
	CId      int `json:"c_id" gorm:"type:bigint;comment:大BID"`
	ShopId   int `json:"shop_id" gorm:"type:bigint;comment:关联客户"`
	Status   int `json:"status" gorm:"type:bigint;comment:配送状态"`
	Money    float64 `json:"money" gorm:"type:double;comment:下单总金额"`
	Number   int `json:"number" gorm:"type:bigint;comment:下单产品数量"`
	Pay   int `json:"pay" gorm:"type:bigint;comment:支付方式"`
	Delivery int `json:"delivery" gorm:"type:bigint;comment:配送周期"`
	models.ModelTime
	models.ControlBy
}

func (Orders) TableName(tableName string) string {
	if tableName == "" {
		return "orders"
	} else {
		return tableName
	}
}

func (e *Orders) GetId() interface{} {
	return e.Id
}
