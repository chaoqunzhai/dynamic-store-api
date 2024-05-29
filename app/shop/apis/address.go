/**
@Author: chaoqun
* @Date: 2024/5/29 17:50
*/
package apis

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"go-admin/app/shop/service/dto"
	"go-admin/cmd/migrate/migration/models"
	cDto "go-admin/common/dto"
	customUser "go-admin/common/jwt/user"
	"go-admin/global"
	"sort"
	"strings"
)

type ShopAddress struct {
	api.Api
}


func (e ShopAddress)GetAllChinaId(ChinnedId int,cache []int)  []int  {

	var object models.ChinaData
	e.Orm.Model(&models.ChinaData{}).Select("id,pid").Where("id = ? ", ChinnedId).Limit(1).Find(&object)
	if object.Id == 0 {
		return cache
	}
	if object.Pid == 0 {
		return  cache
	}
	cache = append(cache,object.Pid)

	return e.GetAllChinaId(object.Pid,cache)

}
func (e ShopAddress) GetPage(c *gin.Context) {
	req := dto.ShopAddressGetPageReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	userDto, err := customUser.GetUserDto(e.Orm, c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}


	list :=make([]models.DynamicUserAddress,0)
	var count int64
	query := e.Orm.Model(&models.DynamicUserAddress{}).Where("c_id = ?",userDto.CId).Scopes(
		cDto.MakeCondition(req.GetNeedSearch()),
		cDto.Paginate(req.GetPageSize(), req.GetPageIndex()),
	).Order(global.OrderTimeKey)


	if req.Filter != ""{
		likeVal:=fmt.Sprintf("%%%v%%",req.Filter)

		userlike := fmt.Sprintf("`name` like '%v'",likeVal)

		phonelike := fmt.Sprintf("`mobile` like '%v'",likeVal)
		addresslike := fmt.Sprintf("`address` like '%v'",likeVal)

		like:=fmt.Sprintf("( %v OR %v OR %v )",userlike,phonelike,addresslike)
		query = query.Where(like)

	}
	err = query.Find(&list).Limit(-1).Offset(-1).Count(&count).Error

	dataList:=make([]interface{},0)
	for _,row:=range list{
		dat :=map[string]interface{}{
			"id":row.Id,
			"username":row.Name,
			"phone":row.Mobile,
			"is_default":row.IsDefault,
			"full_address":row.FullAddress,
			"address":row.Address,
		}
		if row.ChinaId > 0 {

			cache :=make([]int,0)
			cache = append(cache,row.ChinaId)
			chinaIds:=e.GetAllChinaId(row.ChinaId,cache)
			sort.Slice(chinaIds, func(i, j int) bool {
				return i > j
			})
			dat["china_id"] = chinaIds
		}

		dataList = append(dataList,dat)
	}


	e.PageOK(dataList,int(count),req.PageIndex,req.PageSize,"")
	return
}

func (e ShopAddress) Set(c *gin.Context) {

	req := dto.ShopAddressSet{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req, binding.JSON, nil).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	userDto, err := customUser.GetUserDto(e.Orm, c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}
	var object models.DynamicUserAddress
	e.Orm.Model(&object).Where("id = ? and c_id = ?",req.Id,userDto.CId).Limit(1).Find(&object)
	if object.Id == 0 {
		e.Error(-1,nil,"地址不存在")
		return
	}
	e.Orm.Model(&models.DynamicUserAddress{}).Where("user_id = ? and c_id = ?",object.UserId,userDto.CId).Updates(map[string]interface{}{
		"is_default":false,
	})

	object.IsDefault = true
	e.Orm.Save(&object)

	e.OK("","更新成功")

	return

}
func (e ShopAddress) Update(c *gin.Context) {

	req := dto.ShopAddressUpdateReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req, binding.JSON, nil).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	userDto, err := customUser.GetUserDto(e.Orm, c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}
	var object models.DynamicUserAddress
	e.Orm.Model(&object).Where("id = ? and c_id = ?",req.Id,userDto.CId).Limit(1).Find(&object)
	if object.Id == 0 {
		e.Error(-1,nil,"地址不存在")
		return
	}
	fullAddressName:=make([]string,0)


	for _,mapId:=range req.FullAddress{
		var chinaObject models.ChinaData
		e.Orm.Model(&chinaObject).Where("id = ? ",mapId).Limit(1).Find(&chinaObject)
		if chinaObject.Id > 0 {
			fullAddressName = append(fullAddressName,chinaObject.Name)
		}
		//最后一个
		ChinaId := req.FullAddress[len(req.FullAddress)-1]
		FullAddress := strings.Join(fullAddressName,"-")
		object.ChinaId = ChinaId
		object.FullAddress = FullAddress
	}
	object.Address = req.Address
	object.Mobile = req.Phone
	object.Name = req.UserName

	e.Orm.Save(&object)

	e.OK("","更新成功")

	return

}
func (e ShopAddress) Insert(c *gin.Context) {
	req := dto.ShopAddressInsertReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req, binding.JSON, nil).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	userDto, err := customUser.GetUserDto(e.Orm, c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}
	userAddress:=models.DynamicUserAddress{
		Name: req.UserName,
		Address: req.Address,
		Mobile: req.Phone,
		IsDefault: false,
	}
	fullAddressName:=make([]string,0)
	strNumbersIds:=make([]string,0)
	if len(req.FullAddress) > 0 {
		for _,mapId:=range req.FullAddress{
			var chinaObject models.ChinaData
			e.Orm.Model(&chinaObject).Where("id = ? ",mapId).Limit(1).Find(&chinaObject)
			if chinaObject.Id > 0 {
				fullAddressName = append(fullAddressName,chinaObject.Name)
			}
			strNumbersIds = append(strNumbersIds,fmt.Sprintf("%v",mapId))
		}
		//最后一个
		userAddress.ChinaId = req.FullAddress[len(req.FullAddress)-1]
		userAddress.FullAddress = strings.Join(fullAddressName,"-")
	}
	userAddress.UserId = req.UserId
	userAddress.CId = userDto.CId

	e.Orm.Create(&userAddress)

	e.OK("","")
	return

}
func (e ShopAddress) Delete(c *gin.Context) {

	req := dto.ShopAddressDeleteReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	userDto, err := customUser.GetUserDto(e.Orm, c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}

	var data models.DynamicUserAddress
	e.Orm.Model(&data).Where("c_id = ?",userDto.CId).Delete(&data, req.GetId())

	e.OK("","删除成功s")
	return
}

