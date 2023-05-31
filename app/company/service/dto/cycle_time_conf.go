package dto

import (
	"go-admin/app/company/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type CycleTimeConfGetPageReq struct {
	dto.Pagination `search:"-"`
	Layer          string `form:"layer"  search:"type:exact;column:layer;table:cycle_time_conf" comment:"排序"`
	Enable         string `form:"enable"  search:"type:exact;column:enable;table:cycle_time_conf" comment:"开关"`
	CId            string `form:"cId"  search:"type:exact;column:c_id;table:cycle_time_conf" comment:"大BID"`
	Type           string `form:"type"  search:"type:exact;column:type;table:cycle_time_conf" comment:"类型,每天,每周"`
	CycleTimeConfOrder
}

type CycleTimeConfOrder struct {
	Id        string `form:"idOrder"  search:"type:order;column:id;table:cycle_time_conf"`
	CreateBy  string `form:"createByOrder"  search:"type:order;column:create_by;table:cycle_time_conf"`
	UpdateBy  string `form:"updateByOrder"  search:"type:order;column:update_by;table:cycle_time_conf"`
	CreatedAt string `form:"createdAtOrder"  search:"type:order;column:created_at;table:cycle_time_conf"`
	UpdatedAt string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:cycle_time_conf"`
	DeletedAt string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:cycle_time_conf"`
	Layer     string `form:"layerOrder"  search:"type:order;column:layer;table:cycle_time_conf"`
	Enable    string `form:"enableOrder"  search:"type:order;column:enable;table:cycle_time_conf"`
	Desc      string `form:"descOrder"  search:"type:order;column:desc;table:cycle_time_conf"`
	CId       string `form:"cIdOrder"  search:"type:order;column:c_id;table:cycle_time_conf"`
	Type      string `form:"typeOrder"  search:"type:order;column:type;table:cycle_time_conf"`
	StartWeek string `form:"startWeekOrder"  search:"type:order;column:start_week;table:cycle_time_conf"`
	EndWeek   string `form:"endWeekOrder"  search:"type:order;column:end_week;table:cycle_time_conf"`
	StartTime string `form:"startTimeOrder"  search:"type:order;column:start_time;table:cycle_time_conf"`
	EndTime   string `form:"endTimeOrder"  search:"type:order;column:end_time;table:cycle_time_conf"`
	GiveDay   string `form:"giveDayOrder"  search:"type:order;column:give_day;table:cycle_time_conf"`
	GiveTime  string `form:"giveTimeOrder"  search:"type:order;column:give_time;table:cycle_time_conf"`
}

func (m *CycleTimeConfGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type CycleTimeConfInsertReq struct {
	Id     int    `json:"-" comment:"主键编码"` // 主键编码
	Layer  int    `json:"layer" comment:"排序"`
	Enable bool   `json:"enable" comment:"开关"`
	Desc   string `json:"desc" comment:"描述信息"`

	Type      int    `json:"type" comment:"类型,每天,每周"`
	StartWeek int    `json:"start_week" comment:"类型为周,每周开始天"`
	EndWeek   int    `json:"end_week" comment:"类型为周,每周结束天"`
	StartTime string `json:"start_time" comment:"开始下单时间"`
	EndTime   string `json:"end_time" comment:"结束时间"`
	GiveDay   int    `json:"give_day" comment:"跨天值为0是当天,大于0就是当天+天数"`
	GiveTime  string `json:"give_time" comment:"配送时间,例如：15点至19点"`
	common.ControlBy
}

func (s *CycleTimeConfInsertReq) Generate(model *models.CycleTimeConf) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
	model.Layer = s.Layer
	model.Enable = s.Enable
	model.Desc = s.Desc

	model.Type = s.Type
	model.StartWeek = s.StartWeek
	model.EndWeek = s.EndWeek

	model.StartTime = s.StartTime

	model.EndTime = s.EndTime

	model.GiveDay = s.GiveDay
	model.GiveTime = s.GiveTime
}

func (s *CycleTimeConfInsertReq) GetId() interface{} {
	return s.Id
}

type CycleTimeConfUpdateReq struct {
	Id     int    `uri:"id" comment:"主键编码"` // 主键编码
	Layer  int    `json:"layer" comment:"排序"`
	Enable bool   `json:"enable" comment:"开关"`
	Desc   string `json:"desc" comment:"描述信息"`

	Type      int    `json:"type" comment:"类型,每天,每周"`
	StartWeek int    `json:"start_week" comment:"类型为周,每周开始天"`
	EndWeek   int    `json:"end_week" comment:"类型为周,每周结束天"`
	StartTime string `json:"start_time" comment:"开始下单时间"`
	EndTime   string `json:"end_time" comment:"结束时间"`
	GiveDay   int    `json:"give_day" comment:"跨天值为0是当天,大于0就是当天+天数"`
	GiveTime  string `json:"give_time" comment:"配送时间,例如：15点至19点"`
	common.ControlBy
}

func (s *CycleTimeConfUpdateReq) Generate(model *models.CycleTimeConf) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.UpdateBy = s.UpdateBy // 添加这而，需要记录是被谁更新的
	model.Layer = s.Layer
	model.Enable = s.Enable
	model.Desc = s.Desc

	model.Type = s.Type
	model.StartWeek = s.StartWeek
	model.EndWeek = s.EndWeek
	model.StartTime = s.StartTime
	model.EndTime = s.EndTime
	model.GiveDay = s.GiveDay
	model.GiveTime = s.GiveTime
}

func (s *CycleTimeConfUpdateReq) GetId() interface{} {
	return s.Id
}

// CycleTimeConfGetReq 功能获取请求参数
type CycleTimeConfGetReq struct {
	Id int `uri:"id"`
}

func (s *CycleTimeConfGetReq) GetId() interface{} {
	return s.Id
}

// CycleTimeConfDeleteReq 功能删除请求参数
type CycleTimeConfDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *CycleTimeConfDeleteReq) GetId() interface{} {
	return s.Ids
}
