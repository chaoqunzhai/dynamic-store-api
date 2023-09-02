package web_app

//预先定义好一些不变的值

type DiyDefineRow struct {
	Global map[string]interface{}   `json:"global"`
	Value  []map[string]interface{} `json:"value"`
}

type WebApp interface {
	//设置Redis数据
	SetLoadRedis() map[string]interface{}
	//orm数据查询
	SearchRun()
}
// 底部菜单公共配置
var NavDefine = `{
        "iconPath": "icondiy icon-system-home",
        "selectedIconPath": "icondiy icon-system-home-selected",
        "text": "主页",
        "link": {
          "name": "INDEX",
          "title": "主页",
          "wap_url": "/pages/index/index",
          "parent": "MALL_LINK"
        },
        "imgWidth": "40",
        "imgHeight": "40",
        "iconClass": "icon-system-home",
        "icon_type": "icon",
        "selected_icon_type": "icon",
        "style": {
          "fontSize": 100,
          "iconBgColor": [],
          "iconBgColorDeg": 0,
          "iconBgImg": "",
          "bgRadius": 0,
          "iconColor": [
            "#000000"
          ],
          "iconColorDeg": 0
        },
        "selected_style": {
          "fontSize": 100,
          "iconBgColor": [],
          "iconBgColorDeg": 0,
          "iconBgImg": "",
          "bgRadius": 0,
          "iconColor": [
            "#F4391c"
          ],
          "iconColorDeg": 0
        }
      }`

// 商品分类
var DIY_VIEW_GOODS_CATEGORY = `
{
	"global":{
		"title": "商品分类",
		"pageBgColor": "#FFFFFF",
		"topNavColor": "#FFFFFF",
		"topNavBg": false,
		"navBarSwitch": true,
		"textNavColor": "#333333",
		"openBottomNav": true,
		"navStyle": 1,
		"textImgPosLink": "left",
		"mpCollect": false,
		"popWindow": {
			"count": -1
		},
		"template": {
			"textColor": "#303133",
			"componentAngle": "round",
			"elementBgColor": "",
			"elementAngle": "round",
			"margin": {
				"top": 0,
				"bottom": 0,
				"both": 0
			}
		}
	},
	"value":[
		{
			"level": "2",
			"template": "2",
			"quickBuy": 1,
			"search": 1,
			"componentName": "GoodsCategory",
			"componentTitle": "商品分类",
			"isDelete": 1,
			"margin": [],
			"goodsLevel": 1,
			"loadType": "part"
		}
	]
}
`

