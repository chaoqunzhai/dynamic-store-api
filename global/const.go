package global

import (
	"fmt"
	"github.com/ser163/WordBot/generate"
)

const (
	SymBol = "¥"
	StdOut   = "./logs/info.log"
	StdError = "./logs/error.log"
	DebugError = "./logs/debug.log"
	LogIngWeApp = "mobile-client"
	LoginRoleSuper = 0
	LoginRoleCompany = 1
	LoginRoleShop = 2
	LogIngPC = "company-pc"
	LogIngUserType = "用户名登录"
	LogIngPhoneType = "手机号登录"
	ExpressSelf     = 1 //门店自提 门店地址
	ExpressSameCity    = 2 //周期配送 客户地址
	ExpressEms = 3 //物流配送 客户地址

	CacheImage = "cache_image/" //统一临时存放的目录
	//商品目录
	GoodsPath       = "goods" //商品
	AdsPath = "ads" //广告
	AvatarPath = "avatar" //头像
	SysName         = "动创云订货配送"
	Describe        = "致力于解决订货渠道"
	RoleSuper       = 80 //超管
	RoleCompany     = 81 //大B
	RoleCompanyUser = 82 //大B下用户
	RoleShop        = 83 //小B
	RoleUser        = 84 //用户
	RoleSaleMan = 85 //业务员

	RegisterUserVerify = 1 //新用户需要审核,通过后才可以登录
	RegisterUserLogin = 2 //新用户直接注册+登录
	//用户关闭的
	SysUserDisable = 1
	//用户是开启的
	SysUserSuccess = 2

	//大B资源限制
	CompanyVip           = 10   //大B最多可以设置
	CompanyLine          = 1   //默认路线就是1条,也就是基础版的条数
	CompanyOrderRange = 18 //默认订单查询时间范围 单位为月
	CompanyMaxRole       = 10  //大B最多可以设置个角色
	CompanyMaxGoods      = 500 //大B最多可以创建个商品
	CompanyMaxShop       = 10000  //大B最多可以创建个小B客户
	CompanyMaxGoodsClass = 50  //大B最多可以设置分类个数
	CompanyMaxGoodsTag   = 30  //大B最多可以设置标签个数
	CompanyMaxGoodsImage = 4   //大B最多可以设置单个商品做多张图片
	CompanyUserTag       = 30  //大B最多可以设置客户标签个数
	CompanySmsNumber = 100 //大B默认的可用短信条数,因为短信涉及到缴费,所以短信的可用条数是在其他表中(CompanyEmsQuotaCnf)颗粒度控制
	CompanyIndexMessage = 3 //首页消息条目
	CompanyIndexAds = 5 //广告数量
	CompanyLineBindShop = 500 //路线最多支持绑定多少个客户
	CompanyExportWorker = 5 //导出任务队列个数
	CompanySalesmanNumber = 20 //大B默认拥有个业务员
	CompanySmsRecordTag = true //大B的短信消费记录，默认是开启
	OffLinePay = 8 //大B最多可以设置线下支付的个数
	CompanyMaxLocal = 10 //大B支持最多自提设置,不支持动态配置 现在是固定的
	CompanyMaxUnit = 100 //商品单位100个
	CompanyMaxBrand = 200 //商品品牌

	OrderLayerKey    = " layer desc,id desc "
	OrderUserLayerKey    = "layer,user_id desc"
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

	//订单状态
	OrderStatusWaitPay = 0 //默认状态，就是待支付的订单

	OrderStatusWaitSend = 1 //备货中 待配送 当开启审核后 如果审核通过了 那需要查这个状态

	OrderWaitConfirm = 2 //配送中 配送中 到了配送周期后自动成为了这个待收货,也就是已发货 自提订单默认就是这个

	OrderWaitRefunding = 3 //售后处理中

	OrderStatusCancel = 4 //大B操作作废/驳回

	OrderStatusReturn = 5 //售后处理完毕

	OrderPayStatusOfflineSuccess = 6 //线下付款已收款

	OrderPayStatusOfflineDefault = 7 //线下付款默认状态

	OrderStatusPaySuccess = 8 //线上支付成功

	OrderStatusWaitPayDefault = 9 //下单了,但是没有支付的状态,还是放在redis中的

	OrderStatusOver = 10 //订单收尾,那就是收货了,确认了

	OrderApproveOk = 11 //审核通过了

	OrderApproveReject = 12 //审批驳回

	//OrderStatusWaitPay = 0 //默认状态，就是待支付
	//
	//OrderStatusWaitSend = 1//支持成功: 待发货
	//
	//OrderDelivery = 2// 配送中 到了配送周期 默认就成了一个配送中
	//
	//OrderWaitReturn = 3 // 退货 /退款中
	//
	//OrderStatusReject   = 4 //已驳回
	//
	//OrderStatusOver = 5 //订单收尾,那就是收货了,确认了
	//
	//OrderPayStatusOfflineSuccess = 6 	//线下付款已收款
	//
	//OrderPayStatusOfflineDefault = 7 	//线下付款默认状态
	//
	//OrderStatusPaySuccess  = 8 //线上支付成功
	//
	//OrderStatusWaitPayDefault = 9 //下单了,但是没有支付的状态,还是放在redis中的


	//分表的逻辑
	SplitOrder                 = 1
	SplitOrderDefaultTableName = "orders"
	//关联的订单子表,如果进行了订单表的分割,也会默认进行一个分割
	SplitOrderDefaultSubTableName = "order_specs"
	//扩展表
	SplitOrderExtendSubTableName = "order_extend"

	//周期配送下单索引表
	SplitOrderCycleSubTableName = "order_cycle_cnf"

	//订单修改表
	SplitOrderEdit = "order_edit_record"

	//订单退换货表
	SplitOrderReturn = "order_return"

	//出入库记录流水表
	InventoryRecordLog = "inventory_record"

	//Cycle 配送的设置
	//每天
	CyCleTimeDay = 1
	//每周
	CyCleTimeWeek = 2

	CyCleTimeDayMaxNumber = 3 //配送周期为每天时,最多支持3个时间配置创建
	CyCleTimeWeekMaxNumber = 6 //配送周期为每周时,最多支持6个时间配置创建
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

	CouponState1 = 1 //未使用
	CouponState2 = 2 //已使用
	CouponState3 = 3 //过期
	CouponState4 = 4 //作废
	GoodsPreview = 0 //全部用户可以预览
	GoodsAuthVip = 1 //只有VIP可以购买



	OrderSourceApplet = 5 //小程序
	OrderSourceH5 = 6 //H5
	OrderSourceValet = 7 //代客下单
	OrderSourceWeChat = 8 //微信公众号
	OrderSourceAli = 9 //支付宝
	OrderSourceApp  = 10 //app
	PayTypeBalance = 1 //余额支付
	PayTypeCredit = 2 //授信余额支付
	PayTypeOnlineWechat = 3 //线上微信支付
	PayTypeOnlineAli = 4 //线上支付宝支付
	PayTypeOffline = 5 //线下支付
	PayTypeCashOn = 6 //货到付款
	ExportDeliveryOrder = 0 //配送订单
	ExportSelfOrder = 1 //自提订单
	ExportReportOrder =2 //配送报表

	MaxCompanyOrderReturnCnf = 6 //支持配置 最多退货原因


	//售后状态
	RefundDefault = 1 //售后处理中
	RefundOk = 2 //售后处理完成
	RefundOkOverReject = -1 //大B驳回
	RefundOkCancel = -2 //用户主动撤销
	RefundCompanyCancelCType = 3 //大B作废操作

	//大B处理售后
	RefundMoneyOriginal = 1 //原路退还
	RefundMoneyOffline = 2 //线下退款
	RefundMoneyBalance = 3 //退款到余额
	RefundMoneyCredit = 4 //退款到授信余额

	//出入库记录
	InventoryIn = 1 //常规入库
	InventoryOut = 2 //常规出库
	InventoryRefundIn = 3 //退货入库

	InventoryEditIn = 5//商品编辑入库
	InventoryEditOut = 6//商品编辑出库

	InventoryCancelIn = 7 //订单作废入库
	InventoryApproveIn = 8 //审批驳回入库

	InventoryHelpOut = 9 //代客下单出库

	PayEnBalance = "balance"
	PayEnWeChat = "we_chat"
	PayEnCredit = "credit"
	PayEnCashOn = "cash_on"
	PayEnOffline = "offline"
	PayEnAli = "ali"
)

