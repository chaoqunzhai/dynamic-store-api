/**
@Author: chaoqun
* @Date: 2024/8/4 22:27
*/
package router

import (
	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"
	"go-admin/app/company/apis"
	"go-admin/common/actions"
	"go-admin/common/middleware"
)

// 精简API
func init() {

	routerCheckRole = append(routerCheckRole, registerMiniRouter)
}



func registerMiniRouter(v1 *gin.RouterGroup,authMiddleware *jwt.GinJWTMiddleware) {
	api := apis.MiniApi{}

	r := v1.Group("/mini").Use(authMiddleware.MiddlewareFunc()).
		Use(middleware.AuthCheckRole()).Use(actions.PermissionCompanyRole())

	{

		r.GET("/goods",api.GoodsSpec)

		r.GET("/customer",api.CustomerUser)

	}

}