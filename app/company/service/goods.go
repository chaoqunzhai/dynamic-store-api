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

type Goods struct {
	service.Service
}

// GetPage 获取Goods列表
func (e *Goods) GetPage(c *dto.GoodsGetPageReq, p *actions.DataPermission, list *[]models.Goods, count *int64) error {
	var err error
	var data models.Goods

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).Order(global.OrderLayerKey).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("GoodsService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取Goods对象
func (e *Goods) Get(d *dto.GoodsGetReq, p *actions.DataPermission, model *models.Goods) error {
	var data models.Goods

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetGoods error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

func (e *Goods)getTagModels(ids []int) (list []models.GoodsTag) {
	for _, id := range ids {
		var row models.GoodsTag
		e.Orm.Model(&models.GoodsTag{}).Where("id = ?", id).First(&row)
		if row.Id == 0 {
			continue
		}
		list = append(list, row)
	}
	return list
}
func (e *Goods)getClassModels(ids []int) (list []models.GoodsClass) {
	for _, id := range ids {
		var row models.GoodsClass
		e.Orm.Model(&models.GoodsClass{}).Where("id = ?", id).First(&row)
		if row.Id == 0 {
			continue
		}
		list = append(list, row)
	}
	return list
}
// Insert 创建Goods对象
func (e *Goods) Insert(cid int,c *dto.GoodsInsertReq) error {
    var err error
    var data models.Goods
    c.Generate(&data)
    data.CId = cid

    //标签
    if len(c.Tag) > 0 {
		data.Tag = e.getTagModels(c.Tag)
	}
	//分类
	if len(c.Class) > 0 {
		data.Class = e.getClassModels(c.Tag)
	}
	err = e.Orm.Create(&data).Error

	//规格 + vip价格设置存在
	if len(c.Specs) > 0 {

		for _,row:=range c.Specs{
			specsModels :=models.GoodsSpecs{
				Name: row.Name,
				CId: cid,
				Enable: true,
				Layer: 0,
				GoodsId: data.Id,
				Price: row.Price,
				Original: row.Original,
				Inventory: row.Inventory,
				Unit: row.Unit,
				Limit: row.Limit,
			}
			specsModels.CreateBy = data.CreateBy
			e.Orm.Create(&specsModels)
			for _,v:=range row.Vip{
				var gradeRow models.GradeVip
				e.Orm.Model(&models.GradeVip{}).Where("enable = ? and id = ?",true,v.Grade).Limit(1).Find(&gradeRow)
				if gradeRow.Id == 0 {continue}
				vipRow:=models.GoodsVip{
					CId: cid,
					GoodsId:data.Id,
					Enable: true,
					GradeId: gradeRow.Id,
					Layer: 0,
					CustomPrice: v.Price,
				}
				vipRow.CreateBy = data.CreateBy
				e.Orm.Create(&vipRow)
			}
		}
	}

	if err != nil {
		e.Log.Errorf("GoodsService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改Goods对象
func (e *Goods) Update(c *dto.GoodsUpdateReq, p *actions.DataPermission) error {
    var err error
    var data = models.Goods{}
    e.Orm.Scopes(
            actions.Permission(data.TableName(), p),
        ).First(&data, c.GetId())
    c.Generate(&data)

    db := e.Orm.Save(&data)
    if err = db.Error; err != nil {
        e.Log.Errorf("GoodsService Save error:%s \r\n", err)
        return err
    }
    if db.RowsAffected == 0 {
        return errors.New("无权更新该数据")
    }
    return nil
}

// Remove 删除Goods
func (e *Goods) Remove(d *dto.GoodsDeleteReq, p *actions.DataPermission) error {
	var data models.Goods

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
        e.Log.Errorf("Service RemoveGoods error:%s \r\n", err)
        return err
    }
    if db.RowsAffected == 0 {
        return errors.New("无权删除该数据")
    }
	return nil
}
