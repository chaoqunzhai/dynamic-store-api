package dto

import (
     
     
     "go-admin/app/company/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type GoodsBrandGetPageReq struct {
	dto.Pagination     `search:"-"`
    Layer string `form:"layer"  search:"type:exact;column:layer;table:goods_brand" comment:"排序"`
    CId string `form:"cId"  search:"type:exact;column:c_id;table:goods_brand" comment:"大BID"`
    Name string `form:"name"  search:"type:contains;column:name;table:goods_brand" comment:"品牌名称"`
    GoodsBrandOrder
}

type GoodsBrandOrder struct {
    Id string `form:"idOrder"  search:"type:order;column:id;table:goods_brand"`
    Layer string `form:"layerOrder"  search:"type:order;column:layer;table:goods_brand"`
    CId string `form:"cIdOrder"  search:"type:order;column:c_id;table:goods_brand"`
    Name string `form:"nameOrder"  search:"type:order;column:name;table:goods_brand"`
    
}

func (m *GoodsBrandGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type GoodsBrandInsertReq struct {
    Id int `json:"-" comment:"主键编码"` // 主键编码
    Layer int `json:"layer" comment:"排序"`
    CId int `json:"cId" comment:"大BID"`
    Name string `json:"name" comment:"品牌名称"`
    common.ControlBy
}

func (s *GoodsBrandInsertReq) Generate(model *models.GoodsBrand)  {
    if s.Id == 0 {
        model.Model = common.Model{ Id: s.Id }
    }
    model.Layer = s.Layer
    model.CId = s.CId
    model.Name = s.Name
}

func (s *GoodsBrandInsertReq) GetId() interface{} {
	return s.Id
}

type GoodsBrandUpdateReq struct {
    Id int `uri:"id" comment:"主键编码"` // 主键编码
    Layer int `json:"layer" comment:"排序"`
    CId int `json:"cId" comment:"大BID"`
    Name string `json:"name" comment:"品牌名称"`
    common.ControlBy
}

func (s *GoodsBrandUpdateReq) Generate(model *models.GoodsBrand)  {
    if s.Id == 0 {
        model.Model = common.Model{ Id: s.Id }
    }
    model.Layer = s.Layer
    model.CId = s.CId
    model.Name = s.Name
}

func (s *GoodsBrandUpdateReq) GetId() interface{} {
	return s.Id
}

// GoodsBrandGetReq 功能获取请求参数
type GoodsBrandGetReq struct {
     Id int `uri:"id"`
}
func (s *GoodsBrandGetReq) GetId() interface{} {
	return s.Id
}

// GoodsBrandDeleteReq 功能删除请求参数
type GoodsBrandDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *GoodsBrandDeleteReq) GetId() interface{} {
	return s.Ids
}
