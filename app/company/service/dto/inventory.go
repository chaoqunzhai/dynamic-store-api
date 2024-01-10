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


type WarehousingGetPageReq struct {
	dto.Pagination `search:"-"`
	OrderId string `json:"order_id" form:"order_id"`

}

type WarehousingCreateReq struct {
	Desc string `json:"desc"`
	Data map[string]WarehousingRow `json:"data"`
}
type WarehousingRow struct {
	ActionNumber int `json:"action_number"`
	CostPrice float64 `json:"cost_price"`
	Unit string `json:"unit"`
}