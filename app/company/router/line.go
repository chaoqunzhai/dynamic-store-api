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
		//获取一条闲置的路线
		r.GET("unused",api.UnusedOneLine)
		r.GET("",  api.GetPage)
		r.GET("/mini",api.MiniApi)
		r.GET("/:id", api.Get)

		//路线更新,路线名字必须唯一
		r.PUT("/:id",  api.Update)
		//大B路线不能删除了,因为有过期时间
		//r.DELETE("", api.Delete)
		//获取路线下绑定的客户列表
		r.GET("/shop/:id",api.LineBindShopList)
		//更新路线下客户数据
		r.POST("/shop/:id",api.UpdateLineBindShopList)
		//直接绑定客户
		r.POST("/bind_shop",api.BindShop)

	}
}