package dto

import (
	"go-admin/app/company/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type LineGetPageReq struct {
	dto.Pagination `search:"-"`
	Layer          string `form:"layer"  search:"type:exact;column:layer;table:line" comment:"排序"`
	Enable         string `form:"enable"  search:"type:exact;column:enable;table:line" comment:"开关"`
	CId            string `form:"cId"  search:"type:exact;column:c_id;table:line" comment:"大BID"`
	Name           string `form:"name"  search:"type:contains;column:name;table:line" comment:"路线名称"`
	DriverId       string `form:"driver_id"  search:"type:exact;column:driver_id;table:line" comment:"关联司机"`
	BeginTime      string `form:"beginTime" search:"type:gte;column:created_at;table:line" comment:"创建时间"`
	EndTime        string `form:"endTime" search:"type:lte;column:created_at;table:line" comment:"创建时间"`
	LineOrder
}

type LineBindShopGetPageReq struct {
	dto.Pagination `search:"-"`
	Layer          string `form:"layer"  search:"type:exact;column:layer;table:line" comment:"排序"`
	Enable         string `form:"enable"  search:"type:exact;column:enable;table:line" comment:"开关"`
	Name           string `form:"name"  search:"type:contains;column:name;table:line" comment:"路线名称"`
	DriverId       string `form:"driver_id"  search:"type:exact;column:driver_id;table:line" comment:"关联司机"`
}

func (m *LineBindShopGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type UpdateLineBindShopReq struct {
	Id        int     `uri:"id" comment:"主键编码"` // 主键编码
	Layer     int     `form:"layer"`
	Enable    bool    `form:"enable"`
	Address   string  `form:"address" `
	Desc      string  `form:"desc" `
	Longitude float64 `json:"longitude" comment:""`
	Latitude  float64 `json:"latitude" comment:""`
}

type BindLineUserReq struct {
	LineId int   `json:"line_id"`
	ShopId []int `json:"shop_id"`
}
type LineOrder struct {
	Id        string `form:"idOrder"  search:"type:order;column:id;table:line"`
	CreateBy  string `form:"createByOrder"  search:"type:order;column:create_by;table:line"`
	UpdateBy  string `form:"updateByOrder"  search:"type:order;column:update_by;table:line"`
	CreatedAt string `form:"createdAtOrder"  search:"type:order;column:created_at;table:line"`
	UpdatedAt string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:line"`
	DeletedAt string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:line"`
	Layer     string `form:"layerOrder"  search:"type:order;column:layer;table:line"`
	Enable    string `form:"enableOrder"  search:"type:order;column:enable;table:line"`
	Desc      string `form:"descOrder"  search:"type:order;column:desc;table:line"`
	CId       string `form:"cIdOrder"  search:"type:order;column:c_id;table:line"`
	Name      string `form:"nameOrder"  search:"type:order;column:name;table:line"`
	DriverId  string `form:"driverIdOrder"  search:"type:order;column:driver_id;table:line"`
}

func (m *LineGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type LineInsertReq struct {
	Id       int    `json:"-" comment:"主键编码"` // 主键编码
	Layer    int    `json:"layer" comment:"排序"`
	Enable   bool   `json:"enable" comment:"开关"`
	Desc     string `json:"desc" comment:"描述信息"`
	Name     string `json:"name" comment:"路线名称"  binding:"required"`
	DriverId int    `json:"driver_id" comment:"关联司机" `
	common.ControlBy
}

func (s *LineInsertReq) Generate(model *models.Line) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
	model.Layer = s.Layer
	model.Enable = s.Enable
	model.Desc = s.Desc
	model.Name = s.Name
	model.DriverId = s.DriverId
}

func (s *LineInsertReq) GetId() interface{} {
	return s.Id
}

type LineUpdateReq struct {
	Id       int    `uri:"id" comment:"主键编码"` // 主键编码
	Layer    int    `json:"layer" comment:"排序"`
	Enable   bool   `json:"enable" comment:"开关"`
	Desc     string `json:"desc" comment:"描述信息"`
	Name     string `json:"name" comment:"路线名称"  binding:"required"`
	DriverId int    `json:"driver_id" comment:"关联司机" `
	common.ControlBy
}

func (s *LineUpdateReq) Generate(model *models.Line) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.UpdateBy = s.UpdateBy // 添加这而，需要记录是被谁更新的
	model.Layer = s.Layer
	model.Enable = s.Enable
	model.Desc = s.Desc
	model.Name = s.Name
	model.DriverId = s.DriverId
}

func (s *LineUpdateReq) GetId() interface{} {
	return s.Id
}

// LineGetReq 功能获取请求参数
type LineGetReq struct {
	Id int `uri:"id"`
}

func (s *LineGetReq) GetId() interface{} {
	return s.Id
}

