package router

import (
	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"
	"go-admin/app/admin/apis"
	"go-admin/common/middleware"
)

func init() {
	CompanyRouterCheckRole = append(CompanyRouterCheckRole, registerSysMenuRouter)
}

// 需认证的路由代码
func registerSysMenuRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := apis.SysMenu{}

	r := v1.Group("/menu").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		r.GET("", api.GetPage)
		//获取自定义大B的菜单列表
		r.GET("/list", api.GetMenuRole)
		//获取mbm管理端的权限列表
		r.GET("/mbm/list",api.MbmList)
		r.GET("/:id", api.Get)
		r.POST("", api.Insert)
		r.PUT("/:id", api.Update)
		r.DELETE("", api.Delete)
	}
	r1 := v1.Group("").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		//只有admin的权限的时候 就能获取admin-ui后所以得菜单
		r1.GET("/menurole", api.GetAdminMenuRole)
		//r1.GET("/menuids", api.GetMenuIDS)
	}
}
