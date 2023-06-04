package dto

import (
     
     
     
     
     "go-admin/app/company/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
    "strings"
)

type GoodsGetPageReq struct {
	dto.Pagination     `search:"-"`
    Layer string `form:"layer"  search:"type:exact;column:layer;table:goods" comment:"排序"`
    Enable string `form:"enable"  search:"type:exact;column:enable;table:goods" comment:"开关"`
    CId string `form:"cId"  search:"type:exact;column:c_id;table:goods" comment:"大BID"`
    Name string `form:"name"  search:"type:contains;column:name;table:goods" comment:"商品名称"`
    VipSale string `form:"vipSale"  search:"type:exact;column:vip_sale;table:goods" comment:"会员价"`
    GoodsOrder
}

type GoodsOrder struct {
    Id string `form:"idOrder"  search:"type:order;column:id;table:goods"`
    CreateBy string `form:"createByOrder"  search:"type:order;column:create_by;table:goods"`
    UpdateBy string `form:"updateByOrder"  search:"type:order;column:update_by;table:goods"`
    CreatedAt string `form:"createdAtOrder"  search:"type:order;column:created_at;table:goods"`
    UpdatedAt string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:goods"`
    DeletedAt string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:goods"`
    Layer string `form:"layerOrder"  search:"type:order;column:layer;table:goods"`
    Enable string `form:"enableOrder"  search:"type:order;column:enable;table:goods"`
    Desc string `form:"descOrder"  search:"type:order;column:desc;table:goods"`
    CId string `form:"cIdOrder"  search:"type:order;column:c_id;table:goods"`
    Name string `form:"nameOrder"  search:"type:order;column:name;table:goods"`
    Subtitle string `form:"subtitleOrder"  search:"type:order;column:subtitle;table:goods"`
    Quota string `form:"quotaOrder"  search:"type:order;column:quota;table:goods"`
    VipSale string `form:"vipSaleOrder"  search:"type:order;column:vip_sale;table:goods"`
    Code string `form:"codeOrder"  search:"type:order;column:code;table:goods"`
    
}

func (m *GoodsGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type GoodsInsertReq struct {
    Id int `json:"-" comment:"主键编码"` // 主键编码
    Layer int `json:"layer" comment:"排序"`
    Enable bool `json:"enable" comment:"开关"`
    Desc string `json:"desc" comment:"商品详情"`
    Name string `json:"name" comment:"商品名称" binding:"required"`
    Subtitle string `json:"subtitle" comment:"副标题"`
    Image []string `json:"image" comment:"商品图片路径"`
    Quota bool `json:"quota" comment:"是否限购"`
    VipSale bool `json:"vip_sale" comment:"会员价"`
    Code string `json:"code" comment:"条形码"`
    Tag []int `json:"tag" comment:"标签"`
    Class []int `json:"class" comment:"分类"`
    Specs []Specs `json:"specs" comment:"规格"`
    common.ControlBy
}

type Specs struct {
    Name string `json:"name" comment:"规格名称"`
    Price float64 `json:"price" comment:"售价"`
    Layer int `json:"layer"`
    Enable bool `json:"enable"`
    Original float64 `json:"original" comment:"原价"`
    Inventory int `json:"inventory" comment:"库存"`
    Unit string `json:"unit" comment:"单位"`
    Limit int `json:"limit" comment:"起售量"`
    Vip []Vip  `json:"vip" comment:"vip价格设置"`
}
type Vip struct {
    Layer int `json:"layer"`
    Enable bool `json:"enable"`
    Grade int `json:"grade" comment:"登记"`
    Price float64  `json:"price" comment:"售价"`
}

type UpdateSpecs struct {
    Id int `json:"id" `
    Name string `json:"name" comment:"规格名称"`
    Layer int `json:"layer"`
    Enable bool `json:"enable"`
    Price float64 `json:"price" comment:"售价"`
    Original float64 `json:"original" comment:"原价"`
    Inventory int `json:"inventory" comment:"库存"`
    Unit string `json:"unit" comment:"单位"`
    Limit int `json:"limit" comment:"起售量"`
    Vip []UpdateVip  `json:"vip" comment:"vip价格设置"`
}
type UpdateVip struct {
    Id int `json:"id" `
    Layer int `json:"layer"`
    Enable bool `json:"enable"`
    Grade int `json:"grade" comment:"登记"`
    Price float64  `json:"price" comment:"售价"`
}
func (s *GoodsInsertReq) Generate(model *models.Goods)  {
    if s.Id == 0 {
        model.Model = common.Model{ Id: s.Id }
    }
    model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
    model.Layer = s.Layer
    model.Enable = s.Enable
    model.Desc = s.Desc

    model.Name = s.Name
    model.Subtitle = s.Subtitle
    if len(s.Image) > 0 {
        model.Image  = strings.Join(s.Image,",")
    }
    model.Quota = s.Quota
    model.VipSale = s.VipSale
    model.Code = s.Code
}

func (s *GoodsInsertReq) GetId() interface{} {
	return s.Id
}

type GoodsUpdateReq struct {
    Id     int    `uri:"id" comment:""` //
    Layer int `json:"layer" comment:"排序"`
    Enable bool `json:"enable" comment:"开关"`
    Desc string `json:"desc" comment:"商品详情"`
    Name string `json:"name" comment:"商品名称" binding:"required"`
    Subtitle string `json:"subtitle" comment:"副标题"`
    Image []string `json:"image" comment:"商品图片路径"`
    Quota bool `json:"quota" comment:"是否限购"`
    VipSale bool `json:"vip_sale" comment:"会员价"`
    Code string `json:"code" comment:"条形码"`
    Tag []int `json:"tag" comment:"标签"`
    Class []int `json:"class" comment:"分类"`
    Specs []UpdateSpecs `json:"specs" comment:"规格"`
    common.ControlBy
}

func (s *GoodsUpdateReq) Generate(model *models.Goods)  {
    if s.Id == 0 {
        model.Model = common.Model{ Id: s.Id }
    }
    model.UpdateBy = s.UpdateBy // 添加这而，需要记录是被谁更新的
    model.Layer = s.Layer
    model.Enable = s.Enable
    model.Desc = s.Desc
    model.Name = s.Name
    model.Subtitle = s.Subtitle
    if len(s.Image) > 0 {
        model.Image  = strings.Join(s.Image,",")
    }
    model.Quota = s.Quota
    model.VipSale = s.VipSale
    model.Code = s.Code
}

func (s *GoodsUpdateReq) GetId() interface{} {
	return s.Id
}

// GoodsGetReq 功能获取请求参数
type GoodsGetReq struct {
     Id int `uri:"id"`
}
func (s *GoodsGetReq) GetId() interface{} {
	return s.Id
}

// GoodsDeleteReq 功能删除请求参数
type GoodsDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *GoodsDeleteReq) GetId() interface{} {
	return s.Ids
}
