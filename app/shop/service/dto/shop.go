package dto

import (
	"go-admin/app/shop/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type ShopGetPageReq struct {
	dto.Pagination `search:"-"`
	Layer          string `form:"layer"  search:"type:exact;column:layer;table:shop" comment:"排序"`
	Enable         string `form:"enable"  search:"type:exact;column:enable;table:shop" comment:"开关"`
	CId            string `form:"cId"  search:"type:exact;column:c_id;table:shop" comment:"大BID"`
	UserId         string `form:"userId"  search:"type:exact;column:user_id;table:shop" comment:"管理员ID"`
	Name           string `form:"name"  search:"type:exact;column:name;table:shop" comment:"小B名称"`
	Phone          string `form:"phone"  search:"type:contains;column:phone;table:shop" comment:"联系手机号"`
	UserName       string `form:"userName"  search:"type:exact;column:user_name;table:shop" comment:"小B负责人名称"`
	LineId         string `form:"lineId"  search:"type:exact;column:line_id;table:shop" comment:"归属配送路线"`
	ShopOrder
}

type ShopOrder struct {
	Id        string `form:"idOrder"  search:"type:order;column:id;table:shop"`
	CreateBy  string `form:"createByOrder"  search:"type:order;column:create_by;table:shop"`
	UpdateBy  string `form:"updateByOrder"  search:"type:order;column:update_by;table:shop"`
	CreatedAt string `form:"createdAtOrder"  search:"type:order;column:created_at;table:shop"`
	UpdatedAt string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:shop"`
	DeletedAt string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:shop"`
	Layer     string `form:"layerOrder"  search:"type:order;column:layer;table:shop"`
	Enable    string `form:"enableOrder"  search:"type:order;column:enable;table:shop"`
	Desc      string `form:"descOrder"  search:"type:order;column:desc;table:shop"`
	CId       string `form:"cIdOrder"  search:"type:order;column:c_id;table:shop"`
	UserId    string `form:"userIdOrder"  search:"type:order;column:user_id;table:shop"`
	Name      string `form:"nameOrder"  search:"type:order;column:name;table:shop"`
	Phone     string `form:"phoneOrder"  search:"type:order;column:phone;table:shop"`
	UserName  string `form:"userNameOrder"  search:"type:order;column:user_name;table:shop"`
	Address   string `form:"addressOrder"  search:"type:order;column:address;table:shop"`
	Longitude string `form:"longitudeOrder"  search:"type:order;column:longitude;table:shop"`
	Latitude  string `form:"latitudeOrder"  search:"type:order;column:latitude;table:shop"`
	Image     string `form:"imageOrder"  search:"type:order;column:image;table:shop"`
	LineId    string `form:"lineIdOrder"  search:"type:order;column:line_id;table:shop"`
	Amount    string `form:"amountOrder"  search:"type:order;column:amount;table:shop"`
	Integral  string `form:"integralOrder"  search:"type:order;column:integral;table:shop"`
}

func (m *ShopGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type ShopInsertReq struct {
	Id        int     `json:"-" comment:"主键编码"` // 主键编码
	Layer     int     `json:"layer" comment:"排序"`
	Enable    bool    `json:"enable" comment:"开关"`
	Desc      string  `json:"desc" comment:"描述信息"`
	UserId    int     `json:"user_id" comment:"管理员ID"`
	Name      string  `json:"name" comment:"小B名称" binding:"required"`
	Phone     string  `json:"phone" comment:"联系手机号" binding:"required"`
	UserName  string  `json:"username" comment:"小B负责人名称" binding:"required"`
	Address   string  `json:"address" comment:"小B收货地址" `
	Longitude float64     `json:"longitude" comment:""`
	Latitude  float64     `json:"latitude" comment:""`
	Image     string  `json:"image" comment:"图片"`
	LineId    int     `json:"line_id" comment:"归属配送路线"`
	Amount    float64 `json:"amount" comment:"剩余金额"`
	Integral  int     `json:"integral" comment:"可用积分"`
	SalesmanPhone  string     `json:"salesman_phone" comment:"推荐人"`
	Salesman  int     `json:"-" comment:"推荐人"`
	Tags       []int   `json:"tags" comment:"客户标签"`
	common.ControlBy
}

func (s *ShopInsertReq) Generate(model *models.Shop) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
	model.Layer = s.Layer
	model.Enable = s.Enable
	model.Desc = s.Desc
	model.UserId = s.UserId
	model.Name = s.Name
	model.Phone = s.Phone
	model.UserName = s.UserName
	model.Address = s.Address
	model.Longitude = s.Longitude
	model.Latitude = s.Latitude
	model.Image = s.Image
	model.LineId = s.LineId
	model.Amount = s.Amount
	model.Integral = s.Integral
	model.Salesman = s.Salesman
}

func (s *ShopInsertReq) GetId() interface{} {
	return s.Id
}

type ShopIntegralReq struct {
	ShopId int    `json:"shop_id" `
	Number int    `json:"number" `
	Desc   string `json:"desc" `
	Action string `json:"action"`
}
type ShopAmountReq struct {
	ShopId int     `json:"shop_id" `
	Number float64 `json:"number" `
	Desc   string  `json:"desc" `
	Action string  `json:"action"`
}
type ShopUpdateReq struct {
	Id        int     `uri:"id" comment:"主键编码"` // 主键编码
	Layer     int     `json:"layer" comment:"排序"`
	Enable    bool    `json:"enable" comment:"开关"`
	Desc      string  `json:"desc" comment:"描述信息"`
	UserId    int     `json:"userId" comment:"管理员ID"`
	Name      string  `json:"name" comment:"小B名称" binding:"required"`
	Phone     string  `json:"phone" comment:"联系手机号" binding:"required"`
	UserName  string  `json:"username" comment:"小B负责人名称" binding:"required"`
	Address   string  `json:"address" comment:"小B收货地址"`
	Longitude float64     `json:"longitude" comment:""`
	Latitude  float64     `json:"latitude" comment:""`
	Image     string  `json:"image" comment:"图片"`
	LineId    int     `json:"line_id" comment:"归属配送路线"`
	Amount    float64 `json:"amount" comment:"剩余金额"`
	Integral  int     `json:"integral" comment:"可用积分"`
	SalesmanPhone  string     `json:"salesman_phone" comment:"推荐人"`
	Salesman  int     `json:"-" comment:"推荐人"`
	Tags       []int   `json:"tags" comment:"客户标签"`
	common.ControlBy
}

func (s *ShopUpdateReq) Generate(model *models.Shop) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.UpdateBy = s.UpdateBy // 添加这而，需要记录是被谁更新的
	model.Layer = s.Layer
	model.Enable = s.Enable
	model.Desc = s.Desc
	model.UserId = s.UserId
	model.Name = s.Name
	model.Phone = s.Phone
	model.UserName = s.UserName
	model.Address = s.Address
	model.Longitude = s.Longitude
	model.Latitude = s.Latitude
	model.Image = s.Image
	model.LineId = s.LineId
	model.Amount = s.Amount
	model.Integral = s.Integral
	model.Salesman = s.Salesman
}

func (s *ShopUpdateReq) GetId() interface{} {
	return s.Id
}

// ShopGetReq 功能获取请求参数
type ShopGetReq struct {
	Id int `uri:"id"`
}

func (s *ShopGetReq) GetId() interface{} {
	return s.Id
}

// ShopDeleteReq 功能删除请求参数
type ShopDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *ShopDeleteReq) GetId() interface{} {
	return s.Ids
}
