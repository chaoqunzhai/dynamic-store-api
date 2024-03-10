package models

import (
	"gorm.io/gorm"
	"time"
)
type CompanyRegisterRule struct {
	Model
	CreatedAt time.Time      `json:"createdAt" gorm:"comment:创建时间"`
	CId int `gorm:"index;comment:大BID"`
	UserRule int `gorm:"default:1;index;comment:1:需要审核通过才可以登录 2:不需要审核直接注册"`
	Text string `json:"text" gorm:"size:20;comment:消息内容"`

}
func (CompanyRegisterRule) TableName() string {
	return "company_register_rule"
}


//审核表
type CompanyRegisterUserVerify struct {
	Model
	CreatedAt time.Time      `json:"created_at" gorm:"comment:创建时间"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间"`
	CId int `json:"-" gorm:"index;comment:大BID"`
	AdoptTime *time.Time `json:"adopt_time" gorm:"通过时间"`
	AdoptUser string `json:"adopt_user" gorm:"size:11;comment:审批人"`
	Salesman int `json:"salesman" gorm:"comment:推广业务员"`
	Source string `json:"source" gorm:"size:6;comment:注册方式 user | mobile"`
	Value string `json:"value" gorm:"size:15;comment:注册数据,用户名或者手机号"`
	Phone string `json:"phone" gorm:"size:15;comment:注册数据,手机号"`
	AppTypeName string `json:"app_type_name" gorm:"size:20;comment:注册来源例如H5,WECHAT,ALI等"`
	Status int `json:"status" gorm:"default:0;index;comment:0:审核中, 1:通过 -1:驳回"`
	Info string `json:"info" gorm:"size:10;comment:备注"`
}
func (CompanyRegisterUserVerify) TableName() string {
	return "company_register_user_verify"
}
