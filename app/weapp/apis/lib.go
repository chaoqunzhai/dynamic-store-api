package apis

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/captcha"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"
	"go-admin/config"
	"io/ioutil"
	"path"
)

type Lib struct {
	api.Api
}

func (e Lib) Config(c *gin.Context) {
	err := e.MakeContext(c).
		MakeOrm().
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	dat := `{
  "cart_count": 0,
  "style_theme": {
    "id": "121",
    "title": "热情红",
    "name": "default",
    "main_color": "#F4391c",
    "aux_color": "#F7B500"
  },
  "diy_bottom_nav": {
    "type": 1,
    "theme": "diy",
    "backgroundColor": "#FFFFFF",
    "textColor": "#333333",
    "textHoverColor": "#F4391c",
    "bulge": true,
    "list": [
      {
        "iconPath": "icondiy icon-system-home",
        "selectedIconPath": "icondiy icon-system-home-selected",
        "text": "主页",
        "link": {
          "name": "INDEX",
          "title": "主页",
          "wap_url": "/pages/index/index",
          "parent": "MALL_LINK"
        },
        "id": "101a31ebhlb40",
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
      },
      {
        "iconPath": "icondiy icon-system-category",
        "selectedIconPath": "icondiy icon-system-category-selected",
        "text": "商品",
        "link": {
          "name": "SHOP_CATEGORY",
          "title": "商品",
          "wap_url": "/pages/goods/category",
          "parent": "MALL_LINK"
        },
        "imgWidth": "40",
        "imgHeight": "40",
        "id": "10wfzebplyxs0",
        "iconClass": "icon-system-category",
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
      },
      {
        "iconPath": "icondiy icon-system-cart",
        "selectedIconPath": "icondiy icon-system-cart-selected",
        "text": "购物车",
        "link": {
          "name": "SHOPPING_TROLLEY",
          "title": "购物车",
          "wap_url": "/pages/goods/cart",
          "parent": "MALL_LINK"
        },
        "imgWidth": "40",
        "imgHeight": "40",
        "id": "14zphcx1tj340",
        "iconClass": "icon-system-cart",
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
      },
      {
        "iconPath": "icondiy icon-system-my",
        "selectedIconPath": "icondiy icon-system-my-selected",
        "text": "我的",
        "link": {
          "name": "MEMBER_CENTER",
          "title": "会员中心",
          "wap_url": "/pages/member/index",
          "parent": "MALL_LINK"
        },
        "imgWidth": "40",
        "imgHeight": "40",
        "id": "b2fuww1h5jk0",
        "iconClass": "icon-system-my",
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
      }
    ],
    "imgType": 2,
    "iconColor": "#333333",
    "iconHoverColor": "#FF4D4D"
  },
  "addon_is_exist": {
    "fenxiao": 1,
    "pintuan": 1,
    "membersignin": 1,
    "memberrecharge": 1,
    "memberwithdraw": 1,
    "pointexchange": 1,
    "manjian": 1,
    "memberconsume": 1,
    "memberregister": 1,
    "coupon": 1,
    "bundling": 1,
    "discount": 1,
    "seckill": 1,
    "topic": 1,
    "store": 0,
    "groupbuy": 1,
    "bargain": 1,
    "presale": 1,
    "notes": 1,
    "membercancel": 1,
    "servicer": 1,
    "live": 1,
    "cards": 1,
    "egg": 1,
    "turntable": 1,
    "memberrecommend": 1,
    "supermember": 1,
    "giftcard": 1,
    "divideticket": 1,
    "birthdaygift": 1,
    "scenefestival": 1,
    "pinfan": 1,
    "hongbao": 1,
    "blindbox": 1,
    "virtualcard": 1,
    "cardservice": 1,
    "cashier": 1,
    "form": 1
  },
  "default_img": {
    "goods": "public/static/img/default_img/square.png",
    "head": "public/static/img/default_img/head.png",
    "store": "public/static/img/default_img/square.png",
    "article": "public/static/img/default_img/article.png"
  },
  "copyright": {
    "icp": "备案号: 222222",
    "gov_record": "",
    "gov_url": "",
    "market_supervision_url": "",
    "company_name": "动创云",
    "copyright_link": "",
    "copyright_desc": "动创云",
    "auth": true
  },
  "site_info": {
    "site_id": 1,
    "site_domain": "",
    "site_name": "动创云订货软件",
    "logo": "../static/logo.png",
    "seo_title": "",
    "seo_keywords": "动创云订货软件",
    "seo_description": "动创云订货软件",
    "site_tel": "",
    "logo_square": "",
    "shop_status": "1"
  },
  "servicer": {
    "h5": {
      "type": "dongchuangyun",
      "wxwork_url": "https://dongchuangyun.com/",
      "third_url": "https://dongchuangyun.com/"
    },
    "weapp": {
      "type": "dynamic-app",
      "corpid": "",
      "wxwork_url": ""
    },
    "pc": {
      "type": "third",
      "third_url": "http://www.baidu.com"
    },
    "aliapp": {
      "type": "none"
    }
  },
  "store_config": {
    "store_business": "shop"
  }
}`
	row := make(map[string]interface{}, 0)
	marErr:=json.Unmarshal([]byte(dat), &row)
	fmt.Println("marErr",marErr)
	e.OK(row, "successful")
	return
}

