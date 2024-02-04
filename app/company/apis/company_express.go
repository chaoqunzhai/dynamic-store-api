package apis

import (
	"github.com/gin-gonic/gin"
	"go-admin/app/company/service/dto"
	models2 "go-admin/cmd/migrate/migration/models"
	"go-admin/common/actions"
	customUser "go-admin/common/jwt/user"
	"go-admin/common/utils"
	"go-admin/global"
	"gorm.io/gorm"
	"strconv"
)



func (e Company) StoreList(c *gin.Context) {
	err := e.MakeContext(c).
		MakeOrm().
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
	address := make([]map[string]interface{}, 0)
	localAddress := make([]models2.CompanyExpressStore, 0)
	var localObject models2.CompanyExpressStore
	e.Orm.Model(&localObject).Scopes(actions.PermissionSysUser(localObject.TableName(), userDto)).Find(&localAddress)
	for _, r := range localAddress {
		address = append(address, map[string]interface{}{
			"address": r.Address,
			"name":    r.Name,
			"start":   r.Start,
			"end":     r.End,
			"id":r.Id,
		})
	}
	e.OK(address, "successful")
	return
}

func (e Company) GetDelivery(c *gin.Context) {

	err := e.MakeContext(c).
		MakeOrm().
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

	var objectLists []models2.CompanyExpress
	e.Orm.Model(&models2.CompanyExpress{}).Where("c_id = ? ",userDto.CId).Order(global.OrderLayerKey).Find(&objectLists)

	cache:=make([]int,0)
	cacheMap:=make(map[int]interface{},0)
	for _,row:=range objectLists{
		if row.Enable{
			name :=global.GetExpressCn(row.Type)
			cache = append(cache,row.Type)
			cacheMap[row.Type] = map[string]interface{}{
				"name":name,
				"type":row.Type,
			}
		}
	}
	cache = utils.RemoveRepeatInt(cache)
	result:=make([]interface{},0)
	for _,row:=range cache{
		result = append(result,cacheMap[row])
	}
	e.OK(result,"successful")
	return
}
//保证发货方式存在一种,只有在关的时候 才会校验
func ValidExpressLastNumber(cid int,orm *gorm.DB) bool{
	var count int64
	orm.Model(&models2.CompanyExpress{}).Where("c_id = ? and enable = ?",cid,true).Count(&count)

	if count <=1 {
		return false
	}

	return true

}
func (e Company) ExpressList(c *gin.Context) {

	err := e.MakeContext(c).
		MakeOrm().
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	rowKeyValue:=c.Query("type")
	rowKey,_ :=strconv.Atoi(rowKeyValue)
	userDto, err := customUser.GetUserDto(e.Orm, c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}
	//1:查询大B是否有关闭配置
	//2:只返回开启的配置
	//3:如何都没有配置,那就返回全部


	var object models2.CompanyExpress
	e.Orm.Model(&models2.CompanyExpress{}).Scopes(actions.PermissionSysUser(object.TableName(), userDto)).Where("type = ?",rowKey).Limit(1).Find(&object)

	cnf := map[string]interface{}{
		"type": rowKey,
		"desc": global.GetExpressCn(rowKey),
	}
	enable := false
	if object.Id == 0 {
		enable = false
	} else {
		enable = object.Enable
	}

	cnf["enable"] = enable


	var CompanyFreight models2.CompanyFreight
	e.Orm.Model(&models2.CompanyFreight{}).Scopes(actions.PermissionSysUser(CompanyFreight.TableName(), userDto)).Where("type = ?", rowKey).Limit(1).Find(&CompanyFreight)
	if CompanyFreight.Id > 0 {

		Freight := map[string]interface{}{
			"quota_money":   CompanyFreight.QuotaMoney,
			"start_money":   CompanyFreight.StartMoney,
			"freight_money": CompanyFreight.FreightMoney,
		}
		cnf["freight"] = Freight
	}
	if rowKey == global.ExpressSelf {
		address := make([]map[string]interface{}, 0)
		localAddress := make([]models2.CompanyExpressStore, 0)
		var localObject models2.CompanyExpressStore
		e.Orm.Model(&localObject).Scopes(actions.PermissionSysUser(localObject.TableName(), userDto)).Find(&localAddress)
		for _, r := range localAddress {
			address = append(address, map[string]interface{}{
				"address": r.Address,
				"name":    r.Name,
				"start":   r.Start,
				"end":     r.End,
				"id":r.Id,
			})
		}
		cnf["address"] = address
	}

	result:= map[string]interface{}{
		"max_local":global.CompanyMaxLocal,
		"cnf_data":cnf,
	}

	e.OK(result, "successful")
	return
}

func (e Company) ExpressCnfEms(c *gin.Context) {
	req := dto.CompanyExpressCnfReq{}
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

	if !req.Cnf.Enable && !ValidExpressLastNumber(userDto.CId,e.Orm){
		e.Error(500, nil,"必须保留一种发货方式")
		return
	}
	layer :=1
	var object models2.CompanyExpress
	e.Orm.Model(&models2.CompanyExpress{}).Scopes(actions.PermissionSysUser(
		object.TableName(), userDto)).Where("type = ?",global.ExpressEms).Limit(1).Find(&object)

	if object.Id == 0{
		object = models2.CompanyExpress{
			Type: global.ExpressEms,
		}
		object.Desc = global.GetExpressCn(global.ExpressEms)
		object.CId = userDto.CId
		object.Layer = layer
		object.Enable = req.Cnf.Enable
		e.Orm.Create(&object)
	}else {
		e.Orm.Model(&models2.CompanyExpress{}).Where("id = ?",object.Id).Updates(map[string]interface{}{
			"enable":req.Cnf.Enable,
			"layer":layer,
		})
	}
	//快递配置
	var localObject models2.CompanyFreight
	e.Orm.Model(&models2.CompanyFreight{}).Scopes(actions.PermissionSysUser(localObject.TableName(), userDto)).Where("type = ?",global.ExpressEms).Limit(1).Find(&localObject)


	localReq := req.Cnf

	if localObject.Id > 0 {
		e.Orm.Model(&localObject).Updates(map[string]interface{}{
			"quota_money":   localReq.QuotaMoney,
			"start_money":   localReq.StartMoney,
			"freight_Money": localReq.FreightMoney,
		})
		e.OK("更新成功", "successful")
		return
	}
	localObject = models2.CompanyFreight{
		Type:         global.ExpressEms,
		QuotaMoney:   localReq.QuotaMoney,
		StartMoney:   localReq.StartMoney,
		FreightMoney: localReq.FreightMoney,
	}
	localObject.CId = userDto.CId
	localObject.Desc = global.GetExpressCn(global.ExpressEms)
	localObject.Enable = true

	e.Orm.Save(&localObject)
	e.OK("更新成功", "successful")
	return
}


