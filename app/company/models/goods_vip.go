package models

import (
	"go-admin/common/models"
	"gorm.io/gorm"
	"time"
)

type GoodsVip struct {
	models.Model

	Layer       int `json:"layer" gorm:"type:tinyint(4);comment:排序"`
	Enable      bool `json:"enable" gorm:"type:tinyint(1);comment:开关"`
	CId         int `json:"-" gorm:"type:bigint(20);comment:大BID"`
	GoodsId     int `json:"goods_id" gorm:"type:bigint(20);comment:商品ID"`
	SpecsId     int `json:"specs_id" gorm:"type:bigint(20);comment:规格ID"`
	GradeId     int `json:"grade_id" gorm:"type:bigint(20);comment:VipId"`
	CustomPrice float64 `json:"customPrice" gorm:"type:float;comment:自定义价格"`
	CreateBy int `json:"create_by" gorm:"index;comment:创建者"`
	CreatedAt time.Time      `json:"created_at" gorm:"comment:创建时间"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间"`
}

func (GoodsVip) TableName() string {
	return "goods_vip"
}

func (e *GoodsVip) GetId() interface{} {
	return e.Id
}
