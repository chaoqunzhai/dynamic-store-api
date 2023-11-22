package router

import (
	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"

	"go-admin/app/shop/apis"
	"go-admin/common/actions"
	"go-admin/common/middleware"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerShopRegisterUserRouter)
}

// registerShopBalanceLogRouter
func registerShopRegisterUserRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := apis.ShopRegisterList{}
	r := v1.Group("/shop-register-list").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole()).Use(actions.PermissionCompanyRole())
	{
		r.GET("", api.GetPage)
		r.GET("/:id", api.Detail)
		//审核通过,自动创建用户
		r.PUT("/:id", api.Update)
		r.DELETE("", api.Delete)
	}
}