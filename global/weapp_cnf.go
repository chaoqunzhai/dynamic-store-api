/*
*
@Author: chaoqun
* @Date: 2023/7/20 22:54
*/
package global

// 默认配置
const (
	LoginStr    = "username,mobile,wechat"
	RegisterStr = "username,mobile"
)

func LoginCnfToCh(v string) string {
	switch v {
	case "username":
		return "用户名密码"
	case "mobile":
		return "手机号"
	case "wechat":
		return "微信"
	}
	return ""
}
