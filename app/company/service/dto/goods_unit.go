package dto

import (

	"go-admin/app/company/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type GoodsUnitGetPageReq struct {
	dto.Pagination     `search:"-"`
    CId string `form:"cId"  search:"type:exact;column:c_id;table:goods_unit" comment:"大BID"`
    Name string `form:"name"  search:"type:contains;column:name;table:goods_unit" comment:"单位"`
    GoodsUnitOrder
}

type GoodsUnitOrder struct {
    Id string `form:"idOrder"  search:"type:order;column:id;table:goods_unit"`
    Layer string `form:"layerOrder"  search:"type:order;column:layer;table:goods_unit"`
    CId string `form:"cIdOrder"  search:"type:order;column:c_id;table:goods_unit"`
    Name string `form:"nameOrder"  search:"type:order;column:name;table:goods_unit"`
    
}

func (m *GoodsUnitGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type GoodsUnitInsertReq struct {
    Id int `json:"-" comment:"主键编码"` // 主键编码
    Layer int `json:"layer" comment:"排序"`
    CId int `json:"-" comment:"大BID"`
    Name string `json:"name" comment:"单位"`
    common.ControlBy
}

func (s *GoodsUnitInsertReq) Generate(model *models.GoodsUnit)  {
    if s.Id == 0 {
        model.Model = common.Model{ Id: s.Id }
    }
    model.Layer = s.Layer
    model.CId = s.CId
    model.Name = s.Name
}

func (s *GoodsUnitInsertReq) GetId() interface{} {
	return s.Id
}

type GoodsUnitUpdateReq struct {
    Id int `uri:"id" comment:"主键编码"` // 主键编码
    Layer int `json:"layer" comment:"排序"`
    CId int `json:"-" comment:"大BID"`
    Name string `json:"name" comment:"单位"`
    common.ControlBy
}

func (s *GoodsUnitUpdateReq) Generate(model *models.GoodsUnit)  {
    if s.Id == 0 {
        model.Model = common.Model{ Id: s.Id }
    }
    model.Layer = s.Layer
    model.CId = s.CId
    model.Name = s.Name
}

func (s *GoodsUnitUpdateReq) GetId() interface{} {
	return s.Id
}

// GoodsUnitGetReq 功能获取请求参数
type GoodsUnitGetReq struct {
     Id int `uri:"id"`
}
func (s *GoodsUnitGetReq) GetId() interface{} {
	return s.Id
}

// GoodsUnitDeleteReq 功能删除请求参数
type GoodsUnitDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *GoodsUnitDeleteReq) GetId() interface{} {
	return s.Ids
}
