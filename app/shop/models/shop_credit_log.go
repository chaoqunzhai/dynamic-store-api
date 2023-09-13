package models

import (
	"go-admin/common/models"
	"gorm.io/gorm"
)

type ShopCreditLog struct {
	models.Model
	CId  int            `json:"-" gorm:"index;comment:大B"`
	ShopId int `json:"-" gorm:"type:bigint(20);comment:小BID"`
	Number  int `json:"number" gorm:"type:double;comment:变动值"`
	Action string `json:"action" gorm:"type:varchar(10);comment:操作"`
	Type int `json:"type" gorm:"comment:操作类型"`
	Scene  string `json:"scene" gorm:"type:varchar(30);comment:变动场景"`
	Desc   string `json:"desc" gorm:"type:varchar(50);comment:描述/说明"`
	CreatedAt models.XTime      `json:"created_at" gorm:"comment:创建时间"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间"`
	CreateBy int `json:"create_by" gorm:"index;comment:创建者"`
}

func (ShopCreditLog) TableName() string {
	return "shop_credit_log"
}

func (e *ShopCreditLog) GetId() interface{} {
	return e.Id
}