package dto

import (
    "time"
     
     "go-admin/app/company/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type CompanyCouponGetPageReq struct {
	dto.Pagination     `search:"-"`
    Layer string `form:"layer"  search:"type:exact;column:layer;table:company_coupon" comment:"排序"`
    Enable string `form:"enable"  search:"type:exact;column:enable;table:company_coupon" comment:"开关"`
    CId string `form:"cId"  search:"type:exact;column:c_id;table:company_coupon" comment:"大BID"`
    Name string `form:"name"  search:"type:contains;column:name;table:company_coupon" comment:"优惠卷名称"`
    Type string `form:"type"  search:"type:exact;column:type;table:company_coupon" comment:"类型"`
    Range string `form:"range"  search:"type:exact;column:range;table:company_coupon" comment:"使用范围"`
    StartTime time.Time `form:"startTime"  search:"type:gte;column:start_time;table:company_coupon" comment:"开始使用时间"`
    EndTime time.Time `form:"endTime"  search:"type:lte;column:end_time;table:company_coupon" comment:"截止使用时间"`
    Inventory string `form:"inventory"  search:"type:exact;column:inventory;table:company_coupon" comment:"库存"`
    CompanyCouponOrder
}

type CompanyCouponOrder struct {
    Id string `form:"idOrder"  search:"type:order;column:id;table:company_coupon"`
    CreateBy string `form:"createByOrder"  search:"type:order;column:create_by;table:company_coupon"`
    UpdateBy string `form:"updateByOrder"  search:"type:order;column:update_by;table:company_coupon"`
    CreatedAt string `form:"createdAtOrder"  search:"type:order;column:created_at;table:company_coupon"`
    UpdatedAt string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:company_coupon"`
    DeletedAt string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:company_coupon"`
    Layer string `form:"layerOrder"  search:"type:order;column:layer;table:company_coupon"`
    Enable string `form:"enableOrder"  search:"type:order;column:enable;table:company_coupon"`
    Desc string `form:"descOrder"  search:"type:order;column:desc;table:company_coupon"`
    CId string `form:"cIdOrder"  search:"type:order;column:c_id;table:company_coupon"`
    Name string `form:"nameOrder"  search:"type:order;column:name;table:company_coupon"`
    Type string `form:"typeOrder"  search:"type:order;column:type;table:company_coupon"`
    Range string `form:"rangeOrder"  search:"type:order;column:range;table:company_coupon"`
    Money string `form:"moneyOrder"  search:"type:order;column:money;table:company_coupon"`
    Min string `form:"minOrder"  search:"type:order;column:min;table:company_coupon"`
    Max string `form:"maxOrder"  search:"type:order;column:max;table:company_coupon"`
    StartTime string `form:"startTimeOrder"  search:"type:order;column:start_time;table:company_coupon"`
    EndTime string `form:"endTimeOrder"  search:"type:order;column:end_time;table:company_coupon"`
    Inventory string `form:"inventoryOrder"  search:"type:order;column:inventory;table:company_coupon"`
    Limit string `form:"limitOrder"  search:"type:order;column:limit;table:company_coupon"`
    
}

func (m *CompanyCouponGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type CompanyCouponInsertReq struct {
    Id int `json:"-" comment:"主键编码"` // 主键编码
    Layer int `json:"layer" comment:"排序"`
    Enable bool `json:"enable" comment:"开关"`
    Desc string `json:"desc" comment:"描述信息"`
    Name string `json:"name" comment:"优惠卷名称"`
    Type int `json:"type" comment:"类型"`
    Range int `json:"range" comment:"使用范围"`
    Money float64 `json:"money" comment:"优惠卷金额"`
    Min float64 `json:"min" comment:"最低多少钱可以用"`
    Max float64 `json:"max" comment:"满多少钱可以用"`
    StartTime string `json:"start_time" comment:"开始使用时间"`
    EndTime string `json:"end_time" comment:"截止使用时间"`
    Inventory int `json:"inventory" comment:"库存"`
    Limit int `json:"limit" comment:"每个人限领次数"`
    common.ControlBy
}

func (s *CompanyCouponInsertReq) Generate(model *models.CompanyCoupon)  {
    if s.Id == 0 {
        model.Model = common.Model{ Id: s.Id }
    }
    model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
    model.Layer = s.Layer
    model.Enable = s.Enable
    model.Desc = s.Desc

    model.Name = s.Name
    model.Type = s.Type
    model.Range = s.Range
    model.Money = s.Money
    model.Min = s.Min
    model.Max = s.Max

    model.Inventory = s.Inventory
    model.Limit = s.Limit
}

func (s *CompanyCouponInsertReq) GetId() interface{} {
	return s.Id
}

type CompanyCouponUpdateReq struct {
    Id int `uri:"id" comment:"主键编码"` // 主键编码
    Layer int `json:"layer" comment:"排序"`
    Enable bool `json:"enable" comment:"开关"`
    Desc string `json:"desc" comment:"描述信息"`
    Name string `json:"name" comment:"优惠卷名称"`
    Type int `json:"type" comment:"类型"`
    Range int `json:"range" comment:"使用范围"`
    Money float64 `json:"money" comment:"优惠卷金额"`
    Min float64 `json:"min" comment:"最低多少钱可以用"`
    Max float64 `json:"max" comment:"满多少钱可以用"`

    Inventory int `json:"inventory" comment:"库存"`
    Limit int `json:"limit" comment:"每个人限领次数"`
    common.ControlBy
}

func (s *CompanyCouponUpdateReq) Generate(model *models.CompanyCoupon)  {
    if s.Id == 0 {
        model.Model = common.Model{ Id: s.Id }
    }
    model.UpdateBy = s.UpdateBy // 添加这而，需要记录是被谁更新的
    model.Layer = s.Layer
    model.Enable = s.Enable
    model.Desc = s.Desc
    model.Name = s.Name
    model.Type = s.Type
    model.Range = s.Range
    model.Money = s.Money
    model.Min = s.Min
    model.Max = s.Max
    model.Inventory = s.Inventory
    model.Limit = s.Limit
}

func (s *CompanyCouponUpdateReq) GetId() interface{} {
	return s.Id
}

// CompanyCouponGetReq 功能获取请求参数
type CompanyCouponGetReq struct {
     Id int `uri:"id"`
}
func (s *CompanyCouponGetReq) GetId() interface{} {
	return s.Id
}

// CompanyCouponDeleteReq 功能删除请求参数
type CompanyCouponDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *CompanyCouponDeleteReq) GetId() interface{} {
	return s.Ids
}