type DiyInfoRequest struct {
	Name      string     `form:"name" comment:"视角名称"`      //显示名称
}
func (e Lib) DiyInfo(c *gin.Context) {
	req :=DiyInfoRequest{}
	err := e.MakeContext(c).
		MakeOrm().
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	if bindErr := c.ShouldBind(&req); bindErr != nil {
		e.Error(500, bindErr, bindErr.Error())
		return
	}
	MemDat := `{
	  "id": 49,
	  "site_id": 1,
	  "name": "DIY_VIEW_MEMBER_INDEX",
	  "title": "会员中心",
	  "template_id": 55,
	  "template_item_id": 123,
	  "type": "DIY_VIEW_MEMBER_INDEX",
	  "type_name": "会员中心",
	  "value": "{\"global\":{\"title\":\"\\u4f1a\\u5458\\u4e2d\\u5fc3\",\"pageBgColor\":\"#F8F8F8\",\"topNavColor\":\"#FFFFFF\",\"topNavBg\":true,\"navBarSwitch\":true,\"navStyle\":1,\"textNavColor\":\"#333333\",\"topNavImg\":\"\",\"moreLink\":{\"name\":\"\"},\"openBottomNav\":true,\"textImgPosLink\":\"center\",\"mpCollect\":false,\"popWindow\":{\"imageUrl\":\"\",\"count\":-1,\"show\":0,\"link\":{\"name\":\"\"},\"imgWidth\":\"\",\"imgHeight\":\"\"},\"bgUrl\":\"\",\"imgWidth\":\"\",\"imgHeight\":\"\",\"template\":{\"pageBgColor\":\"\",\"textColor\":\"#303133\",\"componentBgColor\":\"\",\"componentAngle\":\"round\",\"topAroundRadius\":0,\"bottomAroundRadius\":0,\"elementBgColor\":\"\",\"elementAngle\":\"round\",\"topElementAroundRadius\":0,\"bottomElementAroundRadius\":0,\"margin\":{\"top\":0,\"bottom\":0,\"both\":0}}},\"value\":[{\"style\":4,\"theme\":\"default\",\"bgColorStart\":\"#FF7230\",\"bgColorEnd\":\"#FF1544\",\"gradientAngle\":\"129\",\"infoMargin\":15,\"id\":\"1tkaoxbhavj4\",\"addonName\":\"\",\"componentName\":\"MemberInfo\",\"componentTitle\":\"\\u4f1a\\u5458\\u4fe1\\u606f\",\"isDelete\":0,\"pageBgColor\":\"\",\"textColor\":\"#303133\",\"componentBgColor\":\"\",\"topAroundRadius\":0,\"bottomAroundRadius\":0,\"elementBgColor\":\"\",\"elementAngle\":\"round\",\"topElementAroundRadius\":0,\"bottomElementAroundRadius\":0,\"margin\":{\"top\":0,\"bottom\":0,\"both\":0}},{\"style\":\"style-12\",\"styleName\":\"\\u98ce\\u683c12\",\"text\":\"\\u6211\\u7684\\u8ba2\\u5355\",\"link\":{\"name\":\"\"},\"fontSize\":17,\"fontWeight\":\"bold\",\"subTitle\":{\"fontSize\":14,\"text\":\"\",\"isElementShow\":true,\"color\":\"#999999\",\"bgColor\":\"#303133\"},\"more\":{\"text\":\"\\u5168\\u90e8\\u8ba2\\u5355\",\"link\":{\"name\":\"ALL_ORDER\",\"title\":\"\\u5168\\u90e8\\u8ba2\\u5355\",\"wap_url\":\"\\/pages\\/order\\/list\",\"parent\":\"MALL_LINK\"},\"isShow\":true,\"isElementShow\":true,\"color\":\"#999999\"},\"id\":\"2txcvx3d5u6\",\"addonName\":\"\",\"componentName\":\"Text\",\"componentTitle\":\"\\u6807\\u9898\",\"isDelete\":0,\"pageBgColor\":\"\",\"textColor\":\"#303133\",\"componentBgColor\":\"#FFFFFF\",\"componentAngle\":\"round\",\"topAroundRadius\":9,\"bottomAroundRadius\":0,\"topElementAroundRadius\":0,\"bottomElementAroundRadius\":0,\"margin\":{\"top\":15,\"bottom\":0,\"both\":15}},{\"color\":\"#EEEEEE\",\"borderStyle\":\"solid\",\"id\":\"3hsh2st470e0\",\"addonName\":\"\",\"componentName\":\"HorzLine\",\"componentTitle\":\"\\u8f85\\u52a9\\u7ebf\",\"isDelete\":0,\"pageBgColor\":\"\",\"topAroundRadius\":0,\"bottomAroundRadius\":0,\"topElementAroundRadius\":0,\"bottomElementAroundRadius\":0,\"margin\":{\"top\":0,\"bottom\":0,\"both\":20}},{\"icon\":{\"waitPay\":{\"title\":\"\\u5f85\\u4ed8\\u6b3e\",\"icon\":\"icondiy icon-system-daifukuan2\",\"style\":{\"bgRadius\":0,\"fontSize\":65,\"iconBgColor\":[],\"iconBgColorDeg\":0,\"iconBgImg\":\"\",\"iconColor\":[\"#fa9c8e\",\"#F4391c\"],\"iconColorDeg\":0}},\"waitSend\":{\"title\":\"\\u5f85\\u53d1\\u8d27\",\"icon\":\"icondiy icon-system-daifahuo2\",\"style\":{\"bgRadius\":0,\"fontSize\":65,\"iconBgColor\":[],\"iconBgColorDeg\":0,\"iconBgImg\":\"\",\"iconColor\":[\"#fa9c8e\",\"#F4391c\"],\"iconColorDeg\":0}},\"waitConfirm\":{\"title\":\"\\u5f85\\u6536\\u8d27\",\"icon\":\"icondiy icon-system-daishouhuo2\",\"style\":{\"bgRadius\":0,\"fontSize\":65,\"iconBgColor\":[],\"iconBgColorDeg\":0,\"iconBgImg\":\"\",\"iconColor\":[\"#fa9c8e\",\"#F4391c\"],\"iconColorDeg\":0}},\"waitUse\":{\"title\":\"\\u5f85\\u4f7f\\u7528\",\"icon\":\"icondiy icon-system-daishiyong2\",\"style\":{\"bgRadius\":0,\"fontSize\":65,\"iconBgColor\":[],\"iconBgColorDeg\":0,\"iconBgImg\":\"\",\"iconColor\":[\"#fa9c8e\",\"#F4391c\"],\"iconColorDeg\":0}},\"refunding\":{\"title\":\"\\u552e\\u540e\",\"icon\":\"icondiy icon-system-shuhou2\",\"style\":{\"bgRadius\":0,\"fontSize\":65,\"iconBgColor\":[],\"iconBgColorDeg\":0,\"iconBgImg\":\"\",\"iconColor\":[\"#fa9c8e\",\"#F4391c\"],\"iconColorDeg\":0}}},\"style\":1,\"id\":\"51h05xpcanw0\",\"addonName\":\"\",\"componentName\":\"MemberMyOrder\",\"componentTitle\":\"\\u6211\\u7684\\u8ba2\\u5355\",\"isDelete\":0,\"pageBgColor\":\"\",\"textColor\":\"#303133\",\"componentBgColor\":\"#FFFFFF\",\"componentAngle\":\"round\",\"topAroundRadius\":0,\"bottomAroundRadius\":9,\"elementBgColor\":\"\",\"elementAngle\":\"round\",\"topElementAroundRadius\":0,\"bottomElementAroundRadius\":0,\"margin\":{\"top\":0,\"bottom\":0,\"both\":15}},{\"style\":\"style-12\",\"styleName\":\"\\u98ce\\u683c12\",\"text\":\"\\u5e38\\u7528\\u5de5\\u5177\",\"link\":{\"name\":\"\"},\"fontSize\":17,\"fontWeight\":\"bold\",\"subTitle\":{\"fontSize\":14,\"text\":\"\",\"isElementShow\":true,\"color\":\"#999999\",\"bgColor\":\"#303133\"},\"more\":{\"text\":\"\",\"link\":{\"name\":\"\"},\"isShow\":0,\"isElementShow\":true,\"color\":\"#999999\"},\"id\":\"405rb6vv3rq0\",\"addonName\":\"\",\"componentName\":\"Text\",\"componentTitle\":\"\\u6807\\u9898\",\"isDelete\":0,\"pageBgColor\":\"\",\"textColor\":\"#303133\",\"componentBgColor\":\"#FFFFFF\",\"componentAngle\":\"round\",\"topAroundRadius\":9,\"bottomAroundRadius\":0,\"topElementAroundRadius\":0,\"bottomElementAroundRadius\":0,\"margin\":{\"top\":15,\"bottom\":0,\"both\":15}},{\"mode\":\"graphic\",\"type\":\"img\",\"showStyle\":\"fixed\",\"ornament\":{\"type\":\"default\",\"color\":\"#EDEDED\"},\"rowCount\":4,\"pageCount\":2,\"carousel\":{\"type\":\"circle\",\"color\":\"#FFFFFF\"},\"imageSize\":30,\"aroundRadius\":0,\"font\":{\"size\":13,\"weight\":\"normal\",\"color\":\"#303133\"},\"list\":[{\"title\":\"\\u4e2a\\u4eba\\u8d44\\u6599\",\"imageUrl\":\"..\\/..\\/static\\/member\\/default_person.png\",\"iconType\":\"img\",\"style\":{\"fontSize\":\"60\",\"iconBgColor\":[],\"iconBgColorDeg\":0,\"iconBgImg\":\"\",\"bgRadius\":0,\"iconColor\":[\"#000000\"],\"iconColorDeg\":0},\"link\":{\"name\":\"MEMBER_INFO\",\"title\":\"\\u4e2a\\u4eba\\u8d44\\u6599\",\"wap_url\":\"\\/pages_tool\\/member\\/info\",\"parent\":\"MALL_LINK\"},\"label\":{\"control\":false,\"text\":\"\\u70ed\\u95e8\",\"textColor\":\"#FFFFFF\",\"bgColorStart\":\"#F83287\",\"bgColorEnd\":\"#FE3423\"},\"icon\":\"\",\"id\":\"10rhv0x6phhc0\"},{\"title\":\"\\u6536\\u8d27\\u5730\\u5740\",\"imageUrl\":\"..\\/..\\/static\\/member\\/default_address.png\",\"iconType\":\"img\",\"style\":{\"fontSize\":\"60\",\"iconBgColor\":[],\"iconBgColorDeg\":0,\"iconBgImg\":\"\",\"bgRadius\":0,\"iconColor\":[\"#000000\"],\"iconColorDeg\":0},\"link\":{\"name\":\"SHIPPING_ADDRESS\",\"title\":\"\\u6536\\u8d27\\u5730\\u5740\",\"wap_url\":\"\\/pages_tool\\/member\\/address\",\"parent\":\"MALL_LINK\"},\"label\":{\"control\":false,\"text\":\"\\u70ed\\u95e8\",\"textColor\":\"#FFFFFF\",\"bgColorStart\":\"#F83287\",\"bgColorEnd\":\"#FE3423\"},\"icon\":\"\",\"id\":\"1n8gycn6xqe80\"},{\"title\":\"\\u6211\\u7684\\u5173\\u6ce8\",\"imageUrl\":\"..\\/..\\/static\\/member\\/default_like.png\",\"iconType\":\"img\",\"style\":{\"fontSize\":\"60\",\"iconBgColor\":[],\"iconBgColorDeg\":0,\"iconBgImg\":\"\",\"bgRadius\":0,\"iconColor\":[\"#000000\"],\"iconColorDeg\":0},\"link\":{\"name\":\"ATTENTION\",\"title\":\"\\u6211\\u7684\\u5173\\u6ce8\",\"wap_url\":\"\\/pages_tool\\/member\\/collection\",\"parent\":\"MALL_LINK\"},\"label\":{\"control\":false,\"text\":\"\\u70ed\\u95e8\",\"textColor\":\"#FFFFFF\",\"bgColorStart\":\"#F83287\",\"bgColorEnd\":\"#FE3423\"},\"icon\":\"\",\"id\":\"cnamoch6cvk0\"},{\"title\":\"\\u6211\\u7684\\u8db3\\u8ff9\",\"imageUrl\":\"..\\/..\\/static\\/member\\/default_toot.png\",\"iconType\":\"img\",\"style\":{\"fontSize\":\"60\",\"iconBgColor\":[],\"iconBgColorDeg\":0,\"iconBgImg\":\"\",\"bgRadius\":0,\"iconColor\":[\"#000000\"],\"iconColorDeg\":0},\"link\":{\"name\":\"FOOTPRINT\",\"title\":\"\\u6211\\u7684\\u8db3\\u8ff9\",\"wap_url\":\"\\/pages_tool\\/member\\/footprint\",\"parent\":\"MALL_LINK\"},\"label\":{\"control\":false,\"text\":\"\\u70ed\\u95e8\",\"textColor\":\"#FFFFFF\",\"bgColorStart\":\"#F83287\",\"bgColorEnd\":\"#FE3423\"},\"icon\":\"\",\"id\":\"drf3hi3slo00\"},{\"title\":\"\\u8d26\\u6237\\u5217\\u8868\",\"imageUrl\":\"..\\/..\\/static\\/member\\/default_cash.png\",\"iconType\":\"img\",\"style\":\"\",\"link\":{\"name\":\"ACCOUNT\",\"title\":\"\\u8d26\\u6237\\u5217\\u8868\",\"wap_url\":\"\\/pages_tool\\/member\\/account\",\"parent\":\"MALL_LINK\"},\"label\":{\"control\":false,\"text\":\"\\u70ed\\u95e8\",\"textColor\":\"#FFFFFF\",\"bgColorStart\":\"#F83287\",\"bgColorEnd\":\"#FE3423\"},\"iconfont\":{\"value\":\"\",\"color\":\"\"},\"id\":\"1l4axfhbayqo0\"},{\"title\":\"\\u4f18\\u60e0\\u5238\",\"imageUrl\":\"..\\/..\\/static\\/member\\/default_discount.png\",\"iconType\":\"img\",\"style\":\"\",\"link\":{\"name\":\"COUPON\",\"title\":\"\\u4f18\\u60e0\\u5238\",\"wap_url\":\"\\/pages_tool\\/member\\/coupon\",\"parent\":\"MALL_LINK\"},\"label\":{\"control\":false,\"text\":\"\\u70ed\\u95e8\",\"textColor\":\"#FFFFFF\",\"bgColorStart\":\"#F83287\",\"bgColorEnd\":\"#FE3423\"},\"iconfont\":{\"value\":\"\",\"color\":\"\"},\"id\":\"1tnu0vihrnq80\"},{\"title\":\"\\u7b7e\\u5230\",\"imageUrl\":\"..\\/..\\/static\\/member\\/default_sign.png\",\"iconType\":\"img\",\"style\":\"\",\"link\":{\"name\":\"SIGN_IN\",\"title\":\"\\u7b7e\\u5230\",\"wap_url\":\"\\/pages_tool\\/member\\/signin\",\"parent\":\"MARKETING_LINK\"},\"label\":{\"control\":false,\"text\":\"\\u70ed\\u95e8\",\"textColor\":\"#FFFFFF\",\"bgColorStart\":\"#F83287\",\"bgColorEnd\":\"#FE3423\"},\"iconfont\":{\"value\":\"\",\"color\":\"\"},\"id\":\"hodjcxowf8g0\"},{\"title\":\"\\u6211\\u7684\\u62fc\\u5355\",\"imageUrl\":\"..\\/..\\/static\\/member\\/default_store.png\",\"iconType\":\"img\",\"style\":\"\",\"link\":{\"name\":\"MY_PINTUAN\",\"title\":\"\\u6211\\u7684\\u62fc\\u56e2\",\"wap_url\":\"\\/pages_promotion\\/pintuan\\/my_spell\",\"parent\":\"MARKETING_LINK\"},\"label\":{\"control\":false,\"text\":\"\\u70ed\\u95e8\",\"textColor\":\"#FFFFFF\",\"bgColorStart\":\"#F83287\",\"bgColorEnd\":\"#FE3423\"},\"iconfont\":{\"value\":\"\",\"color\":\"\"},\"id\":\"uoarcfsleio0\"},{\"title\":\"\\u79ef\\u5206\\u5151\\u6362\",\"imageUrl\":\"..\\/..\\/static\\/member\\/default_point_recond.png\",\"iconType\":\"img\",\"style\":\"\",\"link\":{\"name\":\"INTEGRAL_CONVERSION\",\"title\":\"\\u79ef\\u5206\\u5151\\u6362\",\"wap_url\":\"\\/pages_promotion\\/point\\/order_list\",\"parent\":\"MARKETING_LINK\"},\"label\":{\"control\":false,\"text\":\"\\u70ed\\u95e8\",\"textColor\":\"#FFFFFF\",\"bgColorStart\":\"#F83287\",\"bgColorEnd\":\"#FE3423\"},\"iconfont\":{\"value\":\"\",\"color\":\"\"},\"id\":\"rnyw8xo5rdc0\"},{\"title\":\"\\u5206\\u9500\\u4e2d\\u5fc3\",\"imageUrl\":\"..\\/..\\/static\\/member\\/default_fenxiao.png\",\"iconType\":\"img\",\"style\":\"\",\"link\":{\"name\":\"DISTRIBUTION_CENTRE\",\"title\":\"\\u5206\\u9500\\u4e2d\\u5fc3\",\"wap_url\":\"\\/pages_promotion\\/fenxiao\\/index\",\"parent\":\"MARKETING_LINK\"},\"label\":{\"control\":false,\"text\":\"\\u70ed\\u95e8\",\"textColor\":\"#FFFFFF\",\"bgColorStart\":\"#F83287\",\"bgColorEnd\":\"#FE3423\"},\"iconfont\":{\"value\":\"\",\"color\":\"\"},\"id\":\"yevac1grnlc0\"},{\"title\":\"\\u6211\\u7684\\u780d\\u4ef7\",\"imageUrl\":\"..\\/..\\/static\\/member\\/default_bargain.png\",\"iconType\":\"img\",\"style\":\"\",\"link\":{\"name\":\"MY_BARGAIN\",\"title\":\"\\u6211\\u7684\\u780d\\u4ef7\",\"wap_url\":\"\\/pages_promotion\\/bargain\\/my_bargain\",\"parent\":\"MARKETING_LINK\"},\"label\":{\"control\":false,\"text\":\"\\u70ed\\u95e8\",\"textColor\":\"#FFFFFF\",\"bgColorStart\":\"#F83287\",\"bgColorEnd\":\"#FE3423\"},\"iconfont\":{\"value\":\"\",\"color\":\"\"},\"id\":\"13uz22sbag000\"},{\"title\":\"\\u9080\\u8bf7\\u6709\\u793c\",\"imageUrl\":\"..\\/..\\/static\\/member\\/default_memberrecommend.png\",\"iconType\":\"img\",\"style\":\"\",\"link\":{\"name\":\"MEMBER_RECOMMEND\",\"title\":\"\\u9080\\u8bf7\\u6709\\u793c\",\"wap_url\":\"\\/pages_tool\\/member\\/invite_friends\",\"parent\":\"MARKETING_LINK\"},\"label\":{\"control\":false,\"text\":\"\\u70ed\\u95e8\",\"textColor\":\"#FFFFFF\",\"bgColorStart\":\"#F83287\",\"bgColorEnd\":\"#FE3423\"},\"iconfont\":{\"value\":\"\",\"color\":\"\"},\"id\":\"1h34nmfisge80\"},{\"title\":\"\\u6211\\u7684\\u9884\\u552e\",\"imageUrl\":\"..\\/..\\/static\\/member\\/my_presale.png\",\"iconType\":\"img\",\"style\":\"\",\"link\":{\"name\":\"PRESALE\",\"title\":\"\\u6211\\u7684\\u9884\\u552e\",\"wap_url\":\"\\/pages_promotion\\/presale\\/order_list\",\"parent\":\"MARKETING_LINK\"},\"label\":{\"control\":false,\"text\":\"\\u70ed\\u95e8\",\"textColor\":\"#FFFFFF\",\"bgColorStart\":\"#F83287\",\"bgColorEnd\":\"#FE3423\"},\"iconfont\":{\"value\":\"\",\"color\":\"\"},\"id\":\"1a3cqyziwqdc0\"},{\"title\":\"\\u6211\\u7684\\u793c\\u54c1\\u5361\",\"imageUrl\":\"..\\/..\\/static\\/member\\/my_giftcard.png\",\"iconType\":\"img\",\"style\":\"\",\"link\":{\"name\":\"GIFTCARD\",\"title\":\"\\u6211\\u7684\\u793c\\u54c1\\u5361\",\"wap_url\":\"\\/pages_promotion\\/giftcard\\/member\",\"parent\":\"MARKETING_LINK\"},\"label\":{\"control\":false,\"text\":\"\\u70ed\\u95e8\",\"textColor\":\"#FFFFFF\",\"bgColorStart\":\"#F83287\",\"bgColorEnd\":\"#FE3423\"},\"iconfont\":{\"value\":\"\",\"color\":\"\"},\"id\":\"1es42rgg2mhs0\"},{\"title\":\"\\u597d\\u53cb\\u74dc\\u5206\\u5238\",\"imageUrl\":\"..\\/..\\/static\\/member\\/my_divideticket.png\",\"iconType\":\"img\",\"style\":\"\",\"link\":{\"name\":\"DIVIDETICKET\",\"title\":\"\\u6211\\u7684\\u597d\\u53cb\\u74dc\\u5206\\u5238\",\"wap_url\":\"\\/pages_promotion\\/divideticket\\/my_guafen\",\"parent\":\"MARKETING_LINK\"},\"label\":{\"control\":false,\"text\":\"\\u70ed\\u95e8\",\"textColor\":\"#FFFFFF\",\"bgColorStart\":\"#F83287\",\"bgColorEnd\":\"#FE3423\"},\"iconfont\":{\"value\":\"\",\"color\":\"\"},\"id\":\"14rxg7u5yu2k0\"},{\"title\":\"\\u62fc\\u56e2\\u8fd4\\u5229\",\"imageUrl\":\"..\\/..\\/static\\/member\\/my_pinfan.png\",\"iconType\":\"img\",\"style\":\"\",\"link\":{\"name\":\"PINFAN\",\"title\":\"\\u6211\\u7684\\u62fc\\u56e2\\u8fd4\\u5229\",\"wap_url\":\"\\/pages_promotion\\/pinfan\\/my_rebate\",\"parent\":\"MARKETING_LINK\"},\"label\":{\"control\":false,\"text\":\"\\u70ed\\u95e8\",\"textColor\":\"#FFFFFF\",\"bgColorStart\":\"#F83287\",\"bgColorEnd\":\"#FE3423\"},\"iconfont\":{\"value\":\"\",\"color\":\"\"},\"id\":\"20eeayo377xc0\"},{\"title\":\"\\u88c2\\u53d8\\u7ea2\\u5305\",\"imageUrl\":\"..\\/..\\/static\\/member\\/my_hongbao.png\",\"iconType\":\"img\",\"style\":\"\",\"link\":{\"name\":\"HONGBAO\",\"title\":\"\\u6211\\u7684\\u88c2\\u53d8\\u7ea2\\u5305\",\"wap_url\":\"\\/pages_tool\\/hongbao\\/my_hongbao\",\"parent\":\"MARKETING_LINK\"},\"label\":{\"control\":false,\"text\":\"\\u70ed\\u95e8\",\"textColor\":\"#FFFFFF\",\"bgColorStart\":\"#F83287\",\"bgColorEnd\":\"#FE3423\"},\"iconfont\":{\"value\":\"\",\"color\":\"\"},\"id\":\"nbthjkdt5c00\"},{\"title\":\"\\u76f2\\u76d2\",\"imageUrl\":\"..\\/..\\/static\\/member\\/my_box.png\",\"iconType\":\"img\",\"style\":\"\",\"link\":{\"name\":\"BLINDBOX\",\"title\":\"\\u6211\\u7684\\u76f2\\u76d2\",\"wap_url\":\"\\/pages_promotion\\/blindbox\\/my_box\",\"parent\":\"MARKETING_LINK\"},\"label\":{\"control\":false,\"text\":\"\\u70ed\\u95e8\",\"textColor\":\"#FFFFFF\",\"bgColorStart\":\"#F83287\",\"bgColorEnd\":\"#FE3423\"},\"iconfont\":{\"value\":\"\",\"color\":\"\"},\"id\":\"1i61b1fdeasg0\"},{\"title\":\"\\u6838\\u9500\\u53f0\",\"icon\":\"\",\"imageUrl\":\"..\\/..\\/static\\/member\\/hexiao.png\",\"iconType\":\"img\",\"style\":{\"fontSize\":\"60\",\"iconBgColor\":[],\"iconBgColorDeg\":0,\"iconBgImg\":\"\",\"bgRadius\":0,\"iconColor\":[\"#000000\"],\"iconColorDeg\":0},\"link\":{\"name\":\"VERIFICATION_PLATFORM\",\"title\":\"\\u6838\\u9500\\u53f0\",\"wap_url\":\"\\/pages_tool\\/verification\\/index\",\"parent\":\"MALL_LINK\"},\"label\":{\"control\":false,\"text\":\"\\u70ed\\u95e8\",\"textColor\":\"#FFFFFF\",\"bgColorStart\":\"#F83287\",\"bgColorEnd\":\"#FE3423\"},\"id\":\"12ktkqu49q680\",\"imgWidth\":\"60\",\"imgHeight\":\"60\"}],\"id\":\"5ywbzsnigpw0\",\"addonName\":\"\",\"componentName\":\"GraphicNav\",\"componentTitle\":\"\\u56fe\\u6587\\u5bfc\\u822a\",\"isDelete\":0,\"pageBgColor\":\"\",\"componentBgColor\":\"#FFFFFF\",\"componentAngle\":\"round\",\"topAroundRadius\":0,\"bottomAroundRadius\":9,\"topElementAroundRadius\":0,\"bottomElementAroundRadius\":0,\"margin\":{\"top\":0,\"bottom\":0,\"both\":15}}]}",
	  "is_default": 1
	}`
	IndexDat := `{
    "id": 37,
    "site_id": 1,
    "name": "DIY_VIEW_INDEX",
    "title": "动创云易订货",
    "template_id": 52,
    "template_item_id": 116,
    "type": "DIY_VIEW_INDEX",
    "type_name": "店铺首页",
    "value": "{\"global\":{\"title\":\"\\u65f6\\u5c1a\\u7b80\\u7ea6\\u5546\\u57ce\",\"pageBgColor\":\"#F6F9FF\",\"topNavColor\":\"#FFFFFF\",\"topNavBg\":true,\"navBarSwitch\":true,\"navStyle\":1,\"textNavColor\":\"#333333\",\"topNavImg\":\"\",\"moreLink\":{\"name\":\"\"},\"openBottomNav\":true,\"textImgPosLink\":\"center\",\"mpCollect\":false,\"popWindow\":{\"imageUrl\":\"\",\"count\":-1,\"show\":0,\"link\":{\"name\":\"\"},\"imgWidth\":\"\",\"imgHeight\":\"\"},\"bgUrl\":\"addon\\/diy_default1\\/bg.png\",\"imgWidth\":\"2250\",\"imgHeight\":\"1110\",\"template\":{\"pageBgColor\":\"\",\"textColor\":\"#303133\",\"componentBgColor\":\"\",\"componentAngle\":\"round\",\"topAroundRadius\":0,\"bottomAroundRadius\":0,\"elementBgColor\":\"\",\"elementAngle\":\"round\",\"topElementAroundRadius\":0,\"bottomElementAroundRadius\":0,\"margin\":{\"top\":0,\"bottom\":0,\"both\":12}}},\"value\":[{\"id\":\"1ufba32yuxz4\",\"addonName\":\"\",\"componentName\":\"Search\",\"componentTitle\":\"\\u641c\\u7d22\\u6846\",\"isDelete\":0,\"topAroundRadius\":0,\"bottomAroundRadius\":0,\"topElementAroundRadius\":0,\"bottomElementAroundRadius\":0,\"margin\":{\"top\":10,\"bottom\":10,\"both\":12},\"title\":\"\\u8bf7\\u8f93\\u5165\\u641c\\u7d22\\u5173\\u952e\\u8bcd\",\"textAlign\":\"left\",\"borderType\":2,\"searchImg\":\"\",\"searchStyle\":1,\"searchLink\":{\"name\":\"\"},\"pageBgColor\":\"\",\"textColor\":\"#303133\",\"componentBgColor\":\"\",\"elementBgColor\":\"#F6F9FF\",\"iconType\":\"img\",\"icon\":\"\",\"style\":{\"fontSize\":\"60\",\"iconBgColor\":[],\"iconBgColorDeg\":0,\"iconBgImg\":\"\",\"bgRadius\":0,\"iconColor\":[\"#000000\"],\"iconColorDeg\":0},\"imageUrl\":\"\",\"positionWay\":\"static\"},{\"id\":\"3tzix3re8wo0\",\"list\":[{\"link\":{\"name\":\"\"},\"imageUrl\":\"addon\\/diy_default1\\/banner.png\",\"imgWidth\":\"750\",\"imgHeight\":\"320\",\"id\":\"1iy3xvq2ngf40\",\"imageMode\":\"scaleToFill\"}],\"indicatorIsShow\":true,\"indicatorColor\":\"#ffffff\",\"carouselStyle\":\"circle\",\"indicatorLocation\":\"center\",\"addonName\":\"\",\"componentName\":\"ImageAds\",\"componentTitle\":\"\\u56fe\\u7247\\u5e7f\\u544a\",\"isDelete\":0,\"pageBgColor\":\"\",\"componentBgColor\":\"\",\"componentAngle\":\"round\",\"topAroundRadius\":10,\"bottomAroundRadius\":10,\"topElementAroundRadius\":0,\"bottomElementAroundRadius\":0,\"margin\":{\"top\":0,\"bottom\":12,\"both\":12}},{\"id\":\"42xi3odl9m60\",\"mode\":\"graphic\",\"type\":\"img\",\"showStyle\":\"fixed\",\"ornament\":{\"type\":\"default\",\"color\":\"#EDEDED\"},\"rowCount\":5,\"pageCount\":2,\"carousel\":{\"type\":\"circle\",\"color\":\"#FFFFFF\"},\"imageSize\":40,\"aroundRadius\":25,\"font\":{\"size\":14,\"weight\":\"normal\",\"color\":\"#303133\"},\"list\":[{\"title\":\"\\u56e2\\u8d2d\",\"icon\":\"icondiy icon-system-groupbuy-nav\",\"imageUrl\":\"\",\"iconType\":\"icon\",\"style\":{\"fontSize\":50,\"iconBgColor\":[\"#FF9F3E\",\"#FF4116\"],\"iconBgColorDeg\":90,\"bgRadius\":50,\"iconColor\":[\"#FFFFFF\"],\"iconColorDeg\":0},\"link\":{\"name\":\"GROUPBUY_PREFECTURE\",\"title\":\"\\u56e2\\u8d2d\\u4e13\\u533a\",\"wap_url\":\"\\/pages_promotion\\/groupbuy\\/list\",\"parent\":\"MARKETING_LINK\"},\"label\":{\"control\":false,\"text\":\"\\u70ed\\u95e8\",\"textColor\":\"#FFFFFF\",\"bgColorStart\":\"#F83287\",\"bgColorEnd\":\"#FE3423\"},\"id\":\"ycafod7gfgg0\"},{\"title\":\"\\u62fc\\u56e2\",\"icon\":\"icondiy icon-system-pintuan-nav\",\"imageUrl\":\"\",\"iconType\":\"icon\",\"style\":{\"fontSize\":50,\"iconBgColor\":[\"#58BCFF\",\"#1379FF\"],\"iconBgColorDeg\":90,\"iconBgImg\":\"public\\/static\\/ext\\/diyview\\/img\\/icon_bg\\/bg_06.png\",\"bgRadius\":50,\"iconColor\":[\"#FFFFFF\"],\"iconColorDeg\":0},\"link\":{\"name\":\"PINTUAN_PREFECTURE\",\"title\":\"\\u62fc\\u56e2\\u4e13\\u533a\",\"wap_url\":\"\\/pages_promotion\\/pintuan\\/list\",\"parent\":\"MARKETING_LINK\"},\"label\":{\"control\":false,\"text\":\"\\u70ed\\u95e8\",\"textColor\":\"#FFFFFF\",\"bgColorStart\":\"#F83287\",\"bgColorEnd\":\"#FE3423\"},\"id\":\"wnlf5ak6u8g0\"},{\"title\":\"\\u79d2\\u6740\",\"icon\":\"icondiy icon-system-seckill-time\",\"imageUrl\":\"\",\"iconType\":\"icon\",\"style\":{\"fontSize\":50,\"iconBgColor\":[\"#FFCC26\",\"#FF9F29\"],\"iconBgColorDeg\":90,\"iconBgImg\":\"public\\/static\\/ext\\/diyview\\/img\\/icon_bg\\/bg_06.png\",\"bgRadius\":50,\"iconColor\":[\"#FFFFFF\"],\"iconColorDeg\":0},\"link\":{\"name\":\"SECKILL_PREFECTURE\",\"title\":\"\\u79d2\\u6740\\u4e13\\u533a\",\"wap_url\":\"\\/pages_promotion\\/seckill\\/list\",\"parent\":\"MARKETING_LINK\"},\"label\":{\"control\":true,\"text\":\"\\u70ed\\u95e8\",\"textColor\":\"#FFFFFF\",\"bgColorStart\":\"#F83288\",\"bgColorEnd\":\"#FE3523\"},\"id\":\"lpg2grtvmxo0\"},{\"title\":\" \\u79ef\\u5206\",\"icon\":\"icondiy icon-system-point-nav\",\"imageUrl\":\"\",\"iconType\":\"icon\",\"style\":{\"fontSize\":50,\"iconBgColor\":[\"#02CC96\",\"#43EEC9\"],\"iconBgColorDeg\":90,\"iconBgImg\":\"public\\/static\\/ext\\/diyview\\/img\\/icon_bg\\/bg_06.png\",\"bgRadius\":50,\"iconColor\":[\"#FFFFFF\"],\"iconColorDeg\":0},\"link\":{\"name\":\"INTEGRAL_STORE\",\"title\":\"\\u79ef\\u5206\\u5546\\u57ce\",\"wap_url\":\"\\/pages_promotion\\/point\\/list\",\"parent\":\"MARKETING_LINK\"},\"label\":{\"control\":false,\"text\":\"\\u70ed\\u95e8\",\"textColor\":\"#FFFFFF\",\"bgColorStart\":\"#F83287\",\"bgColorEnd\":\"#FE3423\"},\"id\":\"1jfs721gome8\"},{\"title\":\"\\u4e13\\u9898\\u6d3b\\u52a8\",\"icon\":\"icondiy icon-system-topic-nav\",\"imageUrl\":\"\",\"iconType\":\"icon\",\"style\":{\"fontSize\":50,\"iconBgColor\":[\"#BE79FF\",\"#7B00FF\"],\"iconBgColorDeg\":0,\"iconBgImg\":\"public\\/static\\/ext\\/diyview\\/img\\/icon_bg\\/bg_06.png\",\"bgRadius\":50,\"iconColor\":[\"#FFFFFF\"],\"iconColorDeg\":0},\"link\":{\"name\":\"THEMATIC_ACTIVITIES_LIST\",\"title\":\"\\u4e13\\u9898\\u6d3b\\u52a8\\u5217\\u8868\",\"wap_url\":\"\\/pages_promotion\\/topics\\/list\",\"parent\":\"MARKETING_LINK\"},\"label\":{\"control\":false,\"text\":\"\\u70ed\\u95e8\",\"textColor\":\"#FFFFFF\",\"bgColorStart\":\"#F83287\",\"bgColorEnd\":\"#FE3423\"},\"id\":\"1grejh3c8fwg0\"},{\"title\":\"\\u780d\\u4ef7\",\"icon\":\"icondiy icon-system-bargain-nav\",\"imageUrl\":\"\",\"iconType\":\"icon\",\"style\":{\"fontSize\":50,\"iconBgColor\":[\"#5BBDFF\",\"#2E87FD\"],\"iconBgColorDeg\":90,\"iconBgImg\":\"public\\/static\\/ext\\/diyview\\/img\\/icon_bg\\/bg_06.png\",\"bgRadius\":50,\"iconColor\":[\"#FFFFFF\"],\"iconColorDeg\":0},\"link\":{\"name\":\"BARGAIN_PREFECTURE\",\"title\":\"\\u780d\\u4ef7\\u4e13\\u533a\",\"wap_url\":\"\\/pages_promotion\\/bargain\\/list\",\"parent\":\"MARKETING_LINK\"},\"label\":{\"control\":false,\"text\":\"\\u70ed\\u95e8\",\"textColor\":\"#FFFFFF\",\"bgColorStart\":\"#F83287\",\"bgColorEnd\":\"#FE3423\"},\"id\":\"ycpsnfbaf800\"},{\"title\":\"\\u9886\\u5238\",\"icon\":\"icondiy icon-system-get-coupon\",\"imageUrl\":\"\",\"iconType\":\"icon\",\"style\":{\"fontSize\":50,\"iconBgColor\":[\"#BE79FF\",\"#7B00FF\"],\"iconBgColorDeg\":90,\"iconBgImg\":\"public\\/static\\/ext\\/diyview\\/img\\/icon_bg\\/bg_06.png\",\"bgRadius\":50,\"iconColor\":[\"#FFFFFF\"],\"iconColorDeg\":0},\"link\":{\"name\":\"COUPON_PREFECTURE\",\"title\":\"\\u4f18\\u60e0\\u5238\\u4e13\\u533a\",\"wap_url\":\"\\/pages_tool\\/goods\\/coupon\",\"parent\":\"MARKETING_LINK\"},\"label\":{\"control\":false,\"text\":\"\\u70ed\\u95e8\",\"textColor\":\"#FFFFFF\",\"bgColorStart\":\"#F83287\",\"bgColorEnd\":\"#FE3423\"},\"id\":\"17dcs7xstz400\"},{\"title\":\"\\u6587\\u7ae0\",\"icon\":\"icondiy icon-system-article-nav\",\"imageUrl\":\"\",\"iconType\":\"icon\",\"style\":{\"fontSize\":50,\"iconBgColor\":[\"#FF8052\",\"#FF4830\"],\"iconBgColorDeg\":0,\"iconBgImg\":\"public\\/static\\/ext\\/diyview\\/img\\/icon_bg\\/bg_06.png\",\"bgRadius\":50,\"iconColor\":[\"#FFFFFF\"],\"iconColorDeg\":0},\"link\":{\"name\":\"SHOPPING_ARTICLE\",\"title\":\"\\u6587\\u7ae0\",\"wap_url\":\"\\/pages_tool\\/article\\/list\",\"parent\":\"MALL_LINK\"},\"label\":{\"control\":false,\"text\":\"\\u70ed\\u95e8\",\"textColor\":\"#FFFFFF\",\"bgColorStart\":\"#F83287\",\"bgColorEnd\":\"#FE3423\"},\"id\":\"hg8450mb0hc0\"},{\"title\":\"\\u516c\\u544a\",\"icon\":\"icondiy icon-system-notice-nav\",\"imageUrl\":\"\",\"iconType\":\"icon\",\"style\":{\"fontSize\":50,\"iconBgColor\":[\"#FFCC26\",\"#FF9F29\"],\"iconBgColorDeg\":90,\"iconBgImg\":\"public\\/static\\/ext\\/diyview\\/img\\/icon_bg\\/bg_06.png\",\"bgRadius\":50,\"iconColor\":[\"#FFFFFF\"],\"iconColorDeg\":0},\"link\":{\"name\":\"SHOPPING_NOTICE\",\"title\":\"\\u516c\\u544a\",\"wap_url\":\"\\/pages_tool\\/notice\\/list\",\"parent\":\"MALL_LINK\"},\"label\":{\"control\":false,\"text\":\"\\u70ed\\u95e8\",\"textColor\":\"#FFFFFF\",\"bgColorStart\":\"#F83287\",\"bgColorEnd\":\"#FE3423\"},\"id\":\"1cg964qu9f9c0\"},{\"title\":\"\\u5e2e\\u52a9\",\"icon\":\"icondiy icon-system-help\",\"imageUrl\":\"\",\"iconType\":\"icon\",\"style\":{\"fontSize\":50,\"iconBgColor\":[\"#02CC96\",\"#43EEC9\"],\"iconBgColorDeg\":90,\"iconBgImg\":\"public\\/static\\/ext\\/diyview\\/img\\/icon_bg\\/bg_06.png\",\"bgRadius\":50,\"iconColor\":[\"#FFFFFF\"],\"iconColorDeg\":0},\"link\":{\"name\":\"SHOPPING_HELP\",\"title\":\"\\u5e2e\\u52a9\",\"wap_url\":\"\\/pages_tool\\/help\\/list\",\"parent\":\"MALL_LINK\"},\"label\":{\"control\":false,\"text\":\"\\u70ed\\u95e8\",\"textColor\":\"#FFFFFF\",\"bgColorStart\":\"#F83287\",\"bgColorEnd\":\"#FE3423\"},\"id\":\"1v4budp7jav40\"}],\"addonName\":\"\",\"componentName\":\"GraphicNav\",\"componentTitle\":\"\\u56fe\\u6587\\u5bfc\\u822a\",\"isDelete\":0,\"pageBgColor\":\"\",\"componentBgColor\":\"#FFFFFF\",\"componentAngle\":\"round\",\"topAroundRadius\":10,\"bottomAroundRadius\":10,\"topElementAroundRadius\":0,\"bottomElementAroundRadius\":0,\"margin\":{\"top\":0,\"bottom\":12,\"both\":12}},{\"id\":\"3mkl85oxpdi0\",\"list\":[{\"link\":{\"name\":\"\"},\"imageUrl\":\"addon\\/diy_default1\\/mf_left.png\",\"imgWidth\":\"338\",\"imgHeight\":\"450\",\"previewWidth\":163.5,\"previewHeight\":\"227.68px\",\"imageMode\":\"scaleToFill\"},{\"imageUrl\":\"addon\\/diy_default1\\/mf_right1.png\",\"link\":{\"name\":\"\"},\"imgWidth\":\"354\",\"imgHeight\":\"220\",\"previewWidth\":163.5,\"previewHeight\":\"108.84px\",\"imageMode\":\"scaleToFill\"},{\"imageUrl\":\"addon\\/diy_default1\\/mf_right2.png\",\"imgWidth\":\"354\",\"imgHeight\":\"220\",\"previewWidth\":163.5,\"previewHeight\":\"108.84px\",\"link\":{\"name\":\"\"},\"imageMode\":\"scaleToFill\"}],\"mode\":\"row1-lt-of2-rt\",\"imageGap\":10,\"addonName\":\"\",\"componentName\":\"RubikCube\",\"componentTitle\":\"\\u9b54\\u65b9\",\"isDelete\":0,\"pageBgColor\":\"\",\"componentBgColor\":\"\",\"componentAngle\":\"round\",\"topAroundRadius\":10,\"bottomAroundRadius\":10,\"elementAngle\":\"round\",\"topElementAroundRadius\":0,\"bottomElementAroundRadius\":0,\"margin\":{\"top\":0,\"bottom\":12,\"both\":12}},{\"id\":\"1dbtc1ir8by8\",\"style\":\"style-16\",\"subTitle\":{\"fontSize\":14,\"text\":\"\\u8d85\\u7ea7\\u4f18\\u60e0\",\"isElementShow\":true,\"color\":\"#FFFFFF\",\"bgColor\":\"#FF9F29\",\"icon\":\"icondiy icon-system-coupon\",\"fontWeight\":\"bold\"},\"link\":{\"name\":\"COUPON_PREFECTURE\",\"title\":\"\\u4f18\\u60e0\\u5238\\u4e13\\u533a\",\"wap_url\":\"\\/pages_tool\\/goods\\/coupon\",\"parent\":\"MARKETING_LINK\"},\"fontSize\":16,\"styleName\":\"\\u98ce\\u683c16\",\"fontWeight\":\"bold\",\"more\":{\"text\":\"\",\"link\":{\"name\":\"COUPON_PREFECTURE\",\"title\":\"\\u4f18\\u60e0\\u5238\\u4e13\\u533a\",\"wap_url\":\"\\/pages_tool\\/goods\\/coupon\",\"parent\":\"MARKETING_LINK\"},\"isShow\":true,\"isElementShow\":true,\"color\":\"#999999\"},\"text\":\"\\u4f18\\u60e0\\u4e13\\u533a\",\"addonName\":\"\",\"componentName\":\"Text\",\"componentTitle\":\"\\u6807\\u9898\",\"isDelete\":0,\"pageBgColor\":\"\",\"textColor\":\"#303133\",\"componentBgColor\":\"#FFFFFF\",\"componentAngle\":\"round\",\"topAroundRadius\":10,\"bottomAroundRadius\":0,\"topElementAroundRadius\":0,\"bottomElementAroundRadius\":0,\"margin\":{\"top\":0,\"bottom\":0,\"both\":12}},{\"id\":\"534dml7c3ww0\",\"style\":\"6\",\"sources\":\"initial\",\"styleName\":\"\\u98ce\\u683c\\u516d\",\"couponIds\":[],\"count\":6,\"previewList\":[],\"nameColor\":\"#303133\",\"moneyColor\":\"#FF0000\",\"limitColor\":\"#303133\",\"btnStyle\":{\"textColor\":\"#FFFFFF\",\"bgColor\":\"#303133\",\"text\":\"\\u9886\\u53d6\",\"aroundRadius\":20,\"isBgColor\":true,\"isAroundRadius\":true},\"bgColor\":\"\",\"isName\":true,\"couponBgColor\":\"#FFFFFF\",\"couponBgUrl\":\"\",\"couponType\":\"color\",\"ifNeedBg\":true,\"addonName\":\"coupon\",\"componentName\":\"Coupon\",\"componentTitle\":\"\\u4f18\\u60e0\\u5238\",\"isDelete\":0,\"pageBgColor\":\"\",\"topAroundRadius\":0,\"bottomAroundRadius\":0,\"elementBgColor\":\"\",\"topElementAroundRadius\":0,\"bottomElementAroundRadius\":0,\"margin\":{\"top\":0,\"bottom\":0,\"both\":12}},{\"id\":\"792ss9lts9s0\",\"height\":10,\"addonName\":\"\",\"componentName\":\"HorzBlank\",\"componentTitle\":\"\\u8f85\\u52a9\\u7a7a\\u767d\",\"isDelete\":0,\"pageBgColor\":\"\",\"componentBgColor\":\"#FFFFFF\",\"componentAngle\":\"round\",\"topAroundRadius\":0,\"bottomAroundRadius\":0,\"topElementAroundRadius\":0,\"bottomElementAroundRadius\":0,\"margin\":{\"top\":0,\"bottom\":10,\"both\":12}},{\"id\":\"6vpb6jqadcc0\",\"style\":\"style-1\",\"sources\":\"initial\",\"count\":3,\"goodsId\":[],\"ornament\":{\"type\":\"default\",\"color\":\"#EDEDED\"},\"nameLineMode\":\"single\",\"template\":\"row1-of1\",\"goodsMarginType\":\"default\",\"goodsMarginNum\":10,\"btnStyle\":{\"text\":\"\\u53bb\\u79d2\\u6740\",\"textColor\":\"#FFFFFF\",\"theme\":\"default\",\"aroundRadius\":25,\"control\":true,\"support\":true,\"bgColorStart\":\"#FF7B1D\",\"bgColorEnd\":\"#FF1544\"},\"imgAroundRadius\":10,\"saleStyle\":{\"color\":\"#999CA7\",\"control\":true,\"support\":true},\"progressStyle\":{\"control\":true,\"support\":true,\"currColor\":\"#FDBE6C\",\"bgColor\":\"#FCECD7\"},\"titleStyle\":{\"backgroundImage\":\"addon\\/seckill\\/component\\/view\\/seckill\\/img\\/style_title_3_bg.png\",\"isShow\":true,\"leftStyle\":\"img\",\"leftImg\":\"addon\\/seckill\\/component\\/view\\/seckill\\/img\\/style_title_3_name.png\",\"style\":\"style-3\",\"styleName\":\"\\u98ce\\u683c3\",\"leftText\":\"\\u9650\\u65f6\\u79d2\\u6740\",\"fontSize\":16,\"fontWeight\":true,\"textColor\":\"#FFFFFF\",\"bgColorStart\":\"#FA6400\",\"bgColorEnd\":\"#FF287A\",\"more\":\"\\u66f4\\u591a\",\"moreColor\":\"#FFFFFF\",\"moreFontSize\":12,\"moreSupport\":true,\"timeBgColor\":\"\",\"timeImageUrl\":\"\",\"colonColor\":\"#FFFFFF\",\"numBgColorStart\":\"#FFFFFF\",\"numBgColorEnd\":\"#FFFFFF\",\"numTextColor\":\"#FD3B54\"},\"slideMode\":\"scroll\",\"theme\":\"default\",\"priceStyle\":{\"mainColor\":\"#FF1745\",\"mainControl\":true,\"lineColor\":\"#999CA7\",\"lineControl\":true,\"lineSupport\":true},\"goodsNameStyle\":{\"color\":\"#303133\",\"control\":true,\"fontWeight\":false},\"addonName\":\"seckill\",\"componentName\":\"Seckill\",\"componentTitle\":\"\\u79d2\\u6740\",\"isDelete\":0,\"pageBgColor\":\"\",\"componentBgColor\":\"\",\"componentAngle\":\"round\",\"topAroundRadius\":10,\"bottomAroundRadius\":10,\"elementBgColor\":\"#FFFFFF\",\"elementAngle\":\"round\",\"topElementAroundRadius\":0,\"bottomElementAroundRadius\":0,\"margin\":{\"top\":0,\"bottom\":12,\"both\":12}},{\"id\":\"2ebfwvttwntw\",\"style\":\"style-1\",\"sources\":\"initial\",\"count\":6,\"goodsId\":[],\"ornament\":{\"type\":\"default\",\"color\":\"#EDEDED\"},\"nameLineMode\":\"single\",\"template\":\"horizontal-slide\",\"goodsMarginType\":\"default\",\"goodsMarginNum\":10,\"btnStyle\":{\"text\":\"\\u53bb\\u62fc\\u56e2\",\"textColor\":\"#FFFFFF\",\"theme\":\"default\",\"aroundRadius\":25,\"control\":false,\"support\":false,\"bgColorStart\":\"#FF1544\",\"bgColorEnd\":\"#FF1544\"},\"imgAroundRadius\":10,\"saleStyle\":{\"color\":\"#FF1544\",\"control\":false,\"support\":false},\"groupStyle\":{\"color\":\"#FFFFFF\",\"control\":true,\"support\":true,\"bgColorStart\":\"#FA2379\",\"bgColorEnd\":\"#FF4F61\"},\"priceStyle\":{\"mainColor\":\"#FF1544\",\"mainControl\":true,\"lineColor\":\"#999CA7\",\"lineControl\":true,\"lineSupport\":true},\"slideMode\":\"scroll\",\"theme\":\"default\",\"goodsNameStyle\":{\"color\":\"#303133\",\"control\":true,\"fontWeight\":false},\"titleStyle\":{\"bgColorStart\":\"#9884E3\",\"bgColorEnd\":\"#68B5F0\",\"isShow\":true,\"leftStyle\":\"img\",\"leftImg\":\"addon\\/pintuan\\/component\\/view\\/pintuan\\/img\\/style_2_title.png\",\"style\":\"style-2\",\"styleName\":\"\\u98ce\\u683c2\",\"leftText\":\"\\u8d85\\u503c\\u62fc\\u56e2\",\"fontSize\":16,\"fontWeight\":true,\"textColor\":\"#888888\",\"more\":\"\\u66f4\\u591a\",\"moreColor\":\"#FFFFFF\",\"moreFontSize\":12,\"backgroundImage\":\"\"},\"addonName\":\"pintuan\",\"componentName\":\"Pintuan\",\"componentTitle\":\"\\u62fc\\u56e2\",\"isDelete\":0,\"pageBgColor\":\"\",\"componentBgColor\":\"#FFFFFF\",\"componentAngle\":\"round\",\"topAroundRadius\":0,\"bottomAroundRadius\":10,\"elementBgColor\":\"\",\"elementAngle\":\"round\",\"topElementAroundRadius\":0,\"bottomElementAroundRadius\":0,\"margin\":{\"top\":0,\"bottom\":12,\"both\":12}},{\"id\":\"9vv5m3n3bsg\",\"style\":\"style-1\",\"sources\":\"initial\",\"count\":6,\"goodsId\":[],\"ornament\":{\"type\":\"default\",\"color\":\"#EDEDED\"},\"nameLineMode\":\"single\",\"template\":\"horizontal-slide\",\"goodsMarginType\":\"default\",\"goodsMarginNum\":10,\"btnStyle\":{\"text\":\"\\u7acb\\u5373\\u62a2\\u8d2d\",\"textColor\":\"#FFFFFF\",\"theme\":\"default\",\"aroundRadius\":25,\"control\":false,\"support\":false,\"bgColorStart\":\"#FF7B1D\",\"bgColorEnd\":\"#FF1544\"},\"imgAroundRadius\":5,\"saleStyle\":{\"color\":\"#FFFFFF\",\"control\":true,\"support\":true},\"slideMode\":\"scroll\",\"theme\":\"default\",\"goodsNameStyle\":{\"color\":\"#303133\",\"control\":true,\"fontWeight\":false},\"priceStyle\":{\"mainColor\":\"#FF1745\",\"mainControl\":true,\"lineColor\":\"#999CA7\",\"lineControl\":true,\"lineSupport\":true},\"titleStyle\":{\"bgColorStart\":\"#FF209E\",\"bgColorEnd\":\"#B620E0\",\"isShow\":true,\"leftStyle\":\"img\",\"leftImg\":\"addon\\/bargain\\/component\\/view\\/bargain\\/img\\/row1_of1_style_2_name.png\",\"style\":\"style-1\",\"styleName\":\"\\u98ce\\u683c1\",\"leftText\":\"\\u75af\\u72c2\\u780d\\u4ef7\",\"fontSize\":16,\"fontWeight\":true,\"textColor\":\"#FFFFFF\",\"more\":\"\\u66f4\\u591a\",\"moreColor\":\"#FFFFFF\",\"moreFontSize\":12,\"backgroundImage\":\"addon\\/bargain\\/component\\/view\\/bargain\\/img\\/row1_of1_style_2_bg.png\"},\"addonName\":\"bargain\",\"componentName\":\"Bargain\",\"componentTitle\":\"\\u780d\\u4ef7\",\"isDelete\":0,\"pageBgColor\":\"\",\"componentBgColor\":\"#FFFFFF\",\"componentAngle\":\"round\",\"topAroundRadius\":0,\"bottomAroundRadius\":10,\"elementBgColor\":\"\",\"elementAngle\":\"round\",\"topElementAroundRadius\":0,\"bottomElementAroundRadius\":0,\"margin\":{\"top\":0,\"bottom\":12,\"both\":12}},{\"id\":\"4f7eqfqy07s0\",\"list\":[{\"link\":{\"name\":\"\"},\"imageUrl\":\"addon\\/diy_default1\\/gg.png\",\"imgWidth\":\"702\",\"imgHeight\":\"252\",\"id\":\"1z94aaav9klc0\",\"imageMode\":\"scaleToFill\"}],\"indicatorIsShow\":true,\"indicatorColor\":\"#ffffff\",\"carouselStyle\":\"circle\",\"indicatorLocation\":\"center\",\"addonName\":\"\",\"componentName\":\"ImageAds\",\"componentTitle\":\"\\u56fe\\u7247\\u5e7f\\u544a\",\"isDelete\":0,\"pageBgColor\":\"\",\"componentBgColor\":\"\",\"componentAngle\":\"round\",\"topAroundRadius\":10,\"bottomAroundRadius\":10,\"topElementAroundRadius\":0,\"bottomElementAroundRadius\":0,\"margin\":{\"top\":0,\"bottom\":12,\"both\":12}},{\"style\":\"style-2\",\"ornament\":{\"type\":\"default\",\"color\":\"#EDEDED\"},\"template\":\"row1-of2\",\"goodsMarginType\":\"default\",\"goodsMarginNum\":10,\"count\":6,\"sortWay\":\"default\",\"nameLineMode\":\"single\",\"imgAroundRadius\":0,\"slideMode\":\"scroll\",\"theme\":\"default\",\"btnStyle\":{\"fontWeight\":false,\"padding\":0,\"cartEvent\":\"detail\",\"text\":\"\\u8d2d\\u4e70\",\"textColor\":\"#FFFFFF\",\"theme\":\"default\",\"aroundRadius\":25,\"control\":true,\"support\":true,\"bgColor\":\"#FF6A00\",\"style\":\"button\",\"iconDiy\":{\"iconType\":\"icon\",\"icon\":\"\",\"style\":{\"fontSize\":\"60\",\"iconBgColor\":[],\"iconBgColorDeg\":0,\"iconBgImg\":\"\",\"bgRadius\":0,\"iconColor\":[\"#000000\"],\"iconColorDeg\":0}}},\"tag\":{\"text\":\"\\u9690\\u85cf\",\"value\":\"hidden\"},\"goodsNameStyle\":{\"color\":\"#303133\",\"control\":true,\"fontWeight\":false},\"saleStyle\":{\"color\":\"#999CA7\",\"control\":false,\"support\":true},\"priceStyle\":{\"mainColor\":\"#FF6A00\",\"mainControl\":true,\"lineColor\":\"#999CA7\",\"lineControl\":false,\"lineSupport\":true},\"list\":[{\"title\":\"\\u70ed\\u5356\",\"desc\":\"\\u70ed\\u5356\\u63a8\\u8350\",\"sources\":\"diy\",\"categoryId\":0,\"categoryName\":\"\\u8bf7\\u9009\\u62e9\",\"goodsId\":[\"172\",\"171\",\"170\",\"169\",\"168\",\"167\",\"105\",\"104\"]},{\"title\":\"\\u65b0\\u54c1\",\"desc\":\"\\u65b0\\u54c1\\u63a8\\u8350\",\"sources\":\"category\",\"categoryId\":\"63\",\"categoryName\":\"\\u7bb1\\u5305\\u978b\\u9970\",\"goodsId\":[]},{\"title\":\"\\u7cbe\\u54c1\",\"desc\":\"\\u7cbe\\u54c1\\u63a8\\u8350\",\"sources\":\"category\",\"categoryId\":\"1\",\"categoryName\":\"\\u5bb6\\u7528\\u7535\\u5668\",\"goodsId\":[]},{\"title\":\"\\u4fc3\\u9500\",\"desc\":\"\\u4fc3\\u9500\\u63a8\\u8350\",\"sources\":\"category\",\"categoryId\":\"4\",\"categoryName\":\"\\u7f8e\\u5986\\u4e2a\\u62a4\",\"goodsId\":[]}],\"id\":\"kcvtt9kl7jk\",\"addonName\":\"\",\"componentName\":\"ManyGoodsList\",\"componentTitle\":\"\\u591a\\u5546\\u54c1\\u7ec4\",\"isDelete\":0,\"pageBgColor\":\"\",\"componentBgColor\":\"\",\"componentAngle\":\"round\",\"topAroundRadius\":0,\"bottomAroundRadius\":0,\"elementBgColor\":\"#FFFFFF\",\"elementAngle\":\"round\",\"topElementAroundRadius\":0,\"bottomElementAroundRadius\":0,\"margin\":{\"top\":0,\"bottom\":0,\"both\":12},\"headStyle\":{\"titleColor\":\"#303133\"}}]}",
    "is_default": 1
  }`

	switch req.Name {
	case "DIY_VIEW_MEMBER_INDEX":
		row := make(map[string]interface{}, 0)
		json.Unmarshal([]byte(MemDat), &row)
		e.OK(row, "successful")
		return
	default:
		row := make(map[string]interface{}, 0)
		json.Unmarshal([]byte(IndexDat), &row)
		e.OK(row, "successful")
		return
		
	}
}

