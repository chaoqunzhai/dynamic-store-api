package service

import (
	models2 "go-admin/cmd/migrate/migration/models"
	"gorm.io/gorm"
)

//检测是否开启了库存
func IsOpenInventory(cid int,orm *gorm.DB) bool{
	var Inventory models2.InventoryCnf
	orm.Model(&models2.InventoryCnf{}).Select("id,enable").Where("c_id = ?",cid).Limit(1).Find(&Inventory)
	if Inventory.Id == 0 {
		return false
	}
	return Inventory.Enable
}
