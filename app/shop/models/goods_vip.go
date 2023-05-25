package models

import (
	"go-admin/common/models"
)

type GoodsVip struct {
	models.Model

	Layer       string `json:"layer" gorm:"type:tinyint(4);comment:排序"`
	Enable      string `json:"enable" gorm:"type:tinyint(1);comment:开关"`
	CId         string `json:"cId" gorm:"type:bigint(20);comment:大BID"`
	GoodsId     string `json:"goodsId" gorm:"type:bigint(20);comment:商品ID"`
	SpecsId     string `json:"specsId" gorm:"type:bigint(20);comment:规格ID"`
	GradeId     string `json:"gradeId" gorm:"type:bigint(20);comment:VipId"`
	CustomPrice string `json:"customPrice" gorm:"type:float;comment:自定义价格"`
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
