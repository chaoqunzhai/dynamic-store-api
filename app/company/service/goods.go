package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-admin/common/business"
	"go-admin/common/tx_api"
	"go-admin/common/utils"
	"go-admin/global"
	"go.uber.org/zap"
	"regexp"
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

	if c.Content != ""{
		e.Orm.Create(&models.GoodsDesc{
			GoodsId: data.Id,
			Desc: c.Content,
		})
	}

	specsList := make([]dto.Specs, 0)
	marshErr := json.Unmarshal([]byte(c.Specs), &specsList)
	//商品库存值,所以得规格相加
	inventory := 0
	//规格 + vip价格设置存在
	moneyList := make([]float64, 0)
	if len(specsList) > 0 && marshErr == nil {
		for _, row := range specsList {

			stock :=utils.StringToInt(row.Inventory)

			price := utils.RoundDecimalFlot64(row.Price)

			moneyList = append(moneyList, price)

			inventory += stock

			specsModels := models.GoodsSpecs{
				Name:    row.Name,
				CId:     cid,
				Enable:  row.Enable,
				Layer:   row.Layer,
				GoodsId: data.Id,
				Code:    row.Code,
				Image:   row.Image,
				Price:   price,
				Market: utils.RoundDecimalFlot64(row.Market),
				Original:utils.RoundDecimalFlot64(row.Original),
				Inventory: stock,
				Unit:      row.Unit,
				Limit: utils.StringToInt(row.Limit),
				Max: utils.StringToInt(row.Max),
			}
			specsModels.CreateBy = data.CreateBy
			e.Orm.Create(&specsModels)

			for k, v := range row.Vip {

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
					Enable:  c.VipSale,
					Layer:   0,
					GradeId: gradeRow.Id,
					CustomPrice: utils.StringToFloat64(v),
				}
				vipRow.CreateBy = data.CreateBy
				e.Orm.Create(&vipRow)
			}
		}
	} else {
		zap.S().Errorf("创建商品规格序列化失败%v",marshErr)
	}
	e.Orm.Model(&data).Where("id = ?", data.Id).Updates(map[string]interface{}{
		"inventory": inventory,
		"money": func() string {
			if len(moneyList) > 0 {
				if len(moneyList) == 1 {
					firstMoney := moneyList[0]
					return fmt.Sprintf("¥%v", utils.StringDecimal(firstMoney))
				}
				min1, max2 := utils.MinAndMax(moneyList)
				if min1 == max2 {

					return fmt.Sprintf("¥%v", utils.StringDecimal(min1))
				}
				return fmt.Sprintf("¥%v-%v", utils.StringDecimal(min1), utils.StringDecimal(max2))
			}
			return ""
		}(),
	})

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
		tags := strings.Split(c.Tag, ",")
		data.Tag = e.getTagModels(tags)
	}
	//分类
	e.Orm.Model(&data).Association("Class").Clear()
	if len(c.Class) > 0 {
		class := strings.Split(c.Class, ",")
		data.Class = e.getClassModels(class)
	}
	//商品库存值
	inventory := 0

	specsList := make([]dto.Specs, 0)
	marshErr := json.Unmarshal([]byte(c.Specs), &specsList)
	moneyList := make([]float64, 0)
	//规格更新
	if len(specsList) > 0 && marshErr == nil {
		for _, row := range specsList {
			var specsRow models.GoodsSpecs
			//获取库存量
			stock := utils.StringToInt(row.Inventory)
			price := utils.RoundDecimalFlot64(row.Price)
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
				specsRow.Market = utils.RoundDecimalFlot64(row.Market)
				specsRow.Original = utils.RoundDecimalFlot64(row.Original)
				specsRow.Inventory = stock
				specsRow.Unit = row.Unit
				specsRow.Limit = utils.StringToInt(row.Limit)
				specsRow.Max = utils.StringToInt(row.Max)

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
					Original:utils.RoundDecimalFlot64(row.Original),
					Inventory: stock,
					Unit:      row.Unit,
					Limit: utils.StringToInt(row.Limit),
					Max: utils.StringToInt(row.Max),
				}
				specsRow.CreateBy = data.CreateBy
				e.Orm.Create(&specsRow)
			}

			for k, v := range row.Vip {

				if !strings.HasPrefix(k, "vip_") {
					continue
				}
				gradeInt := strings.Replace(k, "vip_", "", -1)


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
					Enable:  c.VipSale,
					Layer:   0,
					GradeId: gradeRow.Id,
					CustomPrice: utils.RoundDecimalFlot64(v),
				}
				if goodVipRow.Id == 0 {
					vipRow.CreateBy = data.CreateBy
					e.Orm.Create(&vipRow)
				} else {
					e.Orm.Model(&goodVipRow).Updates(&vipRow)
					e.Orm.Model(&models.GoodsVip{}).Where("id = ?",goodVipRow.Id).Updates(map[string]interface{}{
						"enable":c.VipSale,
					})

				}
			}
		}

	} else {
		zap.S().Errorf("更新商品规格序列化失败%v",marshErr)
	}
	data.Money = func() string {
		if len(moneyList) > 0 {
			if len(moneyList) == 1 {
				firstMoney := moneyList[0]
				return fmt.Sprintf("¥%v", utils.StringDecimal(firstMoney))
			}
			min1, max2 := utils.MinAndMax(moneyList)
			if min1 == max2 {

				return fmt.Sprintf("¥%v", utils.StringDecimal(min1))
			}
			return fmt.Sprintf("¥%v-%v", utils.StringDecimal(min1), utils.StringDecimal(max2))
		}
		return ""
	}()

	data.Inventory = inventory
	db := e.Orm.Save(&data)
	if err = db.Error; err != nil {
		e.Log.Errorf("GoodsService Save error:%s \r\n", err)
		return err
	}
	if c.Content != ""{
		var GoodsDesc models.GoodsDesc
		var count int64
		e.Orm.Model(&GoodsDesc).Where("goods_id = ?",data.Id).Count(&count)
		if count == 0 {
			e.Orm.Create(&models.GoodsDesc{
				GoodsId: data.Id,
				Desc: c.Content,
			})
		}else {
			e.Orm.Model(&GoodsDesc).Where("goods_id = ?",data.Id).Updates(map[string]interface{}{
				"desc":c.Content,
			})
		}
	}
	return nil
}

