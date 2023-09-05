package service

import (
	"errors"
	"fmt"
	"github.com/go-admin-team/go-admin-core/sdk/service"
	"go-admin/global"
	"gorm.io/gorm"
	"strings"

	"go-admin/app/shop/models"
	"go-admin/app/shop/service/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
)

type Shop struct {
	service.Service
}

// GetPage 获取Shop列表
func (e *Shop) GetPage(c *dto.ShopGetPageReq, p *actions.DataPermission, list *[]models.Shop, count *int64) error {
	var err error
	var data models.Shop

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).Order(global.OrderLayerKey).Preload("Tag", func(tx *gorm.DB) *gorm.DB {
		return tx.Select("id", "name")
	}).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("ShopService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取Shop对象
func (e *Shop) Get(d *dto.ShopGetReq, p *actions.DataPermission, model *models.Shop) error {
	var data models.Shop

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Preload("Tag", func(tx *gorm.DB) *gorm.DB {
		return tx.Select("id", "name")
	}).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetShop error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}
func (e *Shop) getShopTagModels(ids []int) (list []models.ShopTag) {
	for _, id := range ids {
		var row models.ShopTag
		e.Orm.Model(&models.ShopTag{}).Where("id = ?", id).First(&row)
		if row.Id == 0 {
			continue
		}
		list = append(list, row)
	}
	return list
}

// Insert 创建Shop对象
func (e *Shop) Insert(cid int, c *dto.ShopInsertReq) error {
	var err error
	var data models.Shop
	c.Generate(&data)
	data.CId = cid

	if len(c.Tags) > 0 {
		data.Tag = e.getShopTagModels(c.Tags)
	}
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("ShopService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改Shop对象
func (e *Shop) Update(c *dto.ShopUpdateReq, p *actions.DataPermission) error {
	var err error
	var data = models.Shop{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).First(&data, c.GetId())
	c.Generate(&data)
	//清除关联
	_=e.Orm.Model(&data).Association("Tag").Clear()
	if len(c.Tags) > 0 {
		//增加关联
		fmt.Println("标签", c.Tags)
		data.Tag = e.getShopTagModels(c.Tags)
	}
	db := e.Orm.Save(&data)
	if err = db.Error; err != nil {
		e.Log.Errorf("ShopService Save error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	return nil
}

// Remove 删除Shop
func (e *Shop) Remove(d *dto.ShopDeleteReq, p *actions.DataPermission) error {
	var data models.Shop

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("用户删除失败,", err)
		return err
	}

	removeIds := make([]string, 0)
	for _, t := range d.Ids {
		removeIds = append(removeIds, fmt.Sprintf("%v", t))
	}
	e.Orm.Exec(fmt.Sprintf("DELETE FROM `shop_mark_tag` WHERE `shop_mark_tag`.`shop_id` IN (%v)", strings.Join(removeIds, ",")))
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}
