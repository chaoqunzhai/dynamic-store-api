package router

import (
"github.com/gin-gonic/gin"
jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"

"go-admin/app/shop/apis"
"go-admin/common/actions"
"go-admin/common/middleware"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerShopAddressRouter)
}

// registerShopBalanceLogRouter
func registerShopAddressRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := apis.ShopAddress{}
	r := v1.Group("/shop-address").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole()).Use(actions.PermissionCompanyRole())
	{
		r.GET("", api.GetPage)
		r.POST("", api.Insert)
		r.POST("default", api.Set)
		r.PUT("/:id", api.Update)
		r.DELETE("", api.Delete)
	}
}