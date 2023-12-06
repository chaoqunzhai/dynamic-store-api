package dto

import (
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type CompanyAdsGetPageReq struct {
	dto.Pagination     `search:"-"`
	Layer string `form:"layer"  search:"type:exact;column:layer;table:company_message" comment:"排序"`
	Enable string `form:"enable"  search:"type:exact;column:enable;table:company_message" comment:"开关"`
	CId string `form:"cId"  search:"type:exact;column:c_id;table:company_message" comment:"大BID"`
	BeginTime      string `form:"beginTime" search:"type:gte;column:created_at;table:company_message" comment:"创建时间"`
	EndTime        string `form:"endTime" search:"type:lte;column:created_at;table:company_message" comment:"创建时间"`


}
func (m *CompanyAdsGetPageReq) GetNeedSearch() interface{} {
	return *m
}


type CompanyAdsInsertReq struct {
	Id int `form:"-" comment:"主键编码"` // 主键编码
	Type int `form:"type"`
	Layer int `form:"layer" comment:"排序"`
	Enable bool `form:"enable" comment:"开关"`
	Desc string `form:"desc" comment:"描述信息"`
	CId int `form:"-" comment:"大BID"`
	LinkName string `form:"link_name"`
	LinkUrl string `form:"link_url"`
	ImageUrl string `form:"image_url"`
	common.ControlBy
}
type CompanyAdsDeleteReq struct {
	Id int `json:"id"`
}



type CompanyAdsUpdateReq struct {
	Type int `form:"type"`
	Layer  int    `form:"layer" comment:"排序"`
	Enable bool   `form:"enable" comment:"开关"`
	Desc   string `form:"desc" comment:"备注信息"`
	LinkName string `form:"link_name"`
	LinkUrl string `form:"link_url"`
	ImageUrl string `form:"image_url"`
	FileClear int    `form:"file_clear" comment:"是否清空照片"`
	BaseFiles string `form:"base_files" comment:"原有图片"`
}
