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

type CompanyOfflinePay struct {
	service.Service
}

// GetPage 获取CompanyOfflinePay列表
func (e *CompanyOfflinePay) GetPage(c *dto.CompanyOfflinePayGetPageReq, p *actions.DataPermission, list *[]models.CompanyOfflinePay, count *int64) error {
	var err error
	var data models.CompanyOfflinePay

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("CompanyOfflinePayService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取CompanyOfflinePay对象
func (e *CompanyOfflinePay) Get(d *dto.CompanyOfflinePayGetReq, p *actions.DataPermission, model *models.CompanyOfflinePay) error {
	var data models.CompanyOfflinePay

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetCompanyOfflinePay error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建CompanyOfflinePay对象
func (e *CompanyOfflinePay) Insert(c *dto.CompanyOfflinePayInsertReq) error {
    var err error
    var data models.CompanyOfflinePay
    c.Generate(&data)
    data.Enable = true
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("CompanyOfflinePayService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改CompanyOfflinePay对象
func (e *CompanyOfflinePay) Update(c *dto.CompanyOfflinePayUpdateReq, p *actions.DataPermission) error {
    var err error
    var data = models.CompanyOfflinePay{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).Where("c_id = ? and id = ?",c.CId,c.GetId()).First(&data)
    c.Generate(&data)
	data.Enable = true
    db := e.Orm.Save(&data)
    if err = db.Error; err != nil {
        e.Log.Errorf("CompanyOfflinePayService Save error:%s \r\n", err)
        return err
    }
    if db.RowsAffected == 0 {
        return errors.New("无权更新该数据")
    }
    return nil
}

// Remove 删除CompanyOfflinePay
func (e *CompanyOfflinePay) Remove(d *dto.CompanyOfflinePayDeleteReq, p *actions.DataPermission) error {
	var data models.CompanyOfflinePay
	db := e.Orm.Model(&data).Unscoped().Where("c_id = ? and id = ? ",d.CId,d.GetId()).Delete(&data)
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveCompanyDebitCard error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}
