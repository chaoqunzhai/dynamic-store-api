package global

import (
	"fmt"
	"github.com/ser163/WordBot/generate"
)

const (
	StdOut   = "./logs/info.log"
	StdError = "./logs/error.log"
	DebugError = "./logs/debug.log"
	LogIngWeApp = "移动客户端,weapp-api"
	LoginRoleSuper = 0
	LoginRoleCompany = 1
	LoginRoleShop = 2
	LogIngPC = "PC端,store-api"
	LogIngUserType = "用户名登录"
	LogIngPhoneType = "手机号登录"
	ExpressStore     = 1 //门店自提
	ExpressLocal     = 2 //同城配送
	ExpressLogistics = 3 //物流配送

	//商品目录
	GoodsPath       = "goods"
	SysName         = "动创云订货配送"
	Describe        = "致力于解决订货渠道"
	RoleSuper       = 80 //超管
	RoleCompany     = 81 //大B
	RoleCompanyUser = 82 //大B下用户
	RoleShop        = 83 //小B
	RoleUser        = 84 //用户

	RegisterUserVerify = 1 //新用户需要审核,通过后才可以登录
	RegisterUserLogin = 2 //新用户直接注册+登录
	//用户关闭的
	SysUserDisable = 1
	//用户是开启的
	SysUserSuccess = 2

	//大B资源限制
	CompanyVip           = 6   //大B最多可以设置6个VIP
	CompanyLine          = 3   //大B最多只能有3个路线
	CompanyMaxRole       = 10  //大B最多可以设置10个角色
	CompanyMaxGoods      = 100 //大B最多可以创建50个商品
	CompanyMaxShop       = 30  //大B最多可以创建30个客户
	CompanyMaxGoodsClass = 20  //大B最多可以设置分类个数
	CompanyMaxGoodsTag   = 20  //大B最多可以设置标签个数
	CompanyMaxGoodsImage = 4   //大B最多可以设置单个商品做多6张图片
	CompanyUserTag       = 30  //大B最多可以设置客户标签个数
	CompanyEmsNumber = 100 //大B默认的可用短信条数
	OffLinePay = 6 //大B最多可以设置线下支付的个数
	OrderLayerKey    = "layer desc"
	OrderTimeKey     = "created_at desc"
	UserNumberAdd    = "add"    //增加
	UserNumberReduce = "reduce" //减少
	UserNumberSet    = "set"    //设置

	CouponAppointClass = 1
	CouponGlobal       = 2
	CouponAppointShop  = 3
	CouponUserStateDefault = 0 //可以领取
	CouponUserStateHash = 1 //已领取
	CouponUserStateOver = 2 //已抢光
	CouponTypeFd   = 0
	CouponDiscount = 1

	OrderStatusPayFail = -2  //支付失败

	OrderStatusClose = -1 //订单关闭

	OrderStatusWaitPay = 0 //默认状态，就是待支付

	OrderStatusWaitSend = 1//待发货

	OrderWaitConfirm = 2//待收货

	OrderWaitRefunding  = 3 //售后

	OrderStatusRefund   = 4 //退款

	OrderStatusReturn = 5 //退货


	OrderPayStatusOfflineSuccess = 6 	//线下付款已收款

	OrderPayStatusOfflineDefault = 7 	//线下付款默认状态

	OrderStatusPaySuccess  = 8 //线上支付成功

	OrderStatusWaitPayDefault = 9 //下单了,但是没有支付的状态,还是放在redis中的

	OrderStatusOver = 10 //订单收尾,那就是收货了,确认了

	//分表的逻辑
	SplitOrder                 = 1
	SplitOrderDefaultTableName = "orders"
	//关联的订单子表,如果进行了订单表的分割,也会默认进行一个分割
	SplitOrderDefaultSubTableName = "order_specs"
	//扩展表
	SplitOrderExtendSubTableName = "order_extend"

	//周期配送下单索引表
	SplitOrderCycleSubTableName = "order_cycle_cnf"
	//Cycle 配送的设置
	//每天
	CyCleTimeDay = 1
	//每周
	CyCleTimeWeek = 2


	//支付状态

	OrderToolsActionStatus   = 1
	OrderToolsActionDelivery = 2

	ScanAdmin        = 1 //管理员操作
	ScanShopRecharge = 2 //用户充值
	ScanShopUse      = 3 //用户消费
	ScanShopRefund   = 4 //用户退款



	CouponReceiveType = "wait" //待领取
	ReceiveCoupon1    = 1      //下单时领取的
	ReceiveCoupon2    = 2      //客户自己手动领取的
	ReceiveCoupon3    = 3      //活动领取的


	GoodsPreview = 0 //全部用户可以预览
	GoodsAuthVip = 1 //只有VIP可以购买

	DeductionBalance = 1 //余额抵扣
	DeductionCredit = 2 //授信额抵扣


	OrderSourceApplet = 5 //小程序
	OrderSourceH5 = 6 //H5
	OrderSourceValet = 7 //代客下单
	OrderSourceWeChat = 8 //微信公众号
	OrderSourceAli = 9 //支付宝

	PayTypeBalance = 1 //余额支付
	PayTypeCredit = 2 //授信额支付

	PayTypeOnlineWechat = 3 //线上微信支付
	PayTypeOnlineAli = 4 //线上支付宝支付
	PayTypeOffline = 5 //线下支付

)

