package dto

import (
	"strings"
	"time"

	"go-admin/app/company/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type UpdateEnableReq struct {
	Id             int      `uri:"id" comment:"主键编码"` // 主键编码
	Enable bool `json:"enable"`
}
type CompanyGetPageReq struct {
	dto.Pagination `search:"-"`
	Layer          string    `form:"layer"  search:"type:exact;column:layer;table:company" comment:"排序"`
	Enable         string    `form:"enable"  search:"type:exact;column:enable;table:company" comment:"开关"`
	Name           string    `form:"name"  search:"type:contains;column:name;table:company" comment:"公司(大B)名称"`
	Phone          string    `form:"phone"  search:"type:contains;column:phone;table:company" comment:"负责人联系手机号"`
	UserName       string    `form:"userName"  search:"type:exact;column:user_name;table:company" comment:"大B负责人名称"`
	Shop           string    `form:"shop"  search:"type:exact;column:shop;table:company" comment:"自定义大B系统名称"`
	RenewalTime    time.Time `form:"renewalTime"  search:"type:exact;column:renewal_time;table:company" comment:"续费时间"`
	ExpirationTime time.Time `form:"expirationTime"  search:"type:gt;column:expiration_time;table:company" comment:"到期时间"`
	CompanyOrder
}

type CompanyOrder struct {
	Id             string `form:"idOrder"  search:"type:order;column:id;table:company"`
	CreateBy       string `form:"createByOrder"  search:"type:order;column:create_by;table:company"`
	UpdateBy       string `form:"updateByOrder"  search:"type:order;column:update_by;table:company"`
	CreatedAt      string `form:"createdAtOrder"  search:"type:order;column:created_at;table:company"`
	UpdatedAt      string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:company"`
	DeletedAt      string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:company"`
	Layer          string `form:"layerOrder"  search:"type:order;column:layer;table:company"`
	Enable         string `form:"enableOrder"  search:"type:order;column:enable;table:company"`
	Desc           string `form:"descOrder"  search:"type:order;column:desc;table:company"`
	Name           string `form:"nameOrder"  search:"type:order;column:name;table:company"`
	Phone          string `form:"phoneOrder"  search:"type:order;column:phone;table:company"`
	UserName       string `form:"userNameOrder"  search:"type:order;column:user_name;table:company"`
	Shop           string `form:"shopOrder"  search:"type:order;column:shop;table:company"`
	Address        string `form:"addressOrder"  search:"type:order;column:address;table:company"`
	Longitude      string `form:"longitudeOrder"  search:"type:order;column:longitude;table:company"`
	Latitude       string `form:"latitudeOrder"  search:"type:order;column:latitude;table:company"`
	Image          string `form:"imageOrder"  search:"type:order;column:image;table:company"`
	RenewalTime    string `form:"renewalTimeOrder"  search:"type:order;column:renewal_time;table:company"`
	ExpirationTime string `form:"expirationTimeOrder"  search:"type:order;column:expiration_time;table:company"`
}

