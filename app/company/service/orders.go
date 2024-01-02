package service

import (
	"errors"
	"fmt"
	"github.com/go-admin-team/go-admin-core/sdk/service"
	"github.com/google/uuid"
	sys "go-admin/app/admin/models"
	"go-admin/app/company/models"
	"go-admin/app/company/service/dto"
	"go-admin/common/actions"
	"go-admin/common/business"
	cDto "go-admin/common/dto"
	models2 "go-admin/common/models"
	models3 "go-admin/cmd/migrate/migration/models"
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

//生成核销码

func DeliveryCode() string {
	//9位

	guid,_ := uuid.NewRandom()
	code:=fmt.Sprintf("%v",guid.ID())

	return code[:9]
}
func CalculateTime(day int) (t models2.XTime) {
	//选择的天数，计算出配送周期

	newTime := time.Now().AddDate(0, 0, day)
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
func (e *Orders) GetPage(tableName string, c *dto.OrdersGetPageReq, p *actions.DataPermission, list *[]models.Orders, count *int64) error {
	var err error


	err = e.Orm.Table(tableName).
		Scopes(
			cDto.MakeSplitTableCondition(c.GetNeedSearch(),tableName),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(tableName,p)).Order(global.OrderTimeKey).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
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

func (e *Orders)DetailOrder(orderId string,userDto *sys.SysUser) (result map[string]interface{},err error)  {
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

	result = map[string]interface{}{
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
		"day":nowTimeObj.Format("2006-01-02"),
		"now":nowTimeObj.Format("2006-01-02 15:04:05"),
		"this_user":userDto.Username,
		//https://weapp.dongchuangyun.com/d1#/'
		"url":fmt.Sprintf("%vd%v#/",config.ExtConfig.H5Url,userDto.CId),
		"desc":object.Desc,
		"buyer":object.Buyer,
		"all_money_cn":utils.ConvertNumToCny(object.GoodsMoney),
		"order_money_cn":utils.ConvertNumToCny(object.OrderMoney),
		"run_time":"",
	}
	if shopRow.Id > 0{
		result["shop_name"] =shopRow.Name
		result["shop_username"] =shopRow.UserName
		result["shop_phone"] =shopRow.Phone
		result["shop_address"] =shopRow.Address
	}

	//如果是同城配送那就获取
	switch object.DeliveryType {
	case global.ExpressLocal:
		var userAddress models3.DynamicUserAddress
		e.Orm.Model(&models3.DynamicUserAddress{}).Scopes(actions.PermissionSysUser(userAddress.TableName(),userDto)).Select("id,address").Where(" id = ?",
			object.AddressId).Limit(1).Find(&userAddress)
		if userAddress.Id > 0{
			result["address"] = map[string]interface{}{
				"address":userAddress.Address,
			}
		}

	case global.ExpressStore:
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

	var driverCnf models.Driver
	e.Orm.Model(&driverCnf).Scopes(actions.PermissionSysUser(driverCnf.TableName(),userDto)).Where("id = ? ",object.DriverId).Limit(1).Find(&driverCnf)
	if driverCnf.Id > 0 {
		result["driver_name"] = driverCnf.Name
		result["driver_phone"] = driverCnf.Phone
	}

	var orderSpecs []models.OrderSpecs


	e.Orm.Table(splitTableRes.OrderSpecs).Where("order_id = ?", orderId).Find(&orderSpecs)

	specsList := make([]map[string]interface{}, 0)
	for _, row := range orderSpecs {
		allMoney := utils.RoundDecimalFlot64(row.Money  * float64(row.Number))
		ss := map[string]interface{}{
			"id":         row.Id,
			"name":       row.SpecsName,
			"goods_name":row.GoodsName,
			"created_at": row.CreatedAt.Format("2006-01-02 15:04:05"),
			"specs":      fmt.Sprintf("%v%v", row.Number, row.Unit),
			"unit":row.Unit,
			"money":     utils.StringDecimal(row.Money),
			"number":row.Number,
			"all_money":fmt.Sprintf("%v", utils.StringDecimal(allMoney)),
		}
		specsList = append(specsList, ss)
	}

	result["run_time"] = time.Since(nowTimeObj)
	result["specs_list"] = specsList


	return result,nil

}