/**
@Author: chaoqun
* @Date: 2024/1/11 18:26
*/
package service

import (
	"fmt"
	models2 "go-admin/cmd/migrate/migration/models"
	"go-admin/common/utils"
	"strings"
)



func (e *Goods)GetSpecInventory(cid int,key string) (openInventory bool, stock int ) {
	openInventory = IsOpenInventory(cid,e.Orm)
	if !openInventory{
		return
	}

	var Inventory models2.Inventory

	//inventoryKey = append(inventoryKey,fmt.Sprintf("(goods_id = %v and spec_id = %v)",row.GoodsId,row.Id))

	whereSql:=fmt.Sprintf("c_id = %v and %v",cid,key)
	e.Orm.Model(&models2.Inventory{}).Select("id,stock").Where(whereSql).Limit(1).Find(&Inventory)
	if Inventory.Id == 0 {
		return
	}

	stock = Inventory.Stock
	return

}
//查看规格的库存数量 key:[ goods_id = 1 and specs_id = 1]
func (e *Goods)GetBatchSpecInventory(cid int,inventoryKey []string) (openInventory bool,res map[string]int ){

	openInventory = IsOpenInventory(cid,e.Orm)
	if !openInventory{
		return
	}

	res = make(map[string]int,0)
	var InventoryList []models2.Inventory
	inventoryKey = utils.RemoveRepeatStr(inventoryKey)
	//inventoryKey = append(inventoryKey,fmt.Sprintf("(goods_id = %v and spec_id = %v)",row.GoodsId,row.Id))

	whereSql:=fmt.Sprintf("c_id = %v and %v",cid,strings.Join(inventoryKey," or "))
	e.Orm.Model(&models2.Inventory{}).Select("id,stock").Where(whereSql).Find(&InventoryList)
	if len(InventoryList) == 0 {
		return
	}
	for _,row:=range InventoryList{

		key :=fmt.Sprintf("%v_%v",row.GoodsId,row.SpecId)


		res[key]= row.Stock

	}

	return

}

//商品列表展示时,批量获取库存的数据 返回一个map  商品ID:商品下规格所有的数据数量
func (e *Goods)GetBatchGoodsInventory(cid int,goodsId []int) (openInventory bool, res map[int]int ){
	res =make(map[int]int,0)
	openInventory = IsOpenInventory(cid,e.Orm)
	if !openInventory{
		return
	}

	var InventoryList []models2.Inventory

	e.Orm.Model(&models2.Inventory{}).Select("goods_id,stock").Where("c_id = ? and goods_id in ?",cid,goodsId).Find(&InventoryList)
	if len(InventoryList) == 0 {
		return
	}
	for _,row:=range InventoryList{


		data,ok:=res[row.GoodsId]
		if !ok{
			data = row.Stock
		}else {
			data += row.Stock
		}
		res[row.GoodsId] = data

	}

	return

}