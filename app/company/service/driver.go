package service

import (
	"errors"
	sys "go-admin/app/admin/models"
	"go-admin/global"
	"golang.org/x/crypto/bcrypt"

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
	if data.UserId > 0 {
		var user sys.SysUser
		e.Orm.Model(&sys.SysUser{}).Select("password").Where("user_id = ? and c_id = ?",data.UserId,data.CId).Limit(1).Find(&user)
		data.Password = user.Password
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
		}else {
			userObject := sys.SysUser{
				Username: c.Name,
				Phone:    c.Phone,
				Enable:   true,
				Status:   "2",
				Password: c.PassWord,
				CId:      cid,
				RoleId:   global.RoleDriver,
				Layer:    1,
				AuthExamine: false,
				AuthLoginMbm: false,
			}
			userObject.CreateBy = c.CreateBy
			e.Orm.Create(&userObject)
			data.UserId = userObject.UserId
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
	var userObject sys.SysUser
	e.Orm.Model(&sys.SysUser{}).Where("phone = ? and enable = ? and c_id = ?",c.Phone,true,cid).Limit(1).Find(&userObject)
	if userObject.UserId > 0 {
		data.UserId = userObject.UserId
	}else {
		userObject = sys.SysUser{
			Username: c.Name,
			Phone:    c.Phone,
			Enable:   true,
			Status:   "2",
			Password: c.PassWord,
			CId:      cid,
			RoleId:   global.RoleDriver,
			Layer:    1,
			AuthExamine: false,
			AuthLoginMbm: false,
		}
		userObject.CreateBy = c.CreateBy
		e.Orm.Create(&userObject)
		data.UserId = userObject.UserId
	}
    db := e.Orm.Save(&data)
    if err = db.Error; err != nil {
        e.Log.Errorf("DriverService Save error:%s \r\n", err)
        return err
    }
    if db.RowsAffected == 0 {
        return errors.New("无权更新该数据")
    }

	if c.PassWord != userObject.Password {
		updateMap:=make(map[string]interface{},0)
		var hash []byte
		var hasErr error
		if hash, hasErr = bcrypt.GenerateFromPassword([]byte(c.PassWord), bcrypt.DefaultCost); hasErr != nil {
			return errors.New("密码生成失败")
		} else {
			updateMap["password"] = string(hash)
			e.Orm.Model(&sys.SysUser{}).Where("c_id = ? and user_id = ?", cid, userObject.UserId).Updates(&updateMap)
		}
	}
    return nil
}

// Remove 删除Driver
func (e *Driver) Remove(d *dto.DriverDeleteReq, p *actions.DataPermission,CId int) error {

	for _, dID := range d.Ids {
		var data models.Driver

		e.Orm.Model(&data).Select("user_id").Where("id = ? and c_id = ?",dID,CId).Limit(1).Find(&data)

		e.Orm.Model(&data).Where("id = ? and c_id = ?",dID,CId).Delete(&models.Goods{})
		if data.UserId > 0 {
			e.Orm.Model(&sys.SysUser{}).Where("user_id = ? and c_id = ?",data.UserId,CId).Delete(&models.Driver{})
		}
	}

	return nil
}
