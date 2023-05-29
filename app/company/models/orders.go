package models

import (
	"go-admin/common/models"
)

type Orders struct {
	models.Model

	Layer    string `json:"layer" gorm:"type:tinyint;comment:排序"`
	Enable   string `json:"enable" gorm:"type:tinyint(1);comment:开关"`
	Desc     string `json:"desc" gorm:"type:varchar(35);comment:描述信息"`
	CId      string `json:"cId" gorm:"type:bigint;comment:大BID"`
	ShopId   string `json:"shopId" gorm:"type:bigint;comment:关联客户"`
	Status   string `json:"status" gorm:"type:bigint;comment:配送状态"`
	Money    string `json:"money" gorm:"type:double;comment:金额"`
	Number   string `json:"number" gorm:"type:bigint;comment:下单数量"`
	Delivery string `json:"delivery" gorm:"type:bigint;comment:配送周期"`
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