func (e Lib) ShopImage(c *gin.Context) {
	err := e.MakeContext(c).
		MakeOrm().
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	imageName := c.Param("path")

	pathFile := path.Join(config.ExtConfig.ImageBase, "demo", imageName)

	file, _ := ioutil.ReadFile(pathFile)
	_, _ = c.Writer.WriteString(string(file))

}
func (e Lib) Goodssku(c *gin.Context) {
	err := e.MakeContext(c).
		MakeOrm().
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	dat := `[
    {
      "goods_id": 172,
      "sku_id": 269,
      "price": "2409.00",
      "market_price": "0.00",
      "discount_price": "2409.00",
      "stock": 9986,
      "sale_num": 4,
      "goods_name": "德国施华蔻|健康护理单人套餐",
      "site_id": 1,
      "is_free_shipping": 0,
      "goods_image": "upload/1/common/images/20220929/20220929041943166443958305848.jpg",
      "is_virtual": 1,
      "recommend_way": 0,
      "unit": "",
      "promotion_type": 0,
      "label_name": "",
      "goods_spec_format": ""
    },
    {
      "goods_id": 171,
      "sku_id": 268,
      "price": "799.00",
      "market_price": "0.00",
      "discount_price": "799.00",
      "stock": 576,
      "sale_num": 6,
      "goods_name": "特权女神染发",
      "site_id": 1,
      "is_free_shipping": 0,
      "goods_image": "upload/1/common/images/20220929/20220929041404166443924407445.jpg",
      "is_virtual": 1,
      "recommend_way": 0,
      "unit": "",
      "promotion_type": 0,
      "label_name": "",
      "goods_spec_format": ""
    },
    {
      "goods_id": 170,
      "sku_id": 267,
      "price": "61.00",
      "market_price": "0.00",
      "discount_price": "61.00",
      "stock": 998,
      "sale_num": 4,
      "goods_name": "特权男士精剪",
      "site_id": 1,
      "is_free_shipping": 0,
      "goods_image": "upload/1/common/images/20220929/20220929041404166443924400744.jpg",
      "is_virtual": 1,
      "recommend_way": 0,
      "unit": "",
      "promotion_type": 0,
      "label_name": "",
      "goods_spec_format": ""
    },
    {
      "goods_id": 169,
      "sku_id": 266,
      "price": "59.00",
      "market_price": "0.00",
      "discount_price": "59.00",
      "stock": 997,
      "sale_num": 234,
      "goods_name": "针灸理疗",
      "site_id": 1,
      "is_free_shipping": 0,
      "goods_image": "upload/1/common/images/20220929/20220929040011166443841158628.jpg",
      "is_virtual": 1,
      "recommend_way": 0,
      "unit": "",
      "promotion_type": 0,
      "label_name": "",
      "goods_spec_format": ""
    },
    {
      "goods_id": 168,
      "sku_id": 265,
      "price": "396.00",
      "market_price": "0.00",
      "discount_price": "396.00",
      "stock": 996,
      "sale_num": 4,
      "goods_name": "颈椎/富贵包/肩周-轻度/重度",
      "site_id": 1,
      "is_free_shipping": 0,
      "goods_image": "upload/1/common/images/20220929/20220929035047166443784743275.jpg",
      "is_virtual": 1,
      "recommend_way": 0,
      "unit": "",
      "promotion_type": 0,
      "label_name": "",
      "goods_spec_format": ""
    },
    {
      "goods_id": 167,
      "sku_id": 264,
      "price": "198.00",
      "market_price": "0.00",
      "discount_price": "198.00",
      "stock": 990,
      "sale_num": 10,
      "goods_name": "特色头疗（轻度/重度）",
      "site_id": 1,
      "is_free_shipping": 0,
      "goods_image": "upload/1/common/images/20220929/20220929034419166443745931626.jpg",
      "is_virtual": 1,
      "recommend_way": 0,
      "unit": "",
      "promotion_type": 0,
      "label_name": "",
      "goods_spec_format": ""
    },
    {
      "goods_id": 105,
      "sku_id": 119,
      "price": "59.90",
      "market_price": "69.90",
      "discount_price": "59.90",
      "stock": 972,
      "sale_num": 17,
      "goods_name": "「小黄鸭儿童鞋2022夏季女童运动鞋男童透气网鞋断码网面宝宝学步鞋」",
      "site_id": 1,
      "is_free_shipping": 1,
      "goods_image": "https://b2c-v5-yanshi.oss-cn-hangzhou.aliyuncs.com/upload/1/common/images/20220723/20220723040600165856356023077.jpg,https://b2c-v5-yanshi.oss-cn-hangzhou.aliyuncs.com/upload/1/common/images/20220723/20220723040600165856356015327.jpg,https://b2c-v5-yanshi.oss-cn-hangzhou.aliyuncs.com/upload/1/common/images/20220723/20220723040600165856356012998.jpg",
      "is_virtual": 0,
      "recommend_way": 0,
      "unit": "",
      "promotion_type": 0,
      "label_name": "",
      "goods_spec_format": ""
    },
    {
      "goods_id": 104,
      "sku_id": 118,
      "price": "59.00",
      "market_price": "82.00",
      "discount_price": "59.00",
      "stock": 967,
      "sale_num": 13,
      "goods_name": "「小黄鸭男童鞋子儿童运动鞋2022春秋新款软底轻便休闲网面透气宝宝」",
      "site_id": 1,
      "is_free_shipping": 1,
      "goods_image": "https://b2c-v5-yanshi.oss-cn-hangzhou.aliyuncs.com/upload/1/common/images/20220723/20220723040411165856345130066.jpg,https://b2c-v5-yanshi.oss-cn-hangzhou.aliyuncs.com/upload/1/common/images/20220723/20220723040411165856345133233.jpg,https://b2c-v5-yanshi.oss-cn-hangzhou.aliyuncs.com/upload/1/common/images/20220723/20220723040411165856345134484.jpg,https://b2c-v5-yanshi.oss-cn-hangzhou.aliyuncs.com/upload/1/common/images/20220723/20220723040411165856345117963.jpg",
      "is_virtual": 0,
      "recommend_way": 0,
      "unit": "",
      "promotion_type": 0,
      "label_name": "",
      "goods_spec_format": ""
    }
  ]`

	row := make([]map[string]interface{}, 0)
	json.Unmarshal([]byte(dat), &row)
	e.OK(row, "successful")
	return
}


