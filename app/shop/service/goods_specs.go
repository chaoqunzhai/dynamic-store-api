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

type GoodsSpecs struct {
	service.Service
}

// GetPage 获取GoodsSpecs列表
func (e *GoodsSpecs) GetPage(c *dto.GoodsSpecsGetPageReq, p *actions.DataPermission, list *[]models.GoodsSpecs, count *int64) error {
	var err error
	var data models.GoodsSpecs

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("GoodsSpecsService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取GoodsSpecs对象
func (e *GoodsSpecs) Get(d *dto.GoodsSpecsGetReq, p *actions.DataPermission, model *models.GoodsSpecs) error {
	var data models.GoodsSpecs

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetGoodsSpecs error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建GoodsSpecs对象
func (e *GoodsSpecs) Insert(c *dto.GoodsSpecsInsertReq) error {
    var err error
    var data models.GoodsSpecs
    c.Generate(&data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("GoodsSpecsService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改GoodsSpecs对象
func (e *GoodsSpecs) Update(c *dto.GoodsSpecsUpdateReq, p *actions.DataPermission) error {
    var err error
    var data = models.GoodsSpecs{}
    e.Orm.Scopes(
            actions.Permission(data.TableName(), p),
        ).First(&data, c.GetId())
    c.Generate(&data)

    db := e.Orm.Save(&data)
    if err = db.Error; err != nil {
        e.Log.Errorf("GoodsSpecsService Save error:%s \r\n", err)
        return err
    }
    if db.RowsAffected == 0 {
        return errors.New("无权更新该数据")
    }
    return nil
}

// Remove 删除GoodsSpecs
func (e *GoodsSpecs) Remove(d *dto.GoodsSpecsDeleteReq, p *actions.DataPermission) error {
	var data models.GoodsSpecs

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
        e.Log.Errorf("Service RemoveGoodsSpecs error:%s \r\n", err)
        return err
    }
    if db.RowsAffected == 0 {
        return errors.New("无权删除该数据")
    }
	return nil
}
