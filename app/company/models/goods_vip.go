package models

import (
	"go-admin/common/models"
)

type GoodsVip struct {
	models.Model

	Layer       int `json:"layer" gorm:"type:tinyint(4);comment:排序"`
	Enable      bool `json:"enable" gorm:"type:tinyint(1);comment:开关"`
	CId         int `json:"c_id" gorm:"type:bigint(20);comment:大BID"`
	GoodsId     int `json:"goods_id" gorm:"type:bigint(20);comment:商品ID"`
	SpecsId     int `json:"specs_id" gorm:"type:bigint(20);comment:规格ID"`
	GradeId     int `json:"grade_id" gorm:"type:bigint(20);comment:VipId"`
	CustomPrice float64 `json:"customPrice" gorm:"type:float;comment:自定义价格"`
	models.ModelTime
	models.ControlBy
}

func (GoodsVip) TableName() string {
	return "goods_vip"
}

func (e *GoodsVip) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *GoodsVip) GetId() interface{} {
	return e.Id
}
