package service

import (
	"errors"
	"fmt"
	"github.com/go-admin-team/go-admin-core/sdk/service"
	"github.com/google/uuid"
	sys "go-admin/app/admin/models"
	"go-admin/app/company/models"
	"go-admin/app/company/service/dto"
	models3 "go-admin/cmd/migrate/migration/models"
	"go-admin/common/actions"
	"go-admin/common/business"
	cDto "go-admin/common/dto"
	models2 "go-admin/common/models"
	"go-admin/common/utils"
	"go-admin/config"
	"go-admin/global"
	"gorm.io/gorm"
	"time"
)

type Orders struct {
	service.Service
}

type ValetOrderGoodsRow struct {
	GoodsId int
	SpecsId int
	Number int
}
type TimeConfResponse struct {
	Valid    bool
	ObjectId int
	CycleTime models2.XTime
	CycleStr string
	RandUid  string
	StartTime models2.XTime
	EndTime models2.XTime
}

//生成随机码

func RandomCode() string {
	//9位
	guid, _ := uuid.NewRandom()
	code := fmt.Sprintf("%v", guid.ID())
	if len(code) >= 9 {
		return code[:9]
	}else {
		return code
	}
}
func CalculateTime(day int) (t models2.XTime) {
	//选择的天数，计算出配送周期

	newTime := time.Now().AddDate(0, 0, day)
	x := models2.XTime{}
	x.Time = time.Date(newTime.Year(), newTime.Month(), newTime.Day(), 0, 0, 0, 0, time.Local)
	return x
}
func CalculateSendTime(day int,addTime time.Time) (t models2.XTime) {
	//选择的天数，计算出配送周期

	newTime := addTime.AddDate(0, 0, day)
	x := models2.XTime{}
	x.Time = time.Date(newTime.Year(), newTime.Month(), newTime.Day(), 0, 0, 0, 0, time.Local)
	return x
}
//23:30 to 2023年06月12日23:30:30
func (e *Orders) MakeTime(value string) (x models2.XTime) {
	x = models2.XTime{}
	timeDemo:=fmt.Sprintf("%v %v",
		time.Now().Format("2006-01-02"),
		value)
	tt, err := time.ParseInLocation("2006-01-02 15:04", timeDemo, time.Local)
	if err!=nil{
		x.Time = time.Now()
		return
	}
	x.Time = tt
	return
}
func (e *Orders) ValidTimeConf(cid int) (response *TimeConfResponse) {
	var data models.CycleTimeConf
	var count int64
	response = &TimeConfResponse{}
	e.Orm.Model(&data).Where("c_id = ?", cid).Count(&count)
	if count == 0 {
		response.Valid = false

		return response
	}
	//todo:获取当前时间是周几,
	thisWeek := utils.HasWeekNumber()

	//todo:查询是否配置了每天的时间,查询出大B配置的开始-结束 时间区间的值
	var timeConfDay []models.CycleTimeConf
	e.Orm.Model(&models.CycleTimeConf{}).Where("c_id = ? and type = ?", cid, global.CyCleTimeDay).Find(&timeConfDay)
	for _, d := range timeConfDay {
		//匹配到时间区间直接返回,时间配置的ID即可
		//根据这个N+(deliverDay)=订单的配送时间
		if utils.TimeCheckRange(d.StartTime, d.EndTime) {
			response.Valid = true
			response.RandUid = d.Uid
			response.ObjectId = d.Id
			response.CycleTime = CalculateTime(d.GiveDay)
			response.CycleStr = d.GiveTime
			response.StartTime = e.MakeTime(d.StartTime)
			response.EndTime = e.MakeTime(d.EndTime)
			return

		}
	}

	//todo:检测是否配置了每周的时间
	var timeConfWeek []models.CycleTimeConf
	e.Orm.Model(&models.CycleTimeConf{}).Where("c_id = ? and type = ?", cid, global.CyCleTimeWeek).Find(&timeConfWeek)

	for _, w := range timeConfWeek {
		//当前周在配置的周期中
		if thisWeek >= w.StartWeek && thisWeek <= w.EndWeek {
			if utils.TimeCheckRange(w.StartTime, w.EndTime) {
				response.Valid = true
				response.RandUid = w.Uid
				response.ObjectId = w.Id
				response.CycleTime = CalculateTime(w.GiveDay)
				response.CycleStr = w.GiveTime
				response.StartTime = e.MakeTime(w.StartTime)
				response.EndTime = e.MakeTime(w.EndTime)
				//当前周几换算为时间 + (deliverDay)=订单的配送时间
				return
			}
		}
	}

	return
}

