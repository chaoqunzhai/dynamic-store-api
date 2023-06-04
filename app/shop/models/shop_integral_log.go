package models

import (
	"go-admin/common/models"
	"gorm.io/gorm"
)

type ShopIntegralLog struct {
    models.Model
	CId  int            `gorm:"index;comment:大B"`
    ShopId int `json:"shopId" gorm:"type:bigint(20);comment:小BID"`
    Number int `json:"number" gorm:"type:double;comment:积分变动数值"`
    Scene string `json:"scene" gorm:"type:varchar(30);comment:变动场景"`
	Action string `json:"action" gorm:"type:varchar(10);comment:操作"`
    Desc string `json:"desc" gorm:"type:varchar(50);comment:描述/说明"`
	Type int `json:"type" gorm:"comment:操作类型"`
	CreatedAt models.XTime      `json:"created_at" gorm:"comment:创建时间"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间"`
	CreateBy int `json:"create_by" gorm:"index;comment:创建者"`
}

func (ShopIntegralLog) TableName() string {
    return "shop_integral_log"
}



func (e *ShopIntegralLog) GetId() interface{} {
	return e.Id
}