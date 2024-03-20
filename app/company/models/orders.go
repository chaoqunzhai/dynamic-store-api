package models

import (
	"go-admin/common/models"
	"go-admin/common/utils"
	"gorm.io/gorm"
	"time"
)

type Orders struct {
	models.Model
	HelpBy int `json:"helpBy" gorm:"comment:代客下单用户ID"`
	Uid            string       `json:"uid" gorm:"size:21;index;comment:关联的OrderCycleCnf标识UID"`
	CreateBy       int          `json:"createBy" gorm:"index;comment:创建者"`
	SourceType     int          `json:"source_type" gorm:"type:tinyint(1);default:1;index;comment:订单来源"`
	PayType        int          `json:"pay_type" gorm:"type:tinyint(1);default:1;comment:支付方式"`
	Phone          string       `gorm:"size:11;index;comment:用户联系手机号"`
	CreatedAt      models.XTime `json:"createdAt" gorm:"comment:创建时间"`
	DeliveryRunAt  models.XTime `json:"delivery_run_at" gorm:"column:delivery_run_at; null; comment:开始配送时间"`
	OrderId        string       `json:"order_id" gorm:"index;size:20;comment:订单ID"`
	EmsId string `json:"ems_id"  gorm:"index;size:20;comment:快递单号"`
	OutTradeNo string `json:"out_trade_no" gorm:"size:20;comment:预支付流水号"`
	CId            int          `gorm:"index;comment:大BID"`
	Enable         bool         `gorm:"comment:开关"`
	ShopId         int          `json:"shop_id" gorm:"index;comment:关联客户"`
	OfflinePayId int `json:"offline_pay_id" gorm:"index;comment:线下付款方式"`
	Line           string       `gorm:"size:16;comment:路线名称"`
	LineId         int          `json:"line_id" gorm:"index;type:bigint;comment:线路ID"`
	Status         int          `gorm:"type:tinyint(1);default:0;index;comment:订单的状态"`
	DeliveryCode   string       `json:"delivery_code" gorm:"size:9;index;comment:核销码"`
	WriteOffStatus int          `json:"write_off_status" gorm:"type:tinyint(1);default:0;index;comment:核销状态,0:未核销 1:核销"`
	PayMoney       float64      `gorm:"comment:实际支付价"`
	OrderMoney     float64      `json:"order_money" gorm:"comment:需要支付价"`
	GoodsMoney     float64      `json:"goods_money" gorm:"comment:商品总价格"`
	DeductionMoney float64      `json:"deduction_money" gorm:"comment:抵扣费用"`
	Number         int          `gorm:"comment:下单数量"`
	PayStatus      int          `gorm:"type:tinyint(1);default:0;index;comment:支付状态,0:未付款,1:已付款 2:线下付款，3:下线付款已收款"`
	PayTime        models.XTime `json:"pay_time" gorm:"comment:支付时间"`
	DeliveryTime   models.XTime `json:"delivery_time" gorm:"type:date;comment:下单时计算的配送时间"`
	DeliveryType   int          `json:"delivery_type" gorm:"comment:配送类型"`
	DeliveryID     int          `json:"delivery_id" gorm:"comment:关联的配送ORM-ID,根据配送类型来查询相关数据"`
	DeliveryStr    string       `json:"delivery_str" gorm:"size:25;comment:配送文案"`
	DeliveryMoney  float64      `json:"delivery_money" gorm:"comment:订单总配送费"`
	CouponId       int          `json:"coupon_id" gorm:"comment:使用优惠卷的ID"`
	DriverId       int          `gorm:"index;comment:司机ID"`
	AddressId      int          `json:"user_address_id" gorm:"index;comment:用户的收货地址,或者自提的店家地址"`
	CouponMoney    float64      `json:"coupon_money" gorm:"comment:优惠卷金额/代客下单优惠金额"`
	Buyer          string       `json:"buyer" gorm:"size:24;comment:留言"`
	Desc           string       `json:"desc" gorm:"size:16;comment:备注"`
	Edit bool `json:"edit" gorm:"comment:是否被修改"`
	EditAction string `json:"edit_action" gorm:"size:16;comment:修改"`
	AfterSales     bool         `json:"after_sales" gorm:";comment:是否申请售后,只有小B主动申请生效"`
	AfterStatus    int          `json:"after_status" gorm:"type:tinyint(1);default:0;;comment:售后状态:-2:撤回 -1:驳回, 0:无售后, 1:售后处理中 2:处理完毕  3: 大B退回"`
	ApproveMsg  string    `json:"approve_msg" gorm:"type:varchar(12);comment:审批信息/驳回|作废信息"`
	ApproveStatus int `json:"approve_status" gorm:"type:tinyint(1);default:0;;comment:订单审核状态 1:通过 0:驳回"`
	AcceptMsg string `json:"accept_msg"  gorm:"type:varchar(35);comment:欠账验收信息"`
	AcceptMoney float64      `json:"accept_money" gorm:"comment:收款费用"`
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
	GoodsName string       `json:"goods_name" gorm:"size:30;comment:商品名称"`
	CId       int          `gorm:"index;comment:大BID"`
	SpecId    int          `json:"spec_id" gorm:"index;comment:规格ID"`
	CreatedAt models.XTime `json:"created_at" gorm:"comment:创建时间"`
	OrderId   string       `gorm:"index;size:30;comment:关联订单长ID"`
	GoodsId   int          `json:"goods_id" gorm:"comment:商品ID"`
	SpecsName string       `gorm:"size:30;comment:规格名称"`
	Unit      string       `json:"unit" gorm:"type:varchar(8);comment:单位"`
	Number    int          `gorm:"comment:下单规格数"`
	Status    int          `gorm:"type:tinyint(1);default:1;index;comment:配送状态"`
	Money     float64      `gorm:"comment:规格的单价"`
	Image     string       `gorm:"size:15;comment:商品图片路径"`
	Edit      bool         `json:"edit" gorm:"comment:是否被修改"`
	EditAction string `json:"edit_action" gorm:"size:12;comment:修改/作废说明"`
	AfterStatus    int          `json:"after_status" gorm:"type:tinyint(1);default:0;;comment:售后状态:-2:撤回 -1:驳回, 0:无售后, 1:售后处理中 2:处理完毕  3: 大B退回"`
	AllMoney  float64      `json:"all_money" gorm:"comment:计算的规格价格"` //创建时 计算好

}
func (e *OrderSpecs) BeforeCreate(_ *gorm.DB) error {
	if e.Number > 0 && e.Money > 0 {
		e.AllMoney = utils.RoundDecimalFlot64(e.Money * float64(e.Number))
	}
	return nil
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
	//Buyer string `json:"buyer" gorm:"size:24;comment:留言"`
	CId       int       `gorm:"index;comment:大BID"`
	CreatedAt models.XTime       `json:"createdAt" gorm:"comment:创建时间"`
	OrderId string `json:"order_id" gorm:"index;size:20;comment:订单ID"`
	//DriverId    int         `gorm:"comment:司机ID"`
	////如果是周期配送,就是保存的小B地址的ID
	////如果是自提,那保存的是大B的设置的自提点地址的ID
	//AddressId int `json:"user_address_id" gorm:"comment:用户的收货地址,或者自提的店家地址"`
	//CouponMoney float64 `json:"coupon_money" gorm:"comment:优惠卷金额"`
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


// todo:订单周期列表的记录
type OrderCycleCnf struct {

	models.Model
	CreatedAt models.XTime `json:"createdAt" gorm:"comment:创建时间"`
	CId       int       `gorm:"index;comment:大BID"`
	Uid            string       `json:"uid" gorm:"size:21;index;comment:标识UID"`
	//StartTime 查看下单周期开始时间 EndTime 查看下单结束周期
	//直接算好时间,方便订单查询,因为下单的时候 是可以算出来的
	StartTime time.Time `json:"start_time" gorm:"comment:记录可下单周期开始时间"`
	EndTime   time.Time `json:"end_time" gorm:"comment:记录可下单周期结束时间"`
	//下单周期的文案也是保持最新的
	CreateStr string `json:"create_str" gorm:"size:30;comment:下单日期的文案内容"`
	//配送周期的统一查询
	DeliveryTime models.XTime `json:"delivery_time" gorm:"type:date;comment:计算的配送时间"`
	//展示,也是保持最新的
	DeliveryStr string `json:"delivery_str" gorm:"size:30;comment:配送文案"`
}
func (OrderCycleCnf) TableName(tableName string) string {
	if tableName == "" {
		return "order_cycle_cnf"
	} else {
		return tableName
	}
}

func (e *OrderCycleCnf) GetId() interface{} {
	return e.Id
}