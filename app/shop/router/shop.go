package router

import (
	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"

	"go-admin/app/shop/apis"
	"go-admin/common/actions"
	"go-admin/common/middleware"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerShopRouter)
}

// registerShopRouter
func registerShopRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := apis.Shop{}
	r := v1.Group("/shop").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole()).Use(actions.PermissionCompanyRole())
	{
		r.GET("",  api.GetPage)
		r.GET("/mini",api.MiniApi)
		r.GET("/:id", api.Get)
		r.POST("", api.Insert)
		r.PUT("/:id", api.Update)
		r.DELETE("", api.Delete)
		//积分增加
		r.POST("/integral",api.Integral)
		//授信分增加
		r.POST("/credit",api.Credit)
		//金额增加
		r.POST("/balance",api.Amount)
		//等级修改
		r.POST("/grade",api.Grade)
		//获取客户配置的路线信息
		r.GET("/line/:id",api.GetLine)
	}
}