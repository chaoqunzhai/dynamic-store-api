package models

import (
	"go-admin/common/models"
)

type GradeVip struct {
	models.Model

	Layer    string `json:"layer" gorm:"type:tinyint(4);comment:排序"`
	Enable   string `json:"enable" gorm:"type:tinyint(1);comment:开关"`
	Desc     string `json:"desc" gorm:"type:varchar(25);comment:描述信息"`
	CId      string `json:"cId" gorm:"type:bigint(20);comment:大BID"`
	Name     string `json:"name" gorm:"type:varchar(30);comment:等级名称"`
	Weight   string `json:"weight" gorm:"type:bigint(20);comment:权重,从小到大"`
	Discount string `json:"discount" gorm:"type:float;comment:折扣"`
	Upgrade  string `json:"upgrade" gorm:"type:bigint(20);comment:升级条件,满多少金额,自动升级Weight+1"`
	models.ModelTime
	models.ControlBy
}

func (GradeVip) TableName() string {
	return "grade_vip"
}

func (e *GradeVip) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *GradeVip) GetId() interface{} {
	return e.Id
}
