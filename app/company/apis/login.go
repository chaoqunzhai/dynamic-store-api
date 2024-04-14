package apis

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/captcha"
	sys "go-admin/app/admin/models"
	"go-admin/app/company/models"
	"go-admin/app/company/service"
	"go-admin/common"
	"go-admin/common/business"
	"go-admin/common/systemChan"
	"go-admin/common/web_app"
	"go-admin/global"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"regexp"
	"time"
)
type Login struct {
	api.Api
}
type LoginReq struct {
	//Phone    string `form:"phone" json:"phone" `
	UserName string `form:"username" json:"username"`
	Password string `form:"password" json:"password" binding:"required"`
	Code     string `form:"Code" json:"code" binding:"required"`
	UUID     string `form:"UUID" json:"uuid" binding:"required"`
	Role string `form:"role" json:"role"`
}
type CheckedCompanyReq struct {
	Token string `form:"token" json:"token" binding:"required"`
	SiteId int `form:"site_id" json:"site_id" binding:"required"`
}
var (
	ErrMissingLoginValues   = errors.New("请输入手机号或者密码以及验证码")
	ErrFailedAuthentication = errors.New("手机号或者密码错误")
	ErrInvalidVerification  = errors.New("验证码错误")
	ErrFailedAuthenticationChecked = errors.New("请先登录")
)

func LoginValidCompany(companyId int,tx *gorm.DB) error {

	var companyObject models.Company
	tx.Model(&models.Company{}).Select("id,expiration_time").Where("id = ? and enable = ?", companyId, true).Limit(1).Find(&companyObject)
	if companyObject.Id == 0 {

		return errors.New("账号不存在")
	}
	if companyObject.ExpirationTime.Before(time.Now()) {

		return errors.New("账号已到期,请及时续费")
	}
	return nil
}

func (e Login)UserLogin(c *gin.Context)  {
	req:=LoginReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req, binding.JSON, nil).
		Errors
	if err != nil {
		e.Error(400,ErrMissingLoginValues,ErrMissingLoginValues.Error())
		return
	}


	if !captcha.Verify(req.UUID, req.Code, true) {


		e.Error(400,ErrInvalidVerification,ErrInvalidVerification.Error())
		return
	}

	//查询
	regexPhone:= regexp.MustCompile(`^1[3-9]\d{9}$`)

	//因为大B可能是在多个站点下都有，可能是一个站点管理员，或者是多个站点的员工等等


	if regexPhone.MatchString(req.UserName) {
		//是手机号

		var userList []sys.SysUser
		e.Orm.Model(&sys.SysUser{}).Select("user_id,phone,c_id,password").Where("phone = ? and status = ?",req.UserName,global.SysUserSuccess).Find(&userList)

		if len(userList) > 1 { //多个站点
			//循环 列表 一一校验密码是否对
			//如果密码 正确才可以选择多站点

			fmt.Println("多站点登录")
			loginSuccess:=false
			res:=make(map[string]interface{},0)
			for _,row:=range userList{
				loginOk,_ := pkg.CompareHashAndPassword(row.Password, req.Password)

				if loginOk { //如果检测到密码正确了不往下进行了
					loginSuccess = true

					token, _, tokenErr := service.BuildToken(row.UserId, req.UserName,row.Phone)
					if tokenErr != nil {
						fmt.Println("Mobile build Token Error:", tokenErr.Error())
						e.Error(400,nil,"内部错误")
						return
					}

					//登录成功了 那就返回多个站点
					res=map[string]interface{}{
						"token":  token,
						"select_site":true,
					}
					continue
				}
			}
			if loginSuccess {
				siteList:=make([]map[string]interface{},0)
				for _,row:=range userList{
					var company models.Company
					e.Orm.Model(&company).Where("id = ? and enable = ?",row.CId,true).Limit(1).Find(&company)
					if company.Id == 0 {
						continue
					}
					var subtitle string
					if len(company.Enterprise) > 12 {
						subtitle = company.Enterprise[:12]
					}
					dat :=map[string]interface{}{
						"user_id":row.UserId,
						"id":row.CId,
						"company_name":company.Name,
						"logo":company.Image,
						"subtitle":subtitle,
					}

					if company.Image != ""{
						dat["log"] = business.GetGoodsPathFirst(company.Id,company.Image,global.AvatarPath)
					}
					siteList = append(siteList,dat)
				}
				res["site_list"] = siteList
				e.OK(res, "登录成功")
				return
			}
			e.Error(400,ErrFailedAuthentication,ErrFailedAuthentication.Error())
			return

		}else { //单个站点直接登录

			userRow, role, useErr := e.GetUserPhone(req.UserName,req.Password,e.Orm)
			messageData:=map[string]interface{}{
				"ipaddr":common.GetClientIP(c),
				"user":req.UserName,
				"login_time":time.Now(),
				"source":"PC",
				"client":global.LogIngPC,
				"user_type":global.LogIngPhoneType,
				"role":global.LoginRoleCompany,
			}
			if useErr == nil {
				token, expire, tokenErr := service.BuildToken(userRow.UserId, req.UserName,userRow.Phone)
				if tokenErr != nil {
					fmt.Println("Mobile build Token Error:", tokenErr.Error())
					e.Error(400,nil,"内部错误")
					return
				}
				messageData["c_id"] = userRow.CId
				messageData["user_id"] = userRow.UserId
				systemChan.SendMessage(&systemChan.Message{
					Table: "sys_login_log",
					Data: messageData,
					Orm: e.Orm,
				})
				go func() {
					zap.S().Infof("大B:%v登录,获取移动端配置数据",userRow.CId)
					web_app.SearchAndLoadData(userRow.CId,e.Orm)
				}()

				res:=map[string]interface{}{
					"token":  token,
					"siteId":userRow.CId,
					"expire": expire,
					"role":role,
				}
				e.OK(res, "登录成功")
				return
			}

			e.Error(400,useErr,useErr.Error())
			return


		}


	}else {


		userRow, role, userErr := e.GetUser(req.UserName,req.Password,e.Orm)

		messageData:=map[string]interface{}{
			"ipaddr":common.GetClientIP(c),
			"user":req.UserName,
			"login_time":time.Now(),
			"source":"PC",
			"role":global.LoginRoleCompany,
			"client":global.LogIngPC,
			"user_type":global.LogIngUserType,
		}
		if userErr == nil {
			token, expire, tokenErr := service.BuildToken(userRow.UserId, req.UserName,userRow.Phone)
			if tokenErr != nil {
				fmt.Println("Mobile build Token Error:", tokenErr.Error())

				e.Error(400,nil,"内部错误")
				return
			}
			messageData["c_id"] = userRow.CId
			messageData["user_id"] = userRow.UserId
			systemChan.SendMessage(&systemChan.Message{
				Table: "sys_login_log",
				Data: messageData,
				Orm: e.Orm,
			})
			go func() {
				zap.S().Infof("大B:%v登录,获取移动端配置数据",userRow.CId)
				web_app.SearchAndLoadData(userRow.CId,e.Orm)
			}()

			res:=map[string]interface{}{
				"token":  token,
				"siteId":userRow.CId,
				"expire": expire,
				"role":role,
			}
			e.OK(res, "登录成功")

			return

		}
		e.Error(400,userErr,userErr.Error())
		return
	}

}

