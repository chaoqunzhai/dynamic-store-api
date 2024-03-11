package dto

import (
	"go-admin/app/company/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type GoodsSpecsGetPageReq struct {
	dto.Pagination `search:"-"`
	ShopId 		 int `form:"shop_id" search:"-"`
	Layer          string `form:"layer"  search:"type:exact;column:layer;table:goods_specs" comment:"排序"`
	Enable         string `form:"enable"  search:"type:exact;column:enable;table:goods_specs" comment:"开关"`
	CId            string `form:"cId"  search:"type:exact;column:c_id;table:goods_specs" comment:"大BID"`
	GoodsId        string `form:"goods_id"  search:"type:exact;column:goods_id;table:goods_specs" comment:"商品ID"`
	Name           string `form:"name"  search:"type:contains;column:name;table:goods_specs" comment:"规格名称"`
	Unit           string `form:"unit"  search:"type:exact;column:unit;table:goods_specs" comment:"单位"`
	GoodsSpecsOrder
}

type GoodsSpecsOrder struct {
	Id        string `form:"idOrder"  search:"type:order;column:id;table:goods_specs"`
	CreateBy  string `form:"createByOrder"  search:"type:order;column:create_by;table:goods_specs"`
	CreatedAt string `form:"createdAtOrder"  search:"type:order;column:created_at;table:goods_specs"`
	DeletedAt string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:goods_specs"`
	Layer     string `form:"layerOrder"  search:"type:order;column:layer;table:goods_specs"`
	Enable    string `form:"enableOrder"  search:"type:order;column:enable;table:goods_specs"`
	CId       string `form:"cIdOrder"  search:"type:order;column:c_id;table:goods_specs"`
	GoodsId   string `form:"goodsIdOrder"  search:"type:order;column:goods_id;table:goods_specs"`
	Name      string `form:"nameOrder"  search:"type:order;column:name;table:goods_specs"`
	Price     string `form:"priceOrder"  search:"type:order;column:price;table:goods_specs"`
	Original  string `form:"originalOrder"  search:"type:order;column:original;table:goods_specs"`
	Inventory string `form:"inventoryOrder"  search:"type:order;column:inventory;table:goods_specs"`
	Unit      string `form:"unitOrder"  search:"type:order;column:unit;table:goods_specs"`
	Limit     string `form:"limitOrder"  search:"type:order;column:limit;table:goods_specs"`
}

func (m *GoodsSpecsGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type GoodsSpecsInsertReq struct {
	Id     int  `json:"-" comment:"主键编码"` // 主键编码
	Layer  int  `json:"layer" comment:"排序"`
	Enable bool `json:"enable" comment:"开关"`

	GoodsId   int     `json:"goodsId" comment:"商品ID"`
	Name      string  `json:"name" comment:"规格名称"`
	Price     float64 `json:"price" comment:"售价"`
	Original  float64 `json:"original" comment:"原价"`
	Inventory int     `json:"inventory" comment:"库存"`
	UnitId      int  `json:"unit" comment:"单位"`
	Limit     int     `json:"limit" comment:"起售量"`
	Max int `json:"max"`
	common.ControlBy
}

func (s *GoodsSpecsInsertReq) Generate(model *models.GoodsSpecs) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
	model.Layer = s.Layer
	model.Enable = s.Enable

	model.GoodsId = s.GoodsId
	model.Name = s.Name
	model.Price = s.Price
	model.Original = s.Original
	model.Inventory = s.Inventory
	model.UnitId = s.UnitId
	model.Limit = s.Limit
	model.Max = s.Max
}

func (s *GoodsSpecsInsertReq) GetId() interface{} {
	return s.Id
}

type GoodsSpecsUpdateReq struct {
	Id     int  `uri:"id" comment:"主键编码"` // 主键编码
	Layer  int  `json:"layer" comment:"排序"`
	Enable bool `json:"enable" comment:"开关"`

	GoodsId   int     `json:"goodsId" comment:"商品ID"`
	Name      string  `json:"name" comment:"规格名称"`
	Price     float64 `json:"price" comment:"售价"`
	Original  float64 `json:"original" comment:"原价"`
	Inventory int     `json:"inventory" comment:"库存"`
	UnitId      int  `json:"unit" comment:"单位"`
	Limit     int     `json:"limit" comment:"起售量"`
	Max int `json:"max"`
	common.ControlBy
}

func (s *GoodsSpecsUpdateReq) Generate(model *models.GoodsSpecs) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.Layer = s.Layer
	model.Enable = s.Enable

	model.GoodsId = s.GoodsId
	model.Name = s.Name
	model.Price = s.Price
	model.Original = s.Original
	model.Inventory = s.Inventory
	model.UnitId = s.UnitId
	model.Limit = s.Limit
	model.Max = s.Max
}

func (s *GoodsSpecsUpdateReq) GetId() interface{} {
	return s.Id
}

// GoodsSpecsGetReq 功能获取请求参数
type GoodsSpecsGetReq struct {
	Id int `uri:"id"`
}

func (s *GoodsSpecsGetReq) GetId() interface{} {
	return s.Id
}

// GoodsSpecsDeleteReq 功能删除请求参数
type GoodsSpecsDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *GoodsSpecsDeleteReq) GetId() interface{} {
	return s.Ids
}
