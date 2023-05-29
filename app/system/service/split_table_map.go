package service

import (
	"errors"
	"go-admin/global"

	"github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/gorm"

	"go-admin/app/system/models"
	"go-admin/app/system/service/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
)

type SplitTableMap struct {
	service.Service
}

// GetPage 获取SplitTableMap列表
func (e *SplitTableMap) GetPage(c *dto.SplitTableMapGetPageReq, p *actions.DataPermission, list *[]models.SplitTableMap, count *int64) error {
	var err error
	var data models.SplitTableMap

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).Order(global.OrderLayerKey).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("SplitTableMapService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取SplitTableMap对象
func (e *SplitTableMap) Get(d *dto.SplitTableMapGetReq, p *actions.DataPermission, model *models.SplitTableMap) error {
	var data models.SplitTableMap

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetSplitTableMap error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建SplitTableMap对象
func (e *SplitTableMap) Insert(c *dto.SplitTableMapInsertReq) error {
	var err error
	var data models.SplitTableMap
	c.Generate(&data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("分表创建失败\r\n", err)
		return err
	}

	return nil
}

// Update 修改SplitTableMap对象
func (e *SplitTableMap) Update(c *dto.SplitTableMapUpdateReq, p *actions.DataPermission) error {
	var err error
	var data = models.SplitTableMap{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).First(&data, c.GetId())
	c.Generate(&data)

	db := e.Orm.Save(&data)
	if err = db.Error; err != nil {
		e.Log.Errorf("SplitTableMapService Save error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	return nil
}

// Remove 删除SplitTableMap
func (e *SplitTableMap) Remove(d *dto.SplitTableMapDeleteReq, p *actions.DataPermission) error {
	var data models.SplitTableMap

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveSplitTableMap error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}
