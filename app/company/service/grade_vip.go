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

type GradeVip struct {
	service.Service
}

// GetPage 获取GradeVip列表
func (e *GradeVip) GetPage(c *dto.GradeVipGetPageReq, p *actions.DataPermission, list *[]models.GradeVip, count *int64) error {
	var err error
	var data models.GradeVip

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).Order(global.OrderLayerKey).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("GradeVipService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取GradeVip对象
func (e *GradeVip) Get(d *dto.GradeVipGetReq, p *actions.DataPermission, model *models.GradeVip) error {
	var data models.GradeVip

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetGradeVip error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建GradeVip对象
func (e *GradeVip) Insert(cid int,c *dto.GradeVipInsertReq) error {
    var err error
    var data models.GradeVip
    c.Generate(&data)
    data.CId = cid
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("GradeVipService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改GradeVip对象
func (e *GradeVip) Update(c *dto.GradeVipUpdateReq, p *actions.DataPermission) error {
    var err error
    var data = models.GradeVip{}
    e.Orm.Scopes(
            actions.Permission(data.TableName(), p),
        ).First(&data, c.GetId())
    c.Generate(&data)


    db := e.Orm.Save(&data)
    if err = db.Error; err != nil {
        e.Log.Errorf("GradeVipService Save error:%s \r\n", err)
        return err
    }
    if db.RowsAffected == 0 {
        return errors.New("无权更新该数据")
    }
    return nil
}

// Remove 删除GradeVip
func (e *GradeVip) Remove(d *dto.GradeVipDeleteReq, cid int) error {
	var data models.GradeVip

	db :=e.Orm.Model(&data).Where("c_id = ? and id in ?",cid,d.GetId()).Delete(&data)
	if err := db.Error; err != nil {
        e.Log.Errorf("Service RemoveGradeVip error:%s \r\n", err)
        return err
    }
    if db.RowsAffected == 0 {
        return errors.New("无权删除该数据")
    }
	return nil
}
