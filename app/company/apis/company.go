package apis

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"
	models2 "go-admin/cmd/migrate/migration/models"
	"go-admin/common/business"
	"go-admin/common/jwt/user"
	customUser "go-admin/common/jwt/user"
	"go-admin/global"
	"time"

	"go-admin/app/company/models"
	"go-admin/app/company/service"
	"go-admin/app/company/service/dto"
	"go-admin/common/actions"
)

type Company struct {
	api.Api
}

func (e Company) MonitorData(c *gin.Context) {
	s := service.Company{}
	err := e.MakeContext(c).
		MakeService(&s.Service).
		MakeOrm().
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	//实时订单数据
	overview := make(map[string]interface{}, 0)
	overview = map[string]interface{}{
		"orderTotalPrice": map[string]string{
			"tday": "0",
			"ytd":  "0.00",
		},
		"orderTotal": map[string]string{
			"tday": "0",
			"ytd":  "0.00",
		},
		"newUserTotal": map[string]string{
			"tday": "0",
			"ytd":  "0.00",
		},
		"consumeUserTotal": map[string]string{
			"tday": "0",
			"ytd":  "0.00",
		},
	}
	//统计
	statistics := make(map[string]interface{}, 0)
	statistics = map[string]interface{}{
		"goodsTotal":       "12",
		"userTotal":        "1",
		"orderTotal":       "0",
		"consumeUserTotal": "0",
	}
	//待办
	pending := make(map[string]interface{}, 0)
	pending = map[string]interface{}{
		"goodsTotal":       "12",
		"userTotal":        "1",
		"orderTotal":       "0",
		"consumeUserTotal": "0",
	}
	//近七日交易走势
	tradeTrend := make(map[string]interface{}, 0)
	tradeTrend = map[string]interface{}{
		"date": []string{
			"2023-05-19",
			"2023-05-20",
			"2023-05-21",
			"2023-05-22",
			"2023-05-23",
			"2023-05-24",
			"2023-05-25",
		},
		"orderTotal": []string{
			"0",
			"0",
			"0",
			"0",
			"0",
			"0",
			"0",
		},
		"orderTotalPrice": []string{
			"0.00",
			"0.00",
			"0.00",
			"0.00",
			"0.00",
			"0.00",
			"0.00",
		},
	}
	result := map[string]interface{}{
		"overview":   overview,
		"statistics": statistics,
		"pending":    pending,
		"tradeTrend": tradeTrend,
	}
	e.OK(result, "successful")
	return
}

func (e Company) Demo(c *gin.Context) {

	c.JSON(200, "")
	return

}
func (e Company) SaveCategory(c *gin.Context) {
	req:=CategoryReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	userDto, err := user.GetUserDto(e.Orm, c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}
	var object models2.CompanyCategory

	e.Orm.Model(&models2.CompanyCategory{}).Scopes(actions.PermissionSysUser(object.TableName(), userDto)).Limit(1).Find(&object)
	object.Type = req.Type
	object.CId = userDto.CId
	object.Enable = true
	e.Orm.Save(&object)
	e.OK("","成功")
	return

}
func (e Company) Category(c *gin.Context) {
	err := e.MakeContext(c).
		MakeOrm().
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	userDto, err := user.GetUserDto(e.Orm, c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}
	var object models2.CompanyCategory
	result :=make(map[string]interface{},0)
	e.Orm.Model(&models2.CompanyCategory{}).Scopes(actions.PermissionSysUser(object.TableName(), userDto)).Where("enable = ?",true).Limit(1).Find(&object)
	if object.Id > 0 {
		result["type"] = object.Type
	}else {
		result["type"] = 1
	}
	e.OK(result,"successful")
	return

}

