package models

//微信支付配置
type WeChatPay struct {
	BigBRichGlobal
	MchId string `json:"mch_id" gorm:"index;size:11;comment:商户号"`
	ApiV2 string `json:"api_v2" gorm:"size:50;"`
	ApiV3 string `json:"api_v3" gorm:"size:50;"`
	Refund bool `json:"refund" gorm:"comment:支持退款;"`
	CertPath string `json:"cert_path" gorm:"size:50;comment:支付证书cert路径"`
	KeyPath string  `json:"KeyPath" gorm:"size:50;comment:支付证书key路径"`
}
func (WeChatPay) TableName() string {
	return "company_pay_wechat"
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
	Ali bool `json:"pay_ali" gorm:"size:1;comment:是否开启阿里支付"`
	WeChat bool `json:"we_chat" gorm:"size:1;comment:是否开启微信支付"`
	Credit bool `json:"credit" gorm:"size:1;comment:支持授信减扣"`
}

func (PayCnf) TableName() string {
	return "company_pay_cnf"
}