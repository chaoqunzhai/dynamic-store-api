package models

import (
	"time"

	"go-admin/common/models"
)

type CycleTimeConf struct {
	models.Model

	Layer         string    `json:"layer" gorm:"type:tinyint;comment:排序"`
	Enable        string    `json:"enable" gorm:"type:tinyint(1);comment:开关"`
	Desc          string    `json:"desc" gorm:"type:varchar(35);comment:描述信息"`
	CId           string    `json:"cId" gorm:"type:bigint;comment:大BID"`
	Type          string    `json:"type" gorm:"type:varchar(191);comment:类型,每天,每周"`
	StartWeek     string    `json:"startWeek" gorm:"type:bigint;comment:类型为周,每周开始天"`
	EndWeek       string    `json:"endWeek" gorm:"type:bigint;comment:类型为周,每周结束天"`
	StartTime     time.Time `json:"startTime" gorm:"type:datetime(3);comment:开始下单时间"`
	EndTime       time.Time `json:"endTime" gorm:"type:datetime(3);comment:结束时间"`
	GiveDay       string    `json:"giveDay" gorm:"type:bigint;comment:跨天值为0是当天,大于0就是当天+天数"`
	GiveStartTime time.Time `json:"giveStartTime" gorm:"type:datetime(3);comment:配送开始时间"`
	GiveEndTime   time.Time `json:"giveEndTime" gorm:"type:datetime(3);comment:配送结束时间"`
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
