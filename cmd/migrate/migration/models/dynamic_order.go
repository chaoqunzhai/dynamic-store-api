package models

import (
	"go-admin/common/models"
	"go-admin/common/utils"
	"gorm.io/gorm"
	"time"
)

type CycleTimeConf struct {
	BigBRichGlobal
	Uid       string `json:"uid" gorm:"type:varchar(4);index;comment:周期名称都是天,防止一天可能多个不同周期的配置,加个标识区分周期"`
	Show      bool   `json:"show" gorm:"default:true;comment:是否客户端展示"`
	Type      int    `json:"type" gorm:"type:tinyint(1);default:1;comment:类型,每天,每周"`
	StartWeek int    `json:"start_week" gorm:"type:tinyint(1);default:0;comment:类型为周,每周开始天"`
	EndWeek   int    `json:"end_week" gorm:"type:tinyint(1);default:0;comment:类型为周,每周结束天"`
	StartTime string `json:"start_time" gorm:"size:5;comment:开始下单时间"`
	EndTime   string `json:"end_time" gorm:"size:5;comment:结束时间"`
	GiveDay   int    `json:"give_day" gorm:"type:tinyint(1);default:0;comment:跨天值为0是当天,大于0就是当天+天数"`
	GiveTime  string `json:"give_time" gorm:"size:14;comment:配送时间,例如：15点至19点"`
	Desc      string `json:"desc" gorm:"size:30;comment:描述信息"` //描述
}

func (CycleTimeConf) TableName() string {
	return "cycle_time_conf"
}

// todo:订单,因为订单是一个记录,所有大部分可变的数据都是静态资源,不做关联查询
type Orders struct {
	Model
	Uid            string       `json:"uid" gorm:"size:8;index;comment:关联的配送周期统一UID"`
	CreateBy       int          `json:"createBy" gorm:"index;comment:创建者"`
	SourceType     int          `json:"source_type" gorm:"type:tinyint(1);default:1;index;comment:订单来源"`
	PayType        int          `json:"pay_type" gorm:"type:tinyint(1);default:1;comment:支付方式"`
	Phone          string       `gorm:"size:11;index;comment:用户联系手机号"`
	CreatedAt      models.XTime `json:"createdAt" gorm:"comment:创建时间"`
	UpdatedAt      models.XTime `json:"updatedAt" gorm:"comment:最后更新时间"`
	OrderId        string       `json:"order_id" gorm:"index;size:20;comment:订单ID"`
	EmsId string `json:"ems_id"  gorm:"index;size:20;comment:快递单号"`
	OrderNoId string `json:"order_no_id" gorm:"size:20;comment:预支付流水号"`
	CId            int          `gorm:"index;comment:大BID"`
	Enable         bool         `gorm:"comment:开关"`
	ShopId         int          `gorm:"index;comment:关联客户"`
	OfflinePayId int `json:"offline_pay_id" gorm:"index;comment:线下付款方式"`
	Line           string       `gorm:"size:16;comment:路线名称"`
	LineId         int          `json:"line_id" gorm:"index;type:bigint;comment:线路ID"`
	ParentStatus    int   `gorm:"type:tinyint(1);default:0;comment:保留上次的订单状态"`
	Status         int          `gorm:"type:tinyint(1);default:0;index;comment:订单的状态"`
	DeliveryCode   string       `json:"delivery_code" gorm:"size:9;index;comment:核销码"`
	WriteOffStatus int          `json:"write_off_status" gorm:"type:tinyint(1);default:0;index;comment:核销状态,0:未核销 1:核销"`
	PayMoney       float64      `gorm:"comment:实际支付价"`
	OrderMoney     float64      `json:"order_money" gorm:"comment:需要支付价"`
	GoodsMoney     float64      `json:"goods_money" gorm:"comment:商品总价格"`
	DeductionMoney float64      `json:"deduction_money" gorm:"comment:抵扣费用"`
	Number         int          `gorm:"comment:下单数量"`
	PayStatus      int          `gorm:"type:tinyint(1);default:0;index;comment:支付状态,0:未付款,1:已付款 2:线下付款，3:下线付款已收款"`
	PayTime        models.XTime `json:"pay_time" gorm:"comment:支付时间"`
	DeliveryTime   models.XTime `json:"delivery_time" gorm:"type:date;comment:下单时计算的配送时间"`
	DeliveryType   int          `json:"delivery_type" gorm:"comment:配送类型"`
	DeliveryID     int          `json:"delivery_id" gorm:"comment:关联的配送ID,根据配送类型来查询相关数据"`
	DeliveryStr    string       `json:"delivery_str" gorm:"size:25;comment:配送文案"`
	DeliveryMoney  float64      `json:"delivery_money" gorm:"comment:配送费"`
	CouponId       int          `json:"coupon_id" gorm:"comment:使用优惠卷的ID"`
	DriverId       int          `gorm:"index;comment:司机ID"`
	AddressId      int          `json:"user_address_id" gorm:"index;comment:用户的收货地址,或者自提的店家地址"`
	CouponMoney    float64      `json:"coupon_money" gorm:"comment:优惠卷金额"`
	Buyer          string       `json:"buyer" gorm:"size:24;comment:留言"`
	Desc           string       `json:"desc" gorm:"size:16;comment:备注"`
	Edit bool `json:"edit" gorm:"comment:是否被修改"`
	EditAction string `json:"edit_action" gorm:"size:16;comment:修改"`
	AfterSales     bool         `json:"after_sales" gorm:";comment:是否申请售后,只有小B主动申请生效"`
	AfterStatus    int          `json:"after_status" gorm:"type:tinyint(1);default:0;;comment:售后状态:-2:撤回 -1:驳回, 0:无售后, 1:售后处理中 2:处理完毕  3: 大B退回"`
	ApproveMsg  string    `json:"approve_msg" gorm:"type:varchar(12);comment:审批信息/驳回信息"`
	ApproveStatus int `json:"approve_status" gorm:"type:tinyint(1);default:0;comment:订单审核状态 11:通过 120:驳回"`
	AcceptMsg string `json:"accept_msg"  gorm:"type:varchar(35);comment:欠账验收信息"`
}

