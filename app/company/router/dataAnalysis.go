/**
@Author: chaoqun
* @Date: 2024/8/4 17:08
*/
package router

import (
	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"
	"go-admin/app/company/apis"
	"go-admin/common/actions"
	"go-admin/common/middleware"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerDataRouter)
}

// registerGoodsRouter
func registerDataRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := apis.DataAnalysis{}
	r := v1.Group("/data-analysis").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole()).Use(actions.PermissionCompanyRole())
	{

		r.GET("/goodsList", api.GoodsList)

		r.GET("/goodsClass", api.GoodsClassList)

		r.GET("/goodsBrand", api.GoodsBrandList)

		r.GET("/grossCost", api.Grosslist)


	}



}
