package apis

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"
	sys "go-admin/app/admin/models"
	service2 "go-admin/app/company/service"
	models2 "go-admin/cmd/migrate/migration/models"
	"go-admin/common/business"
	customUser "go-admin/common/jwt/user"
	"go-admin/common/utils"
	"go-admin/global"
	"golang.org/x/crypto/bcrypt"

	"go-admin/app/shop/models"
	"go-admin/app/shop/service"
	"go-admin/app/shop/service/dto"
	"go-admin/common/actions"
)

type Shop struct {
	api.Api
}

func (e Shop) MiniApi(c *gin.Context) {
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
	datalist:=make([]models.Shop,0)
	var object models.Shop
	e.Orm.Model(&models.Shop{}).Scopes(actions.PermissionSysUser(object.TableName(), userDto)).Select("id,name,phone,user_id").Order(global.OrderLayerKey).Find(&datalist)

	result:=make([]map[string]interface{},0)
	for _,row:=range datalist{
		result = append(result, map[string]interface{}{
			"id":row.Id,
			"name":row.Name,
			"phone":row.Phone,
			"shop_user_id":row.UserId,
			"text":fmt.Sprintf("%v/%v",row.Name,row.Phone),
		})
	}
	e.OK(result,"操作成功")
	return
}


func (e Shop) GetPage(c *gin.Context) {
    req := dto.ShopGetPageReq{}
    s := service.Shop{}
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
	userDto, err := customUser.GetUserDto(e.Orm, c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}
	var payCnf models2.PayCnf
	e.Orm.Model(&models2.PayCnf{}).Where("c_id = ?", userDto.CId).Limit(1).Find(&payCnf)

	p := actions.GetPermissionFromContext(c)
	list := make([]models.Shop, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取用户信息失败,%s", err.Error()))
        return
	}
	splitTableRes := business.GetTableName(userDto.CId, e.Orm)
	result :=make([]interface{},0)
	for _,row:=range list{
		GradeName :=""
		if row.GradeId > 0 {
			var gradeRow models2.GradeVip
			e.Orm.Model(&models2.GradeVip{}).Select("name,id").Where("id = ? and enable = ?",row.GradeId,true).Limit(1).Find(&gradeRow)
			if gradeRow.Id > 0 {
				GradeName = gradeRow.Name
			}
		}
		if req.Filter == "mini"{
			row.GradeName = GradeName
			result = append(result,row)
			continue
		}
		cache :=row
		if row.LineId > 0 {
			var lineRow models2.Line
			e.Orm.Model(&models2.Line{}).Select("name,id").Where("id = ? and enable = ?",row.LineId,true).Limit(1).Find(&lineRow)
			if lineRow.Id > 0 {
				cache.LineName = lineRow.Name
			}
		}
		if row.Salesman  > 0 {
			var userRow sys.SysUser
			e.Orm.Model(&sys.SysUser{}).Select("username,phone,user_id").Where("user_id = ? and enable = ?",row.Salesman,true).Limit(1).Find(&userRow)
			if userRow.UserId > 0 {
				cache.SalesmanUser = userRow.Username
				cache.SalesmanPhone = userRow.Phone
			}
		}
		cacheTag:=make([]int,0)
		cacheTagName:=make([]string,0)
		for _,t:=range row.Tag{
			cacheTag = append(cacheTag,t.Id)
			cacheTagName = append(cacheTagName,t.Name)
		}
		cache.GradeName = GradeName
		cache.Tags = cacheTag
		cache.TagName = cacheTagName

		cache.DefaultAddress = row.GetFullAddress()
		//统计欠款数据
		var ShopArrears service.ShopArrears
		e.Orm.Table(splitTableRes.OrderTable).Select("SUM(accept_money) as money").Where("c_id = ? and shop_id = ? and pay_type = ? and status = ?",
			row.CId,row.Id,global.PayTypeCashOn,global.OrderPayStatusOfflineWait).Scan(&ShopArrears)
		//fmt.Println("ShopArrears",ShopArrears.Money)
		cache.Arrears = utils.RoundDecimalFlot64(ShopArrears.Money)
		result = append(result,cache)
	}

	resultData:=map[string]interface{}{
		"payCnf":payCnf,
		"list":result,
		"count":int(count),
		"pageIndex":req.GetPageIndex(),
		"pageSize":req.GetPageSize(),
	}

	e.OK(resultData,"")
	return
}

