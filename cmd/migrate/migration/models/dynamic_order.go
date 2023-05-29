package models

type CycleTimeConf struct {
	BigBRichGlobal
	Type      int    `gorm:"index;comment:类型,每天,每周"`
	StartWeek int    `gorm:"comment:类型为周,每周开始天"`
	EndWeek   int    `gorm:"comment:类型为周,每周结束天"`
	StartTime string `gorm:"size:5;comment:开始下单时间"`
	EndTime   string `gorm:"size:5;comment:结束时间"`
	GiveDay   int    `gorm:"comment:跨天值为0是当天,大于0就是当天+天数"`
	GiveTime  string `gorm:"size:30;comment:配送时间,例如：15点至19点"`
}

func (CycleTimeConf) TableName() string {
	return "cycle_time_conf"
}

// todo:订单
// 订单号是 1000 + id = 订单号
type Orders struct {
	BigBRichGlobal
	ShopId   int     `gorm:"index;comment:关联客户"`
	Status   int     `gorm:"index;comment:配送状态"`
	Money    float64 `gorm:"comment:金额"`
	Number   int     `gorm:"comment:下单数量"`
	Delivery int     `gorm:"comment:配送时间周期"`
}

func (Orders) TableName() string {
	return "orders"
}

// todo:订单规格
type OrderSpecs struct {
	Model
	ControlBy
	ModelTime
	OrderId int     `gorm:"index;comment:关联订单ID"`
	SpecsId int     `gorm:"index;comment:规格ID"`
	Status  int     `gorm:"index;comment:配送状态"`
	Money   float64 `gorm:"comment:价格"`
}

func (OrderSpecs) TableName() string {
	return "order_specs"
}