func (m *CompanyGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type CompanyRenewGetPage struct {
	dto.Pagination `search:"-"`
}

func (m *CompanyRenewGetPage) GetNeedSearch() interface{} {
	return *m
}

type CompanyInsertReq struct {
	Id             int      `json:"-" comment:"主键编码"` // 主键编码
	Layer          int      `json:"layer" comment:"排序"`
	Enable         bool     `json:"enable" comment:"开关"`
	Desc           string   `json:"desc" comment:"描述信息"`
	Name           string   `json:"name" comment:"公司(大B)名称" binding:"required"`
	Phone          string   `json:"phone" comment:"负责人联系手机号"`
	UserName       string   `json:"user_name" comment:"大B负责人名称"`
	ShopName       string   `json:"shop_name" comment:"自定义大B系统名称"`
	Address        string   `json:"address" comment:"大B地址位置"`
	Longitude      float64  `json:"longitude" comment:""`
	Latitude       float64  `json:"latitude" comment:""`
	Image          []string `json:"image" comment:"logo图片"`
	RenewalTime    string   `json:"renewal_time" comment:"续费时间"`
	ExpirationTime string   `json:"expiration_time" comment:"到期时间"`
	common.ControlBy
}

func (s *CompanyInsertReq) Generate(model *models.Company) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
	model.Layer = s.Layer
	model.Enable = s.Enable
	model.Desc = s.Desc
	model.Name = s.Name
	model.Phone = s.Phone
	model.UserName = s.UserName
	model.ShopName = s.ShopName
	model.Address = s.Address
	model.Longitude = s.Longitude
	model.Latitude = s.Latitude
	if len(s.Image) > 0 {
		model.Image = strings.Join(s.Image, ",")
	}

	if s.RenewalTime != "" {
		t, _ := time.Parse("2006-01-02 15:04:05", s.RenewalTime)
		model.RenewalTime = t
	}
	if s.ExpirationTime != "" {
		t, _ := time.Parse("2006-01-02 15:04:05", s.ExpirationTime)
		model.ExpirationTime = t
	}
}

func (s *CompanyInsertReq) GetId() interface{} {
	return s.Id
}

type CompanyUpdateReq struct {
	Id             int      `uri:"id" comment:"主键编码"` // 主键编码
	Layer          int      `json:"layer" comment:"排序"`
	Enable         bool     `json:"enable" comment:"开关"`
	Desc           string   `json:"desc" comment:"描述信息"`
	Name           string   `json:"name" comment:"公司(大B)名称" binding:"required"`
	Phone          string   `json:"phone" comment:"负责人联系手机号"`
	UserName       string   `json:"userName" comment:"大B负责人名称"`
	ShopName       string   `json:"shop_name" comment:"自定义大B系统名称"`
	Address        string   `json:"address" comment:"大B地址位置"`
	Longitude      float64  `json:"longitude" comment:""`
	Latitude       float64  `json:"latitude" comment:""`
	Image          []string `json:"image" comment:"logo图片"`
	RenewalTime    string   `json:"renewalTime" comment:"续费时间"`
	ExpirationTime string   `json:"expirationTime" comment:"到期时间"`
	common.ControlBy
}

func (s *CompanyUpdateReq) Generate(model *models.Company) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.UpdateBy = s.UpdateBy // 添加这而，需要记录是被谁更新的
	model.Layer = s.Layer
	model.Enable = s.Enable
	model.Desc = s.Desc
	model.Name = s.Name
	model.Phone = s.Phone
	model.UserName = s.UserName
	model.ShopName = s.ShopName
	model.Address = s.Address
	model.Longitude = s.Longitude
	model.Latitude = s.Latitude
	if len(s.Image) > 0 {
		model.Image = strings.Join(s.Image, ",")
	}

	if s.RenewalTime != "" {
		t, _ := time.Parse("2006-01-02 15:04:05", s.RenewalTime)
		model.RenewalTime = t
	}
	if s.ExpirationTime != "" {
		t, _ := time.Parse("2006-01-02 15:04:05", s.ExpirationTime)
		model.ExpirationTime = t
	}
}

func (s *CompanyUpdateReq) GetId() interface{} {
	return s.Id
}

// CompanyGetReq 功能获取请求参数
type CompanyGetReq struct {
	Id int `uri:"id"`
}

func (s *CompanyGetReq) GetId() interface{} {
	return s.Id
}

type SmsUseGetPage struct {
	BeginTime      string `form:"beginTime" search:"type:gte;column:created_at;table:company_sms_record_log" comment:"创建时间"`
	EndTime        string `form:"endTime" search:"type:lte;column:created_at;table:company_sms_record_log" comment:"创建时间"`
	Phone        string `form:"phone" search:"type:contains;column:phone;table:company_sms_record_log" comment:"发送人"`
	dto.Pagination `search:"-"`
}
func (m *SmsUseGetPage) GetNeedSearch() interface{} {
	return *m
}

