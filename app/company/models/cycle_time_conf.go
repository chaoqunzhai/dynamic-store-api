package models

import (
	"go-admin/common/models"
)

type CycleTimeConf struct {
	models.Model
	Layer     int    `json:"layer" gorm:"type:tinyint;comment:排序"`
	Enable    bool   `json:"enable" gorm:"type:tinyint(1);comment:开关"`
	Show      bool   `json:"show" gorm:"type:tinyint(1);comment:是否客户端展示"`
	Desc      string `json:"desc" gorm:"type:varchar(35);comment:描述信息"`
	CId       int    `json:"-" gorm:"type:bigint;comment:大BID"`
	Type      int    `json:"type" gorm:";comment:类型,每天,每周"`
	Uid       string `gorm:"type:varchar(4);comment:周期名称都是天,防止一天可能多个不同周期的配置,加个标识区分周期"`
	StartWeek int    `json:"start_week" gorm:"type:bigint;comment:类型为周,每周开始天"`
	EndWeek   int    `json:"end_week" gorm:"type:bigint;comment:类型为周,每周结束天"`
	StartTime string `json:"start_time" gorm:"type:varchar(5);comment:开始下单时间"`
	EndTime   string `json:"end_time" gorm:"type:varchar(5);comment:结束时间"`
	GiveDay   int    `json:"give_day" gorm:"type:bigint;comment:跨天值为0是当天,大于0就是当天+天数"`
	GiveTime  string `json:"give_time" gorm:"type:varchar(30);comment:配送时间,例如：15点至19点"`
	models.ModelTime
	models.ControlBy
}

func (CycleTimeConf) TableName() string {
	return "cycle_time_conf"
}

func (e *CycleTimeConf) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *CycleTimeConf) GetId() interface{} {
	return e.Id
}
