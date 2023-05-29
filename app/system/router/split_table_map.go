package router

import (
	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"

	"go-admin/app/system/apis"
	"go-admin/common/actions"
	"go-admin/common/middleware"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerSplitTableMapRouter)
}

// registerSplitTableMapRouter
func registerSplitTableMapRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := apis.SplitTableMap{}
	r := v1.Group("/split-table-map").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole()).Use(actions.PermissionSuperRole())
	{
		r.GET("", api.GetPage)
		r.GET("/:id", api.Get)
		r.POST("", api.Insert)
		r.PUT("/:id", api.Update)
		r.DELETE("", api.Delete)
	}
}
