package apis

import (
	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/captcha"
)

type System struct {
	api.Api
}

// GenerateCaptchaHandler 获取验证码
// @Summary 获取验证码
// @Description 获取验证码
// @Tags 登陆
// @Success 200 {object} response.Response{data=string,id=string,msg=string} "{"code": 200, "data": [...]}"
// @Router /api/v1/captcha [get]
func (e System) GenerateCaptchaHandler(c *gin.Context) {
	err := e.MakeContext(c).Errors
	if err != nil {
		e.Error(500, err, "服务初始化失败！")
		return
	}
	//driverString:=base64Captcha.DriverString{
	//	Height:          60,
	//	Width:           200,
	//	NoiseCount:      0,     //噪点数
	//	ShowLineOptions: 1, //干扰线
	//	Length:          4,
	//	Source:          "123456789",
	//	Fonts: []string{"wqy-microhei.ttc"},
	//}
	//var driver base64Captcha.Driver = driverString.ConvertFonts()
	//cap1 := base64Captcha.NewCaptcha(driver, store)
	//DriverDigit := base64Captcha.NewDriverDigit(80, 240, 4, 0, 1)
	//
	//cap1 := base64Captcha.NewCaptcha(DriverDigit, base64Captcha.DefaultMemStore)

	id, b64s, err := captcha.DriverDigitFunc()
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
