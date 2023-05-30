package router

import (
	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"
	"go-admin/common/actions"

	"go-admin/app/company/apis"
	"go-admin/common/middleware"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerOrdersRouter)
}

// registerOrdersRouter
func registerOrdersRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := apis.Orders{}
	r := v1.Group("/orders").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole()).Use(actions.PermissionCompanyRole())
	{
		r.GET("", api.GetPage)
		r.PUT("/:id", api.Update)
		r.POST("/valet_order",api.ValetOrder)
		r.DELETE("", api.Delete)
	}

	r2 := v1.Group("/orders").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		r2.GET("/:id", api.Get)
		r2.POST("", api.Insert)
		//todo:校验是否可以下单
		r2.GET("/valid_time",api.ValidTimeConf)

		//todo:获取下单的时间配置
		r2.GET("/times",api.Times)
	}
}
