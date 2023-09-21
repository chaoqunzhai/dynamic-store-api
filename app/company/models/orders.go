package models

import (
	"go-admin/common/models"
)

type Orders struct {
	models.Model

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
	DeductionMoney float64  `json:"deduction_money" gorm:"comment:抵扣额"`
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

// 订单规格
type OrderSpecs struct {
	models.Model
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

// todo:存放一些无关紧要的非必要首次查询的数据
type OrderExtend struct {
	models.Model
	Buyer string `json:"buyer" gorm:"size:27;comment:留言"`
	CId       int       `gorm:"index;comment:大BID"`
	CreatedAt models.XTime       `json:"createdAt" gorm:"comment:创建时间"`
	OrderId string `json:"order_id" gorm:"index;size:20;comment:订单ID"`
	Driver    string         `gorm:"size:12;comment:司机名称"`
	Phone     string         `gorm:"size:11;comment:联系手机号"`
	UserName string `json:"user_name" gorm:"size:20;comment:联系用户名称"`
	UserAddress string `json:"user_address" gorm:"size:70;comment:用户地址"`
	CouponMoney float64 `json:"coupon_money" gorm:"comment:优惠卷金额"`
}

func (OrderExtend) TableName(tableName string) string {
	if tableName == "" {
		return "order_extend"
	} else {
		return tableName
	}
}

func (e *OrderExtend) GetId() interface{} {
	return e.Id
}

// todo:周期列表
type OrderCycleList struct {
	models.Model
	CreatedAt models.XTime `json:"createdAt" gorm:"comment:创建时间"`
	CId       int          `gorm:"index;comment:大BID"`
	Name      string       `gorm:"size:12;comment:当前下单周期日期名称"`
	Uid       string       `gorm:"type:varchar(4);comment:周期名称都是天,防止一天可能多个不同周期的配置,加个标识区分周期"`
	StartTime models.XTime `gorm:"comment:此周期,下单周期开始时间"`
	EndTime   models.XTime `gorm:"comment:此周期,下单周期结束时间"`
	CycleTime models.XTime `json:"cycle_time" gorm:"type:date;comment:计算的配送时间"`
	CycleStr  string       `json:"cycle_str" gorm:"index;size:14;comment:配送时间的文案"`
	SaleMoney float64      `gorm:"comment:销售总额"`
	GoodsAll  int          `gorm:"comment:商品总数"`
	ShopCount int          `gorm:"type:tinyint(3);comment:客户总数"`
}

func (OrderCycleList) TableName() string {
	return "order_cycle_list"
}
