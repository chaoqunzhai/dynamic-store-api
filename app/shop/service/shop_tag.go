package service

import (
	"errors"
	"go-admin/global"

	"github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/gorm"

	"go-admin/app/shop/models"
	"go-admin/app/shop/service/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
)

type ShopTag struct {
	service.Service
}

// GetPage 获取ShopTag列表
func (e *ShopTag) GetPage(c *dto.ShopTagGetPageReq, p *actions.DataPermission, list *[]models.ShopTag, count *int64) error {
	var err error
	var data models.ShopTag

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).Order(global.OrderLayerKey).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("ShopTagService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取ShopTag对象
func (e *ShopTag) Get(d *dto.ShopTagGetReq, p *actions.DataPermission, model *models.ShopTag) error {
	var data models.ShopTag

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetShopTag error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建ShopTag对象
func (e *ShopTag) Insert(cid int,c *dto.ShopTagInsertReq) error {
    var err error
    var data models.ShopTag
    c.Generate(&data)
    data.CId = cid
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("ShopTagService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改ShopTag对象
func (e *ShopTag) Update(c *dto.ShopTagUpdateReq, p *actions.DataPermission) error {
    var err error
    var data = models.ShopTag{}
    e.Orm.Scopes(
            actions.Permission(data.TableName(), p),
        ).First(&data, c.GetId())
    c.Generate(&data)

    db := e.Orm.Save(&data)
    if err = db.Error; err != nil {
        e.Log.Errorf("ShopTagService Save error:%s \r\n", err)
        return err
    }
    if db.RowsAffected == 0 {
        return errors.New("无权更新该数据")
    }
    return nil
}

// Remove 删除ShopTag
func (e *ShopTag) Remove(d *dto.ShopTagDeleteReq, p *actions.DataPermission) error {
	var data models.ShopTag

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
        e.Log.Errorf("Service RemoveShopTag error:%s \r\n", err)
        return err
    }
    if db.RowsAffected == 0 {
        return errors.New("无权删除该数据")
    }
	return nil
}
