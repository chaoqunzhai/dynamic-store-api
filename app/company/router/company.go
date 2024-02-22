package router

import (
	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"

	"go-admin/app/company/apis"
	"go-admin/common/actions"
	"go-admin/common/middleware"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerCompanyRouter)
}

// registerCompanyRouter
// 大B的信息,需要鉴定权限 必须是大B 和超管
func registerCompanyRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := apis.Company{}
	r := v1.Group("/company").Use(authMiddleware.MiddlewareFunc()).
		Use(middleware.AuthCheckRole()).Use(actions.PermissionCompanyRole())
	{
		//密码修改
		r.POST("/renew",api.RenewPass)
		//获取全局开启的支付方式
		r.GET("/pay_cnf",api.PayCnf)

		//大B信息
		r.GET("/info", api.Info)
		r.POST("/information",api.Information)
		//首页的一些数值统计
		r.GET("/count",api.Count)
		//首页的一些pie图
		r.GET("/pie",api.Pie)
		//首页的公告配置
		r.GET("/article",api.Article)
		//一些限制配置
		r.GET("/cnf", api.Cnf)

		//注册登录的配置
		r.POST("/register_cnf",api.RegisterCnf)
		r.GET("/register_cnf",api.RegisterCnfInfo)

		r.GET("/home", api.Demo)

		//大B商城模板配置
		r.POST("/category", api.SaveCategory)
		r.GET("/category", api.Category)
		//门店自提配置
		r.GET("/express/store",api.StoreList)
		//获取开启的配送方式
		r.GET("/delivery",api.GetDelivery)
		//获取指定配送方式的信息
		r.GET("/express", api.ExpressList)
		//大B的配送
		//同城配置
		r.POST("/express/cnf/local", api.ExpressCnfLocal)
		//自提配置
		r.POST("/express/cnf/store", api.ExpressCnfStore)
		//物流配置
		r.POST("/express/cnf/ems", api.ExpressCnfEms)

		//展示一些限制的文案，例如当前可用路线是多少条
		r.GET("/quota/cnf",api.QuotaCnf)

		//协议配置
		r.GET("/agreement",api.AgreementCnf)
		r.POST("/agreement",api.AgreementUpdate)
	}
	//大B用户管理 + 业务员管理 完全可以复用接口
	{
		//返回推广码+登录的地址
		r.GET("/promotionCode",api.PromotionCode)
		//用户列表
		r.GET("/user/list", api.List)
		//更新用户信息
		r.PUT("/user/:id", api.UpdateUser)
		//密码修改
		r.POST("/user/uppass",api.UpPass)
		//对用户进行下线,大B看不到了,但是超管还是可以看到的,更新用户的enable
		r.POST("/user/offline", api.Offline)
		//增加系统用户,防止恶意注册 必须是role =
		r.POST("/user/add", api.CreateUser)
		//增加业务员 role也必须是=85 因为有校验否则注册不过去
		r.POST("/user/sales_man", api.CreateSalesManUser)

		r.GET("/user/mini", api.MiniList)
		r.POST("/user/code", api.MakeCode)
	}

}