func (e Lib) OrderNum(c *gin.Context) {
	err := e.MakeContext(c).
		MakeOrm().
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	e.OK("", "successful")
	return
}
func (e Lib) Captcha(c *gin.Context) {
	err := e.MakeContext(c).
		MakeOrm().
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	id, b64s, err := captcha.DriverDigitFunc()
	if err != nil {
		e.Logger.Errorf("DriverDigitFunc error, %s", err.Error())
		e.Error(500, err, "验证码获取失败")
		return
	}
	e.Custom(gin.H{
		"code": 200,
		"data": map[string]string{
			"id":id,
			"img":b64s,
		},
		"msg":  "success",
	})

}
func (e Lib) CouponList(c *gin.Context) {
	err := e.MakeContext(c).
		MakeOrm().
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	dat := `[
  {
    "coupon_type_id": 1,
    "type": "reward",
    "site_id": 1,
    "coupon_name": "优惠券",
    "money": "20.00",
    "discount": "0.00",
    "max_fetch": 1,
    "at_least": "100.00",
    "end_time": 1674463181,
    "image": "",
    "validity_type": 2,
    "fixed_term": 10,
    "status": 1,
    "is_show": 1,
    "goods_type": 1,
    "discount_limit": "0.00",
    "count": 1000,
    "lead_count": 14,
    "is_remain": 1
  },
  {
    "coupon_type_id": 2,
    "type": "discount",
    "site_id": 1,
    "coupon_name": "优惠券",
    "money": "0.00",
    "discount": "8.00",
    "max_fetch": 1,
    "at_least": "150.00",
    "end_time": 1674463222,
    "image": "",
    "validity_type": 2,
    "fixed_term": 10,
    "status": 1,
    "is_show": 1,
    "goods_type": 1,
    "discount_limit": "50.00",
    "count": 1000,
    "lead_count": 51,
    "is_remain": 1
  },
  {
    "coupon_type_id": 3,
    "type": "reward",
    "site_id": 1,
    "coupon_name": "优惠券",
    "money": "1.00",
    "discount": "0.00",
    "max_fetch": 1,
    "at_least": "10.00",
    "end_time": 1674463249,
    "image": "",
    "validity_type": 2,
    "fixed_term": 10,
    "status": 1,
    "is_show": 1,
    "goods_type": 1,
    "discount_limit": "0.00",
    "count": 10000,
    "lead_count": 4,
    "is_remain": 1
  }
]`
	row := make([]map[string]interface{}, 0)
	json.Unmarshal([]byte(dat), &row)
	e.OK(row, "successful")
	return
}

