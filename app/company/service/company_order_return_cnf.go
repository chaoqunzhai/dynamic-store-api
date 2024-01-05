package service

import (
	"errors"
	"github.com/go-admin-team/go-admin-core/sdk/service"
	"go-admin/global"
	"gorm.io/gorm"

	"go-admin/app/company/models"
	"go-admin/app/company/service/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
)

type CompanyOrderReturnCnf struct {
	service.Service
}

// GetPage 获取CompanyOrderReturnCnf列表
func (e *CompanyOrderReturnCnf) GetPage(c *dto.CompanyOrderReturnCnfGetPageReq, p *actions.DataPermission, list *[]models.CompanyOrderReturnCnf, count *int64) error {
	var err error
	var data models.CompanyOrderReturnCnf

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).Order(global.OrderLayerKey).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("CompanyOrderReturnCnfService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取CompanyOrderReturnCnf对象
func (e *CompanyOrderReturnCnf) Get(d *dto.CompanyOrderReturnCnfGetReq, p *actions.DataPermission, model *models.CompanyOrderReturnCnf) error {
	var data models.CompanyOrderReturnCnf

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetCompanyOrderReturnCnf error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建CompanyOrderReturnCnf对象
func (e *CompanyOrderReturnCnf) Insert(c *dto.CompanyOrderReturnCnfInsertReq) error {
    var err error
    var data models.CompanyOrderReturnCnf
    c.Generate(&data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("CompanyOrderReturnCnfService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改CompanyOrderReturnCnf对象
func (e *CompanyOrderReturnCnf) Update(c *dto.CompanyOrderReturnCnfUpdateReq, p *actions.DataPermission) error {
    var err error
    var data = models.CompanyOrderReturnCnf{}
    e.Orm.Scopes(
            actions.Permission(data.TableName(), p),
        ).First(&data, c.GetId())
    c.Generate(&data)

    db := e.Orm.Save(&data)
    if err = db.Error; err != nil {
        e.Log.Errorf("CompanyOrderReturnCnfService Save error:%s \r\n", err)
        return err
    }
    if db.RowsAffected == 0 {
        return errors.New("无权更新该数据")
    }
    return nil
}

// Remove 删除CompanyOrderReturnCnf
func (e *CompanyOrderReturnCnf) Remove(cid int,d *dto.CompanyOrderReturnCnfDeleteReq) error {
	var data models.CompanyOrderReturnCnf

	db := e.Orm.Model(&data).Where("c_id = ?",cid).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
        e.Log.Errorf("Service RemoveCompanyOrderReturnCnf error:%s \r\n", err)
        return err
    }
    if db.RowsAffected == 0 {
        return errors.New("无权删除该数据")
    }
	return nil
}
