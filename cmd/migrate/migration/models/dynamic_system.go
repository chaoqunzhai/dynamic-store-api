package models

import "time"

type DyNamicMenu struct {
	Model
	ModelTime
	Layer     int    `gorm:"size:1;index;comment:排序"` //排序
	Enable    bool   `gorm:"comment:开关"`
	Role      string `gorm:"size:20;comment:哪个角色菜单"`
	Name      string `gorm:"size:30;comment:英文名称"`
	Path      string `gorm:"size:30;comment:路径,也是权限名称"`
	ParentId  int    `json:"parentId" gorm:"index;size:11;comment:父ID"`
	MetaTitle string `gorm:"size:30;comment:标题"`
	MetaIcon  string `gorm:"size:30;comment:图片"`
	Hidden    bool   `gorm:"comment:是否隐藏"`
	KeepAlive bool   `gorm:"comment:是否缓存"`
	Component string `gorm:"size:50;comment:import路径"`
}

func (DyNamicMenu) TableName() string {
	return "dynamic_menu"
}

// todo: 用户扩展信息
type ExtendUser struct {
	BigBRichGlobal
	UserId int `gorm:"index;comment:用户ID"`
}

func (ExtendUser) TableName() string {
	return "extend_user"
}

// todo: 每个大B设置的角色
// 这里为什么没有使用系统的角色,
// 因为:系统和大B的角色是需要区分隔离开,不能混淆
type CompanyRole struct {
	Model
	CId     int    `gorm:"index;comment:大BID"`
	Id      int    `json:"id" gorm:"primaryKey;autoIncrement"` // 角色编码
	Name    string `json:"roleName" gorm:"size:30;"`           // 角色名称
	Enable  bool
	Layer   int           `gorm:"size:1;index;comment:排序"`        //排序
	Desc    string        `json:"desc" gorm:"size:35;comment:备注"` //备注
	Admin   bool          `json:"admin" gorm:"size:4;"`
	SysMenu []DyNamicMenu `json:"sysMenu" gorm:"many2many:company_role_menu;foreignKey:id;joinForeignKey:role_id;references:id;joinReferences:menu_id;"`
	SysUser []SysUser     `json:"sysUser" gorm:"many2many:company_role_user;foreignKey:id;joinForeignKey:role_id;references:user_id;joinReferences:user_id;"`
	ControlBy
	ModelTime
}

func (CompanyRole) TableName() string {
	return "company_role"
}

// todo:优惠卷
type Coupon struct {
	BigBRichGlobal
	Name       string    `json:"name" gorm:"size:50;comment:优惠卷名称"`
	Type       int       `json:"type" gorm:"type:tinyint(1);default:0;comment:类型,0:满减,1:折扣"`
	Range      int       `gorm:"type:tinyint(1);default:2;comment:使用范围,0:指定商品,1:全场通用 2:全场"`
	Reduce     float64   `gorm:"comment:优惠卷金额"`
	Discount   float64   `gorm:"comment:折扣率"`
	Threshold  float64   `gorm:"comment:满多少钱可以用"`
	First      bool      `gorm:"comment:是否首推,首推的时候下单才会自动领取"`
	Automatic  bool      `gorm:"comment:下单自动领取"`
	ExpireType int       `gorm:"type:tinyint(1);default:0;comment:到期类型,0:领取后生效，1:指定日期生效"`
	ExpireDay  int       `gorm:"type:tinyint(1);default:1;comment:过期多少天"`
	StartTime  time.Time `gorm:"comment:开始使用时间"`
	EndTime    time.Time `gorm:"comment:截止使用时间"`
	ReceiveNum int `json:"receive_num" gorm:"comment:已经领取个数"`
	Inventory  int       `gorm:"comment:库存"`
	Limit      int       `gorm:"comment:每个人限领次数"`
}

func (Coupon) TableName() string {
	return "company_coupon"
}

// todo:优惠卷领取记录
type ReceiveCouponLog struct {
	BigBRichUserGlobal
	CouponId   int     `json:"coupon_id" gorm:"index;comment:优惠卷ID"`
	CouponType int     `json:"coupon_type" gorm:"type:tinyint(1);default:0;comment:优惠卷类型,0:满减,1:折扣"`
	Type       int     `json:"type" gorm:"type:tinyint(1);default:2;comment:类型,1:订单领取 2:自己领取 3:活动领取"`
	Status int `json:"status" gorm:"index;default:1;comment:优惠卷状态,1:未使用 2:已使用 3:已过期"`
}

func (ReceiveCouponLog) TableName() string {
	return "user_receive_coupon_log"
}

//todo:用户余额
type UserAmountStore struct {
	BigBRichUserGlobal
	Balance float64 `json:"balance"  gorm:"comment:余额"`
	Credit float64 `json:"credit" gorm:"comment:授信分"`
	Point int `json:"point" gorm:"comment:积分"`
	GradeId int `json:"grade_id" gorm:"comment:VIP等级"`
}
func (UserAmountStore) TableName() string {
	return "user_amount_store"
}