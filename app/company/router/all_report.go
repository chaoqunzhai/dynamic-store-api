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

		//todo:配送报表
		r.GET("", api.Index)
		//todo:配送报表详情
		r.GET("/:id", api.Detail)



	}
}
