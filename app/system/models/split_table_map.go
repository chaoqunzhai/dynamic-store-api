package models

import (
	"go-admin/common/models"
)

type SplitTableMap struct {
	models.Model

	Layer  int    `json:"layer" gorm:"type:tinyint;comment:排序"`
	Enable bool   `json:"enable" gorm:"type:tinyint(1);comment:开关"`
	CId    int    `json:"cId" gorm:"type:bigint;comment:公司ID"`
	Type   int    `json:"type" gorm:"type:bigint;comment:映射表的类型"`
	Name   string `json:"name" gorm:"type:varchar(60);comment:对应表的名称"`
	Desc   string `json:"desc" gorm:"type:varchar(30);comment:对应表的名称"`
	models.ModelTime
	models.ControlBy
}

func (SplitTableMap) TableName() string {
	return "split_table_map"
}

func (e *SplitTableMap) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *SplitTableMap) GetId() interface{} {
	return e.Id
}
