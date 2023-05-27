package router

import (
	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"

	"go-admin/app/shop/apis"
	"go-admin/common/middleware"
	"go-admin/common/actions"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerShopIntegralLogRouter)
}

// registerShopIntegralLogRouter
func registerShopIntegralLogRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := apis.ShopIntegralLog{}
	r := v1.Group("/shop-integral-log").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole()).Use(actions.PermissionCompanyRole())
	{
		r.GET("", api.GetPage)
		r.GET("/:id", api.Get)
		//不能主动去创建记录
		//r.POST("", api.Insert)
		r.PUT("/:id",  api.Update)
		r.DELETE("", api.Delete)
	}
}