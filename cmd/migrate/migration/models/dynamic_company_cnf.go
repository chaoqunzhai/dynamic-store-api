package models

//大B的配置信息表

import "time"

// 大B的小程序信息配置表
type CompanyWeAppCnf struct {
	BigBRichGlobal
	AppId        string `json:"app_id" gorm:"size:20;comment:小程序的appid"`
	Secret       string `json:"secret" gorm:"size:35;comment:私钥"`
	MchID        string `json:"mch_id" gorm:"size:10;comment:商户号"`
	MchKey       string `json:"mch_key" gorm:"size:10;comment:商户APIv3密钥"`
	SerialNumber string `json:"serial_number" gorm:"size:50;comment:证书序列号"`
	PrivateKey   string `json:"private_key" gorm:"size:50;comment:商户私钥文件路径"`
	Certificate  string `json:"certificate" gorm:"size:50;comment:平台证书文件路径"`
}

func (CompanyWeAppCnf) TableName() string {
	return "company_weapp_cnf"
}

// 大B的小程序注册登录方式
type CompanyRegisterCnf struct {
	BigBRichGlobal
	Login    string `json:"login" gorm:"size:30;comment:登录方式"`
	Register string `json:"register" gorm:"size:30;comment:注册方式"`
}

func (CompanyRegisterCnf) TableName() string {
	return "company_register_cnf"
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

// todo:大B配置表,一般用于配置一些限制,创建资源的配置,
// 默认是读取global的配置,读取这个表配置,如果有使用配置表数据
type CompanyQuotaCnf struct {
	BigBRichGlobal
	Key         string `gorm:"size:16;comment:配置Key"`
	Number      int    `gorm:"comment:限制次数"`
	ExtendValue string `gorm:"size:20;comment:扩展的Value"`
}

func (CompanyQuotaCnf) TableName() string {
	return "company_quota_cnf"
}

// todo:大B短信可用条数配置
type CompanyEmsQuotaCnf struct {
	BigBRichGlobal
	Available int `gorm:"comment:可用次数"`
}

func (CompanyEmsQuotaCnf) TableName() string {
	return "company_ems_cnf"
}

// todo:大B短信充值记录
type CompanyEmsQuotaCnfLog struct {
	BigBRichGlobal
	Number int     `gorm:"comment:充值条数"`
	Money  float32 `gorm:"comment:费用"`
}

func (CompanyEmsQuotaCnfLog) TableName() string {
	return "company_ems_cnf_log"
}
