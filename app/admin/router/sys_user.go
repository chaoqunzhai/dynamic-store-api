package router

import (
	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"
	"go-admin/app/admin/apis"
	"go-admin/common/actions"
	"go-admin/common/middleware"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerSysUserRouter)
}

// 需认证的路由代码
func registerSysUserRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := apis.SysUser{}

	v1auth := v1.Group("").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		v1auth.GET("/getinfo", api.GetInfo)
	}
	v2auth := v1.Group("/user").Use(authMiddleware.MiddlewareFunc()).Use(actions.PermissionCompanyRole())
	{
		v2auth.GET("/info", api.GetUserInfo)
	}
}
