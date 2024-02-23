/**
@Author: chaoqun
* @Date: 2024/2/23 12:37
*/
package router

import (
	"github.com/gin-gonic/gin"
	"go-admin/app/company/apis"
)

func init() {
	routerNoCheckRole = append(routerNoCheckRole, registerNoRouter)

}

func registerNoRouter(v1 *gin.RouterGroup) {
	//手机号密码找回
	forgot := apis.Forgot{}
	r := v1.Group("/company")
	{
		//通过手机号获取验证码 10分钟有效
		r.POST("/forgot/code",forgot.GetPhoneCode)

		//手机号验证码 进行登录请求验证
		r.POST("/forgot/verify_code",forgot.VerifyCode)

		//验证通过 密码修改
		r.POST("/forgot/repass",forgot.Repass)
	}
}