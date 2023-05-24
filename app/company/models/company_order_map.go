package models

import (
     

	"go-admin/common/models"

)

type CompanyOrderMap struct {
    models.Model
    
    Layer string `json:"layer" gorm:"type:tinyint(4);comment:排序"` 
    Enable string `json:"enable" gorm:"type:tinyint(1);comment:开关"` 
    CId string `json:"cId" gorm:"type:bigint(20);comment:公司ID"` 
    Type string `json:"type" gorm:"type:bigint(20);comment:映射表的类型"` 
    OrderTable string `json:"orderTable" gorm:"type:varchar(50);comment:对应表的名称"` 
    models.ModelTime
    models.ControlBy
}

func (CompanyOrderMap) TableName() string {
    return "company_order_map"
}

func (e *CompanyOrderMap) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *CompanyOrderMap) GetId() interface{} {
	return e.Id
}