func GetInventoryInAll() []int  {

	return  []int{InventoryIn,InventoryRefundIn,InventoryEditIn,InventoryCancelIn,InventoryApproveIn}
}

func GetInventoryOutAll() []int  {

	return  []int{InventoryOut,InventoryEditOut,InventoryHelpOut}
}
func GetInventoryActionCn(v int) (bol,val string) {

	switch v {
	case InventoryIn:
		return "+","入库"

	case InventoryOut:

		return "-","出库"
	case InventoryRefundIn:

		return "+","退货入库"
	case InventoryEditIn:

		return "+","订单编辑入库"
	case InventoryEditOut:

		return "-","订单编辑出库"
	case InventoryCancelIn:
		return "+","订单作废入库"
	case InventoryApproveIn:
		return "+","审批驳回入库"
	case InventoryHelpOut:
		return "-","到店开单出库"
	}
	return "",""

}
func OrderEffEct() []int { //配送报表 有效订单状态

	return []int{OrderStatusWaitSend,OrderWaitConfirm,OrderWaitRefunding,OrderStatusOver}
}
func OrderApproveEffEct() []int { //配送报表 有效订单状态

	return []int{OrderStatusWaitSend,OrderWaitConfirm,OrderWaitRefunding,OrderStatusOver}
}
func RefundMoneyTypeStr(v int) string  {
	switch v {
	case RefundMoneyOriginal:
		return "原路退还"
	case RefundMoneyOffline:
		return  "线下退款"
	case RefundMoneyBalance:

		return "退款至余额"
	case RefundMoneyCredit:
		return "退款至授信余额"
	}
	return ""
}
func GetRefundStatus(v int) string {
	switch v {
	case RefundDefault:
		return "处理中"
	case RefundOk: //也就是审核通过
		return "处理完毕"
	case RefundOkOverReject:
		return "驳回"
	case RefundOkCancel:
		return "撤销"
	}
	return "处理中"
}

