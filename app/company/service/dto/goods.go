package dto

import (
	"go-admin/app/company/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type GoodsGetPageReq struct {
	dto.Pagination `search:"-"`

	ShopId 		 int `form:"shop_id" search:"-"`
	Layer          string `form:"layer"  search:"type:exact;column:layer;table:goods" comment:"排序"`
	Enable         string `form:"enable"  search:"type:exact;column:enable;table:goods" comment:"开关"`
	CId            int `form:"-"  search:"-"`
	Name           string `form:"name"  search:"type:contains;column:name;table:goods" comment:"商品名称"`
	VipSale        string `form:"vipSale"  search:"type:exact;column:vip_sale;table:goods" comment:"会员价"`
	Class          string `form:"class"  search:"-" comment:"分类"`
	Brand 		   string `form:"brand"  search:"-" comment:"品牌"`
	BeginTime      string `form:"beginTime" search:"type:gte;column:created_at;table:goods" comment:"创建时间"`
	EndTime        string `form:"endTime" search:"type:lte;column:created_at;table:goods" comment:"创建时间"`
	GoodsOrder
}

type GoodCountOrder struct {
	Count  int64
	GoodId int 
}
type GoodsRemove struct { 
	Image string `json:"image"`
}
type GoodsOrder struct {
	Id        string `form:"idOrder"  search:"type:order;column:id;table:goods"`
	CreateBy  string `form:"createByOrder"  search:"type:order;column:create_by;table:goods"`
	UpdateBy  string `form:"updateByOrder"  search:"type:order;column:update_by;table:goods"`
	CreatedAt string `form:"createdAtOrder"  search:"type:order;column:created_at;table:goods"`
	UpdatedAt string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:goods"`
	DeletedAt string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:goods"`
	Layer     string `form:"layerOrder"  search:"type:order;column:layer;table:goods"`
	Enable    string `form:"enableOrder"  search:"type:order;column:enable;table:goods"`
	Desc      string `form:"descOrder"  search:"type:order;column:desc;table:goods"`
	CId       string `form:"cIdOrder"  search:"type:order;column:c_id;table:goods"`
	Name      string `form:"nameOrder"  search:"type:order;column:name;table:goods"`
	Subtitle  string `form:"subtitleOrder"  search:"type:order;column:subtitle;table:goods"`
	Quota     string `form:"quotaOrder"  search:"type:order;column:quota;table:goods"`
	VipSale   string `form:"vipSaleOrder"  search:"type:order;column:vip_sale;table:goods"`
	Code      string `form:"codeOrder"  search:"type:order;column:code;table:goods"`
}

