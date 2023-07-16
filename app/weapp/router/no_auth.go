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
		//客户端配置
		r.GET("/config/init/:siteId", api.Config)
		//不同菜单 不同的内容
		r.POST("/diyview/info", api.DiyInfo)
		r.POST("/goodssku/components", api.Goodssku)
		r.POST("/goodscategory/tree", api.GoodsTree)
		r.GET("/register/config", api.RegisterCnf)
		r.GET("/config/getCaptchaConfig", api.GetCaptchaConfig)
		//分类列表
		r.POST("/coupon/typelists", api.CouponList)
		r.POST("/captcha", api.Captcha)
		r.POST("/order/num", api.OrderNum)
		//商品列表
		r.POST("/goodssku/page", api.GoodsskuPage)
		//猜你喜欢
		r.POST("/goodssku/recommend", api.GoodsRecommend)
		//商品详情
		r.POST("/goodssku/detail", api.GoodsDetail)
		//
		r.POST("/goods/modifyclicks",api.GoodsModifyclicks)
		//热门搜索
		r.GET("/goods/hotSearchWords",api.HotSearchWords)
		//
		r.POST("/goodsbrand/page",api.GoodsbrandPage)
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