func (e Company) RegisterCnfInfo(c *gin.Context) {
	err := e.MakeContext(c).
		MakeOrm().
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	//获取配置
	userDto, err := user.GetUserDto(e.Orm, c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}
	data :=map[string]interface{}{
		"userRule": 1,
		"text":     "",
	}
	var object models2.CompanyRegisterRule
	e.Orm.Model(&models2.CompanyRegisterRule{}).Where("c_id = ?",userDto.CId).Limit(1).Find(&object)
	if object.Id == 0 {
		e.OK(data, "successful")
		return
	}
	data["userRule"] = object.UserRule
	data["text"] = object.Text
	e.OK(data, "successful")
	return
}
func (e Company) RegisterCnf(c *gin.Context) {
	req:=dto.RegisterRule{}
	err := e.MakeContext(c).
		Bind(&req,binding.JSON,nil).
		MakeOrm().
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	//获取配置
	userDto, err := user.GetUserDto(e.Orm, c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}
	var object models2.CompanyRegisterRule
	e.Orm.Model(&models2.CompanyRegisterRule{}).Where("c_id = ?",userDto.CId).Limit(1).Find(&object)
	if object.Id == 0 {

		e.Orm.Create(&models2.CompanyRegisterRule{
			CId: userDto.CId,
			UserRule: req.Type,
			Text: req.Text,
		})
	}else {
		e.Orm.Model(&models2.CompanyRegisterRule{}).Where("c_id = ?",userDto.CId).Updates(map[string]interface{}{
			"user_rule":req.Type,
			"text":req.Text,
		})
	}

	e.OK("", "successful")
	return
}
func (e Company) Cnf(c *gin.Context) {
	err := e.MakeContext(c).
		MakeOrm().
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	//获取配置
	userDto, err := user.GetUserDto(e.Orm, c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}
	cnf := business.GetCompanyCnf(userDto.CId, "", e.Orm)
	e.OK(cnf, "successful")
	return
}
func (e Company) Info(c *gin.Context) {
	req := dto.CompanyGetReq{}
	s := service.Company{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	userDto, err := user.GetUserDto(e.Orm, c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}
	storeInfo := map[string]interface{}{
		"store_id":      0,
		"store_name":    "暂无管理系统,请耐心等待！",
		"describe":      global.Describe,
		"logo_image_id": 0,
		"sort":          100,
		"is_recycle":    0,
		"is_delete":     0,
		"create_time":   time.Now().Format("2006-01-02 15:04:05"),
		"update_time":   time.Now().Format("2006-01-02 15:04:05"),
		"logoImage":     "",
	}

	var object models.Company
	e.Orm.Model(&models.Company{}).Where("enable = 1 and leader_id = ? ",userDto.UserId).First(&object)

	if object.Id > 0 {
		ShopName:=object.Name
		if object.ShopName != ""{
			ShopName = object.ShopName
		}
		storeInfo = map[string]interface{}{
			"store_id":      object.Id,
			"store_name":    ShopName,
			"describe":      object.Desc,
			"logo_image_id": 0,
			"sort":          object.Layer,
			"is_recycle":    0,
			"is_delete":     0,
			"create_time":   object.CreatedAt.Format("2006-01-02 15:04:05"),
			"update_time":   object.UpdatedAt.Format("2006-01-02 15:04:05"),
			"logoImage":     "",
		}

	}else {
		if userDto.RoleId == global.RoleSuper {
			storeInfo = map[string]interface{}{
				"store_id":      0,
				"store_name":    global.SysName,
				"describe":      global.Describe,
				"logo_image_id": 0,
				"sort":          100,
				"is_recycle":    0,
				"is_delete":     0,
				"create_time":   time.Now().Format("2006-01-02 15:04:05"),
				"update_time":   time.Now().Format("2006-01-02 15:04:05"),
				"logoImage":     "",
			}
		}
	}

	e.OK(storeInfo, "successful")
	return
}


// Insert 创建Company
// @Summary 创建Company
// @Description 创建Company
// @Tags Company
// @Accept application/json
// @Product application/json
// @Param data body dto.CompanyInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/company [post]
// @Security Bearer


func (e Company) QuotaCnf(c *gin.Context)   {

	err := e.MakeContext(c).
		MakeOrm().
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	quotaType:=c.Query("type")
	fmt.Println("quotaType",quotaType)

	userDto, err := customUser.GetUserDto(e.Orm, c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}
	res:=make(map[string]interface{},0)
	MaxNumber:=0
	var dbCount int64
	var msg string
	switch quotaType {
	case "line":
		CompanyCnf := business.GetCompanyCnf(userDto.CId, "line", e.Orm)
		fmt.Printf("CompanyCnf:%v",CompanyCnf)
		MaxNumber = CompanyCnf["line"]
		var object models.Line
		e.Orm.Model(&object).Scopes(actions.PermissionSysUser(object.TableName(),userDto)).Count(&dbCount)
		msg = "条线路可以创建"
	case "goods":
		CompanyCnf := business.GetCompanyCnf(userDto.CId, "goods", e.Orm)
		MaxNumber = CompanyCnf["goods"]
		var object models.Goods
		e.Orm.Model(&object).Scopes(actions.PermissionSysUser(object.TableName(),userDto)).Count(&dbCount)
		msg = "个商品可以创建"
	case "goods_class":
		CompanyCnf := business.GetCompanyCnf(userDto.CId, "goods_class", e.Orm)
		MaxNumber = CompanyCnf["goods_class"]
		var object models.GoodsClass
		e.Orm.Model(&object).Scopes(actions.PermissionSysUser(object.TableName(),userDto)).Count(&dbCount)
		msg = "个商品分类可以创建"
	case "goods_tag":
		CompanyCnf := business.GetCompanyCnf(userDto.CId, "goods_tag", e.Orm)
		MaxNumber = CompanyCnf["goods_tag"]
		var object models.GoodsTag
		e.Orm.Model(&object).Scopes(actions.PermissionSysUser(object.TableName(),userDto)).Count(&dbCount)
		msg = "个商品标签可以创建"
	case "vip":
		CompanyCnf := business.GetCompanyCnf(userDto.CId, "vip", e.Orm)
		MaxNumber = CompanyCnf["vip"]
		var object models.GradeVip
		e.Orm.Model(&object).Scopes(actions.PermissionSysUser(object.TableName(),userDto)).Count(&dbCount)
		msg = "个VIP等级可以创建"
	case "shop_tag":
		CompanyCnf := business.GetCompanyCnf(userDto.CId, "shop_tag", e.Orm)
		MaxNumber = CompanyCnf["shop_tag"]
		var object models2.ShopTag
		e.Orm.Model(&object).Scopes(actions.PermissionSysUser(object.TableName(),userDto)).Count(&dbCount)
		msg = "个客户标签可以创建"
	case "offline_pay":
		CompanyCnf := business.GetCompanyCnf(userDto.CId, "offline_pay", e.Orm)
		MaxNumber = CompanyCnf["offline_pay"]
		var object models2.OfflinePay
		e.Orm.Model(&object).Scopes(actions.PermissionSysUser(object.TableName(),userDto)).Count(&dbCount)
		msg = "个线下支付可以创建"
	case "role":
		CompanyCnf := business.GetCompanyCnf(userDto.CId, "role", e.Orm)
		MaxNumber = CompanyCnf["role"]
		var object models2.CompanyRole
		e.Orm.Model(&object).Scopes(actions.PermissionSysUser(object.TableName(),userDto)).Count(&dbCount)
		msg = "个角色可以创建"
	case "index_message":
		CompanyCnf := business.GetCompanyCnf(userDto.CId, "index_message", e.Orm)
		MaxNumber = CompanyCnf["index_message"]
		var object models2.Message
		e.Orm.Model(&object).Scopes(actions.PermissionSysUser(object.TableName(),userDto)).Count(&dbCount)
		msg = "条公告消息可以创建"
	case "index_ads":
		CompanyCnf := business.GetCompanyCnf(userDto.CId, "index_ads", e.Orm)
		MaxNumber = CompanyCnf["index_ads"]
		var object models2.Ads
		e.Orm.Model(&object).Scopes(actions.PermissionSysUser(object.TableName(),userDto)).Count(&dbCount)
		msg = "条广告可以创建"

	}
	res["msg"] = msg
	if int(dbCount) <= MaxNumber {

		res["show"] = true
		res["count"] =  MaxNumber - int(dbCount)
	}else {
		res["show"] = false
		res["count"] =  0
	}

	e.OK(res,"successful")
	return
}