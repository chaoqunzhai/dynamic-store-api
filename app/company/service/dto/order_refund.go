/**
@Author: chaoqun
* @Date: 2024/1/8 17:46
*/
package dto


type RefundDto struct {
	AuditName string `json:"audit_name"`
	Reason    string    `json:"reason" gorm:"column:reason"`
	ShopName   string   `json:"shop_name" gorm:"column:shop"`
	StatusCn    string    `json:"status_cn" gorm:"column:status_cn"`
	Line    string    `json:"line" gorm:"column:line"`
	ReturnID    string    `json:"return_id" gorm:"column:return_id"`
	RefundApplyMoney    int    `json:"refund_apply_money" gorm:"column:refund_apply_money"`
	LineID    int    `json:"line_id" gorm:"column:line_id"`
	CreatedAt    string    `json:"createdAt" gorm:"column:createdAt"`
	Number    int    `json:"number" gorm:"column:number"`
	Uid    string    `json:"uid" gorm:"column:uid"`
	RefundTypeAction    []RefundTypeAction    `json:"refund_type_action" gorm:"column:refund_type_action"`
	RefundMoney string `json:"refund_money"`
	AuditBy    int    `json:"audit_by" gorm:"column:audit_by"`
	Price    int    `json:"price" gorm:"column:price"`
	CouponMoney    float64    `json:"coupon_money" gorm:"column:coupon_money"`
	PayType    int    `json:"pay_type" gorm:"column:pay_type"`
	RefundDeliveryMoney    string    `json:"refund_delivery_money" gorm:"column:refund_delivery_money"`
	ShopID    int    `json:"ShopId" gorm:"column:ShopId"`
	UpdatedAt    string    `json:"updatedAt" gorm:"column:updatedAt"`
	GoodsName    string    `json:"goods_name" gorm:"column:goods_name"`
	Image    string    `json:"image" gorm:"column:image"`
	Address  RefundAddress `json:"address" gorm:"column:address"`
	GoodsID    int    `json:"goods_id" gorm:"column:goods_id"`
	SpecsName    string    `json:"specs_name" gorm:"column:specs_name"`
	SpecID    int    `json:"SpecId" gorm:"column:SpecId"`
	DriverID    int    `json:"DriverId" gorm:"column:DriverId"`
	CreateBy    int    `json:"createBy" gorm:"column:createBy"`
	Unit    string    `json:"unit" gorm:"column:unit"`
	CDesc    string    `json:"c_desc" gorm:"column:c_desc"`
	Driver    string    `json:"driver" gorm:"column:driver"`
	Money    int    `json:"money" gorm:"column:money"`
	RefundTime     string   `json:"refund_time" gorm:"column:refund_time"`
	RefundTodoMoney    float64    `json:"refund_todo_money" gorm:"column:refund_todo_money"`
	UserAddressID    int    `json:"user_address_id" gorm:"column:user_address_id"`
	SDesc    string    `json:"s_desc" gorm:"column:s_desc"`
	OrderID    string    `json:"order_id" gorm:"column:order_id"`
	RefundMoneyType    string    `json:"refund_money_type" gorm:"column:refund_money_type"`
	Status    int    `json:"status" gorm:"column:status"`
	RefundGoods []RefundOrderRow `json:"refund_goods"`
}
type RefundTypeAction struct {
	Name    string    `json:"name" gorm:"column:name"`
	Value    int    `json:"value" gorm:"column:value"`
}
type RefundAddress struct {
	Address    string    `json:"address" gorm:"column:address"`
	Name    string    `json:"name" gorm:"column:name"`
	Mobile    string    `json:"mobile" gorm:"column:mobile"`

}
type RefundOrderRow struct {
	RefundId int `json:"refund_id"`
	Image string `json:"image"`
	GoodsName string `json:"goods_name"`
	SpecName string `json:"spec_name"`
	Price string `json:"price"`
	Number int `json:"number"` //售后数量
	SourceNumber int `json:"source_number"` //编辑后的原数量
	Edit int `json:"edit"` //为前段编辑 售后数量方便
	Unit string `json:"unit"`
	InNumber int `json:"in_number"` //入库数
	LossNumber int `json:"loss_number"` //损耗数
}