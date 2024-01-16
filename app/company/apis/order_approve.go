/**
@Author: chaoqun
* @Date: 2024/1/16 11:10
*/
package apis

import (
	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"go-admin/app/company/service"
	"go-admin/app/company/service/dto"
	"go-admin/common/business"
	customUser "go-admin/common/jwt/user"
	"go-admin/global"
)
type OrdersApprove struct {
	api.Api
}
//是否开启了审核 + 当前登录用户是否拥有审批权限
func (e OrdersApprove)Config(c *gin.Context) {
	req := dto.OrdersRefundPageReq{}
	s := service.Orders{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	userDto, err := customUser.GetUserDto(e.Orm, c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}

	openApprove,hasApprove:=service.IsHasOpenApprove(userDto,e.Orm)
	result:=map[string]bool{
		"openApprove":openApprove,
		"hasApprove":hasApprove,
	}
	e.OK(result,"successful")

	return
}

func (e OrdersApprove)Approve(c *gin.Context) {
	req := dto.ApproveReq{}
	s := service.Orders{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	userDto, err := customUser.GetUserDto(e.Orm, c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}

	//检测是否有权限
	openApprove,hasApprove:=service.IsHasOpenApprove(userDto,e.Orm)
	if !openApprove || !hasApprove {
		e.Error(500, nil,"无权限操作")
		return

	}
	splitTableRes := business.GetTableName(userDto.CId, e.Orm)

	updateMap:=map[string]interface{}{
		"approve_msg":req.Desc,
	}
	if req.Action == 1{
		updateMap["approve_status"] = global.OrderApproveOk
	}else {
		//驳回了
		updateMap["approve_status"] = global.OrderApproveReject
	}
	//审批操作时 只能是更新/待发货的订单
	e.Orm.Table(splitTableRes.OrderTable).Where("c_id = ? and order_id in ? and status = ?",userDto.CId,req.OrderList,global.OrderStatusWaitSend).Updates(updateMap)


	if req.Action == 0 {
		OrderSpecId:=make([]int,0)
		//这是操作驳回,也就是作废了 需要退库 退钱
		for _,orderId:=range req.OrderList{

			if cancelErr :=s.CancelOrder(global.InventoryApproveIn,true,orderId,OrderSpecId,req.Desc,splitTableRes,userDto);cancelErr!=nil{
				e.Error(500, cancelErr, cancelErr.Error())
				return
			}
		}
	}

	e.OK("","successful")
	return

}
