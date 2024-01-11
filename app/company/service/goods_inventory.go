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

//检测是否开启了库存
func (e *Goods)isOpenInventory(cid int) bool{
	var Inventory models2.InventoryCnf
	e.Orm.Model(&models2.InventoryCnf{}).Select("id,enable").Where("c_id = ?",cid).Limit(1).Find(&Inventory)
	if Inventory.Id == 0 {
		return false
	}
	return Inventory.Enable
}

//创建/编辑商品时,同步修改库存的商品名称 和规格名称
func (e *Goods)SyncInventoryName(cid int,goodsId,specId int,goodsName,specsName string)  {

	if !e.isOpenInventory(cid){
		return
	}
	var Inventory models2.Inventory
	whereSql :=fmt.Sprintf("c_id = %v and goods_id = %v and spec_id = %v",cid,goodsId,specId)
	e.Orm.Model(&models2.Inventory{}).Where(whereSql).Limit(1).Find(&Inventory)
	if Inventory.Id == 0 {
		return
	}
	updateMap:=make(map[string]interface{},0)
	if Inventory.GoodsName != goodsName{
		updateMap["goods_name"] = goodsName
	}
	if Inventory.GoodsSpecName != specsName{
		updateMap["goods_spec_name"] = specsName
	}

	e.Orm.Model(&models2.Inventory{}).Where(whereSql).Updates(updateMap)

}

func (e *Goods)GetSpecInventory(cid int,key string) (openInventory bool, stock int ) {
	openInventory = e.isOpenInventory(cid)
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

	openInventory = e.isOpenInventory(cid)
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
	openInventory = e.isOpenInventory(cid)
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