// Get 获取Shop
// @Summary 获取Shop
// @Description 获取Shop
// @Tags Shop
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.Shop} "{"code": 200, "data": [...]}"
// @Router /api/v1/shop/{id} [get]
// @Security Bearer
func (e Shop) Get(c *gin.Context) {
	req := dto.ShopGetReq{}
	s := service.Shop{}
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
	var object models.Shop
	userDto, err := customUser.GetUserDto(e.Orm, c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}
	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取数据失败,%s", err.Error()))
        return
	}

	if object.LineId > 0 {
		var lineRow models2.Line
		e.Orm.Model(&models2.Line{}).Select("name,id").Where("id = ? and enable = ?",object.LineId,true).Limit(1).Find(&lineRow)
		if lineRow.Id > 0 {
			object.LineName = lineRow.Name
		}
	}
	if object.CreateBy  > 0 {
		var userRow sys.SysUser
		e.Orm.Model(&sys.SysUser{}).Select("username,user_id").Where("user_id = ? and enable = ?",object.CreateBy,true).Limit(1).Find(&userRow)
		if userRow.UserId > 0 {
			object.CreateUser = userRow.Username
		}
	}
	cacheTag:=make([]int,0)
	cacheTagName:=make([]string,0)
	for _,t:=range object.Tag{
		cacheTagName = append(cacheTagName,t.Name)
		cacheTag = append(cacheTag,t.Id)
	}
	object.TagName = cacheTagName
	object.Tags = 	 cacheTag
	var shopOrderCount int64
	splitTableRes := business.GetTableName(userDto.CId, e.Orm)
	e.Orm.Table(splitTableRes.OrderTable).Where("shop_id = ?",object.Id).Count(&shopOrderCount)
	object.OrderCount = shopOrderCount


	e.OK( object, "查询成功")
}

//门店创建 分为2种情况
//1.自主创建 上来的表单,有用户名和密码
//2.通过审批列表通过后,点击创建门店,填的表单上来的数据。
func (e Shop) Insert(c *gin.Context) {
    req := dto.ShopInsertReq{}
    s := service.Shop{}
    err := e.MakeContext(c).
        MakeOrm().
		Bind(&req, binding.JSON, nil).
        MakeService(&s.Service).
        Errors
    if err != nil {
        e.Logger.Error(err)
        e.Error(500, err, err.Error())
        return
    }
	// 设置创建人
	req.SetCreateBy(user.GetUserId(c))

	userDto, err := customUser.GetUserDto(e.Orm, c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}

	var countAll int64
	var object models.Shop
	e.Orm.Model(&models.Shop{}).Scopes(actions.PermissionSysUser(object.TableName(), userDto)).Count(&countAll)

	//限制配置
	CompanyCnf := business.GetCompanyCnf(userDto.CId, "shop", e.Orm)

	MaxNumber:=CompanyCnf["shop"]
	if countAll >= int64(MaxNumber) {
		e.Error(500, errors.New(fmt.Sprintf("客户最多只可创建%v个", MaxNumber)), fmt.Sprintf("客户最多只可创建%v个", MaxNumber))
		return
	}

	//路线是否到期检测
	if msg,ExpiredOrNot :=service2.CheckLineExpire(userDto.CId,req.LineId,e.Orm);!ExpiredOrNot{
		e.Error(500, errors.New(msg), msg)
		return
	}
	CompanyLineCnf := business.GetCompanyCnf(userDto.CId, "line_bind_shop", e.Orm)
	MaxLineBindShopNumber := CompanyLineCnf["line_bind_shop"]

	var LineBindShop int64
	e.Orm.Model(&models.Shop{}).Where("c_id = ? and line_id = ?",userDto.CId,req.LineId).Count(&LineBindShop)

	if LineBindShop >= int64(MaxLineBindShopNumber){
		e.Error(500, errors.New(fmt.Sprintf("一条路线最多可绑定%v个客户", MaxLineBindShopNumber)), fmt.Sprintf("一条路线最多可绑定%v个客户", MaxLineBindShopNumber))
		return
	}


	var count int64
	e.Orm.Model(&models.Shop{}).Where("c_id = ?  and name = ? ",userDto.CId,req.Name).Count(&count)
	if count > 0 {
		e.Error(500, errors.New("名称已经存在"), "名称已经存在")
		return
	}
	var userPhone string
	if req.ApproveId > 0 {

		var data models.CompanyRegisterUserVerify
		e.Orm.Model(&models.CompanyRegisterUserVerify{}).Where("c_id = ? and id = ?",userDto.CId, req.ApproveId).Limit(1).Find(&data)

		if data.Id == 0 {
			e.Error(500, errors.New("无此用户"), "无此用户")
			return
		}
		if data.Status != 1 {
			e.Error(500, errors.New("未审批通过"), "未审批通过")
			return
		}
		userPhone = data.Phone
	}else {
		if req.Phone != ""{
			var userCount int64
			e.Orm.Model(&sys.SysShopUser{}).Where("c_id = ? and phone = ? ",userDto.CId,req.Phone).Count(&userCount)
			if userCount > 0 {
				e.Error(500, errors.New("手机号已经存在"), "手机号已经存在")
				return
			}
		}
		if req.UserName != "" {
			var userNameCount int64
			e.Orm.Model(&sys.SysShopUser{}).Where("c_id = ? and username = ? ",userDto.CId,req.UserName).Count(&userNameCount)
			if userNameCount > 0 {
				e.Error(500, errors.New("用户名已经存在"), "用户名已经存在")
				return
			}
		}
		userPhone = req.Phone
	}


	if userPhone !=""{
		var userShopCount int64
		e.Orm.Model(&models.Shop{}).Where("c_id = ? and phone = ? ",userDto.CId,userPhone).Count(&userShopCount)
		if userShopCount > 0 {
			e.Error(500, errors.New("店铺手机号已经存在"), "店铺手机号已经存在")
			return
		}
	}


	err = s.Insert(userDto,&req)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("用户创建失败,%s", err.Error()))
        return
	}



	e.OK(req.GetId(), "创建成功")
}

