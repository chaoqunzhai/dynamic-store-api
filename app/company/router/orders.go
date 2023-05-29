package router

import (
	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"

	"go-admin/app/company/apis"
	"go-admin/common/actions"
	"go-admin/common/middleware"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerOrdersRouter)
}

// registerLineRouter
func registerOrdersRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := apis.Orders{}
	r := v1.Group("/orders").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole()).Use(actions.PermissionCompanyRole())
	{
		//订单列表
		r.GET("", api.GetPage)
		//订单详情
		r.GET("/:id", api.Get)
		//代客下单
		r.POST("/valet", api.Valet)
	}
}
