package dto

import (
	"go-admin/app/company/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type CompanyRoleGetPageReq struct {
	dto.Pagination `search:"-"`
	Type           int `form:"type"  search:"type:exact;column:type;table:company_role" comment:""`
	Name           string `form:"name"  search:"type:contains;column:name;table:company_role" comment:""`
	Enable         string `form:"enable"  search:"type:exact;column:enable;table:company_role" comment:""`
	Layer          string `form:"layer"  search:"type:exact;column:layer;table:company_role" comment:""`
	Remark         string `form:"remark"  search:"type:exact;column:remark;table:company_role" comment:""`
	Admin          string `form:"admin"  search:"type:exact;column:admin;table:company_role" comment:""`
	CompanyRoleOrder
}

type CompanyRoleOrder struct {
	Id        string `form:"idOrder"  search:"type:order;column:id;table:company_role"`
	Name      string `form:"nameOrder"  search:"type:order;column:name;table:company_role"`
	Enable    string `form:"enableOrder"  search:"type:order;column:enable;table:company_role"`
	Sort      string `form:"sortOrder"  search:"type:order;column:sort;table:company_role"`
	Remark    string `form:"remarkOrder"  search:"type:order;column:remark;table:company_role"`
	Admin     string `form:"adminOrder"  search:"type:order;column:admin;table:company_role"`
	CreateBy  string `form:"createByOrder"  search:"type:order;column:create_by;table:company_role"`
	UpdateBy  string `form:"updateByOrder"  search:"type:order;column:update_by;table:company_role"`
	CreatedAt string `form:"createdAtOrder"  search:"type:order;column:created_at;table:company_role"`
	UpdatedAt string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:company_role"`
	DeletedAt string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:company_role"`
}

func (m *CompanyRoleGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type CompanyRoleInsertReq struct {
	Id     int    `json:"-" comment:""` //
	Type int `json:"type"`
	Name   string `json:"name" comment:""  binding:"required"`
	Enable bool   `json:"enable" comment:""`
	Layer  int    `json:"layer" comment:""`
	Desc   string `json:"desc" comment:""`
	Admin  bool   `json:"admin" comment:""`
	Menus  []int  `json:"menus"`
	MbmMenus  []int  `json:"mbm_menus"`
	User   []int  `json:"user"`
	common.ControlBy
}

func (s *CompanyRoleInsertReq) Generate(model *models.CompanyRole) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.Type = s.Type
	model.Name = s.Name
	model.Enable = s.Enable
	model.Layer = s.Layer
	model.Desc = s.Desc
	model.Admin = s.Admin
	model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
}

func (s *CompanyRoleInsertReq) GetId() interface{} {
	return s.Id
}

type CompanyRoleUpdateReq struct {
	Id     int    `uri:"id" comment:""` //
	Name   string `json:"name" comment:""  binding:"required"`
	Enable bool   `json:"enable" comment:""`
	Layer  int    `json:"layer" comment:""`
	Desc   string `json:"desc" comment:""`
	Admin  bool   `json:"admin" comment:""`
	Menus  []int  `json:"menus"`
	MbmMenus  []int  `json:"mbm_menus"`
	User   []int  `json:"user"`
	common.ControlBy
}

func (s *CompanyRoleUpdateReq) Generate(model *models.CompanyRole) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.Name = s.Name
	model.Enable = s.Enable
	model.Layer = s.Layer
	model.Desc = s.Desc
	model.Admin = s.Admin
	model.UpdateBy = s.UpdateBy // 添加这而，需要记录是被谁更新的
}

func (s *CompanyRoleUpdateReq) GetId() interface{} {
	return s.Id
}

// CompanyRoleGetReq 功能获取请求参数
type CompanyRoleGetReq struct {
	Id int `uri:"id"`
}

func (s *CompanyRoleGetReq) GetId() interface{} {
	return s.Id
}

// CompanyRoleDeleteReq 功能删除请求参数
type CompanyRoleDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *CompanyRoleDeleteReq) GetId() interface{} {
	return s.Ids
}