// Update 修改Shop
// @Summary 修改Shop
// @Description 修改Shop
// @Tags Shop
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.ShopUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/shop/{id} [put]
// @Security Bearer
func (e Shop) Update(c *gin.Context) {
    req := dto.ShopUpdateReq{}
    s := service.Shop{}
    err := e.MakeContext(c).
        MakeOrm().
		Bind(&req, binding.JSON, nil).
        MakeService(&s.Service).
        Errors
    if err != nil {
        e.Logger.Error(err)
        e.Error(500, err, err.Error())
        return
    }
	req.SetUpdateBy(user.GetUserId(c))
	p := actions.GetPermissionFromContext(c)
	userDto, err := customUser.GetUserDto(e.Orm, c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}

	var parentShopRow models.Shop
	e.Orm.Model(&models.Shop{}).Scopes(actions.PermissionSysUser(parentShopRow.TableName(), userDto)).Where("id = ?",req.Id).Limit(1).Find(&parentShopRow)
	if parentShopRow.Id == 0 {
		e.Error(500, errors.New("数据不存在"), "数据不存在")
		return
	}

	//路线是否到期检测
	if msg,ExpiredOrNot :=service2.CheckLineExpire(userDto.CId,req.LineId,e.Orm);!ExpiredOrNot{
		e.Error(500, errors.New(msg), msg)
		return
	}
	//检测路线绑定
	CompanyLineCnf := business.GetCompanyCnf(userDto.CId, "line_bind_shop", e.Orm)
	MaxLineBindShopNumber := CompanyLineCnf["line_bind_shop"]

	var LineBindShop int64
	e.Orm.Model(&models.Shop{}).Where("c_id = ? and line_id = ? and id != ?",userDto.CId,req.LineId,req.Id).Count(&LineBindShop)

	if LineBindShop >= int64(MaxLineBindShopNumber){
		e.Error(500, errors.New(fmt.Sprintf("单路线最多可绑定%v个客户", MaxLineBindShopNumber)), fmt.Sprintf("单路线最多可绑定%v个客户", MaxLineBindShopNumber))
		return
	}

	//名称发生了变化
	if parentShopRow.Name != req.Name {
		var cacheShop models.Shop
		e.Orm.Model(&models.Shop{}).Scopes(actions.PermissionSysUser(cacheShop.TableName(), userDto)).Select("id").Where("name = ?  ",req.Name).Limit(1).Find(&cacheShop)

		if cacheShop.Id != 0 {
			if cacheShop.Id != req.Id {
				e.Error(500, errors.New("名称不可重复"), "名称不可重复")
				return
			}
		}

	}
	//手机号发生了变化
	if req.Phone !="" {

		//检测手机号是否已经存在
		var validUser sys.SysShopUser
		//查询大B下 + 新手机号
		e.Orm.Model(&sys.SysShopUser{}).Select("user_id").Where("phone = ? ",req.Phone).Limit(1).Find(&validUser)

		if validUser.UserId > 0 {
			if validUser.UserId != parentShopRow.UserId {
				e.Error(500, errors.New("手机号已经存在"), "手机号已经存在")
				return
			}
		}
		//商城的手机号也要保持唯一
		var userShop models.Shop
		e.Orm.Model(&userShop).Select("user_id").Where("phone = ? ",req.Phone).Limit(1).Find(&userShop)
		if userShop.UserId > 0 {
			if userShop.UserId != parentShopRow.UserId {
				e.Error(500, errors.New("手机号已经存在"), "手机号已经存在")
				return
			}
		}

	}
	//用户名发生了变化
	if parentShopRow.UserName != req.UserName {

		var validUser sys.SysShopUser

		e.Orm.Model(&sys.SysShopUser{}).Select("user_id").Where("username = ? ",req.UserName).Limit(1).Find(&validUser)
		if validUser.UserId > 0 {
			if validUser.UserId != parentShopRow.UserId {
				e.Error(500, errors.New("用户名已经存在"), "用户名已经存在")
				return
			}
		}
	}

	var userSymanObject sys.SysUser
	e.Orm.Model(&sys.SysUser{}).Select("user_id").Where("phone = ? and enable = ?",req.SalesmanPhone,true).Limit(1).Find(&userSymanObject)
	if userSymanObject.UserId > 0 {
		req.Salesman = userSymanObject.UserId
	}
	err = s.Update(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("修改失败,%s", err.Error()))
        return
	}
	e.OK( req.GetId(), "修改成功")
}

