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

type GoodsBrand struct {
	service.Service
}

// GetPage 获取GoodsBrand列表
func (e *GoodsBrand) GetPage(c *dto.GoodsBrandGetPageReq, p *actions.DataPermission, list *[]models.GoodsBrand, count *int64) error {
	var err error
	var data models.GoodsBrand

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).Order(global.OrderLayerKey).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("GoodsBrandService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取GoodsBrand对象
func (e *GoodsBrand) Get(d *dto.GoodsBrandGetReq, p *actions.DataPermission, model *models.GoodsBrand) error {
	var data models.GoodsBrand

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetGoodsBrand error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建GoodsBrand对象
func (e *GoodsBrand) Insert(c *dto.GoodsBrandInsertReq) error {
    var err error
    var data models.GoodsBrand
    c.Generate(&data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("GoodsBrandService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改GoodsBrand对象
func (e *GoodsBrand) Update(c *dto.GoodsBrandUpdateReq, p *actions.DataPermission) error {
    var err error
    var data = models.GoodsBrand{}
    e.Orm.Scopes(
            actions.Permission(data.TableName(), p),
        ).First(&data, c.GetId())
    c.Generate(&data)

    db := e.Orm.Save(&data)
    if err = db.Error; err != nil {
        e.Log.Errorf("GoodsBrandService Save error:%s \r\n", err)
        return err
    }
    if db.RowsAffected == 0 {
        return errors.New("无权更新该数据")
    }
    return nil
}

// Remove 删除GoodsBrand
func (e *GoodsBrand) Remove(d *dto.GoodsBrandDeleteReq, p *actions.DataPermission) error {
	var data models.GoodsBrand

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
        e.Log.Errorf("Service RemoveGoodsBrand error:%s \r\n", err)
        return err
    }
    if db.RowsAffected == 0 {
        return errors.New("无权删除该数据")
    }
	return nil
}
