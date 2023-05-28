package service

import (
	"errors"

    "github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/gorm"

	"go-admin/app/system/models"
	"go-admin/app/system/service/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
)

type ExtendUser struct {
	service.Service
}

// GetPage 获取ExtendUser列表
func (e *ExtendUser) GetPage(c *dto.ExtendUserGetPageReq, p *actions.DataPermission, list *[]models.ExtendUser, count *int64) error {
	var err error
	var data models.ExtendUser

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("ExtendUserService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取ExtendUser对象
func (e *ExtendUser) Get(d *dto.ExtendUserGetReq, p *actions.DataPermission, model *models.ExtendUser) error {
	var data models.ExtendUser

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetExtendUser error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建ExtendUser对象
func (e *ExtendUser) Insert(cid int,c *dto.ExtendUserInsertReq) error {
    var err error
    var data models.ExtendUser
    c.Generate(&data)
    data.CId  = cid
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("ExtendUserService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改ExtendUser对象
func (e *ExtendUser) Update(c *dto.ExtendUserUpdateReq, p *actions.DataPermission) error {
    var err error
    var data = models.ExtendUser{}
    e.Orm.Scopes(
            actions.Permission(data.TableName(), p),
        ).First(&data, c.GetId())
    c.Generate(&data)

    db := e.Orm.Save(&data)
    if err = db.Error; err != nil {
        e.Log.Errorf("ExtendUserService Save error:%s \r\n", err)
        return err
    }
    if db.RowsAffected == 0 {
        return errors.New("无权更新该数据")
    }
    return nil
}

// Remove 删除ExtendUser
func (e *ExtendUser) Remove(d *dto.ExtendUserDeleteReq, p *actions.DataPermission) error {
	var data models.ExtendUser

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
        e.Log.Errorf("Service RemoveExtendUser error:%s \r\n", err)
        return err
    }
    if db.RowsAffected == 0 {
        return errors.New("无权删除该数据")
    }
	return nil
}
