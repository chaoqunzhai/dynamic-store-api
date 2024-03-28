package dto

import "go-admin/common/dto"

type CompanyInventoryCnfReq struct {

	Enable bool `json:"enable"`
}

type InventoryGoodsReq struct {
	dto.Pagination `search:"-"`
	Class string `json:"class" form:"class" search:"-"`
	Brand string `json:"brand" form:"brand" search:"-"`
	Name string `json:"name" form:"name" search:"type:contains;column:goods_name;table:inventory" comment:"商品名称"`
}
func (m *InventoryGoodsReq) GetNeedSearch() interface{} {
	return *m
}

type GoodsSpecs struct {
	Key string `json:"key"`
	Name string `json:"name"`
	Unit string `json:"unit"`
	CostPrice float64 `json:"cost_price"`
	Stock int `json:"stock"`
	Image string `json:"image"`
	Code string `json:"code"`
	SerialNumber string `json:"serial_number"`
	ArtNo string `json:"art_no"`
}

type ManageListGetPageReq struct {
	dto.Pagination `search:"-"`
	Class string `json:"class" form:"class" search:"-"`
	Name string `json:"name" form:"name" search:"-"`
	BeginTime      string `form:"beginTime" search:"type:gte;column:created_at;table:inventory" comment:"创建时间"`
	EndTime        string `form:"endTime" search:"type:lte;column:created_at;table:inventory" comment:"创建时间"`
	Action string `json:"action"  form:"action" search:"-"`
	Brand string `json:"brand" form:"brand" search:"-"`
}
func (m *ManageListGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type RecordsListGetPageReq struct {
	dto.Pagination `search:"-"`
	Type int `form:"type" search:"-" `
	CreateBy string `form:"create_by" search:"type:exact;column:create_by;table:inventory_record" comment:"操作方法"`
	Action string `form:"action" search:"type:exact;column:action;table:inventory_record" comment:"操作方法"`
	GoodsName string `form:"goods_name" search:"type:contains;column:goods_name;table:inventory_record" comment:"商品名称"`
	BeginTime      string `form:"beginTime" search:"type:gte;column:created_at;table:inventory_record" comment:"创建时间"`
	EndTime        string `form:"endTime" search:"type:lte;column:created_at;table:inventory_record" comment:"创建时间"`
}
func (m *RecordsListGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type OrderListGetPageReq struct {
	dto.Pagination `search:"-"`
	BeginTime      string `form:"beginTime" search:"type:gte;column:created_at;table:inventory_order" comment:"创建时间"`
	EndTime        string `form:"endTime" search:"type:lte;column:created_at;table:inventory_order" comment:"创建时间"`
	OrderId string `json:"order_id" search:"-" form:"order_id"`
	Action string `json:"action"  search:"-" form:"action"` //出入库类型
}
func (m *OrderListGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type EditReq struct {
	Id int `json:"id"`
	EditArtNo string `json:"edit_art_no"`
	EditCode string `json:"edit_code"`
	EditOriginalPrice float64 `json:"edit_original_price"`
}
type InventoryCreateReq struct {
	Desc string `json:"desc"`
	Data map[string]WarehousingRow `json:"data"`
}
type WarehousingRow struct {
	ActionNumber int `json:"action_number"`
	CostPrice float64 `json:"cost_price"`
	Unit string `json:"unit"`
	ArtNo string `json:"art_no" `
	Code      string  `json:"code"`
	SerialNumber string `json:"serial_number"`
}
type GoodsInfo struct {
	Name string `json:"name"`
	SpecName string `json:"spec_name"`
	Unit string `json:"unit"`
	Image string `json:"image"`
}
