package dto

import (
	"go-admin/app/company/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type GoodsTagGetPageReq struct {
	dto.Pagination `search:"-"`
	Layer          string `form:"layer"  search:"type:exact;column:layer;table:goods_tag" comment:"排序"`
	Enable         string `form:"enable"  search:"type:exact;column:enable;table:goods_tag" comment:"开关"`
	CId            string `form:"cId"  search:"type:exact;column:c_id;table:goods_tag" comment:"大BID"`
	Name           string `form:"name"  search:"type:contains;column:name;table:goods_tag" comment:"商品标签名称"`
	GoodsTagOrder
}

type GoodsTagOrder struct {
	Id        string `form:"idOrder"  search:"type:order;column:id;table:goods_tag"`
	CreateBy  string `form:"createByOrder"  search:"type:order;column:create_by;table:goods_tag"`
	UpdateBy  string `form:"updateByOrder"  search:"type:order;column:update_by;table:goods_tag"`
	CreatedAt string `form:"createdAtOrder"  search:"type:order;column:created_at;table:goods_tag"`
	UpdatedAt string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:goods_tag"`
	DeletedAt string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:goods_tag"`
	Layer     string `form:"layerOrder"  search:"type:order;column:layer;table:goods_tag"`
	Enable    string `form:"enableOrder"  search:"type:order;column:enable;table:goods_tag"`
	Desc      string `form:"descOrder"  search:"type:order;column:desc;table:goods_tag"`
	CId       string `form:"cIdOrder"  search:"type:order;column:c_id;table:goods_tag"`
	Name      string `form:"nameOrder"  search:"type:order;column:name;table:goods_tag"`
	Color     string `form:"colorOrder"  search:"type:order;column:color;table:goods_tag"`
}

func (m *GoodsTagGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type GoodsTagInsertReq struct {
	Id     int    `json:"-" comment:"主键编码"` // 主键编码
	Layer  int    `json:"layer" comment:"排序"`
	Enable bool   `json:"enable" comment:"开关"`
	Desc   string `json:"desc" comment:"描述信息"`
	Name   string `json:"name" comment:"商品标签名称" binding:"required"`
	Color  string `json:"color" comment:"标签颜色"`
	common.ControlBy
}

func (s *GoodsTagInsertReq) Generate(model *models.GoodsTag) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
	model.Layer = s.Layer
	model.Enable = s.Enable
	model.Desc = s.Desc
	model.Name = s.Name
	model.Color = s.Color
}

func (s *GoodsTagInsertReq) GetId() interface{} {
	return s.Id
}

type GoodsTagUpdateReq struct {
	Id     int    `uri:"id" comment:"主键编码"` // 主键编码
	Layer  int    `json:"layer" comment:"排序"`
	Enable bool   `json:"enable" comment:"开关"`
	Desc   string `json:"desc" comment:"描述信息"`
	Name   string `json:"name" binding:"required" comment:"商品标签名称"`
	Color  string `json:"color" comment:"标签颜色"`
	common.ControlBy
}

func (s *GoodsTagUpdateReq) Generate(model *models.GoodsTag) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.UpdateBy = s.UpdateBy // 添加这而，需要记录是被谁更新的
	model.Layer = s.Layer
	model.Enable = s.Enable
	model.Desc = s.Desc
	model.Name = s.Name
	model.Color = s.Color
}

func (s *GoodsTagUpdateReq) GetId() interface{} {
	return s.Id
}

// GoodsTagGetReq 功能获取请求参数
type GoodsTagGetReq struct {
	Id int `uri:"id"`
}

func (s *GoodsTagGetReq) GetId() interface{} {
	return s.Id
}

// GoodsTagDeleteReq 功能删除请求参数
type GoodsTagDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *GoodsTagDeleteReq) GetId() interface{} {
	return s.Ids
}
