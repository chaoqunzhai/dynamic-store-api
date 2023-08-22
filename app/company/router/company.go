package router

import (
	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"

	"go-admin/app/company/apis"
	"go-admin/common/actions"
	"go-admin/common/middleware"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerCompanyRouter)
}

// registerCompanyRouter
// 大B的信息,需要鉴定权限 必须是大B 和超管
func registerCompanyRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := apis.Company{}
	r := v1.Group("/company").Use(authMiddleware.MiddlewareFunc()).
		Use(middleware.AuthCheckRole()).Use(actions.PermissionCompanyRole())
	{
		r.GET("/info", api.Info)
		//一些限制配置
		r.GET("/cnf", api.Cnf)
		//r.GET("/home", api.MonitorData)
		r.GET("/home", api.Demo)

		//大B商城模板配置
		r.POST("/category", api.SaveCategory)
		r.GET("/category", api.Category)

		//大B物流信息
		r.GET("/express", api.ExpressList)
		//大B的配送,费用配置,
		r.POST("/express/cnf", api.ExpressCnf)
	}
	//大B用户管理
	{
		//用户列表
		r.GET("/user/list", api.List)
		//更新用户信息
		r.PUT("/user/:id", api.UpdateUser)
		//对用户进行下线,大B看不到了,但是超管还是可以看到的,更新用户的enable
		r.POST("/user/offline", api.Offline)
		//增加用户
		r.POST("/user/add", api.CreateUser)

		r.GET("/user/mini", api.MiniList)
		r.POST("/user/code", api.MakeCode)
	}
	//只有超管有权限
	r2 := v1.Group("/company").Use(authMiddleware.MiddlewareFunc()).
		Use(middleware.AuthCheckRole()).Use(actions.PermissionSuperRole())
	{
		//进行续费,大B的续费最好是通过这个接口来统一进行续费
		r2.POST("/renew", api.Renew)
		//续费日志
		r2.GET("/renew", api.RenewPage)
		r2.GET("", api.GetPage)
		r2.GET("/:id", api.Get)
		r2.POST("", api.Insert)
		r2.PUT("/:id", api.Update)
		r2.DELETE("", api.Delete)
	}

}