// Delete 删除Shop
// @Summary 删除Shop
// @Description 删除Shop
// @Tags Shop
// @Param data body dto.ShopDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/shop [delete]
// @Security Bearer
func (e Shop) Delete(c *gin.Context) {
    s := service.Shop{}
    req := dto.ShopDeleteReq{}
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

	// req.SetUpdateBy(user.GetUserId(c))
	p := actions.GetPermissionFromContext(c)

	err = s.Remove(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("删除Shop失败，\r\n失败信息 %s", err.Error()))
        return
	}
	e.OK( req.GetId(), "删除成功")
}

func (e Shop)Grade(c *gin.Context)  {
	req := dto.ShopGradeReq{}
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
	var count int64
	var object models.Shop
	e.Orm.Model(&models.Shop{}).Scopes(actions.PermissionSysUser(object.TableName(), userDto)).Select("id").Where("id = ? ",req.ShopId).First(&object).Count(&count)
	if count == 0 {
		e.Error(500, errors.New("客户不存在"), "客户不存在")
		return
	}

	e.Orm.Model(&models.Shop{}).Where("id = ?",object.Id).Updates(map[string]interface{}{
		"grade_id":req.GradeId,
		"update_by":user.GetUserId(c),
	})
	e.OK("","操作成功")
	return
}


