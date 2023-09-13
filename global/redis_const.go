package global

const (
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
	//要设置的比预期长点
	PhoneMobileDbTimeOut = 130

	PhoneMobileLogin = "login"
	PhoneMobileFind  = "find"

	SmallBLoginKey  = "login_"
	SmallBConfigKey = "cnf_"
	SmallBConfigExtendKey = "extend_app_"
	SmallBMemberToolsKey = "member_"
	SmallBCategoryKey = "category_"
)
