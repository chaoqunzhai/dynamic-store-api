package models
import "time"
type CompanyRegisterRule struct {
	Model
	CreatedAt time.Time      `json:"createdAt" gorm:"comment:创建时间"`
	CId int `gorm:"index;comment:大BID"`
	UserRule int `gorm:"default:1;index;comment:1:需要审核通过才可以登录 2:不需要审核直接注册"`

}
func (CompanyRegisterRule) TableName() string {
	return "company_register_rule"
}


//审核表

type CompanyRegisterUserVerify struct {
	Model
	ControlBy
	ModelTime
	CId int `gorm:"index;comment:大BID"`
	Source string `gorm:"size:6;comment:注册方式 user | mobile"`
	Value string `gorm:"size:15;comment:注册数据,用户名或者手机号"`
	Status int `gorm:"default:0;index;comment:0:不通过 1:通过"`
}
func (CompanyRegisterUserVerify) TableName() string {
	return "company_register_user_verify"
}
