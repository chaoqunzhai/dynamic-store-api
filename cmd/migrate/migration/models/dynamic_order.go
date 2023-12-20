package models

import (
	"go-admin/common/models"
	"time"
)

type CycleTimeConf struct {
	BigBRichGlobal
	Uid       string `json:"uid" gorm:"type:varchar(4);index;comment:周期名称都是天,防止一天可能多个不同周期的配置,加个标识区分周期"`
	Show      bool   `json:"show" gorm:"default:true;comment:是否客户端展示"`
	Type      int    `json:"type" gorm:"type:tinyint(1);default:1;comment:类型,每天,每周"`
	StartWeek int    `json:"start_week" gorm:"type:tinyint(1);default:0;comment:类型为周,每周开始天"`
	EndWeek   int    `json:"end_week" gorm:"type:tinyint(1);default:0;comment:类型为周,每周结束天"`
	StartTime string `json:"start_time" gorm:"size:5;comment:开始下单时间"`
	EndTime   string `json:"end_time" gorm:"size:5;comment:结束时间"`
	GiveDay   int    `json:"give_day" gorm:"type:tinyint(1);default:0;comment:跨天值为0是当天,大于0就是当天+天数"`
	GiveTime  string `json:"give_time" gorm:"size:14;comment:配送时间,例如：15点至19点"`
	Desc      string `json:"desc" gorm:"size:30;comment:描述信息"` //描述
}

func (CycleTimeConf) TableName() string {
	return "cycle_time_conf"
}

// todo:订单,因为订单是一个记录,所有大部分可变的数据都是静态资源,不做关联查询
type Orders struct {
	Model
	Uid string `json:"uid" gorm:"size:8;index;comment:关联的配送周期统一UID"`
	CreateBy int `json:"createBy" gorm:"index;comment:创建者"`
	SourceType int `json:"source_type" gorm:"type:tinyint(1);default:1;index;comment:订单来源"`
	PayType    int            `json:"pay_type" gorm:"type:tinyint(1);default:1;comment:支付方式"`
	Phone     string         `gorm:"size:11;index;comment:用户联系手机号"`
	CreatedAt models.XTime `json:"createdAt" gorm:"comment:创建时间"`
	UpdatedAt models.XTime `json:"updatedAt" gorm:"comment:最后更新时间"`
	OrderId string `json:"order_id" gorm:"index;size:20;comment:订单ID"`
	OrderNoId string `json:"order_no_id" gorm:"size:20;comment:支付流水号"`
	CId       int       `gorm:"index;comment:大BID"`
	Enable    bool      `gorm:"comment:开关"`
	ShopId    int       `gorm:"index;comment:关联客户"`
	Line      string    `gorm:"size:16;comment:路线名称"`
	LineId    int       `json:"line_id" gorm:"index;type:bigint;comment:线路ID"`
	Status    int       `gorm:"type:tinyint(1);default:0;index;comment:订单的状态,0:待配送,1:配送中,2:已配送,3:退回,4:退款"`
	DeliveryCode string `json:"delivery_code" gorm:"size:9;index;comment:核销码"`
	WriteOffStatus  int `json:"write_off_status" gorm:"type:tinyint(1);default:0;index;comment:核销状态,0:未核销 1:核销"`
	PayMoney     float64   `gorm:"comment:实际支付价"`
	OrderMoney float64  `gorm:"comment:需要支付价"`
	GoodsMoney float64 `gorm:"comment:商品总价格"`
	DeductionMoney float64  `json:"deduction_money" gorm:"comment:抵扣费用"`
	Number    int       `gorm:"comment:下单数量"`
	PayStatus int        `gorm:"type:tinyint(1);default:0;index;comment:支付状态,0:未付款,1:已付款 2:线下付款，3:下线付款已收款"`
	PayTime models.XTime `json:"pay_time" gorm:"comment:支付时间"`
	DeliveryTime models.XTime `json:"delivery_time" gorm:"type:date;comment:下单时计算的配送时间"`
	DeliveryUid  string    `gorm:"type:varchar(4);comment:周期名称都是天,防止一天可能多个不同周期的配置,加个标识区分周期"`
	DeliveryType int `json:"delivery_type" gorm:"comment:配送类型"`
	DeliveryID int `json:"delivery_id" gorm:"comment:关联的配送ID,根据配送类型来查询相关数据"`
	DeliveryStr  string    `json:"delivery_str" gorm:"size:25;comment:配送文案"`
	DeliveryMoney float64 `json:"delivery_money" gorm:"comment:配送费"`
	CouponId int `json:"coupon_id" gorm:"comment:使用优惠卷的ID"`
	DriverId    int         `gorm:"index;comment:司机ID"`
	AddressId int `json:"user_address_id" gorm:"index;comment:用户的收货地址,或者自提的店家地址"`
	CouponMoney float64 `json:"coupon_money" gorm:"comment:优惠卷金额"`
	Buyer string `json:"buyer" gorm:"size:24;comment:留言"`
	Desc string `json:"desc" gorm:"size:16;comment:备注"`
}