// 我的中心
var DIY_VIEW_MEMBER_INDEX = `
{
	"global":{
		"title": "会员中心",
		"pageBgColor": "#F8F8F8",
		"topNavColor": "#FFFFFF",
		"topNavBg": true,
		"navBarSwitch": true,
		"navStyle": 1,
		"textNavColor": "#333333",
		"topNavImg": "",
		"openBottomNav": true,
		"textImgPosLink": "center",
		"mpCollect": false,
        "popWindow":{
			"imageUrl": "",
			"count": -1,
			"show": 0,
			"link": {
				"name": ""
			},
			"imgWidth": "",
			"imgHeight": ""
		},
		"template": {
			"textColor": "#303133",
			"componentAngle": "round",
			"elementBgColor": "",
			"elementAngle": "round",
			"margin": {
				"top": 0,
				"bottom": 0,
				"both": 0
			}
		}
	},
	"value":[
    {
        "style": 4,
        "theme": "default",
        "bgColorStart": "#FF7230",
        "bgColorEnd": "#FF1544",
        "gradientAngle": "129",
        "infoMargin": 15,
        "id": "1tkaoxbhavj4",
        "addonName": "",
        "componentName": "MemberInfo",
        "componentTitle": "会员信息",
        "isDelete": 0,
        "pageBgColor": "",
        "textColor": "#303133",
        "componentBgColor": "",
        "topAroundRadius": 0,
        "bottomAroundRadius": 0,
        "elementBgColor": "",
        "elementAngle": "round",
        "topElementAroundRadius": 0,
        "bottomElementAroundRadius": 0,
        "margin": {
            "top": 0,
            "bottom": 0,
            "both": 0
        }
    },
    {
        "style": "style-12",
        "styleName": "风格12",
        "text": "我的订单",
        "link": {
            "name": ""
        },
        "fontSize": 17,
        "fontWeight": "bold",
        "subTitle": {
            "fontSize": 14,
            "text": "",
            "isElementShow": true,
            "color": "#999999",
            "bgColor": "#303133"
        },
        "more": {
            "text": "全部订单",
            "link": {
                "name": "ALL_ORDER",
                "title": "全部订单",
                "wap_url": "/pages/order/list",
                "parent": "MALL_LINK"
            },
            "isShow": true,
            "isElementShow": true,
            "color": "#999999"
        },
        "id": "2txcvx3d5u6",
        "addonName": "",
        "componentName": "Text",
        "componentTitle": "标题",
        "isDelete": 0,
        "pageBgColor": "",
        "textColor": "#303133",
        "componentBgColor": "#FFFFFF",
        "componentAngle": "round",
        "topAroundRadius": 9,
        "bottomAroundRadius": 0,
        "topElementAroundRadius": 0,
        "bottomElementAroundRadius": 0,
        "margin": {
            "top": 15,
            "bottom": 0,
            "both": 15
        }
    },
    {
        "componentName": "MemberMyOrder",
        "componentTitle": "我的订单",
        "pageBgColor": "",
        "textColor": "#303133",
        "componentBgColor": "#FFFFFF",
        "componentAngle": "round",
        "bottomAroundRadius": 9,
        "elementAngle": "round",
        "margin": {
            "top": 0,
            "bottom": 0,
            "both": 15
        }
    },
    {
        "style": "style-12",
        "styleName": "风格12",
        "text": "常用工具",
        "link": {
            "name": ""
        },
        "fontSize": 17,
        "fontWeight": "bold",
        "subTitle": {
            "fontSize": 14,
            "text": "",
            "isElementShow": true,
            "color": "#999999",
            "bgColor": "#303133"
        },
        "more": {
            "text": "",
            "link": {
                "name": ""
            },
            "isShow": 0,
            "isElementShow": true,
            "color": "#999999"
        },
        "id": "405rb6vv3rq0",
        "addonName": "",
        "componentName": "Text",
        "componentTitle": "标题",
        "isDelete": 0,
        "pageBgColor": "",
        "textColor": "#303133",
        "componentBgColor": "#FFFFFF",
        "componentAngle": "round",
        "topAroundRadius": 9,
        "bottomAroundRadius": 0,
        "topElementAroundRadius": 0,
        "bottomElementAroundRadius": 0,
        "margin": {
            "top": 15,
            "bottom": 0,
            "both": 15
        }
    },
    {
        "mode": "graphic",
        "type": "img",
        "showStyle": "fixed",
        "ornament": {
            "type": "default",
            "color": "#EDEDED"
        },
        "rowCount": 4,
        "pageCount": 2,
        "carousel": {
            "type": "circle",
            "color": "#FFFFFF"
        },
        "imageSize": 30,
        "aroundRadius": 0,
        "font": {
            "size": 13,
            "weight": "normal",
            "color": "#303133"
        },
        "list": [],
        "componentName": "GraphicNav",
        "componentTitle": "图文导航",
        "pageBgColor": "",
        "componentBgColor": "#FFFFFF",
        "componentAngle": "round",
        "bottomAroundRadius": 9,
        "margin": {
            "top": 0,
            "bottom": 0,
            "both": 15
        }
    }
]
}
`