func (Orders) TableName() string {
	return "orders"
}

// todo:订单规格
type OrderSpecs struct {
	Model
	GoodsName string       `json:"goods_name" gorm:"size:30;comment:商品名称"`
	CId       int          `gorm:"index;comment:大BID"`
	SpecId    int          `json:"spec_id" gorm:"index;comment:规格ID"`
	CreatedAt models.XTime `json:"created_at" gorm:"comment:创建时间"`
	OrderId   string       `gorm:"index;size:30;comment:关联订单长ID"`
	GoodsId   int          `json:"goods_id" gorm:"comment:商品ID"`
	SpecsName string       `gorm:"size:30;comment:规格名称"`
	Unit      string       `json:"unit" gorm:"type:varchar(8);comment:单位"`
	Number    int          `gorm:"comment:下单规格数"`
	Status    int          `gorm:"type:tinyint(1);default:1;index;comment:配送状态"`
	Money     float64      `gorm:"comment:规格的单价"`
	Image     string       `gorm:"size:15;comment:商品图片路径"`
	Edit      bool         `json:"edit" gorm:"comment:是否被修改"`
	EditAction string `json:"edit_action" gorm:"size:12;comment:修改/单个作废说明"`
	AfterStatus    int          `json:"after_status" gorm:"type:tinyint(1);default:0;;comment:售后状态:-2:撤回 -1:驳回, 0:无售后, 1:售后处理中 2:处理完毕  3: 大B退回"`
	AllMoney  float64      `json:"all_money" gorm:"comment:计算的规格价格"` //创建时 计算好
}

func (e *OrderSpecs) BeforeCreate(_ *gorm.DB) error {
	if e.Number > 0 && e.Money > 0 {
		e.AllMoney = utils.RoundDecimalFlot64(e.Money * float64(e.Number))
	}
	return nil
}
func (OrderSpecs) TableName() string {
	return "order_specs"
}

// todo:订单扩展信息
type OrderExtend struct {
	Model
	//Buyer string `json:"buyer" gorm:"size:24;comment:留言"`
	CId       int          `gorm:"index;comment:大BID"`
	CreatedAt models.XTime `json:"createdAt" gorm:"comment:创建时间"`
	OrderId   string       `json:"order_id" gorm:"index;size:20;comment:订单ID"`
	//DriverId    int         `gorm:"index;comment:司机ID"`
	//AddressId int `json:"user_address_id" gorm:"index;comment:用户的收货地址,或者自提的店家地址"`
	//CouponMoney float64 `json:"coupon_money" gorm:"comment:优惠卷金额"`
}

func (OrderExtend) TableName() string {
	return "order_extend"
}

// todo:订单周期列表的记录
type OrderCycleCnf struct {
	Model
	CreatedAt models.XTime `json:"createdAt" gorm:"comment:创建时间"`
	CId       int          `gorm:"index;comment:大BID"`
	Uid       string       `gorm:"type:varchar(8);comment:周期名称都是天,防止一天可能多个不同周期的配置,加个标识区分周期"`
	//StartTime 查看下单周期开始时间 EndTime 查看下单结束周期
	//直接算好时间,方便订单查询,因为下单的时候 是可以算出来的
	StartTime time.Time `gorm:"comment:记录可下单周期开始时间"`
	EndTime   time.Time `gorm:"comment:记录可下单周期结束时间"`
	//下单周期的文案也是保持最新的
	CreateStr string `json:"create_str" gorm:"size:30;comment:下单日期的文案内容"`
	//配送周期的统一查询
	DeliveryTime models.XTime `json:"delivery_time" gorm:"type:date;index;comment:计算的配送时间"`
	//展示,也是保持最新的
	DeliveryStr string `json:"delivery_str" gorm:"size:30;comment:配送文案"`

	//应该增加一个周期的价格统计.商品统计 商品分类统计
}

func (OrderCycleCnf) TableName() string {
	return "order_cycle_cnf"
}