// GetPage 获取Orders列表
func (e *Orders) GetPage(openApprove bool,splitTableRes business.TableRow,countMap *dto.CountOrder,
	c *dto.OrdersGetPageReq, p *actions.DataPermission, list *[]models.Orders, count *int64,OrderRangeTime string) error {
	var err error

	orm :=e.Orm.Table(splitTableRes.OrderTable)
	if c.Status > 0{
		//如果开启了审核, 那查待配送的
		if openApprove {

			switch c.Status {
			case global.OrderStatusWaitSend:
				//如果开启了审核，那查询待配送的时候 就需要 满足 待配送 和审批通过
				orm = orm.Table(splitTableRes.OrderTable).Where("status =  ? and approve_status = ?",global.OrderStatusWaitSend,global.OrderApproveOk)
			case global.OrderApproveOk:
				orm = orm.Table(splitTableRes.OrderTable).Where("approve_status = 0")

				
			default:
				orm = orm.Table(splitTableRes.OrderTable).Where("status =  ?",c.Status)
			}

		}else {
			orm = orm.Table(splitTableRes.OrderTable).Where("status =  ?",c.Status)
		}
	}else {
		orm = orm.Table(splitTableRes.OrderTable)
	}


	if c.CycleType == 2{ //下单周期的查询
		//通过uid 获取到选择条目的

		var OrderCycle models3.OrderCycleCnf
		e.Orm.Table(splitTableRes.OrderCycle).Where("uid = ?",c.Uid).Limit(1).Find(&OrderCycle)

		if OrderCycle.Id > 0{
			orm = orm.Table(splitTableRes.OrderTable).Where("`created_at` >= ? and `created_at` <= ?",OrderCycle.StartTime,OrderCycle.EndTime)
		}
	}

	if OrderRangeTime !=""{//时间范围查询
		orm = orm.Where(OrderRangeTime)
	}

	if c.Verify{//根据过滤 然后在统计一次
		//作废的订单不统计
		orm = orm.Where("status != ? and order_money > 0 ", global.OrderStatusCancel)

	}
	err = orm.Scopes(
			cDto.MakeSplitTableCondition(c.GetNeedSearch(),splitTableRes.OrderTable),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(splitTableRes.OrderTable,p)).Order(global.OrderTimeKey).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	countSql:="SUM(accept_money) as all_accept_money,SUM(order_money) as all_order_money,SUM(number) as number,count(*) as 'count',SUM(coupon_money) as all_coupon_money,SUM(goods_money) as all_goods_money "

	//开启对账功能
	if c.Verify{//根据过滤 然后在统计一次,重复利用了orm对象 所以分开执行
		orm.Select(countSql).Scan(&countMap)
	}

	if err != nil {
		e.Log.Errorf("订单不存在", err)
		return err
	}
	return nil
}

