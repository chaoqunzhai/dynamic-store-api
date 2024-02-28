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
	e.OK(result,"操作成功")

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
	//无论是否开启审核 都把这个字段进行一个更新,防止后期开启了审核 是在校验approve_status字段
	if req.Action == 1{ //审核通过
		updateMap["approve_status"] = global.OrderApproveOk
		//updateMap["status"] = global.OrderStatusReturn
		//审批操作时 只能是更新默认状态的订单,不会动订单的状态
		e.Orm.Table(splitTableRes.OrderTable).Where("c_id = ? and order_id in ? ",userDto.CId,req.OrderList).Updates(updateMap)

	}else if req.Action == 0  {
		
		
		OrderSpecId:=make([]int,0)
		//这是操作驳回,也就是作废了 需要退库 退钱
		for _,orderId:=range req.OrderList{

			if cancelErr :=s.CancelOrder(global.InventoryApproveIn,true,orderId,OrderSpecId,req.Desc,splitTableRes,userDto);cancelErr!=nil{
				continue
			}
		}
	}else {
		e.Error(500, nil,"非法操作")
		return
	}


	e.OK("","操作成功")
	return

}
