package dto

import (
	"go-admin/app/system/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type SplitTableMapGetPageReq struct {
	dto.Pagination `search:"-"`
	Layer          string `form:"layer"  search:"type:exact;column:layer;table:split_table_map" comment:"排序"`
	Enable         bool   `form:"enable"  search:"type:exact;column:enable;table:split_table_map" comment:"开关"`
	CId            int    `form:"cId"  search:"type:exact;column:c_id;table:split_table_map" comment:"公司ID"`
	Type           int    `form:"type"  search:"type:exact;column:type;table:split_table_map" comment:"映射表的类型"`
	Table          string `form:"table"  search:"type:exact;column:table;table:split_table_map" comment:"对应表的名称"`
	SplitTableMapOrder
}

type SplitTableMapOrder struct {
	Id        string `form:"idOrder"  search:"type:order;column:id;table:split_table_map"`
	CreateBy  string `form:"createByOrder"  search:"type:order;column:create_by;table:split_table_map"`
	CreatedAt string `form:"createdAtOrder"  search:"type:order;column:created_at;table:split_table_map"`
	DeletedAt string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:split_table_map"`
	Layer     string `form:"layerOrder"  search:"type:order;column:layer;table:split_table_map"`
	Enable    string `form:"enableOrder"  search:"type:order;column:enable;table:split_table_map"`
	CId       string `form:"cIdOrder"  search:"type:order;column:c_id;table:split_table_map"`
	Type      string `form:"typeOrder"  search:"type:order;column:type;table:split_table_map"`
	Table     string `form:"tableOrder"  search:"type:order;column:table;table:split_table_map"`
	Desc      string `form:"descOrder"  search:"type:order;column:desc;table:split_table_map"`
}

func (m *SplitTableMapGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type SplitTableMapInsertReq struct {
	Id     int    `json:"-" comment:"主键编码"` // 主键编码
	Layer  int    `json:"layer" comment:"排序"`
	Enable bool   `json:"enable" comment:"开关"`
	CId    int    `json:"c_id" comment:"公司ID"`
	Type   int    `json:"type" comment:"映射表的类型"`
	Name   string `json:"table" comment:"对应表的名称"`
	Desc   string `json:"desc" comment:"对应表的名称"`
	common.ControlBy
}

func (s *SplitTableMapInsertReq) Generate(model *models.SplitTableMap) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
	model.Layer = s.Layer
	model.Enable = s.Enable
	model.CId = s.CId
	model.Type = s.Type
	model.Name = s.Name
	model.Desc = s.Desc
}

func (s *SplitTableMapInsertReq) GetId() interface{} {
	return s.Id
}

type SplitTableMapUpdateReq struct {
	Id     int    `uri:"id" comment:"主键编码"` // 主键编码
	Layer  int    `json:"layer" comment:"排序"`
	Enable bool   `json:"enable" comment:"开关"`
	CId    int    `json:"cId" comment:"公司ID"`
	Type   int    `json:"type" comment:"映射表的类型"`
	Name   string `json:"name" comment:"对应表的名称"`
	Desc   string `json:"desc" comment:"对应表的名称"`
	common.ControlBy
}

func (s *SplitTableMapUpdateReq) Generate(model *models.SplitTableMap) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.Layer = s.Layer
	model.Enable = s.Enable
	model.CId = s.CId
	model.Type = s.Type
	model.Name = s.Name
	model.Desc = s.Desc
}

func (s *SplitTableMapUpdateReq) GetId() interface{} {
	return s.Id
}

// SplitTableMapGetReq 功能获取请求参数
type SplitTableMapGetReq struct {
	Id int `uri:"id"`
}

func (s *SplitTableMapGetReq) GetId() interface{} {
	return s.Id
}

// SplitTableMapDeleteReq 功能删除请求参数
type SplitTableMapDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *SplitTableMapDeleteReq) GetId() interface{} {
	return s.Ids
}
