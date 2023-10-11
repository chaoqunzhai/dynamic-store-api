package router

import (
	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"

	"go-admin/app/company/apis"
	"go-admin/common/actions"
	"go-admin/common/middleware"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerLineRouter)
}

// registerLineRouter
func registerLineRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := apis.Line{}
	r := v1.Group("/line").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole()).Use(actions.PermissionCompanyRole())
	{
		r.GET("",  api.GetPage)
		r.GET("/mini",api.MiniApi)
		r.GET("/:id", api.Get)
		r.POST("", api.Insert)
		r.PUT("/:id",  api.Update)
		r.DELETE("", api.Delete)
		//获取路线下绑定的客户列表
		r.GET("/shop/:id",api.LineBindShopList)
		//更新路线下客户数据
		r.POST("/shop/:id",api.UpdateLineBindShopList)
		//直接绑定客户
		r.POST("/bind_shop",api.BindShop)

	}
}