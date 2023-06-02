package dto

import (
	"go-admin/app/shop/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type ShopTagGetPageReq struct {
	dto.Pagination `search:"-"`
	Layer          string `form:"layer"  search:"type:exact;column:layer;table:shop_tag" comment:"排序"`
	Enable         string `form:"enable"  search:"type:exact;column:enable;table:shop_tag" comment:"开关"`
	Desc           string `form:"desc"  search:"type:exact;column:desc;table:shop_tag" comment:"描述信息"`
	CId            string `form:"cId"  search:"type:exact;column:c_id;table:shop_tag" comment:"大BID"`
	Name           string `form:"name"  search:"type:exact;column:name;table:shop_tag" comment:"客户标签名称"`
	ShopTagOrder
}

type ShopTagOrder struct {
	Id        string `form:"idOrder"  search:"type:order;column:id;table:shop_tag"`
	CreateBy  string `form:"createByOrder"  search:"type:order;column:create_by;table:shop_tag"`
	UpdateBy  string `form:"updateByOrder"  search:"type:order;column:update_by;table:shop_tag"`
	CreatedAt string `form:"createdAtOrder"  search:"type:order;column:created_at;table:shop_tag"`
	UpdatedAt string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:shop_tag"`
	DeletedAt string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:shop_tag"`
	Layer     string `form:"layerOrder"  search:"type:order;column:layer;table:shop_tag"`
	Enable    string `form:"enableOrder"  search:"type:order;column:enable;table:shop_tag"`
	Desc      string `form:"descOrder"  search:"type:order;column:desc;table:shop_tag"`
	CId       string `form:"cIdOrder"  search:"type:order;column:c_id;table:shop_tag"`
	Name      string `form:"nameOrder"  search:"type:order;column:name;table:shop_tag"`
}

func (m *ShopTagGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type ShopTagInsertReq struct {
	Id     int    `json:"-" comment:"主键编码"` // 主键编码
	Layer  int    `json:"layer" comment:"排序"`
	Enable bool   `json:"enable" comment:"开关"`
	Desc   string `json:"desc" comment:"描述信息"`
	Name   string `json:"name" comment:"客户标签名称" binding:"required"`
	common.ControlBy
}

func (s *ShopTagInsertReq) Generate(model *models.ShopTag) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
	model.Layer = s.Layer
	model.Enable = s.Enable
	model.Desc = s.Desc
	model.Name = s.Name
}

func (s *ShopTagInsertReq) GetId() interface{} {
	return s.Id
}

type ShopTagUpdateReq struct {
	Id     int    `uri:"id" comment:"主键编码"` // 主键编码
	Layer  int    `json:"layer" comment:"排序"`
	Enable bool   `json:"enable" comment:"开关"`
	Desc   string `json:"desc" comment:"描述信息"`
	CId    string `json:"cId" comment:"大BID"`
	Name   string `json:"name" comment:"客户标签名称" binding:"required"`
	common.ControlBy
}

func (s *ShopTagUpdateReq) Generate(model *models.ShopTag) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.UpdateBy = s.UpdateBy // 添加这而，需要记录是被谁更新的
	model.Layer = s.Layer
	model.Enable = s.Enable
	model.Desc = s.Desc
	model.Name = s.Name
}

func (s *ShopTagUpdateReq) GetId() interface{} {
	return s.Id
}

// ShopTagGetReq 功能获取请求参数
type ShopTagGetReq struct {
	Id int `uri:"id"`
}

func (s *ShopTagGetReq) GetId() interface{} {
	return s.Id
}

// ShopTagDeleteReq 功能删除请求参数
type ShopTagDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *ShopTagDeleteReq) GetId() interface{} {
	return s.Ids
}
