package router

import (
	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"

	"go-admin/app/company/apis"
	"go-admin/common/actions"
	"go-admin/common/middleware"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerGoodsRouter)
}

// registerGoodsRouter
func registerGoodsRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := apis.Goods{}
	r := v1.Group("/goods").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole()).Use(actions.PermissionCompanyRole())
	{
		r.GET("", api.GetPage)
		r.GET("/:id", api.Get)
		r.POST("", api.Insert)
		r.POST("/state", api.UpdateState)
		r.PUT("/:id", api.Update)
		r.DELETE("", api.Delete)
		//获取大B所有商品和分类的关联联系
		r.GET("/class_specs", api.ClassSpecs)


	}
}
