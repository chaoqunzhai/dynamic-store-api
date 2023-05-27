package global

const (
	SysName     = "动创云订货配送"
	Describe    = "致力于解决订货渠道"
	RoleSuper   = 80 //超管
	RoleCompany = 81 //大B
	RoleShop    = 82 //小B
	RoleUser    = 83 //用户

	Super   = "admin"
	Company = "company"

	//大B资源限制
	CompanyMaxRole = 5 //大B最多可以设置5个角色
	CompanyMaxGoodsClass = 30 //大B最多可以设置分类个数
	CompanyMaxGoodsTag = 30 //大B最多可以设置标签个数

	CompanyUserTag = 30 //大B最多可以设置客户标签个数

	OrderLayerKey = "layer desc"


	UserNumberAdd = "add" //增加
	UserNumberReduce = "reduce" //减少
	UserNumberSet = "set" //设置

)
