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

type ShopOrderBindRecord struct {
	service.Service
}

// GetPage 获取ShopOrderBindRecord列表
func (e *ShopOrderBindRecord) GetPage(c *dto.ShopOrderBindRecordGetPageReq, p *actions.DataPermission, list *[]models.ShopOrderBindRecord, count *int64) error {
	var err error
	var data models.ShopOrderBindRecord

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("ShopOrderBindRecordService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取ShopOrderBindRecord对象
func (e *ShopOrderBindRecord) Get(d *dto.ShopOrderBindRecordGetReq, p *actions.DataPermission, model *models.ShopOrderBindRecord) error {
	var data models.ShopOrderBindRecord

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetShopOrderBindRecord error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建ShopOrderBindRecord对象
func (e *ShopOrderBindRecord) Insert(c *dto.ShopOrderBindRecordInsertReq) error {
    var err error
    var data models.ShopOrderBindRecord
    c.Generate(&data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("ShopOrderBindRecordService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改ShopOrderBindRecord对象
func (e *ShopOrderBindRecord) Update(c *dto.ShopOrderBindRecordUpdateReq, p *actions.DataPermission) error {
    var err error
    var data = models.ShopOrderBindRecord{}
    e.Orm.Scopes(
            actions.Permission(data.TableName(), p),
        ).First(&data, c.GetId())
    c.Generate(&data)

    db := e.Orm.Save(&data)
    if err = db.Error; err != nil {
        e.Log.Errorf("ShopOrderBindRecordService Save error:%s \r\n", err)
        return err
    }
    if db.RowsAffected == 0 {
        return errors.New("无权更新该数据")
    }
    return nil
}

// Remove 删除ShopOrderBindRecord
func (e *ShopOrderBindRecord) Remove(d *dto.ShopOrderBindRecordDeleteReq, p *actions.DataPermission) error {
	var data models.ShopOrderBindRecord

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
        e.Log.Errorf("Service RemoveShopOrderBindRecord error:%s \r\n", err)
        return err
    }
    if db.RowsAffected == 0 {
        return errors.New("无权删除该数据")
    }
	return nil
}
