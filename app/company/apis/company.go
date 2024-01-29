package apis

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"
	sys "go-admin/app/admin/models"
	models2 "go-admin/cmd/migrate/migration/models"
	"go-admin/common/business"
	"go-admin/common/jwt/user"
	customUser "go-admin/common/jwt/user"
	"go-admin/common/utils"
	"go-admin/global"
	"golang.org/x/crypto/bcrypt"
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
func (e Company)RenewPass(c *gin.Context)  {
	req:=RenewPass{}
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

	var oldRow sys.SysUser
	e.Orm.Model(&sys.SysUser{}).Scopes(actions.PermissionSysUser(oldRow.TableName(),userDto)).Where("username = ? ", req.UserName).Limit(1).Find(&oldRow)

	if oldRow.UserId != 0 {
		if oldRow.UserId != userDto.UserId {
			e.Error(500, errors.New("登录用户名称不可重复"), "登录用户名称不可重复")
			return
		}
	}

	SysUserUpdateMap:=map[string]interface{}{
		"username":req.UserName,
		"nick_name":req.RealName,
	}
	if req.PasswordConfirm != ""{
		hash, GenerateErr := bcrypt.GenerateFromPassword([]byte(req.PasswordConfirm), bcrypt.DefaultCost)
		if GenerateErr!=nil{
			e.Error(500,GenerateErr,"密码生成失败")
			return
		}
		SysUserUpdateMap["password"] = string(hash)
	}

	e.Orm.Model(&sys.SysUser{}).Where("user_id = ?",userDto.UserId).Updates(SysUserUpdateMap)
	e.OK(200,"更新成功")
	return
}

