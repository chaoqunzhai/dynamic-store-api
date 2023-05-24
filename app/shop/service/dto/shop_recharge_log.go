package dto

import (
     
     
     
     
     
     "time"

	"go-admin/app/shop/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type ShopRechargeLogGetPageReq struct {
	dto.Pagination     `search:"-"`
    ShopId string `form:"shopId"  search:"type:exact;column:shop_id;table:shop_recharge_log" comment:"小BID"`
    Uuid string `form:"uuid"  search:"type:exact;column:uuid;table:shop_recharge_log" comment:"订单号"`
    Source string `form:"source"  search:"type:exact;column:source;table:shop_recharge_log" comment:"充值方式"`
    Money string `form:"money"  search:"type:exact;column:money;table:shop_recharge_log" comment:"支付金额"`
    PayStatus string `form:"payStatus"  search:"type:exact;column:pay_status;table:shop_recharge_log" comment:"支付状态"`
    PayTime time.Time `form:"payTime"  search:"type:contains;column:pay_time;table:shop_recharge_log" comment:"支付时间"`
    ShopRechargeLogOrder
}

type ShopRechargeLogOrder struct {
    Id string `form:"idOrder"  search:"type:order;column:id;table:shop_recharge_log"`
    CreateBy string `form:"createByOrder"  search:"type:order;column:create_by;table:shop_recharge_log"`
    CreatedAt string `form:"createdAtOrder"  search:"type:order;column:created_at;table:shop_recharge_log"`
    DeletedAt string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:shop_recharge_log"`
    ShopId string `form:"shopIdOrder"  search:"type:order;column:shop_id;table:shop_recharge_log"`
    Uuid string `form:"uuidOrder"  search:"type:order;column:uuid;table:shop_recharge_log"`
    Source string `form:"sourceOrder"  search:"type:order;column:source;table:shop_recharge_log"`
    Money string `form:"moneyOrder"  search:"type:order;column:money;table:shop_recharge_log"`
    GiveMoney string `form:"giveMoneyOrder"  search:"type:order;column:give_money;table:shop_recharge_log"`
    PayStatus string `form:"payStatusOrder"  search:"type:order;column:pay_status;table:shop_recharge_log"`
    PayTime string `form:"payTimeOrder"  search:"type:order;column:pay_time;table:shop_recharge_log"`
    
}

func (m *ShopRechargeLogGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type ShopRechargeLogInsertReq struct {
    Id int `json:"-" comment:"主键编码"` // 主键编码
    ShopId string `json:"shopId" comment:"小BID"`
    Uuid string `json:"uuid" comment:"订单号"`
    Source string `json:"source" comment:"充值方式"`
    Money string `json:"money" comment:"支付金额"`
    GiveMoney string `json:"giveMoney" comment:"赠送金额"`
    PayStatus string `json:"payStatus" comment:"支付状态"`
    PayTime time.Time `json:"payTime" comment:"支付时间"`
    common.ControlBy
}

func (s *ShopRechargeLogInsertReq) Generate(model *models.ShopRechargeLog)  {
    if s.Id == 0 {
        model.Model = common.Model{ Id: s.Id }
    }
    model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
    model.ShopId = s.ShopId
    model.Uuid = s.Uuid
    model.Source = s.Source
    model.Money = s.Money
    model.GiveMoney = s.GiveMoney
    model.PayStatus = s.PayStatus
    model.PayTime = s.PayTime
}

func (s *ShopRechargeLogInsertReq) GetId() interface{} {
	return s.Id
}

type ShopRechargeLogUpdateReq struct {
    Id int `uri:"id" comment:"主键编码"` // 主键编码
    ShopId string `json:"shopId" comment:"小BID"`
    Uuid string `json:"uuid" comment:"订单号"`
    Source string `json:"source" comment:"充值方式"`
    Money string `json:"money" comment:"支付金额"`
    GiveMoney string `json:"giveMoney" comment:"赠送金额"`
    PayStatus string `json:"payStatus" comment:"支付状态"`
    PayTime time.Time `json:"payTime" comment:"支付时间"`
    common.ControlBy
}

func (s *ShopRechargeLogUpdateReq) Generate(model *models.ShopRechargeLog)  {
    if s.Id == 0 {
        model.Model = common.Model{ Id: s.Id }
    }
    model.ShopId = s.ShopId
    model.Uuid = s.Uuid
    model.Source = s.Source
    model.Money = s.Money
    model.GiveMoney = s.GiveMoney
    model.PayStatus = s.PayStatus
    model.PayTime = s.PayTime
}

func (s *ShopRechargeLogUpdateReq) GetId() interface{} {
	return s.Id
}

// ShopRechargeLogGetReq 功能获取请求参数
type ShopRechargeLogGetReq struct {
     Id int `uri:"id"`
}
func (s *ShopRechargeLogGetReq) GetId() interface{} {
	return s.Id
}

// ShopRechargeLogDeleteReq 功能删除请求参数
type ShopRechargeLogDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *ShopRechargeLogDeleteReq) GetId() interface{} {
	return s.Ids
}
