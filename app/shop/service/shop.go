package service

import (
	"errors"
	"fmt"
	"github.com/go-admin-team/go-admin-core/sdk/service"
	sys "go-admin/app/admin/models"
	"go-admin/app/shop/models"
	"go-admin/app/shop/service/dto"
	models2 "go-admin/cmd/migrate/migration/models"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
	"go-admin/global"
	"gorm.io/gorm"
)

type Shop struct {
	service.Service
}

// GetPage 获取Shop列表
func (e *Shop) GetPage(c *dto.ShopGetPageReq, p *actions.DataPermission, list *[]models.Shop, count *int64) error {
	var err error
	var data models.Shop

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).Order(global.OrderLayerKey).Preload("Tag", func(tx *gorm.DB) *gorm.DB {
		return tx.Select("id", "name")
	}).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("ShopService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取Shop对象
func (e *Shop) Get(d *dto.ShopGetReq, p *actions.DataPermission, model *models.Shop) error {
	var data models.Shop

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Preload("Tag", func(tx *gorm.DB) *gorm.DB {
		return tx.Select("id", "name")
	}).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetShop error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}
func (e *Shop) getShopTagModels(ids []int) (list []models.ShopTag) {
	for _, id := range ids {
		var row models.ShopTag
		e.Orm.Model(&models.ShopTag{}).Where("id = ?", id).First(&row)
		if row.Id == 0 {
			continue
		}
		list = append(list, row)
	}
	return list
}

// Insert 创建Shop对象
func (e *Shop) Insert(userDto *sys.SysUser, c *dto.ShopInsertReq) error {
	//创建小B用户
	shopUserDto:=sys.SysUser{
		Username: c.UserName,
		NickName: c.UserName,
		Phone: c.Phone,
		Password: c.Password,
		Enable: true,
		CId: userDto.CId,
		Status:fmt.Sprintf("%v", global.SysUserSuccess),
		RoleId:global.RoleShop,
	}
	//设置创建用户为大B的名字
	shopUserDto.CreateBy = userDto.UserId

	e.Orm.Create(&shopUserDto)
	var err error
	var data models.Shop
	c.Generate(&data)
	//把创建的用户关联到这个小B上面来
	data.UserId = shopUserDto.UserId
	//关联的ID
	data.CId = userDto.CId

	if len(c.Tags) > 0 {
		data.Tag = e.getShopTagModels(c.Tags)
	}
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("ShopService Insert error:%s \r\n", err)
		return err
	}


	//创建一个小B的一个默认地址,为什么要重复创建2份一样的呢,是为了区分开,方便用户自己去编辑改地址
	//大B创建的只一份
	userAddress:=models2.DynamicUserAddress{
		Source: 1,
		Name: data.Name,
		Address: data.Address,
		Mobile: data.Phone,
		IsDefault: true,
	}
	//存储的是创建小B的用户ID,因为这个地址是小B的
	userAddress.UserId = shopUserDto.UserId
	userAddress.CId = userDto.CId
	e.Orm.Create(&userAddress)

	return nil
}

// Update 修改Shop对象
func (e *Shop) Update(c *dto.ShopUpdateReq, p *actions.DataPermission) error {
	var err error
	var data = models.Shop{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).First(&data, c.GetId())
	c.Generate(&data)
	//清除关联
	_=e.Orm.Model(&data).Association("Tag").Clear()
	if len(c.Tags) > 0 {
		//增加关联
		//fmt.Println("标签", c.Tags)
		data.Tag = e.getShopTagModels(c.Tags)
	}
	db := e.Orm.Save(&data)
	if err = db.Error; err != nil {
		e.Log.Errorf("ShopService Save error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	var userAddress models2.DynamicUserAddress

	var isDefault int64
	e.Orm.Model(&models2.DynamicUserAddress{}).Where("user_id = ? and is_default = 1 and c_id = ?",data.UserId,data.CId).Count(&isDefault)
	//用户把大B创建的地址删了
	e.Orm.Model(&models2.DynamicUserAddress{}).Where("user_id = ? and source = 1 and c_id = ?",data.UserId,data.CId).Limit(1).Find(&userAddress)
	if userAddress.Id == 0  {
		address:=models2.DynamicUserAddress{
			Source: 1,
			Name: data.Name,
			Address: data.Address,
			Mobile: data.Phone,
		}
		//必须存在一个默认地址
		if isDefault > 0 {
			address.IsDefault = false
		}else {
			address.IsDefault = true
		}
		//存储的是创建小B的用户ID,因为这个地址是小B的
		address.UserId = data.UserId
		address.CId = data.CId
		e.Orm.Create(&address)
	}else {
		e.Orm.Model(&models2.DynamicUserAddress{}).Where("user_id = ? and source = 1 and c_id = ?",data.UserId,data.CId).Updates(map[string]interface{}{
			"name":data.Name,
			"address":data.Address,
			"mobile":data.Phone,
		})
	}

	//修改小B用户名更改
	e.Orm.Model(&sys.SysUser{}).Where("c_id = ? and phone = ?",data.CId,data.Phone).Updates(map[string]interface{}{
		"username":data.UserName,
		"phone":data.Phone,
	})

	return nil
}

// Remove 删除Shop
func (e *Shop) Remove(d *dto.ShopDeleteReq, p *actions.DataPermission) error {

	for _,ids:=range d.Ids {
		var data models.Shop
		e.Orm.Model(&data).Scopes(
				actions.Permission(data.TableName(), p),
			).Where("id = ?",ids).Limit(1).Find(&data)
		if data.Id == 0 {
			continue
		}

		e.Orm.Exec(fmt.Sprintf("DELETE FROM `shop_mark_tag` WHERE `shop_mark_tag`.`shop_id` = %v", ids))
		//只删除一个大B给创建的地址
		e.Orm.Model(models2.DynamicUserAddress{}).Where("source = 1 and user_id = ?",data.UserId).Delete(&models2.DynamicUserAddress{})

		//删除小B的用户
		e.Orm.Model(&sys.SysUser{}).Where("c_id = ? and user_id = ?",data.CId,data.UserId).Delete(&sys.SysUser{})

		e.Orm.Model(&data).Scopes(
			actions.Permission(data.TableName(), p),
		).Where("id = ?",ids).Delete(&data)

	}

	return nil
}
