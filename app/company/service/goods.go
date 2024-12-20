package service

import (
	"encoding/json"
	"errors"
	"fmt"
	models2 "go-admin/cmd/migrate/migration/models"
	"go-admin/common/business"
	"go-admin/common/qiniu"
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
)

type Goods struct {
	service.Service

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
	}).Preload("Brand", func(tx *gorm.DB) *gorm.DB {
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
func (e *Goods) getBrandModels(ids []string) (list []models.GoodsBrand) {
	for _, id := range ids {
		var row models.GoodsBrand
		e.Orm.Model(&models.GoodsBrand{}).Where("id = ?", id).First(&row)
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
func (e *Goods) Insert(cid int, c *dto.GoodsInsertReq) (uid int,specDbMap map[int]int, err error) {

	specDbMap = make(map[int]int,0)
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
	//品牌
	if c.Brand  > 0 {
		brandList:=[]string{fmt.Sprintf("%v",c.Brand)}
		data.Brand = e.getBrandModels(brandList)
	}
	err = e.Orm.Create(&data).Error

	if c.Content != ""{
		if utils.CheckStringSize(c.Content){
			e.Orm.Create(&models.GoodsDesc{
				CId: data.CId,
				GoodsId: data.Id,
				Desc: c.Content,
			})
		}
	}

	specsList := make([]dto.Specs, 0)
	marshErr := json.Unmarshal([]byte(c.Specs), &specsList)
	//商品库存值,所以得规格相加
	inventory := 0
	saleAll:=0
	//规格 + vip价格设置存在
	moneyList := make([]float64, 0)
	if len(specsList) > 0 && marshErr == nil {
		for index, row := range specsList {

			stock :=utils.StringToInt(row.Inventory)

			price := utils.RoundDecimalFlot64(row.Price)

			moneyList = append(moneyList, price)

			inventory += stock
			saleAll +=row.VirtuallySale
			specsModels := models.GoodsSpecs{
				Name:    row.Name,
				CId:     cid,
				Enable:  true,
				Layer:   row.Layer,
				GoodsId: data.Id,
				Code:    row.Code,
				VirtuallySale: row.VirtuallySale,
				SerialNumber: row.SerialNumber,
				Image:   row.Image,
				Price:   price,
				Market: utils.RoundDecimalFlot64(row.Market),
				Original:utils.RoundDecimalFlot64(row.Original),
				Inventory: stock,
				UnitId:     utils.StringToInt(row.UnitId),
				Limit: utils.StringToInt(row.Limit),
				Max: utils.StringToInt(row.Max),
			}
			specsModels.CreateBy = data.CreateBy
			e.Orm.Create(&specsModels)
			specDbMap[index] = specsModels.Id
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
		"sale":saleAll,
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
		return 0,specDbMap, err
	}

	return data.Id,specDbMap, err
}

// Update 修改Goods对象
func (e *Goods) Update(cid int,buckClient qiniu.QinUi, c *dto.GoodsUpdateReq, p *actions.DataPermission) (NewSpecImageMap map[int]int,err error ) {
	NewSpecImageMap = make(map[int]int,0)
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
	//fmt.Println("商品的分类",c.Class)
	if len(c.Class) > 0 {
		class := strings.Split(c.Class, ",")
		data.Class = e.getClassModels(class)
		//fmt.Println("获取到的分类的列表",class)
	}
	//品牌
	e.Orm.Model(&data).Association("Brand").Clear()

	if c.Brand  > 0 {
		brandList:=[]string{fmt.Sprintf("%v",c.Brand)}
		data.Brand = e.getBrandModels(brandList)

	}
	//商品库存值
	inventory := 0
	//销售量 虚拟销量 + 实际销量
	saleAll :=0
	specsList := make([]dto.Specs, 0)
	marshErr := json.Unmarshal([]byte(c.Specs), &specsList)
	moneyList := make([]float64, 0)

	//操作的规格列表
	netSpecList:=make([]int,0)

	//原来数据的规格ID
	GoodsSpecsList:=make([]models.GoodsSpecs,0)
	e.Orm.Model(&models.GoodsSpecs{}).Select("id").Where("goods_id = ?",c.Id).Find(&GoodsSpecsList)
	oldSpecList:=make([]int,0)
	for _,spec:=range GoodsSpecsList{
		oldSpecList = append(oldSpecList,spec.Id)
	}

	//{规格ID:图片索引}
	cacheSpecImageMap :=make(map[string]int,0)
	json.Unmarshal([]byte(fmt.Sprintf("%v",c.SpecImageMap)),&cacheSpecImageMap)

	//规格处理
	if len(specsList) > 0 && marshErr == nil {

		//这个是专门用来记录, 文件的位置
		fileIndex:=0
		for index, row := range specsList {

			var specsRow models.GoodsSpecs
			//获取库存量
			stock := utils.StringToInt(row.Inventory)
			price := utils.RoundDecimalFlot64(row.Price)
			moneyList = append(moneyList, price)
			inventory += stock
			//虚拟销量的叠加
			saleAll +=row.VirtuallySale

			if row.Type !="append" {
				//就是一个规格资源的更新
				e.Orm.Model(&models.GoodsSpecs{}).Where("id = ?", row.Id).Limit(1).Find(&specsRow)
				if specsRow.Id == 0 {
					continue
				}
				//原来的销量也叠加上去
				saleAll += specsRow.Sale

				specsRow.VirtuallySale = row.VirtuallySale
				specsRow.Code = row.Code
				specsRow.SerialNumber = row.SerialNumber
				specsRow.Name = row.Name
				specsRow.Enable = row.Enable
				specsRow.Layer = row.Layer
				specsRow.Price = price
				specsRow.Market = utils.RoundDecimalFlot64(row.Market)
				specsRow.Original = utils.RoundDecimalFlot64(row.Original)
				specsRow.Inventory = stock
				specsRow.UnitId = utils.StringToInt(row.UnitId)
				specsRow.Limit = utils.StringToInt(row.Limit)
				specsRow.Max = utils.StringToInt(row.Max)
				netSpecList = append(netSpecList,specsRow.Id)
				e.Orm.Save(&specsRow)


			} else {
				//规格资源的创建
				specsRow = models.GoodsSpecs{
					Name:    row.Name,
					CId:     cid,
					Enable:  true,
					Layer:   row.Layer,
					GoodsId: c.Id,
					VirtuallySale: row.VirtuallySale,
					Code:    row.Code,
					SerialNumber: row.SerialNumber,
					Price:   price,
					Original:utils.RoundDecimalFlot64(row.Original),
					Market: utils.RoundDecimalFlot64(row.Market),
					Inventory: stock,
					UnitId:     utils.StringToInt(row.UnitId),
					Limit: utils.StringToInt(row.Limit),
					Max: utils.StringToInt(row.Max),
				}
				specsRow.CreateBy = data.CreateBy
				e.Orm.Create(&specsRow)
				netSpecList = append(netSpecList,specsRow.Id)

			}
			//这个是新增的一个规格，那这个规格 判断下是否有规格图片
			//cacheSpecImageMap:{规格ID:属于第几个规格}
			for _,specIndex:=range cacheSpecImageMap{
				if specIndex == index{
					//保存到新的这个组中
					NewSpecImageMap[specsRow.Id] =  fileIndex
					//这样才算是一个文件
					fileIndex ++
				}
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

		diffList := utils.DifferenceInt(oldSpecList, netSpecList)


		//fmt.Println("有差别的规格,",diffList)
		//fmt.Println("新的一个map映射",NewSpecImageMap)
		for _,specId:=range diffList{
			var goodsSpec models.GoodsSpecs
			e.Orm.Model(&models.GoodsSpecs{}).Select("id,image").Where("id = ?",specId).Limit(1).Find(&goodsSpec)
			if goodsSpec.Id > 0 && goodsSpec.Image != ""{
				fmt.Println("删除有差别规格的图片",goodsSpec.Image)
				buckClient.RemoveFile(business.GetSiteCosPath(cid,global.GoodsPath,goodsSpec.Image))
			}
			e.Orm.Model(&models.GoodsSpecs{}).Unscoped().Where("id = ?",specId).Delete(&models.GoodsSpecs{})
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
	data.Sale = saleAll
	//fmt.Println("新的库存是，",inventory)
	db := e.Orm.Save(&data)
	if err = db.Error; err != nil {
		e.Log.Errorf("GoodsService Save error:%s \r\n", err)
		return nil, err
	}
	if c.Content != ""{
		if utils.CheckStringSize(c.Content){
			var GoodsDesc models.GoodsDesc
			var count int64
			e.Orm.Model(&GoodsDesc).Where("goods_id = ?",data.Id).Count(&count)
			if count == 0 {
				e.Orm.Create(&models.GoodsDesc{
					CId: data.CId,
					GoodsId: data.Id,
					Desc: c.Content,
				})
			}else {
				e.Orm.Model(&GoodsDesc).Where("goods_id = ?",data.Id).Updates(map[string]interface{}{
					"desc":c.Content,
					"c_id":data.CId,
				})
			}
		}

	}
	return NewSpecImageMap, err
}

// Remove 删除Goods
func (e *Goods) Remove(d *dto.GoodsDeleteReq,cid int, p *actions.DataPermission) (notDelete,okDelete []string) {

	removeIds := make([]string, 0)

	buckClient:=qiniu.QinUi{
		CId: cid,
	}
	buckClient.InitClient()

	notDelete =make([]string,0)
	okDelete = make([]string,0)

	isOpenInventory := IsOpenInventory(cid,e.Orm)

	for _, goodsId := range d.Ids {
		removeFileList:=make([]string,0)
		var goods models.Goods

		//删除商品
		e.Orm.Model(&goods).Where("id = ?",goodsId).Limit(1).Find(&goods)

		//有库存管理
		if isOpenInventory {
			var InventoryList []models2.Inventory
			e.Orm.Model(&models2.Inventory{}).Select("stock").Where("c_id = ? and goods_id = ?",cid,goodsId).Find(&InventoryList)
			allNumber :=0
			for _,row:=range InventoryList{
				allNumber+=row.Stock
			}
			if allNumber > 0 {
				notDelete = append(notDelete,goods.Name)
				continue
			}

		}else {
			//无库存管理
			//库存不为0 不可删除
			if goods.Inventory > 0 {
				notDelete = append(notDelete,goods.Name)
				continue
			}

		}
		//可以删除的时候 才会保留删除ID
		removeIds = append(removeIds, fmt.Sprintf("%v", goodsId))

		//商品的主图
		removeFileList = append(removeFileList, strings.Split(goods.Image,",")...)
		//删除规格图片
		GoodsSpecsList:=make([]models.GoodsSpecs,0)

		e.Orm.Model(&models.GoodsSpecs{}).Select("image").Where("goods_id = ?",goodsId).Find(&GoodsSpecsList)

		for _,spec:=range GoodsSpecsList{
			if spec.Image !="" {
				removeFileList = append(removeFileList,spec.Image)
			}
		}
		fmt.Println("删除商品+规格图片",removeFileList)
		//如果有图片,删除图片
		for _,image :=range removeFileList{
			//_ = os.Remove(business.GetGoodPathName(goods.CId) + image)
			buckClient.RemoveFile(business.GetSiteCosPath(cid,global.GoodsPath,image))
		}
		//如果有商品详细,那就匹配图片路径
		var goodsDesc models.GoodsDesc
		e.Orm.Model(&goodsDesc).Where("goods_id = ?",goodsId).Limit(1).Find(&goodsDesc)
		if goodsDesc.Desc != ""{
			//text:="<p><img src=\"https://dcy-1318497773.cos.ap-nanjing.myqcloud.com/goods/1/088e54a8.jpg\"></p>"
			reImg :=`https?://[^"]+?(\.((jpg)|(png)|(jpeg)|(gif)|(bmp)))`

			re:=regexp.MustCompile(reImg)
			result:=re.FindAllString(goodsDesc.Desc,-1)
			for _,image :=range result{
				cosImagePath :=business.GetDomainSplitFilePath(image)
				buckClient.RemoveFile(cosImagePath)
			}
		}
		//删除商品
		e.Orm.Model(&goods).Where("id = ?",goodsId).Delete(&models.Goods{})

		okDelete = append(okDelete,goods.Name)

	}
	if len(removeIds) > 0 {
		e.Orm.Model(&models.GoodsDesc{}).Where("goods_id in ?", removeIds).Unscoped().Delete(&models.GoodsDesc{})
		//商品删除了关联的一些配置都删除
		e.Orm.Model(&models.GoodsVip{}).Where("goods_id in ?", removeIds).Unscoped().Delete(&models.GoodsVip{})
		e.Orm.Model(&models.GoodsSpecs{}).Where("goods_id in ?", removeIds).Unscoped().Delete(&models.GoodsSpecs{})
		e.Orm.Exec(fmt.Sprintf("DELETE FROM `goods_mark_tag` WHERE `goods_mark_tag`.`goods_id` IN (%v)", strings.Join(removeIds, ",")))
		e.Orm.Exec(fmt.Sprintf("DELETE FROM `goods_mark_class` WHERE `goods_mark_class`.`goods_id` IN (%v)", strings.Join(removeIds, ",")))
		e.Orm.Exec(fmt.Sprintf("DELETE FROM `goods_mark_brand` WHERE `goods_mark_brand`.`goods_id` IN (%v)", strings.Join(removeIds, ",")))

		//删除库存
		if isOpenInventory {
			e.Orm.Model(&models2.Inventory{}).Where("goods_id in ?",removeIds).Unscoped().Delete(&models2.Inventory{})
		}
	}

	return
}
