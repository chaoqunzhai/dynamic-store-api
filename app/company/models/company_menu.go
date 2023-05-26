package models

import (
	"go-admin/common/models"
)

type DyNamicMenu struct {
	models.Model

	Name      string `gorm:"size:30;comment:英文名称"`
	Path      string `gorm:"size:30;comment:路径,也是权限名称"`
	ParentId  int    `json:"parentId" gorm:"index;size:11;comment:父ID"`
	MetaTitle string `gorm:"size:30;comment:标题"`
	MetaIcon  string `gorm:"size:30;comment:图片"`
	Hidden    bool   `gorm:"comment:是否隐藏"`
	KeepAlive bool   `gorm:"comment:是否缓存"`
	Component string `gorm:"size:50;comment:import路径"`
	models.ModelTime
}

func (DyNamicMenu) TableName() string {
	return "dynamic_menu"
}

func (e *DyNamicMenu) GetId() interface{} {
	return e.Id
}
