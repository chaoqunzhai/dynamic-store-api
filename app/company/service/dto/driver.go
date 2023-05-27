package dto

import (
     
     


	"go-admin/app/company/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type DriverGetPageReq struct {
	dto.Pagination     `search:"-"`
    Layer string `form:"layer"  search:"type:exact;column:layer;table:driver" comment:"排序"`
    Enable string `form:"enable"  search:"type:exact;column:enable;table:driver" comment:"开关"`
    CId string `form:"cId"  search:"type:exact;column:c_id;table:driver" comment:"大BID"`
    UserId string `form:"userId"  search:"type:exact;column:user_id;table:driver" comment:"关联的用户ID"`
    Name string `form:"name"  search:"type:contains;column:name;table:driver" comment:"司机名称"`
    Phone string `form:"phone"  search:"type:contains;column:phone;table:driver" comment:"联系手机号"`
    DriverOrder
}

type DriverOrder struct {
    Id string `form:"idOrder"  search:"type:order;column:id;table:driver"`
    CreateBy string `form:"createByOrder"  search:"type:order;column:create_by;table:driver"`
    UpdateBy string `form:"updateByOrder"  search:"type:order;column:update_by;table:driver"`
    CreatedAt string `form:"createdAtOrder"  search:"type:order;column:created_at;table:driver"`
    UpdatedAt string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:driver"`
    DeletedAt string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:driver"`
    Layer string `form:"layerOrder"  search:"type:order;column:layer;table:driver"`
    Enable string `form:"enableOrder"  search:"type:order;column:enable;table:driver"`
    Desc string `form:"descOrder"  search:"type:order;column:desc;table:driver"`
    CId string `form:"cIdOrder"  search:"type:order;column:c_id;table:driver"`
    UserId string `form:"userIdOrder"  search:"type:order;column:user_id;table:driver"`
    Name string `form:"nameOrder"  search:"type:order;column:name;table:driver"`
    Phone string `form:"phoneOrder"  search:"type:order;column:phone;table:driver"`
    
}

func (m *DriverGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type DriverInsertReq struct {
    Id int `json:"-" comment:"主键编码"` // 主键编码
    Layer int `json:"layer" comment:"排序"`
    Enable bool `json:"enable" comment:"开关"`
    Desc string `json:"desc" comment:"备注信息"`
    CId int `json:"cId" comment:"大BID"`
    Name string `json:"name" comment:"司机名称"`
    Phone string `json:"phone" comment:"联系手机号"`
    common.ControlBy
}

func (s *DriverInsertReq) Generate(model *models.Driver)  {
    if s.Id == 0 {
        model.Model = common.Model{ Id: s.Id }
    }
    model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
    model.Layer = s.Layer
    model.Enable = s.Enable
    model.Desc = s.Desc
    model.CId = s.CId
    model.Name = s.Name
    model.Phone = s.Phone
}

func (s *DriverInsertReq) GetId() interface{} {
	return s.Id
}

type DriverUpdateReq struct {
    Id int `uri:"id" comment:"主键编码"` // 主键编码
    Layer int `json:"layer" comment:"排序"`
    Enable bool `json:"enable" comment:"开关"`
    Desc string `json:"desc" comment:"备注信息"`
    CId int `json:"cId" comment:"大BID"`
    Name string `json:"name" comment:"司机名称"`
    Phone string `json:"phone" comment:"联系手机号"`
    common.ControlBy
}

func (s *DriverUpdateReq) Generate(model *models.Driver)  {
    if s.Id == 0 {
        model.Model = common.Model{ Id: s.Id }
    }
    model.UpdateBy = s.UpdateBy // 添加这而，需要记录是被谁更新的
    model.Layer = s.Layer
    model.Enable = s.Enable
    model.Desc = s.Desc
    model.Name = s.Name
    model.Phone = s.Phone
}

func (s *DriverUpdateReq) GetId() interface{} {
	return s.Id
}

// DriverGetReq 功能获取请求参数
type DriverGetReq struct {
     Id int `uri:"id"`
}
func (s *DriverGetReq) GetId() interface{} {
	return s.Id
}

// DriverDeleteReq 功能删除请求参数
type DriverDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *DriverDeleteReq) GetId() interface{} {
	return s.Ids
}
