package apis

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"go-admin/cmd/migrate/migration/models"
	"go-admin/common/actions"
	customUser "go-admin/common/jwt/user"
	"go-admin/global"
)

type Trade struct {
	api.Api
}
type TradeInsertReq struct {
	CloseHours int `json:"close_hours" gorm:"size:1;comment:是否开启余额支付"`
	ReceiveDays int `json:"receive_days" gorm:"size:1;comment:是否开启阿里支付"`
	RefundDays int `json:"refund_days" gorm:"size:1;comment:是否开启微信支付"`

}
func (e *Trade) Create(c *gin.Context) {
	req := TradeInsertReq{}
	err := e.MakeContext(c).
		Bind(&req, binding.JSON, nil).
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

	var object models.OrderTrade
	e.Orm.Model(&object).Scopes(actions.PermissionSysUser(object.TableName(), userDto)).Limit(1).Find(&object)


	if object.Id > 0 {

		object.CloseHours = req.CloseHours
		object.RefundDays = req.RefundDays
		object.ReceiveDays = req.ReceiveDays
		e.Orm.Save(&object)
	}else {
		trade:=models.OrderTrade{
			CloseHours: req.CloseHours,
			RefundDays: req.RefundDays,
			ReceiveDays: req.ReceiveDays,
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
	var object models.OrderTrade
	e.Orm.Model(&object).Scopes(actions.PermissionSysUser(object.TableName(), userDto)).Limit(1).Find(&object)

	if object.Id == 0 {
		object = models.OrderTrade{
			CloseHours: int(global.OrderExpirationTime.Minutes()),
			ReceiveDays: global.OrderReceiveDays,
			RefundDays: global.OrderRefundDays,
		}
		object.Enable = true
		object.CId = userDto.CId
		e.Orm.Create(&object)
		e.OK(object,"successful")
		return
	}

	e.OK(object,"successful")
	return
}
