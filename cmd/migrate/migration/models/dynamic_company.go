package models

//大B的业务信息
import (
	"go-admin/common/models"
	"time"
)

// todo:大B信息
type Company struct {
	RichGlobal
	LeaderId uint `json:"leader_id"`
	Enterprise string `json:"enterprise" gorm:"size:20;comment:企业名称"`
	Filings string `json:"filings" gorm:"size:20;comment:备案号"`
	Image          string    `gorm:"size:20;comment:logo图片地址"`
	Name           string    `gorm:"index;size:20;comment:公司名称"`
	Phone          string    `gorm:"size:11;comment:负责人联系手机号"`
	UserName       string    `gorm:"size:12;comment:大B负责人名称"`
	ShopName       string    `gorm:"size:16;comment:自定义大B系统名称"`
	ShopStatus int `json:"shop_status" gorm:"default:1;index;comment:营业状态"`
	Address        string    `gorm:"size:80;comment:大B地址位置"`
	Longitude      float64   //经度
	Latitude       float64   //纬度
	RenewalTime    time.Time `json:"renewal_time" gorm:"comment:续费时间"`
	ExpirationTime time.Time `json:"expiration_time" gorm:"comment:到期时间"`
	LoginTime models.XTime     `json:"login_time" gorm:"type:datetime(3);comment:登录时间"`
	CopyrightEnable bool `json:"copyright_enable" gorm:"comment:开启版本展示"`
	HelpPhone string `json:"help_phone" gorm:"type:varchar(11);comment:联系我们的电话"`
	HelpMessage string `json:"help_message"  gorm:"type:varchar(120);comment:联系我们信息"`
	Vip 			int 		`json:"vip" gorm:"type:tinyint(1);comment:vip版本"`
	InventoryModule bool `json:"inventory_module" gorm:"comment:仓库功能"`
	SaleUserModule bool `json:"sale_user_module" gorm:"comment:业务员功能"`
}

func (Company) TableName() string {
	return "company"
}

// todo:大B设置的VIP等级
type GradeVip struct {
	BigBRichGlobal
	Name     string  `gorm:"size:30;comment:等级名称"`
	Weight   int     `gorm:"type:tinyint(1);default:1;comment:等级,从小到大"`
	Discount float32 `gorm:"comment:折扣"`
	Upgrade  int     `gorm:"default:0;comment:升级条件,满多少金额,自动升级Weight+1"`
}

func (GradeVip) TableName() string {
	return "grade_vip"
}

// todo:路线信息,被司机关联
// 每个路线就是一个配送员
// 大B下有很多路线
type Line struct {
	BigBRichGlobal
	Name     string `gorm:"index;size:16;comment:路线名称"`
	RenewalTime    models.XTime     `json:"renewal_time" gorm:"type:datetime(3);comment:续费时间"`
	ExpirationTime models.XTime      `json:"expiration_time" gorm:"type:datetime(3);comment:到期时间"`
	DriverId int    `gorm:"index;comment:关联司机"`
}

func (Line) TableName() string {
	return "line"
}

// todo:司机信息
type Driver struct {
	BigBRichGlobal
	UserId int    `gorm:"index;comment:关联的用户ID"`
	Name   string `gorm:"size:12;comment:司机名称"`
	Phone  string `gorm:"index;size:11;comment:联系手机号"`
	Desc   string `gorm:"size:50;comment:备注信息"`
}

func (Driver) TableName() string {
	return "driver"
}

// todo:大B物流配送方式,负责一些特殊的自定义配置
type CompanyExpress struct {
	BigBRichGlobal
	Type int `gorm:"type:tinyint(1);default:2;comment:物流类型"`
}

func (CompanyExpress) TableName() string {
	return "company_express"
}

// 门店
type CompanyExpressStore struct {
	BigBRichGlobal
	Name    string `gorm:"size:20;comment:门店名称"`
	Phone string `json:"phone" gorm:"size:11;comment:电话"`
	Address string `json:"address" gorm:"size:120;comment:大B门店地址位置"`
	Start   string `gorm:"size:12;comment:营业开始时间"`
	End     string `gorm:"size:12;comment:营业结束时间"`
}

func (CompanyExpressStore) TableName() string {
	return "company_express_store"
}

// todo: 运费配置
type CompanyFreight struct {
	BigBRichGlobal
	Type         int `gorm:"type:tinyint(1);default:2;comment:物流类型"`
	QuotaMoney   int `gorm:"comment:达到多少钱可以免运费"`
	StartMoney   int `gorm:"comment:最低的起送金额"`
	FreightMoney int `gorm:"comment:没有达到QuotaMoney,运费多少钱"`
}

func (CompanyFreight) TableName() string {
	return "company_freight"
}

//任务队列状态
type CompanyTasks struct {
	Model
	UserName string `json:"user_name" gorm:"size:25;comment:操作用户名称"`
	Key string `json:"key" gorm:"size:30;comment:唯一的key"` // 可以用这个key 来进行校验, 同一个key 如果还在处理中就不在接受
	Title string `json:"title" gorm:"size:30;"` //保存一些标题信息
	CreateBy  int            `json:"createBy" gorm:"index;comment:创建者"`
	CreatedAt time.Time      `json:"createdAt" gorm:"comment:创建时间"`
	UpdatedAt time.Time      `json:"updatedAt" gorm:"comment:最后更新时间"`
	CId int `json:"c_id" gorm:"index;comment:大BID"`
	Type int `json:"type" gorm:"type:tinyint(1);comment:任务类型 0:订单数据导出 1:汇总 2:路线报表 3:路线配送表"`
	Path string `json:"path" gorm:"size:60;comment:路径"`
	Msg string `json:"msg" gorm:"size:60;"`
	Status int `json:"status" gorm:"type:tinyint(1);comment:任务状态,0:执行中 1:成功 2:失败"`
}
func (CompanyTasks) TableName() string {
	return "company_tasks"
}

//退货配置

type OrderReturnCnf struct {
	BigBRichGlobal
	Value string `json:"value" gorm:"size:15;comment:配送文案"`
	Cost float64 `json:"cost" gorm:";comment:配送费用"`
}
func (OrderReturnCnf) TableName() string {
	return "company_order_return_cnf"
}

//协议配置 agreement
type Agreement struct {
	BigBRichGlobal
	Value string `json:"value" gorm:"comment:协议内容"`
}
func (Agreement) TableName() string {
	return "company_agreement"
}