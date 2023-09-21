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
    UserId     int       `json:"user_id" gorm:"type:bigint(20);comment:用户ID"`
    Salesman   int       `json:"salesman" gorm:"type:bigint(20);comment:推荐人"`
    Name       string    `json:"name" gorm:"type:varchar(30);comment:小B名称"`
    Phone      string    `json:"phone" gorm:"type:varchar(11);comment:联系手机号"`
    UserName   string    `json:"username" gorm:"type:varchar(20);comment:小B负责人名称"`
    Address    string    `json:"address" gorm:"type:varchar(60);comment:小B收货地址"`
    Longitude  float64       `json:"longitude" gorm:"type:double;comment:Longitude"`
    Latitude   float64       `json:"latitude" gorm:"type:double;comment:Latitude"`
    Image      string    `json:"image" gorm:"type:varchar(80);comment:图片"`
    LineId     int       `json:"line_id" gorm:"type:bigint(20);comment:归属配送路线"`
    GradeId   int    `json:"grade_id" gorm:"index;comment:会员等级"`
    Platform  string `json:"platform" gorm:"size:10;comment:注册来源"`
    SuggestId int    `json:"suggest_id" gorm:"index;comment:推荐人ID"`
    Balance    float64   `json:"balance" gorm:"comment:金额"`
    Integral   int       `json:"integral" gorm:"type:bigint(20);comment:可用积分"`
    Tag        []ShopTag `json:"-" gorm:"many2many:shop_mark_tag;foreignKey:id;joinForeignKey:shop_id;references:id;joinReferences:tag_id;"`
    Credit float64   `gorm:"comment:授信分" json:"credit"`
    CreateUser string    `json:"create_user" gorm:"-"`
    SalesmanUser string `json:"salesman_user" gorm:"-"`
    SalesmanPhone string `json:"salesman_phone" gorm:"-"`
    LineName       string    `json:"line_name" gorm:"-"`
    GradeName string `json:"grade_name" gorm:"-"`
    Tags []int `json:"tags" gorm:"-"`
    TagName []string `json:"tag_name" gorm:"-"`
    OrderCount int64 `json:"order_count" gorm:"-"`
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