package models

//大B的配置信息表

import "time"


// todo:大B续费信息
type CompanyRenewalTimeLog struct {
	MiniLog
	CId            int       `json:"-" gorm:"index;comment:公司(大B)ID"`
	Money          float64   `json:"money" gorm:"comment:续费金额"`
	ExpirationTime time.Time `json:"expiration_time" gorm:"comment:续费到期时间"`
	ExpirationStr string `json:"expiration_str" gorm:"size:12;comment:续费时长,文案"`
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
//todo: 大B可用线路
type CompanyLineCnf struct {
	BigBRichGlobal
	Number      int    `gorm:"comment:可用的线路"`
}
func (CompanyLineCnf) TableName() string {
	return "company_line_cnf"
}
// todo:大B可用线路充值记录
type CompanyLineCnfLog struct {
	MiniLog
	Number int     `gorm:"comment:购买条数"`
	Money  float32 `gorm:"comment:费用"`
}

func (CompanyLineCnfLog) TableName() string {
	return "company_line_cnf_log"
}
// todo:大B短信可用条数配置
type CompanyEmsQuotaCnf struct {
	BigBRichGlobal
	Available int `gorm:"comment:可用次数"`
	Record bool `gorm:"comment:是否开启消费记录"`
}

func (CompanyEmsQuotaCnf) TableName() string {
	return "company_ems_cnf"
}

// todo:大B短信充值记录
type CompanyEmsQuotaCnfLog struct {
	MiniLog
	Number int     `gorm:"comment:充值条数"`
	Money  float32 `gorm:"comment:费用"`
}

func (CompanyEmsQuotaCnfLog) TableName() string {
	return "company_ems_cnf_log"
}

//todo:大B短信消费记录, 这个是开关,如果大B需要记录就
type CompanyEmsRecordLog struct {
	MiniLog
	CId int `gorm:"index;comment:大BID"`
	Source string `gorm:"size:15;comment:发送源头"`
	Phone string `gorm:"size:11;comment:手机号"`
	Code string `gorm:"size:6;comment:验证码"`
}
func (CompanyEmsRecordLog) TableName() string {
	return "company_ems_record_log"
}