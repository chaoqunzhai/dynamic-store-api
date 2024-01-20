package models

import "time"

type WeChatAppIdCnf struct {
	BigBRichGlobal
	AppId string `json:"app_id" gorm:"index;size:20;comment:应用ID"`
	AppSecret string `json:"app_secret"  gorm:"index;size:50;comment:应用ID"`
}
func (WeChatAppIdCnf) TableName() string {
	return "company_wechat_app_cnf"
}
//微信支付配置
type WeChatOfficialPay struct {
	BigBRichGlobal
	MchId string `json:"mch_id" gorm:"index;size:20;comment:商户号"`
	ApiV2 string `json:"api_v2" gorm:"size:32;"`
	ApiV3 string `json:"api_v3" gorm:"size:32;"`
	Refund bool `json:"refund" gorm:"comment:支持退款;"`
	CertText string `json:"cert_text" gorm:"comment:支付证书cert内容"`
	KeyText string  `json:"key_text" gorm:"comment:支付证书key路径"`
	OfficialAppId string `json:"official_app_id" gorm:"size:20;comment:微信公众号APPID"`
	SerialNumber string `json:"serial_number" gorm:"size:60;comment:证书序列号"`
}
func (WeChatOfficialPay) TableName() string {
	return "company_wechat_official_pay"
}


//支付宝配置

type AliPay struct {
	BigBRichGlobal
	AppId string `json:"app_id" gorm:"index;size:20;comment:应用ID"`
	PrivateKey string `json:"private_key" gorm:"comment:应用私钥"`
	PublicKey string `json:"public_key"  gorm:"comment:应用公钥"`
	AlipayPublicKey string `json:"alipay_public_key" gorm:"comment:支付宝公钥"`
	Refund bool `json:"refund" gorm:"comment:支持退款;"`
}

func (AliPay) TableName() string {
	return "company_pay_ali"
}
//支付的开关配置
type PayCnf struct {
	BigBRichGlobal
	BalanceDeduct bool `json:"balance_deduct" gorm:"size:1;comment:是否开启余额支付"`
	//Ali bool `json:"pay_ali" gorm:"size:1;comment:是否开启阿里支付"`
	//WeChat bool `json:"we_chat" gorm:"size:1;comment:是否开启微信支付"`
	Credit bool `json:"credit" gorm:"size:1;comment:支持授信减扣"`
}

func (PayCnf) TableName() string {
	return "company_pay_cnf"
}

//借记卡管理,

type DebitCard struct {
	BigBRichGlobal
	Bank string `json:"bank" gorm:"size:20;comment:银行名称"`
	Name string `json:"name" gorm:"size:20;comment:持卡人名称"`
	BankName string `json:"bank_name" gorm:"size:15;comment:开户行"`
	CardNumber string `json:"card_number" gorm:"size:25;comment:银行卡号"`
}

func (DebitCard) TableName() string {
	return "company_debit_card"
}


type OfflinePay struct {
	BigBRichGlobal
	CreatedAt time.Time      `json:"createdAt" gorm:"comment:创建时间"`
	Name string `json:"back_name" gorm:"size:12;comment:线下支付名称"`
}
func (OfflinePay) TableName() string {
	return "company_offline_pay"
}

//用户主动提交的付款单

type UserApplyPaymentOrder struct {
	BigBRichGlobal

	TransferDate time.Time `json:"transfer_date"`
	Money float64 `json:"money"`
	UseTo int `json:"use_to" gorm:"size:1;comment:用途 0:记录 1:计入余额 2:计入授信额"`
	Status int `json:"status" gorm:"size:1;comment:付款单状态 0:提交中  1:确认到账 2:有问题"`
	Bank string `json:"bank" gorm:"size:20;comment:银行名称"`
	Name string `json:"name" gorm:"size:20;comment:持卡人名称"`
	BankName string `json:"bank_name" gorm:"size:15;comment:开户行"`
	CardNumber string `json:"card_number" gorm:"size:25;comment:银行卡号"`
	ApproveMsg string `json:"approve_msg" gorm:"size:20;comment:大B审批写的内容"`
}
func (UserApplyPaymentOrder) TableName() string {
	return "company_apply_payment_order"
}