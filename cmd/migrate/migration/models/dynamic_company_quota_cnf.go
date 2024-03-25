package models

//大B的配置信息表

import (
	"gorm.io/gorm"
	"time"
)


// todo:大B续费信息
type CompanyRenewalTimeLog struct {
	MiniLog
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
	LineId int `json:"line_id" gorm:"comment:线路ID"`
	Money  float64 `gorm:"comment:费用"`
	BuyType     int    `json:"buy_type" gorm:"size:1;index;"` //排序,默认是0:购买 1:续费
	ExpirationTime time.Time `json:"expiration_time" gorm:"comment:续费到期时间"`
	ExpirationStr string `json:"expiration_str" gorm:"size:12;comment:续费时长,文案"`
}

func (CompanyLineCnfLog) TableName() string {
	return "company_line_cnf_log"
}
// todo:大B短信可用条数配置
type CompanySmsQuotaCnf struct {
	BigBRichGlobal
	Available int `json:"available" gorm:"comment:可用次数"`
	Record bool `json:"record" gorm:"comment:是否开启消费记录"`
	OrderNotice bool `json:"order_notice" gorm:"default:false;comment:是否开启下订单短信通知"`
}

func (CompanySmsQuotaCnf) TableName() string {
	return "company_sms_cnf"
}

// todo:大B短信充值记录
type CompanySmsQuotaCnfLog struct {
	MiniLog
	Number int     `gorm:"comment:充值条数"`
	Money  float64 `gorm:"comment:费用"`
}

func (CompanySmsQuotaCnfLog) TableName() string {
	return "company_sms_cnf_log"
}

//todo:大B短信消费记录, 这个是开关,如果大B需要记录就
type CompanySmsRecordLog struct {
	Model
	CId            int       `json:"c_id" gorm:"index;comment:公司(大B)ID"`
	CreatedAt time.Time      `json:"createdAt" gorm:"comment:创建时间"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间"`
	Status bool `json:"status" gorm:"default:true;comment:发送状态"`
	Source string `gorm:"size:18;comment:发送源头"`
	Phone string `gorm:"size:11;comment:手机号"`
	Code string `gorm:"size:6;comment:验证码"`
	Msg string `json:"msg" gorm:"size:80;comment:其他情况下发送的内容"`
}
func (CompanySmsRecordLog) TableName() string {
	return "company_sms_record_log"
}