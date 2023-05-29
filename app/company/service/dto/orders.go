package dto

import "go-admin/common/dto"

type OrdersGetPageReq struct {
	dto.Pagination `search:"-"`
	Layer          string `form:"layer"  search:"type:exact;column:layer;table:line" comment:"排序"`
	Enable         string `form:"enable"  search:"type:exact;column:enable;table:line" comment:"开关"`
	CId            string `form:"cId"  search:"type:exact;column:c_id;table:line" comment:"大BID"`
	Name           string `form:"name"  search:"type:contains;column:name;table:line" comment:"路线名称"`
	DriverId       string `form:"driverId"  search:"type:exact;column:driver_id;table:line" comment:"关联司机"`
	LineOrder
}

func (m *OrdersGetPageReq) GetNeedSearch() interface{} {
	return *m
}