func (e Shop)Credit(c *gin.Context)  {
	req := dto.ShopCreditReq{}
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
	if req.Value < 0 {
		e.Error(500, nil,"不可为负数")
		return
	}

	var count int64
	var object models.Shop
	e.Orm.Model(&models.Shop{}).Scopes(actions.PermissionSysUser(object.TableName(), userDto)).Select("credit,id,credit_quota").Where("id = ? ",req.ShopId).First(&object).Count(&count)
	if count == 0 {
		e.Error(500, errors.New("客户不存在"), "客户不存在")
		return
	}
	Scene:=""
	updateMap:=map[string]interface{}{
		"update_by":user.GetUserId(c),
	}
	switch req.Mode {
	case global.UserNumberAdd:
		//增加时 是给授信额度 增加
		//同时在授信额度增加
		object.Credit += float64(req.Value)
		Scene = fmt.Sprintf("手动增加%v授信额度和余额",req.Value)
		object.CreditQuota += float64(req.Value)
		updateMap["credit_quota"] = utils.RoundDecimalFlot64(object.CreditQuota)
		updateMap["credit"] = utils.RoundDecimalFlot64(object.Credit)
	case global.UserNumberReduce:
		//减少授信额
		if float64(req.Value) > object.Credit {
			e.Error(500, errors.New("授信余额不足"), "授信余额不足")
			return
		}
		object.CreditQuota -=float64(req.Value)
		object.Credit -=float64(req.Value)
		Scene = fmt.Sprintf("手动减少%v授信额度和余额",req.Value)
		updateMap["credit_quota"] = utils.RoundDecimalFlot64(object.CreditQuota)
		updateMap["credit"] = utils.RoundDecimalFlot64(object.Credit)
	//case global.UserNumberSet:
	//	object.Credit = float64(req.Value)
	//	Scene = fmt.Sprintf("手动设置为%v授信余额",req.Value)
	default:
		e.Error(500, nil,"操作不合法")
		return
	}
	e.Orm.Model(&models.Shop{}).Scopes(actions.PermissionSysUser(object.TableName(), userDto)).Where("id = ?",object.Id).Updates(updateMap)
	row:=models.ShopCreditLog{
		CId: userDto.CId,
		ShopId: req.ShopId,
		Desc: req.Desc,
		Number: req.Value,
		Scene:fmt.Sprintf("管理员[%v] %v",userDto.Username,Scene),
		Action: req.Mode,
		Type: global.ScanAdmin,
	}
	row.CreateBy = user.GetUserId(c)
	e.Orm.Create(&row)
	e.OK("","操作成功")
	return
}
func (e Shop)Amount(c *gin.Context)  {
	req := dto.ShopAmountReq{}
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
	if req.Value < 0 {
		e.Error(500, nil,"不可为负数")
		return
	}

	var count int64
	var object models.Shop
	e.Orm.Model(&models.Shop{}).Scopes(actions.PermissionSysUser(object.TableName(), userDto)).Select("balance,id").Where("id = ? ",req.ShopId).First(&object).Count(&count)
	if count == 0 {
		e.Error(500, errors.New("客户不存在"), "客户不存在")
		return
	}
	Scene:=""
	switch req.Mode {
	case global.UserNumberAdd:
		object.Balance += req.Value
		Scene = fmt.Sprintf("手动增加%v元",req.Value)
	case global.UserNumberReduce:
		if req.Value > object.Balance {
			e.Error(500, errors.New("数值非法"), "数值非法")
			return
		}
		object.Balance -=req.Value
		Scene = fmt.Sprintf("手动减少%v元",req.Value)
	case global.UserNumberSet:
		object.Balance = req.Value
		Scene = fmt.Sprintf("手动设置为%v元",req.Value)
	default:
		e.Error(500, nil,"操作不合法")
		return
	}
	e.Orm.Model(&models.Shop{}).Scopes(actions.PermissionSysUser(object.TableName(), userDto)).Where("id = ?",object.Id).Updates(map[string]interface{}{
		"balance":utils.RoundDecimalFlot64(object.Balance),
		"update_by":user.GetUserId(c),
	})
	row:=models.ShopBalanceLog{
		CId: userDto.CId,
		ShopId: req.ShopId,
		Desc: req.Desc,
		Money: req.Value,
		Scene:fmt.Sprintf("管理员[%v]%v",userDto.Username,Scene),
		Action: req.Mode,
		Type: global.ScanAdmin,
	}
	row.CreateBy = user.GetUserId(c)
	e.Orm.Create(&row)
	e.OK("","操作成功")
	return
}
func (e Shop)Integral(c *gin.Context)  {
	req := dto.ShopIntegralReq{}
	err := e.MakeContext(c).
		Bind(&req).
		MakeOrm().
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	if req.Value < 0 {
		e.Error(500, nil,"不可为负数")
		return
	}
	userDto, err := customUser.GetUserDto(e.Orm, c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}
	var count int64
	var object models.Shop
	e.Orm.Model(&models.Shop{}).Scopes(actions.PermissionSysUser(object.TableName(), userDto)).Select("integral,id").Where("id = ?",req.ShopId).First(&object).Count(&count)
	if count == 0 {
		e.Error(500, errors.New("客户不存在"), "客户不存在")
		return
	}
	Scene:=""
	switch req.Mode {
	case global.UserNumberAdd:
		object.Integral += req.Value
		Scene = fmt.Sprintf("手动增加%v个积分",req.Value)
	case global.UserNumberReduce:
		if req.Value > object.Integral {
			e.Error(500, errors.New("数值非法"), "数值非法")
			return
		}
		object.Integral -=req.Value
		Scene = fmt.Sprintf("手动减少%v个积分",req.Value)
	case global.UserNumberSet:
		object.Integral = req.Value
		Scene = fmt.Sprintf("手动积分设置为%v",req.Value)
	default:
		e.Error(500, nil,"操作不合法")
		return
	}
	e.Orm.Model(&models.Shop{}).Scopes(actions.PermissionSysUser(object.TableName(), userDto)).Where("id = ?",object.Id).Updates(map[string]interface{}{
		"integral":object.Integral,
		"update_by":user.GetUserId(c),
	})
	row:=models.ShopIntegralLog{
		CId: userDto.CId,
		ShopId: req.ShopId,
		Desc: req.Desc,
		Number: req.Value,
		Scene:fmt.Sprintf("管理员[%v]%v",userDto.Username,Scene),
		Action: req.Mode,
		Type: global.ScanAdmin,
	}
	row.CreateBy = user.GetUserId(c)
	e.Orm.Create(&row)
	e.OK("","操作成功")
	return

}