func (e Company) ExpressCnfLocal(c *gin.Context) {
	req := dto.CompanyExpressCnfReq{}
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

	if !req.Cnf.Enable && !ValidExpressLastNumber(userDto.CId,e.Orm){
		e.Error(500, nil,"必须保留一种发货方式")
		return
	}
	layer:=3
	var object models2.CompanyExpress
	e.Orm.Model(&models2.CompanyExpress{}).Scopes(actions.PermissionSysUser(
		object.TableName(), userDto)).Where("type = ?",global.ExpressSameCity).Limit(1).Find(&object)

	if object.Id == 0{
		object = models2.CompanyExpress{
			Type: global.ExpressSameCity,
		}
		object.Layer =layer
		object.Desc = global.GetExpressCn(global.ExpressSameCity)
		object.CId = userDto.CId
		object.Enable = req.Cnf.Enable
		e.Orm.Create(&object)
	}else {
		e.Orm.Model(&models2.CompanyExpress{}).Where("id = ?",object.Id).Updates(map[string]interface{}{
			"enable":req.Cnf.Enable,
			"layer":layer,
		})
	}

	//快递配置
	var localObject models2.CompanyFreight
	e.Orm.Model(&models2.CompanyFreight{}).Scopes(actions.PermissionSysUser(localObject.TableName(), userDto)).Where("type = ?",global.ExpressSameCity).Limit(1).Find(&localObject)

	localReq := req.Cnf

	if localObject.Id > 0 {
		e.Orm.Model(&localObject).Updates(map[string]interface{}{
			"quota_money":   localReq.QuotaMoney,
			"start_money":   localReq.StartMoney,
			"freight_Money": localReq.FreightMoney,
		})
		e.OK("更新成功", "successful")
		return
	}
	localObject = models2.CompanyFreight{
		Type:         global.ExpressSameCity,
		QuotaMoney:   localReq.QuotaMoney,
		StartMoney:   localReq.StartMoney,
		FreightMoney: localReq.FreightMoney,
	}
	localObject.CId = userDto.CId
	localObject.Desc = global.GetExpressCn(global.ExpressSameCity)
	localObject.Enable = true

	e.Orm.Save(&localObject)
	e.OK("更新成功", "successful")
	return
}

func (e Company) ExpressCnfStore(c *gin.Context) {
	req := dto.CompanyExpressCnfReq{}
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
	if !req.Store.Enable && !ValidExpressLastNumber(userDto.CId,e.Orm){
		e.Error(500, nil,"必须保留一种发货方式")
		return
	}
	layer :=2
	var objectStore models2.CompanyExpress
	e.Orm.Model(&models2.CompanyExpress{}).Scopes(actions.PermissionSysUser(objectStore.TableName(), userDto)).Where("type = ?", global.ExpressSelf).Limit(1).Find(&objectStore)
	if objectStore.Id > 0 {
		objectStore.Enable = req.Store.Enable
		objectStore.Layer = layer
		e.Orm.Save(&objectStore)
	} else {
		store := models2.CompanyExpress{}
		store.Enable = req.Store.Enable
		store.CId = userDto.CId
		store.Layer = layer
		store.Type = global.ExpressSelf
		store.Desc = global.GetExpressCn(global.ExpressSelf)
		e.Orm.Create(&store)
	}
	var objectStoreList []models2.CompanyExpressStore

	e.Orm.Model(&models2.CompanyExpressStore{}).Select("id").Where("c_id = ?",userDto.CId).Find(&objectStoreList)
	sourceList:=make([]int,0)
	for _,k:=range objectStoreList{
		sourceList = append(sourceList,k.Id)
	}
	newList:=make([]int,0)
	for _, row := range req.Store.Address {
		newList = append(newList,row.Id)
	}

	diffList := utils.DifferenceInt(sourceList,newList)

	for _,row:=range diffList{
		e.Orm.Model(&models2.CompanyExpressStore{}).Where("id = ?",row).Delete(&models2.CompanyExpressStore{})
	}
	for _, row := range req.Store.Address {

		if row.Name == "" {
			continue
		}
		updateRow:=models2.CompanyExpressStore{
			Address: row.Address,
			Name:    row.Name,
			Start:   row.Start,
			End:     row.End,
		}
		if row.Id > 0 {
			var objectStore2 models2.CompanyExpressStore
			e.Orm.Model(&objectStore2).Scopes(actions.PermissionSysUser(objectStore2.TableName(), userDto)).Where("id = ?",row.Id).Updates(updateRow)
		}else {
			//自提配置
			//先清空
			updateRow.CId = userDto.CId
			e.Orm.Create(&updateRow)
		}

	}
	e.OK("更新成功", "successful")
	return
}
