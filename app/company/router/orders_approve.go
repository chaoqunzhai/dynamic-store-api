/**
@Author: chaoqun
* @Date: 2024/1/2 11:37
*/
package router

import (
	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"
	"go-admin/app/company/apis"
	"go-admin/common/actions"
	"go-admin/common/middleware"
)

//订单审批

func init() {
	routerCheckRole = append(routerCheckRole, registerOrdersApproveRouter)
}

// registerOrdersRouter
func registerOrdersApproveRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := apis.OrdersApprove{}
	r := v1.Group("/orders_approve").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole()).Use(actions.PermissionCompanyRole())
	{
		//获取是否开启配置
		r.GET("",api.Config)
		//进行审批 可以批量审批 驳回(也就是作废)
		r.POST("",api.Approve)
	}
}