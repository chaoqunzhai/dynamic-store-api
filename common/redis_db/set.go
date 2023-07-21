package redis_db

import (
	"encoding/json"
	"fmt"
	"go-admin/global"
	"go.uber.org/zap"
)

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
func SetConfigInit(siteId int, value interface{}) (val map[string]interface{}, err error) {
	RedisCli.Do(Ctx, "select", global.SmallBConfigDB)

	data, _ := json.Marshal(value)
	res, err := RedisCli.Set(Ctx, fmt.Sprintf("%v%v", global.SmallBConfigKey, siteId), string(data), 0).Result()
	if err != nil {
		zap.S().Errorf("Redis操作,设置大B小程序Config失败,原因:%v", err.Error())
	}
	return Marsh(res)
}
