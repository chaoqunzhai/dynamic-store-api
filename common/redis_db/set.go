package redis_db

import (
	"encoding/json"
	"fmt"
	"go-admin/global"
	"go.uber.org/zap"
	"time"
)
// 设置验证码
// 默认设置120S过期时间
func SetPhoneCode(T, phone, Code string) (val string, err error) {
	RedisCli.Do(Ctx, "select", global.PhoneMobileCodeDB)

	phoneKey := fmt.Sprintf("%v_%v", T, phone)

	return RedisCli.Set(Ctx, phoneKey, Code, global.PhoneMobileDbTimeOut*time.Second).Result()

}
// 设置大B客户端登录方式
func SetLoginCnf(siteId int, value interface{}) (val map[string]interface{}, err error) {
	RedisCli.Do(Ctx, "select", global.SmallBLoginCnfDB)

	data, _ := json.Marshal(value)
	res, err := RedisCli.Set(Ctx, fmt.Sprintf("%v%v", global.SmallBLoginKey, siteId), string(data), 0).Result()
	if err != nil {
		zap.S().Errorf("Redis操作,设置大B登录方式失败,原因:%v", err.Error())
	}
	return Marsh(res)
}

// 设置获取菜单配置,插件配置 颜色配置等
func SetConfigManyInit(siteId int,redisKey string, value interface{}) (res string, err error) {
	RedisCli.Do(Ctx, "select", global.SmallBConfigDB)

	data, _ := json.Marshal(value)
	res, err = RedisCli.Set(Ctx, fmt.Sprintf("%v%v", redisKey, siteId), string(data), 0).Result()
	if err != nil {
		zap.S().Errorf("Redis操作,设置大B小程序Config失败,原因:%v", err.Error())
	}
	return
}


//设置个人中心配置
func SetMemberInfo(siteId int, value interface{}) (val map[string]interface{}, err error) {
	RedisCli.Do(Ctx, "select", global.SmallBMemberToolsDB)

	data, _ := json.Marshal(value)
	res, err := RedisCli.Set(Ctx, fmt.Sprintf("%v%v", global.SmallBMemberToolsKey, siteId), string(data), 0).Result()
	if err != nil {
		zap.S().Errorf("Redis操作,设置大B个人中心失败,原因:%v", err.Error())
	}
	return Marsh(res)
}
//设置商品分类数据
func SetGoodsCategoryTree(siteId int, value interface{}) (val map[string]interface{}, err error) {
	//选择DB
	RedisCli.Do(Ctx, "select", global.SmallBCategoryDB)

	//数据反序列化
	data, _ := json.Marshal(value)
	res, err := RedisCli.Set(Ctx, fmt.Sprintf("%v%v", global.SmallBCategoryKey, siteId), string(data), 0).Result()
	if err != nil {
		zap.S().Errorf("Redis操作,设置大B小程序商品分类失败,原因:%v", err.Error())
	}
	return Marsh(res)
}

//设置分表配置

func SetCompanyTableSplitCnf(siteId int, value interface{}) {
	//选择DB
	RedisCli.Do(Ctx, "select", global.CompanySplitTableCnf)

	//数据反序列化
	data, _ := json.Marshal(value)
	res, err := RedisCli.Set(Ctx, fmt.Sprintf("%v%v", global.CompanySplitKey, siteId), string(data), 0).Result()
	if err != nil {
		zap.S().Errorf("Redis操作,设置大B分表记录失败,原因:%v", err.Error())
	}
	if res != "ok"{
		zap.S().Errorf("Redis操作,设置大B分表记录失败,返回:%v", res)
	}
	return
}