package models

import (
	"time"
)

type SysLoginLog struct {
	Model
	UserId int `json:"user_id" gorm:"index;comment:登录用户ID"`
	CId int `json:"c_id" gorm:"index;comment:大BID"`
	User      string    `json:"user" gorm:"type:varchar(18);comment:登录的用户名称/手机号"`
	UserType string `json:"user_type" gorm:"type:varchar(10);comment:登录的类型"`
	Client string `json:"client" gorm:"type:varchar(15);comment:端"`
	Ipaddr        string    `json:"ipaddr" gorm:"type:varchar(15);comment:ip地址"`
	Source        string    `json:"source" gorm:"type:varchar(12);comment:登录来源"`
	LoginTime     time.Time `json:"login_time" gorm:"type:timestamp;comment:登录时间"`
	Remark        string    `json:"remark" gorm:"type:varchar(6);comment:备注"`
	Role int `json:"role" gorm:"type:varchar(1);comment:0:超管 1:大B 2:小B"`
}

func (SysLoginLog) TableName() string {
	return "sys_login_log"
}
