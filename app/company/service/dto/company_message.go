package dto

import (
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type CompanyMessageGetPageReq struct {
	dto.Pagination     `search:"-"`
	Layer string `form:"layer"  search:"type:exact;column:layer;table:company_message" comment:"排序"`
	Enable string `form:"enable"  search:"type:exact;column:enable;table:company_message" comment:"开关"`
	CId string `form:"cId"  search:"type:exact;column:c_id;table:company_message" comment:"大BID"`
	BeginTime      string `form:"beginTime" search:"type:gte;column:created_at;table:company_message" comment:"创建时间"`
	EndTime        string `form:"endTime" search:"type:lte;column:created_at;table:company_message" comment:"创建时间"`


}
func (m *CompanyMessageGetPageReq) GetNeedSearch() interface{} {
	return *m
}


type CompanyMessageInsertReq struct {
	Id int `json:"-" comment:"主键编码"` // 主键编码
	Layer int `json:"layer" comment:"排序"`
	Enable bool `json:"enable" comment:"开关"`
	Desc string `json:"desc" comment:"描述信息"`
	CId int `json:"-" comment:"大BID"`
	Link string `json:"class" comment:"跳转链接"`
	Context string `json:"context" comment:"文章内容"`
	common.ControlBy
}
type CompanyMessageDeleteReq struct {
	Ids []int `json:"ids"`
}
func (s *CompanyMessageDeleteReq) GetId() interface{} {
	return s.Ids
}


type CompanyMessageUpdateReq struct {
	Id     int    `uri:"id" comment:"主键编码"` // 主键编码
	Layer  int    `json:"layer" comment:"排序"`
	Enable bool   `json:"enable" comment:"开关"`
	Desc   string `json:"desc" comment:"备注信息"`
	Link string `json:"link" comment:"跳转链接"`
	Context string `json:"context" comment:"文章内容"`
}
