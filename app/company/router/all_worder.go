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
	routerCheckRole = append(routerCheckRole, registerReportRouter)
}


func registerReportRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := apis.Worker{}

	r := v1.Group("/q/worker").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole()).Use(actions.PermissionCompanyRole())
	{
		//获取导出的记录
		r.GET("/export/list",api.Get)
		//获取下载链接
		r.GET("/download/:uid",api.Download)
		//导出任务记录删除
		r.DELETE("/:uid",api.Remove)

		//创建数据导出任务 通用请求入口，
		//支持选中导出，[OK]
		//汇总导出,[OK]
		//路线导出,
		//
		r.POST("/export/worker",api.Create)

	}



}
