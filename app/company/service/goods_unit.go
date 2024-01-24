package service

import (
	"errors"
	"go-admin/global"

	"github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/gorm"

	"go-admin/app/company/models"
	"go-admin/app/company/service/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
)

type GoodsUnit struct {
	service.Service
}

// GetPage 获取GoodsUnit列表
func (e *GoodsUnit) GetPage(c *dto.GoodsUnitGetPageReq, p *actions.DataPermission, list *[]models.GoodsUnit, count *int64) error {
	var err error
	var data models.GoodsUnit

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).Order(global.OrderLayerKey).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("GoodsUnitService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取GoodsUnit对象
func (e *GoodsUnit) Get(d *dto.GoodsUnitGetReq, p *actions.DataPermission, model *models.GoodsUnit) error {
	var data models.GoodsUnit

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetGoodsUnit error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建GoodsUnit对象
func (e *GoodsUnit) Insert(c *dto.GoodsUnitInsertReq) error {
    var err error
    var data models.GoodsUnit
    c.Generate(&data)
	data.CId = c.CId
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("GoodsUnitService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改GoodsUnit对象
func (e *GoodsUnit) Update(c *dto.GoodsUnitUpdateReq, p *actions.DataPermission) error {
    var err error
    var data = models.GoodsUnit{}
    e.Orm.Scopes(
            actions.Permission(data.TableName(), p),
        ).First(&data, c.GetId())
    c.Generate(&data)
	data.CId = c.CId
    db := e.Orm.Save(&data)
    if err = db.Error; err != nil {
        e.Log.Errorf("GoodsUnitService Save error:%s \r\n", err)
        return err
    }
    if db.RowsAffected == 0 {
        return errors.New("无权更新该数据")
    }
    return nil
}

// Remove 删除GoodsUnit
func (e *GoodsUnit) Remove(d *dto.GoodsUnitDeleteReq, p *actions.DataPermission) error {
	var data models.GoodsUnit

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())


	if err := db.Error; err != nil {
        e.Log.Errorf("Service RemoveGoodsUnit error:%s \r\n", err)
        return err
    }
    if db.RowsAffected == 0 {
        return errors.New("无权删除该数据")
    }
	return nil
}
