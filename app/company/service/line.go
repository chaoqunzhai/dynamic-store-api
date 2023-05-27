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

type Line struct {
	service.Service
}

// GetPage 获取Line列表
func (e *Line) GetPage(c *dto.LineGetPageReq, p *actions.DataPermission, list *[]models.Line, count *int64) error {
	var err error
	var data models.Line

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).Order(global.OrderLayerKey).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("LineService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取Line对象
func (e *Line) Get(d *dto.LineGetReq, p *actions.DataPermission, model *models.Line) error {
	var data models.Line

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetLine error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建Line对象
func (e *Line) Insert(cid int,c *dto.LineInsertReq) error {
    var err error
    var data models.Line
    c.Generate(&data)
    data.CId = cid
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("LineService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改Line对象
func (e *Line) Update(c *dto.LineUpdateReq, p *actions.DataPermission) error {
    var err error
    var data = models.Line{}
    e.Orm.Scopes(
            actions.Permission(data.TableName(), p),
        ).First(&data, c.GetId())
    c.Generate(&data)

    db := e.Orm.Save(&data)
    if err = db.Error; err != nil {
        e.Log.Errorf("LineService Save error:%s \r\n", err)
        return err
    }
    if db.RowsAffected == 0 {
        return errors.New("无权更新该数据")
    }
    return nil
}

// Remove 删除Line
func (e *Line) Remove(d *dto.LineDeleteReq, p *actions.DataPermission) error {
	var data models.Line

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
        e.Log.Errorf("Service RemoveLine error:%s \r\n", err)
        return err
    }
    if db.RowsAffected == 0 {
        return errors.New("无权删除该数据")
    }
	return nil
}
