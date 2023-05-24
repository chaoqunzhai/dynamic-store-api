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

type ShopRechargeLog struct {
	service.Service
}

// GetPage 获取ShopRechargeLog列表
func (e *ShopRechargeLog) GetPage(c *dto.ShopRechargeLogGetPageReq, p *actions.DataPermission, list *[]models.ShopRechargeLog, count *int64) error {
	var err error
	var data models.ShopRechargeLog

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("ShopRechargeLogService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取ShopRechargeLog对象
func (e *ShopRechargeLog) Get(d *dto.ShopRechargeLogGetReq, p *actions.DataPermission, model *models.ShopRechargeLog) error {
	var data models.ShopRechargeLog

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetShopRechargeLog error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建ShopRechargeLog对象
func (e *ShopRechargeLog) Insert(c *dto.ShopRechargeLogInsertReq) error {
    var err error
    var data models.ShopRechargeLog
    c.Generate(&data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("ShopRechargeLogService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改ShopRechargeLog对象
func (e *ShopRechargeLog) Update(c *dto.ShopRechargeLogUpdateReq, p *actions.DataPermission) error {
    var err error
    var data = models.ShopRechargeLog{}
    e.Orm.Scopes(
            actions.Permission(data.TableName(), p),
        ).First(&data, c.GetId())
    c.Generate(&data)

    db := e.Orm.Save(&data)
    if err = db.Error; err != nil {
        e.Log.Errorf("ShopRechargeLogService Save error:%s \r\n", err)
        return err
    }
    if db.RowsAffected == 0 {
        return errors.New("无权更新该数据")
    }
    return nil
}

// Remove 删除ShopRechargeLog
func (e *ShopRechargeLog) Remove(d *dto.ShopRechargeLogDeleteReq, p *actions.DataPermission) error {
	var data models.ShopRechargeLog

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
        e.Log.Errorf("Service RemoveShopRechargeLog error:%s \r\n", err)
        return err
    }
    if db.RowsAffected == 0 {
        return errors.New("无权删除该数据")
    }
	return nil
}
