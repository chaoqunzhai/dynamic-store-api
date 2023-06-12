package service

import (
	"errors"
	"fmt"
	"github.com/go-admin-team/go-admin-core/sdk/service"
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

func (e *Orders) CalculateTime(day int) (delivery_time models2.XTime) {
	//选择的天数，计算出配送周期

	newTime := time.Now().AddDate(0, 0, day)
	x := models2.XTime{}
	x.Time = time.Date(newTime.Year(), newTime.Month(), newTime.Day(), 0, 0, 0, 0, time.Local)
	return x
}
func (e *Orders) ValidTimeConf(cid int) (v bool, uid int, delivery_time models2.XTime, deliver_str string) {
	var data models.CycleTimeConf
	var count int64
	e.Orm.Model(&data).Where("c_id = ?", cid).Count(&count)
	if count == 0 {
		return false, 0, models2.XTime{}, ""
	}
	//todo:获取当前时间是周几,
	thisWeek := utils.HasWeekNumber()

	//todo:查询是否配置了每天的时间,查询出大B配置的开始-结束 时间区间的值
	var timeConfDay []models.CycleTimeConf
	e.Orm.Model(&models.CycleTimeConf{}).Where("c_id = ? and type = ?", cid, global.CyCleTimeDay).Find(&timeConfDay)
	dayId := 0
	deliverDay := 0
	deliverStr := ""
	for _, d := range timeConfDay {
		//匹配到时间区间直接返回,时间配置的ID即可
		if utils.TimeCheckRange(d.StartTime, d.EndTime) {
			dayId = d.Id
			deliverStr = d.GiveTime
			deliverDay = d.GiveDay
		}
	}
	//根据这个N+(deliverDay)=订单的配送时间
	if dayId > 0 {

		return true, dayId, e.CalculateTime(deliverDay), deliverStr
	}

	//todo:检测是否配置了每周的时间
	var timeConfWeek []models.CycleTimeConf
	e.Orm.Model(&models.CycleTimeConf{}).Where("c_id = ? and type = ?", cid, global.CyCleTimeWeek).Find(&timeConfWeek)
	weekId := 0
	weekDeliverDay := 0
	weekDeliverStr := ""
	for _, w := range timeConfWeek {
		//当前周在配置的周期中
		if thisWeek >= w.StartWeek && thisWeek <= w.EndWeek {
			if utils.TimeCheckRange(w.StartTime, w.EndTime) {
				weekId = w.Id
				weekDeliverDay = w.GiveDay
				weekDeliverStr = w.GiveTime
			}
		}
	}
	//当前周几换算为时间 + (deliverDay)=订单的配送时间
	if weekId > 0 {
		return true, weekId, e.CalculateTime(weekDeliverDay), weekDeliverStr
	}

	return false, 0, models2.XTime{}, ""
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
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(tableName), p),
		).
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
