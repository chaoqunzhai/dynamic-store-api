package models

import "time"

// todo:大B信息
type Company struct {
	RichGlobal
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
}

func (Company) TableName() string {
	return "company"
}

// todo:大B续费信息
type CompanyRenewalTimeLog struct {
	Model
	ModelTime
	CreateBy       int       `json:"create_by" gorm:"index;comment:创建者"`
	CId            int       `json:"-" gorm:"index;comment:公司(大B)ID"`
	Money          float64   `json:"money" gorm:"comment:续费金额"`
	Desc           string    `json:"desc" gorm:"size:50;comment:描述信息"`
	ExpirationTime time.Time `json:"expiration_time" gorm:"comment:续费到期时间"`
}

func (CompanyRenewalTimeLog) TableName() string {
	return "company_renewal_time_log"
}

// todo:大B设置的VIP等级
type GradeVip struct {
	BigBRichGlobal
	Name     string  `gorm:"size:30;comment:等级名称"`
	Weight   int     `gorm:"type:tinyint(1);default:1;comment:权重,从小到大"`
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

// todo:大B配置表,一般用于配置一些限制,
// 默认是读取global的配置,读取这个表配置,如果有使用配置表数据
type CompanyCnf struct {
	BigBRichGlobal
	Key         string `gorm:"size:20;comment:配置Key"`
	Value       string `gorm:"size:20;comment:限制的Value"`
	ExtendValue string `gorm:"size:20;comment:扩展的Value"`
}

func (CompanyCnf) TableName() string {
	return "company_cnf"
}

//todo:大B店铺设计类型保存

type CompanyCategory struct {
	BigBRichGlobal
	Type int `gorm:"type:tinyint(1);default:1;comment:模板类型"`
}

func (CompanyCategory) TableName() string {
	return "company_category"
}