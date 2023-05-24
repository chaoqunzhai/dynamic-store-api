package service

import (
	"errors"

    "github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/gorm"

	"go-admin/app/company/models"
	"go-admin/app/company/service/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
)

type Company struct {
	service.Service
}

// GetPage 获取Company列表
func (e *Company) GetPage(c *dto.CompanyGetPageReq, p *actions.DataPermission, list *[]models.Company, count *int64) error {
	var err error
	var data models.Company

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("CompanyService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取Company对象
func (e *Company) Get(d *dto.CompanyGetReq, p *actions.DataPermission, model *models.Company) error {
	var data models.Company

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetCompany error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建Company对象
func (e *Company) Insert(c *dto.CompanyInsertReq) error {
    var err error
    var data models.Company
    c.Generate(&data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("CompanyService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改Company对象
func (e *Company) Update(c *dto.CompanyUpdateReq, p *actions.DataPermission) error {
    var err error
    var data = models.Company{}
    e.Orm.Scopes(
            actions.Permission(data.TableName(), p),
        ).First(&data, c.GetId())
    c.Generate(&data)

    db := e.Orm.Save(&data)
    if err = db.Error; err != nil {
        e.Log.Errorf("CompanyService Save error:%s \r\n", err)
        return err
    }
    if db.RowsAffected == 0 {
        return errors.New("无权更新该数据")
    }
    return nil
}

// Remove 删除Company
func (e *Company) Remove(d *dto.CompanyDeleteReq, p *actions.DataPermission) error {
	var data models.Company

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
        e.Log.Errorf("Service RemoveCompany error:%s \r\n", err)
        return err
    }
    if db.RowsAffected == 0 {
        return errors.New("无权删除该数据")
    }
	return nil
}
