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

type CompanyDebitCard struct {
	service.Service
}

// GetPage 获取CompanyDebitCard列表
func (e *CompanyDebitCard) GetPage(c *dto.CompanyDebitCardGetPageReq, p *actions.DataPermission, list *[]models.CompanyDebitCard, count *int64) error {
	var err error
	var data models.CompanyDebitCard

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("CompanyDebitCardService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取CompanyDebitCard对象
func (e *CompanyDebitCard) Get(d *dto.CompanyDebitCardGetReq, p *actions.DataPermission, model *models.CompanyDebitCard) error {
	var data models.CompanyDebitCard

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetCompanyDebitCard error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建CompanyDebitCard对象
func (e *CompanyDebitCard) Insert(c *dto.CompanyDebitCardInsertReq) error {
    var err error
    var data models.CompanyDebitCard
    c.Generate(&data)
    data.Layer = "0"
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("CompanyDebitCardService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改CompanyDebitCard对象
func (e *CompanyDebitCard) Update(c *dto.CompanyDebitCardUpdateReq, p *actions.DataPermission) error {
    var err error
    var data = models.CompanyDebitCard{}
    e.Orm.Scopes(
            actions.Permission(data.TableName(), p),
        ).Where("c_id = ? and id = ?",c.CId,c.GetId()).First(&data)
    c.Generate(&data)
	data.Enable = true
	data.Layer = "0"
    db := e.Orm.Save(&data)
    if err = db.Error; err != nil {
        e.Log.Errorf("CompanyDebitCardService Save error:%s \r\n", err)
        return err
    }
    if db.RowsAffected == 0 {
        return errors.New("无权更新该数据")
    }
    return nil
}

// Remove 删除CompanyDebitCard
func (e *CompanyDebitCard) Remove(d *dto.CompanyDebitCardDeleteReq, p *actions.DataPermission) error {
	var data models.CompanyDebitCard

	db := e.Orm.Model(&data).Scopes(actions.Permission(data.TableName(), p)).Unscoped().Where("id = ? ",d.GetId()).Delete(&data)
	if err := db.Error; err != nil {
        e.Log.Errorf("Service RemoveCompanyDebitCard error:%s \r\n", err)
        return err
    }
    if db.RowsAffected == 0 {
        return errors.New("无权删除该数据")
    }
	return nil
}
