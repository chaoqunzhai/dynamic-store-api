package dto

import (
     
     


	"go-admin/app/shop/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type ShopOrderBindRecordGetPageReq struct {
	dto.Pagination     `search:"-"`
    ShopId string `form:"shopId"  search:"type:exact;column:shop_id;table:shop_order_bind_record" comment:"关联的小B客户"`
    RecordId string `form:"recordId"  search:"type:exact;column:record_id;table:shop_order_bind_record" comment:"每次记录的总ID"`
    OrderId string `form:"orderId"  search:"type:exact;column:order_id;table:shop_order_bind_record" comment:"订单ID"`
    ShopOrderBindRecordOrder
}

type ShopOrderBindRecordOrder struct {
    Id string `form:"idOrder"  search:"type:order;column:id;table:shop_order_bind_record"`
    CreateBy string `form:"createByOrder"  search:"type:order;column:create_by;table:shop_order_bind_record"`
    CreatedAt string `form:"createdAtOrder"  search:"type:order;column:created_at;table:shop_order_bind_record"`
    DeletedAt string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:shop_order_bind_record"`
    ShopId string `form:"shopIdOrder"  search:"type:order;column:shop_id;table:shop_order_bind_record"`
    RecordId string `form:"recordIdOrder"  search:"type:order;column:record_id;table:shop_order_bind_record"`
    OrderId string `form:"orderIdOrder"  search:"type:order;column:order_id;table:shop_order_bind_record"`
    
}

func (m *ShopOrderBindRecordGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type ShopOrderBindRecordInsertReq struct {
    Id int `json:"-" comment:"主键编码"` // 主键编码
    ShopId string `json:"shopId" comment:"关联的小B客户"`
    RecordId string `json:"recordId" comment:"每次记录的总ID"`
    OrderId string `json:"orderId" comment:"订单ID"`
    common.ControlBy
}

func (s *ShopOrderBindRecordInsertReq) Generate(model *models.ShopOrderBindRecord)  {
    if s.Id == 0 {
        model.Model = common.Model{ Id: s.Id }
    }
    model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
    model.ShopId = s.ShopId
    model.RecordId = s.RecordId
    model.OrderId = s.OrderId
}

func (s *ShopOrderBindRecordInsertReq) GetId() interface{} {
	return s.Id
}

type ShopOrderBindRecordUpdateReq struct {
    Id int `uri:"id" comment:"主键编码"` // 主键编码
    ShopId string `json:"shopId" comment:"关联的小B客户"`
    RecordId string `json:"recordId" comment:"每次记录的总ID"`
    OrderId string `json:"orderId" comment:"订单ID"`
    common.ControlBy
}

func (s *ShopOrderBindRecordUpdateReq) Generate(model *models.ShopOrderBindRecord)  {
    if s.Id == 0 {
        model.Model = common.Model{ Id: s.Id }
    }
    model.ShopId = s.ShopId
    model.RecordId = s.RecordId
    model.OrderId = s.OrderId
}

func (s *ShopOrderBindRecordUpdateReq) GetId() interface{} {
	return s.Id
}

// ShopOrderBindRecordGetReq 功能获取请求参数
type ShopOrderBindRecordGetReq struct {
     Id int `uri:"id"`
}
func (s *ShopOrderBindRecordGetReq) GetId() interface{} {
	return s.Id
}

// ShopOrderBindRecordDeleteReq 功能删除请求参数
type ShopOrderBindRecordDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *ShopOrderBindRecordDeleteReq) GetId() interface{} {
	return s.Ids
}
