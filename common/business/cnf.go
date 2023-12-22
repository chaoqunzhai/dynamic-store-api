package business

import (
	"fmt"
	"go-admin/cmd/migrate/migration/models"
	"go-admin/global"
	"gorm.io/gorm"
)


func GetCompanyCnf(cid int, key string, orm *gorm.DB) map[string]int {
	//默认的配置
	defaultCnf := map[string]int{
		"vip":        global.CompanyVip,
		"role":       global.CompanyMaxRole,
		"goods_image": global.CompanyMaxGoodsImage,
		"goods_class": global.CompanyMaxGoodsClass,
		"goods_tag":   global.CompanyMaxGoodsTag,
		"shop_tag":   global.CompanyUserTag,
		"shop":       global.CompanyMaxShop,
		"goods":       global.CompanyMaxGoods,
		"offline_pay":global.OffLinePay,
		"index_message":global.CompanyIndexMessage,
		"index_ads":global.CompanyIndexAds,
		"export_worker":global.CompanyExportWorker,
	}

	//如果查询线路,线路是单独的表中 做的限制
	//线路,短信 都是单独表做的限制,因为是花钱单独购买
	var lineNumber int
	var smsNumber int
	switch key {
	case "line":
		var lineCnf models.CompanyLineCnf
		orm.Model(&models.CompanyLineCnf{}).Select("number,id").Where("c_id = ?",cid).Limit(1).Find(&lineCnf)
		if lineCnf.Id == 0 {
			lineNumber = global.CompanyLine
		}else {
			lineNumber = lineCnf.Number
		}
		defaultCnf["line"] = lineNumber
	case "sms":
		var smsCnf models.CompanySmsQuotaCnf
		orm.Model(&models.CompanySmsQuotaCnf{}).Select("available,id").Where("c_id = ?",cid).Limit(1).Find(&smsCnf)
		if smsCnf.Id == 0 {
			smsNumber = global.CompanySmsNumber
		}else {
			smsNumber = smsCnf.Available
		}
		defaultCnf["sms"] = smsNumber
	}
	var cnf []models.CompanyQuotaCnf
	var sql string
	if key != "" {
		sql = fmt.Sprintf("c_id = %v and enable = %v and `key` = '%v'", cid, true, key)
	} else {
		sql = fmt.Sprintf("c_id = %v and enable = %v", cid, true)
	}
	orm.Model(&models.CompanyQuotaCnf{}).Where(sql).Find(&cnf)
	//没有进行特殊配置,那就都返回系统初始化配置的值

	if len(cnf) == 0 {
		return defaultCnf
	}
	result := make(map[string]int, 0)
	for _, row := range cnf {
		result[row.Key] = row.Number
	}
	if key != "" {
		//如果DB没有配置一些特殊的配置,那就使用global的配置
		_, ok := result[key]
		if !ok {
			v := 0
			switch key {
			case "role":
				v = global.CompanyMaxRole
			case "goods_class":
				v = global.CompanyMaxGoodsClass
			case "goods_tag":
				v = global.CompanyMaxGoodsTag
			case "shop_tag":
				v = global.CompanyUserTag
			case "goods":
				v = global.CompanyMaxGoods
			case "shop":
				v = global.CompanyMaxShop
			case "offline_pay":
				v = global.OffLinePay
			case "index_message":
				v=  global.CompanyIndexMessage
			case "index_ads":
				v = global.CompanyIndexAds
			case "export_worker":
				v = global.CompanyExportWorker
			case "line":
				v = lineNumber

			}
			result[key] = v
		}
	}
	return result
}
