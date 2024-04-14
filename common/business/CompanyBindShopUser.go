package business

import (
	"errors"
	"go-admin/cmd/migrate/migration/models"
	"gorm.io/gorm"
)

//用于大B和小B直接绑定关系的函数


type CompanyBindUser struct {
	SiteId int `json:"site_id"` //为了更明显的区分 专门用SiteId来查询
	Orm *gorm.DB
}

//是否有绑定,也就是通过这个函数得知 移动端(H5)携带的站点ID 和当前用户查询 是否有绑定

func (b *CompanyBindUser)IsBind(shopUser int)  bool {
	var object models.CompanyBindShopUser

	b.Orm.Model(&object).Where("c_id = ? and shop_user_id = ?",b.SiteId,shopUser).Limit(1).Find(&object)

	if object.Id == 0 {
		return  false
	}

	return true
}

//增加绑定

func (b *CompanyBindUser)AddBind(shopUser int)  (ok bool,err error) {

	if b.IsBind(shopUser){
		return true,errors.New("已经存在")
	}

	b.Orm.Create(&models.CompanyBindShopUser{
		CId: b.SiteId,
		ShopUserId: shopUser,
		Enable: true,
	})

	return true,nil
}

//解除绑定

func (b *CompanyBindUser)UnBind(shopUser int)  (ok bool,err error) {
	var object models.CompanyBindShopUser
	b.Orm.Model(&object).Where("c_id = ? and shop_user_id = ?",b.SiteId,shopUser).Delete(&object)

	return true,nil
}

