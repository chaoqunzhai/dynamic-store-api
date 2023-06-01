package dto

import (
	"go-admin/app/company/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type OrdersGetPageReq struct {
	dto.Pagination `search:"-"`
	Layer          string `form:"layer"  search:"type:exact;column:layer;table:orders" comment:"排序"`
	Enable         string `form:"enable"  search:"type:exact;column:enable;table:orders" comment:"开关"`
	CId            int    `form:"cId"  search:"type:exact;column:c_id;table:orders" comment:"大BID"`
	ShopId         string `form:"shopId"  search:"type:exact;column:shop_id;table:orders" comment:"关联客户"`
	Status         string `form:"status"  search:"type:exact;column:status;table:orders" comment:"配送状态"`
	Number         string `form:"number"  search:"type:exact;column:number;table:orders" comment:"下单数量"`
	Delivery       string `form:"delivery"  search:"type:exact;column:delivery;table:orders" comment:"配送周期"`
	OrdersOrder
}

type OrdersOrder struct {
	Id        string `form:"idOrder"  search:"type:order;column:id;table:orders"`
	CreateBy  string `form:"createByOrder"  search:"type:order;column:create_by;table:orders"`
	UpdateBy  string `form:"updateByOrder"  search:"type:order;column:update_by;table:orders"`
	CreatedAt string `form:"createdAtOrder"  search:"type:order;column:created_at;table:orders"`
	UpdatedAt string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:orders"`
	DeletedAt string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:orders"`
	Layer     string `form:"layerOrder"  search:"type:order;column:layer;table:orders"`
	Enable    string `form:"enableOrder"  search:"type:order;column:enable;table:orders"`
	Desc      string `form:"descOrder"  search:"type:order;column:desc;table:orders"`
	CId       string `form:"cIdOrder"  search:"type:order;column:c_id;"`
	ShopId    string `form:"shopIdOrder"  search:"type:order;column:shop_id;table:orders"`
	Status    string `form:"statusOrder"  search:"type:order;column:status;table:orders"`
	Money     string `form:"moneyOrder"  search:"type:order;column:money;table:orders"`
	Number    string `form:"numberOrder"  search:"type:order;column:number;table:orders"`
	Delivery  string `form:"deliveryOrder"  search:"type:order;column:delivery;table:orders"`
}

func (m *OrdersGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type OrdersInsertReq struct {
	Id      int    `json:"-" comment:"主键编码"` // 主键编码
	Layer   int    `json:"layer" comment:"排序"`
	Enable  bool   `json:"enable" comment:"开关"`
	Desc    string `json:"desc" comment:"描述信息"`
	ShopId  int    `json:"shop_id" comment:"关联客户"`
	ClassId int    `json:"class_id"`
	Status  int    `json:"status" comment:"配送状态"`
	Number  int    `json:"number" comment:"下单数量"`
	GoodsId int    `json:"goods_id"  comment:"商品ID"`
	//Delivery int `json:"delivery" comment:"配送周期"`
	GoodsSpecs []OrderGoodsSpecs `json:"goods_specs" comment:"商品规格"`
	common.ControlBy
}
type OrderGoodsSpecs struct {
	SpecsId int     `json:"specs_id" comment:"规格ID"`
	Name    string  `json:"name" comment:"产品名称"`
	Spec    string  `json:"spec" comment:"规格"`
	Money   float64 `json:"money" comment:"金额"`
}

func (s *OrdersInsertReq) Generate(model *models.Orders) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
	model.Layer = s.Layer
	model.Enable = s.Enable
	model.Desc = s.Desc
	model.ShopId = s.ShopId
	model.Status = s.Status

	model.Number = s.Number
	//model.Delivery = s.Delivery
}

func (s *OrdersInsertReq) GetId() interface{} {
	return s.Id
}

type ValetOrderReq struct {
	ShopId     int                `json:"shop_id"`
	DeliveryId int                `json:"delivery_id"` //代客下单,只需要获取选择的时间段就行
	Goods      []*valetOrderSpecs `json:"goods"`
	Desc       string             `json:"desc"`
}
type valetOrderSpecs struct {
	ClassId int           `json:"class_id"`
	Specs   []*valetSpecs `json:"specs"`
}
type valetSpecs struct {
	Id      int     `json:"id"`
	Number  int     `json:"number"`
	Money   float64 `json:"money"`
	GoodsId int     `json:"goods_id"`
}
type ToolsOrdersUpdateReq struct {
	Id       int    `uri:"id" comment:"主键编码"` // 主键编码
	Type     int    `json:"type"`
	Status   int    `json:"status"`
	Desc     string `json:"desc"`
	Delivery int    `json:"delivery"`
}
type OrdersUpdateReq struct {
	Id     int     `uri:"id" comment:"主键编码"` // 主键编码
	Layer  int     `json:"layer" comment:"排序"`
	Enable bool    `json:"enable" comment:"开关"`
	Desc   string  `json:"desc" comment:"描述信息"`
	ShopId int     `json:"shop_id" comment:"关联客户"`
	Status int     `json:"status" comment:"配送状态"`
	Money  float64 `json:"money" comment:"金额"`
	Number int     `json:"number" comment:"下单数量"`
	//Delivery int `json:"delivery" comment:"配送周期"`
	common.ControlBy
}

func (s *OrdersUpdateReq) Generate(model *models.Orders) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.UpdateBy = s.UpdateBy // 添加这而，需要记录是被谁更新的
	model.Layer = s.Layer
	model.Enable = s.Enable
	model.Desc = s.Desc
	model.ShopId = s.ShopId
	model.Status = s.Status
	model.Money = s.Money
	model.Number = s.Number
	//model.Delivery = s.Delivery
}

func (s *OrdersUpdateReq) GetId() interface{} {
	return s.Id
}

// OrdersGetReq 功能获取请求参数
type OrdersGetReq struct {
	Id int `uri:"id"`
}

func (s *OrdersGetReq) GetId() interface{} {
	return s.Id
}

// OrdersDeleteReq 功能删除请求参数
type OrdersDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *OrdersDeleteReq) GetId() interface{} {
	return s.Ids
}
