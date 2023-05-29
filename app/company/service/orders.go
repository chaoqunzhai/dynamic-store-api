package service

import (
	"errors"

	"github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/gorm"

	"go-admin/app/company/models"
	"go-admin/app/company/service/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
)

type Orders struct {
	service.Service
}

// GetPage 获取Orders列表
func (e *Orders) GetPage(tableName string, c *dto.OrdersGetPageReq, p *actions.DataPermission, list *[]models.Orders, count *int64) error {
	var err error
	var data models.Orders

	err = e.Orm.Table(tableName).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
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
