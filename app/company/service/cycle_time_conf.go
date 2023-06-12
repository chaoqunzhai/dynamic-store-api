package service

import (
	"errors"
	"fmt"
	"github.com/go-admin-team/go-admin-core/sdk/service"
	"go-admin/common/utils"
	"go-admin/global"
	"gorm.io/gorm"

	"go-admin/app/company/models"
	"go-admin/app/company/service/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
)

type CycleTimeConf struct {
	service.Service
}

// GetPage 获取CycleTimeConf列表
func (e *CycleTimeConf) GetPage(c *dto.CycleTimeConfGetPageReq, p *actions.DataPermission, list *[]models.CycleTimeConf, count *int64) error {
	var err error
	var data models.CycleTimeConf

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).Order("created_at desc").
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("CycleTimeConfService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取CycleTimeConf对象
func (e *CycleTimeConf) Get(d *dto.CycleTimeConfGetReq, p *actions.DataPermission, model *models.CycleTimeConf) error {
	var data models.CycleTimeConf

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetCycleTimeConf error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建CycleTimeConf对象
func (e *CycleTimeConf) Insert(cid int, c *dto.CycleTimeConfInsertReq) error {
	var err error
	var data models.CycleTimeConf
	c.Generate(&data)
	data.CId = cid
	data.Enable = true
	data.Uid = utils.CreateCode()
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("CycleTimeConfService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改CycleTimeConf对象
func (e *CycleTimeConf) Update(c *dto.CycleTimeConfUpdateReq, p *actions.DataPermission) error {
	var err error
	var data = models.CycleTimeConf{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).First(&data, c.GetId())
	//判断时间是否有变动,如果有变动更新uid标记
	uidTag := false
	if c.StartTime != data.StartTime || c.EndTime != data.EndTime {
		uidTag = true
	}
	switch c.Type {
	case global.CyCleTimeWeek:
		if c.StartWeek != data.StartWeek || c.EndWeek != data.EndWeek {
			uidTag = true
		}
	}
	if uidTag {
		data.Uid = utils.CreateCode()
		fmt.Println("数据发生变更更新code")
	}
	c.Generate(&data)
	data.Enable = true
	db := e.Orm.Save(&data)
	if err = db.Error; err != nil {
		e.Log.Errorf("CycleTimeConfService Save error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	return nil
}

// Remove 删除CycleTimeConf
func (e *CycleTimeConf) Remove(d *dto.CycleTimeConfDeleteReq, p *actions.DataPermission) error {
	var data models.CycleTimeConf

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveCycleTimeConf error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}
