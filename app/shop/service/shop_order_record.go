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

type ShopOrderRecord struct {
	service.Service
}

// GetPage 获取ShopOrderRecord列表
func (e *ShopOrderRecord) GetPage(c *dto.ShopOrderRecordGetPageReq, p *actions.DataPermission, list *[]models.ShopOrderRecord, count *int64) error {
	var err error
	var data models.ShopOrderRecord

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("ShopOrderRecordService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取ShopOrderRecord对象
func (e *ShopOrderRecord) Get(d *dto.ShopOrderRecordGetReq, p *actions.DataPermission, model *models.ShopOrderRecord) error {
	var data models.ShopOrderRecord

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetShopOrderRecord error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建ShopOrderRecord对象
func (e *ShopOrderRecord) Insert(c *dto.ShopOrderRecordInsertReq) error {
    var err error
    var data models.ShopOrderRecord
    c.Generate(&data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("ShopOrderRecordService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改ShopOrderRecord对象
func (e *ShopOrderRecord) Update(c *dto.ShopOrderRecordUpdateReq, p *actions.DataPermission) error {
    var err error
    var data = models.ShopOrderRecord{}
    e.Orm.Scopes(
            actions.Permission(data.TableName(), p),
        ).First(&data, c.GetId())
    c.Generate(&data)

    db := e.Orm.Save(&data)
    if err = db.Error; err != nil {
        e.Log.Errorf("ShopOrderRecordService Save error:%s \r\n", err)
        return err
    }
    if db.RowsAffected == 0 {
        return errors.New("无权更新该数据")
    }
    return nil
}

// Remove 删除ShopOrderRecord
func (e *ShopOrderRecord) Remove(d *dto.ShopOrderRecordDeleteReq, p *actions.DataPermission) error {
	var data models.ShopOrderRecord

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
        e.Log.Errorf("Service RemoveShopOrderRecord error:%s \r\n", err)
        return err
    }
    if db.RowsAffected == 0 {
        return errors.New("无权删除该数据")
    }
	return nil
}
