/**
@Author: chaoqun
* @Date: 2023/7/20 22:32
*/
package apis

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"go-admin/app/company/models"
	"go-admin/app/company/service/dto"
	models2 "go-admin/cmd/migrate/migration/models"
	customUser "go-admin/common/jwt/user"
	"go-admin/common/redis_db"
	"go-admin/global"
	"strings"
	"time"
)
type WeApp struct {
	api.Api
}

func (e WeApp)LoginList(c *gin.Context)  {
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
	loginResult:=make([]map[string]interface{},0)
	registerResult:=make([]map[string]interface{},0)
	//查询是否存在配置

	for _,val:=range strings.Split(global.LoginStr,","){
		var registerCnf models.CompanyRegisterCnf
		e.Orm.Model(&models.CompanyRegisterCnf{}).Where("c_id = ? and type = 0 and value = ?" ,userDto.CId,val).Limit(1).Find(&registerCnf)
		enable:=true
		if registerCnf.Id > 0 {
			enable = registerCnf.Enable
		}
		loginResult = append(loginResult, map[string]interface{}{
			"cn":global.LoginCnfToCh(val) + "登录",
			"value":val,
			"enable":enable,
		})
	}

	for _,val:=range strings.Split(global.RegisterStr,","){
		var registerCnf models.CompanyRegisterCnf
		e.Orm.Model(&models.CompanyRegisterCnf{}).Where("c_id = ? and type = 1 and value = ?" ,userDto.CId,val).Limit(1).Find(&registerCnf)
		enable:=true
		if registerCnf.Id > 0 {
			enable = registerCnf.Enable
		}
		registerResult = append(registerResult, map[string]interface{}{
			"cn":global.LoginCnfToCh(val) + "注册",
			"value":val,
			"enable":enable,
		})
	}
	result :=map[string]interface{}{
		"login":loginResult,
		"register":registerResult,
	}
	e.OK(result,"successful")
	return
}


func (e WeApp)UpdateLoginList(c *gin.Context)  {
	req:=dto.UpdateLogin{}
	err := e.MakeContext(c).
		Bind(&req).
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
	if req.Val == ""{
		e.Error(500, nil, "请选择类型")
		return
	}

	var registerCnf models.CompanyRegisterCnf
	e.Orm.Model(&models.CompanyRegisterCnf{}).Where("c_id = ? and type = ? and value = ?",userDto.CId,req.T,req.Val).Limit(1).Find(&registerCnf)
	if registerCnf.Id == 0 {
		registerCnf = models.CompanyRegisterCnf{
			CId: userDto.CId,
			Type: req.T,
			Value: req.Val,
			Enable: req.Enable,
		}
		registerCnf.CreateBy = userDto.UserId
		e.Orm.Create(&registerCnf)
	}else {
		registerCnf.Enable = req.Enable
		e.Orm.Save(&registerCnf)
	}
	//根据前段的开启和关闭进行一个数据缓存

	//需要写入到redis中,实现配置
	l1:=make([]string,0)
	r1:=make([]string,0)
	for _,val:=range strings.Split(global.RegisterStr,","){
		var registerRow models.CompanyRegisterCnf
		e.Orm.Model(&models.CompanyRegisterCnf{}).Select("id,enable").Where("c_id = ?  and value = ? and type = 1",userDto.CId,val).Limit(1).Find(&registerRow)
		enable := false
		if registerRow.Id == 0 {
			enable = true
		}else {
			enable = registerRow.Enable
		}
		if enable{
			r1 = append(r1,val)
		}

	}
	for _,val:=range strings.Split(global.LoginStr,","){
		var registerRow models.CompanyRegisterCnf
		e.Orm.Model(&models.CompanyRegisterCnf{}).Select("id,enable").Where("c_id = ? and value = ? and type = 0",userDto.CId,val).Limit(1).Find(&registerRow)
		enable := false
		if registerRow.Id == 0 {
			enable = true
		}else {
			enable = registerRow.Enable
		}

		if enable{
			l1 = append(l1,val)
		}

	}

	fmt.Println("r1",r1)
	fmt.Println("l1",l1)
	value :=redis_db.LoginValue{
		PwdLen: 6,
		PwdComplexity: "number",
		Register: strings.Join(r1,","),
		Login: strings.Join(l1,","),
	}
	redisData :=redis_db.RedisLoginCnf{
		ConfigDesc: "注册规则",
		CreateTime: time.Now().Unix(),
		IsUse:1,
		Value: value,
	}

	redis_db.SetLoginCnf(userDto.CId,redisData)

	e.OK("","successful")
	return
}


func (e WeApp)Navbar(c *gin.Context) {
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
	//查询是否有特殊配置

	navList :=make([]models2.WeAppGlobalNavCnf,0)
	e.Orm.Model(&models2.WeAppGlobalNavCnf{}).Where("enable = true").Find(&navList)

	result:=make([]models2.WeAppGlobalNavCnf,0)
	for _,row:=range navList{
		var object models2.CompanyNavCnf
		e.Orm.Model(&models2.CompanyNavCnf{}).Where("c_id = ? and g_id = ?",userDto.CId,row.Id).Limit(1).Find(&object)

		if object.Id > 0 {
			row.UserEnable = object.Enable

		}else {
			row.UserEnable = row.Enable
		}
		result = append(result,row)

	}


	//如果没有特殊配置
	e.OK(result,"successful")
	return
}

func (e WeApp)UpdateNavbar(c *gin.Context)  {
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
	navList :=make([]models2.WeAppGlobalNavCnf,0)
	e.Orm.Model(&models2.WeAppGlobalNavCnf{}).Where("enable = true").Find(&navList)

	navLibMap:=make(map[string]interface{},0)
	_=json.Unmarshal([]byte(dto.NavLib),&navLibMap)
	for _,row:=range navList{
		var object models2.CompanyNavCnf
		e.Orm.Model(&models2.CompanyNavCnf{}).Where("c_id = ? and g_id = ?",userDto.CId,row.Id).Limit(1).Find(&object)

		if object.Id > 0 {
			//配置了并且是关闭的,那就返回吧
			if !object.Enable{
				continue
			}
		}
		navLibMap["iconPath"] = row.IconPath
		navLibMap["selectedIconPath"] = row.SelectedIconPath
		navLibMap["text"] = row.Text
		navLibMap["iconClass"] = row.IconClass
		navLibMap["link"] = map[string]interface{}{
			"name":row.Name,
			"title":row.Text,
			"wap_url":row.WapUrl,
			"parent":"MALL_LINK",
		}

	}

}