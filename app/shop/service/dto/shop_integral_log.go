package dto

import (
	"go-admin/app/shop/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type ShopIntegralLogGetPageReq struct {
	dto.Pagination `search:"-"`
	Type           string `form:"type"  search:"type:exact;column:type;table:shop_integral_log" comment:"变动类型"`
	ShopId         string `form:"shopId"  search:"type:exact;column:shop_id;table:shop_integral_log" comment:"小BID"`
	BeginTime      string `form:"beginTime" search:"type:gte;column:created_at;table:shop_integral_log" comment:"创建时间"`
	EndTime        string `form:"endTime" search:"type:lte;column:created_at;table:shop_integral_log" comment:"创建时间"`
	ShopIntegralLogOrder
}

type ShopIntegralLogOrder struct {
	Id        string `form:"idOrder"  search:"type:order;column:id;table:shop_integral_log"`
	CreateBy  string `form:"createByOrder"  search:"type:order;column:create_by;table:shop_integral_log"`
	CreatedAt string `form:"createdAtOrder"  search:"type:order;column:created_at;table:shop_integral_log"`
	DeletedAt string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:shop_integral_log"`
	ShopId    string `form:"shopIdOrder"  search:"type:order;column:shop_id;table:shop_integral_log"`
	Number    string `form:"numberOrder"  search:"type:order;column:number;table:shop_integral_log"`
	Scene     string `form:"sceneOrder"  search:"type:order;column:scene;table:shop_integral_log"`
	Desc      string `form:"descOrder"  search:"type:order;column:desc;table:shop_integral_log"`
}

func (m *ShopIntegralLogGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type ShopIntegralLogInsertReq struct {
	Id     int    `json:"-" comment:"主键编码"` // 主键编码
	ShopId int    `json:"shopId" comment:"小BID"`
	Number int    `json:"number" comment:"积分变动数值"`
	Scene  string `json:"scene" comment:"变动场景"`
	Desc   string `json:"desc" comment:"描述/说明"`
	common.ControlBy
}

func (s *ShopIntegralLogInsertReq) Generate(model *models.ShopIntegralLog) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
	model.ShopId = s.ShopId
	model.Number = s.Number
	model.Scene = s.Scene
	model.Desc = s.Desc
}

func (s *ShopIntegralLogInsertReq) GetId() interface{} {
	return s.Id
}

type ShopIntegralLogUpdateReq struct {
	Id     int    `uri:"id" comment:"主键编码"` // 主键编码
	ShopId int    `json:"shopId" comment:"小BID"`
	Number int    `json:"number" comment:"积分变动数值"`
	Scene  string `json:"scene" comment:"变动场景"`
	Desc   string `json:"desc" comment:"描述/说明"`
	common.ControlBy
}

func (s *ShopIntegralLogUpdateReq) Generate(model *models.ShopIntegralLog) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.ShopId = s.ShopId
	model.Number = s.Number
	model.Scene = s.Scene
	model.Desc = s.Desc
}

func (s *ShopIntegralLogUpdateReq) GetId() interface{} {
	return s.Id
}

// ShopIntegralLogGetReq 功能获取请求参数
type ShopIntegralLogGetReq struct {
	Id int `uri:"id"`
}

func (s *ShopIntegralLogGetReq) GetId() interface{} {
	return s.Id
}

// ShopIntegralLogDeleteReq 功能删除请求参数
type ShopIntegralLogDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *ShopIntegralLogDeleteReq) GetId() interface{} {
	return s.Ids
}
