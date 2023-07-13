package router

import (
	"github.com/gin-gonic/gin"
	"go-admin/app/weapp/apis"
)

func init() {
	routerNoCheckRole = append(routerNoCheckRole, registerNoRouter)
	routerNoPubCheckRole = append(routerNoPubCheckRole, registerPubNoRouter)
}

func registerNoRouter(v1 *gin.RouterGroup) {
	api := apis.Lib{}
	r := v1.Group("")
	{
		r.GET("/config/init", api.Config)
		r.POST("/diyview/info", api.DiyInfo)
		r.POST("/goodssku/components", api.Goodssku)
		r.POST("/goodscategory/tree", api.GoodsTree)
		r.GET("/register/config", api.RegisterCnf)
		r.GET("/config/getCaptchaConfig", api.GetCaptchaConfig)
		r.POST("/coupon/typelists", api.CouponList)

	}
}

func registerPubNoRouter(v1 *gin.RouterGroup) {
	api := apis.Lib{}
	r := v1.Group("")
	{
		//不同的商户图片不一样
		r.GET("/addon/:shopName/:path", api.ShopImage)

	}
}
