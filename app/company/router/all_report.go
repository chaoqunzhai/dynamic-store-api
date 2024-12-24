package router

import (
	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"
	"go-admin/common/actions"

	"go-admin/app/company/apis"
	"go-admin/common/middleware"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerReportOrdersRouter)
	routerNoCheckRole = append(routerNoCheckRole, registerNoAuthReportOrdersRouter)
}

// registerOrdersRouter
func registerReportOrdersRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := apis.Orders{}
	r := v1.Group("/report").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole()).Use(actions.PermissionCompanyRole())

	{
		//配送周期 下汇总的商品数
		r.GET("/summary",api.Summary)

		//配送周期下的路线列表
		r.GET("/line",api.Line)
		//配送报表
		r.GET("", api.Index)

		//获取指定配送路线下 小B的列表
		r.GET("/detail/:line_id", api.Detail)
		//获取配送路线下->小B ->商品列表
		r.GET("/detail/line_shop_goods/:line_id",api.DetailShopGoods)

		//点击配送路线下 -> 配送商品信息
		r.GET("/detail/line_goods_detail",api.LineGoodsDetail)

	}
}
func registerNoAuthReportOrdersRouter(v1 *gin.RouterGroup) {

	api := apis.Orders{}
	r := v1.Group("/report")

	{
		r.GET("/file/:path/:xlsx",api.LineSummary)

		r.GET("/:path/:xlsx",api.ExportDownload)
	}


}
