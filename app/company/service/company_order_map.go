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

type CompanyOrderMap struct {
	service.Service
}

// GetPage 获取CompanyOrderMap列表
func (e *CompanyOrderMap) GetPage(c *dto.CompanyOrderMapGetPageReq, p *actions.DataPermission, list *[]models.CompanyOrderMap, count *int64) error {
	var err error
	var data models.CompanyOrderMap

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).Order(global.OrderLayerKey).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("CompanyOrderMapService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取CompanyOrderMap对象
func (e *CompanyOrderMap) Get(d *dto.CompanyOrderMapGetReq, p *actions.DataPermission, model *models.CompanyOrderMap) error {
	var data models.CompanyOrderMap

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetCompanyOrderMap error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建CompanyOrderMap对象
func (e *CompanyOrderMap) Insert(c *dto.CompanyOrderMapInsertReq) error {
    var err error
    var data models.CompanyOrderMap
    c.Generate(&data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("CompanyOrderMapService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改CompanyOrderMap对象
func (e *CompanyOrderMap) Update(c *dto.CompanyOrderMapUpdateReq, p *actions.DataPermission) error {
    var err error
    var data = models.CompanyOrderMap{}
    e.Orm.Scopes(
            actions.Permission(data.TableName(), p),
        ).First(&data, c.GetId())
    c.Generate(&data)

    db := e.Orm.Save(&data)
    if err = db.Error; err != nil {
        e.Log.Errorf("CompanyOrderMapService Save error:%s \r\n", err)
        return err
    }
    if db.RowsAffected == 0 {
        return errors.New("无权更新该数据")
    }
    return nil
}

// Remove 删除CompanyOrderMap
func (e *CompanyOrderMap) Remove(d *dto.CompanyOrderMapDeleteReq, p *actions.DataPermission) error {
	var data models.CompanyOrderMap

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
        e.Log.Errorf("Service RemoveCompanyOrderMap error:%s \r\n", err)
        return err
    }
    if db.RowsAffected == 0 {
        return errors.New("无权删除该数据")
    }
	return nil
}
