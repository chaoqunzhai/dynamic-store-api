/**
@Author: chaoqun
* @Date: 2024/1/17 11:15
*/
package router

import (
	"github.com/gin-gonic/gin"

	"go-admin/app/company/apis"

)

func init() {
	routerNoCheckRole = append(routerNoCheckRole, registerCrontabRouter)
}


func registerCrontabRouter(v1 *gin.RouterGroup) {
	api := apis.Crontab{}

	r := v1.Group("/crontab")
	{
		//每天凌晨 进行订单状态同步修改
		//1.没有开启审批: 查询待配送 更新为 配送中
		//2.开启审批: 审批通过 更新为配送中
		r.GET("/sync_order",api.SyncOrder)


	}



}
