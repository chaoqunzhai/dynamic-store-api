package router

import (
	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"

	"go-admin/app/company/apis"
	"go-admin/common/actions"
	"go-admin/common/middleware"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerPayMetOrderRouter)
}

// registerCompanyArticleRouter
func registerPayMetOrderRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := apis.PaymentOrder{}
	r := v1.Group("/company-payment-order").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole()).Use(actions.PermissionCompanyRole())
	{
		r.GET("", actions.PermissionAction(), api.GetPage)
		r.POST("/:id", actions.PermissionAction(), api.Update)
	}
}