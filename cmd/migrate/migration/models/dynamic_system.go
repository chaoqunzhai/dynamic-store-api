package models

// todo: 用户扩展信息
type ExtendUser struct {
	MiniGlobal
	Platform       string `json:"platform" gorm:"size:12;comment:注册来源"`
	GradeId        int    `gorm:"index;comment:会员等级"`
	SuggestId      int    `gorm:"index;comment:推荐人ID"`
	InvitationCode string `gorm:"size:10;comment:本人邀请码"`
}

func (ExtendUser) TableName() string {
	return "extend_user"
}

// todo: 每个大B设置的角色
// 这里为什么没有使用系统的角色,
// 因为:系统和大B的角色是需要区分隔离开
type CompanyRole struct {
	Id      int    `json:"id" gorm:"primaryKey;autoIncrement"` // 角色编码
	Name    string `json:"roleName" gorm:"size:30;"`           // 角色名称
	Enable  int
	Sort    int       //角色排序
	Remark  string    `json:"remark" gorm:"size:50;"` //备注
	Admin   bool      `json:"admin" gorm:"size:4;"`
	SysMenu []SysMenu `json:"sysMenu" gorm:"many2many:company_role_menu;foreignKey:id;joinForeignKey:role_id;references:MenuId;joinReferences:menu_id;"`
	ControlBy
	ModelTime
}

func (CompanyRole) TableName() string {
	return "company_role"
}
