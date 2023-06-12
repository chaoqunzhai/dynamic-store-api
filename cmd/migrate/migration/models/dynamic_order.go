package models

import (
	"gorm.io/gorm"
	"time"
)

type CycleTimeConf struct {
	BigBRichGlobal
	Show      bool   `gorm:"default:true;comment:是否客户端展示"`
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

// todo:订单,因为订单是一个记录,所有大部分可变的数据都是静态资源,不做关联查询
type Orders struct {
	Model
	ControlBy
	ModelTime
	CId       int       `gorm:"index;comment:大BID"`
	Enable    bool      `gorm:"comment:开关"`
	GoodsId   int       `gorm:"index;comment:商品ID"`
	GoodsName string    `json:"goods_name" gorm:"size:35;comment:商品名称+广告"`
	ShopId    int       `gorm:"index;comment:关联客户"`
	Line      string    `gorm:"index;size:16;comment:路线名称"`
	LineId    int       `json:"line_id" gorm:"type:bigint;comment:线路ID"`
	Status    int       `gorm:"type:tinyint(1);default:1;index;comment:配送状态"`
	Money     float64   `gorm:"comment:下单费用"`
	Number    int       `gorm:"comment:下单数量"`
	Pay       int       `gorm:"type:tinyint(1);default:0;index;comment:支付方式,0:线上,1:线下"`
	PayStatus int       `gorm:"type:tinyint(1);default:0;index;comment:支付状态,0:未付款,1:已付款 2:线下付款，3:下线付款已收款"`
	CycleTime time.Time `json:"cycle_time" gorm:"type:date;comment:计算配送时间"`
	CycleStr  string    `json:"cycle_str" gorm:"index;size:14;comment:配送时间的文案"`
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
	SpecsName string         `gorm:"size:30;comment:规格名称"`
	Unit      string         `json:"unit" gorm:"type:varchar(8);comment:单位"`
	Number    int            `gorm:"comment:下单规格数"`
	Status    int            `gorm:"type:tinyint(1);default:1;index;comment:配送状态"`
	Money     float64        `gorm:"comment:规格的价格"`
}

func (OrderSpecs) TableName() string {
	return "order_specs"
}

// todo:订单扩展信息,存放一些其他无关紧要数据
type OrderExtend struct {
	Model
	CreatedAt time.Time      `json:"createdAt" gorm:"comment:创建时间"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间"`
	OrderId   int            `gorm:"index;comment:关联订单ID"`
	Driver    string         `gorm:"size:12;comment:司机名称"`
	Phone     string         `gorm:"size:11;comment:联系手机号"`
	Source    int            `gorm:"type:tinyint(1);default:0;index;comment:订单来源,客户下单还是,代客下单"`
	Desc      string         `gorm:"size:50;comment:描述信息"`
}

func (OrderExtend) TableName() string {
	return "order_extend"
}
