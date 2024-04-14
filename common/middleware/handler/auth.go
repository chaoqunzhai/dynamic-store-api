package handler

import (
	"errors"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/captcha"
	"go-admin/app/admin/models"
	"go-admin/common"
	"go-admin/common/systemChan"
	"go-admin/common/web_app"
	"go-admin/global"
	"go.uber.org/zap"
	"net/http"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/response"
)

var (
	ErrMissingLoginValues   = errors.New("请输入手机号或者密码以及验证码")
	ErrFailedAuthentication = errors.New("手机号或者密码错误")
	ErrInvalidVerification  = errors.New("验证码错误")
)

// 设置完权限后,需要重新登录,因为一些信息是从token中解析的
func PayloadFunc(data interface{}) jwt.MapClaims {
	if v, ok := data.(map[string]interface{}); ok {
		u, _ := v["user"].(SysUser)
		r, _ := v["role"].(SysRole)
		return jwt.MapClaims{
			jwt.IdentityKey:  u.UserId,
			jwt.RoleIdKey:    r.RoleId,
			jwt.RoleKey:      r.RoleKey,
			jwt.NiceKey:      u.Username,
			jwt.DataScopeKey: r.DataScope,
			jwt.RoleNameKey:  r.RoleName,
		}
	}

	return jwt.MapClaims{}
}

func IdentityHandler(c *gin.Context) interface{} {
	claims := jwt.ExtractClaims(c)
	return map[string]interface{}{
		"IdentityKey": claims["identity"],
		"UserName":    claims["nice"],
		"RoleKey":     claims["rolekey"],
		"UserId":      claims["identity"],
		"RoleIds":     claims["roleid"],
		"DataScope":   claims["datascope"],
	}
}

// 登录函数
func Authenticator(c *gin.Context) (interface{}, error) {
	log := api.GetRequestLogger(c)
	db, err := pkg.GetOrm(c)
	if err != nil {
		log.Errorf("get db error, %s", err.Error())
		response.Error(c, 500, err, "数据库连接获取失败")
		return nil, ErrFailedAuthentication
	}

	var loginVals Login

	if err = c.ShouldBind(&loginVals); err != nil {

		return nil, ErrMissingLoginValues
	}
	if !captcha.Verify(loginVals.UUID, loginVals.Code, true) {

		return nil, ErrInvalidVerification
	}

	//查询
	regexPhone:= regexp.MustCompile(`^1[3-9]\d{9}$`)

	//因为大B可能是在多个站点下都有，可能是一个站点管理员，或者是多个站点的员工等等




	if regexPhone.MatchString(loginVals.UserName) {
		//是手机号

		var userList []SysUser
		db.Model(&SysUser{}).Select("user_id,phone,c_id").Where("phone = ? and status = ?",loginVals.UserName,global.SysUserSuccess).Find(&userList)

		if len(userList) > 0 {

			return map[string]interface{}{"user_phone": loginVals.UserName, "select": true}, nil
		}
		userRow, role, e := loginVals.GetUserPhone(db)
		messageData:=map[string]interface{}{
			"ipaddr":common.GetClientIP(c),
			"user":loginVals.UserName,
			"login_time":time.Now(),
			"source":"PC",
			"client":global.LogIngPC,
			"user_type":global.LogIngPhoneType,
			"role":global.LoginRoleCompany,
		}
		if e == nil {
			messageData["c_id"] = userRow.CId
			messageData["user_id"] = userRow.UserId
			systemChan.SendMessage(&systemChan.Message{
				Table: "sys_login_log",
				Data: messageData,
				Orm: db,
			})
			go func() {
				zap.S().Infof("大B:%v登录,获取移动端配置数据",userRow.CId)
				web_app.SearchAndLoadData(userRow.CId,db)
			}()
			return map[string]interface{}{"user": userRow, "role": role}, nil
		} else {
			return nil, e
		}


	}else {


		userRow, role, userErr := loginVals.GetUser(db)

		messageData:=map[string]interface{}{
			"ipaddr":common.GetClientIP(c),
			"user":loginVals.UserName,
			"login_time":time.Now(),
			"source":"PC",
			"role":global.LoginRoleCompany,
			"client":global.LogIngPC,
			"user_type":global.LogIngUserType,
		}
		if userErr == nil {
			messageData["c_id"] = userRow.CId
			messageData["user_id"] = userRow.UserId
			systemChan.SendMessage(&systemChan.Message{
				Table: "sys_login_log",
				Data: messageData,
				Orm: db,
			})
			go func() {
				zap.S().Infof("大B:%v登录,获取移动端配置数据",userRow.CId)
				web_app.SearchAndLoadData(userRow.CId,db)
			}()
			return map[string]interface{}{"user": userRow, "role": role}, nil
		} else {
			return nil, userErr
		}
	}


}

// LogOut
// @Summary 退出登录
// @Description 获取token
// LoginHandler can be used by clients to get a jwt token.
// Reply will be of the form {"token": "TOKEN"}.
// @Accept  application/json
// @Product application/json
// @Success 200 {string} string "{"code": 200, "msg": "成功退出系统" }"
// @Router /logout [post]
// @Security Bearer
func LogOut(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "退出成功",
	})

}

func Authorizator(data interface{}, c *gin.Context) bool {

	if v, ok := data.(map[string]interface{}); ok {
		u, _ := v["user"].(models.SysUser)
		//r, _ := v["role"].(models.SysRole)
		//c.Set("role", r.RoleName)
		//c.Set("roleIds", r.RoleId)
		c.Set("cId", u.CId)
		c.Set("userId", u.UserId)
		c.Set("userName", u.Username)

		return true
	}
	return false
}

func Unauthorized(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  message,
	})
}
