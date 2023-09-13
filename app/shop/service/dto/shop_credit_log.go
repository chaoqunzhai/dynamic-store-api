package dto

import (
	"go-admin/app/shop/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type ShopCreditLogGetPageReq struct {
	dto.Pagination `search:"-"`
	Id             string `form:"id"  search:"type:exact;column:id;table:shop_credit_log" comment:"订单ID"`
	ShopId         string `form:"shopId"  search:"type:exact;column:shop_id;table:shop_credit_log" comment:"小BID"`
	Money          string `form:"money"  search:"type:exact;column:money;table:shop_credit_log" comment:"变动金额"`
	Scene          string `form:"scene"  search:"type:exact;column:scene;table:shop_credit_log" comment:"变动场景"`
	Type           string `form:"type"  search:"type:exact;column:type;table:shop_credit_log" comment:"变动类型"`
	BeginTime      string `form:"beginTime" search:"type:gte;column:created_at;table:shop_credit_log" comment:"创建时间"`
	EndTime        string `form:"endTime" search:"type:lte;column:created_at;table:shop_credit_log" comment:"创建时间"`
	ShopBalanceLogOrder
}

type ShopCreditLogOrder struct {
	Id        string `form:"idOrder"  search:"type:order;column:id;table:shop_credit_log"`
	CreateBy  string `form:"createByOrder"  search:"type:order;column:create_by;table:shop_credit_log"`
	CreatedAt string `form:"createdAtOrder"  search:"type:order;column:created_at;table:shop_credit_log"`
	DeletedAt string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:shop_credit_log"`
	ShopId    string `form:"shopIdOrder"  search:"type:order;column:shop_id;table:shop_credit_log"`
	Money     string `form:"moneyOrder"  search:"type:order;column:money;table:shop_credit_log"`
	Scene     string `form:"sceneOrder"  search:"type:order;column:scene;table:shop_credit_log"`
	Desc      string `form:"descOrder"  search:"type:order;column:desc;table:shop_credit_log"`
}

func (m *ShopCreditLogGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type ShopCreditLogInsertReq struct {
	Id     int     `json:"-" comment:"主键编码"` // 主键编码
	ShopId int     `json:"shopId" comment:"小BID"`
	Number  int `json:"number" comment:"变动金额"`
	Scene  string  `json:"scene" comment:"变动场景"`
	Desc   string  `json:"desc" comment:"描述/说明"`
	common.ControlBy
}

func (s *ShopCreditLogInsertReq) Generate(model *models.ShopCreditLog) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
	model.ShopId = s.ShopId
	model.Number = s.Number
	model.Scene = s.Scene
	model.Desc = s.Desc
}

func (s *ShopCreditLogInsertReq) GetId() interface{} {
	return s.Id
}

type ShopCreditLogUpdateReq struct {
	Id     int     `uri:"id" comment:"主键编码"` // 主键编码
	ShopId int     `json:"shopId" comment:"小BID"`
	Number  int `json:"number" comment:"变动金额"`
	Scene  string  `json:"scene" comment:"变动场景"`
	Desc   string  `json:"desc" comment:"描述/说明"`
	common.ControlBy
}

func (s *ShopCreditLogUpdateReq) Generate(model *models.ShopCreditLog) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.ShopId = s.ShopId
	model.Number = s.Number
	model.Scene = s.Scene
	model.Desc = s.Desc
}

func (s *ShopCreditLogUpdateReq) GetId() interface{} {
	return s.Id
}

// ShopBalanceLogGetReq 功能获取请求参数
type ShopCreditLogGetReq struct {
	Id int `uri:"id"`
}

func (s *ShopCreditLogGetReq) GetId() interface{} {
	return s.Id
}

// ShopBalanceLogDeleteReq 功能删除请求参数
type ShopCreditLogDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *ShopCreditLogDeleteReq) GetId() interface{} {
	return s.Ids
}
