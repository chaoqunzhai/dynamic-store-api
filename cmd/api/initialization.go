/*
*
@Author: chaoqun
* @Date: 2022/7/27 10:12
*/
package api

import (
	"fmt"
	"github.com/go-admin-team/go-admin-core/sdk"
	"go-admin/cmd/migrate/migration/models"
	"go-admin/global"
	"strconv"
)


func InitializationWeApp()  {
	dbs := sdk.Runtime.GetDb()
	diyBottomNav := []map[string]interface{}{
		{
			"icon_path":"icondiy icon-system-home",
			"selected_icon_path":"icondiy icon-system-home-selected",
			"text":"主页",
			"name":"INDEX",
			"wap_url":"/pages/index/index",
			"icon_class":"icon-system-home",
			"layer":1,
			"enable":true,
		},
		{
			"icon_path":"icondiy icon-system-category",
			"selected_icon_path":"icondiy icon-system-category-selected",
			"text":"商品",
			"name":"SHOP_CATEGORY",
			"wap_url":"/pages/goods/category",
			"icon_class":"icon-system-category",
			"layer":2,
			"enable":true,
		},
		{
			"icon_path":"icondiy icondiy icon-system-broadcast-fill",
			"selected_icon_path":"icondiy icon-system-category-selected",
			"text":"咨询",
			"name":"SHOP_INFO",
			"wap_url":"/pages_tool/article/list",
			"icon_class":"icon-system-broadcast-fill",
			"layer":3,
			"enable":false,
		},
		{
			"icon_path":"icondiy icon-system-cart",
			"selected_icon_path":"icondiy icon-system-cart-selected",
			"text":"购物车",
			"name":"SHOPPING_TROLLEY",
			"wap_url":"/pages/goods/cart",
			"icon_class":"icon-system-cart",
			"layer":4,
			"enable":true,
		},
		{
			"icon_path":"icondiy icon-system-my",
			"selected_icon_path":"icondiy icon-system-my-selected",
			"text":"我的",
			"name":"MEMBER_CENTER",
			"wap_url":"/pages/member/index",
			"icon_class":"icon-system-my",
			"layer":5,
			"enable":true,
		},
	}
	for _, db := range dbs {
		for _, row := range diyBottomNav {
			var count int64
			if db.Model(&models.WeAppGlobalNavCnf{}).Where("name = ?", row["name"]).Count(&count); count > 0 {
				continue
			}
			db.Model(&models.WeAppGlobalNavCnf{}).Create(&row)
		}
	}
	diyTools:=[]map[string]interface{}{
		{
			"name":"个人资料",
			"image_url":"../../static/member/default_person.png",
			"wap_url":"/pages_tool/member/info",
			"default_show":1,
		},
		{
			"name":"收货地址",
			"image_url":"../../static/member/default_address.png",
			"wap_url":"/pages_tool/member/address",
			"default_show":1,
		},
		{
			"name":"优惠卷",
			"image_url":"../../static/member/default_discount.png",
			"wap_url":"/pages_tool/member/coupon",
			"default_show":1,
		},
		{
			"name":"我的收藏",
			"image_url":"../../static/member/default_like.png",
			"wap_url":"/pages_tool/member/collection",
			"default_show":1,
		},
		{
			"name":"核销台",
			"image_url":"../../static/member/default_fenxiao.png",
			"wap_url":"/pages_tool/verification/index",
			"default_show":1, //设置是否默认可以显示的
		},
		{
			"name":"付款单",
			"image_url":"../../static/member/default_fukuandan.png",
			"wap_url":"/pages_tool/verification/index",
			"default_show":0, //设置是否默认可以显示的
		},
	}

	for _, db := range dbs {
		for _, row := range diyTools {
			var count int64
			if db.Model(&models.WeAppQuickTools{}).Where("name = ?", row["name"]).Count(&count); count > 0 {
				continue
			}

			db.Model(&models.WeAppQuickTools{}).Create(&row)
		}
	}
}
func Initialization() {

	fmt.Println("开始录入系统初始化配置")
	dbs := sdk.Runtime.GetDb()
	ComQuotaCnf := []map[string]interface{}{
		{
			"key":   "line",
			"value": global.CompanyLine,
		},
		{
			"key":   "vip",
			"value": global.CompanyVip,
		},
		{
			"key":   "role",
			"value": global.CompanyMaxRole,
		},
		{
			"key":   "goods",
			"value": global.CompanyMaxGoods,
		},
		{
			"key":   "shop",
			"value": global.CompanyMaxShop,
		},
		{
			"key":   "goods_class",
			"value": global.CompanyMaxGoodsClass,
		},
		{
			"key":   "goods_tag",
			"value": global.CompanyMaxGoodsTag,
		},
		{
			"key":   "goods_image",
			"value": global.CompanyMaxGoodsImage,
		},
		{
			"key":   "shop_tag",
			"value": global.CompanyUserTag,
		},
		{
			"key":   "offline_pay",
			"value": global.OffLinePay ,
		},
	}
	for _, db := range dbs {
		for _, row := range ComQuotaCnf {
			var thisRow models.CompanyQuotaCnf
			var count int64
			if db.Model(&models.CompanyQuotaCnf{}).Where("`key` = ?", row["key"]).First(&thisRow).Count(&count); count > 0 {
				continue
			}
			Number, _ := strconv.Atoi(fmt.Sprintf("%v", row["value"]))
			rows := &models.CompanyQuotaCnf{
				Key:    fmt.Sprintf("%v", row["key"]),
				Number: Number,
				BigBRichGlobal: models.BigBRichGlobal{
					RichGlobal: models.RichGlobal{
						Enable: true,
					},
				},
			}
			db.Create(&rows)
		}
	}
	//抽成默认配置
	//需要注意,如果有一个菜单下有children,那这个菜单需要有一个layer来声明出场顺序
	Menus := []map[string]interface{}{
		{
			"Name":       "statistics",
			"Path":       "/statistics",
			"Component":  "@/views/statistics/Index",
			"MetaTitle":  "数据统计",
			"KeepAlive":  false,
			"MetaIcon":   "Icons.statistics",
			"Hidden":     false,
			"ParentName": "",
			"Layer":      95,
		},
		{
			"Name":       "setting",
			"Path":       "/setting",
			"Component":  "@/views/setting/Index",
			"MetaTitle":  "系统设置",
			"KeepAlive":  false,
			"MetaIcon":   "Icons.setting",
			"Hidden":     false,
			"ParentName": "",
			"Layer":      91,
		},
		{
			"Name":       "/setting/manage",
			"Path":       "/setting/manage",
			"Component":  "",
			"MetaTitle":  "管理员",
			"KeepAlive":  false,
			"MetaIcon":   "",
			"Hidden":     false,
			"ParentName": "setting",
			"Layer":      90,
		},
		{
			"Name":       "/setting/manage/user/index",
			"Path":       "/setting/manage/user/index",
			"Component":  "@/views/setting/manage/user/Index",
			"MetaTitle":  "员工列表",
			"KeepAlive":  false,
			"MetaIcon":   "",
			"Hidden":     false,
			"ParentName": "/setting/manage",
		},
		{
			"Name":       "/setting/manage/role/index",
			"Path":       "/setting/manage/role/index",
			"Component":  "@/views/setting/manage/role/Index",
			"MetaTitle":  "角色管理",
			"KeepAlive":  false,
			"MetaIcon":   "",
			"Hidden":     false,
			"ParentName": "/setting/manage",
		},
		{
			"Name":       "store",
			"Path":       "/store",
			"Component":  "",
			"MetaTitle":  "商城管理",
			"KeepAlive":  false,
			"MetaIcon":   "Icons.shop",
			"Hidden":     false,
			"ParentName": "",
			"Layer":      93,
		},
		{
			"Name":       "/store/page",
			"Path":       "/store/page",
			"Component":  "",
			"MetaTitle":  "店铺页面",
			"KeepAlive":  false,
			"MetaIcon":   "",
			"Hidden":     false,
			"ParentName": "store",
		},
		{
			"Name":       "/store/page/index",
			"Path":       "/store/page/index",
			"Component":  "@/views/store/page/Index",
			"MetaTitle":  "页面设计",
			"KeepAlive":  false,
			"MetaIcon":   "",
			"Hidden":     false,
			"ParentName": "/store/page",
		},
		{
			"Name":       "line",
			"Path":       "/line",
			"Component":  "",
			"MetaTitle":  "路线管理",
			"KeepAlive":  false,
			"MetaIcon":   "Icons.line",
			"Hidden":     false,
			"ParentName": "",
			"Layer":      94,
		},
		{
			"Name":       "/line/index",
			"Path":       "/line/index",
			"Component":  "",
			"MetaTitle":  "路线列表",
			"KeepAlive":  false,
			"MetaIcon":   "Icons.line",
			"Hidden":     false,
			"ParentName": "line",
		},
		{
			"Name":       "/line/create",
			"Path":       "/line/create",
			"Component":  "@/views/line/Create",
			"MetaTitle":  "路线创建",
			"KeepAlive":  false,
			"MetaIcon":   "",
			"Hidden":     true,
			"ParentName": "line",
		},
		{
			"Name":       "/line/update",
			"Path":       "/line/update",
			"Component":  "@/views/line/Update",
			"MetaTitle":  "路线更新",
			"KeepAlive":  false,
			"MetaIcon":   "",
			"Hidden":     true,
			"ParentName": "line",
		},

		{
			"Name":       "driver",
			"Path":       "/driver/index",
			"Component":  "@/views/driver/Index",
			"MetaTitle":  "司机列表",
			"KeepAlive":  false,
			"MetaIcon":   "",
			"Hidden":     false,
			"ParentName": "line",
		},

		{
			"Name":       "/driver/create",
			"Path":       "/driver/create",
			"Component":  "@/views/driver/Create",
			"MetaTitle":  "司机创建",
			"KeepAlive":  false,
			"MetaIcon":   "",
			"Hidden":     true,
			"ParentName": "line",
		},
		{
			"Name":       "/driver/update",
			"Path":       "/driver/update",
			"Component":  "@/views/driver/Update",
			"MetaTitle":  "司机更新",
			"KeepAlive":  false,
			"MetaIcon":   "",
			"Hidden":     true,
			"ParentName": "line",
		},
		{
			"Name":       "market",
			"Path":       "/market",
			"Component":  "",
			"MetaTitle":  "营销中心",
			"KeepAlive":  false,
			"MetaIcon":   "Icons.market",
			"Hidden":     false,
			"ParentName": "",
			"Layer":      92,
		},
		{
			"Name":       "/market/coupon",
			"Path":       "/market/coupon",
			"Component":  "",
			"MetaTitle":  "优惠券管理",
			"KeepAlive":  false,
			"MetaIcon":   "",
			"Hidden":     false,
			"ParentName": "market",
		},

		{
			"Name":       "/market/coupon/index",
			"Path":       "/market/coupon/index",
			"Component":  "",
			"MetaTitle":  "优惠券列表",
			"KeepAlive":  false,
			"MetaIcon":   "",
			"Hidden":     false,
			"ParentName": "/market/coupon",
		},

		{
			"Name":       "/market/coupon/create",
			"Path":       "/market/coupon/create",
			"Component":  "@/views/market/coupon/Create",
			"MetaTitle":  "创建优惠券",
			"KeepAlive":  false,
			"MetaIcon":   "",
			"Hidden":     false,
			"ParentName": "/market/coupon",
		},

		{
			"Name":       "/market/coupon/update",
			"Path":       "/market/coupon/update",
			"Component":  "@/views/market/coupon/Update",
			"MetaTitle":  "编辑优惠券",
			"KeepAlive":  false,
			"MetaIcon":   "",
			"Hidden":     false,
			"ParentName": "/market/coupon",
		},
		{
			"Name":       "/market/coupon/receive/index",
			"Path":       "/market/coupon/receive/index",
			"Component":  "@/views/market/coupon/Receive",
			"MetaTitle":  "领券记录",
			"KeepAlive":  false,
			"MetaIcon":   "",
			"Hidden":     false,
			"ParentName": "/market/coupon",
		},
		{
			"Name":       "index",
			"Path":       "/index",
			"Component":  "@/views/index/Index",
			"MetaTitle":  "今日概况",
			"KeepAlive":  false,
			"MetaIcon":   "Icons.home",
			"Hidden":     false,
			"ParentName": "",
			"Layer":      100,
		},
		{
			"Name":       "report",
			"Path":       "/report",
			"Component":  "@/views/report/Index",
			"MetaTitle":  "配送报表",
			"KeepAlive":  false,
			"MetaIcon":   "Icons.give",
			"Hidden":     false,
			"ParentName": "",
			"Layer":      99,
		},
		{
			"Name":       "order",
			"Path":       "/order",
			"Component":  "",
			"MetaTitle":  "订单管理",
			"KeepAlive":  false,
			"MetaIcon":   "Icons.order",
			"Hidden":     false,
			"ParentName": "",
			"Layer":      98,
		},
		{
			"Name":       "/order/index",
			"Path":       "/order/index",
			"Component":  "@/views/order/Index",
			"MetaTitle":  "订单列表",
			"KeepAlive":  false,
			"MetaIcon":   "",
			"Hidden":     false,
			"ParentName": "order",
		},
		{
			"Name":       "/order/help",
			"Path":       "/order/help",
			"Component":  "@/views/order/Help",
			"MetaTitle":  "代客下单",
			"KeepAlive":  false,
			"MetaIcon":   "",
			"Hidden":     false,
			"ParentName": "order",
		},
		{
			"Name":       "/order/interval",
			"Path":       "/order/interval",
			"Component":  "@/views/order/interval/index",
			"MetaTitle":  "周期配置",
			"KeepAlive":  false,
			"MetaIcon":   "",
			"Hidden":     false,
			"ParentName": "order",
		},
		{
			"Name":       "/order/detail",
			"Path":       "/order/detail",
			"Component":  "@/views/order/Detail",
			"MetaTitle":  "订单详情",
			"KeepAlive":  false,
			"MetaIcon":   "",
			"Hidden":     true,
			"ParentName": "order",
		},
		{
			"Name":       "/order/shop/orders",
			"Path":       "/order/shop/orders",
			"Component":  "@/views/order/ShopOrderDetail",
			"MetaTitle":  "更多记录",
			"KeepAlive":  false,
			"MetaIcon":   "",
			"Hidden":     true,
			"ParentName": "order",
		},
		{
			"Name":       "goods",
			"Path":       "/goods",
			"Component":  "",
			"MetaTitle":  "商品管理",
			"KeepAlive":  false,
			"MetaIcon":   "Icons.goods",
			"Hidden":     false,
			"ParentName": "",
			"Layer":      97,
		},
		{
			"Name":       "/goods/index",
			"Path":       "/goods/index",
			"Component":  "@/views/goods/Index",
			"MetaTitle":  "商品列表",
			"KeepAlive":  false,
			"MetaIcon":   "",
			"Hidden":     false,
			"ParentName": "goods",
		},
		{
			"Name":       "/goods/create",
			"Path":       "/goods/create",
			"Component":  "@/views/goods/Create",
			"MetaTitle":  "创建商品",
			"KeepAlive":  false,
			"MetaIcon":   "",
			"Hidden":     true,
			"ParentName": "goods",
		},
		{
			"Name":       "/goods/update",
			"Path":       "/goods/update",
			"Component":  "@/views/goods/Update",
			"MetaTitle":  "编辑商品",
			"KeepAlive":  false,
			"MetaIcon":   "",
			"Hidden":     true,
			"ParentName": "goods",
		},
		{
			"Name":       "/goods/class/index",
			"Path":       "/goods/class/index",
			"Component":  "@/views/goods/class/Index",
			"MetaTitle":  "商品分类",
			"KeepAlive":  false,
			"MetaIcon":   "",
			"Hidden":     false,
			"ParentName": "goods",
		},
		{
			"Name":       "/goods/class/create",
			"Path":       "/goods/class/create",
			"Component":  "@/views/goods/class/Create",
			"MetaTitle":  "分类创建",
			"KeepAlive":  false,
			"MetaIcon":   "",
			"Hidden":     true,
			"ParentName": "goods",
		},
		{
			"Name":       "/goods/class/update",
			"Path":       "/goods/class/update",
			"Component":  "@/views/goods/class/Update",
			"MetaTitle":  "分类更新",
			"KeepAlive":  false,
			"MetaIcon":   "",
			"Hidden":     true,
			"ParentName": "goods",
		},
		{
			"Name":       "/goods/tag/index",
			"Path":       "/goods/tag/index",
			"Component":  "@/views/goods/tag/Index",
			"MetaTitle":  "商品标签",
			"KeepAlive":  false,
			"MetaIcon":   "",
			"Hidden":     false,
			"ParentName": "goods",
		},
		{
			"Name":       "/goods/tag/create",
			"Path":       "/goods/tag/create",
			"Component":  "@/views/goods/tag/Create",
			"MetaTitle":  "标签创建",
			"KeepAlive":  false,
			"MetaIcon":   "",
			"Hidden":     true,
			"ParentName": "goods",
		},
		{
			"Name":       "/goods/tag/update",
			"Path":       "/goods/tag/update",
			"Component":  "@/views/goods/tag/Update",
			"MetaTitle":  "标签更新",
			"KeepAlive":  false,
			"MetaIcon":   "",
			"Hidden":     true,
			"ParentName": "goods",
		},
		{
			"Name":       "user",
			"Path":       "/user",
			"Component":  "",
			"MetaTitle":  "客户管理",
			"KeepAlive":  false,
			"MetaIcon":   "Icons.user",
			"Hidden":     false,
			"ParentName": "",
			"Layer":      96,
		},
		{
			"Name":       "/user/index",
			"Path":       "/user/index",
			"Component":  "",
			"MetaTitle":  "客户管理列表",
			"KeepAlive":  false,
			"MetaIcon":   "Icons.user",
			"Hidden":     false,
			"ParentName": "user",
			"Layer":      80,
		},
		{
			"Name":       "/user/list",
			"Path":       "/user/list",
			"Component":  "@/views/user/Index",
			"MetaTitle":  "客户列表",
			"KeepAlive":  false,
			"MetaIcon":   "",
			"Hidden":     false,
			"ParentName": "/user/index",
		},
		{
			"Name":       "/user/tag/index",
			"Path":       "/user/tag/index",
			"Component":  "@/views/user/tag/Index",
			"MetaTitle":  "客户标签",
			"KeepAlive":  false,
			"MetaIcon":   "",
			"Hidden":     false,
			"ParentName": "/user/index",
		},
		{
			"Name":       "/user/tag/create",
			"Path":       "/user/tag/create",
			"Component":  "@/views/user/tag/Create",
			"MetaTitle":  "标签创建",
			"KeepAlive":  false,
			"MetaIcon":   "",
			"Hidden":     true,
			"ParentName": "/user/index",
		},
		{
			"Name":       "/user/tag/update",
			"Path":       "/user/tag/update",
			"Component":  "@/views/user/tag/Update",
			"MetaTitle":  "标签更新",
			"KeepAlive":  false,
			"MetaIcon":   "",
			"Hidden":     true,
			"ParentName": "/user/index",
		},
		{
			"Name":       "/user/grade/index",
			"Path":       "/user/grade/index",
			"Component":  "@/views/user/grade/Index",
			"MetaTitle":  "客户等级管理",
			"KeepAlive":  false,
			"MetaIcon":   "",
			"Hidden":     false,
			"ParentName": "/user/index",
		},
		{
			"Name":       "/user/grade/create",
			"Path":       "/user/grade/create",
			"Component":  "@/views/user/grade/Create",
			"MetaTitle":  "客户等级创建",
			"KeepAlive":  false,
			"MetaIcon":   "",
			"Hidden":     true,
			"ParentName": "/user/index",
		},
		{
			"Name":       "/user/grade/update",
			"Path":       "/user/grade/update",
			"Component":  "@/views/user/grade/Update",
			"MetaTitle":  "客户等级更新",
			"KeepAlive":  false,
			"MetaIcon":   "",
			"Hidden":     true,
			"ParentName": "/user/index",
		},
		{
			"Name":       "/user/salesman",
			"Path":       "/user/salesman",
			"Component":  "",
			"MetaTitle":  "业务员管理",
			"KeepAlive":  false,
			"MetaIcon":   "",
			"Hidden":     false,
			"ParentName": "user",
			"Layer":      79,
		},
		{
			"Name":       "/user/salesman/index",
			"Path":       "/user/salesman/index",
			"Component":  "",
			"MetaTitle":  "业务员列表",
			"KeepAlive":  false,
			"MetaIcon":   "",
			"Hidden":     false,
			"ParentName": "/user/salesman",
		},
		{
			"Name":       "/user/balance",
			"Path":       "/user/balance",
			"Component":  "",
			"MetaTitle":  "余额记录",
			"KeepAlive":  false,
			"MetaIcon":   "",
			"Hidden":     false,
			"ParentName": "user",
			"Layer":      78,
		},
		{
			"Name":       "/user/recharge/index",
			"Path":       "/user/recharge/index",
			"Component":  "@/views/user/recharge/Index",
			"MetaTitle":  "充值记录",
			"KeepAlive":  false,
			"MetaIcon":   "",
			"Hidden":     false,
			"ParentName": "/user/balance",
		},
		{
			"Name":       "/user/balance/index",
			"Path":       "/user/balance/index",
			"Component":  "@/views/user/balance/Index",
			"MetaTitle":  "余额明细",
			"KeepAlive":  false,
			"MetaIcon":   "",
			"Hidden":     false,
			"ParentName": "/user/balance",
		},
		{
			"Name":       "/user/integral/index",
			"Path":       "/user/integral/index",
			"Component":  "@/views/user/integral/Index",
			"MetaTitle":  "积分明细",
			"KeepAlive":  false,
			"MetaIcon":   "",
			"Hidden":     false,
			"ParentName": "/user/balance",
		},
	}
	for _, db := range dbs {

		for _, row := range Menus {
			var thisRow models.DyNamicMenu
			ParentName := fmt.Sprintf("%v", row["ParentName"])
			var count int64
			if db.Model(&models.DyNamicMenu{}).Where("name = ?", row["Name"]).First(&thisRow).Count(&count); count > 0 {

				if ParentName != "" {
					//已经存在的,跳过
					if thisRow.ParentId > 0 {
						continue
					}
					var parentRow models.DyNamicMenu
					db.Model(&models.DyNamicMenu{}).Where("name = ?", ParentName).First(&parentRow)
					if parentRow.Id == 0 {
						continue
					}
					fmt.Println("更新",row["Name"],"的父节点是",ParentName,"id是",parentRow.Id)
					db.Model(&models.DyNamicMenu{}).Where("name = ?",
						row["Name"]).Update("parent_id", parentRow.Id)
				}
				continue
			}

			KeepAlive, _ := strconv.ParseBool(fmt.Sprintf("%v", row["KeepAlive"]))
			hidden, _ := strconv.ParseBool(fmt.Sprintf("%v", row["Hidden"]))
			LayerInt, _ := strconv.Atoi(fmt.Sprintf("%v", row["Layer"]))
			rows := &models.DyNamicMenu{
				Name:      fmt.Sprintf("%v", row["Name"]),
				Path:      fmt.Sprintf("%v", row["Path"]),
				Component: fmt.Sprintf("%v", row["Component"]),
				MetaTitle: fmt.Sprintf("%v", row["MetaTitle"]),
				KeepAlive: KeepAlive,
				MetaIcon:  fmt.Sprintf("%v", row["MetaIcon"]),
				Hidden:    hidden,
				Role:      "admin,company",
				Enable:    true,
				Layer:     LayerInt,
			}

			if ParentName != "" {
				var parentRow models.DyNamicMenu
				db.Model(&models.DyNamicMenu{}).Where("name = ?", ParentName).First(&parentRow)
				if parentRow.Id == 0 {
					continue
				}
				rows.ParentId = parentRow.ParentId
			}
			db.Create(&rows)
		}
	}

}
