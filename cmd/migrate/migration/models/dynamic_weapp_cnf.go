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

//	{
//	               "title": "邀请有礼",
//	               "imageUrl": "../../static/member/default_memberrecommend.png",
//	               "iconType": "img",
//	               "style": "",
//	               "link": {
//	                   "name": "MEMBER_RECOMMEND",
//	                   "title": "邀请有礼",
//	                   "wap_url": "/pages_tool/member/invite_friends",
//	                   "parent": "MARKETING_LINK"
//	               },
//	               "label": {
//	                   "control": false,
//	                   "text": "热门",
//	                   "textColor": "#FFFFFF",
//	                   "bgColorStart": "#F83287",
//	                   "bgColorEnd": "#FE3423"
//	               },
//	               "iconfont": {
//	                   "value": "",
//	                   "color": ""
//	               },
//	               "id": "1h34nmfisge80"
//	           },
//
// 小程序个人中心,快捷导航配置
type WeAppQuickTools struct {
	Model
	Name     string `json:"name"  gorm:"size:30;comment:名称"`
	Enable   bool   `json:"-" gorm:"type:tinyint(1);default:1;comment:是否开启"`
	ImageUrl string `json:"image_url" gorm:"size:30;comment:图片路径"`
	WapUrl   string `gorm:"size:50;" json:"wap_url"`
}

func (WeAppQuickTools) TableName() string {
	return "weapp_quick_tools"
}

// 大B和快捷导航配置关联配置
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
