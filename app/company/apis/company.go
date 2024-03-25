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
	cDto "go-admin/common/dto"
	"go-admin/common/jwt/user"
	customUser "go-admin/common/jwt/user"
	"go-admin/common/qiniu"
	"go-admin/common/utils"
	"go-admin/config"
	"go-admin/global"
	"golang.org/x/crypto/bcrypt"
	"os"
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
	e.OK(result, "操作成功")
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


func (e Company)Pie(c *gin.Context) {
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
	splitTableRes := business.GetTableName(userDto.CId, e.Orm)
	countSql:=fmt.Sprintf("SELECT  SUM(order_money) as all_money,DATE_FORMAT(created_at, '%%Y-%%m-%%d') AS date, COUNT(*) AS count FROM %v  GROUP BY  DATE_FORMAT(created_at, '%%Y-%%m-%%d') ORDER BY date;",splitTableRes.OrderTable)

	orderChat:=dto.ResponseOrderData{
		Date: make([]string,0),
		OrderTotalPrice: make([]float64,0),
		OrderTotal: make([]int64,0),
	}
	var rowList []dto.DateCount

	e.Orm.Table(splitTableRes.OrderTable).Raw(countSql).Scan(&rowList)
	for _,row:=range rowList{
		orderChat.Date = append(orderChat.Date,row.Date)
		orderChat.OrderTotalPrice = append(orderChat.OrderTotalPrice,row.AllMoney)
		orderChat.OrderTotal = append(orderChat.OrderTotal,row.Count)
	}
	e.OK(map[string]interface{}{
		"order":orderChat,
	},"操作成功")
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
	OrderRangeTime :=business.GetOrderRangeTime(userDto.CId,e.Orm)
	thisDayS:=fmt.Sprintf("%v 00:00:00",time.Now().Format("2006-01-02"))
	thisDayE:=fmt.Sprintf("%v 23.59.59",time.Now().Format("2006-01-02"))
	thisDaySql:=fmt.Sprintf("c_id = '%v' and  created_at >= '%v' AND created_at <= '%v'",userDto.CId,thisDayS,thisDayE)

	openApprove,_:=service.IsHasOpenApprove(userDto,e.Orm)

	//检测是否开启了审核,如果开启了审核 必须是审核通过后的订单
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
		Line:func() int64{
			var count int64
			e.Orm.Model(&models.Line{}).Where("c_id = ? and enable = ? ",
				userDto.CId,true).Count(&count)
			return count
		}(),
		Salesman:func() int64{
			var count int64
			e.Orm.Model(&sys.SysUser{}).Where("c_id = ? and role_id = ? ",userDto.CId,global.RoleSaleMan).Count(&count)
			return count
		}(),
		WaitOrder:func() int64{
			var count int64
			orm := e.Orm.Table(splitTableRes.OrderTable).Where("c_id = ? and status =  ?",userDto.CId,global.OrderStatusWaitSend)
			if  OrderRangeTime != ""{
				orm = orm.Where(OrderRangeTime)
			}
			//待配送查询时 需要检测是否开启了审批
			if openApprove{
				orm.Where("approve_status = ?",global.OrderApproveOk).Count(&count)
			}else {
				orm.Count(&count)
			}
			//fmt.Println("查询待发货",count)
			return count
		}(),
		RefundWaitOrder:func() int64{
			var count int64

			orm :=e.Orm.Table(splitTableRes.OrderReturn).Where("c_id = ? and status = ? ",userDto.CId,global.RefundDefault)
			if  OrderRangeTime != ""{
				orm = orm.Where(OrderRangeTime)
			}
			orm.Count(&count)

			return count
		}(),
		WaitSelfOrder:func() int64{
			var count int64
			orm :=e.Orm.Table(splitTableRes.OrderTable).Where("c_id = ? and delivery_type = 1 and status =  ?",userDto.CId,global.OrderWaitConfirm).Count(&count)

			if  OrderRangeTime != ""{
				orm = orm.Where(OrderRangeTime)
			}
			orm.Count(&count)

			return count
		}(),
		ThisDayPayOkOrder: func() int64{
			var count int64
			e.Orm.Table(splitTableRes.OrderTable).Where("c_id = ?",userDto.CId).Where(thisDaySql).Count(&count)
			return count
		}(),
		ThisDayNewShop: func() int64{
			var count int64
			e.Orm.Model(&models2.Shop{}).Where("c_id = ? ",userDto.CId).Where(thisDaySql).Count(&count)
			return count
		}(),
		ThisDayPayOkShopUser:func() int64{
			var count int64
			e.Orm.Table(splitTableRes.OrderTable).Where("c_id = ? and 'status' = ? ",userDto.CId,global.OrderStatusWaitSend).Where(thisDaySql).Count(&count)
			return count
		}(),
		ThisDayPayAll: func() string {

			return "0.00"
		}(),
		SurplusEms: func() int{
			var CompanySmsQuotaCnf models2.CompanySmsQuotaCnf

			e.Orm.Model(&models2.CompanySmsQuotaCnf{}).Where("c_id = ?",userDto.CId).Limit(1).Find(&CompanySmsQuotaCnf)
			return  CompanySmsQuotaCnf.Available

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
	e.OK(result,"操作成功")
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
		e.OK(data, "操作成功")
		return
	}
	data["userRule"] = object.UserRule
	data["text"] = object.Text
	e.OK(data, "操作成功")
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

	e.OK("", "操作成功")
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
	e.OK(cnf, "操作成功")
	return
}

