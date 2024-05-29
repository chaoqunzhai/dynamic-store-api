/**
@Author: chaoqun
* @Date: 2024/5/29 17:56
*/
package dto

import (
	"go-admin/common/dto"
)

type ShopAddressGetPageReq struct {
	dto.Pagination `search:"-"`
	Filter string `form:"filter" search:"-"`
	UserId int `form:"user_id" search:"type:exact;column:user_id;table:dynamic_user_address" comment:"客户"`
}


func (m *ShopAddressGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type ShopAddressInsertReq struct {
	UserId            int     `json:"user_id" comment:"客户用户ID"`
	Phone         string  `json:"phone" comment:"联系手机号"` //小B的手机号
	UserName      string  `json:"username" comment:"小B负责人名称"`
	Address       string  `json:"address" comment:"小B收货地址" `
	FullAddress []int `json:"full_address" comment:"省市区的id"`
}
type ShopAddressSet struct {
	Id            int     `json:"id" comment:"客户用户ID"`
}
type ShopAddressUpdateReq struct {
	Id            int     `uri:"id" comment:"主键编码"`
	Phone         string  `json:"phone" comment:"联系手机号"` //小B的手机号
	UserName      string  `json:"username" comment:"小B负责人名称"`
	Address       string  `json:"address" comment:"小B收货地址" `
	FullAddress []int `json:"full_address" comment:"省市区的id"`
}
type ShopAddressDeleteReq struct {
	Ids []int `json:"ids"`
}
func (s *ShopAddressDeleteReq) GetId() interface{} {
	return s.Ids
}