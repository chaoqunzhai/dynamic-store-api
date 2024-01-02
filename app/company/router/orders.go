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
		//todo:订单列表
		r.GET("", api.GetPage)

		//todo:获取下单时创建的配送周期列表
		r.GET("/cycle_tables",api.Cycle)
		//todo:代客下单
		r.POST("/valet_order", api.ValetOrder)

		//todo:订单状态更新,周期延后 等
		r.PUT("/tools/:id", api.ToolsOrders)

		//todo:获取下订单创建的周期列表和配送列表
		r.GET("/cycle_lists", api.OrderCycleList)

		//todo:获取商家的更多订单记录
		r.GET("/shop/:id", api.ShopOrderList)
	}

	r2 := v1.Group("/orders").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		//todo:订单详情
		r2.GET("/:orderId", api.Get)

		//todo: 多个订单数据
		r2.POST("/rich_data", api.RichData)

		//todo:校验是否可以下单
		r2.GET("/valid_time", api.ValidTimeConf)

		//todo:获取下单的时间配置
		r2.GET("/times", api.Times)

	}

	//订单修改操作
	r3  := v1.Group("/orders").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())

	{
		//todo:对订单进行修改
		r3.POST("/edit/:orderId")

	}
}
