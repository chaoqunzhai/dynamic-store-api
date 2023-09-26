package models

import (
     
     
     
     
     
     "go-admin/common/models"

)

type CompanyDebitCard struct {
    models.Model
    
    Layer string `json:"layer" gorm:"type:tinyint(4);comment:排序"` 
    Enable bool `json:"enable" gorm:"type:tinyint(1);comment:开关"`
    Desc string `json:"desc" gorm:"type:varchar(35);comment:描述信息"`
    CId int `json:"-" gorm:"type:bigint(20);comment:大BID"`
    Bank string `json:"bank" gorm:"type:varchar(20);comment:银行名称"`
    Name string `json:"name" gorm:"type:varchar(20);comment:持卡人名称"`
    CardNumber string `json:"card_number" gorm:"type:varchar(25);comment:银行卡号"`
	BankName string `json:"bank_name" gorm:"type:varchar(15);comment:开户行"`
    models.ModelTime
    models.ControlBy
}

func (CompanyDebitCard) TableName() string {
    return "company_debit_card"
}

func (e *CompanyDebitCard) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *CompanyDebitCard) GetId() interface{} {
	return e.Id
}