/**
@Author: chaoqun
* @Date: 2024/2/23 12:21
*/
package apis

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/captcha"
	sys "go-admin/app/admin/models"
	"go-admin/common"
	"go-admin/common/redis_db"
	"go-admin/global"
	"golang.org/x/crypto/bcrypt"
)

type Forgot struct {
	api.Api
}

type GetCodeReq struct {
	Phone string `json:"phone" binding:"required"`
	Code string `json:"code" binding:"required"`
	Uuid   string `json:"uuid"`
}

type VerifyCodeReq struct {
	Phone string `json:"phone" binding:"required"`
	PhoneCode string `json:"phone_code" binding:"required"`
}

type RepassReq struct {
	PasswordConfirm    string    `json:"password_confirm" binding:"required" `
	PasswordNew    string    `json:"password_new" binding:"required"`
	VerifyCode    string    `json:"verify_code" binding:"required"`
	VerifyPhone    string    `json:"verify_phone" binding:"required"`
}

func (e Forgot) GetPhoneCode(c *gin.Context) {
	req := GetCodeReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req, binding.JSON, nil).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	//验证 验证码是否对
	if !captcha.Verify(req.Uuid, req.Code, false) {
		e.Error(500, nil, "验证码错误")
		return
	}
	//1.验证手机号是否存在系统中
	if req.Phone == ""{
		e.Error(500, nil,"请输入手机号")
		return
	}

	var user sys.SysUser
	e.Orm.Model(&user).Where("enable = ? and phone = ?",true,req.Phone).Limit(1).Find(&user)

	if user.UserId == 0{
		e.Error(500, nil,"手机号不存在")
		return
	}
	//2.判断手机号是否在redis中.如果有直接返回,如果无，触发第三步
	redisValue, _ := redis_db.GetPhoneCode(global.ForgotPrefix, req.Phone)
	if redisValue != "" {
		e.OK(200,"验证码已发送")
		return
	}

	//3.调用api发送验证码
	code, sendMsgErr := common.SendSms("密码找回", req.Phone,user.CId, e.Orm)
	if sendMsgErr != nil {
		e.Error(500, nil,sendMsgErr.Error())
		return
	}
	//4.验证码存在redis中 10分钟有效期
	_, _ = redis_db.SetPhoneCode(global.ForgotPrefix, req.Phone, code)

	e.OK(code, "手机验证码发送成功")
	return

}


func (e Forgot) VerifyCode(c *gin.Context) {

	req := VerifyCodeReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req, binding.JSON, nil).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	//1.验证手机号是否存在系统中
	if req.Phone == ""{
		e.Error(500, nil,"请输入手机号")
		return
	}

	var user sys.SysUser
	e.Orm.Model(&user).Where("enable = ? and phone = ?",true,req.Phone).Limit(1).Find(&user)

	if user.UserId == 0{
		e.Error(500, nil,"手机号不存在")
		return
	}

	//2.判断验证码是否在redis中.并且一致。如果不一致 直接返回错误
	//调试阶段 打开这个
	redisValue, _ := redis_db.GetPhoneCode(global.ForgotPrefix, req.Phone)

	if redisValue != req.PhoneCode {
		e.Error(500, nil,"手机号验证码已过期")
		return
	}
	//验证没有问题
	e.OK( req.PhoneCode, "验证成功")
	return


}


func (e Forgot) Repass(c *gin.Context) {
	//1.验证码手机号是否存在
	//2.验证 手机号验证码是否正确
	//3.密码修改

	req := RepassReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req, binding.JSON, nil).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	//1.验证手机号是否存在系统中
	if req.PasswordConfirm == ""{
		e.Error(500, nil,"请输入6-11位密码")
		return
	}

	var user sys.SysUser
	e.Orm.Model(&user).Where("enable = ? and phone = ?",true,req.VerifyPhone).Limit(1).Find(&user)

	if user.UserId == 0{
		e.Error(500, nil,"手机号不存在")
		return
	}

	//2.判断验证码是否在redis中.并且一致。如果不一致 直接返回错误
	redisValue, _ := redis_db.GetPhoneCode(global.ForgotPrefix, req.VerifyPhone)

	if redisValue != req.VerifyCode {
		e.Error(500, nil,"手机验证码已过期")
		return
	}
	//验证没有问题
	updateMap := make(map[string]interface{}, 0)

	var hash []byte
	var hasErr error
	if hash, hasErr = bcrypt.GenerateFromPassword([]byte(req.PasswordConfirm), bcrypt.DefaultCost); hasErr != nil {
		e.Error(500, err, "密码生成失败")
		return
	} else {
		updateMap["password"] = string(hash)
	}

	e.Orm.Model(&user).Where("enable = ? and phone = ?",true,req.VerifyPhone).Updates(&updateMap)


	e.OK("", "密码修改成功")
	return
}