package apis

import (
	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/captcha"
)

type GoAdminSystem struct {
	api.Api
}
func (e GoAdminSystem)GoAdmin(c *gin.Context) {
	c.String(200, "dcy-store.api")
}

func (e GoAdminSystem)GenerateCaptchaHandler(c *gin.Context) {

	err := e.MakeContext(c).Errors
	if err != nil {
		e.Error(500, err, "服务初始化失败！")
		return
	}
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