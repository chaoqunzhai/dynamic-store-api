package models

type Cycle struct {
	BigBRichGlobal
	Name  string `gorm:"size:15;comment:周期名称"`
	Type  int    `gorm:"index;comment:类型"`
	Start string `gorm:"size:15;comment:开始下单时间"`
	End   string `gorm:"size:15;comment:结束时间"`
}

func (Cycle) TableName() string {
	return "cycle"
}

// todo:订单
// 订单号是 1000 + id = 订单号
type Orders struct {
	BigBRichGlobal
	ShopId   int     `gorm:"index;comment:关联客户"`
	Status   int     `gorm:"index;comment:配送状态"`
	Money    float64 `gorm:"comment:金额"`
	Number   int     `gorm:"comment:下单数量"`
	Delivery int     `gorm:"comment:配送周期"`
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