func (Orders) TableName() string {
	return "orders"
}

// todo:订单规格
type OrderSpecs struct {
	Model
	GoodsName string        `json:"goods_name" gorm:"size:50;comment:商品名称"`
	CId       int       `gorm:"index;comment:大BID"`
	SpecId int `gorm:"index;comment:规格ID"`
	CreatedAt models.XTime  `json:"created_at" gorm:"comment:创建时间"`
	OrderId   string            `gorm:"index;size:30;comment:关联订单长ID"`
	GoodsId int `json:"goods_id" gorm:"comment:商品ID"`
	SpecsName string         `gorm:"size:30;comment:规格名称"`
	Unit      string         `json:"unit" gorm:"type:varchar(8);comment:单位"`
	Number    int            `gorm:"comment:下单规格数"`
	Status    int            `gorm:"type:tinyint(1);default:1;index;comment:配送状态"`
	Money     float64        `gorm:"comment:规格的价格"`
	Image     string  `gorm:"size:15;comment:商品图片路径"`
}

func (OrderSpecs) TableName() string {
	return "order_specs"
}

// todo:订单扩展信息
type OrderExtend struct {
	Model
	//Buyer string `json:"buyer" gorm:"size:24;comment:留言"`
	CId       int       `gorm:"index;comment:大BID"`
	CreatedAt models.XTime       `json:"createdAt" gorm:"comment:创建时间"`
	OrderId string `json:"order_id" gorm:"index;size:20;comment:订单ID"`
	//DriverId    int         `gorm:"index;comment:司机ID"`
	//AddressId int `json:"user_address_id" gorm:"index;comment:用户的收货地址,或者自提的店家地址"`
	//CouponMoney float64 `json:"coupon_money" gorm:"comment:优惠卷金额"`
}

func (OrderExtend) TableName() string {
	return "order_extend"
}


// todo:订单周期列表的记录
type OrderCycleCnf struct {
	Model
	CreatedAt models.XTime `json:"createdAt" gorm:"comment:创建时间"`
	CId       int       `gorm:"index;comment:大BID"`
	Uid       string    `gorm:"type:varchar(8);comment:周期名称都是天,防止一天可能多个不同周期的配置,加个标识区分周期"`
	//StartTime 查看下单周期开始时间 EndTime 查看下单结束周期
	//直接算好时间,方便订单查询,因为下单的时候 是可以算出来的
	StartTime time.Time `gorm:"comment:记录可下单周期开始时间"`
	EndTime   time.Time `gorm:"comment:记录可下单周期结束时间"`
	//下单周期的文案也是保持最新的
	CreateStr string `json:"create_str" gorm:"size:30;comment:下单日期的文案内容"`
	//配送周期的统一查询
	DeliveryTime time.Time `json:"delivery_time" gorm:"type:date;index;comment:计算的配送时间"`
	//展示,也是保持最新的
	DeliveryStr string `json:"delivery_str" gorm:"size:30;comment:配送文案"`

	//应该增加一个周期的价格统计.商品统计 商品分类统计
}

func (OrderCycleCnf) TableName() string {
	return "order_cycle_cnf"
}


//todo:订单和Redis的映射表
type OrderToRedisMap struct {
	Model
	CreatedAt time.Time      `json:"createdAt" gorm:"comment:创建时间"`
	CId int `json:"c_id" gorm:"index;comment:大BID"`
	UserId int `json:"user_id" gorm:"index;comment:用户ID"`
	RedisKey string `json:"redis_key" gorm:"size:20;comment:redis的订单ID"`
	Status int `json:"status" gorm:"index;comment:订单状态"`
	RandomId string `json:"random_id" gorm:"size:10;comment:随机计算的ID,防止恶意请求订单,payInfo会产生"`

}
func (OrderToRedisMap) TableName() string {
	return "order_to_redis_map"
}