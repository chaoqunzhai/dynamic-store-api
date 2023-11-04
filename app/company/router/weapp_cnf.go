/*
*
@Author: chaoqun
* @Date: 2023/7/20 22:31
*/
package router

import (
	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"
	"go-admin/common/actions"

	"go-admin/app/company/apis"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerWeAppCnfRouter)
}

//对大B小程序配置,那就只有超管有这个权限

// registerOrdersRouter
func registerWeAppCnfRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := apis.WeApp{}
	r := v1.Group("/weappp/conf").Use(authMiddleware.MiddlewareFunc()).Use(actions.PermissionCompanyRole())
	{

		//todo:登录列表,暂时关闭
		r.GET("/login", api.LoginList)
		//todo:修改登录方式,暂时关闭
		r.POST("/login", api.UpdateLoginList)

		//todo:底栏菜单配置
		r.GET("/navbar", api.Navbar)
		//todo:修改底栏菜单配置
		r.POST("/navbar", api.UpdateNavbar)

		//todo:个人中心配置工具集合列表
		r.GET("/quick", api.Quick)

		r.POST("/quick", api.UpdateQuick)
	}

}
