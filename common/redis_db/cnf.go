/**
@Author: chaoqun
* @Date: 2023/7/20 23:52
*/
package redis_db

type RedisLoginCnf struct {
	ConfigDesc string `json:"config_desc" `
	CreateTime int64 `json:"create_time" `
	Value LoginValue `json:"value" `
	IsUse int `json:"is_use" `

}
type LoginValue struct {
	Login string `json:"login" `
	Register string `json:"register" `
	PwdLen int `json:"pwd_len" `
	PwdComplexity string `json:"pwd_complexity" `
}
