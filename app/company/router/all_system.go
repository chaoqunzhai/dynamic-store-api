/**
@Author: chaoqun
* @Date: 2023/12/22 10:42
*/
package router

import (
	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"
	"go-admin/app/company/apis"
	"go-admin/common/actions"
	"go-admin/common/middleware"
)

//存放系统配置
func init() {
	routerCheckRole = append(routerCheckRole, registerSystemRouter)
}

// registerLineRouter
func registerSystemRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := apis.Worker{}
	r := v1.Group("/q/order").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole()).Use(actions.PermissionCompanyRole())
	{
		//数据导出队列
		r.GET("/export/worker",api.Get)
		//获取下载链接
		r.GET("/download/:uid",api.Download)
		//创建数据导出任务
		r.POST("/export/worker",api.Create)
		r.DELETE("/:uid",api.Remove)
	}
}
