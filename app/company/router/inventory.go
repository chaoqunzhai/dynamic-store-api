package router

import (
	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"

	"go-admin/app/company/apis"
	"go-admin/common/actions"
	"go-admin/common/middleware"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerInventoryRouter)
}

// registerGradeVipRouter
func registerInventoryRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := apis.CompanyInventory{}
	r := v1.Group("/inventory").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole()).Use(actions.PermissionCompanyRole())
	{
		//获取信息
		r.GET("/cnf_info", api.Info)
		//开启库存
		r.POST("/cnf_update", api.UpdateCnf)

		//仓库商品列表
		r.GET("", api.GetPage)

		//商品选择
		r.GET("/goods",api.Goods)
		//获取仓库商品流水
		r.GET("/records", api.InRecords)
	}




	//入库管理
	warehousing := v1.Group("/in/warehousing").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole()).Use(actions.PermissionCompanyRole())
	{
		//入库单列表
		warehousing.GET("", api.Warehousing)
		//入库单创建
		warehousing.POST("/create", api.WarehousingCreate)
		//入库单记录
		warehousing.GET("/detail/:orderId", api.WarehousingDetail)
	}

	//出库管理
	outbound := v1.Group("/in/outbound").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole()).Use(actions.PermissionCompanyRole())
	{
		//出库单列表
		outbound.GET("", api.Outbound)
		//出库单创建
		outbound.POST("/create", api.OutboundCreate)
		//出库单记录
		outbound.GET("/detail", api.OutboundDetail)
	}


}
