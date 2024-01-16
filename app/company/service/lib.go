package service

import (
	sys "go-admin/app/admin/models"
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
//检测是否开启了 订单审批 + 用户有订单权限
func IsHasOpenApprove(user *sys.SysUser,orm *gorm.DB) (openApprove,hasApprove bool) {
	var Approve models2.OrderApproveCnf
	orm.Model(&Approve).Where("c_id = ?",user.CId).Limit(1).Find(&Approve)

	if Approve.Id == 0 {
		return false,false
	}
	return Approve.Enable,user.AuthExamine
}