func (e Lib) GetCaptchaConfig(c *gin.Context) {
	err := e.MakeContext(c).
		MakeOrm().
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	dat := `{
  "site_id": 1,
  "app_module": "shop",
  "config_key": "CAPTCHA_CONFIG",
  "value": {
    "shop_login": "1",
    "shop_reception_login": "1"
  },
  "config_desc": "验证码设置",
  "is_use": 1,
  "create_time": 1606205334,
  "modify_time": 1661913454
}`
	row := make(map[string]interface{}, 0)
	json.Unmarshal([]byte(dat), &row)
	e.OK(row, "successful")
	return
}

func (e Lib) RegisterCnf(c *gin.Context) {
	err := e.MakeContext(c).
		MakeOrm().
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	dat := `{
  "site_id": 1,
  "app_module": "shop",
  "config_key": "REGISTER_CONFIG",
  "value": {
    "login": "username,mobile",
    "register": "username,mobile",
    "pwd_len": "6",
    "pwd_complexity": "number",
    "third_party": "1",
    "bind_mobile": "1"
  },
  "config_desc": "注册规则",
  "is_use": 1,
  "create_time": 1603973110,
  "modify_time": 1669962484
}`

	row := make(map[string]interface{}, 0)
	json.Unmarshal([]byte(dat), &row)
	e.OK(row, "successful")
	return
}

