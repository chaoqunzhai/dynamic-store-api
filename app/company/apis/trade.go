package apis

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"go-admin/cmd/migrate/migration/models"
	customUser "go-admin/common/jwt/user"
)

type Trade struct {
	api.Api
}
type TradeInsertReq struct {
	CloseHours int `json:"close_hours"`
	ReceiveDays  int `json:"receive_days"`
	RefundDays int `json:"refund_days"`


}
func (e Trade) Create(c *gin.Context) {
	req := TradeInsertReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
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
	fmt.Printf("refund_days:%v",req)
	var OrderTrade models.OrderTrade
	e.Orm.Model(&models.OrderTrade{}).Where("c_id = ?",userDto.CId).Limit(1).Find(&OrderTrade)


	if OrderTrade.Id > 0 {

		OrderTrade.ReceiveDays = req.ReceiveDays
		OrderTrade.RefundDays = req.RefundDays
		OrderTrade.CloseHours= req.CloseHours
		e.Orm.Save(&OrderTrade)
	}else {
		trade:=models.OrderTrade{
			ReceiveDays: req.ReceiveDays,
			RefundDays: req.RefundDays,
			CloseHours: req.CloseHours,
		}
		trade.CId = userDto.CId
		trade.Enable = true
		e.Orm.Create(&trade)
	}
	e.OK("","successful")
	return
}
func (e Trade) Detail(c *gin.Context) {
	err := e.MakeContext(c).
		MakeOrm().
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
	var OrderTrade models.OrderTrade
	e.Orm.Model(&models.OrderTrade{}).Where("c_id = ?",userDto.CId).Limit(1).Find(&OrderTrade)


	e.OK(OrderTrade,"successful")
	return
}
