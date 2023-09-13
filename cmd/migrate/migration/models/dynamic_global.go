package models

import (
	"gorm.io/gorm"
	"time"
)

// todo: 比较简单的函数,节约字段
type MiniGlobal struct {
	Model
	CreateBy  int            `json:"createBy" gorm:"index;comment:创建者"`
	CreatedAt time.Time      `json:"createdAt" gorm:"comment:创建时间"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间"`
	Layer     int            `gorm:"size:1;index;comment:排序"` //排序,默认是0,从0开始倒序排列
	Enable    bool           `gorm:"comment:开关"`
}

// todo: 公共字段更丰富的
type RichGlobal struct {
	Model
	ControlBy
	ModelTime
	Layer  int    `json:"layer" gorm:"size:1;default:1;index;comment:排序"` //排序
	Enable bool   `json:"enable" gorm:"default:true;comment:开关"`
	Desc   string `json:"desc" gorm:"size:35;comment:描述信息"` //描述
}

// todo: 包含大BID的公共函数
type BigBRichGlobal struct {
	RichGlobal
	CId int `json:"c_id" gorm:"index;comment:大BID"`
}
type BigBRichUserGlobal struct {
	RichGlobal
	CId    int `gorm:"index;comment:大BID"`
	UserId int `gorm:"index;comment:管理员ID"`
}

// todo: 大B下小B的简约公共函数
type BigBMiniGlobal struct {
	MiniGlobal
	CId int `gorm:"index;comment:大BID"`
}

// todo:公司映射表,起到分表作用
type SplitTableMap struct {
	RichGlobal
	CId  int    `gorm:"index;comment:公司ID"`                             //公司ID
	Type int    `gorm:"type:tinyint(1);default:1;index;comment:映射表的类型"` //根据不同的类型,来做细分
	Name string `gorm:"size:60;index;comment:对应表的名称"`                   //分表的名称
	Desc string `gorm:"size:30;index;comment:对应表的名称"`
}

func (SplitTableMap) TableName() string {
	return "split_table_map"
}