func (e Company) Information(c *gin.Context) {
	req := dto.UpdateInfo{}
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
	var object models.Company
	e.Orm.Model(&models.Company{}).Select("image,id").Where("enable = 1 and leader_id = ? ",userDto.UserId).First(&object)

	if object.Id == 0 {

		e.Error(500,nil,"not company")
		return
	}

	updateMap:=map[string]interface{}{
		"enterprise":req.Enterprise,
		"help_phone":req.HelpPhone,
		"help_message":req.HelpMessage,
		"shop_name":req.ShopName,
		"address":req.Address,
		"filings":req.Filings,
		"shop_status":req.ShopStatus,
	}
	if req.ActionImage {
		file, fileErr := c.FormFile("file")
		buckClient :=qiniu.QinUi{CId: userDto.CId}
		buckClient.InitClient()
		var imageUrl string
		if fileErr == nil {
			_, goodsImagePath := GetCosImagePath(global.AvatarPath, file.Filename, userDto.CId)

			if saveErr := c.SaveUploadedFile(file, goodsImagePath); saveErr == nil {

				//1.上传到cos中
				fileName, cosErr := buckClient.PostImageFile(goodsImagePath)
				if cosErr == nil {
					//上传成功了 那就是新的名字
					imageUrl = fileName
				}
				//本地删除
				_ = os.RemoveAll(goodsImagePath)
			}
		}

		//有原头像的,那就需要删除原头像
		if object.Image != "" { //原来是有头像的
			buckClient.RemoveFile(business.GetSiteCosPath(userDto.CId, global.AvatarPath, object.Image))
		}
		updateMap["image"] = imageUrl
	}

	e.Orm.Model(&models.Company{}).Where("id = ?",object.Id).Updates(updateMap)

	e.OK("","更新成功")
	return
}


