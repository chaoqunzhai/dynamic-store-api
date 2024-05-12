package router

import (
	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"

	"go-admin/app/company/apis"
	"go-admin/common/actions"
	"go-admin/common/middleware"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerGoodsRouter)
}

// registerGoodsRouter
func registerGoodsRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := apis.Goods{}
	r := v1.Group("/goods").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole()).Use(actions.PermissionCompanyRole())
	{
		r.GET("", api.GetPage)
		r.GET("/mini",api.MiniApi)
		r.GET("/:id", api.Get)
		r.POST("", api.Insert)
		r.POST("/state", api.UpdateState)

		r.PUT("/:id", api.Update)
		//商品的删除 会校验一次库存,如果库存>0 不可删除
		r.DELETE("", api.Delete)
		//根据分类获取分类下的商品
		r.GET("/class_specs", api.ClassSpecs)
		//商品的图片上传
		r.POST("/image",api.CosSaveImage)
		//商品图片的删除
		r.DELETE("/image",api.CosRemoveImage)

		//商品数据导出
		r.GET("/export",api.Export)

	}



}
