package models

import (
	"gorm.io/gorm"
	"time"
)

// 所有的软关联 必须是
// 1.字段名称Id 组合
// 2.必须增加index索引


// todo:小B商家信息,也就是客户
type Shop struct {
	BigBRichGlobal
	UserId    int     `gorm:"index;comment:管理员ID"`
	Name      string  `gorm:"size:30;comment:小B名称"`
	Phone     string  `gorm:"size:11;comment:联系手机号"`
	UserName  string  `gorm:"size:20;comment:小B负责人名称"`
	Address   string  `gorm:"size:200;comment:小B收货地址"`
	Longitude float64 //经度
	Latitude  float64 //纬度
	Image     string  `gorm:"size:80;comment:图片"`
	Salesman int `json:"salesman" gorm:"index;comment:推荐人"`
	//给小B打标签
	Tag      []ShopTag `json:"tag" gorm:"many2many:shop_mark_tag;foreignKey:id;joinForeignKey:shop_id;references:id;joinReferences:tag_id;"`
	LineId   int       `gorm:"index;comment:归属配送路线"`
	Amount   float64   `gorm:"comment:剩余金额"`
	Integral int       `gorm:"comment:可用积分"`
}

func (Shop) TableName() string {
	return "shop"
}

// todo: 小B充值记录
type ShopRechargeLog struct {
	Model
	CreateBy  int            `json:"createBy" gorm:"index;comment:创建者"`
	CreatedAt time.Time      `json:"createdAt" gorm:"comment:创建时间"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间"`
	ShopId    int            `gorm:"index;comment:小BID"`
	Uuid      string         `json:"uuid" gorm:"index;size:20;comment:订单号"` //订单号
	Source    string         `json:"source" gorm:"size:16;comment:充值方式"`
	Money     float64        `gorm:"comment:支付金额"`
	GiveMoney float64        `gorm:"comment:赠送金额"`
	PayStatus bool           `gorm:"comment:支付状态"`
	PayTime   time.Time      `gorm:"comment:支付时间"`
}

func (ShopRechargeLog) TableName() string {
	return "shop_recharge_log"
}

// todo:小B余额变动明细
type ShopBalanceLog struct {
	Model
	CreateBy  int            `json:"createBy" gorm:"index;comment:创建者"`
	CreatedAt time.Time      `json:"createdAt" gorm:"comment:创建时间"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间"`
	ShopId    int            `gorm:"index;comment:小BID"`
	Money     float64        `gorm:"comment:变动金额"`
	Scene     string         ` gorm:"size:30;comment:变动场景"`
	Desc string   ` gorm:"size:50;comment:描述/说明"`
}

func (ShopBalanceLog) TableName() string {
	return "shop_balance_log"
}

// todo:客户每次订单的统计日志
// 专门用来数据统计
type ShopOrderRecord struct {
	Model
	CreateBy  int            `json:"createBy" gorm:"index;comment:创建者"`
	CreatedAt time.Time      `json:"createdAt" gorm:"comment:创建时间"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间"`
	ShopId    int            `gorm:"index;comment:关联的小B客户"`
	ShopName  string         `gorm:"size:30;comment:客户名称"`
	Money     float64        `gorm:"comment:订单金额"`
	Number    int            `gorm:"comment:订单量"`
}

func (ShopOrderRecord) TableName() string {
	return "shop_order_record"
}

// todo:每次订单统计关联的具体订单
type ShopOrderBindRecord struct {
	Model
	CreateBy  int            `json:"createBy" gorm:"index;comment:创建者"`
	CreatedAt time.Time      `json:"createdAt" gorm:"comment:创建时间"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间"`
	ShopId    int            `gorm:"index;comment:关联的小B客户"`
	RecordId  int            `gorm:"index;comment:每次记录的总ID"`
	OrderId   int            `gorm:"index;comment:订单ID"`
}

func (ShopOrderBindRecord) TableName() string {
	return "shop_order_bind_record"
}

// todo:客户标签
type ShopTag struct {
	BigBRichGlobal
	Name string `gorm:"index;size:35;comment:客户标签名称"`
}

func (ShopTag) TableName() string {
	return "shop_tag"
}
