package business

import (
	"fmt"
	"go-admin/cmd/migrate/migration/models"
	"go-admin/global"
	"gorm.io/gorm"
)

func GetCompanyCnf(cid int, key string, orm *gorm.DB) map[string]int {
	defaultCnf := map[string]int{
		"line":       global.CompanyLine,
		"vip":        global.CompanyVip,
		"role":       global.CompanyMaxRole,
		"good_image": global.CompanyMaxGoodsImage,
		"good_class": global.CompanyMaxGoodsClass,
		"good_tag":   global.CompanyMaxGoodsTag,
		"shop_tag":   global.CompanyUserTag,
		"shop":       global.CompanyMaxShop,
		"good":       global.CompanyMaxGoods,
	}
	var cnf []models.CompanyQuotaCnf
	var sql string
	if key != "" {
		sql = fmt.Sprintf("c_id = %v and enable = %v and `key` = '%v'", cid, true, key)
	} else {
		sql = fmt.Sprintf("c_id = %v and enable = %v", cid, true)
	}
	orm.Model(&models.CompanyQuotaCnf{}).Where(sql).Find(&cnf)
	//没有进行特殊配置,那就都返回配置即可
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
			case "good_class":
				v = global.CompanyMaxGoodsClass
			case "good_tag":
				v = global.CompanyMaxGoodsTag
			case "shop_tag":
				v = global.CompanyUserTag
			case "good":
				v = global.CompanyMaxGoods
			case "shop":
				v = global.CompanyMaxShop
			case "line":
				v = global.CompanyLine

			}
			result[key] = v
		}
	}
	return result
}
