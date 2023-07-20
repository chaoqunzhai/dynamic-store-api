package redis_db

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go-admin/config"
	"runtime"
	"time"
)

var (
	RedisCli *redis.Client
	Ctx      context.Context
)

func init() {
	Ctx = context.Background()
}
func RedisConn() {
	redisCnf := config.ExtConfig.Redis
	RedisCli = redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", redisCnf.Ip, redisCnf.Port), // Redis服务器地址和端口
		Password:     redisCnf.Password,                                // Redis密码，如果没有密码则为空字符串
		PoolSize:     4 * runtime.NumCPU(),                             //链接池
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolTimeout:  4 * time.Second,
	})
	// Ping命令，检查是否成功连接到Redis
	_, err := RedisCli.Ping(Ctx).Result()
	if err != nil {
		fmt.Println("redis链接失败", err)
		panic(err)
	}
	fmt.Println("成功连接到Redis")
}
