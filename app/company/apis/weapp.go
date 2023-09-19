/*
*
@Author: chaoqun
* @Date: 2023/7/20 22:32
*/
package apis

import (
	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"go-admin/app/company/models"
	"go-admin/app/company/service/dto"
	models2 "go-admin/cmd/migrate/migration/models"
	cDto "go-admin/common/dto"
	customUser "go-admin/common/jwt/user"
	"go-admin/common/redis_db"
	"go-admin/common/web_app"
	"go-admin/global"
	"strings"
	"time"
)

type WeApp struct {
	api.Api
}

func (e WeApp) LoginList(c *gin.Context) {
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
	loginResult := make([]map[string]interface{}, 0)
	registerResult := make([]map[string]interface{}, 0)
	//查询是否存在配置

	for _, val := range strings.Split(global.LoginStr, ",") {
		var registerCnf models.CompanyRegisterCnf
		e.Orm.Model(&models.CompanyRegisterCnf{}).Where("c_id = ? and type = 0 and value = ?", userDto.CId, val).Limit(1).Find(&registerCnf)
		enable := true
		if registerCnf.Id > 0 {
			enable = registerCnf.Enable
		}
		loginResult = append(loginResult, map[string]interface{}{
			"cn":     global.LoginCnfToCh(val) + "登录",
			"value":  val,
			"enable": enable,
		})
	}

	for _, val := range strings.Split(global.RegisterStr, ",") {
		var registerCnf models.CompanyRegisterCnf
		e.Orm.Model(&models.CompanyRegisterCnf{}).Where("c_id = ? and type = 1 and value = ?", userDto.CId, val).Limit(1).Find(&registerCnf)
		enable := true
		if registerCnf.Id > 0 {
			enable = registerCnf.Enable
		}
		registerResult = append(registerResult, map[string]interface{}{
			"cn":     global.LoginCnfToCh(val) + "注册",
			"value":  val,
			"enable": enable,
		})
	}
	result := map[string]interface{}{
		"login":    loginResult,
		"register": registerResult,
	}
	e.OK(result, "successful")
	return
}

func (e WeApp) UpdateLoginList(c *gin.Context) {
	req := dto.UpdateLogin{}
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
	if req.Val == "" {
		e.Error(500, nil, "请选择类型")
		return
	}

	var registerCnf models.CompanyRegisterCnf
	e.Orm.Model(&models.CompanyRegisterCnf{}).Where("c_id = ? and type = ? and value = ?", userDto.CId, req.T, req.Val).Limit(1).Find(&registerCnf)
	if registerCnf.Id == 0 {
		registerCnf = models.CompanyRegisterCnf{
			CId:    userDto.CId,
			Type:   req.T,
			Value:  req.Val,
			Enable: req.Enable,
		}
		registerCnf.CreateBy = userDto.UserId
		e.Orm.Create(&registerCnf)
	} else {
		registerCnf.Enable = req.Enable
		e.Orm.Save(&registerCnf)
	}
	//根据前段的开启和关闭进行一个数据缓存

	//需要写入到redis中,实现配置
	l1 := make([]string, 0)
	r1 := make([]string, 0)
	for _, val := range strings.Split(global.RegisterStr, ",") {
		var registerRow models.CompanyRegisterCnf
		e.Orm.Model(&models.CompanyRegisterCnf{}).Select("id,enable").Where("c_id = ?  and value = ? and type = 1", userDto.CId, val).Limit(1).Find(&registerRow)
		enable := false
		if registerRow.Id == 0 {
			enable = true
		} else {
			enable = registerRow.Enable
		}
		if enable {
			r1 = append(r1, val)
		}

	}
	for _, val := range strings.Split(global.LoginStr, ",") {
		var registerRow models.CompanyRegisterCnf
		e.Orm.Model(&models.CompanyRegisterCnf{}).Select("id,enable").Where("c_id = ? and value = ? and type = 0", userDto.CId, val).Limit(1).Find(&registerRow)
		enable := false
		if registerRow.Id == 0 {
			enable = true
		} else {
			enable = registerRow.Enable
		}

		if enable {
			l1 = append(l1, val)
		}

	}

	value := redis_db.LoginValue{
		PwdLen:        6,
		PwdComplexity: "number",
		Register:      strings.Join(r1, ","),
		Login:         strings.Join(l1, ","),
	}
	redisData := redis_db.RedisLoginCnf{
		ConfigDesc: "注册规则",
		CreateTime: time.Now().Unix(),
		IsUse:      1,
		Value:      value,
	}

	redis_db.SetLoginCnf(userDto.CId, redisData)

	e.OK("", "successful")
	return
}

