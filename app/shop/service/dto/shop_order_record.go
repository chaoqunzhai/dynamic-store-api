package dto

import (
     


	"go-admin/app/shop/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type ShopOrderRecordGetPageReq struct {
	dto.Pagination     `search:"-"`
    ShopId string `form:"shopId"  search:"type:exact;column:shop_id;table:shop_order_record" comment:"关联的小B客户"`
    ShopName string `form:"shopName"  search:"type:contains;column:shop_name;table:shop_order_record" comment:"客户名称"`
    Number string `form:"number"  search:"type:exact;column:number;table:shop_order_record" comment:"订单量"`
    ShopOrderRecordOrder
}

type ShopOrderRecordOrder struct {
    Id string `form:"idOrder"  search:"type:order;column:id;table:shop_order_record"`
    CreateBy string `form:"createByOrder"  search:"type:order;column:create_by;table:shop_order_record"`
    CreatedAt string `form:"createdAtOrder"  search:"type:order;column:created_at;table:shop_order_record"`
    DeletedAt string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:shop_order_record"`
    ShopId string `form:"shopIdOrder"  search:"type:order;column:shop_id;table:shop_order_record"`
    ShopName string `form:"shopNameOrder"  search:"type:order;column:shop_name;table:shop_order_record"`
    Money string `form:"moneyOrder"  search:"type:order;column:money;table:shop_order_record"`
    Number string `form:"numberOrder"  search:"type:order;column:number;table:shop_order_record"`
    
}

func (m *ShopOrderRecordGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type ShopOrderRecordInsertReq struct {
    Id int `json:"-" comment:"主键编码"` // 主键编码
    ShopId string `json:"shopId" comment:"关联的小B客户"`
    ShopName string `json:"shopName" comment:"客户名称"`
    Money string `json:"money" comment:"订单金额"`
    Number string `json:"number" comment:"订单量"`
    common.ControlBy
}

func (s *ShopOrderRecordInsertReq) Generate(model *models.ShopOrderRecord)  {
    if s.Id == 0 {
        model.Model = common.Model{ Id: s.Id }
    }
    model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
    model.ShopId = s.ShopId
    model.ShopName = s.ShopName
    model.Money = s.Money
    model.Number = s.Number
}

func (s *ShopOrderRecordInsertReq) GetId() interface{} {
	return s.Id
}

type ShopOrderRecordUpdateReq struct {
    Id int `uri:"id" comment:"主键编码"` // 主键编码
    ShopId string `json:"shopId" comment:"关联的小B客户"`
    ShopName string `json:"shopName" comment:"客户名称"`
    Money string `json:"money" comment:"订单金额"`
    Number string `json:"number" comment:"订单量"`
    common.ControlBy
}

func (s *ShopOrderRecordUpdateReq) Generate(model *models.ShopOrderRecord)  {
    if s.Id == 0 {
        model.Model = common.Model{ Id: s.Id }
    }
    model.ShopId = s.ShopId
    model.ShopName = s.ShopName
    model.Money = s.Money
    model.Number = s.Number
}

func (s *ShopOrderRecordUpdateReq) GetId() interface{} {
	return s.Id
}

// ShopOrderRecordGetReq 功能获取请求参数
type ShopOrderRecordGetReq struct {
     Id int `uri:"id"`
}
func (s *ShopOrderRecordGetReq) GetId() interface{} {
	return s.Id
}

// ShopOrderRecordDeleteReq 功能删除请求参数
type ShopOrderRecordDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *ShopOrderRecordDeleteReq) GetId() interface{} {
	return s.Ids
}
