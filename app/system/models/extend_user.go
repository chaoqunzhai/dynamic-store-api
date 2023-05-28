package models

import (
	"go-admin/common/models"
)

type ExtendUser struct {
    models.Model
    UserId int `gorm:"index;comment:用户ID"`
	CId int `gorm:"index;comment:归属大B"`
    Layer int `json:"layer" gorm:"type:tinyint(4);comment:排序"`
    Enable bool `json:"enable" gorm:"type:tinyint(1);comment:开关"`
    Platform string `json:"platform" gorm:"type:varchar(12);comment:注册来源"` 
    GradeId int `json:"gradeId" gorm:"type:bigint(20);comment:会员等级"`
    SuggestId int `json:"suggestId" gorm:"type:bigint(20);comment:推荐人ID"`
    InvitationCode string `json:"invitationCode" gorm:"type:varchar(10);comment:本人邀请码"`
    Desc string
	models.ModelTime
	models.ControlBy
}

func (ExtendUser) TableName() string {
    return "extend_user"
}

func (e *ExtendUser) Generate() models.ActiveRecord {
	o := *e
	return &o
}
func (e *ExtendUser) GetId() interface{} {
	return e.Id
}