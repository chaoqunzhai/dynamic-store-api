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
	customUser "go-admin/common/jwt/user"
)


type PayWechat struct {
	api.Api
}
type WechatPayReq struct {
	MchId string  `json:"mch_id" `
	ApiV2 string `json:"api_v2"`
	ApiV3 string `json:"api_v3"`
	CertPath string `json:"cert_path" `
	KeyPath string  `json:"key_path" `
	Enable bool `json:"enable"`
	Refund bool `json:"refund"`
}
func (e *PayWechat) Create(c *gin.Context) {
	req := WechatPayReq{}
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
	var PayCnf models.WeChatPay

	e.Orm.Model(&PayCnf).Where("c_id = ? ",userDto.CId).Limit(1).First(&PayCnf)

	if PayCnf.Id > 0 {

		PayCnf.Enable = req.Enable
		PayCnf.Refund = req.Refund
		PayCnf.ApiV2 = req.ApiV2
		PayCnf.ApiV3 = req.ApiV3
		PayCnf.MchId = req.MchId
		PayCnf.CertPath = req.CertPath
		PayCnf.KeyPath = req.KeyPath
		e.Orm.Save(&PayCnf)
	}else {
		trade:=models.WeChatPay{
			MchId: req.MchId,
			ApiV3: req.ApiV3,
			ApiV2: req.ApiV2,
			CertPath: req.CertPath,
			KeyPath: req.KeyPath,
		}
		trade.CreateBy = userDto.UserId
		trade.CId = userDto.CId
		trade.Enable = req.Enable
		trade.Refund = req.Refund
		e.Orm.Create(&trade)
	}
	e.OK("","successful")

	return
}

func (e *PayWechat) Detail(c *gin.Context) {
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

	var data models.WeChatPay
	e.Orm.Model(&models.WeChatPay{}).Where("c_id = ? ",userDto.CId).Limit(1).Find(&data)

	e.OK(data,"successful")
	return
}
