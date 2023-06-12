package models

import (
	"go-admin/common/models"
)

type Orders struct {
	models.Model
	Enable    bool         `json:"enable" gorm:"type:tinyint(1);comment:开关"`
	GoodsName string       `json:"goods_name" gorm:"size:35;comment:商品名称+广告"`
	CId       int          `json:"c_id" gorm:"type:bigint;comment:大BID"`
	ShopId    int          `json:"shop_id" gorm:"type:bigint;comment:关联客户"`
	LineId    int          `json:"line_id" gorm:"type:bigint;comment:线路ID"`
	Line      string       `json:"line" gorm:"index;size:16;comment:路线名称"`
	GoodsId   int          `gorm:"index;comment:商品表ID"`
	Status    int          `json:"status" gorm:"type:bigint;default:1;comment:配送状态"`
	Money     float64      `json:"money" gorm:"type:double;comment:下单总金额"`
	Number    int          `json:"number" gorm:"type:bigint;comment:下单产品数量"`
	Pay       int          `gorm:"type:tinyint(1);default:0;index;comment:支付方式,0:线上,1:线下"`
	PayStatus int          `gorm:"type:tinyint(1);default:1;index;comment:支付状态,0:未付款,1:已付款 2:线下付款，3:下线付款已收款"`
	CycleTime models.XTime `json:"cycle_time" gorm:"type:date;comment:计算的配送时间"`
	CycleStr string       `json:"cycle_str" gorm:"index;size:16;comment:配送周期文案"`
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
	SpecsName string        `gorm:"size:30;comment:规格名称"`
	Unit      string        `json:"unit" gorm:"type:varchar(8);comment:单位"`
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

// todo:存放一些无关紧要的非必要首次查询的数据
type OrderExtend struct {
	models.Model

	OrderId   int           `json:"orderId" gorm:"type:bigint(20);comment:关联订单ID"`
	CreatedAt models.XTime  `json:"created_at" gorm:"comment:创建时间"`
	DeletedAt *models.XTime `json:"-" gorm:"index;comment:删除时间"`
	Desc      string        `json:"desc" gorm:"type:varchar(35);comment:描述信息"`
	Driver    string        `gorm:"size:12;comment:司机名称"`
	Phone     string        `gorm:"size:11;comment:联系手机号"`
	Source    int           `gorm:"type:tinyint(1);default:0;index;comment:订单来源,0:客户下单,1:代客下单"`
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
	Name      string       `gorm:"size:12;comment:下单周期日期名称"`
	Uid       string       `gorm:"type:varchar(4);comment:周期名称都是天,防止一天可能多个不同周期的配置,加个标识区分周期"`
	StartTime models.XTime `gorm:"comment:此周期,下单周期开始时间"`
	EndTime   models.XTime `gorm:"comment:此周期,下单周期结束时间"`
	CycleTime models.XTime `json:"cycle_time" gorm:"type:date;comment:计算的配送时间"`
	CycleStr  string       `json:"cycle_str" gorm:"index;size:14;comment:配送时间的文案"`
	SoldMoney float64      `gorm:"comment:销售总额"`
	GoodsAll  int      `gorm:"comment:商品总数"`
	ShopCount int          `gorm:"type:tinyint(3);comment:客户总数"`
}

func (OrderCycleList) TableName() string {
	return "order_cycle_list"
}
