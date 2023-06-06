package service

import (
	"errors"
	"fmt"
	"go-admin/global"
	"strings"

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

func (e *Goods) getTagModels(ids []int) (list []models.GoodsTag) {
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
func (e *Goods) getClassModels(ids []int) (list []models.GoodsClass) {
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
func (e *Goods) Insert(cid int, c *dto.GoodsInsertReq) (uid int, err error) {

	var data models.Goods
	c.Generate(&data)
	data.CId = cid

	//标签
	if len(c.Tag) > 0 {
		data.Tag = e.getTagModels(c.Tag)
	}
	//分类
	if len(c.Class) > 0 {
		data.Class = e.getClassModels(c.Class)
	}
	err = e.Orm.Create(&data).Error

	//规格 + vip价格设置存在
	if len(c.Specs) > 0 {

		for _, row := range c.Specs {
			specsModels := models.GoodsSpecs{
				Name:      row.Name,
				CId:       cid,
				Enable:    row.Enable,
				Layer:     row.Layer,
				GoodsId:   data.Id,
				Price:     row.Price,
				Original:  row.Original,
				Inventory: row.Inventory,
				Unit:      row.Unit,
				Limit:     row.Limit,
			}
			specsModels.CreateBy = data.CreateBy
			e.Orm.Create(&specsModels)
			var vipEnable bool
			for k, v := range row.Vip {
				if k == "enable" {
					vipEnable = v.(bool)
				}
				if !strings.HasPrefix(k, "vip_") {
					continue
				}
				gradeInt := strings.Replace(k, "vip_", "", -1)

				var gradeRow models.GradeVip
				e.Orm.Model(&models.GradeVip{}).Where("enable = ? and id = ? and c_id = ?", true, gradeInt, cid).Limit(1).Find(&gradeRow)
				if gradeRow.Id == 0 {
					continue
				}
				vipRow := models.GoodsVip{
					CId:         cid,
					GoodsId:     data.Id,
					Enable:      vipEnable,
					Layer:       0,
					GradeId:     gradeRow.Id,
					CustomPrice: v.(float64),
				}
				vipRow.CreateBy = data.CreateBy
				e.Orm.Create(&vipRow)
			}
		}
	}

	if err != nil {
		e.Log.Errorf("GoodsService Insert error:%s \r\n", err)
		return 0, err
	}
	return data.Id, err
}

// Update 修改Goods对象
func (e *Goods) Update(cid int, c *dto.GoodsUpdateReq, p *actions.DataPermission) error {
	var err error
	var data = models.Goods{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).First(&data, c.GetId())
	c.Generate(&data)

	//标签
	e.Orm.Model(&data).Association("Tag").Clear()
	if len(c.Tag) > 0 {
		data.Tag = e.getTagModels(c.Tag)
	}
	//分类
	e.Orm.Model(&data).Association("Class").Clear()
	if len(c.Class) > 0 {
		data.Class = e.getClassModels(c.Class)
	}

	//规格更新
	if len(c.Specs) > 0 {
		for _, row := range c.Specs {
			if row.Id > 0 {
				//就是一个规格资源的更新
				var specsRow models.GoodsSpecs
				e.Orm.Model(&models.GoodsSpecs{}).Where("id = ?", row.Id).First(&specsRow)

				specsRow.Name = row.Name
				specsRow.Enable = row.Enable
				specsRow.Layer = row.Layer
				specsRow.Price = row.Price
				specsRow.Original = row.Original
				specsRow.Inventory = row.Inventory
				specsRow.Unit = row.Unit
				specsRow.Limit = row.Limit
				e.Orm.Save(&specsRow)
			} else {
				//规格资源的创建
				specsModels := models.GoodsSpecs{
					Name:      row.Name,
					CId:       cid,
					Enable:    row.Enable,
					Layer:     row.Layer,
					GoodsId:   data.Id,
					Price:     row.Price,
					Original:  row.Original,
					Inventory: row.Inventory,
					Unit:      row.Unit,
					Limit:     row.Limit,
				}
				specsModels.CreateBy = data.CreateBy
				e.Orm.Create(&specsModels)
			}
			for _, v := range row.Vip {
				if v.Id > 0 {
					//vip价格的更新
				} else {
					//vip价格的创建
					var gradeRow models.GradeVip
					e.Orm.Model(&models.GradeVip{}).Where("enable = ? and id = ?", true, v.Grade).Limit(1).Find(&gradeRow)
					if gradeRow.Id == 0 {
						continue
					}
					vipRow := models.GoodsVip{
						CId:         cid,
						GoodsId:     data.Id,
						Enable:      v.Enable,
						GradeId:     gradeRow.Id,
						Layer:       v.Layer,
						CustomPrice: v.Price,
					}
					vipRow.CreateBy = data.CreateBy
					e.Orm.Create(&vipRow)
				}
			}
		}

	}
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
		e.Log.Errorf("删除失败,", err)
		return err
	}

	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	removeIds := make([]string, 0)
	for _, t := range d.Ids {
		removeIds = append(removeIds, fmt.Sprintf("%v", t))
	}
	e.Orm.Exec(fmt.Sprintf("DELETE FROM `goods_mark_tag` WHERE `goods_mark_tag`.`goods_id` IN (%v)", strings.Join(removeIds, ",")))
	e.Orm.Exec(fmt.Sprintf("DELETE FROM `goods_mark_class` WHERE `goods_mark_class`.`goods_id` IN (%v)", strings.Join(removeIds, ",")))
	return nil
}
