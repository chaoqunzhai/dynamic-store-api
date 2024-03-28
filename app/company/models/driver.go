package models

import (
     


	"go-admin/common/models"

)

type Driver struct {
    models.Model
    
    Layer int `json:"layer" gorm:"type:tinyint(4);comment:排序"`
    Enable bool `json:"enable" gorm:"type:tinyint(1);comment:开关"`
    Desc string `json:"desc" gorm:"type:varchar(50);comment:备注信息"` 
    CId int `json:"cId" gorm:"type:bigint(20);comment:大BID"`
    UserId int `json:"userId" gorm:"type:bigint(20);comment:关联的用户ID"`
    Name string `json:"name" gorm:"type:varchar(12);comment:司机名称"` 
    Phone string `json:"phone" gorm:"type:varchar(11);comment:联系手机号"`
    LineName string `json:"line_name" gorm:"-"`
    Disable bool `json:"disable" gorm:"-"`
    Password string `json:"password" gorm:"-"`
	models.ModelTime
    models.ControlBy
}

func (Driver) TableName() string {
    return "driver"
}

func (e *Driver) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *Driver) GetId() interface{} {
	return e.Id
}