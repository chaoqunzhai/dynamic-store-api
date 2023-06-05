package business

import (
	"fmt"
	"go-admin/cmd/migrate/migration/models"
	"go-admin/global"
	"gorm.io/gorm"
)

func GetCompanyCnf(cid int, key string, orm *gorm.DB) map[string]string {
	defaultCnf := map[string]string{
		"role":       fmt.Sprintf("%v", global.CompanyMaxRole),
		"good_class": fmt.Sprintf("%v", global.CompanyMaxGoodsClass),
		"good_tag":   fmt.Sprintf("%v", global.CompanyMaxGoodsTag),
		"shop_tag":   fmt.Sprintf("%v", global.CompanyUserTag),
	}
	var cnf []models.CompanyCnf
	var sql string
	if key != "" {
		sql = fmt.Sprintf("c_id = %v and enable = %v and key = %v", cid, true, key)
	} else {
		sql = fmt.Sprintf("c_id = %v and enable = %v", cid, true)
	}
	orm.Model(&models.CompanyCnf{}).Where(sql).Find(&cnf)
	//没有进行特殊配置,那就都返回配置即可
	if len(cnf) == 0 {
		return defaultCnf
	}
	result := make(map[string]string, 0)
	for _, row := range cnf {
		result[row.Key] = row.Value
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
			}
			result[key] = fmt.Sprintf("%v", v)
		}
	}
	return result
}
