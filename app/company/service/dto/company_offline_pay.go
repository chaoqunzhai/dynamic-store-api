package dto

import (
     
     
     
     "go-admin/app/company/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type CompanyOfflinePayGetPageReq struct {
	dto.Pagination     `search:"-"`
    CId string `form:"cId"  search:"type:exact;column:c_id;table:company_offline_pay" comment:"大BID"`
    Name string `form:"name"  search:"type:exact;column:name;table:company_offline_pay" comment:"线下支付名称"`
    Layer string `form:"layer"  search:"type:exact;column:layer;table:company_offline_pay" comment:"排序"`
    Enable string `form:"enable"  search:"type:exact;column:enable;table:company_offline_pay" comment:"开关"`
    CompanyOfflinePayOrder
}

type CompanyOfflinePayOrder struct {
    Id string `form:"idOrder"  search:"type:order;column:id;table:company_offline_pay"`
    CreatedAt string `form:"createdAtOrder"  search:"type:order;column:created_at;table:company_offline_pay"`
    CId string `form:"cIdOrder"  search:"type:order;column:c_id;table:company_offline_pay"`
    Name string `form:"nameOrder"  search:"type:order;column:name;table:company_offline_pay"`
    CreateBy string `form:"createByOrder"  search:"type:order;column:create_by;table:company_offline_pay"`
    UpdateBy string `form:"updateByOrder"  search:"type:order;column:update_by;table:company_offline_pay"`
    UpdatedAt string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:company_offline_pay"`
    DeletedAt string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:company_offline_pay"`
    Layer string `form:"layerOrder"  search:"type:order;column:layer;table:company_offline_pay"`
    Enable string `form:"enableOrder"  search:"type:order;column:enable;table:company_offline_pay"`
    Desc string `form:"descOrder"  search:"type:order;column:desc;table:company_offline_pay"`
    
}

func (m *CompanyOfflinePayGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type CompanyOfflinePayInsertReq struct {
    Id int `json:"-" comment:"主键编码"` // 主键编码
    CId int `json:"-" comment:"大BID"`
    Name string `json:"name" comment:"线下支付名称"`
    Layer string `json:"layer" comment:"排序"`

    Desc string `json:"desc" comment:"描述信息"`
    common.ControlBy
}

func (s *CompanyOfflinePayInsertReq) Generate(model *models.CompanyOfflinePay)  {
    if s.Id == 0 {
        model.Model = common.Model{ Id: s.Id }
    }
    model.CId = s.CId
    model.Name = s.Name
    model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
    model.Layer = s.Layer

    model.Desc = s.Desc
}

func (s *CompanyOfflinePayInsertReq) GetId() interface{} {
	return s.Id
}

type CompanyOfflinePayUpdateReq struct {
    Id int `uri:"id" comment:"主键编码"` // 主键编码
    CId int `json:"cId" comment:"大BID"`
    Name string `json:"name" comment:"线下支付名称"`
    Layer string `json:"layer" comment:"排序"`

    Desc string `json:"desc" comment:"描述信息"`
    common.ControlBy
}

func (s *CompanyOfflinePayUpdateReq) Generate(model *models.CompanyOfflinePay)  {
    if s.Id == 0 {
        model.Model = common.Model{ Id: s.Id }
    }
    model.CId = s.CId
    model.Name = s.Name
    model.UpdateBy = s.UpdateBy // 添加这而，需要记录是被谁更新的
    model.Layer = s.Layer

    model.Desc = s.Desc
}

func (s *CompanyOfflinePayUpdateReq) GetId() interface{} {
	return s.Id
}

// CompanyOfflinePayGetReq 功能获取请求参数
type CompanyOfflinePayGetReq struct {
     Id int `uri:"id"`
}
func (s *CompanyOfflinePayGetReq) GetId() interface{} {
	return s.Id
}

// CompanyOfflinePayDeleteReq 功能删除请求参数
type CompanyOfflinePayDeleteReq struct {
    CId int
    Id int `json:"id"`
}

func (s *CompanyOfflinePayDeleteReq) GetId() interface{} {
	return s.Id
}
