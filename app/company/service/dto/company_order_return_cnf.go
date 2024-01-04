package dto

import (
     
     
     "go-admin/app/company/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type CompanyOrderReturnCnfGetPageReq struct {
	dto.Pagination     `search:"-"`
    Enable string `form:"enable"  search:"type:exact;column:enable;table:company_order_return_cnf" comment:"开关"`
    CId string `form:"cId"  search:"type:exact;column:c_id;table:company_order_return_cnf" comment:"大BID"`
    Value string `form:"value"  search:"type:exact;column:value;table:company_order_return_cnf" comment:"配送文案"`
    CompanyOrderReturnCnfOrder
}

type CompanyOrderReturnCnfOrder struct {
    Id string `form:"idOrder"  search:"type:order;column:id;table:company_order_return_cnf"`
    CreateBy string `form:"createByOrder"  search:"type:order;column:create_by;table:company_order_return_cnf"`
    UpdateBy string `form:"updateByOrder"  search:"type:order;column:update_by;table:company_order_return_cnf"`
    CreatedAt string `form:"createdAtOrder"  search:"type:order;column:created_at;table:company_order_return_cnf"`
    UpdatedAt string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:company_order_return_cnf"`
    DeletedAt string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:company_order_return_cnf"`
    Layer string `form:"layerOrder"  search:"type:order;column:layer;table:company_order_return_cnf"`
    Enable string `form:"enableOrder"  search:"type:order;column:enable;table:company_order_return_cnf"`
    Desc string `form:"descOrder"  search:"type:order;column:desc;table:company_order_return_cnf"`
    CId string `form:"cIdOrder"  search:"type:order;column:c_id;table:company_order_return_cnf"`
    Value string `form:"valueOrder"  search:"type:order;column:value;table:company_order_return_cnf"`
    Cost string `form:"costOrder"  search:"type:order;column:cost;table:company_order_return_cnf"`
    
}

func (m *CompanyOrderReturnCnfGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type CompanyOrderReturnCnfInsertReq struct {
    Id int `json:"-" comment:"主键编码"` // 主键编码
    Layer string `json:"layer" comment:"排序"`
    Enable bool `json:"enable" comment:"开关"`
    Desc string `json:"desc" comment:"描述信息"`
    CId int `json:"-" comment:"大BID"`
    Value string `json:"value" comment:"配送文案"`
    Cost float64 `json:"cost" comment:"配送费用"`
    common.ControlBy
}

func (s *CompanyOrderReturnCnfInsertReq) Generate(model *models.CompanyOrderReturnCnf)  {
    if s.Id == 0 {
        model.Model = common.Model{ Id: s.Id }
    }
    model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
    model.Layer = s.Layer
    model.Enable = s.Enable
    model.Desc = s.Desc
    model.CId = s.CId
    model.Value = s.Value
    model.Cost = s.Cost
}

func (s *CompanyOrderReturnCnfInsertReq) GetId() interface{} {
	return s.Id
}

type CompanyOrderReturnCnfUpdateReq struct {
    Id int `uri:"id" comment:"主键编码"` // 主键编码
    Layer string `json:"layer" comment:"排序"`
    Enable bool `json:"enable" comment:"开关"`
    Desc string `json:"desc" comment:"描述信息"`
    CId int `json:"cId" comment:"大BID"`
    Value string `json:"value" comment:"配送文案"`
    Cost float64 `json:"cost" comment:"配送费用"`
    common.ControlBy
}

func (s *CompanyOrderReturnCnfUpdateReq) Generate(model *models.CompanyOrderReturnCnf)  {
    if s.Id == 0 {
        model.Model = common.Model{ Id: s.Id }
    }
    model.UpdateBy = s.UpdateBy // 添加这而，需要记录是被谁更新的
    model.Layer = s.Layer
    model.Enable = s.Enable
    model.Desc = s.Desc
    model.CId = s.CId
    model.Value = s.Value
    model.Cost = s.Cost
}

func (s *CompanyOrderReturnCnfUpdateReq) GetId() interface{} {
	return s.Id
}

// CompanyOrderReturnCnfGetReq 功能获取请求参数
type CompanyOrderReturnCnfGetReq struct {
     Id int `uri:"id"`
}
func (s *CompanyOrderReturnCnfGetReq) GetId() interface{} {
	return s.Id
}

// CompanyOrderReturnCnfDeleteReq 功能删除请求参数
type CompanyOrderReturnCnfDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *CompanyOrderReturnCnfDeleteReq) GetId() interface{} {
	return s.Ids
}
