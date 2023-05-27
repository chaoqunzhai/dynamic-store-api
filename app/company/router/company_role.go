package router

import (
	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"

	"go-admin/app/company/apis"
	"go-admin/common/middleware"
	"go-admin/common/actions"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerCompanyRoleRouter)
}

// registerCompanyRoleRouter
func registerCompanyRoleRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := apis.CompanyRole{}
	r := v1.Group("/company-role").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole()).Use(actions.PermissionCompanyRole())
	{
		r.GET("", api.GetPage)
		r.GET("/:id",  api.Get)
		r.POST("", api.Insert)
		r.PUT("/:id",api.Update)
		r.DELETE("", api.Delete)
	}
}