func (e *Login) GetUser(username,password string,tx *gorm.DB) (user sys.SysUser, role models.CompanyRole, err error) {
	orm := tx.Table("sys_user").Where("username = ?  and status = '2' and enable = ? ", username,true)
	err = orm.First(&user).Error
	if err != nil {
		err = errors.New("用户不存在")
		return
	}
	_, err = pkg.CompareHashAndPassword(user.Password, password)
	if err != nil {
		err = errors.New("用户名或密码错误")
		return
	}
	orm.Updates(map[string]interface{}{
		"login_time":time.Now(),
	})
	err =LoginValidCompany(user.CId,tx)
	return
}

func (e *Login) GetUserPhone(phone,password string,tx *gorm.DB) (user sys.SysUser, role models.CompanyRole, err error) {
	orm := tx.Table("sys_user").Where("phone = ?  and status = ? and enable = ?",
		phone, global.SysUserSuccess, true)
	err =orm.Limit(1).First(&user).Error
	if err != nil {
		err = errors.New("手机号不存在")
		return
	}

	_, loginErr := pkg.CompareHashAndPassword(user.Password, password)
	if loginErr != nil {
		err = ErrFailedAuthentication
		return
	}

	orm.Updates(map[string]interface{}{
		"login_time":time.Now(),
	})

	err =LoginValidCompany(user.CId,tx)
	return
}
func (e Login)CompanyChecked(c *gin.Context)  {
	req:=CheckedCompanyReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req, binding.JSON, nil).
		Errors
	if err != nil {
		e.Error(400,ErrFailedAuthenticationChecked,ErrFailedAuthenticationChecked.Error())
		return
	}

	//解析token 是否合法
	Claims,tokenErr:=service.ParseToken(req.Token)
	if tokenErr!=nil{

		e.Error(400,nil,"token不合法")
		return
	}
	fmt.Printf("获取到用户的ID:%v 手机号:%v 切换到站点:%v",Claims.UserId,Claims.Phone,req.SiteId)


	//查询站点ID + 手机号获取到的用户做一个token,这个token
	var userObject sys.SysUser
	e.Orm.Model(&userObject).Where("c_id = ? and phone= ? and status = ?",
		req.SiteId,Claims.Phone,global.SysUserSuccess).Limit(1).Find(&userObject)

	if userObject.UserId == 0 {
		e.Error(400,nil,"切换失败,在此商户下无用户配置")
		return
	}
	messageData:=map[string]interface{}{
		"ipaddr":common.GetClientIP(c),
		"user":userObject.Username,
		"login_time":time.Now(),
		"source":"PC",
		"client":global.LogIngPC,
		"user_type":global.LogIngPhoneTypeCheckSite,
		"role":global.LoginRoleCompany,
	}

	token, expire, tokenErr := service.BuildToken(userObject.UserId, userObject.Username,userObject.Phone)
	if tokenErr != nil {
		fmt.Println("Mobile build Token Error:", tokenErr.Error())
		e.Error(400,nil,"内部错误")
		return
	}
	messageData["c_id"] = userObject.CId
	messageData["user_id"] = userObject.UserId
	systemChan.SendMessage(&systemChan.Message{
		Table: "sys_login_log",
		Data: messageData,
		Orm: e.Orm,
	})
	go func() {
		zap.S().Infof("大B:%v切换站点登录,获取移动端配置数据",userObject.CId)
		web_app.SearchAndLoadData(userObject.CId,e.Orm)
	}()
	res:=map[string]interface{}{
		"token":  token,
		"siteId":userObject.CId,
		"expire": expire,
		"role":userObject.RoleId,
	}
	e.OK(res, "登录成功")
	return
}