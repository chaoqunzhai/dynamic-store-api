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
	"strings"
)


type PayWechat struct {
	api.Api
}
type WechatPayReq struct {
	AppId string `json:"app_id"`
	AppSecret string `json:"app_secret"`
	MchId string  `json:"mch_id" `
	ApiV2 string `json:"api_v2"`
	ApiV3 string `json:"api_v3"`
	CertText string `json:"cert_text" `
	KeyText string  `json:"key_text" `
	OfficialAppId string `json:"official_app_id"`
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

	var AppCnf models.WeChatAppIdCnf
	e.Orm.Model(&AppCnf).Scopes(actions.PermissionSysUser(AppCnf.TableName(),userDto)).Limit(1).First(&AppCnf)


	if AppCnf.Id > 0 {

		AppCnf.AppId = req.AppId
		AppCnf.AppSecret = req.AppSecret
		e.Orm.Save(&AppCnf)
	}else {
		cnf :=models.WeChatAppIdCnf{
			AppId: req.AppId,
			AppSecret: req.AppSecret,
		}
		cnf.Enable = true
		cnf.Layer = 0
		cnf.CreateBy = userDto.UserId
		cnf.CId = userDto.CId
		e.Orm.Create(&cnf)
	}
	var PayCnf models.WeChatOfficialPay

	e.Orm.Model(&PayCnf).Scopes(actions.PermissionSysUser(PayCnf.TableName(),userDto)).Limit(1).Find(&PayCnf)

	if PayCnf.Id > 0 {

		PayCnf.Enable = req.Enable
		PayCnf.Refund = req.Refund
		PayCnf.ApiV2 = strings.TrimSpace(req.ApiV2)
		PayCnf.ApiV3 = strings.TrimSpace(req.ApiV3)
		PayCnf.MchId = strings.TrimSpace(req.MchId)
		PayCnf.CertText = strings.TrimSpace(req.CertText)
		PayCnf.KeyText = strings.TrimSpace(req.KeyText)
		PayCnf.OfficialAppId = strings.TrimSpace(req.OfficialAppId)
		e.Orm.Save(&PayCnf)
	}else {
		trade:=models.WeChatOfficialPay{
			MchId: strings.TrimSpace(req.MchId),
			ApiV3: strings.TrimSpace(req.ApiV3),
			ApiV2: strings.TrimSpace(req.ApiV2),
			CertText: strings.TrimSpace(req.CertText),
			KeyText: strings.TrimSpace(req.KeyText),
			OfficialAppId: strings.TrimSpace(req.OfficialAppId),
		}
		trade.CreateBy = userDto.UserId
		trade.CId = userDto.CId
		trade.Enable = req.Enable
		trade.Refund = req.Refund
		trade.Layer = 0
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

	var data models.WeChatOfficialPay
	e.Orm.Model(&models.WeChatOfficialPay{}).Scopes(actions.PermissionSysUser(data.TableName(),userDto)).Limit(1).Find(&data)

	//var appCnf models.WeChatAppIdCnf
	//e.Orm.Model(&models.WeChatAppIdCnf{}).Scopes(actions.PermissionSysUser(appCnf.TableName(),userDto)).Limit(1).Find(&appCnf)

	//result:=map[string]interface{}{
	//	"pay":data,
	//	"app":appCnf,
	//}
	//data.AppId = appCnf.AppId
	//data.AppSecret = appCnf.AppSecret
	e.OK(data,"successful")
	return
}