// 主页
var DIY_VIEW_INDEX = `
{
	"global":{
		"title": "",
		"pageBgColor": "#F6F9FF",
		"topNavColor": "#FFFFFF",
		"topNavBg": true,
		"navBarSwitch": true,
		"navStyle": 1,
		"textNavColor": "#333333",
		"openBottomNav": true,
		"textImgPosLink": "center",
		"mpCollect": false,
		"popWindow": {
			"count": -1,
			"show": 0
		},
		"bgUrl": "addon/diy_default1/bg.png",
		"imgWidth": "2250",
		"imgHeight": "1110",
		"template": {
			"textColor": "#303133",
			"componentAngle": "round",
			"elementAngle": "round",
			"margin": {
				"top": 0,
				"bottom": 0,
				"both": 12
			}
		}
	}
}
`
//如果没有对大B进行特殊配置,那就默认已经配置的
var MemIndexData = `{
      "id": 0,
      "site_id": 0,
      "default":1,
	  "name": "DIY_VIEW_MEMBER_INDEX",
	  "title": "个人中心",
	  "template_id": 55,
	  "template_item_id": 123,
	  "type": "DIY_VIEW_MEMBER_INDEX",
	  "type_name": "会员中心",
	  "value": "{\"global\":{\"title\":\"会员中心\",\"pageBgColor\":\"#F8F8F8\",\"topNavColor\":\"#FFFFFF\",\"topNavBg\":true,\"navBarSwitch\":true,\"navStyle\":1,\"textNavColor\":\"#333333\",\"topNavImg\":\"\",\"openBottomNav\":true,\"textImgPosLink\":\"center\",\"mpCollect\":false,\"template\":{\"textColor\":\"#303133\",\"componentAngle\":\"round\",\"elementBgColor\":\"\",\"elementAngle\":\"round\",\"margin\":{\"top\":0,\"bottom\":0,\"both\":0}}},\"value\":[{\"style\":4,\"theme\":\"default\",\"bgColorStart\":\"#FF7230\",\"bgColorEnd\":\"#FF1544\",\"gradientAngle\":\"129\",\"infoMargin\":15,\"id\":\"1tkaoxbhavj4\",\"addonName\":\"\",\"componentName\":\"MemberInfo\",\"componentTitle\":\"会员信息\",\"isDelete\":0,\"pageBgColor\":\"\",\"textColor\":\"#303133\",\"componentBgColor\":\"\",\"topAroundRadius\":0,\"bottomAroundRadius\":0,\"elementBgColor\":\"\",\"elementAngle\":\"round\",\"topElementAroundRadius\":0,\"bottomElementAroundRadius\":0,\"margin\":{\"top\":0,\"bottom\":0,\"both\":0}},{\"style\":\"style-12\",\"styleName\":\"风格12\",\"text\":\"我的订单\",\"link\":{\"name\":\"\"},\"fontSize\":17,\"fontWeight\":\"bold\",\"subTitle\":{\"fontSize\":14,\"text\":\"\",\"isElementShow\":true,\"color\":\"#999999\",\"bgColor\":\"#303133\"},\"more\":{\"text\":\"全部订单\",\"link\":{\"name\":\"ALL_ORDER\",\"title\":\"全部订单\",\"wap_url\":\"/pages/order/list\",\"parent\":\"MALL_LINK\"},\"isShow\":true,\"isElementShow\":true,\"color\":\"#999999\"},\"id\":\"2txcvx3d5u6\",\"addonName\":\"\",\"componentName\":\"Text\",\"componentTitle\":\"标题\",\"isDelete\":0,\"pageBgColor\":\"\",\"textColor\":\"#303133\",\"componentBgColor\":\"#FFFFFF\",\"componentAngle\":\"round\",\"topAroundRadius\":9,\"bottomAroundRadius\":0,\"topElementAroundRadius\":0,\"bottomElementAroundRadius\":0,\"margin\":{\"top\":15,\"bottom\":0,\"both\":15}},{\"componentName\":\"MemberMyOrder\",\"componentTitle\":\"我的订单\",\"pageBgColor\":\"\",\"textColor\":\"#303133\",\"componentBgColor\":\"#FFFFFF\",\"componentAngle\":\"round\",\"bottomAroundRadius\":9,\"elementAngle\":\"round\",\"margin\":{\"top\":0,\"bottom\":0,\"both\":15}},{\"style\":\"style-12\",\"styleName\":\"风格12\",\"text\":\"常用工具\",\"link\":{\"name\":\"\"},\"fontSize\":17,\"fontWeight\":\"bold\",\"subTitle\":{\"fontSize\":14,\"text\":\"\",\"isElementShow\":true,\"color\":\"#999999\",\"bgColor\":\"#303133\"},\"more\":{\"text\":\"\",\"link\":{\"name\":\"\"},\"isShow\":0,\"isElementShow\":true,\"color\":\"#999999\"},\"id\":\"405rb6vv3rq0\",\"addonName\":\"\",\"componentName\":\"Text\",\"componentTitle\":\"标题\",\"isDelete\":0,\"pageBgColor\":\"\",\"textColor\":\"#303133\",\"componentBgColor\":\"#FFFFFF\",\"componentAngle\":\"round\",\"topAroundRadius\":9,\"bottomAroundRadius\":0,\"topElementAroundRadius\":0,\"bottomElementAroundRadius\":0,\"margin\":{\"top\":15,\"bottom\":0,\"both\":15}},{\"mode\":\"graphic\",\"type\":\"img\",\"showStyle\":\"fixed\",\"ornament\":{\"type\":\"default\",\"color\":\"#EDEDED\"},\"rowCount\":4,\"pageCount\":2,\"carousel\":{\"type\":\"circle\",\"color\":\"#FFFFFF\"},\"imageSize\":30,\"aroundRadius\":0,\"font\":{\"size\":13,\"weight\":\"normal\",\"color\":\"#303133\"},\"list\":[{\"title\":\"个人资料\",\"imageUrl\":\"../../static/member/default_person.png\",\"iconType\":\"img\",\"style\":{\"fontSize\":\"60\",\"iconBgColor\":[],\"iconBgColorDeg\":0,\"iconBgImg\":\"\",\"bgRadius\":0,\"iconColor\":[\"#000000\"],\"iconColorDeg\":0},\"link\":{\"name\":\"MEMBER_INFO\",\"title\":\"个人资料\",\"wap_url\":\"/pages_tool/member/info\",\"parent\":\"MALL_LINK\"},\"label\":{\"control\":false,\"text\":\"热门\",\"textColor\":\"#FFFFFF\",\"bgColorStart\":\"#F83287\",\"bgColorEnd\":\"#FE3423\"},\"icon\":\"\",\"id\":\"10rhv0x6phhc0\"},{\"title\":\"收货地址\",\"imageUrl\":\"../../static/member/default_address.png\",\"iconType\":\"img\",\"style\":{\"fontSize\":\"60\",\"iconBgColor\":[],\"iconBgColorDeg\":0,\"iconBgImg\":\"\",\"bgRadius\":0,\"iconColor\":[\"#000000\"],\"iconColorDeg\":0},\"link\":{\"name\":\"SHIPPING_ADDRESS\",\"title\":\"收货地址\",\"wap_url\":\"/pages_tool/member/address\",\"parent\":\"MALL_LINK\"},\"label\":{\"control\":false,\"text\":\"热门\",\"textColor\":\"#FFFFFF\",\"bgColorStart\":\"#F83287\",\"bgColorEnd\":\"#FE3423\"},\"icon\":\"\",\"id\":\"1n8gycn6xqe80\"},{\"title\":\"优惠券\",\"imageUrl\":\"../../static/member/default_discount.png\",\"iconType\":\"img\",\"style\":\"\",\"link\":{\"name\":\"COUPON\",\"title\":\"优惠券\",\"wap_url\":\"/pages_tool/member/coupon\",\"parent\":\"MALL_LINK\"},\"label\":{\"control\":false,\"text\":\"热门\",\"textColor\":\"#FFFFFF\",\"bgColorStart\":\"#F83287\",\"bgColorEnd\":\"#FE3423\"},\"iconfont\":{\"value\":\"\",\"color\":\"\"},\"id\":\"1tnu0vihrnq80\"}],\"componentName\":\"GraphicNav\",\"componentTitle\":\"图文导航\",\"pageBgColor\":\"\",\"componentBgColor\":\"#FFFFFF\",\"componentAngle\":\"round\",\"bottomAroundRadius\":9,\"margin\":{\"top\":0,\"bottom\":0,\"both\":15}}]}"	}`

