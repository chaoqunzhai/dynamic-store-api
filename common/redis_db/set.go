package redis_db

import (
	"encoding/json"
	"fmt"
	"go-admin/global"
	"go.uber.org/zap"
)


//设置大B客户端登录方式
func SetLoginCnf(siteId int, value interface{}) (val map[string]interface{}, err error) {
	RedisCli.Do(Ctx, "select", global.SmallBLoginCnf)

	data,_:=json.Marshal(value)
	res, err := RedisCli.Set(Ctx, fmt.Sprintf("%v",siteId), string(data), 0).Result()
	if err !=nil{
		zap.S().Errorf("Redis操作,设置大B登录方式失败,原因:%v",err.Error())
	}
	return Marsh(res)
}
