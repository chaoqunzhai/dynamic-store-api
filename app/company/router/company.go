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
		Use(middleware.AuthCheckRole()).Use(actions.PermissionCompanyRole())
	{
		r.GET("/info", api.Info)
		//r.GET("/home", api.MonitorData)
		r.GET("/home", api.Demo)
	}

	//只有超管有权限
	r2 := v1.Group("/company").Use(authMiddleware.MiddlewareFunc()).
		Use(middleware.AuthCheckRole()).Use(actions.PermissionSuperRole())
	{
		//进行续费,大B的续费最好是通过这个接口来统一进行续费
		r2.POST("/renew", api.Renew)
		//续费日志
		r2.GET("/renew", api.RenewPage)
		r2.GET("", api.GetPage)
		r2.GET("/:id", api.Get)
		r2.POST("", api.Insert)
		r2.PUT("/:id", api.Update)
		r2.DELETE("", api.Delete)
	}

}
