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

// Authenticator 获取token
// @Summary 登陆
// @Description 获取token
// @Description LoginHandler can be used by clients to get a jwt token.
// @Description Payload needs to be json in the form of {"username": "USERNAME", "password": "PASSWORD"}.
// @Description Reply will be of the form {"token": "TOKEN"}.
// @Description dev mode：It should be noted that all fields cannot be empty, and a value of 0 can be passed in addition to the account password
// @Description 注意：开发模式：需要注意全部字段不能为空，账号密码外可以传入0值
// @Tags 登陆
// @Accept  application/json
// @Product application/json
// @Param account body Login  true "account"
// @Success 200 {string} string "{"code": 200, "expire": "2019-08-07T12:45:48+08:00", "token": ".eyJleHAiOjE1NjUxNTMxNDgsImlkIjoiYWRtaW4iLCJvcmlnX2lhdCI6MTU2NTE0OTU0OH0.-zvzHvbg0A" }"
// @Router /api/v1/login [post]
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
	if loginVals.Phone != "" { //手机号登录
		userRow, role, e := loginVals.GetUserPhone(db)
		messageData:=map[string]interface{}{
			"ipaddr":common.GetClientIP(c),
			"user":loginVals.Phone,
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
	} else {//用户名或者手机号登录
		//先查询下用户名
		userRow, role, userErr := loginVals.GetUser(db)
		var lastErr error
		if userErr !=nil{//如果查询用户名没有,继续查询手机号
			//然后查询下手机号
			loginVals.Phone = loginVals.UserName
			userRow, role, lastErr = loginVals.GetUserPhone(db)
		}

		messageData:=map[string]interface{}{
			"ipaddr":common.GetClientIP(c),
			"user":loginVals.UserName,
			"login_time":time.Now(),
			"source":"PC",
			"role":global.LoginRoleCompany,
			"client":global.LogIngPC,
			"user_type":global.LogIngUserType,
		}
		if lastErr == nil {
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
			return nil, lastErr
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
		//c.Set("dataScope", r.DataScope)
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