func (e Lib) GoodsTree(c *gin.Context) {
	err := e.MakeContext(c).
		MakeOrm().
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	dat := `[
  {
    "category_id": 1,
    "category_name": "家用电器",
    "short_name": "",
    "pid": 0,
    "level": 1,
    "image": "upload/1/common/images/20221226/20221226025643167203780325143.png",
    "category_id_1": 1,
    "category_id_2": 0,
    "category_id_3": 0,
    "image_adv": "",
    "link_url": "",
    "is_recommend": 0,
    "icon": "",
    "child_list": [
      {
        "category_id": 13,
        "category_name": "养生壶",
        "short_name": "",
        "pid": 1,
        "level": 2,
        "image": "https://b2c-v5-yanshi.oss-cn-hangzhou.aliyuncs.com/upload/1/common/images/20220722/20220722031236165847395654153.png",
        "category_id_1": 1,
        "category_id_2": 13,
        "category_id_3": 0,
        "image_adv": "",
        "link_url": "",
        "is_recommend": 0,
        "icon": ""
      },
      {
        "category_id": 8,
        "category_name": "厨房冰箱",
        "short_name": "",
        "pid": 1,
        "level": 2,
        "image": "https://b2c-v5-yanshi.oss-cn-hangzhou.aliyuncs.com/upload/1/common/images/20220722/20220722031008165847380800645.png",
        "category_id_1": 1,
        "category_id_2": 8,
        "category_id_3": 0,
        "image_adv": "",
        "link_url": "",
        "is_recommend": 0,
        "icon": ""
      },
      {
        "category_id": 10,
        "category_name": "洗衣机",
        "short_name": "",
        "pid": 1,
        "level": 2,
        "image": "https://b2c-v5-yanshi.oss-cn-hangzhou.aliyuncs.com/upload/1/common/images/20220722/20220722031058165847385810759.png",
        "category_id_1": 1,
        "category_id_2": 10,
        "category_id_3": 0,
        "image_adv": "",
        "link_url": "",
        "is_recommend": 0,
        "icon": ""
      },
      {
        "category_id": 7,
        "category_name": "智能电视",
        "short_name": "",
        "pid": 1,
        "level": 2,
        "image": "https://b2c-v5-yanshi.oss-cn-hangzhou.aliyuncs.com/upload/1/common/images/20220722/20220722030935165847377563385.png",
        "category_id_1": 1,
        "category_id_2": 7,
        "category_id_3": 0,
        "image_adv": "",
        "link_url": "",
        "is_recommend": 0,
        "icon": ""
      },
      {
        "category_id": 11,
        "category_name": "电饭煲",
        "short_name": "",
        "pid": 1,
        "level": 2,
        "image": "https://b2c-v5-yanshi.oss-cn-hangzhou.aliyuncs.com/upload/1/common/images/20220722/20220722031148165847390818480.png",
        "category_id_1": 1,
        "category_id_2": 11,
        "category_id_3": 0,
        "image_adv": "",
        "link_url": "",
        "is_recommend": 0,
        "icon": ""
      },
      {
        "category_id": 12,
        "category_name": "破壁机",
        "short_name": "",
        "pid": 1,
        "level": 2,
        "image": "https://b2c-v5-yanshi.oss-cn-hangzhou.aliyuncs.com/upload/1/common/images/20220722/20220722031209165847392955665.png",
        "category_id_1": 1,
        "category_id_2": 12,
        "category_id_3": 0,
        "image_adv": "",
        "link_url": "",
        "is_recommend": 0,
        "icon": ""
      },
      {
        "category_id": 15,
        "category_name": "吸尘器",
        "short_name": "",
        "pid": 1,
        "level": 2,
        "image": "https://b2c-v5-yanshi.oss-cn-hangzhou.aliyuncs.com/upload/1/common/images/20220722/20220722031255165847397593363.png",
        "category_id_1": 1,
        "category_id_2": 15,
        "category_id_3": 0,
        "image_adv": "",
        "link_url": "",
        "is_recommend": 0,
        "icon": ""
      },
      {
        "category_id": 17,
        "category_name": "除螨仪",
        "short_name": "",
        "pid": 1,
        "level": 2,
        "image": "https://b2c-v5-yanshi.oss-cn-hangzhou.aliyuncs.com/upload/1/common/images/20220722/20220722031352165847403297981.png",
        "category_id_1": 1,
        "category_id_2": 17,
        "category_id_3": 0,
        "image_adv": "",
        "link_url": "",
        "is_recommend": 0,
        "icon": ""
      },
      {
        "category_id": 9,
        "category_name": "电热水器",
        "short_name": "",
        "pid": 1,
        "level": 2,
        "image": "https://b2c-v5-yanshi.oss-cn-hangzhou.aliyuncs.com/upload/1/common/images/20220722/20220722031031165847383189441.png",
        "category_id_1": 1,
        "category_id_2": 9,
        "category_id_3": 0,
        "image_adv": "",
        "link_url": "",
        "is_recommend": 0,
        "icon": ""
      }
    ]
  },
  {
    "category_id": 4,
    "category_name": "美妆个护",
    "short_name": "",
    "pid": 0,
    "level": 1,
    "image": "upload/1/common/images/20221226/20221226025325167203760521699.png",
    "category_id_1": 4,
    "category_id_2": 0,
    "category_id_3": 0,
    "image_adv": "",
    "link_url": "",
    "is_recommend": 0,
    "icon": "",
    "child_list": [
      {
        "category_id": 45,
        "category_name": "护肤套装",
        "short_name": "",
        "pid": 4,
        "level": 2,
        "image": "https://b2c-v5-yanshi.oss-cn-hangzhou.aliyuncs.com/upload/1/common/images/20220722/20220722032153165847451319741.png",
        "category_id_1": 4,
        "category_id_2": 45,
        "category_id_3": 0,
        "image_adv": "",
        "link_url": "",
        "is_recommend": 0,
        "icon": ""
      },
      {
        "category_id": 53,
        "category_name": "洁面",
        "short_name": "",
        "pid": 4,
        "level": 2,
        "image": "https://b2c-v5-yanshi.oss-cn-hangzhou.aliyuncs.com/upload/1/common/images/20220722/20220722032310165847459075429.png",
        "category_id_1": 4,
        "category_id_2": 53,
        "category_id_3": 0,
        "image_adv": "",
        "link_url": "",
        "is_recommend": 0,
        "icon": ""
      },
      {
        "category_id": 55,
        "category_name": "香水",
        "short_name": "",
        "pid": 4,
        "level": 2,
        "image": "https://b2c-v5-yanshi.oss-cn-hangzhou.aliyuncs.com/upload/1/common/images/20220722/20220722032330165847461098265.png",
        "category_id_1": 4,
        "category_id_2": 55,
        "category_id_3": 0,
        "image_adv": "",
        "link_url": "",
        "is_recommend": 0,
        "icon": ""
      },
      {
        "category_id": 56,
        "category_name": "彩妆",
        "short_name": "",
        "pid": 4,
        "level": 2,
        "image": "https://b2c-v5-yanshi.oss-cn-hangzhou.aliyuncs.com/upload/1/common/images/20220722/20220722032351165847463125166.png",
        "category_id_1": 4,
        "category_id_2": 56,
        "category_id_3": 0,
        "image_adv": "",
        "link_url": "",
        "is_recommend": 0,
        "icon": ""
      },
      {
        "category_id": 57,
        "category_name": "粉底",
        "short_name": "",
        "pid": 4,
        "level": 2,
        "image": "https://b2c-v5-yanshi.oss-cn-hangzhou.aliyuncs.com/upload/1/common/images/20220722/20220722032422165847466266320.png",
        "category_id_1": 4,
        "category_id_2": 57,
        "category_id_3": 0,
        "image_adv": "",
        "link_url": "",
        "is_recommend": 0,
        "icon": ""
      },
      {
        "category_id": 59,
        "category_name": "隔离",
        "short_name": "",
        "pid": 4,
        "level": 2,
        "image": "https://b2c-v5-yanshi.oss-cn-hangzhou.aliyuncs.com/upload/1/common/images/20220722/20220722032449165847468913675.png",
        "category_id_1": 4,
        "category_id_2": 59,
        "category_id_3": 0,
        "image_adv": "",
        "link_url": "",
        "is_recommend": 0,
        "icon": ""
      },
      {
        "category_id": 60,
        "category_name": "洗发水",
        "short_name": "",
        "pid": 4,
        "level": 2,
        "image": "https://b2c-v5-yanshi.oss-cn-hangzhou.aliyuncs.com/upload/1/common/images/20220722/20220722032508165847470875637.png",
        "category_id_1": 4,
        "category_id_2": 60,
        "category_id_3": 0,
        "image_adv": "",
        "link_url": "",
        "is_recommend": 0,
        "icon": ""
      },
      {
        "category_id": 61,
        "category_name": "沐浴露",
        "short_name": "",
        "pid": 4,
        "level": 2,
        "image": "https://b2c-v5-yanshi.oss-cn-hangzhou.aliyuncs.com/upload/1/common/images/20220722/20220722032532165847473299325.png",
        "category_id_1": 4,
        "category_id_2": 61,
        "category_id_3": 0,
        "image_adv": "",
        "link_url": "",
        "is_recommend": 0,
        "icon": ""
      },
      {
        "category_id": 62,
        "category_name": "口腔清理",
        "short_name": "",
        "pid": 4,
        "level": 2,
        "image": "https://b2c-v5-yanshi.oss-cn-hangzhou.aliyuncs.com/upload/1/common/images/20220722/20220722032556165847475635278.png",
        "category_id_1": 4,
        "category_id_2": 62,
        "category_id_3": 0,
        "image_adv": "",
        "link_url": "",
        "is_recommend": 0,
        "icon": ""
      }
    ]
  },
  {
    "category_id": 74,
    "category_name": "童装童鞋",
    "short_name": "",
    "pid": 0,
    "level": 1,
    "image": "upload/1/common/images/20221226/20221226025626167203778614525.png",
    "category_id_1": 74,
    "category_id_2": 0,
    "category_id_3": 0,
    "image_adv": "",
    "link_url": "",
    "is_recommend": 0,
    "icon": "",
    "child_list": [
      {
        "category_id": 79,
        "category_name": "裙子",
        "short_name": "",
        "pid": 74,
        "level": 2,
        "image": "https://b2c-v5-yanshi.oss-cn-hangzhou.aliyuncs.com/upload/1/common/images/20220722/20220722062316165848539679675.png",
        "category_id_1": 74,
        "category_id_2": 79,
        "category_id_3": 0,
        "image_adv": "",
        "link_url": "",
        "is_recommend": 0,
        "icon": ""
      },
      {
        "category_id": 81,
        "category_name": "汉服",
        "short_name": "",
        "pid": 74,
        "level": 2,
        "image": "https://b2c-v5-yanshi.oss-cn-hangzhou.aliyuncs.com/upload/1/common/images/20220722/20220722062341165848542118545.png",
        "category_id_1": 74,
        "category_id_2": 81,
        "category_id_3": 0,
        "image_adv": "",
        "link_url": "",
        "is_recommend": 0,
        "icon": ""
      },
      {
        "category_id": 82,
        "category_name": "演出服",
        "short_name": "",
        "pid": 74,
        "level": 2,
        "image": "https://b2c-v5-yanshi.oss-cn-hangzhou.aliyuncs.com/upload/1/common/images/20220722/20220722062404165848544426369.png",
        "category_id_1": 74,
        "category_id_2": 82,
        "category_id_3": 0,
        "image_adv": "",
        "link_url": "",
        "is_recommend": 0,
        "icon": ""
      },
      {
        "category_id": 85,
        "category_name": "T恤",
        "short_name": "",
        "pid": 74,
        "level": 2,
        "image": "https://b2c-v5-yanshi.oss-cn-hangzhou.aliyuncs.com/upload/1/common/images/20220722/20220722062432165848547228958.png",
        "category_id_1": 74,
        "category_id_2": 85,
        "category_id_3": 0,
        "image_adv": "",
        "link_url": "",
        "is_recommend": 0,
        "icon": ""
      },
      {
        "category_id": 88,
        "category_name": "外套",
        "short_name": "",
        "pid": 74,
        "level": 2,
        "image": "https://b2c-v5-yanshi.oss-cn-hangzhou.aliyuncs.com/upload/1/common/images/20220722/20220722062538165848553845056.png",
        "category_id_1": 74,
        "category_id_2": 88,
        "category_id_3": 0,
        "image_adv": "",
        "link_url": "",
        "is_recommend": 0,
        "icon": ""
      },
      {
        "category_id": 89,
        "category_name": "裤子",
        "short_name": "",
        "pid": 74,
        "level": 2,
        "image": "https://b2c-v5-yanshi.oss-cn-hangzhou.aliyuncs.com/upload/1/common/images/20220722/20220722062601165848556181668.png",
        "category_id_1": 74,
        "category_id_2": 89,
        "category_id_3": 0,
        "image_adv": "",
        "link_url": "",
        "is_recommend": 0,
        "icon": ""
      },
      {
        "category_id": 91,
        "category_name": "凉鞋",
        "short_name": "",
        "pid": 74,
        "level": 2,
        "image": "https://b2c-v5-yanshi.oss-cn-hangzhou.aliyuncs.com/upload/1/common/images/20220722/20220722062648165848560880872.png",
        "category_id_1": 74,
        "category_id_2": 91,
        "category_id_3": 0,
        "image_adv": "",
        "link_url": "",
        "is_recommend": 0,
        "icon": ""
      },
      {
        "category_id": 92,
        "category_name": "皮鞋",
        "short_name": "",
        "pid": 74,
        "level": 2,
        "image": "https://b2c-v5-yanshi.oss-cn-hangzhou.aliyuncs.com/upload/1/common/images/20220722/20220722062741165848566150791.png",
        "category_id_1": 74,
        "category_id_2": 92,
        "category_id_3": 0,
        "image_adv": "",
        "link_url": "",
        "is_recommend": 0,
        "icon": ""
      },
      {
        "category_id": 93,
        "category_name": "运动鞋",
        "short_name": "",
        "pid": 74,
        "level": 2,
        "image": "https://b2c-v5-yanshi.oss-cn-hangzhou.aliyuncs.com/upload/1/common/images/20220722/20220722062804165848568431422.png",
        "category_id_1": 74,
        "category_id_2": 93,
        "category_id_3": 0,
        "image_adv": "",
        "link_url": "",
        "is_recommend": 0,
        "icon": ""
      }
    ]
  },
  {
    "category_id": 5,
    "category_name": "品牌女装",
    "short_name": "",
    "pid": 0,
    "level": 1,
    "image": "upload/1/common/images/20221226/20221226025511167203771195782.png",
    "category_id_1": 5,
    "category_id_2": 0,
    "category_id_3": 0,
    "image_adv": "",
    "link_url": "",
    "is_recommend": 0,
    "icon": "",
    "child_list": [
      {
        "category_id": 22,
        "category_name": "时尚套装",
        "short_name": "",
        "pid": 5,
        "level": 2,
        "image": "https://b2c-v5-yanshi.oss-cn-hangzhou.aliyuncs.com/upload/1/common/images/20220722/20220722031545165847414596612.png",
        "category_id_1": 5,
        "category_id_2": 22,
        "category_id_3": 0,
        "image_adv": "",
        "link_url": "",
        "is_recommend": 0,
        "icon": ""
      },
      {
        "category_id": 24,
        "category_name": "外套",
        "short_name": "",
        "pid": 5,
        "level": 2,
        "image": "https://b2c-v5-yanshi.oss-cn-hangzhou.aliyuncs.com/upload/1/common/images/20220722/20220722031615165847417590477.png",
        "category_id_1": 5,
        "category_id_2": 24,
        "category_id_3": 0,
        "image_adv": "",
        "link_url": "",
        "is_recommend": 0,
        "icon": ""
      },
      {
        "category_id": 25,
        "category_name": "连衣裙",
        "short_name": "",
        "pid": 5,
        "level": 2,
        "image": "https://b2c-v5-yanshi.oss-cn-hangzhou.aliyuncs.com/upload/1/common/images/20220722/20220722031640165847420097081.png",
        "category_id_1": 5,
        "category_id_2": 25,
        "category_id_3": 0,
        "image_adv": "",
        "link_url": "",
        "is_recommend": 0,
        "icon": ""
      },
      {
        "category_id": 26,
        "category_name": "衬衫",
        "short_name": "",
        "pid": 5,
        "level": 2,
        "image": "https://b2c-v5-yanshi.oss-cn-hangzhou.aliyuncs.com/upload/1/common/images/20220722/20220722031711165847423106336.png",
        "category_id_1": 5,
        "category_id_2": 26,
        "category_id_3": 0,
        "image_adv": "",
        "link_url": "",
        "is_recommend": 0,
        "icon": ""
      },
      {
        "category_id": 28,
        "category_name": "毛衣",
        "short_name": "",
        "pid": 5,
        "level": 2,
        "image": "https://b2c-v5-yanshi.oss-cn-hangzhou.aliyuncs.com/upload/1/common/images/20220722/20220722031730165847425025615.png",
        "category_id_1": 5,
        "category_id_2": 28,
        "category_id_3": 0,
        "image_adv": "",
        "link_url": "",
        "is_recommend": 0,
        "icon": ""
      },
      {
        "category_id": 30,
        "category_name": "卫衣",
        "short_name": "",
        "pid": 5,
        "level": 2,
        "image": "https://b2c-v5-yanshi.oss-cn-hangzhou.aliyuncs.com/upload/1/common/images/20220722/20220722031749165847426973340.png",
        "category_id_1": 5,
        "category_id_2": 30,
        "category_id_3": 0,
        "image_adv": "",
        "link_url": "",
        "is_recommend": 0,
        "icon": ""
      },
      {
        "category_id": 31,
        "category_name": "牛仔裤",
        "short_name": "",
        "pid": 5,
        "level": 2,
        "image": "https://b2c-v5-yanshi.oss-cn-hangzhou.aliyuncs.com/upload/1/common/images/20220722/20220722031816165847429652387.png",
        "category_id_1": 5,
        "category_id_2": 31,
        "category_id_3": 0,
        "image_adv": "",
        "link_url": "",
        "is_recommend": 0,
        "icon": ""
      },
      {
        "category_id": 33,
        "category_name": "休闲裤",
        "short_name": "",
        "pid": 5,
        "level": 2,
        "image": "https://b2c-v5-yanshi.oss-cn-hangzhou.aliyuncs.com/upload/1/common/images/20220722/20220722031842165847432282617.png",
        "category_id_1": 5,
        "category_id_2": 33,
        "category_id_3": 0,
        "image_adv": "",
        "link_url": "",
        "is_recommend": 0,
        "icon": ""
      },
      {
        "category_id": 35,
        "category_name": "半身裙",
        "short_name": "",
        "pid": 5,
        "level": 2,
        "image": "https://b2c-v5-yanshi.oss-cn-hangzhou.aliyuncs.com/upload/1/common/images/20220722/20220722031912165847435244935.png",
        "category_id_1": 5,
        "category_id_2": 35,
        "category_id_3": 0,
        "image_adv": "",
        "link_url": "",
        "is_recommend": 0,
        "icon": ""
      }
    ]
  },
  {
    "category_id": 3,
    "category_name": "手机配件",
    "short_name": "",
    "pid": 0,
    "level": 1,
    "image": "upload/1/common/images/20221226/20221226025231167203755136149.png",
    "category_id_1": 3,
    "category_id_2": 0,
    "category_id_3": 0,
    "image_adv": "",
    "link_url": "",
    "is_recommend": 0,
    "icon": "",
    "child_list": [
      {
        "category_id": 32,
        "category_name": "苹果",
        "short_name": "",
        "pid": 3,
        "level": 2,
        "image": "https://b2c-v5-yanshi.oss-cn-hangzhou.aliyuncs.com/upload/1/common/images/20220722/20220722031821165847430141636.png",
        "category_id_1": 3,
        "category_id_2": 32,
        "category_id_3": 0,
        "image_adv": "",
        "link_url": "",
        "is_recommend": 0,
        "icon": ""
      },
      {
        "category_id": 34,
        "category_name": "华为",
        "short_name": "",
        "pid": 3,
        "level": 2,
        "image": "https://b2c-v5-yanshi.oss-cn-hangzhou.aliyuncs.com/upload/1/common/images/20220722/20220722031910165847435031785.png",
        "category_id_1": 3,
        "category_id_2": 34,
        "category_id_3": 0,
        "image_adv": "",
        "link_url": "",
        "is_recommend": 0,
        "icon": ""
      },
      {
        "category_id": 36,
        "category_name": "小米",
        "short_name": "",
        "pid": 3,
        "level": 2,
        "image": "https://b2c-v5-yanshi.oss-cn-hangzhou.aliyuncs.com/upload/1/common/images/20220722/20220722031935165847437511231.png",
        "category_id_1": 3,
        "category_id_2": 36,
        "category_id_3": 0,
        "image_adv": "",
        "link_url": "",
        "is_recommend": 0,
        "icon": ""
      },
      {
        "category_id": 38,
        "category_name": "手机壳",
        "short_name": "",
        "pid": 3,
        "level": 2,
        "image": "https://b2c-v5-yanshi.oss-cn-hangzhou.aliyuncs.com/upload/1/common/images/20220722/20220722032009165847440930564.png",
        "category_id_1": 3,
        "category_id_2": 38,
        "category_id_3": 0,
        "image_adv": "",
        "link_url": "",
        "is_recommend": 0,
        "icon": ""
      },
      {
        "category_id": 41,
        "category_name": "手机耳机",
        "short_name": "",
        "pid": 3,
        "level": 2,
        "image": "https://b2c-v5-yanshi.oss-cn-hangzhou.aliyuncs.com/upload/1/common/images/20220722/20220722032045165847444549238.png",
        "category_id_1": 3,
        "category_id_2": 41,
        "category_id_3": 0,
        "image_adv": "",
        "link_url": "",
        "is_recommend": 0,
        "icon": ""
      },
      {
        "category_id": 43,
        "category_name": "充电宝",
        "short_name": "",
        "pid": 3,
        "level": 2,
        "image": "https://b2c-v5-yanshi.oss-cn-hangzhou.aliyuncs.com/upload/1/common/images/20220722/20220722032132165847449241219.png",
        "category_id_1": 3,
        "category_id_2": 43,
        "category_id_3": 0,
        "image_adv": "",
        "link_url": "",
        "is_recommend": 0,
        "icon": ""
      },
      {
        "category_id": 46,
        "category_name": "手机贴膜",
        "short_name": "",
        "pid": 3,
        "level": 2,
        "image": "https://b2c-v5-yanshi.oss-cn-hangzhou.aliyuncs.com/upload/1/common/images/20220722/20220722032156165847451603610.png",
        "category_id_1": 3,
        "category_id_2": 46,
        "category_id_3": 0,
        "image_adv": "",
        "link_url": "",
        "is_recommend": 0,
        "icon": ""
      },
      {
        "category_id": 48,
        "category_name": "手机支架",
        "short_name": "",
        "pid": 3,
        "level": 2,
        "image": "https://b2c-v5-yanshi.oss-cn-hangzhou.aliyuncs.com/upload/1/common/images/20220722/20220722032220165847454080473.png",
        "category_id_1": 3,
        "category_id_2": 48,
        "category_id_3": 0,
        "image_adv": "",
        "link_url": "",
        "is_recommend": 0,
        "icon": ""
      },
      {
        "category_id": 52,
        "category_name": "数据线",
        "short_name": "",
        "pid": 3,
        "level": 2,
        "image": "https://b2c-v5-yanshi.oss-cn-hangzhou.aliyuncs.com/upload/1/common/images/20220722/20220722032306165847458661033.png",
        "category_id_1": 3,
        "category_id_2": 52,
        "category_id_3": 0,
        "image_adv": "",
        "link_url": "",
        "is_recommend": 0,
        "icon": ""
      }
    ]
  },
  {
    "category_id": 2,
    "category_name": "电子数码",
    "short_name": "",
    "pid": 0,
    "level": 1,
    "image": "upload/1/common/images/20221226/20221226024815167203729502763.png",
    "category_id_1": 2,
    "category_id_2": 0,
    "category_id_3": 0,
    "image_adv": "",
    "link_url": "",
    "is_recommend": 0,
    "icon": "",
    "child_list": [
      {
        "category_id": 20,
        "category_name": "音箱",
        "short_name": "",
        "pid": 2,
        "level": 2,
        "image": "https://b2c-v5-yanshi.oss-cn-hangzhou.aliyuncs.com/upload/1/common/images/20220722/20220722031509165847410973430.png",
        "category_id_1": 2,
        "category_id_2": 20,
        "category_id_3": 0,
        "image_adv": "",
        "link_url": "",
        "is_recommend": 0,
        "icon": ""
      },
      {
        "category_id": 18,
        "category_name": "平板电脑",
        "short_name": "",
        "pid": 2,
        "level": 2,
        "image": "https://b2c-v5-yanshi.oss-cn-hangzhou.aliyuncs.com/upload/1/common/images/20220722/20220722031359165847403995698.png",
        "category_id_1": 2,
        "category_id_2": 18,
        "category_id_3": 0,
        "image_adv": "",
        "link_url": "",
        "is_recommend": 0,
        "icon": ""
      },
      {
        "category_id": 16,
        "category_name": "笔记本电脑",
        "short_name": "",
        "pid": 2,
        "level": 2,
        "image": "https://b2c-v5-yanshi.oss-cn-hangzhou.aliyuncs.com/upload/1/common/images/20220722/20220722031333165847401363992.png",
        "category_id_1": 2,
        "category_id_2": 16,
        "category_id_3": 0,
        "image_adv": "",
        "link_url": "",
        "is_recommend": 0,
        "icon": ""
      },
      {
        "category_id": 19,
        "category_name": "耳机耳麦",
        "short_name": "",
        "pid": 2,
        "level": 2,
        "image": "https://b2c-v5-yanshi.oss-cn-hangzhou.aliyuncs.com/upload/1/common/images/20220722/20220722031441165847408113136.png",
        "category_id_1": 2,
        "category_id_2": 19,
        "category_id_3": 0,
        "image_adv": "",
        "link_url": "",
        "is_recommend": 0,
        "icon": ""
      },
      {
        "category_id": 14,
        "category_name": "数码相机",
        "short_name": "",
        "pid": 2,
        "level": 2,
        "image": "https://b2c-v5-yanshi.oss-cn-hangzhou.aliyuncs.com/upload/1/common/images/20220722/20220722031255165847397549513.png",
        "category_id_1": 2,
        "category_id_2": 14,
        "category_id_3": 0,
        "image_adv": "",
        "link_url": "",
        "is_recommend": 0,
        "icon": ""
      },
      {
        "category_id": 21,
        "category_name": "麦克风",
        "short_name": "",
        "pid": 2,
        "level": 2,
        "image": "https://b2c-v5-yanshi.oss-cn-hangzhou.aliyuncs.com/upload/1/common/images/20220722/20220722031535165847413543313.png",
        "category_id_1": 2,
        "category_id_2": 21,
        "category_id_3": 0,
        "image_adv": "",
        "link_url": "",
        "is_recommend": 0,
        "icon": ""
      },
      {
        "category_id": 23,
        "category_name": "游戏机",
        "short_name": "",
        "pid": 2,
        "level": 2,
        "image": "https://b2c-v5-yanshi.oss-cn-hangzhou.aliyuncs.com/upload/1/common/images/20220722/20220722031609165847416948993.png",
        "category_id_1": 2,
        "category_id_2": 23,
        "category_id_3": 0,
        "image_adv": "",
        "link_url": "",
        "is_recommend": 0,
        "icon": ""
      },
      {
        "category_id": 27,
        "category_name": "投影仪",
        "short_name": "",
        "pid": 2,
        "level": 2,
        "image": "https://b2c-v5-yanshi.oss-cn-hangzhou.aliyuncs.com/upload/1/common/images/20220722/20220722031712165847423229461.png",
        "category_id_1": 2,
        "category_id_2": 27,
        "category_id_3": 0,
        "image_adv": "",
        "link_url": "",
        "is_recommend": 0,
        "icon": ""
      },
      {
        "category_id": 29,
        "category_name": "无人机",
        "short_name": "",
        "pid": 2,
        "level": 2,
        "image": "https://b2c-v5-yanshi.oss-cn-hangzhou.aliyuncs.com/upload/1/common/images/20220722/20220722031738165847425816952.png",
        "category_id_1": 2,
        "category_id_2": 29,
        "category_id_3": 0,
        "image_adv": "",
        "link_url": "",
        "is_recommend": 0,
        "icon": ""
      }
    ]
  },
  {
    "category_id": 73,
    "category_name": "母婴用品",
    "short_name": "",
    "pid": 0,
    "level": 1,
    "image": "upload/1/common/images/20221226/20221226025414167203765412090.png",
    "category_id_1": 73,
    "category_id_2": 0,
    "category_id_3": 0,
    "image_adv": "",
    "link_url": "",
    "is_recommend": 0,
    "icon": "",
    "child_list": [
      {
        "category_id": 75,
        "category_name": "婴儿装",
        "short_name": "",
        "pid": 73,
        "level": 2,
        "image": "https://b2c-v5-yanshi.oss-cn-hangzhou.aliyuncs.com/upload/1/common/images/20220722/20220722062204165848532450558.png",
        "category_id_1": 73,
        "category_id_2": 75,
        "category_id_3": 0,
        "image_adv": "",
        "link_url": "",
        "is_recommend": 0,
        "icon": ""
      },
      {
        "category_id": 76,
        "category_name": "早教机",
        "short_name": "",
        "pid": 73,
        "level": 2,
        "image": "https://b2c-v5-yanshi.oss-cn-hangzhou.aliyuncs.com/upload/1/common/images/20220722/20220722062227165848534756991.png",
        "category_id_1": 73,
        "category_id_2": 76,
        "category_id_3": 0,
        "image_adv": "",
        "link_url": "",
        "is_recommend": 0,
        "icon": ""
      },
      {
        "category_id": 77,
        "category_name": "洋娃娃",
        "short_name": "",
        "pid": 73,
        "level": 2,
        "image": "https://b2c-v5-yanshi.oss-cn-hangzhou.aliyuncs.com/upload/1/common/images/20220722/20220722062249165848536976891.png",
        "category_id_1": 73,
        "category_id_2": 77,
        "category_id_3": 0,
        "image_adv": "",
        "link_url": "",
        "is_recommend": 0,
        "icon": ""
      },
      {
        "category_id": 78,
        "category_name": "毛绒玩具",
        "short_name": "",
        "pid": 73,
        "level": 2,
        "image": "https://b2c-v5-yanshi.oss-cn-hangzhou.aliyuncs.com/upload/1/common/images/20220722/20220722062311165848539152412.png",
        "category_id_1": 73,
        "category_id_2": 78,
        "category_id_3": 0,
        "image_adv": "",
        "link_url": "",
        "is_recommend": 0,
        "icon": ""
      },
      {
        "category_id": 80,
        "category_name": "洗澡用具",
        "short_name": "",
        "pid": 73,
        "level": 2,
        "image": "https://b2c-v5-yanshi.oss-cn-hangzhou.aliyuncs.com/upload/1/common/images/20220722/20220722062335165848541594819.png",
        "category_id_1": 73,
        "category_id_2": 80,
        "category_id_3": 0,
        "image_adv": "",
        "link_url": "",
        "is_recommend": 0,
        "icon": ""
      },
      {
        "category_id": 83,
        "category_name": "宝宝个护",
        "short_name": "",
        "pid": 73,
        "level": 2,
        "image": "https://b2c-v5-yanshi.oss-cn-hangzhou.aliyuncs.com/upload/1/common/images/20220722/20220722062408165848544838821.png",
        "category_id_1": 73,
        "category_id_2": 83,
        "category_id_3": 0,
        "image_adv": "",
        "link_url": "",
        "is_recommend": 0,
        "icon": ""
      },
      {
        "category_id": 84,
        "category_name": "儿童餐具",
        "short_name": "",
        "pid": 73,
        "level": 2,
        "image": "https://b2c-v5-yanshi.oss-cn-hangzhou.aliyuncs.com/upload/1/common/images/20220722/20220722062426165848546617264.png",
        "category_id_1": 73,
        "category_id_2": 84,
        "category_id_3": 0,
        "image_adv": "",
        "link_url": "",
        "is_recommend": 0,
        "icon": ""
      },
      {
        "category_id": 86,
        "category_name": "儿童水杯",
        "short_name": "",
        "pid": 73,
        "level": 2,
        "image": "https://b2c-v5-yanshi.oss-cn-hangzhou.aliyuncs.com/upload/1/common/images/20220722/20220722062447165848548748611.png",
        "category_id_1": 73,
        "category_id_2": 86,
        "category_id_3": 0,
        "image_adv": "",
        "link_url": "",
        "is_recommend": 0,
        "icon": ""
      },
      {
        "category_id": 87,
        "category_name": "婴儿推车",
        "short_name": "",
        "pid": 73,
        "level": 2,
        "image": "https://b2c-v5-yanshi.oss-cn-hangzhou.aliyuncs.com/upload/1/common/images/20220722/20220722062511165848551169932.png",
        "category_id_1": 73,
        "category_id_2": 87,
        "category_id_3": 0,
        "image_adv": "",
        "link_url": "",
        "is_recommend": 0,
        "icon": ""
      }
    ]
  },
  {
    "category_id": 6,
    "category_name": "品牌男装",
    "short_name": "",
    "pid": 0,
    "level": 1,
    "image": "upload/1/common/images/20221226/20221226025457167203769783038.jpg",
    "category_id_1": 6,
    "category_id_2": 0,
    "category_id_3": 0,
    "image_adv": "",
    "link_url": "",
    "is_recommend": 0,
    "icon": "",
    "child_list": [
      {
        "category_id": 37,
        "category_name": "运动套装",
        "short_name": "",
        "pid": 6,
        "level": 2,
        "image": "https://b2c-v5-yanshi.oss-cn-hangzhou.aliyuncs.com/upload/1/common/images/20220722/20220722031956165847439619396.png",
        "category_id_1": 6,
        "category_id_2": 37,
        "category_id_3": 0,
        "image_adv": "",
        "link_url": "",
        "is_recommend": 0,
        "icon": ""
      },
      {
        "category_id": 39,
        "category_name": "时尚套装",
        "short_name": "",
        "pid": 6,
        "level": 2,
        "image": "https://b2c-v5-yanshi.oss-cn-hangzhou.aliyuncs.com/upload/1/common/images/20220722/20220722032016165847441628347.png",
        "category_id_1": 6,
        "category_id_2": 39,
        "category_id_3": 0,
        "image_adv": "",
        "link_url": "",
        "is_recommend": 0,
        "icon": ""
      },
      {
        "category_id": 40,
        "category_name": "西服套装",
        "short_name": "",
        "pid": 6,
        "level": 2,
        "image": "https://b2c-v5-yanshi.oss-cn-hangzhou.aliyuncs.com/upload/1/common/images/20220722/20220722032044165847444423922.png",
        "category_id_1": 6,
        "category_id_2": 40,
        "category_id_3": 0,
        "image_adv": "",
        "link_url": "",
        "is_recommend": 0,
        "icon": ""
      },
      {
        "category_id": 42,
        "category_name": "外套",
        "short_name": "",
        "pid": 6,
        "level": 2,
        "image": "https://b2c-v5-yanshi.oss-cn-hangzhou.aliyuncs.com/upload/1/common/images/20220722/20220722032113165847447322989.png",
        "category_id_1": 6,
        "category_id_2": 42,
        "category_id_3": 0,
        "image_adv": "",
        "link_url": "",
        "is_recommend": 0,
        "icon": ""
      },
      {
        "category_id": 44,
        "category_name": "T恤",
        "short_name": "",
        "pid": 6,
        "level": 2,
        "image": "https://b2c-v5-yanshi.oss-cn-hangzhou.aliyuncs.com/upload/1/common/images/20220722/20220722032141165847450108458.png",
        "category_id_1": 6,
        "category_id_2": 44,
        "category_id_3": 0,
        "image_adv": "",
        "link_url": "",
        "is_recommend": 0,
        "icon": ""
      },
      {
        "category_id": 47,
        "category_name": "polo",
        "short_name": "",
        "pid": 6,
        "level": 2,
        "image": "https://b2c-v5-yanshi.oss-cn-hangzhou.aliyuncs.com/upload/1/common/images/20220722/20220722032158165847451876813.png",
        "category_id_1": 6,
        "category_id_2": 47,
        "category_id_3": 0,
        "image_adv": "",
        "link_url": "",
        "is_recommend": 0,
        "icon": ""
      },
      {
        "category_id": 49,
        "category_name": "短裤",
        "short_name": "",
        "pid": 6,
        "level": 2,
        "image": "https://b2c-v5-yanshi.oss-cn-hangzhou.aliyuncs.com/upload/1/common/images/20220722/20220722032226165847454688941.png",
        "category_id_1": 6,
        "category_id_2": 49,
        "category_id_3": 0,
        "image_adv": "",
        "link_url": "",
        "is_recommend": 0,
        "icon": ""
      },
      {
        "category_id": 51,
        "category_name": "长裤",
        "short_name": "",
        "pid": 6,
        "level": 2,
        "image": "https://b2c-v5-yanshi.oss-cn-hangzhou.aliyuncs.com/upload/1/common/images/20220722/20220722032258165847457860493.png",
        "category_id_1": 6,
        "category_id_2": 51,
        "category_id_3": 0,
        "image_adv": "",
        "link_url": "",
        "is_recommend": 0,
        "icon": ""
      },
      {
        "category_id": 54,
        "category_name": "工装裤",
        "short_name": "",
        "pid": 6,
        "level": 2,
        "image": "https://b2c-v5-yanshi.oss-cn-hangzhou.aliyuncs.com/upload/1/common/images/20220722/20220722032319165847459992946.png",
        "category_id_1": 6,
        "category_id_2": 54,
        "category_id_3": 0,
        "image_adv": "",
        "link_url": "",
        "is_recommend": 0,
        "icon": ""
      }
    ]
  },
  {
    "category_id": 63,
    "category_name": "箱包鞋饰",
    "short_name": "",
    "pid": 0,
    "level": 1,
    "image": "upload/1/common/images/20221226/20221226025301167203758152276.png",
    "category_id_1": 63,
    "category_id_2": 0,
    "category_id_3": 0,
    "image_adv": "",
    "link_url": "",
    "is_recommend": 0,
    "icon": "",
    "child_list": [
      {
        "category_id": 65,
        "category_name": "男鞋",
        "short_name": "",
        "pid": 63,
        "level": 2,
        "image": "https://b2c-v5-yanshi.oss-cn-hangzhou.aliyuncs.com/upload/1/common/images/20220722/20220722061631165848499121588.png",
        "category_id_1": 63,
        "category_id_2": 65,
        "category_id_3": 0,
        "image_adv": "",
        "link_url": "",
        "is_recommend": 0,
        "icon": ""
      },
      {
        "category_id": 64,
        "category_name": "女鞋",
        "short_name": "",
        "pid": 63,
        "level": 2,
        "image": "https://b2c-v5-yanshi.oss-cn-hangzhou.aliyuncs.com/upload/1/common/images/20220722/20220722061609165848496937580.png",
        "category_id_1": 63,
        "category_id_2": 64,
        "category_id_3": 0,
        "image_adv": "",
        "link_url": "",
        "is_recommend": 0,
        "icon": ""
      },
      {
        "category_id": 66,
        "category_name": "旅行箱",
        "short_name": "",
        "pid": 63,
        "level": 2,
        "image": "https://b2c-v5-yanshi.oss-cn-hangzhou.aliyuncs.com/upload/1/common/images/20220722/20220722061648165848500867214.png",
        "category_id_1": 63,
        "category_id_2": 66,
        "category_id_3": 0,
        "image_adv": "",
        "link_url": "",
        "is_recommend": 0,
        "icon": ""
      },
      {
        "category_id": 68,
        "category_name": " 男士包袋",
        "short_name": "",
        "pid": 63,
        "level": 2,
        "image": "https://b2c-v5-yanshi.oss-cn-hangzhou.aliyuncs.com/upload/1/common/images/20220722/20220722061731165848505132557.png",
        "category_id_1": 63,
        "category_id_2": 68,
        "category_id_3": 0,
        "image_adv": "",
        "link_url": "",
        "is_recommend": 0,
        "icon": ""
      },
      {
        "category_id": 67,
        "category_name": "女士包袋",
        "short_name": "",
        "pid": 63,
        "level": 2,
        "image": "https://b2c-v5-yanshi.oss-cn-hangzhou.aliyuncs.com/upload/1/common/images/20220722/20220722061708165848502872696.png",
        "category_id_1": 63,
        "category_id_2": 67,
        "category_id_3": 0,
        "image_adv": "",
        "link_url": "",
        "is_recommend": 0,
        "icon": ""
      },
      {
        "category_id": 69,
        "category_name": "服配",
        "short_name": "",
        "pid": 63,
        "level": 2,
        "image": "https://b2c-v5-yanshi.oss-cn-hangzhou.aliyuncs.com/upload/1/common/images/20220722/20220722061818165848509877345.png",
        "category_id_1": 63,
        "category_id_2": 69,
        "category_id_3": 0,
        "image_adv": "",
        "link_url": "",
        "is_recommend": 0,
        "icon": ""
      },
      {
        "category_id": 70,
        "category_name": "饰品",
        "short_name": "",
        "pid": 63,
        "level": 2,
        "image": "https://b2c-v5-yanshi.oss-cn-hangzhou.aliyuncs.com/upload/1/common/images/20220722/20220722061853165848513357137.png",
        "category_id_1": 63,
        "category_id_2": 70,
        "category_id_3": 0,
        "image_adv": "",
        "link_url": "",
        "is_recommend": 0,
        "icon": ""
      },
      {
        "category_id": 71,
        "category_name": "手表",
        "short_name": "",
        "pid": 63,
        "level": 2,
        "image": "https://b2c-v5-yanshi.oss-cn-hangzhou.aliyuncs.com/upload/1/common/images/20220722/20220722061913165848515366829.png",
        "category_id_1": 63,
        "category_id_2": 71,
        "category_id_3": 0,
        "image_adv": "",
        "link_url": "",
        "is_recommend": 0,
        "icon": ""
      },
      {
        "category_id": 72,
        "category_name": "眼镜",
        "short_name": "",
        "pid": 63,
        "level": 2,
        "image": "https://b2c-v5-yanshi.oss-cn-hangzhou.aliyuncs.com/upload/1/common/images/20220722/20220722061932165848517270359.png",
        "category_id_1": 63,
        "category_id_2": 72,
        "category_id_3": 0,
        "image_adv": "",
        "link_url": "",
        "is_recommend": 0,
        "icon": ""
      }
    ]
  }
]`

	row := make([]map[string]interface{}, 0)
	json.Unmarshal([]byte(dat), &row)
	e.OK(row, "successful")
	return
}
