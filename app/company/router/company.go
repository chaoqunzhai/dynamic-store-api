package router

import (
	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"

	"go-admin/app/company/apis"
	"go-admin/common/actions"
	"go-admin/common/middleware"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerCompanyRouter)
}

// registerCompanyRouter
// 大B的信息,需要鉴定权限 必须是大B 和超管
func registerCompanyRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := apis.Company{}
	r := v1.Group("/company").Use(authMiddleware.MiddlewareFunc()).
		Use(middleware.AuthCheckRole()).Use(actions.PermissionCompanyRole()).Use(middleware.SentinelContext())
	{
		r.GET("/info", api.Info)
		r.GET("/home/data", api.MonitorData)
		r.GET("", api.GetPage)
		r.GET("/:id", api.Get)
		r.POST("", api.Insert)
		r.PUT("/:id", api.Update)
		r.DELETE("", api.Delete)
	}
}
