/**
@Author: chaoqun
* @Date: 2024/1/7 10:37
*/
package router

import (
	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"
	"go-admin/app/company/apis"
	"go-admin/common/actions"
	"go-admin/common/middleware"
)

//退货API


func init() {
	routerCheckRole = append(routerCheckRole, registerRefundOrdersRouter)
}

// registerOrdersRouter
func registerRefundOrdersRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := apis.OrdersRefund{}
	r := v1.Group("/orders/refund").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole()).Use(actions.PermissionCompanyRole())

	{
		//todo:订单列表
		r.GET("", api.GetPage)
	}

}