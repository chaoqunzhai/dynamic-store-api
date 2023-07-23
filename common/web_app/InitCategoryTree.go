/**
@Author: chaoqun
* @Date: 2023/7/22 22:52
*/
package web_app

import (
	"go-admin/cmd/migrate/migration/models"
	"go-admin/common/redis_db"
	"gorm.io/gorm"
)

type MakeWeAppGoodsTree struct {
	CId       int
	Orm *gorm.DB

}

func (m *MakeWeAppGoodsTree) SetLoadRedis() map[string]interface{} {

	dat := map[string]interface{}{

	}
	return dat
}

//查询到数据加载到大BRedis中
func (m *MakeWeAppGoodsTree) SearchRun() []map[string]interface{} {
	//查询大B下所有商品分类 + 分类下的商品列表

	//查询大B下所有的分类
	goodsClass :=make([]models.GoodsClass,0)
	m.Orm.Model(&models.GoodsClass{}).Select("id,name").Where("c_id = ? and enable = ? ",m.CId,true).Order("layer desc").Find(&goodsClass)

	dat :=make([]map[string]interface{},0)
	for _,row:=range goodsClass {
		//查询分类
		category:=map[string]interface{}{
			"category_id":row.Id,
			"category_name":row.Name,
			"image":row.Image,
			"child_list":make([]string,0), //分类下的细分种类,根据业务扩展
		}
		dat = append(dat,category)
	}

	redis_db.SetGoodsCategoryTree(m.CId,dat)
	return dat
}