// todo:订单和Redis的映射表 用于待支付缓存订单
type OrderToRedisMap struct {
	Model
	CreatedAt time.Time `json:"createdAt" gorm:"comment:创建时间"`
	CId       int       `json:"c_id" gorm:"index;comment:大BID"`
	UserId    int       `json:"user_id" gorm:"index;comment:用户ID"`
	RedisKey  string    `json:"redis_key" gorm:"size:20;comment:redis的订单ID"`
	Status    int       `json:"status" gorm:"index;comment:订单状态"`
	RandomId  string    `json:"random_id" gorm:"size:10;comment:随机计算的ID,防止恶意请求订单,payInfo会产生"`
}

func (OrderToRedisMap) TableName() string {
	return "order_to_redis_map"
}

//退货表

type OrderReturn struct {
	Model
	DeliveryType   int          `json:"delivery_type" gorm:"comment:配送类型"`
	PayType        int          `json:"pay_type" gorm:"type:tinyint(1);default:1;comment:订单支付方式"`
	Uid        string       `json:"uid" gorm:"size:8;index;comment:关联的配送周期统一UID"`
	CreateBy   int          `json:"createBy" gorm:"index;comment:退货人"`
	AuditBy int  `json:"audit_by" gorm:"index;comment:审批人"`
	CreatedAt  models.XTime `json:"createdAt" gorm:"comment:退货日期"`
	UpdatedAt  models.XTime `json:"updatedAt" gorm:"comment:操作时间"`
	OrderId    string       `json:"order_id" gorm:"index;size:20;comment:与退货单相关的原始订单的编号"`
	ReturnId   string       `json:"return_id" gorm:"index;size:20;comment:退货单号"`
	CId        int          `gorm:"index;comment:大BID"`
	ShopId     int          `gorm:"index;comment:关联小B客户"`
	LineId     int          `json:"line_id" gorm:"index;type:bigint;comment:退货线路ID"`
	DriverId   int          `gorm:"index;comment:退货处理司机ID"`
	AddressId  int          `json:"user_address_id" gorm:"index;comment:用户的收货地址"`
	GoodsId    int          `json:"goods_id" gorm:"comment:退货商品ID"`
	SpecId     int          `gorm:"index;comment::退货商品规格ID"`
	GoodsName string       `json:"goods_name" gorm:"size:30;comment:商品名称"`
	SpecsName  string       `json:"specs_name" gorm:"size:20;comment:规格名称"`
	Unit       string       `json:"unit" gorm:"type:varchar(8);comment:单位"`
	Number     int          `json:"number" gorm:"comment:退货商品数量"`
	Price      float64      `json:"price" gorm:"comment:商品单价"`
	Image      string       `json:"image" gorm:"size:15;comment:商品图片"`
	RefundDeliveryMoney    float64   `json:"refund_delivery_money" gorm:"comment:支付运费"` //支付运费
	RefundApplyMoney float64 `json:"refund_apply_money" gorm:"comment:退款金额"`
	RefundMoneyType int `json:"refund_money_type" gorm:"type:tinyint(1);default:0;index;comment:退款路径 默认处理中"`
	RefundTime models.XTime `json:"refund_time" gorm:"comment:处理时间"`
	SDesc      string       `json:"s_desc" gorm:"size:24;comment:退货原因"`
	CDesc      string       `json:"c_desc" gorm:"size:24;comment:大B处理信息"`
	Reason string `json:"reason"   gorm:"size:15;comment:售后发起原因"`
	Status     int          `json:"status" gorm:"type:tinyint(1);default:1;index;comment:退货状态, 1:处理中 2:处理完成"`
	Edit bool `json:"edit" gorm:"default:false;comment:是否被编辑"`
	Source int `json:"source" gorm:"comment:原数量"`
	InNumber int `json:"in_number" gorm:"comment:入库数"`
	LossNumber int `json:"loss_number" gorm:"comment:损耗数"`
}

func (OrderReturn) TableName() string {
	return "order_return"
}

// 订单修改记录表
// 只记录数量,其他元素只需去订单中查找即可

type OrderEdit struct {
	Model
	CId          int          `gorm:"index;comment:大BID"`
	CreateBy     string          `json:"createBy" gorm:"size:20;comment:修改人"`
	CreatedAt    models.XTime `json:"createdAt" gorm:"comment:修改日期"`
	OrderId      string       `json:"order_id" gorm:"index;size:20;comment:订单ID"`
	SpecsName string       `gorm:"size:30;comment:规格名称"`
	SourerNumber int          `json:"sourer_number" gorm:"comment:原数量"`
	SourerMoney  float64      `json:"sourer_money" gorm:"comment:原价格"`
	Action int `json:"action" gorm:"comment:减少/新增 0:减少 1:新增"`
	Number       int          `gorm:"comment:新数量"`
	Money        float64      `gorm:"comment:新价格"`
	Status       int          `gorm:"type:tinyint(1);default:1;index;comment:修改结果"`
	Desc         string       `json:"c_desc" gorm:"size:24;comment:修改内容描述"`
}

func (OrderEdit) TableName() string {
	return "order_edit_record"
}
