/*
*
@Author: chaoqun
* @Date: 2023/7/20 22:46
*/
package models

import (
	"go-admin/common/models"
	"gorm.io/gorm"
	"time"
)

type CompanyRegisterCnf struct {
	models.Model
	Layer  int    `json:"layer" gorm:"type:tinyint(4);comment:排序"`
	Enable bool   `json:"enable" gorm:"type:tinyint(1);comment:开关"`
	Desc   string `json:"desc" gorm:"type:varchar(35);comment:描述信息"`
	CId    int    `json:"-" gorm:"type:bigint(20);comment:大BID"`
	Type   int    `json:"type" gorm:"type:tinyint(1);comment:类型"` //0:登录  1:注册
	Value  string `json:"login" gorm:"size:12;comment:登录方式"`      //username,mobile,wechat 代表用户名,手机号,微信

	models.ModelTime
	models.ControlBy
}

func (CompanyRegisterCnf) TableName() string {
	return "company_register_cnf"
}

type CompanyNavCnf struct {
	models.Model
	Enable    bool
	GId       int            `gorm:"index;type:tinyint(1);comment:关联的菜单配置ID"`
	CId       int            `gorm:"index;comment:大B"`
	CreatedAt time.Time      `json:"created_at" gorm:"comment:创建时间"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间"`
}

func (CompanyNavCnf) TableName() string {
	return "company_nav_cnf"
}


type CompanyQuickTools struct {
	models.Model
	Enable    bool
	QuickId int `gorm:"index;comment:关联的导航配置"`
	CId       int            `gorm:"index;comment:大B"`
	CreatedAt time.Time      `json:"created_at" gorm:"comment:创建时间"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间"`
}

func (CompanyQuickTools) TableName() string {
	return "company_quick_tools"
}