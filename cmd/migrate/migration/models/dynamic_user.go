/*
*
@Author: chaoqun
* @Date: 2023/8/7 00:37
*/
package models

//用户只要进行更新,发现地址不一样,那就创建新的条目，这样也不会影响旧的历史订单。因为地址是被订单关联的，
//海量的订单去存储收货地址,是非常不现实的，关联地址ID即可，
type DynamicUserAddress struct {
	BigBNoCreateByRichGlobal
	Name      string `json:"name" gorm:"size:20;comment:收件人姓名"`
	Mobile    string `json:"mobile"  gorm:"size:12;comment:收件人电话"`
	Address   string `json:"address"  gorm:"size:100;comment:收件人地址"`
	IsDefault bool   `json:"is_default" gorm:"comment:是否默认地址"`
}

func (DynamicUserAddress) TableName() string {
	return "dynamic_user_address"
}