package router

import (
	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerRouter)
}

func registerRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	//api := apis.Lib{}
	//r := v1.Group("").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	//{
	//
	//
	//}
}
