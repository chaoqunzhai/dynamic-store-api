package models

import (
	"time"

	"go-admin/common/models"
)

type Company struct {
	models.Model
	LeaderId uint `json:"leader_id"`
	Enterprise string `json:"enterprise" gorm:"size:20;comment:企业名称"`
	Filings string `json:"filings" gorm:"size:20;comment:备案号"`
	NewPhone string   `gorm:"size:11;comment:联系手机号"`
	Name           string        `json:"name" gorm:"type:varchar(30);comment:公司(大B)名称"`
	Address        string        `json:"address" gorm:"type:varchar(155);comment:大B地址位置"`
	ShopName       string        `json:"shop_name" gorm:"type:varchar(50);comment:店铺名称"`

	Layer          int           `json:"layer" gorm:"type:tinyint(4);comment:排序"`
	Enable         bool          `json:"enable" gorm:"type:tinyint(1);comment:开关"`
	Desc           string        `json:"desc" gorm:"type:varchar(25);comment:描述信息"`

	Phone          string        `json:"phone" gorm:"type:varchar(11);comment:负责人联系手机号"`
	UserName       string        `json:"user_name" gorm:"type:varchar(20);comment:大B负责人名称"`

	Longitude      float64       `json:"longitude" gorm:"type:double;comment:Longitude"`
	Latitude       float64       `json:"latitude" gorm:"type:double;comment:Latitude"`
	Image          string        `json:"image" gorm:"type:varchar(80);comment:logo图片"`
	RenewalTime    time.Time     `json:"renewal_time" gorm:"type:datetime(3);comment:续费时间"`
	ExpirationTime time.Time     `json:"expiration_time" gorm:"type:datetime(3);comment:到期时间"`
	LoginTime models.XTime      `json:"login_time" gorm:"type:datetime(3);comment:登录时间"`

	NavList        []interface{} `json:"nav_list" gorm:"-"` //展示大B的菜单
	models.ModelTime
	models.ControlBy
}

func (Company) TableName() string {
	return "company"
}

func (e *Company) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *Company) GetId() interface{} {
	return e.Id
}


type UserApplyPaymentOrder struct {
	models.Model
	models.ModelTime
	models.ControlBy
	Layer int `json:"layer" gorm:"type:tinyint(4);comment:排序"`
	Enable bool `json:"enable" gorm:"type:tinyint(1);comment:开关"`
	Desc string `json:"desc" gorm:"type:varchar(35);comment:描述信息"`
	CId int `json:"-" gorm:"type:bigint(20);comment:大BID"`
	TransferDate models.XTime `json:"transfer_date"`
	Money float64 `json:"money"`
	UseTo int `json:"use_to" gorm:"size:1;comment:用途 0:记录 1:计入余额 2:计入授信额"`
	Status int `json:"status" gorm:"size:1;comment:付款单状态 0:提交中  1:确认到账 2:驳回"`
	Bank string `json:"bank" gorm:"size:20;comment:银行名称"`
	Name string `json:"name" gorm:"size:20;comment:持卡人名称"`
	BankName string `json:"bank_name" gorm:"size:15;comment:开户行"`
	CardNumber string `json:"card_number" gorm:"size:25;comment:银行卡号"`
	ApproveMsg string `json:"approve_msg" gorm:"size:20;comment:大B审批写的内容"`
	UserName string `json:"user_name" gorm:"-"`
}
func (UserApplyPaymentOrder) TableName() string {
	return "company_apply_payment_order"
}