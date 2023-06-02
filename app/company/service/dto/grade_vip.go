package dto

import (
    "go-admin/app/company/models"
    "go-admin/common/dto"
    common "go-admin/common/models"
)

type GradeVipGetPageReq struct {
    dto.Pagination `search:"-"`
    Layer          string `form:"layer"  search:"type:exact;column:layer;table:grade_vip" comment:"排序"`
    Enable         string `form:"enable"  search:"type:exact;column:enable;table:grade_vip" comment:"开关"`
    CId            string `form:"cId"  search:"type:exact;column:c_id;table:grade_vip" comment:"大BID"`
    Name           string `form:"name"  search:"type:contains;column:name;table:grade_vip" comment:"等级名称"`
    Weight         string `form:"weight"  search:"type:exact;column:weight;table:grade_vip" comment:"权重,从小到大"`
    Upgrade        string `form:"upgrade"  search:"type:exact;column:upgrade;table:grade_vip" comment:"升级条件,满多少金额,自动升级Weight+1"`
    GradeVipOrder
}

type GradeVipOrder struct {
    Id        string `form:"idOrder"  search:"type:order;column:id;table:grade_vip"`
    CreateBy  string `form:"createByOrder"  search:"type:order;column:create_by;table:grade_vip"`
    UpdateBy  string `form:"updateByOrder"  search:"type:order;column:update_by;table:grade_vip"`
    CreatedAt string `form:"createdAtOrder"  search:"type:order;column:created_at;table:grade_vip"`
    UpdatedAt string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:grade_vip"`
    DeletedAt string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:grade_vip"`
    Layer     string `form:"layerOrder"  search:"type:order;column:layer;table:grade_vip"`
    Enable    string `form:"enableOrder"  search:"type:order;column:enable;table:grade_vip"`
    Desc      string `form:"descOrder"  search:"type:order;column:desc;table:grade_vip"`
    CId       string `form:"cIdOrder"  search:"type:order;column:c_id;table:grade_vip"`
    Name      string `form:"nameOrder"  search:"type:order;column:name;table:grade_vip"`
    Weight    string `form:"weightOrder"  search:"type:order;column:weight;table:grade_vip"`
    Discount  string `form:"discountOrder"  search:"type:order;column:discount;table:grade_vip"`
    Upgrade   string `form:"upgradeOrder"  search:"type:order;column:upgrade;table:grade_vip"`
}

func (m *GradeVipGetPageReq) GetNeedSearch() interface{} {
    return *m
}

type GradeVipInsertReq struct {
    Id       int     `json:"-" comment:"主键编码"` // 主键编码
    Layer    int     `json:"layer" comment:"排序"`
    Enable   bool    `json:"enable" comment:"开关"`
    Desc     string  `json:"desc" comment:"描述信息"`
    Name     string  `json:"name" comment:"等级名称"  binding:"required"`
    Weight   int     `json:"weight" comment:"权重,从小到大"`
    Discount float64 `json:"discount" comment:"折扣"`
    Upgrade  int     `json:"upgrade" comment:"升级条件,满多少金额,自动升级Weight+1"`
    common.ControlBy
}

func (s *GradeVipInsertReq) Generate(model *models.GradeVip) {
    if s.Id == 0 {
        model.Model = common.Model{Id: s.Id}
    }
    model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
    model.Layer = s.Layer
    model.Enable = s.Enable
    model.Desc = s.Desc
    model.Name = s.Name
    model.Weight = s.Weight
    model.Discount = s.Discount
    model.Upgrade = s.Upgrade
}

func (s *GradeVipInsertReq) GetId() interface{} {
    return s.Id
}

type GradeVipUpdateReq struct {
    Id       int     `uri:"id" comment:"主键编码"` // 主键编码
    Layer    int     `json:"layer" comment:"排序"`
    Enable   bool    `json:"enable" comment:"开关"`
    Desc     string  `json:"desc" comment:"描述信息"`
    Name     string  `json:"name" comment:"等级名称"  binding:"required"`
    Weight   int     `json:"weight" comment:"权重,从小到大"`
    Discount float64 `json:"discount" comment:"折扣"`
    Upgrade  int     `json:"upgrade" comment:"升级条件,满多少金额,自动升级Weight+1"`
    common.ControlBy
}

func (s *GradeVipUpdateReq) Generate(model *models.GradeVip) {
    if s.Id == 0 {
        model.Model = common.Model{Id: s.Id}
    }
    model.UpdateBy = s.UpdateBy // 添加这而，需要记录是被谁更新的
    model.Layer = s.Layer
    model.Enable = s.Enable
    model.Desc = s.Desc
    model.Name = s.Name
    model.Weight = s.Weight
    model.Discount = s.Discount
    model.Upgrade = s.Upgrade
}

func (s *GradeVipUpdateReq) GetId() interface{} {
    return s.Id
}

// GradeVipGetReq 功能获取请求参数
type GradeVipGetReq struct {
    Id int `uri:"id"`
}

func (s *GradeVipGetReq) GetId() interface{} {
    return s.Id
}

// GradeVipDeleteReq 功能删除请求参数
type GradeVipDeleteReq struct {
    Ids []int `json:"ids"`
}

func (s *GradeVipDeleteReq) GetId() interface{} {
    return s.Ids
}
