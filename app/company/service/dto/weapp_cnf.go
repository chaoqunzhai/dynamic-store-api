/**
@Author: chaoqun
* @Date: 2023/7/20 23:38
*/
package dto

type UpdateLogin struct {

	Enable bool   `json:"enable" comment:"开关"`
	T   int `json:"t" comment:"类型"`
	Val string  `json:"val" comment:"值"`
}
