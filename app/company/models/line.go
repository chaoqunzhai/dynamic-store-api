package models

import (
	"go-admin/common/models"
)

type Line struct {
	models.Model

	Layer    int    `json:"layer" gorm:"type:tinyint(4);comment:排序"`
	Enable   bool   `json:"enable" gorm:"type:tinyint(1);comment:开关"`
	Desc     string `json:"desc" gorm:"type:varchar(25);comment:描述信息"`
	CId      int    `json:"cId" gorm:"type:bigint(20);comment:大BID"`
	Name     string `json:"name" gorm:"type:varchar(16);comment:路线名称"`
	DriverId int    `json:"driver_id" gorm:"type:bigint(20);comment:关联司机"`
	DriverName string  `json:"driver_name" gorm:"-"`
	ExpiresTime models.XTime  `json:"expires_time" gorm:"路线到期时间"`
	ShopCount int64 `json:"shop_count" gorm:"-"`
	models.ModelTime
	models.ControlBy
}

func (Line) TableName() string {
	return "line"
}

func (e *Line) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *Line) GetId() interface{} {
	return e.Id
}
