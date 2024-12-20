/**
@Author: chaoqun
* @Date: 2023/9/25 23:49
*/
package apis

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"go-admin/cmd/migrate/migration/models"
	"go-admin/common/actions"
	customUser "go-admin/common/jwt/user"
	"os"
	"os/exec"
	"strings"
)


type PayWechat struct {
	api.Api
}
type WechatAppPayReq struct {
	AppId string `json:"app_id"`
	AppSecret string `json:"app_secret"`
	MchId string  `json:"mch_id" `
	ApiV2 string `json:"api_v2"`
	ApiV3 string `json:"api_v3"`
	CertText string `json:"cert_text" `
	KeyText string  `json:"key_text" `
	Enable bool `json:"enable"`
	Refund bool `json:"refund"`
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
	OfficialAppSecret string `json:"official_app_secret"`
	Enable bool `json:"enable"`
	Refund bool `json:"refund"`
}

func analysisCert(cid int,cert string) (searNumber string,err error)  {
	//写入到缓存demo文件中
	cacheFile:=fmt.Sprintf("%v_cert.pem",cid)

	file, _ := os.OpenFile(cacheFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if _, err = file.WriteString(cert);err!=nil{

		return "", err
	}
	cmd := exec.Command("bash","-c",fmt.Sprintf("openssl x509 -in %v -noout -serial",cacheFile))

	output, err := cmd.Output()
	if err != nil {
		fmt.Println("执行错误",err)
		return
	}
	//serial=
	data :=strings.TrimSpace(string(output))
	data = strings.Replace(data,"serial=","",-1)
	_=os.Remove(cacheFile)
	return data,nil

}

func (e *PayWechat) CreateWechatAppPay(c *gin.Context) {
	req := WechatAppPayReq{}
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

	var PayCnf models.WeChatAppPay

	e.Orm.Model(&PayCnf).Scopes(actions.PermissionSysUser(PayCnf.TableName(),userDto)).Limit(1).Find(&PayCnf)


	searchNumber,err:=analysisCert(userDto.CId,req.CertText)
	if err!=nil{
		e.Error(500, nil,"证书解析失败")
		return
	}
	if PayCnf.Id > 0 {
		PayCnf.SerialNumber = searchNumber
		PayCnf.Refund = req.Refund
		PayCnf.ApiV2 = strings.TrimSpace(req.ApiV2)
		PayCnf.ApiV3 = strings.TrimSpace(req.ApiV3)
		PayCnf.MchId = strings.TrimSpace(req.MchId)
		PayCnf.CertText = strings.TrimSpace(req.CertText)
		PayCnf.KeyText = strings.TrimSpace(req.KeyText)
		PayCnf.AppId = strings.TrimSpace(req.AppId)
		PayCnf.AppSecret = strings.TrimSpace(req.AppSecret)
		PayCnf.Enable = req.Enable
		e.Orm.Save(&PayCnf)
	}else {
		trade:=models.WeChatAppPay{
			MchId: strings.TrimSpace(req.MchId),
			ApiV3: strings.TrimSpace(req.ApiV3),
			ApiV2: strings.TrimSpace(req.ApiV2),
			CertText: strings.TrimSpace(req.CertText),
			KeyText: strings.TrimSpace(req.KeyText),
			AppId: strings.TrimSpace(req.AppId),
			AppSecret:strings.TrimSpace(req.AppSecret),
			SerialNumber: searchNumber,
		}
		trade.CreateBy = userDto.UserId
		trade.CId = userDto.CId
		trade.Enable = req.Enable
		trade.Refund = req.Refund
		trade.Layer = 0
		e.Orm.Create(&trade)
	}
	e.OK("","操作成功")

	return
}
func (e *PayWechat) WechatAppDetail(c *gin.Context) {
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

	var data models.WeChatAppPay
	e.Orm.Model(&models.WeChatAppPay{}).Scopes(actions.PermissionSysUser(data.TableName(),userDto)).Limit(1).Find(&data)

	e.OK(data,"操作成功")
	return
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

	var PayCnf models.WeChatOfficialPay

	e.Orm.Model(&PayCnf).Scopes(actions.PermissionSysUser(PayCnf.TableName(),userDto)).Limit(1).Find(&PayCnf)


	searchNumber,err:=analysisCert(userDto.CId,req.CertText)
	if err!=nil{
		e.Error(500, nil,"证书解析失败")
		return
	}
	if PayCnf.Id > 0 {
		PayCnf.SerialNumber = searchNumber
		PayCnf.Refund = req.Refund
		PayCnf.ApiV2 = strings.TrimSpace(req.ApiV2)
		PayCnf.ApiV3 = strings.TrimSpace(req.ApiV3)
		PayCnf.MchId = strings.TrimSpace(req.MchId)
		PayCnf.CertText = strings.TrimSpace(req.CertText)
		PayCnf.KeyText = strings.TrimSpace(req.KeyText)
		PayCnf.OfficialAppId = strings.TrimSpace(req.OfficialAppId)
		PayCnf.OfficialAppSecret = strings.TrimSpace(req.OfficialAppSecret)
		PayCnf.Enable = req.Enable
		e.Orm.Save(&PayCnf)
	}else {
		trade:=models.WeChatOfficialPay{
			MchId: strings.TrimSpace(req.MchId),
			ApiV3: strings.TrimSpace(req.ApiV3),
			ApiV2: strings.TrimSpace(req.ApiV2),
			CertText: strings.TrimSpace(req.CertText),
			KeyText: strings.TrimSpace(req.KeyText),
			OfficialAppId: strings.TrimSpace(req.OfficialAppId),
			OfficialAppSecret:strings.TrimSpace(req.OfficialAppSecret),
			SerialNumber: searchNumber,
		}
		trade.CreateBy = userDto.UserId
		trade.CId = userDto.CId
		trade.Enable = req.Enable
		trade.Refund = req.Refund
		trade.Layer = 0
		e.Orm.Create(&trade)
	}
	e.OK("","操作成功")

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

	e.OK(data,"操作成功")
	return
}
