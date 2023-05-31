package models

import (
    "go-admin/common/models"
)

type Shop struct {
    models.Model
    Layer      int       `json:"layer" gorm:"type:tinyint(4);comment:排序"`
    Enable     bool      `json:"enable" gorm:"type:tinyint(1);comment:开关"`
    Desc       string    `json:"desc" gorm:"type:varchar(25);comment:描述信息"`
    CId        int       `json:"-" gorm:"type:bigint(20);comment:大BID"`
    UserId     int       `json:"user_id" gorm:"type:bigint(20);comment:管理员ID"`
    Salesman   int       `json:"salesman" gorm:"type:bigint(20);comment:推荐人"`
    Name       string    `json:"name" gorm:"type:varchar(30);comment:小B名称"`
    Phone      string    `json:"phone" gorm:"type:varchar(11);comment:联系手机号"`
    UserName   string    `json:"userName" gorm:"type:varchar(20);comment:小B负责人名称"`
    Address    string    `json:"address" gorm:"type:varchar(200);comment:小B收货地址"`
    Longitude  int       `json:"longitude" gorm:"type:double;comment:Longitude"`
    Latitude   int       `json:"latitude" gorm:"type:double;comment:Latitude"`
    Image      string    `json:"image" gorm:"type:varchar(80);comment:图片"`
    LineId     int       `json:"-" gorm:"type:bigint(20);comment:归属配送路线"`
    CreateUser string    `json:"create_user" gorm:"-"`
    Line       string    `json:"line" gorm:"-"`
    Amount     float64   `json:"amount" gorm:"type:double;comment:剩余金额"`
    Integral   int       `json:"integral" gorm:"type:bigint(20);comment:可用积分"`
    Tag        []ShopTag `json:"tag" gorm:"many2many:shop_mark_tag;foreignKey:id;joinForeignKey:shop_id;references:id;joinReferences:tag_id;"`
    models.ModelTime
    models.ControlBy
}

func (Shop) TableName() string {
    return "shop"
}

func (e *Shop) Generate() models.ActiveRecord {
    o := *e
    return &o
}

func (e *Shop) GetId() interface{} {
    return e.Id
}