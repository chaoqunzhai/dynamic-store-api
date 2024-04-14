package router

import (
	"github.com/gin-gonic/gin"
	"go-admin/app/company/apis"
)
func init() {

	routerNoCheckRole = append(routerNoCheckRole, registerNoLoginRouter)
}
func registerNoLoginRouter(v1 *gin.RouterGroup) {
	api := apis.Login{}
	//登录
	l := v1.Group("/login")
	{
		//大B登录
		l.POST("", api.UserLogin)

		//选择大B进入
		l.POST("/checked",api.CompanyChecked)

	}
}