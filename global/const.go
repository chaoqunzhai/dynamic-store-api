package global

const (
	StdOut   = "./logs/info.log"
	StdError = "./logs/error.log"

	//商品目录
	GoodsPath       = "goods"
	SysName         = "动创云订货配送"
	Describe        = "致力于解决订货渠道"
	RoleSuper       = 80 //超管
	RoleCompany     = 81 //大B
	RoleCompanyUser = 82 //大B下用户
	RoleShop        = 83 //小B
	RoleUser        = 84 //用户

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

	OrderLayerKey    = "layer desc"
	OrderTimeKey     = "created_at desc"
	UserNumberAdd    = "add"    //增加
	UserNumberReduce = "reduce" //减少
	UserNumberSet    = "set"    //设置

	CouponGlobal       = 1
	CouponAppointShop  = 2
	CouponAppointClass = 3

	CouponTypeFd   = 1
	CouponDiscount = 2

	//未付款
	OrderPayStatusDefault = 0
	//已付款
	OrderPayStatusSuccess = 1
	//开始线下付款
	OrderPayStatusOfflineDefault = 2
	//下线付款已收款
	OrderPayStatusOfflineSuccess = 3

	//待配送
	OrderStatusWait = 0
	//配送中
	OrderStatusLoading = 1

	//已配送
	OrderStatusOk = 2
	//退回
	OrderStatusReturn = 3
	//退款
	OrderStatusRefund = 4

	//分表的逻辑
	SplitOrder                 = 1
	SplitOrderDefaultTableName = "orders"
	//关联的订单子表,如果进行了订单表的分割,也会默认进行一个分割
	SplitOrderDefaultSubTableName = "order_specs"
	//扩展表
	SplitOrderExtendSubTableName = "order_extend"
	//Cycle 配送的设置
	//每天
	CyCleTimeDay = 1
	//每周
	CyCleTimeWeek = 2
	//支付方式
	PayWechat  = 0 //微信支付
	PayAmount  = 1 //余额支付
	PayCollect = 2 //到付

	//支付状态

	OrderToolsActionStatus   = 1
	OrderToolsActionDelivery = 2

	ScanAdmin        = 1 //管理员操作
	ScanShopRecharge = 2 //用户充值
	ScanShopUse      = 3 //用户消费
	ScanShopRefund   = 4 //用户退款
	ExpressSelf      = 1 //门店自提
	ExpressLocal     = 2 //同城配送
	ExpressLogistics = 3 //物流配送
)

func GetOrderPayStatus(v int) string {
	switch v {
	case OrderPayStatusDefault:
		return "未付款"
	case OrderPayStatusSuccess:
		return "已付款"
	case OrderPayStatusOfflineDefault:
		return "线下付款"
	case OrderPayStatusOfflineSuccess:
		return "下线付款已收款"

	}
	return "未知"
}
func OrderStatus(v int) string {
	switch v {
	case OrderStatusWait:
		return "待配送"
	case OrderStatusLoading:
		return "配送中"
	case OrderStatusOk:
		return "已配送"
	case OrderStatusReturn:
		return "退回"
	case OrderStatusRefund:
		return "退款"

	}
	return "未知"
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
		return "用用户退款户充值"

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

func GetPayTypeStr(v int) string {
	switch v {
	case PayWechat:
		return "微信支付"
	case PayAmount:
		return "余额支付"
	case PayCollect:
		return "到付"

	}
	return "未支付"
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