func (e Company) SmsUseList(c *gin.Context) {
	req := dto.SmsUseGetPage{}
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

	list := make([]models2.CompanySmsRecordLog, 0)
	var count int64

	e.Orm.Model(&models2.CompanySmsRecordLog{}).Unscoped().
		Scopes(
			cDto.MakeCondition(req.GetNeedSearch()),
			cDto.Paginate(req.GetPageSize(), req.GetPageIndex())).Where("c_id = ?",userDto.CId).Order("id desc").
		Find(&list).Limit(-1).Offset(-1).
		Count(&count)
	result:=make([]interface{},0)
	companyList:=make([]int,0)
	for _,row:=range list{
		companyList = append(companyList,row.CId)
	}


	for _,row:=range list{
		body:=""
		if row.Msg != ""{
			body = row.Msg
		}else {
			body = row.Code
		}
		mm := map[string]interface{}{
			"body":body,
			"source":row.Source,
			"phone":row.Phone,
			"id":row.Id,
			"create_time":row.CreatedAt.Format("2006-01-02 15:04:05"),
		}

		result = append(result,mm)
	}
	e.PageOK(result, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
	return

}

func (e Company) SmsCnfUpdate(c *gin.Context) {
	req := dto.CompanySmsUpdate{}
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
	e.Orm.Model(&models2.CompanySmsQuotaCnf{}).Where("c_id = ?",userDto.CId).Updates(map[string]interface{}{
		"record":req.Record,
		"order_notice":req.Enable,
	})
	e.OK("","successful")
	return
}

func (e Company) SmsCnf(c *gin.Context) {
	s := service.Company{}
	err := e.MakeContext(c).
		MakeOrm().
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
	var CompanySmsQuotaCnf models2.CompanySmsQuotaCnf
	e.Orm.Model(&models2.CompanySmsQuotaCnf{}).Where("c_id = ?",userDto.CId).Limit(1).Find(&CompanySmsQuotaCnf)


	e.OK(CompanySmsQuotaCnf,"successful")
	return
}
func (e Company) PayCnf(c *gin.Context) {
	req := dto.CompanyPayReq{}
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
	//fmt.Println("payCnf",req.Source)
	userDto, err := user.GetUserDto(e.Orm, c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}
	var object models2.PayCnf

	e.Orm.Model(&object).Where("c_id = ?",userDto.CId).Limit(1).Find(&object)

	data :=make([]map[string]interface{},0)

	if object.BalanceDeduct {
		data = append(data, map[string]interface{}{
			"value":"余额支付",
			"label":object.BalanceDeduct,
			"key":global.PayEnBalance,
			"type":global.PayTypeBalance,
		})
	}
	if object.Credit {
		data = append(data, map[string]interface{}{
			"value":"授信余额支付",
			"label":object.Credit,
			"key":global.PayEnCredit,
			"type":global.PayTypeCredit,
		})
	}
	if object.CashOn {
		data = append(data, map[string]interface{}{
			"value":"货到付款",
			"label":object.CashOn,
			"key":global.PayEnCashOn,
			"type":global.PayTypeCashOn,
		})
	}
	//代客下单时 多返回的支付方式
	if req.Source == "valet"{
		data = append(data, map[string]interface{}{
			"value":"线下支付",
			"label":true,
			"key":global.PayEnOffline,
			"type":global.PayTypeOffline,
		})
	}else {
		if object.WeChat{
			data = append(data, map[string]interface{}{
				"value":"微信支付",
				"label":object.WeChat,
				"key":global.PayEnWeChat,
				"type":global.PayTypeOnlineWechat,
			})
		}
	}

	e.OK(data,"")
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
		"sort":          100,
		"is_recycle":    0,
		"is_delete":     0,
		"create_time":   time.Now().Format("2006-01-02 15:04:05"),
		"update_time":   time.Now().Format("2006-01-02 15:04:05"),
		"logoImage":     "",
	}

	var object models.Company
	e.Orm.Model(&models.Company{}).Where("enable = 1 and id = ? ",userDto.CId).First(&object)

	if object.Id > 0 {
		var CompanyLineCnf models2.CompanyLineCnf
		//获取配置的最大路线
		e.Orm.Model(&CompanyLineCnf).Where("c_id = ?",userDto.CId).Limit(1).Find(&CompanyLineCnf)
		//当前路线的使用数

		e.Orm.Model(&models2.Line{}).Where("")
		var logoImage  string
		if object.Image != ""{
			logoImage = business.GetGoodsPathFirst(userDto.CId,object.Image,global.AvatarPath)
		}
		storeInfo = map[string]interface{}{
			"store_id":      object.Id,
			"phone":object.Phone,
			"name":object.Name,
			"shop_name":object.ShopName,
			"sys_name":    "动创云",
			"describe":      object.Desc,
			"sort":          object.Layer,
			"create_time":   object.CreatedAt.Format("2006-01-02 15:04:05"),
			"update_time":   object.UpdatedAt.Format("2006-01-02 15:04:05"),
			"start_time":object.CreatedAt.Format("2006-01-02"), //创建时间
			"end_time":object.ExpirationTime.Format("2006-01-02"), //到期时间
			"logoImage":    logoImage,
			"enterprise":object.Enterprise,
			"filings":object.Filings,
			"address":object.Address,
			"shop_status":object.ShopStatus,
			"help_phone":object.HelpPhone,
			"help_message":object.HelpMessage,
			"inventory_module":object.InventoryModule,
			"sale_user_module":object.SaleUserModule,
			"url":config.ExtConfig.H5Url + fmt.Sprintf("?siteId=%v",object.Id),
		}

	}else {
		if userDto.RoleId == global.RoleSuper {
			storeInfo = map[string]interface{}{
				"store_id":      1,
				"store_name":    "动创云",
				"name":"动创云",
				"describe":      global.Describe,
				"sort":          100,
				"create_time":   time.Now().Format("2006-01-02 15:04:05"),
				"update_time":   time.Now().Format("2006-01-02 15:04:05"),
				"logoImage":     "",
			}
		}else {

			e.Error(500,nil,"not company")
			return
		}
	}

	e.OK(storeInfo, "操作成功")
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


func (e Company) AgreementCnf(c *gin.Context)   {

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
	var object models2.Agreement
	e.Orm.Model(&object).Where("c_id = ?",userDto.CId).Limit(1).Find(&object)
	e.OK(object.Value,"")
	return

}

func (e Company) AgreementUpdate(c *gin.Context)   {
	req:=dto.Agreement{}
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

	var object models2.Agreement
	e.Orm.Model(&object).Where("c_id = ?",userDto.CId).Limit(1).Find(&object)
	if object.Id == 0 {
		object.CId = userDto.CId
		object.Value = req.Value
		object.Layer = 0
		object.Enable = true
		e.Orm.Save(&object)
	}else {
		e.Orm.Model(&object).Where("c_id = ?",userDto.CId).Updates(map[string]interface{}{
			"value":req.Value,
		})
	}
	e.OK("","")
	return

}
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
	//fmt.Println("quotaType",quotaType)

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
		msg = "个客户等级可以创建"
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

	e.OK(res,"操作成功")
	return
}