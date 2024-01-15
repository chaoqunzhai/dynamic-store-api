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
//todo:更简单的字段，一般用来记录日志
type MiniLog struct {
	Model
	CId            int       `json:"c_id" gorm:"index;comment:公司(大B)ID"`
	CreateBy  int            `json:"createBy" gorm:"index;comment:创建者"`
	CreatedAt time.Time      `json:"createdAt" gorm:"comment:创建时间"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间"`
	Desc   string `json:"desc" gorm:"size:30;comment:描述信息"` //描述
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
//没有创建者和更新者
type BigBNoCreateByRichGlobal struct {
	Model
	ModelTime
	CId int `gorm:"index;comment:大BID"`
	UserId int `gorm:"index;comment:用户ID"`
	Enable bool   `gorm:"default:true;comment:开关"`
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
	Model
	Enable bool   `gorm:"default:true;comment:开关"`
	CreateBy  int            `json:"createBy" gorm:"index;comment:创建者"`
	CreatedAt time.Time      `json:"createdAt" gorm:"comment:创建时间"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间"`
	CId  int    `gorm:"index;comment:公司ID"`                             //公司ID
	OrderTable string `gorm:"size:30;index;comment:订单表"`                   //分表的名称
	OrderSpecs string `gorm:"size:30;index;comment:订单规格表"`
	OrderCycle string `gorm:"size:30;index;comment:周期配送下单索引表"`
	OrderEdit string  `gorm:"size:30;index;comment:订单修改表"`
	OrderReturn string `gorm:"size:30;index;comment:订单退货表"`
	InventoryRecordLog string `gorm:"size:30;index;comment:出入库记录流水表"`
	CycleCnfList []string `json:"cycle_cnf_list" gorm:"-"` //只是定时任务需要用这个字段 做一个更新,gorm不会更新到DB中
}

func (SplitTableMap) TableName() string {
	return "split_table_map"
}