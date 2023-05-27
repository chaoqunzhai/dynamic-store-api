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

type Shop struct {
	service.Service
}

// GetPage 获取Shop列表
func (e *Shop) GetPage(c *dto.ShopGetPageReq, p *actions.DataPermission, list *[]models.Shop, count *int64) error {
	var err error
	var data models.Shop

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("ShopService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取Shop对象
func (e *Shop) Get(d *dto.ShopGetReq, p *actions.DataPermission, model *models.Shop) error {
	var data models.Shop

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetShop error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建Shop对象
func (e *Shop) Insert(c *dto.ShopInsertReq) error {
    var err error
    var data models.Shop
    c.Generate(&data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("ShopService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改Shop对象
func (e *Shop) Update(c *dto.ShopUpdateReq, p *actions.DataPermission) error {
    var err error
    var data = models.Shop{}
    e.Orm.Scopes(
            actions.Permission(data.TableName(), p),
        ).First(&data, c.GetId())
    c.Generate(&data)

    db := e.Orm.Save(&data)
    if err = db.Error; err != nil {
        e.Log.Errorf("ShopService Save error:%s \r\n", err)
        return err
    }
    if db.RowsAffected == 0 {
        return errors.New("无权更新该数据")
    }
    return nil
}

// Remove 删除Shop
func (e *Shop) Remove(d *dto.ShopDeleteReq, p *actions.DataPermission) error {
	var data models.Shop

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
        e.Log.Errorf("Service RemoveShop error:%s \r\n", err)
        return err
    }
    if db.RowsAffected == 0 {
        return errors.New("无权删除该数据")
    }
	return nil
}
