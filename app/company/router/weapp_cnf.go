/**
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

// registerOrdersRouter
func registerWeAppCnfRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := apis.WeApp{}
	r := v1.Group("/weappp/conf").Use(authMiddleware.MiddlewareFunc()).Use(actions.PermissionCompanyRole())
	{
		//todo:登录列表
		r.GET("/login", api.LoginList)
		//todo:修改登录方式
		r.POST("/login", api.UpdateLoginList)

		//todo:底栏菜单配置
		//r.GET("/navbar", api.LoginList)
		//todo:修改底栏菜单配置
		//r.POST("/navbar", api.LoginList)

		//todo:个人中心配置工具集合列表

	}

}