// Get 获取Orders对象
func (e *Orders) Get(tableName string, d *dto.OrdersGetReq, p *actions.DataPermission, model *models.Orders) error {
	var data models.Orders

	err := e.Orm.Table(tableName).
		Scopes(
			actions.Permission(data.TableName(tableName), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetOrders error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建Orders对象
func (e *Orders) Insert(tableName string, c *dto.OrdersInsertReq) error {
	var err error
	var data models.Orders
	c.Generate(&data)
	err = e.Orm.Table(tableName).Create(&data).Error
	if err != nil {
		e.Log.Errorf("OrdersService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改Orders对象
func (e *Orders) Update(tableName string, c *dto.OrdersUpdateReq, p *actions.DataPermission) error {
	var err error
	var data = models.Orders{}
	e.Orm.Table(tableName).Scopes(
		actions.Permission(data.TableName(tableName), p),
	).First(&data, c.GetId())
	c.Generate(&data)

	db := e.Orm.Save(&data)
	if err = db.Error; err != nil {
		e.Log.Errorf("OrdersService Save error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	return nil
}

// Remove 删除Orders
func (e *Orders) Remove(tableName string, d *dto.OrdersDeleteReq, p *actions.DataPermission) error {
	var data models.Orders

	db := e.Orm.Table(tableName).
		Scopes(
			actions.Permission(data.TableName(tableName), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveOrders error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}

func (e *Orders)DetailOrder(orderId string,userDto *sys.SysUser,req dto.DetailReq) (result map[string]interface{},err error)  {
	nowTimeObj :=time.Now()
	var object models.Orders

	splitTableRes := business.GetTableName(userDto.CId, e.Orm)
	orderErr := e.Orm.Table(splitTableRes.OrderTable).Where("order_id = ?",orderId).First(&object).Error
	if orderErr != nil && errors.Is(orderErr, gorm.ErrRecordNotFound) {

		return nil,errors.New("订单不存在")
	}
	if orderErr != nil {
		return nil,errors.New("订单不存在")
	}

	var shopRow models3.Shop
	e.Orm.Model(&models3.Shop{}).Scopes(actions.PermissionSysUser(shopRow.TableName(),userDto)).Where("id = ? ", object.ShopId).Limit(1).Find(&shopRow)
	openApprove,hasApprove:=IsHasOpenApprove(userDto,e.Orm)
	result = map[string]interface{}{
		"order":object.OrderId,
		"order_money":fmt.Sprintf("%v", utils.StringDecimal(object.OrderMoney)),
		"order_goods_money":fmt.Sprintf("%v", utils.StringDecimal(object.GoodsMoney)),
		"coupon_money":fmt.Sprintf("%v", utils.StringDecimal(object.CouponMoney)),
		"delivery_money":fmt.Sprintf("%v", utils.StringDecimal(object.DeliveryMoney)),
		"order_number":object.Number,
		"order_id":       object.Id,
		"created_at":     object.CreatedAt.Format("2006-01-02 15:04:05"),
		"cycle_time":object.CreatedAt.Format("2006-01-02"),
		"delivery_time":     object.DeliveryTime.Format("2006-01-02"),
		"delivery_str":      object.DeliveryStr,
		"pay":            global.GetPayType(object.PayType),
		"pay_int":object.PayType,
		"pay_status_str": global.GetOrderPayStatus(object.PayStatus),
		"pay_status":     object.PayStatus,
		"status_int":object.Status,
		"status":     global.OrderStatus(object.Status),
		"delivery_type":object.DeliveryType,
		"delivery_type_cn":global.GetExpressCn(object.DeliveryType),
		"day":nowTimeObj.Format("2006-01-02"),
		"now":nowTimeObj.Format("2006-01-02 15:04:05"),
		"this_user":userDto.Username,
		"phone":userDto.Phone,
		//https://weapp.dongchuangyun.com/d1#/'
		"url":fmt.Sprintf("%v?siteId=%v",config.ExtConfig.H5Url,userDto.CId),
		"desc":object.Desc,
		"buyer":object.Buyer,
		"all_money_cn":utils.ConvertNumToCny(object.GoodsMoney),
		"order_money_cn":utils.ConvertNumToCny(object.OrderMoney),
		"run_time":"",
		"source_type_int":object.SourceType,
		"accept_msg":object.AcceptMsg,
		"balance":map[string]interface{}{
			"credit":shopRow.Credit,
			"balance":shopRow.Balance,
		},
		"edit_action":object.EditAction,
		"approve_msg":object.ApproveMsg,
		"approve_status":object.ApproveStatus,
		"openApprove":openApprove,
		"hasApprove":hasApprove,
		"ems_id":object.EmsId,
		"coupon_id":object.CouponId,
		"line_id":object.LineId,
	}
	if object.LineId > 0 {
		var lineObj models.Line
		e.Orm.Model(&models.Line{}).Select("name,driver_id,id").Where("id = ? and enable = ? ", object.LineId, true).Limit(1).Find(&lineObj)
		if lineObj.Id > 0 {
			result["line_name"] = lineObj.Name
		}
	}
	if object.CouponId > 0 {
		var couponObj models3.Coupon
		e.Orm.Model(&couponObj).Where("id = ? and c_id = ?",object.CouponId,object.CId).Limit(1).Find(&couponObj)
		result["coupon_name"] = couponObj.Name
	}
	var acceptInt int64
	e.Orm.Model(&models3.OrderAccept{}).Where("c_id = ? and order_id = ?",userDto.CId,object.OrderId).Count(&acceptInt)
	result["accept_count"] = acceptInt
	if object.HelpBy > 0 {
		var user *sys.SysUser
		e.Orm.Model(&user).Select("username").Where("user_id = ?",object.HelpBy).Limit(1).Find(&user)
		result["help_user"] = user.Username
	}
	if shopRow.GradeId > 0 {
		var gradeRow models3.GradeVip
		e.Orm.Model(&models3.GradeVip{}).Select("name,id").Where("id = ? and enable = ?",shopRow.GradeId,true).Limit(1).Find(&gradeRow)
		if gradeRow.Id > 0 {
			result["grade_name"] = gradeRow.Name
		}
	}
	if object.OfflinePayId > 0 {
		var OfflinePay models3.OfflinePay
		e.Orm.Model(&models3.OfflinePay{}).Where("c_id = ? and id = ?",object.CId,object.OfflinePayId).Limit(1).Find(&OfflinePay)
		if OfflinePay.Id > 0 {
			result["offline_pay"] = OfflinePay.Name
		}
	}
	if openApprove {
		//开启了审核
		if object.ApproveStatus == 0 { //还没审核,这个订单就是审核中
			result["status"] = "待审核"
		}
	}
	if object.ApproveStatus == global.OrderApproveReject {
		result["status"] = "已驳回"
	}
	if object.DeliveryType == global.ExpressSelf {

		if object.Status == global.OrderWaitConfirm{
			result["status"] = "配送中"
		}
	}
	if shopRow.Id > 0{
		result["shop_name"] =shopRow.Name
		result["shop_username"] =shopRow.UserName
		result["shop_phone"] =shopRow.Phone
		result["shop_address"] =shopRow.Address
	}

	//如果是同城配送那就获取
	switch object.DeliveryType {
	case global.ExpressSameCity,global.ExpressEms: //同城和物流都是读取客户的地址
		var userAddress models3.DynamicUserAddress
		e.Orm.Model(&models3.DynamicUserAddress{}).Scopes(actions.PermissionSysUser(userAddress.TableName(),userDto)).Where(" id = ?",
			object.AddressId).Limit(1).Find(&userAddress)
		if userAddress.Id > 0{
			result["address"] = map[string]interface{}{
				"address":userAddress.AddressAll(),
				"name":userAddress.Name,
				"phone":userAddress.Mobile,
			}
		}

	case global.ExpressSelf: //自提
		var expressStore models3.CompanyExpressStore
		e.Orm.Model(&models3.CompanyExpressStore{}).Scopes(actions.PermissionSysUser(expressStore.TableName(),userDto)).Select("id,address,name").Where(" id = ?",
			object.AddressId).Limit(1).Find(&expressStore)
		if expressStore.Id > 0{
			result["address"] = map[string]interface{}{
				"name":expressStore.Name,
				"address":expressStore.Address,
			}
		}

	}

	if object.DriverId > 0 {
		var driverCnf models.Driver
		e.Orm.Model(&driverCnf).Scopes(actions.PermissionSysUser(driverCnf.TableName(),userDto)).Where("id = ? ",object.DriverId).Limit(1).Find(&driverCnf)
		if driverCnf.Id > 0 {
			result["driver_name"] = driverCnf.Name
			result["driver_phone"] = driverCnf.Phone
		}
	}


	var orderSpecs []models.OrderSpecs


	e.Orm.Table(splitTableRes.OrderSpecs).Where("order_id = ?", orderId).Find(&orderSpecs)


	specsList := make([]map[string]interface{}, 0)
	//因为编辑需要有库存上限的

	isOpenInventory:=IsOpenInventory(userDto.CId,e.Orm)
	for _, row := range orderSpecs {
		if row.Number == 0 {continue}
		var stock int
		if isOpenInventory{ //去仓库中拿数据
			var Inventory models3.Inventory
			e.Orm.Model(&models3.Inventory{}).Select("stock").Where("c_id = ? and goods_id = ? and spec_id = ?",userDto.CId,row.GoodsId,row.SpecId).Limit(1).Find(&Inventory)
			stock = Inventory.Stock
		}else {
			var goodsSpecs models3.GoodsSpecs
			e.Orm.Model(&goodsSpecs).Select("inventory").Where("c_id = ? and goods_id = ? and id = ?",userDto.CId,row.GoodsId,row.SpecId).Limit(1).Find(&goodsSpecs)
			stock = goodsSpecs.Inventory
		}
		if row.AllMoney == 0 {
			row.AllMoney = utils.RoundDecimalFlot64(row.Money  * float64(row.Number))
		}
		ss := map[string]interface{}{
			"id":         row.Id,
			"name":       row.SpecsName,
			"goods_name":row.GoodsName,
			"created_at": row.CreatedAt.Format("2006-01-02 15:04:05"),
			"specs":      fmt.Sprintf("%v%v", row.Number, row.Unit),
			"unit":row.Unit,
			"money":     utils.StringDecimal(row.Money),
			"number":row.Number,
			"status":row.Status,
			"after_status":row.AfterStatus,
			"edit_action":row.EditAction,
			"stock":stock,//库存量
			"all_money": utils.StringDecimal(row.AllMoney),
		}
		specsList = append(specsList, ss)
	}

	result["run_time"] = time.Since(nowTimeObj)
	result["specs_list"] = specsList


	return result,nil

}