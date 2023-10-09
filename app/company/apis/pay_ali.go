/**
@Author: chaoqun
* @Date: 2023/9/25 23:49
*/
package apis

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"go-admin/cmd/migrate/migration/models"
	"go-admin/common/actions"
	customUser "go-admin/common/jwt/user"
)


type PayALI struct {
	api.Api
}
type AliPayReq struct {
	AppId string `json:"app_id"`
	PrivateKey string `json:"private_key"`
	PublicKey string `json:"public_key"`
	AlipayPublicKey string `json:"alipay_public_key"`
	Enable bool `json:"enable"`
	Refund bool `json:"refund"`
}

func (e *PayALI) Create(c *gin.Context) {
	req := AliPayReq{}
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
	var PayCnf models.AliPay

	e.Orm.Model(&PayCnf).Scopes(actions.PermissionSysUser(PayCnf.TableName(), userDto)).Limit(1).First(&PayCnf)

	if PayCnf.Id > 0 {
		PayCnf.AppId = req.AppId
		PayCnf.Enable = req.Enable
		PayCnf.Refund = req.Refund
		PayCnf.PrivateKey = req.PrivateKey
		PayCnf.PublicKey = req.PublicKey
		PayCnf.AlipayPublicKey = req.AlipayPublicKey
		e.Orm.Save(&PayCnf)
	}else {
		trade:=models.AliPay{
			PublicKey: req.PublicKey,
			AppId: req.AppId,
			PrivateKey: req.PrivateKey,
			AlipayPublicKey: req.AlipayPublicKey,
		}
		trade.CId = userDto.CId
		trade.Enable = req.Enable
		trade.Refund = req.Refund
		e.Orm.Create(&trade)
	}
	e.OK("","successful")

	return
}

func (e *PayALI) Detail(c *gin.Context) {
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

	var data models.AliPay
	e.Orm.Model(&models.AliPay{}).Scopes(
		actions.PermissionSysUser(data.TableName(),userDto)).Limit(1).Find(&data)

	e.OK(data,"successful")
	return
}
