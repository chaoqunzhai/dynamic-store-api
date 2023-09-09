package router

import (
	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"
	"go-admin/common/actions"

	"go-admin/app/company/apis"
	"go-admin/common/middleware"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerThemeRouter)
}

func registerThemeRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := apis.Theme{}
	r := v1.Group("/theme").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole()).Use(actions.PermissionCompanyRole())
	{
		//todo:更新
		r.POST("", api.Create)
		//todo:获取配置
		r.GET("", api.Detail)

	}
}
