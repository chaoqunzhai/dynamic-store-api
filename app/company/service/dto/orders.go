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
	ShopId         string `form:"shop_id"  search:"type:exact;column:shop_id;table:orders" comment:"关联客户"`
	Status         string `form:"status"  search:"type:exact;column:status;table:orders" comment:"配送状态"`
	Number         string `form:"number"  search:"type:exact;column:number;table:orders" comment:"下单数量"`
	Delivery       string `form:"delivery"  search:"type:exact;column:delivery;table:orders" comment:"配送周期"`
	BeginTime      string `form:"beginTime" search:"type:gte;column:created_at;table:orders" comment:"创建时间"`
	EndTime        string `form:"endTime" search:"type:lte;column:created_at;table:orders" comment:"创建时间"`
	DeliveryType int `form:"delivery_type" search:"type:exact;column:delivery_type;table:orders" comment:""`
	SourceType int `form:"source_type" search:"type:exact;column:source_type;table:orders" comment:""`
	PayType int `form:"pay_type" search:"type:exact;column:pay_type;table:orders" comment:""`
	Line string `form:"line" search:"type:exact;column:line_id;table:orders" comment:""`
	DeliveryTime string `form:"delivery_time" search:"type:exact;column:delivery_time;table:orders" comment:""`
	Uid string `form:"uid" search:"type:exact;column:uid;table:orders"`
	OrdersOrder
}

type CyClePageReq struct {
	dto.Pagination `search:"-"`


}
func (m *CyClePageReq) GetNeedSearch() interface{} {
	return *m
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


type OrdersShopGetPageReq struct {
	dto.Pagination `search:"-"`

}

func (m *OrdersShopGetPageReq) GetNeedSearch() interface{} {
	return *m
}


type OrdersInsertReq struct {
	Id int `json:"-" comment:"主键编码"` // 主键编码
	//Layer   int    `json:"layer" comment:"排序"`
	//Enable  bool   `json:"enable" comment:"开关"`
	Desc       string            `json:"desc" comment:"描述信息"`
	ShopId     int               `json:"shop_id" comment:"关联客户"`
	Status     int               `json:"status" comment:"配送状态"`
	Number     int               `json:"number" comment:"下单数量"`
	GoodsId    int               `json:"goods_id"  comment:"商品ID"`
	GoodsSpecs []OrderGoodsSpecs `json:"goods_specs" comment:"商品规格"`
	common.ControlBy
}
type OrderGoodsSpecs struct {
	SpecsId   int     `json:"specs_id" comment:"规格ID"`
	Name      string  `json:"name" comment:"产品名称"`
	Unit      string  `json:"unit" comment:"单位`
	Money     float64 `json:"money" comment:"金额"`
	Number    int     `json:"number" comment:"数量"`
	Inventory int     `json:"-" comment:"查询后规格实际的库存"`
}

func (s *OrdersInsertReq) Generate(model *models.Orders) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的

	model.Enable = true

	model.ShopId = s.ShopId
	model.Status = s.Status

	model.Number = s.Number
	//model.Delivery = s.Delivery
}

func (s *OrdersInsertReq) GetId() interface{} {
	return s.Id
}

type RichOrderDataReq struct {
	OrderId []string `json:"order_id"`
}


type OrderCyCleReq struct {
	dto.Pagination `search:"-"`
	CyCle int ` json:"cycle" form:"cycle" search:"-"`
	DeliveryStr           string `json:"delivery_str" form:"delivery_str"  search:"type:contains;column:delivery_str;table:order_cycle_cnf" comment:"商品名称"`
}

func (m *OrderCyCleReq) GetNeedSearch() interface{} {
	return *m
}


type ValetOrderReq struct {
	Shop  int                     `json:"shop"`
	Cycle int                     `json:"cycle"` //代客下单,只需要获取选择的时间段就行
	Goods map[string][]valetSpecs `json:"goods"`
	Desc  string                  `json:"desc"`
	DeductionType int `json:"deduction_type"`
}
type valetSpecs struct {
	Id      int     `json:"id"`
	Number  int     `json:"number"`
	Price   float64 `json:"price"`
	Unit    string  `json:"unit"`
	GoodsId int     `json:"goods_id"`
	GoodsName string `json:"goods_name"`
}
type ToolsOrdersUpdateReq struct {
	Id       int    `uri:"id" comment:"主键编码"` // 主键编码
	Type     int    `json:"type"`
	Status   int    `json:"status"`
	Desc     string `json:"desc"`
	Delivery int    `json:"delivery"`
}
type ShopOrder struct {
	dto.Pagination `search:"-"`
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

	model.Enable = s.Enable

	model.ShopId = s.ShopId
	model.Status = s.Status
	model.OrderMoney = s.Money
	model.Number = s.Number
	//model.Delivery = s.Delivery
}

func (s *OrdersUpdateReq) GetId() interface{} {
	return s.Id
}

// OrdersGetReq 功能获取请求参数
type OrdersGetReq struct {
	OrderId int `uri:"orderId"`
}

func (s *OrdersGetReq) GetId() interface{} {
	return s.OrderId
}

// OrdersDeleteReq 功能删除请求参数
type OrdersDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *OrdersDeleteReq) GetId() interface{} {
	return s.Ids
}

type OrderSpecsRow struct {
	Number int `json:"number"`

}

type OrdersEditReq struct {
	EditList []struct {
		Id       int `json:"id"` //规格ID
		NewAllNumber   int `json:"new_all_number"` //单规格新的总数量
		NewAllMoney float64 `json:"new_all_money"` //单规格新的总价格
	} `json:"editList"`
	Deduction int    `json:"deduction"` //扣款方式|退还方式 1:余额 2:授信额
	Reduce    bool   `json:"reduce"` //减少数量
	Increase  bool   `json:"increase"` //新增数量
	Number    int    `json:"number"` //数量
	Money     float64    `json:"money"` //费用 正数:需要补缴  负数:退还费用
	MoneyStr  string `json:"money_str"` //文字描述
	Desc      string `json:"desc"` //修改描述
}

type OrdersReturnReq struct {
	OrderId string `json:"order_id"` //订单ID
	SpecsId int `json:"specs_id"` //规格ID
}


type OrdersRefundPageReq struct {
	dto.Pagination `search:"-"`
	OrderId          string `form:"order_id"  search:"type:exact;column:order_id;table:order_return" comment:"排序"`
	ReturnId          string `form:"return_id"  search:"type:exact;column:return_id;table:order_return" comment:"排序"`
	ShopId         string `form:"shop_id"  search:"type:exact;column:shop_id;table:order_return" comment:"关联客户"`
	BeginTime      string `form:"beginTime" search:"type:gte;column:created_at;table:order_return" comment:"创建时间"`
	EndTime        string `form:"endTime" search:"type:lte;column:created_at;table:order_return" comment:"创建时间"`
	Status         int `form:"status"  search:"type:exact;column:status;table:order_return" comment:"配送状态"`
}
func (m *OrdersRefundPageReq) GetNeedSearch() interface{} {
	return *m
}

type RefundAuditReq struct {
	RefundId string `json:"refund_id"` //订单ID
	CDesc string `json:"c_desc"`//描述
	Status int `json:"status"` //审批状态
	RefundMoney float64 `json:"refund_money"` //退款金额
	RefundMoneyType int `json:"refund_money_type"` //退款方式

}

type RefundEditReq struct {
	RefundOrderId string `json:"refund_order_id"`
	EditList []EditList `json:"edit_list"`
}
type EditList struct {
	EditNumber int `json:"edit_number"`
	RefundId string `json:"refund_id"`
	SourceNumber int `json:"source_number"`
}