package dto

import (
     
     
     
     
     
     "go-admin/app/company/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type CompanyDebitCardGetPageReq struct {
	dto.Pagination     `search:"-"`
    Layer string `form:"layer"  search:"type:exact;column:layer;table:company_debit_card" comment:"排序"`
    Enable string `form:"enable"  search:"type:exact;column:enable;table:company_debit_card" comment:"开关"`

    Name string `form:"name"  search:"type:exact;column:name;table:company_debit_card" comment:"持卡人名称"`
    BackName string `form:"back_name"  search:"type:exact;column:back_name;table:company_debit_card" comment:"开户行"`
    CardNumber string `form:"card_number"  search:"type:exact;column:card_number;table:company_debit_card" comment:"银行卡号"`
    CompanyDebitCardOrder
}

type CompanyDebitCardOrder struct {
    Id string `form:"idOrder"  search:"type:order;column:id;table:company_debit_card"`
    CreateBy string `form:"createByOrder"  search:"type:order;column:create_by;table:company_debit_card"`
    UpdateBy string `form:"updateByOrder"  search:"type:order;column:update_by;table:company_debit_card"`
    CreatedAt string `form:"createdAtOrder"  search:"type:order;column:created_at;table:company_debit_card"`
    UpdatedAt string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:company_debit_card"`
    DeletedAt string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:company_debit_card"`
    Layer string `form:"layerOrder"  search:"type:order;column:layer;table:company_debit_card"`
    Enable string `form:"enableOrder"  search:"type:order;column:enable;table:company_debit_card"`
    Desc string `form:"descOrder"  search:"type:order;column:desc;table:company_debit_card"`
    CId string `form:"cIdOrder"  search:"type:order;column:c_id;table:company_debit_card"`
    Name string `form:"nameOrder"  search:"type:order;column:name;table:company_debit_card"`
    BackName string `form:"backNameOrder"  search:"type:order;column:back_name;table:company_debit_card"`
    CardNumber string `form:"cardNumberOrder"  search:"type:order;column:card_number;table:company_debit_card"`
    
}

func (m *CompanyDebitCardGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type CompanyDebitCardInsertReq struct {
    Id int `json:"-" comment:"主键编码"` // 主键编码
    Layer string `json:"layer" comment:"排序"`

    Desc string `json:"desc" comment:"描述信息"`
    CId int `json:"-" comment:"大BID"`
    Bank string `json:"bank" comment:"用户名称"`
    Name string `json:"name" comment:"持卡人名称"`
    BankName string `json:"bank_name" comment:"开户行"`
    CardNumber string `json:"card_number" comment:"银行卡号"`
    common.ControlBy
}

func (s *CompanyDebitCardInsertReq) Generate(model *models.CompanyDebitCard)  {
    if s.Id == 0 {
        model.Model = common.Model{ Id: s.Id }
    }
    model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
    model.Layer = s.Layer

    model.Desc = s.Desc
    model.CId = s.CId
    model.Name = s.Name
    model.Bank = s.Bank
    model.BankName = s.BankName
    model.CardNumber = s.CardNumber
}

func (s *CompanyDebitCardInsertReq) GetId() interface{} {
	return s.Id
}

type CompanyDebitCardUpdateReq struct {
    Id int `uri:"id" comment:"主键编码"` // 主键编码
    Layer string `json:"layer" comment:"排序"`

    Desc string `json:"desc" comment:"描述信息"`
    CId int `json:"-" comment:"大BID"`
    Bank string `json:"bank" comment:"用户名称"`
    Name string `json:"name" comment:"持卡人名称"`
    BankName string `json:"bank_name" comment:"开户行"`
    CardNumber string `json:"card_number" comment:"银行卡号"`
    common.ControlBy
}

func (s *CompanyDebitCardUpdateReq) Generate(model *models.CompanyDebitCard)  {
    if s.Id == 0 {
        model.Model = common.Model{ Id: s.Id }
    }
    model.UpdateBy = s.UpdateBy // 添加这而，需要记录是被谁更新的
    model.Layer = s.Layer

    model.Desc = s.Desc
    model.CId = s.CId
    model.Name = s.Name
    model.Bank = s.Bank
    model.BankName = s.BankName
    model.CardNumber = s.CardNumber
}

func (s *CompanyDebitCardUpdateReq) GetId() interface{} {
	return s.Id
}

// CompanyDebitCardGetReq 功能获取请求参数
type CompanyDebitCardGetReq struct {
     Id int `uri:"id"`
}
func (s *CompanyDebitCardGetReq) GetId() interface{} {
	return s.Id
}

// CompanyDebitCardDeleteReq 功能删除请求参数
type CompanyDebitCardDeleteReq struct {
    CId int
	Id int `json:"id"`
}

func (s *CompanyDebitCardDeleteReq) GetId() interface{} {
	return s.Id
}
