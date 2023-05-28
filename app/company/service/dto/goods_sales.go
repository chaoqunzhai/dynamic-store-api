package dto

import (
     
     
     
     
     "go-admin/app/company/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type GoodsSalesGetPageReq struct {
	dto.Pagination     `search:"-"`
    Layer string `form:"layer"  search:"type:exact;column:layer;table:goods_sales" comment:"排序"`
    Enable string `form:"enable"  search:"type:exact;column:enable;table:goods_sales" comment:"开关"`
    CId string `form:"cId"  search:"type:exact;column:c_id;table:goods_sales" comment:"大BID"`
    ProductId string `form:"productId"  search:"type:exact;column:product_id;table:goods_sales" comment:"产品ID"`
    ProductName string `form:"productName"  search:"type:contains;column:product_name;table:goods_sales" comment:"产品名称"`
    GoodsSalesOrder
}

type GoodsSalesOrder struct {
    Id string `form:"idOrder"  search:"type:order;column:id;table:goods_sales"`
    CreateBy string `form:"createByOrder"  search:"type:order;column:create_by;table:goods_sales"`
    CreatedAt string `form:"createdAtOrder"  search:"type:order;column:created_at;table:goods_sales"`
    DeletedAt string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:goods_sales"`
    Layer string `form:"layerOrder"  search:"type:order;column:layer;table:goods_sales"`
    Enable string `form:"enableOrder"  search:"type:order;column:enable;table:goods_sales"`
    CId string `form:"cIdOrder"  search:"type:order;column:c_id;table:goods_sales"`
    ProductId string `form:"productIdOrder"  search:"type:order;column:product_id;table:goods_sales"`
    ProductName string `form:"productNameOrder"  search:"type:order;column:product_name;table:goods_sales"`
    Sales string `form:"salesOrder"  search:"type:order;column:sales;table:goods_sales"`
    Inventory string `form:"inventoryOrder"  search:"type:order;column:inventory;table:goods_sales"`
    
}

func (m *GoodsSalesGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type GoodsSalesInsertReq struct {
    Id int `json:"-" comment:"主键编码"` // 主键编码
    Layer string `json:"layer" comment:"排序"`
    Enable string `json:"enable" comment:"开关"`
    CId string `json:"cId" comment:"大BID"`
    ProductId string `json:"productId" comment:"产品ID"`
    ProductName string `json:"productName" comment:"产品名称"`
    Sales string `json:"sales" comment:"当时销量"`
    Inventory string `json:"inventory" comment:"当时剩余库存"`
    common.ControlBy
}

func (s *GoodsSalesInsertReq) Generate(model *models.GoodsSales)  {
    if s.Id == 0 {
        model.Model = common.Model{ Id: s.Id }
    }
    model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
    model.Layer = s.Layer
    model.Enable = s.Enable
    model.CId = s.CId
    model.ProductId = s.ProductId
    model.ProductName = s.ProductName
    model.Sales = s.Sales
    model.Inventory = s.Inventory
}

func (s *GoodsSalesInsertReq) GetId() interface{} {
	return s.Id
}

type GoodsSalesUpdateReq struct {
    Id int `uri:"id" comment:"主键编码"` // 主键编码
    Layer string `json:"layer" comment:"排序"`
    Enable string `json:"enable" comment:"开关"`
    CId string `json:"cId" comment:"大BID"`
    ProductId string `json:"productId" comment:"产品ID"`
    ProductName string `json:"productName" comment:"产品名称"`
    Sales string `json:"sales" comment:"当时销量"`
    Inventory string `json:"inventory" comment:"当时剩余库存"`
    common.ControlBy
}

func (s *GoodsSalesUpdateReq) Generate(model *models.GoodsSales)  {
    if s.Id == 0 {
        model.Model = common.Model{ Id: s.Id }
    }
    model.Layer = s.Layer
    model.Enable = s.Enable
    model.CId = s.CId
    model.ProductId = s.ProductId
    model.ProductName = s.ProductName
    model.Sales = s.Sales
    model.Inventory = s.Inventory
}

func (s *GoodsSalesUpdateReq) GetId() interface{} {
	return s.Id
}

// GoodsSalesGetReq 功能获取请求参数
type GoodsSalesGetReq struct {
     Id int `uri:"id"`
}
func (s *GoodsSalesGetReq) GetId() interface{} {
	return s.Id
}

// GoodsSalesDeleteReq 功能删除请求参数
type GoodsSalesDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *GoodsSalesDeleteReq) GetId() interface{} {
	return s.Ids
}
