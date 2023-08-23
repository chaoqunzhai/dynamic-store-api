package models

import (
	"database/sql"
	"go-admin/common/models"
)

type CompanyCoupon struct {
	models.Model

	Layer  int    `json:"layer" gorm:"type:tinyint(4);comment:排序"`
	Enable bool   `json:"enable" gorm:"type:tinyint(1);comment:开关"`
	Desc   string `json:"desc" gorm:"type:varchar(35);comment:描述信息"`
	CId    int    `json:"-" gorm:"type:bigint(20);comment:大BID"`

	Name       string       `json:"name" gorm:"size:50;comment:优惠卷名称"`
	Type       int          `json:"type" gorm:"type:tinyint(1);default:0;comment:类型,0:满减,1:折扣"`
	Range      int          `json:"range" gorm:"type:tinyint(1);default:2;comment:使用范围,0:指定商品,1:指定分类 2:全场"`
	Reduce     float64      `json:"reduce" gorm:"comment:优惠卷金额"`
	Discount   float64      `json:"discount" gorm:"comment:折扣率"`
	Threshold  float64      `json:"threshold" gorm:"comment:满多少钱可以用"`
	ExpireType int          `json:"expire_type" gorm:"type:tinyint(1);default:0;comment:到期类型,0:领取后生效，1:指定日期生效"`
	ExpireDay  int          `json:"expire_day" gorm:"type:tinyint(1);default:1;comment:过期多少天"`
	First      bool         `json:"first"`
	Automatic  bool         `json:"automatic"`
	StartTime  sql.NullTime `json:"start_time" gorm:"comment:开始使用时间"`
	EndTime    sql.NullTime `json:"end_time" gorm:"comment:截止使用时间"`
	Inventory  int          `json:"inventory" gorm:"comment:库存"`
	Limit      int          `json:"limit" gorm:"comment:每个人限领次数"`
	Start      string       `json:"start" gorm:"-"`
	End        string       `json:"end" gorm:"-"`
	models.ModelTime
	models.ControlBy
}

func (CompanyCoupon) TableName() string {
	return "company_coupon"
}

func (e *CompanyCoupon) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *CompanyCoupon) GetId() interface{} {
	return e.Id
}
