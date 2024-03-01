package web_app

import (
	"encoding/json"
	"go-admin/cmd/migrate/migration/models"
	"go-admin/common/redis_db"
	"go-admin/global"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type MakeWeAppInitConf struct {
	BottomNav []interface{}
	Company   models.Company
	CId       int
	Orm *gorm.DB
}

type StyleTheme struct {
	Title     string `json:"title"`
	Name      string `json:"name"`
	MainColor string `json:"main_color"`
	AuxColor  string `json:"aux_color"`
}

// 小程序的配置
// 购物车数量
func (m *MakeWeAppInitConf) cartCount() int {

	return 0
}

// 底栏样式,颜色
func (m *MakeWeAppInitConf) styleTheme() map[string]interface{} {
	dat := map[string]interface{}{
		"name":       "default", //根据这个name来在前段做一个区分设置
	}
	return dat
}

// 底栏菜单
func (m *MakeWeAppInitConf) bottomNav() map[string]interface{} {
	dat := map[string]interface{}{
		"type":            "1",
		"theme":           "diy",
		"backgroundColor": "#FFFFFF",
		"textColor":       "#333333",
		"textHoverColor":  "#F4391c",
		"bulge":           true,
		"imgType":         2,
		"iconColor":       "#333333",
		"iconHoverColor":  "#FF4D4D",
	}
	return dat
}

// 默认图片
func (m *MakeWeAppInitConf) defaultImg() map[string]string {
	return map[string]string{
		"goods":   "",
		"head":    "",
		"store":   "",
		"article": "",
	}
}

// 返回copyright
func (m *MakeWeAppInitConf) copyright() map[string]interface{} {
	return map[string]interface{}{
		"icp":                    "备案号: 陕ICP备2022002980号-1",
		"gov_record":             "",
		"gov_url":                "",
		"market_supervision_url": "",
		"company_name":           "动创云",
		"copyright_link":         "",
		"copyright_desc":         "动创云",
		"auth":                   true,
	}
}

// 返回大B信息
func (m *MakeWeAppInitConf) companyInfo(company models.Company) map[string]interface{} {
	return map[string]interface{}{
		"site_id":         company.Id,
		"site_domain":     "",
		"site_name":       company.Name,
		"logo":            company.Image,
		"seo_title":       "",
		"seo_keywords":    company.Name,
		"seo_description": company.Name,
		"site_tel":        "",
		"logo_square":     "",
		"shop_status":     "1",
	}
}

// 返回server配置
func (m *MakeWeAppInitConf) serviceCnf() map[string]interface{} {
	return map[string]interface{}{
		"h5": map[string]string{
			"type":       "dongchuangyun",
			"wxwork_url": "https://dongchuangyun.com/",
			"third_url":  "https://dongchuangyun.com/",
		},
		"weapp": map[string]string{
			"type":       "dynamic-app",
			"corpid":     "",
			"wxwork_url": "",
		},
		"pc": map[string]string{
			"type":      "third",
			"third_url": "https://dongchuangyun.com",
		},
		"aliapp": map[string]string{
			"type": "none",
		},
	}
}

// 返回插件配置
func (m *MakeWeAppInitConf) addonExist() map[string]interface{} {
	return map[string]interface{}{
		"fenxiao":         1,
		"pintuan":         1,
		"membersignin":    1,
		"memberrecharge":  1,
		"memberwithdraw":  1,
		"pointexchange":   1,
		"manjian":         1,
		"memberconsume":   1,
		"memberregister":  1,
		"coupon":          1,
		"bundling":        1,
		"discount":        1,
		"seckill":         1,
		"topic":           1,
		"store":           0,
		"groupbuy":        1,
		"bargain":         1,
		"presale":         1,
		"notes":           1,
		"membercancel":    1,
		"servicer":        1,
		"live":            1,
		"cards":           1,
		"egg":             1,
		"turntable":       1,
		"memberrecommend": 1,
		"supermember":     1,
		"giftcard":        1,
		"divideticket":    1,
		"birthdaygift":    1,
		"scenefestival":   1,
		"pinfan":          1,
		"hongbao":         1,
		"blindbox":        1,
		"virtualcard":     1,
		"cardservice":     1,
		"cashier":         1,
		"form":            1,
	}
}

func (m *MakeWeAppInitConf) SetLoadRedis() map[string]interface{} {

	diyBottomNav := m.bottomNav()

	diyBottomNav["list"] = m.BottomNav

	dat := map[string]interface{}{
		"cart_count":     m.cartCount(),
		"style_theme":    m.styleTheme(),
		"diy_bottom_nav": diyBottomNav,
		"addon_is_exist": m.addonExist(),
		"default_img":    m.defaultImg(),
		"copyright":      m.copyright(),
		"site_info":      m.companyInfo(m.Company),
		"servicer":       m.serviceCnf(),
		"store_config": map[string]string{
			"store_business": "shop",
		},
	}
	if _,setRedisErr:=redis_db.SetConfigManyInit(m.CId,  global.SmallBConfigKey,dat);setRedisErr!=nil{
		zap.S().Errorf("大B %v,Redis操作,设置大B小程序Config失败,原因:%v",m.CId,setRedisErr.Error())
	}
	return dat

}
func (m *MakeWeAppInitConf)SearchRun()  {
	var companyObject models.Company
	m.Orm.Model(&models.Company{}).Where("id = ?", m.CId).Limit(1).Find(&companyObject)
	m.Company = companyObject
	//首次进来,那就进行一次数据的初始化
	navList := make([]models.WeAppGlobalNavCnf, 0)
	m.Orm.Model(&models.WeAppGlobalNavCnf{}).Order("layer asc").Find(&navList)

	navListData := make([]interface{}, 0)
	for _, row := range navList {
		var object models.CompanyNavCnf
		m.Orm.Model(&models.CompanyNavCnf{}).Where("c_id = ? and g_id = ?", m.CId, row.Id).Limit(1).Find(&object)

		if object.Id > 0 {
			//配置了并且是关闭的,那就返回吧
			if !object.Enable {
				continue
			}
		}else {
			//没有单独配置,那就走全局的策略
			if !row.Enable{
				continue
			}
		}
		navLibMap := make(map[string]interface{}, 0)
		_ = json.Unmarshal([]byte(NavDefine), &navLibMap)
		navLibMap["iconPath"] = row.IconPath
		navLibMap["selectedIconPath"] = row.SelectedIconPath
		navLibMap["text"] = row.Text
		navLibMap["iconClass"] = row.IconClass
		navLibMap["link"] = map[string]interface{}{
			"name":    row.Name,
			"title":   row.Text,
			"wap_url": row.WapUrl,
			"parent":  "MALL_LINK",
		}
		navListData = append(navListData, navLibMap)

	}

	m.BottomNav = navListData

	return
}
func NewMakeWeAppInitConf(siteId int,orm *gorm.DB) *MakeWeAppInitConf {
	m := &MakeWeAppInitConf{
		BottomNav: make([]interface{}, 0),
		CId: siteId,
		Orm: orm,
	}

	return m
}
//查询DB中的配置,查询到的数据写入到redis中
func SearchAndLoadData(siteId int,orm *gorm.DB) map[string]interface{} {
	makeCnf := NewMakeWeAppInitConf(siteId,orm)
	makeCnf.SearchRun()
	return makeCnf.SetLoadRedis()

}
