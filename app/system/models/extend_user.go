package models

import (

	"go-admin/common/models"

)

type ExtendUser struct {
    models.Model
    
    Layer string `json:"layer" gorm:"type:tinyint(4);comment:排序"` 
    Enable string `json:"enable" gorm:"type:tinyint(1);comment:开关"` 
    Platform string `json:"platform" gorm:"type:varchar(12);comment:注册来源"` 
    GradeId string `json:"gradeId" gorm:"type:bigint(20);comment:会员等级"` 
    SuggestId string `json:"suggestId" gorm:"type:bigint(20);comment:推荐人ID"` 
    InvitationCode string `json:"invitationCode" gorm:"type:varchar(10);comment:本人邀请码"` 
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