func (e WeApp) Navbar(c *gin.Context) {
	req := dto.CompanyGetPageReq{}
	err := e.MakeContext(c).
		Bind(&req).
		MakeOrm().
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	//查询是否有特殊配置
	var data models.Company
	list := make([]models.Company, 0)
	var count int64
	//获取所有的大B
	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(req.GetNeedSearch()),
			cDto.Paginate(req.GetPageSize(), req.GetPageIndex()),
		).Order(global.OrderLayerKey).
		Find(&list).Limit(-1).Offset(-1).
		Count(&count).Error
	navList := make([]models2.WeAppGlobalNavCnf, 0)
	e.Orm.Model(&models2.WeAppGlobalNavCnf{}).Where("enable = true").Find(&navList)

	result := make([]interface{}, 0)
	for _, row := range list {

		navCnf := make([]interface{}, 0)
		for _, nav := range navList {
			var object models.CompanyNavCnf
			e.Orm.Model(&models.CompanyNavCnf{}).Where("c_id = ? and g_id = ?", row.Id, nav.Id).Limit(1).Find(&object)

			if object.Id > 0 {
				nav.UserEnable = object.Enable

			} else {
				nav.UserEnable = row.Enable
			}
			navCnf = append(navCnf, nav)
		}
		row.NavList = navCnf
		result = append(result, row)
	}

	e.PageOK(result, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
	return
}

func (e WeApp) UpdateNavbar(c *gin.Context) {
	req := dto.UpdateNav{}
	err := e.MakeContext(c).
		Bind(&req).
		MakeOrm().
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	var companyObject models2.Company
	e.Orm.Model(&models2.Company{}).Where("id = ?", req.CId).Limit(1).Find(&companyObject)
	if companyObject.Id == 0 {
		e.Error(500, nil, "不存在")
		return
	}
	var count int64
	e.Orm.Model(&models2.WeAppGlobalNavCnf{}).Where("id = ?", req.NavId).Count(&count)
	if count == 0 {
		e.Error(500, nil, "不存在")
		return
	}
	var CompanyNavCnf models.CompanyNavCnf
	e.Orm.Model(&models.CompanyNavCnf{}).Where("g_id = ? and c_id = ?", req.NavId, req.CId).Limit(1).Find(&CompanyNavCnf)
	if CompanyNavCnf.Id == 0 {
		rr := models.CompanyNavCnf{
			Enable: req.Enable,
			CId:    req.CId,
			GId:    req.NavId,
		}

		e.Orm.Create(&rr)
	} else {
		CompanyNavCnf.Enable = req.Enable
		e.Orm.Save(&CompanyNavCnf)
	}

	web_app.SearchAndLoadData(req.CId,e.Orm)
	e.OK("", "successful")
	return
}

