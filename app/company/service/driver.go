package service

import (
	"errors"
	sys "go-admin/app/admin/models"
	"go-admin/global"

	"github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/gorm"

	"go-admin/app/company/models"
	"go-admin/app/company/service/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
)

type Driver struct {
	service.Service
}

// GetPage 获取Driver列表
func (e *Driver) GetPage(c *dto.DriverGetPageReq, p *actions.DataPermission, list *[]models.Driver, count *int64) error {
	var err error
	var data models.Driver

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).Order(global.OrderLayerKey).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("DriverService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取Driver对象
func (e *Driver) Get(d *dto.DriverGetReq, p *actions.DataPermission, model *models.Driver) error {
	var data models.Driver

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetDriver error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建Driver对象
func (e *Driver) Insert(cid int,c *dto.DriverInsertReq) error {
    var err error
    var data models.Driver
    c.Generate(&data)
    data.CId = cid

    if c.Phone != "" {
    	var sysUser sys.SysUser
    	e.Orm.Model(&sys.SysUser{}).Where("phone = ? and enable = ? and c_id = ?",c.Phone,true,cid).Limit(1).Find(&sysUser)
    	if sysUser.UserId > 0 {
    		data.UserId = sysUser.UserId
		}
	}
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("DriverService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改Driver对象
func (e *Driver) Update(cid int,c *dto.DriverUpdateReq, p *actions.DataPermission) error {
    var err error
    var data = models.Driver{}
    e.Orm.Scopes(
            actions.Permission(data.TableName(), p),
        ).First(&data, c.GetId())
    c.Generate(&data)
	if c.Phone != "" {
		var sysUser sys.SysUser
		e.Orm.Model(&sys.SysUser{}).Where("phone = ? and enable = ? and c_id = ?",c.Phone,true,cid).Limit(1).Find(&sysUser)
		if sysUser.UserId > 0 {
			data.UserId = sysUser.UserId
		}
	}
    db := e.Orm.Save(&data)
    if err = db.Error; err != nil {
        e.Log.Errorf("DriverService Save error:%s \r\n", err)
        return err
    }
    if db.RowsAffected == 0 {
        return errors.New("无权更新该数据")
    }
    return nil
}

// Remove 删除Driver
func (e *Driver) Remove(d *dto.DriverDeleteReq, p *actions.DataPermission) error {
	var data models.Driver

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
        e.Log.Errorf("Service RemoveDriver error:%s \r\n", err)
        return err
    }
    if db.RowsAffected == 0 {
        return errors.New("无权删除该数据")
    }
	return nil
}
