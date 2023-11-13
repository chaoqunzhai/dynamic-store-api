package dto

import (
    "go-admin/common/dto"
    common "go-admin/common/models"
)

type PayMetOrderGetPageReq struct {
	dto.Pagination     `search:"-"`
    Status string `form:"status"  search:"type:exact;column:status;table:company_apply_payment_order" comment:"状态"`
    UseTo string `form:"use_to"  search:"type:exact;column:use_to;table:company_apply_payment_order" comment:"用途"`
    Layer string `form:"layer"  search:"type:exact;column:layer;table:company_apply_payment_order" comment:"排序"`
    Enable string `form:"enable"  search:"type:exact;column:enable;table:company_apply_payment_order" comment:"开关"`
    CId string `form:"cId"  search:"type:exact;column:c_id;table:company_apply_payment_order" comment:"大BID"`
    Desc string `form:"desc"  search:"type:contains;column:desc;table:company_apply_payment_order" comment:"描述内容"`
    BeginTime      string `form:"beginTime" search:"type:gte;column:created_at;table:company_apply_payment_order" comment:"创建时间"`
    EndTime        string `form:"endTime" search:"type:lte;column:created_at;table:company_apply_payment_order" comment:"创建时间"`
    Before int `form:"before"  search:"-" comment:"状态"`
}


func (m *PayMetOrderGetPageReq) GetNeedSearch() interface{} {
	return *m
}



type PayMetOrderUpdateReq struct {
    Id int `uri:"id" comment:"主键编码"` // 主键编码
    Status int `json:"status"`
    UseTo int `json:"use_to"`
    Desc string `json:"desc" comment:"审批内容"`
    common.ControlBy
}

