package models

//todo:大B的商场H5|小程序配置

// todo:大B店铺设计类型保存
type CompanyCategory struct {
	BigBRichGlobal
	Type int `gorm:"type:tinyint(1);default:1;comment:模板类型"`
}

func (CompanyCategory) TableName() string {
	return "company_category"
}

// 大B的小程序信息配置表
type CompanyWeAppCnf struct {
	BigBRichGlobal
	AppId        string `json:"app_id" gorm:"size:20;comment:小程序的appid"`
	Secret       string `json:"secret" gorm:"size:35;comment:小程序AppSecret"`
	MchID        string `json:"mch_id" gorm:"size:10;comment:商户号"`
	MchKey       string `json:"mch_key" gorm:"size:10;comment:商户APIKEY密钥"`
	SerialNumber string `json:"serial_number" gorm:"size:50;comment:证书序列号"`
	CertPem      string `json:"cert_pem" gorm:"size:100;comment:证书文件cert"`
	KeyPem       string `json:"key_pem" gorm:"size:100;comment:证书文件key"`
}

func (CompanyWeAppCnf) TableName() string {
	return "company_weapp_cnf"
}

// 大B小程序注册登录方式
type CompanyRegisterCnf struct {
	BigBRichGlobal
	Type int `json:"type" gorm:"type:tinyint(1);comment:类型"` //0:登录  1:注册
	Value    string `json:"login" gorm:"size:12;comment:登录方式"` //username,mobile,wechat 代表用户名,手机号,微信
}

func (CompanyRegisterCnf) TableName() string {
	return "company_register_cnf"
}

// 小程序个人中心,快捷导航配置
type WeAppQuickTools struct {
	RichGlobal
	Name string `json:"name"  gorm:"size:30;comment:导航名称"`
	Icon string `json:"icon" gorm:"size:30;comment:图片路径"`
	Url  string `json:"url" gorm:"size:30;comment:跳转路径"`
}

func (WeAppQuickTools) TableName() string {
	return "we_app_quick_tools"
}

// 大B和快照导航关联的配置
type CompanyQuickTools struct {
	Model
	Cid     int `gorm:"index;comment:关联的大BID"`
	QuickId int `gorm:"index;comment:关联的导航配置"`
}

func (CompanyQuickTools) TableName() string {
	return "company_quick_tools"
}

// 小程序个人中心配置
type CompanyMemberIndex struct {
	Model
	Type         int    `gorm:"type:tinyint(0);default:1;comment:主题类型.0是默认,1是自定义"`
	BgColorStart string `json:"bg_color_start" gorm:"size:10;comment:主题开始颜色"`
	BgColorEnd   string `json:"bg_color_end" gorm:"size:10;comment:主题结尾颜色"`
}

func (CompanyMemberIndex) TableName() string {
	return "company_member_index"
}

//小程序我的订单配置

//小程序