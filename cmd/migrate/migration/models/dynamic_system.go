package models

type DyNamicMenu struct {
	Model
	ModelTime
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
// 因为:系统和大B的角色是需要区分隔离开,不能混淆
type CompanyRole struct {
	Model
	CId     int    `gorm:"index;comment:大BID"`
	Id      int    `json:"id" gorm:"primaryKey;autoIncrement"` // 角色编码
	Name    string `json:"roleName" gorm:"size:30;"`           // 角色名称
	Enable  int
	Sort    int           //角色排序
	Remark  string        `json:"remark" gorm:"size:50;"` //备注
	Admin   bool          `json:"admin" gorm:"size:4;"`
	SysMenu []DyNamicMenu `json:"sysMenu" gorm:"many2many:company_role_menu;foreignKey:id;joinForeignKey:role_id;references:id;joinReferences:menu_id;"`
	SysUser []SysUser `json:"sysUser" gorm:"many2many:company_role_user;foreignKey:id;joinForeignKey:role_id;references:user_id;joinReferences:user_id;"`
	ControlBy
	ModelTime
}

func (CompanyRole) TableName() string {
	return "company_role"
}