func (e Company)Article(c *gin.Context) {
	err := e.MakeContext(c).
		MakeOrm().
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	_, err = user.GetUserDto(e.Orm, c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}
	var GlobalArticle []models2.GlobalArticle
	e.Orm.Model(&models2.GlobalArticle{}).Where("enable = ?",true).Order(global.OrderLayerKey).Find(&GlobalArticle)
	Notice:=make([]dto.NoticeRow,0)
	document:=make([]dto.NoticeRow,0)
	for _,row:=range GlobalArticle{
		d:=dto.NoticeRow{
			Name: row.Name,
			Subtitle: row.Subtitle,
			Link: row.Link,
			Time: row.CreatedAt.Format("2006-01-02"),
		}
		if row.Type == 1 {
			Notice = append(Notice,d)
		}else {
			document = append(document,d)
		}
	}

	result:=map[string]interface{}{
		"notice":Notice,
		"document":document,
	}

	e.OK(result,"")
	return
}
func (e Company)Count(c *gin.Context)  {
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
	thisDayS:=fmt.Sprintf("%v 00:00:00",time.Now().Format("2006-01-02"))
	thisDayE:=fmt.Sprintf("%v 23.59.59",time.Now().Format("2006-01-02"))
	thisDaySql:=fmt.Sprintf("c_id = '%v' and  created_at >= '%v' AND created_at <= '%v'",userDto.CId,thisDayS,thisDayE)
	splitTableRes := business.GetTableName(userDto.CId, e.Orm)
	countResponse :=dto.IndexCount{
		Goods: func() int64{
			var count int64
			e.Orm.Model(&models2.Goods{}).Where("c_id = ? ",userDto.CId).Count(&count)
			return count
		}(),
		Shop: func() int64{
			var count int64
			e.Orm.Model(&models2.Shop{}).Where("c_id = ? ",userDto.CId).Count(&count)
			return count
		}(),
		Order:func() int64{
			var count int64
			e.Orm.Table(splitTableRes.OrderTable).Where("c_id = ? ",userDto.CId).Count(&count)
			return count
		}(),
		SelfOrder:func() int64{
			var count int64
			e.Orm.Table(splitTableRes.OrderTable).Where("c_id = ? and delivery_type = 1",userDto.CId).Count(&count)
			return count
		}(),
		WaitOrder:func() int64{
			var count int64
			e.Orm.Table(splitTableRes.OrderTable).Where("c_id = ? and pay_status = ?",userDto.CId,global.OrderStatusWaitSend).Count(&count)
			return count
		}(),
		RefundOrder:func() int64{
			var count int64
			e.Orm.Table(splitTableRes.OrderReturn).Where("c_id = ? ",userDto.CId).Count(&count)
			return count
		}(),
		WaitSelfOrder:func() int64{
			var count int64
			e.Orm.Table(splitTableRes.OrderTable).Where("c_id = ? and delivery_type = 1 and 'status' =  ?",userDto.CId,global.OrderWaitConfirm).Count(&count)
			return count
		}(),
		ThisDayPayOkOrder: func() int64{
			var count int64
			e.Orm.Table(splitTableRes.OrderTable).Where("c_id = ? and ?",userDto.CId,thisDaySql).Count(&count)
			return count
		}(),
		ThisDayNewShop: func() int64{
			var count int64
			e.Orm.Model(&models2.Shop{}).Where("c_id = ? and ? ",userDto.CId,thisDaySql).Count(&count)
			return count
		}(),
		ThisDayPayOkShopUser:func() int64{
			var count int64
			e.Orm.Table(splitTableRes.OrderTable).Where("c_id = ? and ? and 'status' = ? ",userDto.CId,thisDaySql,global.OrderStatusWaitSend).Count(&count)
			return count
		}(),
		ThisDayPayAll: func() string {



			return "0.00"
		}(),
	}
	list:=make([]models2.Orders,0)
	e.Orm.Table(splitTableRes.OrderTable).Select("order_money,after_sales,after_status").Where(thisDaySql).Find(&list)
	var sumMoney float64
	for _,row:=range list{
		if row.AfterSales && row.AfterStatus == global.RefundOk{
			continue
		}
		sumMoney +=row.OrderMoney
	}

	countResponse.ThisDayPayAll = utils.StringDecimal(sumMoney)
	isOpenInventory:=service.IsOpenInventory(userDto.CId,e.Orm)
	var goodsSellOut int64
	if isOpenInventory {
		e.Orm.Model(&models2.Inventory{}).Where("c_id = ? and stock = 0",userDto.CId).Count(&goodsSellOut)
	}else {

		e.Orm.Model(&models2.Goods{}).Where("c_id = ? and inventory = 0",userDto.CId).Count(&goodsSellOut)

	}
	countResponse.GoodsSellOut = goodsSellOut

	e.OK(countResponse,"")
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
		"name":"",
		"sys_name":    "动创云",
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
		//ShopName:=object.Name
		//if object.ShopName != ""{
		//	ShopName = object.ShopName
		//}


		storeInfo = map[string]interface{}{
			"store_id":      object.Id,
			"phone":object.Phone,
			"name":object.ShopName,
			"sys_name":    "动创云",
			"describe":      object.Desc,
			"logo_image_id": 0,
			"sort":          object.Layer,
			"is_recycle":    0,
			"is_delete":     0,
			"create_time":   object.CreatedAt.Format("2006-01-02 15:04:05"),
			"update_time":   object.UpdatedAt.Format("2006-01-02 15:04:05"),
			"start_time":object.CreatedAt.Format("2006-01-02 15:04"), //创建时间
			"end_time":object.ExpirationTime.Format("2006-01-02 15:04"), //到期时间
			"logoImage":     "",
		}

	}else {
		if userDto.RoleId == global.RoleSuper {
			storeInfo = map[string]interface{}{
				"store_id":      0,
				"store_name":    "动创云",
				"name":"动创云",
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
	case "export_worker":
		CompanyCnf := business.GetCompanyCnf(userDto.CId, "export_worker", e.Orm)
		MaxNumber = CompanyCnf["export_worker"]
		msg = fmt.Sprintf("最多同时支持%v个任务执行",MaxNumber)
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