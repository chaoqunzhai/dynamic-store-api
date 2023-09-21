package service

import (
	"errors"
	"fmt"
	"github.com/go-admin-team/go-admin-core/sdk/service"
	"github.com/google/uuid"
	"go-admin/app/company/models"
	"go-admin/app/company/service/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
	models2 "go-admin/common/models"
	"go-admin/common/utils"
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
func (e *Orders) CalculateTime(day int) (t models2.XTime) {
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
			response.CycleTime = e.CalculateTime(d.GiveDay)
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
				response.CycleTime = e.CalculateTime(w.GiveDay)
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
	var data models.Orders
	whereSQL := fmt.Sprintf("")
	if c.CId > 0 {
		whereSQL = fmt.Sprintf("c_id = %v", c.CId)
	}
	err = e.Orm.Table(tableName).Where(whereSQL).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(tableName), p),
		).Order(global.OrderTimeKey).
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
