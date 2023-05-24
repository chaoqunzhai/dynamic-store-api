package models

import (
     


	"go-admin/common/models"

)

type CompanyRole struct {
    models.Model
    
    Name string `json:"name" gorm:"type:varchar(30);comment:Name"` 
    Enable string `json:"enable" gorm:"type:bigint(20);comment:Enable"` 
    Sort string `json:"sort" gorm:"type:bigint(20);comment:Sort"` 
    Remark string `json:"remark" gorm:"type:varchar(50);comment:Remark"` 
    Admin string `json:"admin" gorm:"type:tinyint(1);comment:Admin"` 
    models.ModelTime
    models.ControlBy
}

func (CompanyRole) TableName() string {
    return "company_role"
}

func (e *CompanyRole) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *CompanyRole) GetId() interface{} {
	return e.Id
}