package dto

import (
	"go-admin/app/company/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type GoodsClassGetPageReq struct {
	dto.Pagination `search:"-"`
	Layer          string `form:"layer"  search:"type:exact;column:layer;table:goods_class" comment:"排序"`
	Enable         string `form:"enable"  search:"type:exact;column:enable;table:goods_class" comment:"开关"`
	CId            string `form:"cId"  search:"type:exact;column:c_id;table:goods_class" comment:"大BID"`
	Name           string `form:"name"  search:"type:contains;column:name;table:goods_class" comment:"商品分类名称"`
	GoodsClassOrder
}

type GoodsClassOrder struct {
	Id        string `form:"idOrder"  search:"type:order;column:id;table:goods_class"`
	CreateBy  string `form:"createByOrder"  search:"type:order;column:create_by;table:goods_class"`
	UpdateBy  string `form:"updateByOrder"  search:"type:order;column:update_by;table:goods_class"`
	CreatedAt string `form:"createdAtOrder"  search:"type:order;column:created_at;table:goods_class"`
	UpdatedAt string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:goods_class"`
	DeletedAt string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:goods_class"`
	Layer     string `form:"layerOrder"  search:"type:order;column:layer;table:goods_class"`
	Enable    string `form:"enableOrder"  search:"type:order;column:enable;table:goods_class"`
	Desc      string `form:"descOrder"  search:"type:order;column:desc;table:goods_class"`
	CId       string `form:"cIdOrder"  search:"type:order;column:c_id;table:goods_class"`
	Name      string `form:"nameOrder"  search:"type:order;column:name;table:goods_class"`
}

func (m *GoodsClassGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type GoodsClassInsertReq struct {
	Id     int    `json:"-" comment:"主键编码"` // 主键编码
	Layer  int    `json:"layer" comment:"排序"`
	Enable bool   `json:"enable" comment:"开关"`
	Desc   string `json:"desc" comment:"描述信息"`
	Name   string `json:"name" comment:"商品分类名称" binding:"required"`
	common.ControlBy
}

func (s *GoodsClassInsertReq) Generate(model *models.GoodsClass) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
	model.Layer = s.Layer
	model.Enable = s.Enable
	model.Desc = s.Desc
	model.Name = s.Name
}

func (s *GoodsClassInsertReq) GetId() interface{} {
	return s.Id
}

type GoodsClassUpdateReq struct {
	Id     int    `uri:"id" comment:"主键编码"` // 主键编码
	Layer  int    `json:"layer" comment:"排序"`
	Enable bool   `json:"enable" comment:"开关"`
	Desc   string `json:"desc" comment:"描述信息"`
	Name   string `json:"name" comment:"商品分类名称" binding:"required"`
	common.ControlBy
}

func (s *GoodsClassUpdateReq) Generate(model *models.GoodsClass) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.UpdateBy = s.UpdateBy // 添加这而，需要记录是被谁更新的
	model.Layer = s.Layer
	model.Enable = s.Enable
	model.Desc = s.Desc
	model.Name = s.Name
}

func (s *GoodsClassUpdateReq) GetId() interface{} {
	return s.Id
}

// GoodsClassGetReq 功能获取请求参数
type GoodsClassGetReq struct {
	Id int `uri:"id"`
}

func (s *GoodsClassGetReq) GetId() interface{} {
	return s.Id
}

// GoodsClassDeleteReq 功能删除请求参数
type GoodsClassDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *GoodsClassDeleteReq) GetId() interface{} {
	return s.Ids
}
