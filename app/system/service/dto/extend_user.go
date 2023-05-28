package dto

import (
     


	"go-admin/app/system/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type ExtendUserGetPageReq struct {
	dto.Pagination     `search:"-"`
    Layer string `form:"layer"  search:"type:exact;column:layer;table:extend_user" comment:"排序"`
    Enable string `form:"enable"  search:"type:exact;column:enable;table:extend_user" comment:"开关"`
    Platform string `form:"platform"  search:"type:exact;column:platform;table:extend_user" comment:"注册来源"`
    GradeId string `form:"gradeId"  search:"type:exact;column:grade_id;table:extend_user" comment:"会员等级"`
    SuggestId string `form:"suggestId"  search:"type:exact;column:suggest_id;table:extend_user" comment:"推荐人ID"`
    InvitationCode string `form:"invitationCode"  search:"type:contains;column:invitation_code;table:extend_user" comment:"本人邀请码"`
    ExtendUserOrder
}

type ExtendUserOrder struct {
    Id string `form:"idOrder"  search:"type:order;column:id;table:extend_user"`
    CreateBy string `form:"createByOrder"  search:"type:order;column:create_by;table:extend_user"`
    CreatedAt string `form:"createdAtOrder"  search:"type:order;column:created_at;table:extend_user"`
    DeletedAt string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:extend_user"`
    Layer string `form:"layerOrder"  search:"type:order;column:layer;table:extend_user"`
    Enable string `form:"enableOrder"  search:"type:order;column:enable;table:extend_user"`
    Platform string `form:"platformOrder"  search:"type:order;column:platform;table:extend_user"`
    GradeId string `form:"gradeIdOrder"  search:"type:order;column:grade_id;table:extend_user"`
    SuggestId string `form:"suggestIdOrder"  search:"type:order;column:suggest_id;table:extend_user"`
    InvitationCode string `form:"invitationCodeOrder"  search:"type:order;column:invitation_code;table:extend_user"`
    
}

func (m *ExtendUserGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type ExtendUserInsertReq struct {
    Id int `json:"-" comment:"主键编码"` // 主键编码
    UserId int `json:"user_id" comment:""`
    Layer int `json:"layer" comment:"排序"`
    Enable bool `json:"enable" comment:"开关"`
    Platform string `json:"platform" comment:"注册来源"`
    GradeId int `json:"grade_id" comment:"会员等级"`
    SuggestId int `json:"suggest_id" comment:"推荐人ID"`
    InvitationCode string `json:"invitationCode" comment:"本人邀请码"`
    Desc string `json:"desc" comment:""`
    common.ControlBy
}

func (s *ExtendUserInsertReq) Generate(model *models.ExtendUser)  {
    if s.Id == 0 {
        model.Model = common.Model{ Id: s.Id }
    }
    model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
    model.UserId = s.UserId
    model.Layer = s.Layer
    model.Enable = s.Enable
    model.Platform = s.Platform
    model.GradeId = s.GradeId
    model.SuggestId = s.SuggestId
    model.InvitationCode = s.InvitationCode
    model.Desc = s.Desc
}

func (s *ExtendUserInsertReq) GetId() interface{} {
	return s.Id
}

type ExtendUserUpdateReq struct {
    Id int `uri:"id" comment:"主键编码"` // 主键编码
    Layer int `json:"layer" comment:"排序"`
    UserId int `json:"user_id" comment:""`
    Enable bool `json:"enable" comment:"开关"`
    Platform string `json:"platform" comment:"注册来源"`
    GradeId int `json:"gradeId" comment:"会员等级"`
    SuggestId int `json:"suggestId" comment:"推荐人ID"`
    InvitationCode string `json:"invitationCode" comment:"本人邀请码"`
    Desc string `json:"desc" comment:""`
    common.ControlBy
}

func (s *ExtendUserUpdateReq) Generate(model *models.ExtendUser)  {
    if s.Id == 0 {
        model.Model = common.Model{ Id: s.Id }
    }
    model.UserId = s.UserId
    model.Layer = s.Layer
    model.Enable = s.Enable
    model.Platform = s.Platform
    model.GradeId = s.GradeId
    model.SuggestId = s.SuggestId
    model.InvitationCode = s.InvitationCode
    model.Desc = s.Desc
}

func (s *ExtendUserUpdateReq) GetId() interface{} {
	return s.Id
}

// ExtendUserGetReq 功能获取请求参数
type ExtendUserGetReq struct {
     Id int `uri:"id"`
}
func (s *ExtendUserGetReq) GetId() interface{} {
	return s.Id
}


type ExtendUserGradeReq struct {
    Grade int `json:"grade"`
    Ids []int `json:"ids"`
}
// ExtendUserDeleteReq 功能删除请求参数
type ExtendUserDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *ExtendUserDeleteReq) GetId() interface{} {
	return s.Ids
}
