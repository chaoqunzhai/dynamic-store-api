package middleware

import (
	"errors"
	"fmt"
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/alibaba/sentinel-golang/core/flow"
	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/response"
)

func init() {

	SentinelInit()
}

func SentinelInit() {
	// 务必先进行初始化
	err := sentinel.InitDefault()
	if err != nil {
		fmt.Println("限流器初始化失败")
	}
	fmt.Println("初始化测试流控配置成功！！！！！")
	_, _ = flow.LoadRules([]*flow.Rule{
		//下面测试就是: 1秒内最多10个请求,超出的请求 都需要等待800毫秒
		{
			//Threshold 是 10，Sentinel 默认使用1s作为控制周期，表示1秒内10个请求匀速排队，所以排队时间就是 1000ms/10 = 100ms；
			Resource:               "throttling",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Throttling, // 流控效果为匀速排队
			// 请求的间隔控制在  例如 1000(StatIntervalInMs:周期)/10(10个请求Threshold)=100 ms (100毫秒通过一个请求)
			Threshold:              10,		 //StatIntervalInMs周期内多个请求通过
			MaxQueueingTimeMs:      1000,             // 最长排队等待时间,客户请求等待时间 等待1000秒
			StatIntervalInMs:       1000,            //指定多少秒为统计一个周期
		},
		{
			// Threshold + StatIntervalInMs 可组合出多长时间限制通过多少请求，这里相当于限制为 10 qps
			Resource: "limit",
			//Direct表示直接使用字段 Threshold 作为阈值；WarmUp表示使用预热方式计算Token的阈值。
			TokenCalculateStrategy: flow.Direct,
			//Reject表示超过阈值直接拒绝，Throttling表示匀速排队
			ControlBehavior:  flow.Reject, //直接拒绝
			Threshold:        100,          //阈值
			StatIntervalInMs: 1000,        //一秒作为统计周期，1秒内最多100个请求过来
		},
	})
}
func SentinelContext() gin.HandlerFunc {

	return func(c *gin.Context) {
		//使用排队算法
		sentE, b := sentinel.Entry("throttling", sentinel.WithTrafficType(base.Inbound))
		if b != nil {
			// 请求被流控，可以从 BlockError 中获取限流详情
			//fmt.Println("限流生效",b.Error(),b.BlockMsg())
			response.Error(c, 500, errors.New("请求过于频繁,请稍后重试"), "请求过于频繁,请稍后重试")
			c.Abort()
			return
		}
		sentE.Exit()
		// 处理请求
		c.Next()
	}
}
