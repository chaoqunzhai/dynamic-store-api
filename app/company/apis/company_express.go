package apis

import (
	"github.com/gin-gonic/gin"
	"go-admin/app/company/service/dto"
	models2 "go-admin/cmd/migrate/migration/models"
	"go-admin/common/actions"
	customUser "go-admin/common/jwt/user"
	"go-admin/global"
)

func (e Company) ExpressList(c *gin.Context) {
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
	//1:查询大B是否有关闭配置
	//2:只返回开启的配置
	//3:如何都没有配置,那就返回全部

	expressList := make([]map[string]interface{}, 0)
	for _, row := range global.CompanyGlobalExpress() {

		var object models2.CompanyExpress
		e.Orm.Model(&models2.CompanyExpress{}).Scopes(actions.PermissionSysUser(object.TableName(), userDto)).Where("type = ?",row).Limit(1).Find(&object)

		cnf := map[string]interface{}{
			"type": row,
			"desc": global.GetExpressCn(row),
		}
		enable := false
		if object.Id == 0 {
			enable = true
		} else {
			enable = object.Enable
		}

		cnf["enable"] = enable

		if enable {
			var CompanyFreight models2.CompanyFreight
			e.Orm.Model(&models2.CompanyFreight{}).Scopes(actions.PermissionSysUser(CompanyFreight.TableName(), userDto)).Where("type = ?", row).Limit(1).Find(&CompanyFreight)
			if CompanyFreight.Id > 0 {

				Freight := map[string]interface{}{
					"quota_money":   CompanyFreight.QuotaMoney,
					"start_money":   CompanyFreight.StartMoney,
					"freight_money": CompanyFreight.FreightMoney,
				}
				cnf["freight"] = Freight
			}

		}
		if row == global.ExpressStore {
			address := make([]map[string]string, 0)
			localAddress := make([]models2.CompanyExpressStore, 0)
			var localObject models2.CompanyExpressStore
			e.Orm.Model(&localObject).Scopes(actions.PermissionSysUser(localObject.TableName(), userDto)).Find(&localAddress)
			for _, r := range localAddress {
				address = append(address, map[string]string{
					"address": r.Address,
					"name":    r.Name,
					"start":   r.Start,
					"end":     r.End,
				})
			}
			cnf["address"] = address
		}

		expressList = append(expressList, cnf)
	}
	e.OK(expressList, "successful")
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

	//快递配置
	var localObject models2.CompanyFreight
	e.Orm.Model(&models2.CompanyFreight{}).Scopes(actions.PermissionSysUser(localObject.TableName(), userDto)).Where("type = ?",global.ExpressLocal).Limit(1).Find(&localObject)

	localReq := req.Local

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
		Type:         global.ExpressLocal,
		QuotaMoney:   localReq.QuotaMoney,
		StartMoney:   localReq.StartMoney,
		FreightMoney: localReq.FreightMoney,
	}
	localObject.CId = userDto.CId
	localObject.Desc = global.GetExpressCn(global.ExpressLocal)
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

	var objectStore models2.CompanyExpress
	e.Orm.Model(&models2.CompanyExpress{}).Scopes(actions.PermissionSysUser(objectStore.TableName(), userDto)).Where("type = ?", global.ExpressStore).Limit(1).Find(&objectStore)
	if objectStore.Id > 0 {
		objectStore.Enable = req.Store.Enable
		e.Orm.Save(&objectStore)
	} else {
		store := models2.CompanyExpress{}
		store.Enable = req.Store.Enable
		store.CId = userDto.CId
		store.Type = global.ExpressStore
		store.Desc = global.GetExpressCn(global.ExpressStore)
		e.Orm.Create(&store)
	}

	//自提配置
	//先清空
	var object models2.CompanyExpressStore
	e.Orm.Model(&models2.CompanyExpressStore{}).Scopes(actions.PermissionSysUser(object.TableName(), userDto)).Unscoped().Delete(&models2.CompanyExpressStore{})
	//后增加
	for _, row := range req.Store.Address {
		rr := models2.CompanyExpressStore{
			Address: row.Address,
			Name:    row.Name,
			Start:   row.Start,
			End:     row.End,
		}
		rr.CId = userDto.CId
		e.Orm.Create(&rr)
	}
	e.OK("更新成功", "successful")
	return
}
