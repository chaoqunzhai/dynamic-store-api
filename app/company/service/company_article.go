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

type CompanyArticle struct {
	service.Service
}

// GetPage 获取CompanyArticle列表
func (e *CompanyArticle) GetPage(c *dto.CompanyArticleGetPageReq, p *actions.DataPermission, list *[]models.CompanyArticle, count *int64) error {
	var err error
	var data models.CompanyArticle

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("CompanyArticleService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取CompanyArticle对象
func (e *CompanyArticle) Get(d *dto.CompanyArticleGetReq, p *actions.DataPermission, model *models.CompanyArticle) error {
	var data models.CompanyArticle

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetCompanyArticle error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建CompanyArticle对象
func (e *CompanyArticle) Insert(c *dto.CompanyArticleInsertReq) error {
    var err error
    var data models.CompanyArticle
    c.Generate(&data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("CompanyArticleService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改CompanyArticle对象
func (e *CompanyArticle) Update(c *dto.CompanyArticleUpdateReq, p *actions.DataPermission) error {
    var err error
    var data = models.CompanyArticle{}
    e.Orm.Scopes(
            actions.Permission(data.TableName(), p),
        ).First(&data, c.GetId())
    c.Generate(&data)

    db := e.Orm.Save(&data)
    if err = db.Error; err != nil {
        e.Log.Errorf("CompanyArticleService Save error:%s \r\n", err)
        return err
    }
    if db.RowsAffected == 0 {
        return errors.New("无权更新该数据")
    }
    return nil
}

// Remove 删除CompanyArticle
func (e *CompanyArticle) Remove(d *dto.CompanyArticleDeleteReq, p *actions.DataPermission) error {
	var data models.CompanyArticle

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
        e.Log.Errorf("Service RemoveCompanyArticle error:%s \r\n", err)
        return err
    }
    if db.RowsAffected == 0 {
        return errors.New("无权删除该数据")
    }
	return nil
}
