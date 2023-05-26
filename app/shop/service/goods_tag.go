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

type GoodsTag struct {
	service.Service
}

// GetPage 获取GoodsTag列表
func (e *GoodsTag) GetPage(c *dto.GoodsTagGetPageReq, p *actions.DataPermission, list *[]models.GoodsTag, count *int64) error {
	var err error
	var data models.GoodsTag

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("GoodsTagService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取GoodsTag对象
func (e *GoodsTag) Get(d *dto.GoodsTagGetReq, p *actions.DataPermission, model *models.GoodsTag) error {
	var data models.GoodsTag

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetGoodsTag error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建GoodsTag对象
func (e *GoodsTag) Insert(cId int,c *dto.GoodsTagInsertReq) error {
    var err error
    var data models.GoodsTag
    c.Generate(&data)
	data.CId = cId
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("GoodsTagService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改GoodsTag对象
func (e *GoodsTag) Update(c *dto.GoodsTagUpdateReq, p *actions.DataPermission) error {
    var err error
    var data = models.GoodsTag{}
    e.Orm.Scopes(
            actions.Permission(data.TableName(), p),
        ).First(&data, c.GetId())
    c.Generate(&data)

    db := e.Orm.Save(&data)
    if err = db.Error; err != nil {
        e.Log.Errorf("修改标签失败,", err)
        return err
    }
    if db.RowsAffected == 0 {
        return errors.New("无权更新该数据")
    }
    return nil
}

// Remove 删除GoodsTag
func (e *GoodsTag) Remove(d *dto.GoodsTagDeleteReq, p *actions.DataPermission) error {
	var data models.GoodsTag

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
        e.Log.Errorf("Service RemoveGoodsTag error:%s \r\n", err)
        return err
    }
    if db.RowsAffected == 0 {
        return errors.New("无权删除该数据")
    }
	return nil
}
