package dto

import (
     
     
     
     
     
     
     
     "time"

	"go-admin/app/company/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type CompanyGetPageReq struct {
	dto.Pagination     `search:"-"`
    Layer string `form:"layer"  search:"type:exact;column:layer;table:company" comment:"排序"`
    Enable string `form:"enable"  search:"type:exact;column:enable;table:company" comment:"开关"`
    Name string `form:"name"  search:"type:contains;column:name;table:company" comment:"公司(大B)名称"`
    Phone string `form:"phone"  search:"type:contains;column:phone;table:company" comment:"负责人联系手机号"`
    UserName string `form:"userName"  search:"type:exact;column:user_name;table:company" comment:"大B负责人名称"`
    Shop string `form:"shop"  search:"type:exact;column:shop;table:company" comment:"自定义大B系统名称"`
    RenewalTime time.Time `form:"renewalTime"  search:"type:exact;column:renewal_time;table:company" comment:"续费时间"`
    ExpirationTime time.Time `form:"expirationTime"  search:"type:gt;column:expiration_time;table:company" comment:"到期时间"`
    CompanyOrder
}

type CompanyOrder struct {
    Id string `form:"idOrder"  search:"type:order;column:id;table:company"`
    CreateBy string `form:"createByOrder"  search:"type:order;column:create_by;table:company"`
    UpdateBy string `form:"updateByOrder"  search:"type:order;column:update_by;table:company"`
    CreatedAt string `form:"createdAtOrder"  search:"type:order;column:created_at;table:company"`
    UpdatedAt string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:company"`
    DeletedAt string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:company"`
    Layer string `form:"layerOrder"  search:"type:order;column:layer;table:company"`
    Enable string `form:"enableOrder"  search:"type:order;column:enable;table:company"`
    Desc string `form:"descOrder"  search:"type:order;column:desc;table:company"`
    Name string `form:"nameOrder"  search:"type:order;column:name;table:company"`
    Phone string `form:"phoneOrder"  search:"type:order;column:phone;table:company"`
    UserName string `form:"userNameOrder"  search:"type:order;column:user_name;table:company"`
    Shop string `form:"shopOrder"  search:"type:order;column:shop;table:company"`
    Address string `form:"addressOrder"  search:"type:order;column:address;table:company"`
    Longitude string `form:"longitudeOrder"  search:"type:order;column:longitude;table:company"`
    Latitude string `form:"latitudeOrder"  search:"type:order;column:latitude;table:company"`
    Image string `form:"imageOrder"  search:"type:order;column:image;table:company"`
    RenewalTime string `form:"renewalTimeOrder"  search:"type:order;column:renewal_time;table:company"`
    ExpirationTime string `form:"expirationTimeOrder"  search:"type:order;column:expiration_time;table:company"`
    
}

func (m *CompanyGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type CompanyInsertReq struct {
    Id int `json:"-" comment:"主键编码"` // 主键编码
    Layer string `json:"layer" comment:"排序"`
    Enable string `json:"enable" comment:"开关"`
    Desc string `json:"desc" comment:"描述信息"`
    Name string `json:"name" comment:"公司(大B)名称"`
    Phone string `json:"phone" comment:"负责人联系手机号"`
    UserName string `json:"userName" comment:"大B负责人名称"`
    Shop string `json:"shop" comment:"自定义大B系统名称"`
    Address string `json:"address" comment:"大B地址位置"`
    Longitude string `json:"longitude" comment:""`
    Latitude string `json:"latitude" comment:""`
    Image string `json:"image" comment:"logo图片"`
    RenewalTime time.Time `json:"renewalTime" comment:"续费时间"`
    ExpirationTime time.Time `json:"expirationTime" comment:"到期时间"`
    common.ControlBy
}

func (s *CompanyInsertReq) Generate(model *models.Company)  {
    if s.Id == 0 {
        model.Model = common.Model{ Id: s.Id }
    }
    model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
    model.Layer = s.Layer
    model.Enable = s.Enable
    model.Desc = s.Desc
    model.Name = s.Name
    model.Phone = s.Phone
    model.UserName = s.UserName
    model.Shop = s.Shop
    model.Address = s.Address
    model.Longitude = s.Longitude
    model.Latitude = s.Latitude
    model.Image = s.Image
    model.RenewalTime = s.RenewalTime
    model.ExpirationTime = s.ExpirationTime
}

func (s *CompanyInsertReq) GetId() interface{} {
	return s.Id
}

type CompanyUpdateReq struct {
    Id int `uri:"id" comment:"主键编码"` // 主键编码
    Layer string `json:"layer" comment:"排序"`
    Enable string `json:"enable" comment:"开关"`
    Desc string `json:"desc" comment:"描述信息"`
    Name string `json:"name" comment:"公司(大B)名称"`
    Phone string `json:"phone" comment:"负责人联系手机号"`
    UserName string `json:"userName" comment:"大B负责人名称"`
    Shop string `json:"shop" comment:"自定义大B系统名称"`
    Address string `json:"address" comment:"大B地址位置"`
    Longitude string `json:"longitude" comment:""`
    Latitude string `json:"latitude" comment:""`
    Image string `json:"image" comment:"logo图片"`
    RenewalTime time.Time `json:"renewalTime" comment:"续费时间"`
    ExpirationTime time.Time `json:"expirationTime" comment:"到期时间"`
    common.ControlBy
}

func (s *CompanyUpdateReq) Generate(model *models.Company)  {
    if s.Id == 0 {
        model.Model = common.Model{ Id: s.Id }
    }
    model.UpdateBy = s.UpdateBy // 添加这而，需要记录是被谁更新的
    model.Layer = s.Layer
    model.Enable = s.Enable
    model.Desc = s.Desc
    model.Name = s.Name
    model.Phone = s.Phone
    model.UserName = s.UserName
    model.Shop = s.Shop
    model.Address = s.Address
    model.Longitude = s.Longitude
    model.Latitude = s.Latitude
    model.Image = s.Image
    model.RenewalTime = s.RenewalTime
    model.ExpirationTime = s.ExpirationTime
}

func (s *CompanyUpdateReq) GetId() interface{} {
	return s.Id
}

// CompanyGetReq 功能获取请求参数
type CompanyGetReq struct {
     Id int `uri:"id"`
}
func (s *CompanyGetReq) GetId() interface{} {
	return s.Id
}

// CompanyDeleteReq 功能删除请求参数
type CompanyDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *CompanyDeleteReq) GetId() interface{} {
	return s.Ids
}
