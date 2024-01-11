package dto

import "go-admin/common/dto"

type CompanyInventoryCnfReq struct {

	Enable bool `json:"enable"`
}

type InventoryGoodsReq struct {
	dto.Pagination `search:"-"`
	Name string `json:"name" form:"name"`

}

type GoodsSpecs struct {
	Key string `json:"key"`
	Name string `json:"name"`
	Unit string `json:"unit"`
	Price float64 `json:"price"`
	Stock int `json:"stock"`
	Image string `json:"image"`
}

type ManageListGetPageReq struct {
	dto.Pagination `search:"-"`
	GoodsName string `form:"goods_name" search:"type:contains;column:goods_name;table:inventory" comment:"商品名称"`
	BeginTime      string `form:"beginTime" search:"type:gte;column:created_at;table:inventory" comment:"创建时间"`
	EndTime        string `form:"endTime" search:"type:lte;column:created_at;table:inventory" comment:"创建时间"`
}
func (m *ManageListGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type RecordsListGetPageReq struct {
	dto.Pagination `search:"-"`
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
type InventoryCreateReq struct {
	Desc string `json:"desc"`
	Data map[string]WarehousingRow `json:"data"`
}
type WarehousingRow struct {
	ActionNumber int `json:"action_number"`
	CostPrice float64 `json:"cost_price"`
	Unit string `json:"unit"`
}