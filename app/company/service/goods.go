package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-admin/common/utils"
	"go-admin/global"
	"strconv"
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

	query := e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).Order(global.OrderLayerKey).Preload("Class", func(tx *gorm.DB) *gorm.DB {
		return tx.Select("id,name")
	})

	if c.Class != "" {
		query = query.Joins("LEFT JOIN goods_mark_class ON goods.id = goods_mark_class.goods_id").Where("goods_mark_class.class_id in ?",
			strings.Split(c.Class, ","))
	}
	err = query.Find(list).Limit(-1).Offset(-1).
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
		).Preload("Tag", func(tx *gorm.DB) *gorm.DB {
		return tx.Select("id,name")
	}).Preload("Class", func(tx *gorm.DB) *gorm.DB {
		return tx.Select("id,name")
	}).
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

func (e *Goods) getTagModels(ids []string) (list []models.GoodsTag) {
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
func (e *Goods) getClassModels(ids []string) (list []models.GoodsClass) {
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
		tags := strings.Split(c.Tag, ",")
		data.Tag = e.getTagModels(tags)
	}
	//分类
	if len(c.Class) > 0 {
		class := strings.Split(c.Class, ",")
		data.Class = e.getClassModels(class)
	}
	err = e.Orm.Create(&data).Error

	specsList := make([]dto.Specs, 0)
	marshErr := json.Unmarshal([]byte(c.Specs), &specsList)
	//商品库存值,所以得规格相加
	inventory := 0
	//规格 + vip价格设置存在
	moneyList := make([]float64, 0)
	if len(specsList) > 0 && marshErr == nil {
		for _, row := range specsList {
			stock := func() int {

				n, _ := strconv.Atoi(fmt.Sprintf("%v", row.Inventory))
				return n
			}()
			price := func() float64 {

				if row.Price == "" {
					return 0
				}
				n, _ := strconv.ParseFloat(fmt.Sprintf("%v", row.Price), 64)
				return n
			}()
			moneyList = append(moneyList, price)
			inventory += stock
			specsModels := models.GoodsSpecs{
				Name:    row.Name,
				CId:     cid,
				Enable:  row.Enable,
				Layer:   row.Layer,
				GoodsId: data.Id,
				Code:    row.Code,
				Price:   price,
				Original: func() float64 {

					if row.Original == "" {
						return 0
					}
					n, _ := strconv.ParseFloat(fmt.Sprintf("%v", row.Original), 64)
					return n
				}(),
				Inventory: stock,
				Unit:      row.Unit,
				Limit: func() int {

					n, _ := strconv.Atoi(fmt.Sprintf("%v", row.Limit))
					return n
				}(),
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
					CId:     cid,
					GoodsId: data.Id,
					SpecsId: specsModels.Id,
					Enable:  vipEnable,
					Layer:   0,
					GradeId: gradeRow.Id,
					CustomPrice: func() float64 {
						if v == "" {
							return 0
						}
						n, _ := strconv.ParseFloat(fmt.Sprintf("%v", v), 64)
						return n
					}(),
				}
				vipRow.CreateBy = data.CreateBy
				e.Orm.Create(&vipRow)
			}
		}
	} else {
		fmt.Println("规格数据序列化失败", marshErr)
	}

	if err != nil {
		e.Log.Errorf("GoodsService Insert error:%s \r\n", err)
		return 0, err
	}

	e.Orm.Model(&data).Where("id = ?", data.Id).Updates(map[string]interface{}{
		"inventory": inventory,
		"money": func() string {
			if len(moneyList) > 0 {
				if len(moneyList) == 1 {
					return fmt.Sprintf("¥%v", moneyList[0])
				}
				min, max := utils.MinAndMax(moneyList)
				return fmt.Sprintf("¥%v-%v", min, max)
			}
			return ""
		}(),
	})
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
		tags := strings.Split(c.Tag, ",")
		data.Tag = e.getTagModels(tags)
	}
	//分类
	e.Orm.Model(&data).Association("Class").Clear()
	if len(c.Class) > 0 {
		class := strings.Split(c.Class, ",")
		data.Class = e.getClassModels(class)
	}
	specsList := make([]dto.Specs, 0)
	marshErr := json.Unmarshal([]byte(c.Specs), &specsList)
	//商品库存值,所以得规格相加
	inventory := 0
	moneyList := make([]float64, 0)
	//规格更新
	if len(specsList) > 0 && marshErr == nil {
		for _, row := range specsList {
			var specsRow models.GoodsSpecs
			//获取库存量
			stock := func() int {
				n, _ := strconv.Atoi(fmt.Sprintf("%v", row.Inventory))
				return n
			}()
			price := func() float64 {

				if row.Price == "" {
					return 0
				}
				n, _ := strconv.ParseFloat(fmt.Sprintf("%v", row.Price), 64)
				return n
			}()
			moneyList = append(moneyList, price)
			inventory += stock
			if row.Id > 0 {
				//就是一个规格资源的更新
				e.Orm.Model(&models.GoodsSpecs{}).Where("id = ?", row.Id).Limit(1).Find(&specsRow)
				if specsRow.Id == 0 {
					continue
				}
				specsRow.Code = row.Code
				specsRow.Name = row.Name
				specsRow.Enable = row.Enable
				specsRow.Layer = row.Layer
				specsRow.Price = price
				specsRow.Original = func() float64 {

					if row.Original == "" {
						return 0
					}
					n, _ := strconv.ParseFloat(fmt.Sprintf("%v", row.Original), 64)
					return n
				}()
				specsRow.Inventory = stock
				specsRow.Unit = row.Unit
				specsRow.Limit = func() int {
					n, _ := strconv.Atoi(fmt.Sprintf("%v", row.Limit))
					return n
				}()
				e.Orm.Save(&specsRow)
			} else {
				//规格资源的创建
				specsRow = models.GoodsSpecs{
					Name:    row.Name,
					CId:     cid,
					Enable:  row.Enable,
					Layer:   row.Layer,
					GoodsId: c.Id,
					Code:    row.Code,
					Price:   price,
					Original: func() float64 {

						if row.Price == "" {
							return 0
						}
						n, _ := strconv.ParseFloat(fmt.Sprintf("%v", row.Original), 64)
						return n
					}(),
					Inventory: stock,
					Unit:      row.Unit,
					Limit: func() int {
						n, _ := strconv.Atoi(fmt.Sprintf("%v", row.Limit))
						return n
					}(),
				}
				specsRow.CreateBy = data.CreateBy
				e.Orm.Create(&specsRow)
			}
			var vipEnable bool
			for k, v := range row.Vip {
				if k == "enable" {
					vipEnable = v.(bool)
				}
				if !strings.HasPrefix(k, "vip_") {
					continue
				}
				gradeInt := strings.Replace(k, "vip_", "", -1)

				fmt.Println("gradeIntgradeInt", gradeInt, "k", k)
				var gradeRow models.GradeVip
				e.Orm.Model(&models.GradeVip{}).Where("enable = ? and id = ? and c_id = ?",
					true, gradeInt, cid).Limit(1).Find(&gradeRow)
				if gradeRow.Id == 0 {
					continue
				}
				var goodVipRow models.GoodsVip
				e.Orm.Model(&goodVipRow).Where("goods_id = ? and specs_id = ? and grade_id =?",
					specsRow.GoodsId, specsRow.Id, gradeRow.Id).Limit(1).Find(&goodVipRow)
				vipRow := models.GoodsVip{
					CId:     cid,
					GoodsId: data.Id,
					SpecsId: specsRow.Id,
					Enable:  vipEnable,
					Layer:   0,
					GradeId: gradeRow.Id,
					CustomPrice: func() float64 {
						if v == "" {
							return 0
						}
						n, _ := strconv.ParseFloat(fmt.Sprintf("%v", v), 64)
						return n
					}(),
				}
				if goodVipRow.Id == 0 {
					vipRow.CreateBy = data.CreateBy
					e.Orm.Create(&vipRow)
				} else {
					e.Orm.Model(&goodVipRow).Updates(&vipRow)
				}
			}
		}

	} else {
		fmt.Println("规格数据序列化失败", marshErr)
	}
	data.Money = func() string {
		if len(moneyList) > 0 {
			if len(moneyList) == 1 {
				return fmt.Sprintf("¥%v", moneyList[0])
			}
			min, max := utils.MinAndMax(moneyList)
			return fmt.Sprintf("¥%v-%v", min, max)
		}

		return ""
	}()
	data.Inventory = inventory
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
	//商品删除了关联的一些配置都删除
	e.Orm.Model(&models.GoodsVip{}).Where("goods_id in ?", removeIds).Unscoped().Delete(&models.GoodsVip{})
	e.Orm.Model(&models.GoodsSpecs{}).Where("goods_id in ?", removeIds).Unscoped().Delete(&models.GoodsSpecs{})
	e.Orm.Exec(fmt.Sprintf("DELETE FROM `goods_mark_tag` WHERE `goods_mark_tag`.`goods_id` IN (%v)", strings.Join(removeIds, ",")))
	e.Orm.Exec(fmt.Sprintf("DELETE FROM `goods_mark_class` WHERE `goods_mark_class`.`goods_id` IN (%v)", strings.Join(removeIds, ",")))
	return nil
}
