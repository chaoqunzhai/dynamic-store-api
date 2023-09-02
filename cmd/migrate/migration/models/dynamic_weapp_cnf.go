package models

//todo:大B的商场H5|小程序配置
//配置的菜单表

type WeAppGlobalNavCnf struct {
	Model
	Enable           bool   `json:"-" gorm:"type:tinyint(1);default:1;"`
	UserEnable       bool   `gorm:"-" json:"user_enable"` //只是用来渲染,不做创建字段
	IconPath         string `gorm:"size:50;" json:"icon_path"`
	SelectedIconPath string `gorm:"size:50;" json:"selected_icon_path"`
	Text             string `gorm:"size:50;" json:"text"`
	Name             string `gorm:"size:50;" json:"name"`
	WapUrl           string `gorm:"size:50;" json:"wap_url"`
	IconClass        string `gorm:"size:50;" json:"icon_class"`
}

func (WeAppGlobalNavCnf) TableName() string {
	return "weapp_global_nav_cnf"
}

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
	Type  int    `json:"type" gorm:"type:tinyint(1);comment:类型"` //0:登录  1:注册
	Value string `json:"login" gorm:"size:12;comment:登录方式"`      //username,mobile,wechat 代表用户名,手机号,微信
}

func (CompanyRegisterCnf) TableName() string {
	return "company_register_cnf"
}

// 大B底栏菜单配置
type CompanyNavCnf struct {
	BigBMiniGlobal
	GId int `gorm:"index;type:tinyint(1);comment:关联的菜单配置ID"`
}

func (CompanyNavCnf) TableName() string {
	return "company_nav_cnf"
}

// 小程序个人中心,快捷导航配置
type WeAppQuickTools struct {
	Model
	Name     string `json:"name"  gorm:"size:30;comment:名称"`
	UserEnable       bool   `gorm:"-" json:"user_enable"` //只是用来渲染,不做创建字段
	Enable   bool   `json:"-" gorm:"type:tinyint(1);default:1;comment:是否开启"`
	ImageUrl string `json:"image_url" gorm:"size:60;comment:图片路径"`
	WapUrl   string `gorm:"size:50;" json:"wap_url"`
}

func (WeAppQuickTools) TableName() string {
	return "weapp_quick_tools"
}

//TODO:我得页面中 常用工具 导航配置关联配置
type CompanyQuickTools struct {
	BigBMiniGlobal
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

//TODO:配置一些扩展的客户端一些
type WeAppExtendCnf struct {
	BigBRichGlobal
	DetailAddCart string `json:"detail_add_cart" gorm:"size:10;comment:加入购物车按钮的重命名"` //详情页面中,加入购物车的文案
	DetailAddCartColor string `json:"detail_add_cart_color" gorm:"size:8;comment:加入购物车按钮的颜色"` //详情页面中,加入购物车的颜色
	DetailAddCartShow bool `json:"detail_add_cart_show" gorm:"default:1"`
	DetailByNow string `json:"detail_by_now" gorm:"size:10;comment:立即购买的重命名"` //详情页面中,立即购买的文案
	DetailByNowColor string `json:"detail_by_now_color" gorm:"size:8;comment:立即购买颜色"` //详情页面中,立即购买的文案
	DetailByNowShow bool `json:"detail_by_now_show" gorm:"default:1"`
}
func (WeAppExtendCnf) TableName() string {
	return "company_member_index"
}

//小程序我的订单配置

//小程序
