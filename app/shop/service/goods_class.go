package service

import (
	"errors"

    "github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/gorm"

	"go-admin/app/shop/models"
	"go-admin/app/shop/service/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
)

type GoodsClass struct {
	service.Service
}

// GetPage 获取GoodsClass列表
func (e *GoodsClass) GetPage(c *dto.GoodsClassGetPageReq, p *actions.DataPermission, list *[]models.GoodsClass, count *int64) error {
	var err error
	var data models.GoodsClass

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("GoodsClassService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取GoodsClass对象
func (e *GoodsClass) Get(d *dto.GoodsClassGetReq, p *actions.DataPermission, model *models.GoodsClass) error {
	var data models.GoodsClass

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetGoodsClass error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建GoodsClass对象
func (e *GoodsClass) Insert(cId int,c *dto.GoodsClassInsertReq) error {
    var err error
    var data models.GoodsClass
    c.Generate(&data)
    data.CId = cId
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("分类创建失败,", err)
		return err
	}
	return nil
}

// Update 修改GoodsClass对象
func (e *GoodsClass) Update(c *dto.GoodsClassUpdateReq, p *actions.DataPermission) error {
    var err error
    var data = models.GoodsClass{}
    e.Orm.Scopes(
            actions.Permission(data.TableName(), p),
        ).First(&data, c.GetId())
    c.Generate(&data)

    db := e.Orm.Save(&data)
    if err = db.Error; err != nil {
        e.Log.Errorf("分类更新失败,%s", err)
        return err
    }
    if db.RowsAffected == 0 {
        return errors.New("无权更新该数据")
    }
    return nil
}

// Remove 删除GoodsClass
func (e *GoodsClass) Remove(d *dto.GoodsClassDeleteReq, p *actions.DataPermission) error {
	var data models.GoodsClass

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
        e.Log.Errorf("Service RemoveGoodsClass error:%s \r\n", err)
        return err
    }
    if db.RowsAffected == 0 {
        return errors.New("无权删除该数据")
    }
	return nil
}
