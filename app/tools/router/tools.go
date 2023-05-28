package router

import (
	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"

	"go-admin/app/tools/apis"
	"go-admin/common/actions"
	"go-admin/common/middleware"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerCompanyRouter)
}
//todo:系统所有的工具类路由都在这里存放
func registerCompanyRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := apis.Tools{}
	r := v1.Group("/tools").Use(authMiddleware.MiddlewareFunc()).
		Use(middleware.AuthCheckRole()).Use(actions.PermissionCompanyRole())
	{
		r.GET("/image", api.ShowImage)
		r.POST("/image/:t/:name", api.SaveImage)
		//小B注册
		r.POST("/apply/shop")
	}
}
