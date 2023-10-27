package router

import (
	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"
	"go-admin/app/admin/apis"

	"go-admin/common/middleware/handler"
	_ "go-admin/docs/admin"
)

func InitSysRouter(r *gin.Engine, authMiddleware *jwt.GinJWTMiddleware) *gin.RouterGroup {
	g := r.Group("")
	sysBaseRouter(g)

	// 需要认证
	sysCheckRoleRouterInit(g, authMiddleware)
	return g
}

func sysBaseRouter(r *gin.RouterGroup) {

	adminBase :=apis.GoAdminSystem{}
	r.GET("", adminBase.GoAdmin)
	r.GET("/info", handler.Ping)
	r.GET("/api/v1/captcha", adminBase.GenerateCaptchaHandler)
	r.GET("/image/:type/:cid/:name", handler.ImageShow)
}



func sysCheckRoleRouterInit(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {

	v1 := r.Group("/api/v1")
	{
		v1.POST("/login", authMiddleware.LoginHandler)

	}
	registerBaseRouter(v1, authMiddleware)
}

func registerBaseRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {

	v2auth := v1.Group("")

	{
		v2auth.POST("/logout", handler.LogOut)
	}
}