func GetRoleCname(v int) string  {
	switch v {
	case RoleSuper:
		return "超管"
	case RoleCompany:
		return "大B"
	case RoleCompanyUser:
		return "大B员工"
	case RoleShop:
		return "小B"
	case RoleUser:
		return "用户"

	}
	return "非法角色"
}
func GetOrderPayStatus(v int) string {
	switch v {
	case OrderStatusWaitPay:
		return "未付款"
	case OrderStatusPaySuccess:
		return "已付款"
	case OrderPayStatusOfflineDefault:
		return "线下付款"
	case OrderPayStatusOfflineSuccess:
		return "线下已收款"

	}
	return "未知"
}
func OrderStatus(v int) string {

	switch v {
	case OrderStatusPayFail:
		return "支付失败"
	case OrderStatusClose:
		return "订单关闭"
	case OrderStatusWaitPay:
		return "待支付"
	case OrderStatusWaitSend:
		return "待发货"
	case OrderWaitConfirm:
		return "待收货"
	case OrderWaitRefunding:
		return "售后"
	case OrderStatusRefund:
		return "退款"
	case OrderStatusReturn:
		return "退货"
	case OrderPayStatusOfflineSuccess:
		return "线下付款已收款"
	case OrderPayStatusOfflineDefault:
		return "线下付款"
	case OrderStatusPaySuccess:
		return "支付成功"
	case OrderStatusWaitPayDefault:
		return "待支付"
	case OrderStatusOver:
		return "完成"
	}
	return fmt.Sprintf("%v",v)
}
func WeekIntToMsg(v int) string {
	switch v {
	case 1:
		return "一"
	case 2:
		return "二"
	case 3:
		return "三"
	case 4:
		return "四"
	case 5:
		return "五"
	case 6:
		return "六"
	case 0:
		return "日"

	}
	return "未知"
}
func GetScanStr(v int) string {
	switch v {
	case ScanAdmin:
		return "管理员操作"
	case ScanShopRecharge:
		return "用户充值"
	case ScanShopUse:
		return "用户消费"
	case ScanShopRefund:
		return "用户退款"

	}

	return "未知"

}
func GetCouponTypeEn(v int) string {
	switch v {
	case CouponTypeFd:
		return "reduction"
	case CouponDiscount:
		return "discount"
	}

	return "未知"

}
func GetCouponType(v int) string {
	switch v {
	case CouponTypeFd:
		return "满减卷"
	case CouponDiscount:
		return "折扣卷"
	}

	return "未知"

}
//支付方式
func GetPayType(v int)  string {
	switch  v{
	case PayTypeOnlineWechat:
		return "线上微信支付"
	case PayTypeOnlineAli:
		return "线上支付宝支付"
	case PayTypeOffline:
		return "线下支付"
	case PayTypeBalance:
		return "余额支付"
	case PayTypeCredit:
		return "授信额支付"
	}
	return "线上支付"
}
func GetExpressCn(v int) string {
	switch v {
	case ExpressStore:
		return "门店自提"
	case ExpressLocal:
		return "同城配送"
	case ExpressLogistics:
		return "物流配送"
	}
	return "同城配送"
}
func GetOrderSource(v int) string {
	switch v {
	case OrderSourceApplet:
		return "微信小程序"
	case OrderSourceH5:
		return "H5"
	case OrderSourceValet:
		return "代客下单"
	case OrderSourceWeChat:
		return "微信公众号"
	case OrderSourceAli:
		return "支付宝"
	}
	return "H5"
}
func GetCouponStr(v int) string {

	switch v {
	case CouponGlobal:
		return "全场通用"
	case CouponAppointShop:
		return "指定商品"
	case CouponAppointClass:
		return "指定分类"

	}
	return "未知"
}
func RandomName(phone string) string {
	startValue := ""
	if len(phone) > 4 {
		startValue = phone[len(phone)-4:]
	} else {
		startValue = "小白"
	}
	value, _ := generate.GenRandomWorld(3, "mix")
	return fmt.Sprintf("%v_%v", startValue, value.Word)
}
