package router

import (
	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"

	"go-admin/app/company/apis"
	"go-admin/common/actions"
	"go-admin/common/middleware"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerAdsRouter)
}

// registerCompanyArticleRouter
func registerAdsRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := apis.CompanyAds{}
	r := v1.Group("/company-ads").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole()).Use(actions.PermissionCompanyRole())
	{
		r.GET("",  api.GetPage)

		r.POST("/enable/:id",api.Enable)

		r.POST("", api.Insert)

		r.PUT("/:id",api.Update)

		r.DELETE("", api.Delete)

	}
}