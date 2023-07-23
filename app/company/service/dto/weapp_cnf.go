/*
*
@Author: chaoqun
* @Date: 2023/7/20 23:38
*/
package dto

type UpdateLogin struct {
	Enable bool   `json:"enable" comment:"开关"`
	T      int    `json:"t" comment:"类型"`
	Val    string `json:"val" comment:"值"`
}

type UpdateNav struct {
	NavId  int  `json:"nav_id"` //菜单ID
	CId    int  `json:"c_id"`   //大B
	Enable bool `json:"enable" comment:"开关"`
}

type UpdateQuick struct {
	QuickId  int  `json:"quick_id"` //菜单ID
	CId    int  `json:"c_id"`   //大B
	Enable bool `json:"enable" comment:"开关"`
}
