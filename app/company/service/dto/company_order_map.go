package dto

import (
     
     
     

	"go-admin/app/company/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type CompanyOrderMapGetPageReq struct {
	dto.Pagination     `search:"-"`
    Layer string `form:"layer"  search:"type:exact;column:layer;table:company_order_map" comment:"排序"`
    Enable string `form:"enable"  search:"type:exact;column:enable;table:company_order_map" comment:"开关"`
    CId string `form:"cId"  search:"type:exact;column:c_id;table:company_order_map" comment:"公司ID"`
    Type string `form:"type"  search:"type:exact;column:type;table:company_order_map" comment:"映射表的类型"`
    OrderTable string `form:"orderTable"  search:"type:exact;column:order_table;table:company_order_map" comment:"对应表的名称"`
    CompanyOrderMapOrder
}

type CompanyOrderMapOrder struct {
    Id string `form:"idOrder"  search:"type:order;column:id;table:company_order_map"`
    CreateBy string `form:"createByOrder"  search:"type:order;column:create_by;table:company_order_map"`
    CreatedAt string `form:"createdAtOrder"  search:"type:order;column:created_at;table:company_order_map"`
    DeletedAt string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:company_order_map"`
    Layer string `form:"layerOrder"  search:"type:order;column:layer;table:company_order_map"`
    Enable string `form:"enableOrder"  search:"type:order;column:enable;table:company_order_map"`
    CId string `form:"cIdOrder"  search:"type:order;column:c_id;table:company_order_map"`
    Type string `form:"typeOrder"  search:"type:order;column:type;table:company_order_map"`
    OrderTable string `form:"orderTableOrder"  search:"type:order;column:order_table;table:company_order_map"`
    
}

func (m *CompanyOrderMapGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type CompanyOrderMapInsertReq struct {
    Id int `json:"-" comment:"主键编码"` // 主键编码
    Layer string `json:"layer" comment:"排序"`
    Enable string `json:"enable" comment:"开关"`
    CId string `json:"cId" comment:"公司ID"`
    Type string `json:"type" comment:"映射表的类型"`
    OrderTable string `json:"orderTable" comment:"对应表的名称"`
    common.ControlBy
}

func (s *CompanyOrderMapInsertReq) Generate(model *models.CompanyOrderMap)  {
    if s.Id == 0 {
        model.Model = common.Model{ Id: s.Id }
    }
    model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
    model.Layer = s.Layer
    model.Enable = s.Enable
    model.CId = s.CId
    model.Type = s.Type
    model.OrderTable = s.OrderTable
}

func (s *CompanyOrderMapInsertReq) GetId() interface{} {
	return s.Id
}

type CompanyOrderMapUpdateReq struct {
    Id int `uri:"id" comment:"主键编码"` // 主键编码
    Layer string `json:"layer" comment:"排序"`
    Enable string `json:"enable" comment:"开关"`
    CId string `json:"cId" comment:"公司ID"`
    Type string `json:"type" comment:"映射表的类型"`
    OrderTable string `json:"orderTable" comment:"对应表的名称"`
    common.ControlBy
}

func (s *CompanyOrderMapUpdateReq) Generate(model *models.CompanyOrderMap)  {
    if s.Id == 0 {
        model.Model = common.Model{ Id: s.Id }
    }
    model.Layer = s.Layer
    model.Enable = s.Enable
    model.CId = s.CId
    model.Type = s.Type
    model.OrderTable = s.OrderTable
}

func (s *CompanyOrderMapUpdateReq) GetId() interface{} {
	return s.Id
}

// CompanyOrderMapGetReq 功能获取请求参数
type CompanyOrderMapGetReq struct {
     Id int `uri:"id"`
}
func (s *CompanyOrderMapGetReq) GetId() interface{} {
	return s.Id
}

// CompanyOrderMapDeleteReq 功能删除请求参数
type CompanyOrderMapDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *CompanyOrderMapDeleteReq) GetId() interface{} {
	return s.Ids
}
