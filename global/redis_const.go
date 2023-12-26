package global

import (
	"context"
	"time"
)

const (
	ExportOrderFilePath = "order_export" //订单导出目录 这个目录 会在云端也存在
	//手机号验证DB
	PhoneMobileCodeDB = iota //0
	SmallBLoginCnfDB         //1
	//小B小程序颜色插件
	//底部菜单配置
	//配置按钮文案和商品库存是否展示的
	SmallBConfigDB // 2

	//小B首页
	SmallBIndexDB //3
	//商品分类
	SmallBCategoryDB //4
	//购物车
	SmallBCartDB // 5
	//个人中心菜单展示  +  详情页面中的底栏展示
	SmallBMemberToolsDB //6

	//待支付的订单详细
	OrderDetailDB //7

	//全局的公司一些信息配置
	AllGlobalCnf //8

	AllQueueChannel //9
	//订单的过期时间,设置为半个小时
	OrderExpirationTime = 30 * time.Minute
	//如果在期间未确认收货，系统自动完成收货，默认7天自动收货
	OrderReceiveDays = 7

	//订单完成后，用户在指定期限内可申请售后，设置0天不允许申请
	OrderRefundDays = 3
	//关闭的订单只保留20分钟即可
	OrderCloseExpirationTime = 20 *time.Minute

	//要设置的比预期长点
	PhoneMobileDbTimeOut = 130

	PhoneMobileLogin = "login"
	PhoneMobileFind  = "find"

	SmallBLoginKey  = "login_"
	SmallBConfigKey = "cnf_"
	SmallBConfigExtendKey = "extend_app_"
	SmallBMemberToolsKey = "member_"
	SmallBCategoryKey = "category_"

	WorkerOrderStartName = "order" //订单选中导出
	WorkerReportStartName = "report" //配送报表

)

type ExportReq struct {
	Queue string `json:"queue"`
	Order []string `json:"order"`
	CId int `json:"c_id"`
	OrmId int `json:"orm_id"`
	ExportUser string `json:"export_user"`
	Type string `json:"type"` //类型 0:配送订单导出 1:自提订单导出 2:总汇总表导出 3:基于路线导出
}
type GetQueueReq struct {
	CId int `json:"c_id"`
	Name string `json:"name"`
}

var (

	QueueGroup []string
	RedisCtx context.Context

)
func init()  {
	RedisCtx = context.Background()
	QueueGroup =[]string{
		WorkerOrderStartName,WorkerReportStartName,
	}
}

