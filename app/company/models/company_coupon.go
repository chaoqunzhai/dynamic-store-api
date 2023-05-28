package models

import (
     
     
     
     "time"
     
     
     
     
     "go-admin/common/models"

)

type CompanyCoupon struct {
    models.Model
    
    Layer int `json:"layer" gorm:"type:tinyint(4);comment:排序"`
    Enable bool `json:"enable" gorm:"type:tinyint(1);comment:开关"`
    Desc string `json:"desc" gorm:"type:varchar(35);comment:描述信息"` 
    CId int `json:"c_id" gorm:"type:bigint(20);comment:大BID"`
    Name string `json:"name" gorm:"type:varchar(50);comment:优惠卷名称"` 
    Type int `json:"type" gorm:"type:bigint(20);comment:类型"`
    Range int `json:"range" gorm:"type:bigint(20);comment:使用范围"`
    Money float64 `json:"money" gorm:"type:bigint(20);comment:优惠卷金额"`
    Min float64 `json:"min" gorm:"type:double;comment:最低多少钱可以用"`
    Max float64 `json:"max" gorm:"type:double;comment:满多少钱可以用"`
    StartTime time.Time `json:"start_time" gorm:"type:datetime(3);comment:开始使用时间"`
    EndTime time.Time `json:"end_time" gorm:"type:datetime(3);comment:截止使用时间"`
    Inventory int `json:"inventory" gorm:"type:bigint(20);comment:库存"`
    Limit int `json:"limit" gorm:"type:bigint(20);comment:每个人限领次数"`
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