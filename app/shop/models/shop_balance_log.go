package models

import (
	"go-admin/common/models"
	"gorm.io/gorm"
	"time"
)

type ShopBalanceLog struct {
	models.Model
	CId  int            `gorm:"index;comment:大B"`
	ShopId int `json:"shopId" gorm:"type:bigint(20);comment:小BID"`
	Money  float64 `json:"money" gorm:"type:double;comment:变动金额"`
	Action string `json:"action" gorm:"type:varchar(10);comment:操作"`
	Scene  string `json:"scene" gorm:"type:varchar(30);comment:变动场景"`
	Desc   string `json:"desc" gorm:"type:varchar(50);comment:描述/说明"`
	CreatedAt time.Time      `json:"created_at" gorm:"comment:创建时间"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间"`
	CreateBy int `json:"create_by" gorm:"index;comment:创建者"`
}

func (ShopBalanceLog) TableName() string {
	return "shop_balance_log"
}


func (e *ShopBalanceLog) GetId() interface{} {
	return e.Id
}