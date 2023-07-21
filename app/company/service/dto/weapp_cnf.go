/*
*
@Author: chaoqun
* @Date: 2023/7/20 23:38
*/
package dto

type UpdateLogin struct {
	Enable bool   `json:"enable" comment:"开关"`
	T      int    `json:"t" comment:"类型"`
	Val    string `json:"val" comment:"值"`
}

type UpdateNav struct {
	NavId  int  `json:"nav_id"` //菜单ID
	CId    int  `json:"c_id"`   //大B
	Enable bool `json:"enable" comment:"开关"`
}

var NavLib = `{
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
