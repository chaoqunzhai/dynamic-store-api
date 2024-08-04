/**
@Author: chaoqun
* @Date: 2024/8/4 22:49
*/
package apis

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"
	"go-admin/app/company/models"
	models2 "go-admin/cmd/migrate/migration/models"
	"go-admin/common/jwt/user"
	"go-admin/global"
)

type DtoGoodsRow struct {
	
	Name string `json:"name"`
}
// 精简API 快速返回全部资源数据

type MiniApi struct {
	api.Api
}




func (e MiniApi) CustomerUser(c *gin.Context) {
	err := e.MakeContext(c).
		MakeOrm().
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	userDto, err := user.GetUserDto(e.Orm, c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}
	var shopList []models2.Shop
	e.Orm.Model(&models2.Shop{}).Select("id,name").Where("c_id = ? ",userDto.CId).Find(&shopList)
	var result []map[string]interface{}

	for _,row:=range shopList{
		result = append(result, map[string]interface{}{
			"id":row.Id,
			"label":row.Name,
		})
	}

	e.OK(result,"")
	return
}

func (e MiniApi) GoodsSpec(c *gin.Context) {
	err := e.MakeContext(c).
		MakeOrm().
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	userDto, err := user.GetUserDto(e.Orm, c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}
	var goodsList []models.Goods
	e.Orm.Model(&models.Goods{}).Select("id,name").
		Where("c_id = ?",userDto.CId).Order(global.OrderLayerKey).Find(&goodsList)
	goodsMap:=make(map[int]DtoGoodsRow,0)
	goodsId:=make([]int,0)
	for _,row:=range goodsList{
		goodsMap[row.Id] = DtoGoodsRow{
			Name: row.Name,
		}
		goodsId = append(goodsId,row.Id)
	}
	var goodsSpecsList []models.GoodsSpecs
	e.Orm.Model(&models.GoodsSpecs{}).Select("goods_id,name,id").
		Where("c_id = ? and goods_id in ?",userDto.CId,goodsId).Find(&goodsSpecsList)

	result:=make([]map[string]interface{},0)

	for _,row:=range goodsSpecsList{
		goodsInfo,ok:=goodsMap[row.GoodsId]
		if !ok{continue}
		dto:=map[string]interface{}{
			"goods_id":row.GoodsId,
			"specs_id":row.Id,
			"label":fmt.Sprintf("%v-%v",goodsInfo.Name,row.Name),
		}
		result = append(result,dto)
	}

	e.OK(result,"")
	return
}
