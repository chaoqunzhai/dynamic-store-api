package web_app

//预先定义好一些不变的值

type DiyDefineRow struct {
	Global map[string]interface{}   `json:"global"`
	Value  []map[string]interface{} `json:"value"`
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
	}
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