// Remove 删除Goods
func (e *Goods) Remove(d *dto.GoodsDeleteReq,CId interface{}, p *actions.DataPermission) error {

	removeIds := make([]string, 0)
	txClient :=tx_api.TxCos{}
	txClient.InitClient()
	for _, t := range d.Ids {
		removeIds = append(removeIds, fmt.Sprintf("%v", t))

		var goods models.Goods

		//删除商品
		e.Orm.Model(&goods).Select("image").Where("id = ?",t).Limit(1).Find(&goods)

		//如果有图片,删除图片
		for _,image :=range strings.Split(goods.Image,","){
			//_ = os.Remove(business.GetGoodPathName(goods.CId) + image)
			txClient.RemoveFile(business.GetSiteGoodsPath(CId,image))
		}
		//如果有商品详细,那就匹配图片路径
		var goodsDesc models.GoodsDesc
		e.Orm.Model(&goodsDesc).Where("goods_id = ?",t).Limit(1).Find(&goodsDesc)
		if goodsDesc.Desc != ""{
			//text:="<p><img src=\"https://dcy-1318497773.cos.ap-nanjing.myqcloud.com/goods/1/088e54a8.jpg\"></p>"
			reImg :=`https?://[^"]+?(\.((jpg)|(png)|(jpeg)|(gif)|(bmp)))`

			re:=regexp.MustCompile(reImg)
			result:=re.FindAllString(goodsDesc.Desc,-1)
			for _,image :=range result{
				cosImagePath :=business.GetDomainSplitFilePath(image)
				txClient.RemoveFile(cosImagePath)
			}
		}
		//删除商品
		e.Orm.Model(&goods).Where("id = ?",t).Delete(&models.Goods{})


	}
	e.Orm.Model(&models.GoodsDesc{}).Where("goods_id in ?", removeIds).Unscoped().Delete(&models.GoodsDesc{})
	//商品删除了关联的一些配置都删除
	e.Orm.Model(&models.GoodsVip{}).Where("goods_id in ?", removeIds).Unscoped().Delete(&models.GoodsVip{})
	e.Orm.Model(&models.GoodsSpecs{}).Where("goods_id in ?", removeIds).Unscoped().Delete(&models.GoodsSpecs{})
	e.Orm.Exec(fmt.Sprintf("DELETE FROM `goods_mark_tag` WHERE `goods_mark_tag`.`goods_id` IN (%v)", strings.Join(removeIds, ",")))
	e.Orm.Exec(fmt.Sprintf("DELETE FROM `goods_mark_class` WHERE `goods_mark_class`.`goods_id` IN (%v)", strings.Join(removeIds, ",")))



	return nil
}
