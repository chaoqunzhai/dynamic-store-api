package apis

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/google/uuid"
	"github.com/mojocn/base64Captcha"
	"go-admin/common/global"
)
type configJsonBody struct {
	Id            string
	CaptchaType   string
	VerifyValue   string
	DriverAudio   *base64Captcha.DriverAudio
	DriverString  *base64Captcha.DriverString
	DriverChinese *base64Captcha.DriverChinese
	DriverMath    *base64Captcha.DriverMath
	DriverDigit   *base64Captcha.DriverDigit
}
type GoAdminSystem struct {
	api.Api
}
func (e GoAdminSystem)GoAdmin(c *gin.Context) {
	c.String(200, fmt.Sprintf( "dcy-store %s",global.Version))
}
func DriverDigitFunc() (id, b64s string, err error) {
	e := configJsonBody{}
	e.Id = uuid.New().String()

	e.DriverDigit = base64Captcha.NewDriverDigit(80, 240, 4, 0.2, 16)
	driver := e.DriverDigit

	return base64Captcha.NewCaptcha(driver, base64Captcha.DefaultMemStore).Generate()
}
func (e GoAdminSystem)GenerateCaptchaHandler(c *gin.Context) {

	err := e.MakeContext(c).Errors
	if err != nil {
		e.Error(500, err, "服务初始化失败！")
		return
	}
	id, b64s, err := DriverDigitFunc()
	if err != nil {
		e.Logger.Errorf("DriverDigitFunc error, %s", err.Error())
		e.Error(500, err, "验证码获取失败")

		return
	}
	e.Custom(gin.H{
		"code": 200,
		"data": b64s,
		"id":   id,
		"msg":  "success",
	})
}
