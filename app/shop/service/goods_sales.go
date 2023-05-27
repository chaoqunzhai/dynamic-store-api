package service

import (
	"errors"
	"go-admin/global"

	"github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/gorm"

	"go-admin/app/shop/models"
	"go-admin/app/shop/service/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
)

type GoodsSales struct {
	service.Service
}

// GetPage 获取GoodsSales列表
func (e *GoodsSales) GetPage(c *dto.GoodsSalesGetPageReq, p *actions.DataPermission, list *[]models.GoodsSales, count *int64) error {
	var err error
	var data models.GoodsSales

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).Order(global.OrderLayerKey).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("GoodsSalesService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取GoodsSales对象
func (e *GoodsSales) Get(d *dto.GoodsSalesGetReq, p *actions.DataPermission, model *models.GoodsSales) error {
	var data models.GoodsSales

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetGoodsSales error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建GoodsSales对象
func (e *GoodsSales) Insert(c *dto.GoodsSalesInsertReq) error {
    var err error
    var data models.GoodsSales
    c.Generate(&data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("GoodsSalesService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改GoodsSales对象
func (e *GoodsSales) Update(c *dto.GoodsSalesUpdateReq, p *actions.DataPermission) error {
    var err error
    var data = models.GoodsSales{}
    e.Orm.Scopes(
            actions.Permission(data.TableName(), p),
        ).First(&data, c.GetId())
    c.Generate(&data)

    db := e.Orm.Save(&data)
    if err = db.Error; err != nil {
        e.Log.Errorf("GoodsSalesService Save error:%s \r\n", err)
        return err
    }
    if db.RowsAffected == 0 {
        return errors.New("无权更新该数据")
    }
    return nil
}

// Remove 删除GoodsSales
func (e *GoodsSales) Remove(d *dto.GoodsSalesDeleteReq, p *actions.DataPermission) error {
	var data models.GoodsSales

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
        e.Log.Errorf("Service RemoveGoodsSales error:%s \r\n", err)
        return err
    }
    if db.RowsAffected == 0 {
        return errors.New("无权删除该数据")
    }
	return nil
}
