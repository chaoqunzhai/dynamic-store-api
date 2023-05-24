package dto

import (
     
     
     
     
     
     
     "go-admin/app/shop/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type GoodsVipGetPageReq struct {
	dto.Pagination     `search:"-"`
    Layer string `form:"layer"  search:"type:exact;column:layer;table:goods_vip" comment:"排序"`
    Enable string `form:"enable"  search:"type:exact;column:enable;table:goods_vip" comment:"开关"`
    CId string `form:"cId"  search:"type:exact;column:c_id;table:goods_vip" comment:"大BID"`
    GoodsId string `form:"goodsId"  search:"type:exact;column:goods_id;table:goods_vip" comment:"商品ID"`
    SpecsId string `form:"specsId"  search:"type:exact;column:specs_id;table:goods_vip" comment:"规格ID"`
    GradeId string `form:"gradeId"  search:"type:exact;column:grade_id;table:goods_vip" comment:"VipId"`
    CustomPrice string `form:"customPrice"  search:"type:exact;column:custom_price;table:goods_vip" comment:"自定义价格"`
    GoodsVipOrder
}

type GoodsVipOrder struct {
    Id string `form:"idOrder"  search:"type:order;column:id;table:goods_vip"`
    CreateBy string `form:"createByOrder"  search:"type:order;column:create_by;table:goods_vip"`
    CreatedAt string `form:"createdAtOrder"  search:"type:order;column:created_at;table:goods_vip"`
    DeletedAt string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:goods_vip"`
    Layer string `form:"layerOrder"  search:"type:order;column:layer;table:goods_vip"`
    Enable string `form:"enableOrder"  search:"type:order;column:enable;table:goods_vip"`
    CId string `form:"cIdOrder"  search:"type:order;column:c_id;table:goods_vip"`
    GoodsId string `form:"goodsIdOrder"  search:"type:order;column:goods_id;table:goods_vip"`
    SpecsId string `form:"specsIdOrder"  search:"type:order;column:specs_id;table:goods_vip"`
    GradeId string `form:"gradeIdOrder"  search:"type:order;column:grade_id;table:goods_vip"`
    CustomPrice string `form:"customPriceOrder"  search:"type:order;column:custom_price;table:goods_vip"`
    
}

func (m *GoodsVipGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type GoodsVipInsertReq struct {
    Id int `json:"-" comment:"主键编码"` // 主键编码
    Layer string `json:"layer" comment:"排序"`
    Enable string `json:"enable" comment:"开关"`
    CId string `json:"cId" comment:"大BID"`
    GoodsId string `json:"goodsId" comment:"商品ID"`
    SpecsId string `json:"specsId" comment:"规格ID"`
    GradeId string `json:"gradeId" comment:"VipId"`
    CustomPrice string `json:"customPrice" comment:"自定义价格"`
    common.ControlBy
}

func (s *GoodsVipInsertReq) Generate(model *models.GoodsVip)  {
    if s.Id == 0 {
        model.Model = common.Model{ Id: s.Id }
    }
    model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
    model.Layer = s.Layer
    model.Enable = s.Enable
    model.CId = s.CId
    model.GoodsId = s.GoodsId
    model.SpecsId = s.SpecsId
    model.GradeId = s.GradeId
    model.CustomPrice = s.CustomPrice
}

func (s *GoodsVipInsertReq) GetId() interface{} {
	return s.Id
}

type GoodsVipUpdateReq struct {
    Id int `uri:"id" comment:"主键编码"` // 主键编码
    Layer string `json:"layer" comment:"排序"`
    Enable string `json:"enable" comment:"开关"`
    CId string `json:"cId" comment:"大BID"`
    GoodsId string `json:"goodsId" comment:"商品ID"`
    SpecsId string `json:"specsId" comment:"规格ID"`
    GradeId string `json:"gradeId" comment:"VipId"`
    CustomPrice string `json:"customPrice" comment:"自定义价格"`
    common.ControlBy
}

func (s *GoodsVipUpdateReq) Generate(model *models.GoodsVip)  {
    if s.Id == 0 {
        model.Model = common.Model{ Id: s.Id }
    }
    model.Layer = s.Layer
    model.Enable = s.Enable
    model.CId = s.CId
    model.GoodsId = s.GoodsId
    model.SpecsId = s.SpecsId
    model.GradeId = s.GradeId
    model.CustomPrice = s.CustomPrice
}

func (s *GoodsVipUpdateReq) GetId() interface{} {
	return s.Id
}

// GoodsVipGetReq 功能获取请求参数
type GoodsVipGetReq struct {
     Id int `uri:"id"`
}
func (s *GoodsVipGetReq) GetId() interface{} {
	return s.Id
}

// GoodsVipDeleteReq 功能删除请求参数
type GoodsVipDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *GoodsVipDeleteReq) GetId() interface{} {
	return s.Ids
}
