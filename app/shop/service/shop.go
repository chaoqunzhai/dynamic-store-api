package service

import (
	"errors"
	"fmt"
	"github.com/go-admin-team/go-admin-core/sdk/service"
	sys "go-admin/app/admin/models"
	"go-admin/app/shop/models"
	"go-admin/app/shop/service/dto"
	models2 "go-admin/cmd/migrate/migration/models"
	"go-admin/common"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
	"go-admin/global"
	"gorm.io/gorm"
	"strings"
)

type Shop struct {
	service.Service
}
type ShopArrears struct {
	Money float64 `json:"money"`
}
// GetPage 获取Shop列表
func (e *Shop) GetPage(c *dto.ShopGetPageReq, p *actions.DataPermission, list *[]models.Shop, count *int64) error {
	var err error
	var data models.Shop

	query :=e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).Order(global.OrderLayerKey)
	if c.Filter == ""{
		query = query.Order(global.OrderLayerKey).Preload("Tag", func(tx *gorm.DB) *gorm.DB {
			return tx.Select("id", "name")
		})

	}
	err = query.Find(list).Limit(-1).Offset(-1).Count(count).Error
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


func (e *Shop) Insert(userDto *sys.SysUser, c *dto.ShopInsertReq) error {
	var userId int
	var err error
	var data models.Shop

	//先设置请求的默认值
	c.Generate(&data)

	for _,pay:=range c.SelectPay{
		switch pay {
		case global.PayEnWeChat:
			data.IsWeChat = true
		case global.PayEnBalance:
			data.IsBalanceDeduct = true
		case global.PayEnAli:
			data.IsAli = true
		case global.PayEnCredit:
			data.IsCredit = true
		case global.PayEnCashOn:
			data.IsCashOn = true
		}
	}
	//后面根据选择进行重新赋值
	//到这里一定是检测通过的
	if c.ApproveId > 0 {
		var RegisterUser models.CompanyRegisterUserVerify
		e.Orm.Model(&RegisterUser).Where("c_id = ? and id = ?",userDto.CId, c.ApproveId).Limit(1).Find(&RegisterUser)

		data.Phone = RegisterUser.Phone
		data.UserName = RegisterUser.Value
		//创建的用户ID赋值给大B
		userId = RegisterUser.ShopUserId
	}else {
		//创建小B用户
		shopUserDto:=sys.SysShopUser{
			Username: c.UserName,
			NickName: c.UserName,
			Phone: c.Phone,
			Password: c.Password,
			Enable: true,
			CId: userDto.CId,
			Status:global.SysUserSuccess,
			RoleId:global.RoleShop,
		}
		//设置创建用户为大B的名字
		shopUserDto.CreateBy = userDto.UserId
		e.Orm.Create(&shopUserDto)
		userId = shopUserDto.UserId
	}

	//把创建的用户关联到这个小B上面来 重要，否则会检测不存在
	data.UserId = userId
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
	fullAddressName:=make([]string,0)
	strNumbersIds:=make([]string,0)
	for _,mapId:=range c.FullAddress{
		var chinaObject models2.ChinaData
		e.Orm.Model(&chinaObject).Where("id = ? ",mapId).Limit(1).Find(&chinaObject)
		if chinaObject.Id > 0 {
			fullAddressName = append(fullAddressName,chinaObject.Name)
		}
		strNumbersIds = append(strNumbersIds,fmt.Sprintf("%v",mapId))
	}
	//最后一个
	userAddress.ChinaId = c.FullAddress[len(c.FullAddress)-1]
	userAddress.FullAddress = strings.Join(fullAddressName,"-")
	//存储的是创建小B的用户ID,因为这个地址是小B的
	userAddress.UserId = userId
	userAddress.CId = userDto.CId
	createErr:=e.Orm.Create(&userAddress).Error
	if createErr==nil{
		//更新地址
		e.Orm.Model(&models.Shop{}).Where("id = ?",data.Id).Updates(map[string]interface{}{
			"china_id":strings.Join(strNumbersIds,","), //存ID
			"full_address":userAddress.FullAddress,
		})
		if c.ApproveId > 0 {
			//更新审批状态
			var RegisterUserVerify models.CompanyRegisterUserVerify
			e.Orm.Model(&RegisterUserVerify).Where("id = ?",c.ApproveId).Limit(1).Find(&RegisterUserVerify)
			//给用户发送短信通知
			if RegisterUserVerify.Phone !="" { //有手机号 开始发送门店创建成功的短信
				//获取大B的名称
				var company models2.Company
				e.Orm.Model(&company).Select("shop_name").Where("id = ?",userDto.CId).Limit(1).Find(&company)
				common.SendRegisterSuccessSms("门店创建成功通知",company.ShopName,RegisterUserVerify.Phone,userDto.CId,e.Orm)
			}
			e.Orm.Model(&RegisterUserVerify).Where("id = ?",c.ApproveId).Updates(map[string]interface{}{
				"status":2,
			})
		}
	}

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
	for _,pay:=range c.SelectPay{
		switch pay {
		case global.PayEnWeChat:
			data.IsWeChat = true
		case global.PayEnBalance:
			data.IsBalanceDeduct = true
		case global.PayEnAli:
			data.IsAli = true
		case global.PayEnCredit:
			data.IsCredit = true
		case global.PayEnCashOn:
			data.IsCashOn = true
		}
	}
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

	//修改小B用户名更改

	var shopUserObject sys.SysShopUser
	e.Orm.Model(&shopUserObject).Where("c_id = ? and user_id = ?",data.CId,data.UserId).Limit(1).Find(&shopUserObject)
	//保存的是小B的用户ID
	var shopUserId int
	if shopUserObject.UserId > 0  {
		shopUserId = shopUserObject.UserId
		e.Orm.Model(&sys.SysShopUser{}).Where("c_id = ? and user_id = ?",data.CId,data.UserId).Updates(map[string]interface{}{
			"username":data.UserName,
			"phone":data.Phone,
		})
	}
	var userAddress models2.DynamicUserAddress

	var isDefault int64
	e.Orm.Model(&models2.DynamicUserAddress{}).Where("user_id = ? and is_default = 1 and c_id = ?",shopUserId,data.CId).Count(&isDefault)
	//用户把大B创建的地址删了
	e.Orm.Model(&models2.DynamicUserAddress{}).Where("user_id = ? and source = 1 and c_id = ?",shopUserId,data.CId).Limit(1).Find(&userAddress)

	fullAddressName:=make([]string,0)
	strNumbersIds:=make([]string,0)
	for _,mapId:=range c.FullAddress{
		var chinaObject models2.ChinaData
		e.Orm.Model(&chinaObject).Where("id = ? ",mapId).Limit(1).Find(&chinaObject)
		if chinaObject.Id > 0 {
			fullAddressName = append(fullAddressName,chinaObject.Name)
		}
		strNumbersIds = append(strNumbersIds,fmt.Sprintf("%v",mapId))
	}
	//最后一个
	ChinaId := c.FullAddress[len(c.FullAddress)-1]
	FullAddress := strings.Join(fullAddressName,"-")
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
		address.UserId = shopUserId
		address.CId = data.CId
		address.ChinaId  = ChinaId
		address.FullAddress = FullAddress
		e.Orm.Create(&address)
	}else {
		e.Orm.Model(&models2.DynamicUserAddress{}).Where("user_id = ? and source = 1 and c_id = ?",shopUserId,data.CId).Updates(map[string]interface{}{
			"name":data.Name,
			"address":data.Address,
			"mobile":data.Phone,
			"china_id":ChinaId,
			"full_address":FullAddress,
		})
	}

	e.Orm.Model(&models.Shop{}).Where("id = ?",data.Id).Updates(map[string]interface{}{
		"china_id":strings.Join(strNumbersIds,","), //存ID
		"full_address":FullAddress,
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
		e.Orm.Model(&sys.SysShopUser{}).Where("c_id = ? and user_id = ?",data.CId,data.UserId).Delete(&sys.SysShopUser{})

		e.Orm.Model(&data).Scopes(
			actions.Permission(data.TableName(), p),
		).Where("id = ?",ids).Delete(&data)

	}

	return nil
}
