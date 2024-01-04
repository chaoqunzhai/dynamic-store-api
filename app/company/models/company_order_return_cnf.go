package models

import (
     
     
     "go-admin/common/models"

)

type CompanyOrderReturnCnf struct {
    models.Model
    
    Layer string `json:"layer" gorm:"type:tinyint(4);comment:排序"` 
    Enable bool `json:"enable" gorm:"type:tinyint(1);comment:开关"`
    Desc string `json:"desc" gorm:"type:varchar(35);comment:描述信息"` 
    CId int `json:"cId" gorm:"type:bigint(20);comment:大BID"`
    Value string `json:"value" gorm:"type:varchar(15);comment:配送文案"` 
    Cost float64 `json:"cost" gorm:"type:double;comment:配送费用"`
    models.ModelTime
    models.ControlBy
}

func (CompanyOrderReturnCnf) TableName() string {
    return "company_order_return_cnf"
}

func (e *CompanyOrderReturnCnf) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *CompanyOrderReturnCnf) GetId() interface{} {
	return e.Id
}