func GetActionStr(v string) string  {
	switch v {
	case UserNumberAdd:
		return "增加"
	case UserNumberSet:
		return "重设"
	case UserNumberReduce:
		return "减少"

	}
	return "未知"
}
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
		return "备货中"
	case OrderWaitConfirm:
		return "配送中"
	case OrderWaitRefunding:
		return "售后处理"
	case OrderStatusCancel:
		return "已作废"
	case OrderStatusReturn:
		return "售后完毕"
	case OrderPayStatusOfflineSuccess:
		return "线下已收款"
	case OrderPayStatusOfflineDefault:
		return "线下付款"
	case OrderStatusPaySuccess:
		return "支付成功"
	case OrderStatusWaitPayDefault:
		return "待支付"
	case OrderStatusOver:
		return "已完成"
	case OrderApproveReject:
		return "审批驳回"
	case OrderApproveOk:
		return "审批通过"

	}
	return fmt.Sprintf("%v", v)
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
		return "微信支付"
	case PayTypeOnlineAli:
		return "支付宝支付"
	case PayTypeOffline:
		return "线下支付"
	case PayTypeBalance:
		return "余额支付"
	case PayTypeCredit:
		return "授信余额支付"
	case PayTypeCashOn:
		return "货到付款"
	}
	return "线上支付"
}
func GetExpressCn(v int) string {
	switch v {
	case ExpressSelf:
		return "门店自提"
	case ExpressSameCity:
		return "周期配送"
	case ExpressEms:
		return "快递物流"
	}
	return "周期配送"
}
func GetOrderSource(v int) string {
	switch v {
	case OrderSourceApplet:
		return "微信小程序"
	case OrderSourceH5:
		return "H5"
	case OrderSourceValet:
		return "到店开单"
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
