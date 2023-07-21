package models

import (
	"time"

	"go-admin/common/models"
)

type Company struct {
	models.Model

	Layer          int           `json:"layer" gorm:"type:tinyint(4);comment:排序"`
	Enable         bool          `json:"enable" gorm:"type:tinyint(1);comment:开关"`
	Desc           string        `json:"desc" gorm:"type:varchar(25);comment:描述信息"`
	Name           string        `json:"name" gorm:"type:varchar(30);comment:公司(大B)名称"`
	Phone          string        `json:"phone" gorm:"type:varchar(11);comment:负责人联系手机号"`
	UserName       string        `json:"user_name" gorm:"type:varchar(20);comment:大B负责人名称"`
	ShopName       string        `json:"shop_name" gorm:"type:varchar(50);comment:自定义大B系统名称"`
	Address        string        `json:"address" gorm:"type:varchar(155);comment:大B地址位置"`
	Longitude      float64       `json:"longitude" gorm:"type:double;comment:Longitude"`
	Latitude       float64       `json:"latitude" gorm:"type:double;comment:Latitude"`
	Image          string        `json:"image" gorm:"type:varchar(80);comment:logo图片"`
	RenewalTime    time.Time     `json:"renewal_time" gorm:"type:datetime(3);comment:续费时间"`
	ExpirationTime time.Time     `json:"expiration_time" gorm:"type:datetime(3);comment:到期时间"`
	NavList        []interface{} `json:"nav_list" gorm:"-"` //展示大B的菜单
	models.ModelTime
	models.ControlBy
}

func (Company) TableName() string {
	return "company"
}

func (e *Company) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *Company) GetId() interface{} {
	return e.Id
}
