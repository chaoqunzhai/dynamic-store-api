/*
*
@Author: chaoqun
* @Date: 2023/6/1 00:41
*/
package apis

import (
	"fmt"
	"github.com/gin-gonic/gin"
	models2 "go-admin/app/company/models"
	"go-admin/app/shop/models"
	customUser "go-admin/common/jwt/user"
	"go-admin/global"
)

type IndexReq struct {
	Day string `json:"day" form:"day"`
}
type DetailReq struct {
	Id int `uri:"id" comment:"主键编码"`
}
type OrderShopResult struct {
	ShopId int `json:"shop_id"`
	GoodId int `json:"good_id"`
	LineId int `json:"line_id"`
}
type ReportResult struct {
	Line      string        `json:"line"`
	Driver    string        `json:"driver"`
	Id        int           `json:"id"`
	ShopName  string        `json:"shop_name"`
	ShopImage string        `json:"shop_image"`
	Goods     []reportGoods `json:"goods"`
}
type reportGoods struct {
	Name   string `json:"name"`
	Image  string `json:"image"`
	Number int `json:"number"`
}

// 获取指定日期的报表
// 按配送员区分,每个配送员
// 下订单是和商家关联的，而且商家都有一个关联的路线,所以反查即可
// 是根据配送周期
//[
//
//  {
//    "line":"丈八",
//    "driver":"张山",
//    "id":1,
//    "shop_name":"生鲜超市",
//    "shop_image":"",
//    "goods":[
//      {
//        "name":"红枣",
//        "image":"",
//        "number":"20"
//      }
//    ]
//  }
//]
func (e Orders) Index(c *gin.Context) {
	req := IndexReq{}
	err := e.MakeContext(c).
		Bind(&req).
		MakeOrm().
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	userDto, err := customUser.GetUserDto(e.Orm, c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}
	if req.Day == "" {
		e.Error(500, nil, "必须选择时间")
		return
	}
	//根据选择的日期 + 大B配置的自定义配送时间
	orderTableName := e.getTableName(userDto.CId)
	//获取指定天数的订单的商家列表
	//大B + 选择天数 + 待送 + 有用的单子 +
	//只聚合查询出，哪些客户=哪些路线  哪些商品=商品的配送
	whereSql := fmt.Sprintf("select shop_id,good_id,line_id from orders where c_id = %v and enable = %v and delivery_time = '%v' and status ='%v' GROUP BY shop_id,good_id,line_id",
		userDto.CId, true, req.Day, global.OrderStatusWait)
	orderResult := make([]OrderShopResult, 0)
	e.Orm.Table(orderTableName).Raw(whereSql).Scan(&orderResult)

	//todo:统一聚合查询,统一查询资源
	shopList := make([]int, 0)
	lineList := make([]int, 0)
	goodsList := make([]int, 0)
	for _, row := range orderResult {
		goodsList = append(goodsList, row.GoodId)
		shopList = append(shopList, row.ShopId)
		lineList = append(lineList, row.LineId)
	}
	//todo:商品信息,要把查询到对的商品放到指定的路线下
	goodsModelLists := make([]models2.Goods, 0)
	e.Orm.Model(&models2.Goods{}).Select("name,image,id").Where("enable = ? and c_id = ? and id in ?", true, userDto.CId, goodsList).Find(&goodsModelLists)
	goodsMapData:=make(map[int]models2.Goods,0)
	for _,g:=range goodsModelLists{
		goodsMapData[g.Id] =g
	}
	//todo:统一查询,做map的key
	//todo:商家的信息
	shopModelLists := make([]models.Shop, 0)
	e.Orm.Model(&models.Shop{}).Select("name,image,line_id,id").Where("enable = ? and c_id = ? and id in ?", true, userDto.CId, shopList).Find(&shopModelLists)
	shopInfoMap := make(map[int]models.Shop)
	for _, s := range shopModelLists {
		shopInfoMap[s.LineId] = s
	}
	//todo:路线信息
	lineModelLists := make([]models2.Line, 0)
	e.Orm.Model(&models2.Line{}).Select("name,driver_id,id").Where("enable = ? and c_id = ? and id in ?", true, userDto.CId, lineList).Find(&lineModelLists)

	fmt.Println("lineIds", lineModelLists)

	reportCache :=make(map[int]ReportResult,0)
	for _, line := range lineModelLists {
		if _, ok := shopInfoMap[line.Id]; !ok {
			fmt.Println("路线和商家数据严重不符合！！！")
			continue
		}
		var DriverObject models2.Driver
		e.Orm.Model(&DriverObject).Where("id = ?",line.DriverId).Limit(1).Find(&DriverObject)
		re := ReportResult{
			Line:      line.Name,
			Id:        line.Id,
			ShopName:  shopInfoMap[line.Id].Name,
			ShopImage: shopInfoMap[line.Id].Image,
			Driver:DriverObject.Name,
		}
		reportCache[line.Id] = re
	}

	var list []models2.Orders
	result := make([]ReportResult, 0)
	e.Orm.Table(orderTableName).Select("number,good_id,line_id").Where("c_id = ? and enable = ? and delivery_time = ? and status =? ",userDto.CId,true,req.Day,global.OrderStatusWait).Find(&list)

	for _,row:=range list{
		fmt.Println("商品ID",row.GoodId,"路线ID",row.LineId)
		goodsRow,ok:=goodsMapData[row.GoodId]
		if !ok{
			fmt.Println("订单中的商品不在统一数据中！")
			continue
		}

		reportRow,okReport:=reportCache[row.LineId]
		if !okReport{
			fmt.Println("订单中的路线不在统一数据中！")
			continue
		}
		report :=reportGoods{
			Name: goodsRow.Name,
			Image: goodsRow.Image,
			Number: row.Number,
		}
		reportRow.Goods = append(reportRow.Goods,report)
		result = append(result,reportRow)
	}

	e.OK(result,"successful")
	return
}
func (e Orders) Detail(c *gin.Context) {
	req := DetailReq{}
	err := e.MakeContext(c).
		Bind(&req).
		MakeOrm().
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

}