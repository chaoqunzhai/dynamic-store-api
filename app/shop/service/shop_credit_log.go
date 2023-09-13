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

type ShopCreditLog struct {
	service.Service
}

// GetPage 获取ShopBalanceLog列表
func (e *ShopCreditLog) GetPage(c *dto.ShopCreditLogGetPageReq, p *actions.DataPermission, list *[]models.ShopCreditLog, count *int64) error {
	var err error
	var data models.ShopCreditLog

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).Order("created_at desc").
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("ShopBalanceLogService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取ShopBalanceLog对象
func (e *ShopCreditLog) Get(d *dto.ShopCreditLogGetReq, p *actions.DataPermission, model *models.ShopCreditLog) error {
	var data models.ShopCreditLog

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetShopBalanceLog error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建ShopBalanceLog对象
func (e *ShopCreditLog) Insert(c *dto.ShopCreditLogInsertReq) error {
    var err error
    var data models.ShopCreditLog
    c.Generate(&data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("ShopBalanceLogService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改ShopBalanceLog对象 ShopBalanceLogUpdateReq
func (e *ShopCreditLog) Update(c *dto.ShopCreditLogUpdateReq, p *actions.DataPermission) error {
    var err error
    var data = models.ShopCreditLog{}
    e.Orm.Scopes(
            actions.Permission(data.TableName(), p),
        ).First(&data, c.GetId())
    c.Generate(&data)

    db := e.Orm.Save(&data)
    if err = db.Error; err != nil {
        e.Log.Errorf("ShopBalanceLogService Save error:%s \r\n", err)
        return err
    }
    if db.RowsAffected == 0 {
        return errors.New("无权更新该数据")
    }
    return nil
}

// Remove 删除ShopBalanceLog ShopBalanceLogDeleteReq
func (e *ShopCreditLog) Remove(d *dto.ShopCreditLogDeleteReq, p *actions.DataPermission) error {
	var data models.ShopCreditLog

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
        e.Log.Errorf("Service RemoveShopBalanceLog error:%s \r\n", err)
        return err
    }
    if db.RowsAffected == 0 {
        return errors.New("无权删除该数据")
    }
	return nil
}
