package models

import (
	"go-admin/common/models"
)

type Orders struct {
	models.Model

	Layer        int          `json:"layer" gorm:"type:tinyint;comment:排序"`
	Enable       bool         `json:"enable" gorm:"type:tinyint(1);comment:开关"`
	Desc         string       `json:"desc" gorm:"type:varchar(35);comment:描述信息"`
	CId          int          `json:"c_id" gorm:"type:bigint;comment:大BID"`
	ShopId       int          `json:"shop_id" gorm:"type:bigint;comment:关联客户"`
	ClassId      int          `json:"class_id"`
	Status       int          `json:"status" gorm:"type:bigint;comment:配送状态"`
	Money        float64      `json:"money" gorm:"type:double;comment:下单总金额"`
	Number       int          `json:"number" gorm:"type:bigint;comment:下单产品数量"`
	Pay          int          `json:"pay" gorm:"type:bigint;comment:支付方式"`
	DeliveryId   int          `json:"delivery_id" gorm:"type:bigint;comment:配送周期"`
	DeliveryTime models.XTime `json:"delivery_time"`
	DeliveryStr  string       `json:"delivery_str" gorm:"type:varchar(14);comment:配送时间"`
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

type OrderSpecs struct {
	models.Model

	OrderId   int           `json:"orderId" gorm:"type:bigint(20);comment:关联订单ID"`
	SpecsId   int           `json:"specsId" gorm:"type:bigint(20);comment:规格ID"`
	Status    int           `json:"status" gorm:"type:bigint(20);comment:配送状态"`
	Money     float64       `json:"money" gorm:"type:double;comment:价格"`
	Number    int           `json:"number" gorm:"type:bigint(20);comment:下单产品数"`
	CreatedAt models.XTime  `json:"created_at" gorm:"comment:创建时间"`
	DeletedAt *models.XTime `json:"-" gorm:"index;comment:删除时间"`
}

func (OrderSpecs) TableName(tableName string) string {
	if tableName == "" {
		return "order_specs"
	} else {
		return tableName
	}
}

func (e *OrderSpecs) GetId() interface{} {
	return e.Id
}