func (e Shop) cacheChildren(pid int,cache []interface{},lists []models2.ChinaData) (data []interface{},hs bool)  {

	filterList:=make([]models2.ChinaData,0)

	for _,row:=range lists{
		if row.Pid == pid {
			filterList = append(filterList,row)
		}
	}
	if len(filterList) == 0 {

		return cache,false
	}

	for _,dat:=range filterList{

		CascaderRow:=dto.CascaderRow{
			Code: dat.Id,
			Name: dat.Name,
		}
		cacheChildren:=make([]interface{},0)
		Children,valid:=e.cacheChildren(dat.Id,cacheChildren,lists)
		if valid{
			CascaderRow.Children = Children
		}
		cache = append(cache,CascaderRow)
	}

	return cache,true
}
func (e Shop) Cascader(c *gin.Context) {
	err := e.MakeContext(c).

		MakeOrm().
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	var lists []models2.ChinaData
	e.Orm.Model(&models2.ChinaData{}).Find(&lists)



	cacheChildren:=make([]interface{},0)
	result,_ :=e.cacheChildren(0,cacheChildren,lists)
	e.OK(result,"")
	return
}
func (e Shop) UpPass(c *gin.Context) {
	req :=dto.UpPass{}
	err := e.MakeContext(c).
		Bind(&req).
		MakeOrm().
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	if req.Pass == ""{
		e.Error(500, nil, "请输入密码")
		return
	}
	userDto, err := customUser.GetUserDto(e.Orm, c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}
	//更新的是小B的用户
	var sysDto sys.SysShopUser
	e.Orm.Model(&sys.SysShopUser{}).Scopes(actions.PermissionSysUser(sysDto.TableName(), userDto)).Where("user_id = ? ",req.Id).Limit(1).Find(&sysDto)

	if sysDto.UserId == 0 {
		e.Error(500,nil,"用户不存在")
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Pass), bcrypt.DefaultCost)

	if err!=nil{
		e.Error(500,err,"密码更新失败")
		return
	}
	e.Orm.Model(&sys.SysShopUser{}).Scopes(actions.PermissionSysUser(sysDto.TableName(), userDto)).Where("user_id = ? ",req.Id).Updates(map[string]interface{}{
		"password":string(hash),
	})

	e.OK("","操作成功")
	return
}
func (e Shop) GetLine(c *gin.Context) {

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
	shopId:=c.Param("id")

	var object models.Shop
	e.Orm.Model(&models.Shop{}).Scopes(actions.PermissionSysUser(object.TableName(), userDto)).Where("id = ? and enable = ?",shopId,true).Limit(1).Find(&object)
	if object.Id == 0 {
		e.Error(500,nil,"数据不存在")
		return
	}

	result :=map[string]interface{}{
		"username":object.UserName,
		"address":object.Address,
		"line":"",
		"grade":"",
		"driver":"",
	}
	var lineObject models2.Line
	e.Orm.Model(&models2.Line{}).Scopes(actions.PermissionSysUser(object.TableName(), userDto)).Where("id = ? and enable = ?",object.LineId,true).Limit(1).Find(&lineObject)
	if lineObject.Id  > 0 {
		result["line"] = lineObject.Name
	}
	var driverObject models2.Driver
	e.Orm.Model(&models2.Driver{}).Scopes(actions.PermissionSysUser(object.TableName(), userDto)).Where("id = ? and enable = ?",lineObject.DriverId,true).Limit(1).Find(&driverObject)
	if driverObject.Id  > 0 {
		result["driver_name"] = driverObject.Name
		result["driver_phone"] = driverObject.Phone
	}
	if object.GradeId  > 0 {
		var gradeVip models2.GradeVip
		e.Orm.Model(&models2.GradeVip{}).Scopes(actions.PermissionSysUser(object.TableName(), userDto)).Where("id = ? and enable = ?",object.GradeId,true).Limit(1).Find(&gradeVip)
		if gradeVip.Id  > 0 {
			result["grade"] = gradeVip.Name
		}

	}

	e.OK(result,"操作成功")
	return

}