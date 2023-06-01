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

		//todo:代客下单
		r.POST("/valet_order", api.ValetOrder)
		//不支持对订单数据的直接更新,因为是客户下单的
		//r.PUT("/:id", api.Update)
		//todo:订单状态更新,周期延后 等
		r.PUT("/tools/:id", api.ToolsOrders)
		//暂时不可进行订单删除
		r.DELETE("", api.Delete)
	}

	r2 := v1.Group("/orders").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		//todo:订单详情
		r2.GET("/:id", api.Get)
		//todo:创建订单
		r2.POST("", api.Insert)
		//todo:校验是否可以下单
		r2.GET("/valid_time", api.ValidTimeConf)

		//todo:获取下单的时间配置
		r2.GET("/times", api.Times)
	}
}
