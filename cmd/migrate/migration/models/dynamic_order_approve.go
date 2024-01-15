/**
@Author: chaoqun
* @Date: 2024/1/15 18:50
*/
package models
//订单审批

type OrderApproveCnf struct { //是否开启订单审批
	MiniGlobal
	CId            int       `json:"c_id" gorm:"index;comment:公司(大B)ID"`
}
func (OrderApproveCnf) TableName() string {
	return "order_approve_cnf"
}