//示例demo代码
var IndexDat = `{
    "id": 0,
    "site_id": 0,
    "default":1,
    "name": "DIY_VIEW_INDEX",
    "title": "店铺首页",
    "type": "DIY_VIEW_INDEX",
    "type_name": "店铺首页",
    "value": "{\"global\":{\"title\":\"动创云易订货\",\"pageBgColor\":\"#F6F9FF\",\"topNavColor\":\"#FFFFFF\",\"topNavBg\":true,\"navBarSwitch\":true,\"navStyle\":1,\"textNavColor\":\"#333333\",\"openBottomNav\":true,\"textImgPosLink\":\"center\",\"mpCollect\":false,\"popWindow\":{\"count\":-1,\"show\":0},\"bgUrl\":\"addon/diy_default1/bg.png\",\"imgWidth\":\"2250\",\"imgHeight\":\"1110\",\"template\":{\"textColor\":\"#303133\",\"componentAngle\":\"round\",\"elementAngle\":\"round\",\"margin\":{\"top\":0,\"bottom\":0,\"both\":12}}},\"value\":[{\"componentName\":\"Search\",\"componentTitle\":\"搜索框\",\"margin\":{\"top\":10,\"bottom\":10,\"both\":12},\"title\":\"请输入搜索关键词\",\"textAlign\":\"left\",\"borderType\":2,\"searchStyle\":1,\"textColor\":\"#303133\",\"elementBgColor\":\"#F6F9FF\",\"iconType\":\"img\",\"style\":{\"fontSize\":\"60\",\"iconColor\":[\"#000000\"]},\"positionWay\":\"static\"},{\"list\":[{\"imageUrl\":\"addon/diy_default1/banner.png\",\"imgWidth\":\"750\",\"imgHeight\":\"320\",\"imageMode\":\"scaleToFill\"}],\"indicatorIsShow\":true,\"indicatorColor\":\"#ffffff\",\"carouselStyle\":\"circle\",\"indicatorLocation\":\"center\",\"componentName\":\"ImageAds\",\"componentTitle\":\"图片广告\",\"isDelete\":0,\"componentAngle\":\"round\",\"topAroundRadius\":10,\"bottomAroundRadius\":10,\"margin\":{\"top\":0,\"bottom\":12,\"both\":12}},{\"style\":\"style-16\",\"subTitle\":{\"fontSize\":14,\"text\":\"超级优惠\",\"isElementShow\":true,\"color\":\"#FFFFFF\",\"bgColor\":\"#FF9F29\",\"icon\":\"icondiyicon-system-coupon\",\"fontWeight\":\"bold\"},\"link\":{\"name\":\"COUPON_PREFECTURE\",\"title\":\"优惠券专区\",\"wap_url\":\"/pages_tool/goods/coupon\",\"parent\":\"MARKETING_LINK\"},\"fontSize\":16,\"styleName\":\"风格16\",\"fontWeight\":\"bold\",\"more\":{\"text\":\"\",\"link\":{\"name\":\"COUPON_PREFECTURE\",\"title\":\"优惠券专区\",\"wap_url\":\"/pages_tool/goods/coupon\",\"parent\":\"MARKETING_LINK\"},\"isShow\":true,\"isElementShow\":true,\"color\":\"#999999\"},\"text\":\"优惠专区\",\"componentName\":\"Text\",\"componentTitle\":\"标题\",\"textColor\":\"#303133\",\"componentBgColor\":\"#FFFFFF\",\"componentAngle\":\"round\",\"topAroundRadius\":10,\"margin\":{\"top\":0,\"bottom\":0,\"both\":12}},{\"style\":\"6\",\"sources\":\"initial\",\"styleName\":\"风格六\",\"couponIds\":[],\"count\":6,\"previewList\":[],\"nameColor\":\"#303133\",\"moneyColor\":\"#FF0000\",\"limitColor\":\"#303133\",\"btnStyle\":{\"textColor\":\"#FFFFFF\",\"bgColor\":\"#303133\",\"text\":\"领取\",\"aroundRadius\":20,\"isBgColor\":true,\"isAroundRadius\":true},\"isName\":true,\"couponBgColor\":\"#FFFFFF\",\"couponType\":\"color\",\"ifNeedBg\":true,\"addonName\":\"coupon\",\"componentName\":\"Coupon\",\"componentTitle\":\"优惠券\",\"margin\":{\"top\":0,\"bottom\":0,\"both\":12}},{\"id\":\"4f7eqfqy07s0\",\"list\":[{\"imageUrl\":\"addon/diy_default1/gg.png\",\"imgWidth\":\"702\",\"imgHeight\":\"252\",\"id\":\"1z94aaav9klc0\",\"imageMode\":\"scaleToFill\"}],\"indicatorIsShow\":true,\"indicatorColor\":\"#ffffff\",\"carouselStyle\":\"circle\",\"indicatorLocation\":\"center\",\"componentName\":\"ImageAds\",\"componentTitle\":\"图片广告\",\"componentAngle\":\"round\",\"topAroundRadius\":10,\"bottomAroundRadius\":10,\"margin\":{\"top\":0,\"bottom\":12,\"both\":12}}]}",
    "is_default": 1
  }`