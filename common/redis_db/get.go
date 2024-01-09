package redis_db

import (
	"encoding/json"
	"errors"
	"fmt"
	sys "go-admin/app/admin/models"
	"go-admin/global"
	"strconv"
)

func Marsh(data string) (m map[string]interface{}, err error) {
	m = make(map[string]interface{}, 0)

	err = json.Unmarshal([]byte(data), &m)

	fmt.Println("err!",err)
	if err!=nil{
		return nil, errors.New(fmt.Sprintf("Unmarshal 数据异常 %v",err))
	}
	return
}

func ArrayMarsh(data string) (m []map[string]interface{}, err error) {
	m = make([]map[string]interface{}, 0)

	err = json.Unmarshal([]byte(data), &m)

	return
}

// 获取小B注册规则配置
func GetRegisterConf(siteId string) (val map[string]interface{}, err error) {
	RedisCli.Do(Ctx, "select", global.SmallBLoginCnfDB)
	res, err := RedisCli.Get(Ctx, fmt.Sprintf("login_%v", siteId)).Result()

	return Marsh(res)
}

// 通过手机号获取验证码
func GetPhoneCode(T string, phone string) (val string, err error) {

	phoneKey := fmt.Sprintf("%v_%v", T, phone)
	RedisCli.Do(Ctx, "select", global.PhoneMobileCodeDB)

	return RedisCli.Get(Ctx, phoneKey).Result()

}

// 获取菜单配置,插件配置 颜色配置等
func GetSmallBConfigInit(siteId interface{}) (val map[string]interface{}, err error) {
	RedisKey := fmt.Sprintf("%v%v", global.SmallBConfigKey, siteId)
	RedisCli.Do(Ctx, "select", global.SmallBConfigDB)
	res, err := RedisCli.Get(Ctx, RedisKey).Result()
	if res == "" {
		return nil, errors.New("不存在")
	}
	return Marsh(res)
}

func GetSmallBConfigExtendKey(siteId interface{}) (val string, err error) {
	RedisKey := fmt.Sprintf("%v%v", global.SmallBConfigKey, siteId)
	RedisCli.Do(Ctx, "select", global.SmallBConfigExtendKey)
	res, err := RedisCli.Get(Ctx, RedisKey).Result()
	if res == "" {
		return "", errors.New("不存在")
	}
	return res,nil
}
// 获取小B个人中心配置
func GetSmallBMemberTools(siteId interface{}) (val string, err error) {
	//set key
	RedisKey := fmt.Sprintf("%v%v", global.SmallBMemberToolsKey, siteId)
	//select db
	RedisCli.Do(Ctx, "select", global.SmallBMemberToolsDB)
	res, err := RedisCli.Get(Ctx, RedisKey).Result()
	if res == "" {
		return "", errors.New("不存在")
	}
	return res, nil
}

// 获取商品分类
func GetSmallBCategoryTree(siteId interface{}) (val []map[string]interface{}, err error) {
	//set key
	RedisKey := fmt.Sprintf("%v%v", global.SmallBCategoryKey, siteId)
	//select db
	RedisCli.Do(Ctx, "select", global.SmallBCategoryDB)

	res, err := RedisCli.Get(Ctx, RedisKey).Result()
	if res == "" {
		return nil, errors.New("不存在")
	}
	return ArrayMarsh(res)
}

// 获取小B购物车
func GetSmallBCart(siteId string) (val string, err error) {
	return "", err
}

// 获取小B购物车
func GetSmallBTools(siteId string) (val string, err error) {
	return "", err
}

// 获取购物车hash数据
func GetCartList(redisKey string) (val map[string]string, err error) {
	RedisCli.Do(Ctx, "select", global.SmallBCartDB)

	result, _, _ := RedisCli.HScan(Ctx, redisKey, 0, "", 0).Result()
	val = make(map[string]string, 0)
	cacheVal := ""
	for index, row := range result {

		//2位为一个配置 [39_22 22]
		if (index+1)%2 == 0 {
			val[cacheVal] = row
		}
		cacheVal = row
	}
	//字典排序,防止每次都会变

	return val, err
}

//获取购物车指定key的数据，redisKey是站点ID_用户ID
//key是 商品_规格或者商品  这样的字符串
func GetCartKey(redisKey, key string) int {
	RedisCli.Do(Ctx, "select", global.SmallBCartDB)

	result, _, _ := RedisCli.HScan(Ctx, redisKey, 0, key, 0).Result()
	numStr := ""
	for index, row := range result {
		if index == 1 {
			numStr = row
		}
	}

	num, _ := strconv.Atoi(numStr)

	return num
}

func GetOrderUserKey(userDto *sys.SysShopUser) string {

	return fmt.Sprintf("%v_%v",userDto.CId,userDto.UserId)
}

//通过订单号获取订单的详细内容
func GetOrderDetail(redisKey,orderId string) (res string, err error)  {
	RedisCli.Do(Ctx, "select", global.OrderDetailDB)

	orderKey :=fmt.Sprintf("%v:%v",redisKey,orderId)
	res, err = RedisCli.Get(Ctx, orderKey).Result()
	if res == "" {
		return "", errors.New("订单不存在")
	}

	return res,err
}

//模糊查询,统计个数
func GetOrderLikeCount(redisKey string) int {
	RedisCli.Do(Ctx, "select", global.OrderDetailDB)

	result,_,error2:=RedisCli.Scan(Ctx,0,fmt.Sprintf("%v*",redisKey),0).Result()

	if error2 !=nil{
		return 0
	}
	return len(result)

}

//获取指定DB下所有的订单列表
func GetOrderCacheList(redisKey string) []string {
	RedisCli.Do(Ctx, "select", global.OrderDetailDB)

	result,_,error2:=RedisCli.Scan(Ctx,0,fmt.Sprintf("%v*",redisKey),0).Result()

	if error2 !=nil{
		return nil
	}


	dat:=make([]string,0)
	for _,row:=range result{
		res,resErr:=RedisCli.Get(Ctx,row).Result()
		if resErr!=nil{
			continue
		}

		dat = append(dat,res)

	}

	return dat
}
//获取全局的动创云配置  字符串key获取即可
func GetAllGlobalCnf(RedisKey string) string {
	RedisCli.Do(Ctx, "select", global.AllGlobalCnf)

	res, _ := RedisCli.Get(Ctx, RedisKey).Result()
	if res == "" {
		return ""
	}
	return res

}

//获取redis中配置的大B 分表配置
func GetSplitTableCnf(siteId interface{}) (dat string,err error)  {
	RedisKey := fmt.Sprintf("%v%v", global.CompanySplitKey, siteId)
	RedisCli.Do(Ctx, "select", global.CompanySplitTableCnf)

	res, _ := RedisCli.Get(Ctx, RedisKey).Result()
	if res == "" {
		return "", errors.New("暂无数据")
	}
	return res,nil
}