package router

import (
	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"

	"go-admin/app/shop/apis"
	"go-admin/common/actions"
	"go-admin/common/middleware"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerShopCreditLogRouter)
}

// registerShopBalanceLogRouter
func registerShopCreditLogRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := apis.ShopCreditLog{}
	r := v1.Group("/shop-credit-log").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole()).Use(actions.PermissionCompanyRole())
	{
		r.GET("", api.GetPage)
		r.GET("/:id", api.Get)
		//不能主动去创建记录
		//r.POST("", api.Insert)
		//不能更新
		//r.PUT("/:id", api.Update)
		r.DELETE("", api.Delete)
	}
}