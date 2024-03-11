package router

import (
	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"
	"go-admin/common/actions"

	"go-admin/app/company/apis"
	"go-admin/common/middleware"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerPayCnfRouter)
}

func registerPayCnfRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := apis.PayApi{}
	r := v1.Group("/pay_cnf").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole()).Use(actions.PermissionCompanyRole())
	{
		//todo:更新
		r.POST("", api.Create)
		//todo:获取配置
		r.GET("", api.Detail)
	}
	r1 := v1.Group("/wechat_app_pay").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole()).Use(actions.PermissionCompanyRole())
	apiWechatPay := apis.PayWechat{}
	{
		r1.POST("", apiWechatPay.CreateWechatAppPay)


		r1.GET("",apiWechatPay.WechatAppDetail)
	}
	r2 := v1.Group("/wechat_pay").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole()).Use(actions.PermissionCompanyRole())

	{

		r2.POST("", apiWechatPay.Create)


		r2.GET("",apiWechatPay.Detail)
	}

	r3 := v1.Group("/ali_pay").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole()).Use(actions.PermissionCompanyRole())
	apiAliPay := apis.PayALI{}
	{

		r3.POST("", apiAliPay.Create)


		r3.GET("",apiAliPay.Detail)
	}


}