func (e WeApp) Quick(c *gin.Context) {
	req := dto.CompanyGetPageReq{}
	err := e.MakeContext(c).
		Bind(&req).
		MakeOrm().
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	var data models.Company
	list := make([]models.Company, 0)
	var count int64
	//获取所有的大B
	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(req.GetNeedSearch()),
			cDto.Paginate(req.GetPageSize(), req.GetPageIndex()),
		).Order(global.OrderLayerKey).
		Find(&list).Limit(-1).Offset(-1).
		Count(&count).Error
	navList := make([]models2.WeAppQuickTools, 0)
	e.Orm.Model(&models2.WeAppQuickTools{}).Where("enable = true").Find(&navList)

	result := make([]interface{}, 0)
	for _, row := range list {

		navCnf := make([]interface{}, 0)
		for _, nav := range navList {
			var object models.CompanyQuickTools
			e.Orm.Model(&models.CompanyQuickTools{}).Select("id,enable").Where("c_id = ? and quick_id = ?", row.Id, nav.Id).Limit(1).Find(&object)

			if object.Id > 0 {
				nav.UserEnable = object.Enable

			} else {

				nav.UserEnable = nav.DefaultShow
			}
			navCnf = append(navCnf, nav)
		}
		row.NavList = navCnf
		result = append(result, row)
	}

	e.PageOK(result, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
	return

}


func (e WeApp) UpdateQuick(c *gin.Context) {
	req := dto.UpdateQuick{}
	err := e.MakeContext(c).
		Bind(&req).
		MakeOrm().
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	var companyObject models2.Company
	e.Orm.Model(&models2.Company{}).Where("id = ?", req.CId).Limit(1).Find(&companyObject)
	if companyObject.Id == 0 {
		e.Error(500, nil, "不存在")
		return
	}
	var count int64
	e.Orm.Model(&models2.WeAppQuickTools{}).Where("id = ?", req.QuickId).Count(&count)
	if count == 0 {
		e.Error(500, nil, "不存在")
		return
	}
	var CompanyNavCnf models.CompanyQuickTools
	e.Orm.Model(&models.CompanyQuickTools{}).Where("quick_id = ? and c_id = ?", req.QuickId, req.CId).Limit(1).Find(&CompanyNavCnf)
	if CompanyNavCnf.Id == 0 {
		rr := models.CompanyQuickTools{
			Enable: req.Enable,
			CId:    req.CId,
			QuickId:    req.QuickId,
		}

		e.Orm.Create(&rr)
	} else {
		CompanyNavCnf.Enable = req.Enable
		e.Orm.Save(&CompanyNavCnf)
	}
	navList := make([]models2.WeAppQuickTools, 0)
	e.Orm.Model(&models2.WeAppQuickTools{}).Where("enable = true").Find(&navList)


	quickToolsData := make([]interface{}, 0)
	for _, row := range navList {
		var object models.CompanyQuickTools
		e.Orm.Model(&models.CompanyQuickTools{}).Where("c_id = ? and quick_id = ?", req.CId, row.Id).Limit(1).Find(&object)
		if object.Id > 0 {
			//配置了并且是关闭的,那就返回吧
			if !object.Enable {
				continue
			}
		}
		//大部分都是一些后期需要加的配置,先保存默认
		//后台DB只需要配置路径即可
		quickRow:=map[string]interface{}{
			"title":row.Name,
			"imageUrl":row.ImageUrl,
			"iconType":"img",
			"style":map[string]interface{}{
				"fontSize":"60",
				"iconBgColorDeg":0,
				"iconBgImg":"",
				"bgRadius":0,
				"iconColor":[]string{	"#000000"},
				"iconColorDeg":0,
			},
			"link":map[string]interface{}{
				"name":row.Name,
				"title":row.Name,
				"wap_url":row.WapUrl,
				"parent":"MALL_LINK",
			},
			"label":map[string]interface{}{
				"control":false,
				"text":"热门",
				"textColor":"#FFFFFF",
				"bgColorStart":"#F83287",
				"bgColorEnd":"#FE3423",
			},
			"icon":"",
		}
		quickToolsData = append(quickToolsData,quickRow)
	}
	makeCnf := web_app.NewMakeWeAppQuickTools(req.CId)

	makeCnf.ToolsData = quickToolsData
	makeCnf.LoadRedis()

	e.OK("", "successful")
}