func (m *GoodsGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type GoodsStateReq struct {
	Goods  []int `json:"goods" comment:"主键编码"` // 主键编码
	Enable bool  `json:"enable"  comment:"开关"`
}
type GoodsInsertReq struct {
	Id       int    `form:"-" comment:"主键编码"` // 主键编码
	Layer    int    `form:"layer"  comment:"排序"`
	Enable   bool   `form:"enable"  comment:"开关"`
	Brand int `form:"brand" comment:"品牌"`
	SpecName string  `form:"spec_name"  comment:"规格名称"`
	Name     string `form:"name"  comment:"商品名称" binding:"required"`
	Subtitle string `form:"subtitle"  comment:"副标题"`
	Quota    int    `form:"quota"  comment:"活动类型"`
	EnjoyVipSale  bool   `form:"enjoy_vip_sale"  comment:"享受会员功能"`
	VipSale  bool   `form:"vip_sale"  comment:"单独会员价"`
	Code     string `form:"code" comment:"条形码"`
	Tag      string `form:"tag" comment:"标签"`
	Class    string `form:"class"  comment:"分类"`
	Specs    string `form:"specs"  comment:"分类"`
	Recommend bool `form:"recommend" json:"recommend"`
	RubikCube bool  `form:"rubik_cube" json:"rubik_cube"`
	Content string `form:"content"  json:"content"`

	common.ControlBy
}

/*
	{
	           "key": 1686048835986,
	           "name": "2",
	           "price": 2,
	           "original": 3,
	           "inventory": 4,
	           "limit": 100,
	           "enable": true,
	           "layer": 1,
	           "unit": "件",
	           "vip": {
	               "key": 1686048835986,
	               "name": "2",
	               "price": 2,
	               "vip_2": 35,
	               "vip_1": 35
	           }
	       }
*/
type Specs struct {
	Id     int         `json:"id" form:"id" `
	Key    interface{} `json:"key" form:"key"`
	Name   string      `json:"name" form:"name" comment:"规格名称"`
	Market interface{} `form:"market"  json:"market" comment:"市场价"`
	Price  interface{} `json:"price" form:"price" comment:"销售价"`
	Layer  int         `json:"layer" form:"layer"`
	Enable bool        `json:"enable" form:"enable"`
	Code   string      `json:"code" form:"code"`

	VirtuallySale int `json:"virtually_sale" form:"virtually_sale"`

	SerialNumber    string                 `json:"serial_number" form:"number"`
	Image     string                 `json:"image" form:"image"`
	Type      string                 `json:"type"`
	Original  interface{}            `json:"original" form:"original" comment:"原价"`
	Inventory interface{}            `json:"inventory" form:"inventory" comment:"库存"`
	UnitId    interface{}            `json:"unit_id" form:"unit_id" comment:"单位"`
	Limit     interface{}            `json:"limit" form:"limit" comment:"起售量"`
	Max       interface{}            `json:"max" form:"limit" comment:"起售量"`
	Vip       map[string]interface{} `json:"vip" form:"vip" comment:"vip价格设置"`
}

type UpdateSpecs struct {
	Id        int         `form:"id" `
	Name      string      `form:"name" comment:"规格名称"`
	Layer     int         `form:"layer"`
	Enable    bool        `form:"enable"`
	Price     float64     `form:"price" comment:"售价"`
	Original  float64     `form:"original" comment:"原价"`
	Inventory int         `form:"inventory" comment:"库存"`
	Unit      string      `form:"unit" comment:"单位"`
	Limit     int         `form:"limit" comment:"起售量"`
	Vip       []UpdateVip `form:"vip" comment:"vip价格设置"`
}
type UpdateVip struct {
	Id     int     `json:"id" `
	Layer  int     `json:"layer"`
	Enable bool    `json:"enable"`
	Grade  int     `json:"grade" comment:"登记"`
	Price  float64 `json:"price" comment:"售价"`
}

func (s *GoodsInsertReq) Generate(model *models.Goods) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
	model.Layer = s.Layer
	model.Enable = s.Enable

	model.SpecName = s.SpecName
	model.Name = s.Name
	model.Subtitle = s.Subtitle
	model.Quota = s.Quota
	model.VipSale = s.VipSale
	model.EnjoyVipSale = s.EnjoyVipSale
	model.Recommend = s.Recommend
	model.RubikCube = s.RubikCube
}

func (s *GoodsInsertReq) GetId() interface{} {
	return s.Id
}

type GoodsUpdateReq struct {
	Id        int    `uri:"id" comment:""` //
	Layer     int    `form:"layer"  comment:"排序"`
	Enable    bool   `form:"enable"  comment:"开关"`
	Name      string `form:"name"  comment:"商品名称" binding:"required"`
	Subtitle  string `form:"subtitle"  comment:"副标题"`
	Quota     int    `form:"quota"  comment:"活动类型"`
	EnjoyVipSale  bool   `form:"enjoy_vip_sale"  comment:"享受会员功能"`
	VipSale   bool   `form:"vip_sale"  comment:"会员价"`
	Brand int `form:"brand" comment:"品牌"`
	Code      string `form:"code" comment:"条形码"`
	Tag       string `form:"tag" comment:"标签"`
	Class     string `form:"class"  comment:"分类"`
	Specs     string `form:"specs"  comment:"规格"`//是一个json的字符串
	SpecName string  `form:"spec_name"  comment:"规格名称"`
	FileClear int    `form:"file_clear" comment:"是否清空照片"`
	SpecFileClear int    `form:"spec_file_clear" comment:"是否清空规格照片"`
	BaseFiles string `form:"base_files" comment:"原有图片"`
	Recommend bool `form:"recommend" json:"recommend"`
	RubikCube bool  `form:"rubik_cube" json:"rubik_cube"`
	Content string `form:"content"  json:"content"  comment:"商品详情"`
	SpecImageMap string `form:"spec_image_map" comment:"规格图片映射表"`
	common.ControlBy
}

func (s *GoodsUpdateReq) Generate(model *models.Goods) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.UpdateBy = s.UpdateBy // 添加这而，需要记录是被谁更新的
	model.Layer = s.Layer
	model.Enable = s.Enable

	model.Name = s.Name
	model.Subtitle = s.Subtitle
	model.Quota = s.Quota
	model.EnjoyVipSale = s.EnjoyVipSale
	model.VipSale = s.VipSale
	model.SpecName = s.SpecName
	model.Recommend = s.Recommend
	model.RubikCube = s.RubikCube
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
