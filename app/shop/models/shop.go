package models

import (
    "go-admin/common/models"
    "gorm.io/gorm"
)

type Shop struct {
    models.Model
    Layer      int       `json:"layer" gorm:"type:tinyint(4);comment:排序"`
    Enable     bool      `json:"enable" gorm:"type:tinyint(1);comment:开关"`
    Desc       string    `json:"desc" gorm:"type:varchar(25);comment:描述信息"`
    CId        int       `json:"-" gorm:"type:bigint(20);comment:大BID"`
    UserId     int       `json:"user_id" gorm:"type:bigint(20);comment:用户ID"` //是在sys_shop_user的ID
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
    CreditQuota float64 `gorm:"comment:授信额度" json:"credit_quota"`
    Credit float64   `gorm:"comment:授信余额" json:"credit"`
    CreateUser string    `json:"create_user" gorm:"-"`
    SalesmanUser string `json:"salesman_user" gorm:"-"`
    SalesmanPhone string `json:"salesman_phone" gorm:"-"`
    LineName       string    `json:"line_name" gorm:"-"`
    GradeName string `json:"grade_name" gorm:"-"`
    Tags []int `json:"tags" gorm:"-"`
    TagName []string `json:"tag_name" gorm:"-"`
    DefaultAddress string `json:"default_address" gorm:"-"`
    OrderCount int64 `json:"order_count" gorm:"-"`
    LoginTime models.XTime     `json:"login_time" gorm:"type:datetime(3);comment:登录时间"`

    IsBalanceDeduct bool `json:"is_balance_deduct" gorm:"size:1;comment:是否开启余额支付"`
    IsAli bool `json:"is_pay_ali" gorm:"size:1;comment:是否开启阿里支付"`
    IsCashOn bool `json:"is_cash_on"  gorm:"size:1;comment:是否开启货到付款"`
    IsWeChat bool `json:"is_we_chat" gorm:"size:1;comment:是否开启微信支付"`
    IsCredit bool `json:"is_credit" gorm:"size:1;comment:支持授信支付"`

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


type CompanyRegisterUserVerify struct {
    models.Model
    CreatedAt models.XTime `json:"created_at" gorm:"comment:创建时间"`
    DeletedAt gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间"`
    CId int `json:"-" gorm:"index;comment:大BID"`
    AdoptTime models.XTime  `json:"adopt_time" gorm:"通过时间"`
    AdoptUser string `json:"adopt_user" gorm:"size:11;comment:审批人"`
    Source string `json:"source" gorm:"size:6;comment:注册方式 user | mobile"`
    Value string `json:"value" gorm:"size:15;comment:注册数据,用户名"`
    Phone string `json:"phone" gorm:"size:15;comment:注册数据,手机号"`
    Password string `json:"password" gorm:"size:20;comment:密码"`
    AppTypeName string `json:"app_type_name" gorm:"size:6;comment:注册来源例如H5,WECHAT,ALI等"`
    Status int `json:"status" gorm:"default:0;index;comment:0:审核中, 1:通过 -1:驳回 2:已经创建门店"`
    Info string `json:"info" gorm:"size:10;comment:备注"`
    ShopUserId int `json:"shop_user_id" gorm:"comment:创建成功的大B用户ID" `
}
func (CompanyRegisterUserVerify) TableName() string {
    return "company_register_user_verify"
}