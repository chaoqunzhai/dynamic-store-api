package service

import (
	"errors"

    "github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/gorm"

	"go-admin/app/shop/models"
	"go-admin/app/shop/service/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
)

type ShopIntegralLog struct {
	service.Service
}

// GetPage 获取ShopIntegralLog列表
func (e *ShopIntegralLog) GetPage(c *dto.ShopIntegralLogGetPageReq, p *actions.DataPermission, list *[]models.ShopIntegralLog, count *int64) error {
	var err error
	var data models.ShopIntegralLog

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("ShopIntegralLogService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取ShopIntegralLog对象
func (e *ShopIntegralLog) Get(d *dto.ShopIntegralLogGetReq, p *actions.DataPermission, model *models.ShopIntegralLog) error {
	var data models.ShopIntegralLog

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetShopIntegralLog error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建ShopIntegralLog对象
func (e *ShopIntegralLog) Insert(c *dto.ShopIntegralLogInsertReq) error {
    var err error
    var data models.ShopIntegralLog
    c.Generate(&data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("ShopIntegralLogService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改ShopIntegralLog对象
func (e *ShopIntegralLog) Update(c *dto.ShopIntegralLogUpdateReq, p *actions.DataPermission) error {
    var err error
    var data = models.ShopIntegralLog{}
    e.Orm.Scopes(
            actions.Permission(data.TableName(), p),
        ).First(&data, c.GetId())
    c.Generate(&data)

    db := e.Orm.Save(&data)
    if err = db.Error; err != nil {
        e.Log.Errorf("ShopIntegralLogService Save error:%s \r\n", err)
        return err
    }
    if db.RowsAffected == 0 {
        return errors.New("无权更新该数据")
    }
    return nil
}

// Remove 删除ShopIntegralLog
func (e *ShopIntegralLog) Remove(d *dto.ShopIntegralLogDeleteReq, p *actions.DataPermission) error {
	var data models.ShopIntegralLog

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
        e.Log.Errorf("Service RemoveShopIntegralLog error:%s \r\n", err)
        return err
    }
    if db.RowsAffected == 0 {
        return errors.New("无权删除该数据")
    }
	return nil
}