type CompanySmsUpdate struct {
	Enable bool  `form:"source"`
	Record bool `json:"record" form:"record"`
	AcceptPhone string `json:"accept_phone" form:"record"`
}
type CompanyPayReq struct {
	Source string `form:"source"`
}
type CompanyRenewReq struct {
	Time  string  `json:"time"`
	Money float64 `json:"money"`
	Ids   []int   `json:"ids"`
	Desc  string  `json:"desc"`
}

// CompanyDeleteReq 功能删除请求参数
type CompanyDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *CompanyDeleteReq) GetId() interface{} {
	return s.Ids
}
type RegisterRule struct {
	Type int  `json:"type"`
	Text string `json:"text"`
}

type IndexCount struct {
	ThisDayPayAll string `json:"this_day_pay_all"` //今日销售
	ThisDayNewShop int64 `json:"this_day_new_shop"` //今日新增客户数
	ThisDayPayOkOrder int64 `json:"this_day_pay_ok_order"` //今日订单成交量
	ThisDayPayOkShopUser int64 `json:"this_day_pay_ok_shop_user"` //今日付款客户数
	
	Goods int64 `json:"goods"`  //商品总数
	Shop int64 `json:"shop"` //小B总数
	Line int64 `json:"line"` //路线数量
	Salesman int64 `json:"salesman"` //业务员数据量

	WaitOrder int64 `json:"wait_order"` //待发货订单
	RefundWaitOrder int64 `json:"refund_wait_order"` //待处理的售后单
	WaitSelfOrder int64 `json:"wait_self_order"` //待自提
	GoodsSellOut int64 `json:"goods_sell_out"` //已售罄的商品
	SurplusEms int `json:"surplus_ems"` //剩余短信
}

type NoticeRow struct {
	Name string `json:"name"`
	Link string `json:"link"`
	Subtitle string `json:"subtitle"`
	Time string `json:"time"`
}
type DateCount struct {
	Date string `json:"date"`
	Count int64 `json:"count"`
	AllMoney  float64 `json:"all_money"`
}
type ResponseOrderData struct {
	Date []string `json:"date"`
	OrderTotalPrice []float64 `json:"order_total_price"`
	OrderTotal []int64 `json:"order_total"`
}
type UpdateInfo struct {
	ActionImage bool `form:"action_image"` //是否操作了头像
	Enterprise string `form:"enterprise" gorm:"size:20;comment:企业名称"`
	Filings string `form:"filings" gorm:"size:20;comment:备案号"`
	HelpPhone string   `form:"help_phone" gorm:"size:11;comment:联系手机号"`
	HelpMessage string   `form:"help_message" gorm:"size:50;comment:联系我们"`
	Address        string        `form:"address" gorm:"type:varchar(155);comment:大B地址位置"`
	ShopName       string        `form:"shop_name" gorm:"type:varchar(50);comment:店铺名称"`
	ShopStatus int `form:"shop_status"`
}

type LineDetailReq struct {
	dto.Pagination `search:"-"`
	Cycle int `json:"cycle" form:"cycle"`
}
type ShopLineDetailReq struct {
	dto.Pagination `search:"-"`
	Cycle int `json:"cycle" form:"cycle"`
	ShopId int	`json:"shop_id" form:"shop_id"`
}
type LineGoodsDetail struct {
	dto.Pagination `search:"-"`
 	CycleId int `form:"cycle_id"`
 	LineId int `form:"line_id"`
 	GoodsId int `form:"goods_id"`
}
type DetailCount struct {
	Count int `json:"count"`
	Money float64 `json:"money"`
}
type DetailGoodsRow struct {
	Id int `json:"id"`
	Name string `json:"name"`
	GoodsName string `json:"goods_name"`
	CreatedAt string `json:"created_at"`
	Money string `json:"money"`
	Image string `json:"image"`
	Number int `json:"number"`
	Unit string `json:"unit"`
	AllMoney float64 `json:"-"`
	AllMoneyValue string `json:"all_money_value"`
}

type Agreement struct {
	Value string `json:"value" form:"value"`
}