package router

import (
	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"

	"go-admin/app/company/apis"
	"go-admin/common/actions"
	"go-admin/common/middleware"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerGoodsSpecsRouter)
}

// registerGoodsSpecsRouter
func registerGoodsSpecsRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := apis.GoodsSpecs{}
	r := v1.Group("/goods-specs").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole()).Use(actions.PermissionCompanyRole())
	{
		r.GET("",  api.GetPage)
		r.GET("/:id",  api.Get)
		r.POST("", api.Insert)
		r.PUT("/:id",  api.Update)
		r.DELETE("", api.Delete)
	}
}