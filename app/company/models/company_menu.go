package models

import (
	"go-admin/common/models"
)

type DyNamicMenu struct {
	models.Model

	Name      string `json:"name" gorm:"size:30;comment:英文名称"`
	Path      string `json:"path" gorm:"size:30;comment:路径,也是权限名称"`
	ParentId  int    `json:"parent_id" gorm:"index;size:11;comment:父ID"`
	MetaTitle string `json:"meta_title" gorm:"size:30;comment:标题"`
	MetaIcon  string `json:"meta_icon" gorm:"size:30;comment:图片"`
	Hidden    bool   `json:"hidden" gorm:"comment:是否隐藏"`
	KeepAlive bool   `json:"keep_alive" gorm:"comment:是否缓存"`
	Component string `json:"component" gorm:"size:50;comment:import路径"`
	models.ModelTime
}

func (DyNamicMenu) TableName() string {
	return "dynamic_menu"
}

func (e *DyNamicMenu) GetId() interface{} {
	return e.Id
}

type DyMbmMenu struct {
	models.Model
	Layer     int    `gorm:"index;comment:排序"` //排序
	Enable    bool   `gorm:"comment:开关"`
	Name      string `gorm:"size:30;comment:英文名称"`
	MetaTitle string `gorm:"size:30;comment:标题"`
	ParentId  int    `json:"parentId" gorm:"index;size:11;comment:父ID"`
}

func (DyMbmMenu) TableName() string {
	return "dynamic_mbm_menu"
}


