package models

import (
	"gorm.io/gorm"
	"time"
)

type CycleTimeConf struct {
	BigBRichGlobal
	Type      int    `gorm:"type:tinyint(1);default:1;comment:类型,每天,每周"`
	StartWeek int    `gorm:"type:tinyint(1);default:0;comment:类型为周,每周开始天"`
	EndWeek   int    `gorm:"type:tinyint(1);default:0;comment:类型为周,每周结束天"`
	StartTime string `gorm:"size:5;comment:开始下单时间"`
	EndTime   string `gorm:"size:5;comment:结束时间"`
	GiveDay   int    `gorm:"type:tinyint(1);default:0;comment:跨天值为0是当天,大于0就是当天+天数"`
	GiveTime  string `gorm:"size:14;comment:配送时间,例如：15点至19点"`
	Desc      string `gorm:"size:30;comment:描述信息"` //描述
}

func (CycleTimeConf) TableName() string {
	return "cycle_time_conf"
}

// todo:订单
type Orders struct {
	BigBRichGlobal
	ShopId       int       `gorm:"index;comment:关联客户"`
	ClassId      int       `gorm:"index;comment:商品分类"`
	GoodId       int       `gorm:"index;comment:商品ID"`
	LineId       int       `gorm:"index;comment:线路ID"`
	Status       int       `gorm:"type:tinyint(1);default:1;index;comment:配送状态"`
	Money        float64   `gorm:"comment:金额"`
	Number       int       `gorm:"comment:下单数量"`
	Pay          int       `gorm:"type:tinyint(1);default:1;index;comment:支付方式"`
	DeliveryId   int       `gorm:"index;comment:配送时间周期"`
	DeliveryTime time.Time `json:"delivery_time" gorm:"type:date;comment:计算配送时间"`
	DeliveryStr  string    `gorm:"size:14;comment:配送时间,例如：15点至19点" json:"delivery_str"`
}

func (Orders) TableName() string {
	return "orders"
}

// todo:订单规格
type OrderSpecs struct {
	Model
	CreatedAt time.Time      `json:"createdAt" gorm:"comment:创建时间"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间"`
	OrderId   int            `gorm:"index;comment:关联订单ID"`
	SpecsId   int            `gorm:"index;comment:规格ID"`
	Number    int            `gorm:"comment:下单产品数"`
	Status    int            `gorm:"type:tinyint(1);default:1;index;comment:配送状态"`
	Money     float64        `gorm:"comment:价格"`
}

func (OrderSpecs) TableName() string {
	return "order_specs"
}
