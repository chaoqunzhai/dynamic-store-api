package service

import (
	sys "go-admin/app/admin/models"
	models2 "go-admin/cmd/migrate/migration/models"
	"gorm.io/gorm"
	"time"
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

//检测线路是否过期
func CheckLineExpire(cid,lineId int,orm *gorm.DB) (msg string,ExpiredOrNot bool) {
	var line models2.Line
	orm.Model(&line).Where("c_id = ? and id = ?",cid,lineId).Limit(1).Find(&line)

	if line.Id == 0 {
		return "暂无路线",false
	}
	if !line.ExpirationTime.Time.IsZero() { //有时间配置
		if line.ExpirationTime.Before(time.Now()){
			return "路线已过期",false
		}
	}
	return "路线可用",true
}
