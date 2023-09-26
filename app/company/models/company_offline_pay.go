package models

import (
     
     
     
     "go-admin/common/models"

)

type CompanyOfflinePay struct {
    models.Model
    
    CId int `json:"-" gorm:"type:bigint(20);comment:大BID"`
    Name string `json:"name" gorm:"type:varchar(12);comment:线下支付名称"` 
    Layer string `json:"layer" gorm:"type:tinyint(4);comment:排序"` 
    Enable bool `json:"enable" gorm:"type:tinyint(1);comment:开关"`
    Desc string `json:"desc" gorm:"type:varchar(35);comment:描述信息"` 
    models.ModelTime
    models.ControlBy
}

func (CompanyOfflinePay) TableName() string {
    return "company_offline_pay"
}

func (e *CompanyOfflinePay) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *CompanyOfflinePay) GetId() interface{} {
	return e.Id
}