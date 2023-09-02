package global

const (
	//手机号验证DB
	PhoneMobileCodeDB = iota //0
	SmallBLoginCnfDB         //1
	//小B小程序颜色插件,底部菜单配置
	SmallBConfigDB // 2

	//小B首页
	SmallBIndexDB //3
	//商品分类
	SmallBCategoryDB //4
	//购物车
	SmallBCartDB // 5
	//个人中心工具
	SmallBMemberToolsDB //6
	//要设置的比预期长点
	PhoneMobileDbTimeOut = 130

	PhoneMobileLogin = "login"
	PhoneMobileFind  = "find"

	SmallBLoginKey  = "login_"
	SmallBConfigKey = "cnf_"

	SmallBMemberToolsKey = "member_"
	SmallBCategoryKey = "category_"
)
