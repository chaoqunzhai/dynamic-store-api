package redis_db

import (
	"encoding/json"
	"fmt"
	"go-admin/global"
)

func Marsh(data string) (m map[string]interface{}, err error) {
	m = make(map[string]interface{}, 0)

	err = json.Unmarshal([]byte(data), &m)

	return
}

// 获取小B注册规则配置
func GetRegisterConf(siteId string) (val map[string]interface{}, err error) {
	RedisCli.Do(Ctx, "select", global.SmallBLoginCnfDB)
	res, err := RedisCli.Get(Ctx, siteId).Result()

	return Marsh(res)
}

// 通过手机号获取验证码
func GetPhoneCode(T string, phone string) (val string, err error) {

	phoneKey := fmt.Sprintf("%v_%v", T, phone)
	RedisCli.Do(Ctx, "select", global.PhoneMobileCodeDB)

	return RedisCli.Get(Ctx, phoneKey).Result()

}

// 获取小B菜单配置
func GetSmallBNavbar(siteId string) (val string, err error) {
	return "", err
}

// 获取小B首页数据配置
func GetSmallBIndex(siteId string) (val string, err error) {
	return "", err
}

// 获取商品分类
func GetSmallBCategoty(siteId string) (val string, err error) {
	return "", err
}

// 获取小B购物车
func GetSmallBCart(siteId string) (val string, err error) {
	return "", err
}

// 获取小B购物车
func GetSmallBTools(siteId string) (val string, err error) {
	return "", err
}
