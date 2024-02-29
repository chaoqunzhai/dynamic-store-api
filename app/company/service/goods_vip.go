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

type GoodsVip struct {
	service.Service
}

// GetPage 获取GoodsVip列表
func (e *GoodsVip) GetPage(c *dto.GoodsVipGetPageReq, p *actions.DataPermission, list *[]models.GoodsVip, count *int64) error {
	var err error
	var data models.GoodsVip

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).Order(global.OrderLayerKey).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("GoodsVipService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取GoodsVip对象
func (e *GoodsVip) Get(d *dto.GoodsVipGetReq, p *actions.DataPermission, model *models.GoodsVip) error {
	var data models.GoodsVip

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetGoodsVip error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建GoodsVip对象
func (e *GoodsVip) Insert(c *dto.GoodsVipInsertReq) error {
    var err error
    var data models.GoodsVip
    c.Generate(&data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("GoodsVipService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改GoodsVip对象
func (e *GoodsVip) Update(c *dto.GoodsVipUpdateReq, p *actions.DataPermission) error {
    var err error
    var data = models.GoodsVip{}
    e.Orm.Scopes(
            actions.Permission(data.TableName(), p),
        ).First(&data, c.GetId())
    c.Generate(&data)

    db := e.Orm.Save(&data)
    if err = db.Error; err != nil {
        e.Log.Errorf("GoodsVipService Save error:%s \r\n", err)
        return err
    }
    if db.RowsAffected == 0 {
        return errors.New("无权更新该数据")
    }
    return nil
}

// Remove 删除GoodsVip
func (e *GoodsVip) Remove(d *dto.GoodsVipDeleteReq, cid int) error {
	var data models.GoodsVip

	db :=e.Orm.Model(&data).Where("c_id = ? and id in ?",cid,d.GetId()).Delete(&data)
	if err := db.Error; err != nil {
        e.Log.Errorf("Service RemoveGoodsVip error:%s \r\n", err)
        return err
    }
    if db.RowsAffected == 0 {
        return errors.New("无权删除该数据")
    }
	return nil
}
