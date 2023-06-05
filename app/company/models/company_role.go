package models

import (
	sys "go-admin/app/admin/models"
	"go-admin/common/models"
)

type CompanyRole struct {
	models.Model

	CId     int           `json:"-" gorm:"index;comment:大BID"`
	Id      int           `json:"id" gorm:"primaryKey;autoIncrement"` // 角色编码
	Name    string        `json:"name" gorm:"size:30;"`               // 角色名称
	Enable  bool          `json:"enable"`
	Layer   int           `json:"layer"`                //角色排序
	Desc    string        `json:"desc" gorm:"size:50;"` //备注
	Admin   bool          `json:"admin" gorm:"size:4;"`
	SysMenu []DyNamicMenu `json:"sys_menu" gorm:"many2many:company_role_menu;foreignKey:id;joinForeignKey:role_id;references:id;joinReferences:menu_id;"`
	SysUser []sys.SysUser `json:"sys_user" gorm:"many2many:company_role_user;foreignKey:id;joinForeignKey:role_id;references:user_id;joinReferences:user_id;"`
	models.ModelTime
	models.ControlBy
}

func (CompanyRole) TableName() string {
	return "company_role"
}

func (e *CompanyRole) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *CompanyRole) GetId() interface{} {
	return e.Id
}
