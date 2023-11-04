package models

//大B的业务信息
import (
	"go-admin/common/models"
	"time"
)

// todo:大B信息
type Company struct {
	RichGlobal
	LeaderId uint `json:"leader_id"`
	Name           string    `gorm:"index;size:30;comment:公司(大B)名称"`
	Phone          string    `gorm:"size:11;comment:负责人联系手机号"`
	UserName       string    `gorm:"size:20;comment:大B负责人名称"`
	ShopName       string    `gorm:"size:50;comment:自定义大B系统名称"`
	Address        string    `gorm:"size:120;comment:大B地址位置"`
	Longitude      float64   //经度
	Latitude       float64   //纬度
	Image          string    `gorm:"size:80;comment:logo图片"`
	RenewalTime    time.Time `json:"renewal_time" gorm:"comment:续费时间"`
	ExpirationTime time.Time `json:"expiration_time" gorm:"comment:到期时间"`
	LoginTime time.Time     `json:"login_time" gorm:"type:datetime(3);comment:登录时间"`
}

func (Company) TableName() string {
	return "company"
}

// todo:大B设置的VIP等级
type GradeVip struct {
	BigBRichGlobal
	Name     string  `gorm:"size:30;comment:等级名称"`
	Weight   int     `gorm:"type:tinyint(1);default:1;comment:等级,从小到大"`
	Discount float32 `gorm:"comment:折扣"`
	Upgrade  int     `gorm:"default:0;comment:升级条件,满多少金额,自动升级Weight+1"`
}

func (GradeVip) TableName() string {
	return "grade_vip"
}

// todo:路线信息,被司机关联
// 每个路线就是一个配送员
// 大B下有很多路线
type Line struct {
	BigBRichGlobal
	Name     string `gorm:"index;size:16;comment:路线名称"`
	RenewalTime    models.XTime     `json:"renewal_time" gorm:"type:datetime(3);comment:续费时间"`
	ExpirationTime models.XTime      `json:"expiration_time" gorm:"type:datetime(3);comment:到期时间"`
	DriverId int    `gorm:"index;comment:关联司机"`
}

func (Line) TableName() string {
	return "line"
}

// todo:司机信息
type Driver struct {
	BigBRichGlobal
	UserId int    `gorm:"index;comment:关联的用户ID"`
	Name   string `gorm:"size:12;comment:司机名称"`
	Phone  string `gorm:"index;size:11;comment:联系手机号"`
	Desc   string `gorm:"size:50;comment:备注信息"`
}

func (Driver) TableName() string {
	return "driver"
}

// todo:大B物流配送方式,负责一些特殊的自定义配置
type CompanyExpress struct {
	BigBRichGlobal
	Type int `gorm:"type:tinyint(1);default:2;comment:物流类型"`
}

func (CompanyExpress) TableName() string {
	return "company_express"
}

// 门店
type CompanyExpressStore struct {
	BigBRichGlobal
	Name    string `gorm:"size:20;comment:门店名称"`
	Phone string `json:"phone" gorm:"size:11;comment:电话"`
	Address string `json:"address" gorm:"size:120;comment:大B门店地址位置"`
	Start   string `gorm:"size:12;comment:营业开始时间"`
	End     string `gorm:"size:12;comment:营业结束时间"`
}

func (CompanyExpressStore) TableName() string {
	return "company_express_store"
}

// todo: 运费配置
type CompanyFreight struct {
	BigBRichGlobal
	Type         int `gorm:"type:tinyint(1);default:2;comment:物流类型"`
	QuotaMoney   int `gorm:"comment:达到多少钱可以免运费"`
	StartMoney   int `gorm:"comment:最低的起送金额"`
	FreightMoney int `gorm:"comment:没有达到QuotaMoney,运费多少钱"`
}

func (CompanyFreight) TableName() string {
	return "company_freight"
}
