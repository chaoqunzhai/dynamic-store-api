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
	"strconv"
)

func Initialization() {

	fmt.Println("开始录入系统初始化配置")
	dbs := sdk.Runtime.GetDb()
	//抽成默认配置
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
			"Layer":95,
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
			"Layer":90,
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
			"Layer":93,
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
			"Layer":94,
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
			"Layer":92,
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
			"Layer":100,
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
			"Layer":99,
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
			"Layer":98,
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
			"Name":       "goods",
			"Path":       "/goods",
			"Component":  "",
			"MetaTitle":  "商品管理",
			"KeepAlive":  false,
			"MetaIcon":   "Icons.goods",
			"Hidden":     false,
			"ParentName": "",
			"Layer":97,
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
			"Layer":96,
		},
		{
			"Name":       "/user/index",
			"Path":       "/user/index",
			"Component":  "@/views/user/Index",
			"MetaTitle":  "客户列表",
			"KeepAlive":  false,
			"MetaIcon":   "",
			"Hidden":     false,
			"ParentName": "user",
		},
		{
			"Name":       "/user/tag/index",
			"Path":       "/user/tag/index",
			"Component":  "@/views/user/tag/Index",
			"MetaTitle":  "客户标签",
			"KeepAlive":  false,
			"MetaIcon":   "",
			"Hidden":     false,
			"ParentName": "user",
		},
		{
			"Name":       "/user/tag/create",
			"Path":       "/user/tag/create",
			"Component":  "@/views/user/tag/Create",
			"MetaTitle":  "标签创建",
			"KeepAlive":  false,
			"MetaIcon":   "",
			"Hidden":     true,
			"ParentName": "user",
		},
		{
			"Name":       "/user/tag/update",
			"Path":       "/user/tag/update",
			"Component":  "@/views/user/tag/Update",
			"MetaTitle":  "标签更新",
			"KeepAlive":  false,
			"MetaIcon":   "",
			"Hidden":     true,
			"ParentName": "user",
		},
		{
			"Name":       "/user/grade/index",
			"Path":       "/user/grade/index",
			"Component":  "@/views/user/grade/Index",
			"MetaTitle":  "客户等级管理",
			"KeepAlive":  false,
			"MetaIcon":   "",
			"Hidden":     false,
			"ParentName": "user",
		},
		{
			"Name":       "/user/grade/create",
			"Path":       "/user/grade/create",
			"Component":  "@/views/user/grade/Create",
			"MetaTitle":  "客户等级创建",
			"KeepAlive":  false,
			"MetaIcon":   "",
			"Hidden":     true,
			"ParentName": "user",
		},
		{
			"Name":       "/user/grade/update",
			"Path":       "/user/grade/update",
			"Component":  "@/views/user/grade/Update",
			"MetaTitle":  "客户等级更新",
			"KeepAlive":  false,
			"MetaIcon":   "",
			"Hidden":     true,
			"ParentName": "user",
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
		{
			"Name":       "manage",
			"Path":       "manage",
			"Component":  "",
			"MetaTitle":  "管理员",
			"KeepAlive":  false,
			"MetaIcon":   "",
			"Hidden":     false,
			"ParentName": "",
			"Layer":91,
		},
		{
			"Name":       "/manage/user/index",
			"Path":       "/manage/user/index",
			"Component":  "@/views/manage/user/Index",
			"MetaTitle":  "管理员列表",
			"KeepAlive":  false,
			"MetaIcon":   "",
			"Hidden":     false,
			"ParentName": "manage",
		},
		{
			"Name":       "/manage/role/index",
			"Path":       "/manage/role/index",
			"Component":  "@/views/manage/role/Index",
			"MetaTitle":  "角色管理",
			"KeepAlive":  false,
			"MetaIcon":   "",
			"Hidden":     false,
			"ParentName": "manage",
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
					db.Model(&models.DyNamicMenu{}).Where("name = ?",
						row["Name"]).Update("parent_id", parentRow.Id)
				}
				continue
			}

			KeepAlive, _ := strconv.ParseBool(fmt.Sprintf("%v", row["KeepAlive"]))
			hidden, _ := strconv.ParseBool(fmt.Sprintf("%v", row["Hidden"]))
			LayerInt,_:=strconv.Atoi(fmt.Sprintf("%v", row["Layer"]))
			rows := &models.DyNamicMenu{
				Name:      fmt.Sprintf("%v", row["Name"]),
				Path:      fmt.Sprintf("%v", row["Path"]),
				Component: fmt.Sprintf("%v", row["Component"]),
				MetaTitle: fmt.Sprintf("%v", row["MetaTitle"]),
				KeepAlive: KeepAlive,
				MetaIcon:  fmt.Sprintf("%v", row["MetaIcon"]),
				Hidden:    hidden,
				Role:      "admin,company",
				Enable: true,
				Layer:LayerInt,
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
