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
	Number string `json:"number"`
}

// 获取指定日期的报表
// 按配送员区分,每个配送员
// 下订单是和商家关联的，而且商家都有一个关联的路线,所以反查即可
// 是根据配送周期
// [
//
//	{
//	  "line":"丈八",
//	  "driver":"张山",
//	  "id":1,
//	  "goods":[
//	    {
//	      "name":"红枣",
//	      "image":"",
//	      "number":"20"
//	    }
//	  ]
//	}
//
// ]
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

	//通过商家列表获取到各路线
	shopList := make([]int, 0)
	for _, row := range orderResult {
		shopList = append(shopList, row.ShopId)
	}
	lineList := make([]int, 0)
	for _, row := range orderResult {
		lineList = append(lineList, row.LineId)
	}
	goodsList := make([]int, 0)
	goodMapline := make(map[int]interface{}, 0)
	for _, row := range orderResult {
		goodsList = append(goodsList, row.GoodId)
		goodMapline[row.GoodId] = row.LineId
	}
	//todo:商品信息,要把查询到对的商品放到指定的路线下

	goodsModelLists := make([]models2.Goods, 0)
	e.Orm.Model(&models.Shop{}).Select("name,image").Where("enable = ? and c_id = ? and id in ?", true, userDto.CId, goodsList).Find(&goodsModelLists)

	//todo:商家的信息
	shopModelLists := make([]models.Shop, 0)
	e.Orm.Model(&models.Shop{}).Select("name,image,line_id").Where("enable = ? and c_id = ? and id in ?", true, userDto.CId, shopList).Find(&shopModelLists)

	shopInfoMap := make(map[int]models.Shop)
	for _, r := range shopModelLists {
		shopInfoMap[r.LineId] = r
	}
	//todo:路线信息
	lineModelLists := make([]models2.Line, 0)
	e.Orm.Model(&models2.Line{}).Select("name,driver_id").Where("enable = ? and c_id = ? and id in ?", true, userDto.CId, lineList).Find(&lineModelLists)

	fmt.Println("lineIds", lineModelLists)

	result := make([]ReportResult, 0)

	for _, line := range lineModelLists {

		if _, ok := shopInfoMap[line.Id]; !ok {
			fmt.Println("路线和商家数据严重不符合！！！")
			continue
		}
		re := ReportResult{
			Line:      line.Name,
			Id:        line.Id,
			ShopName:  shopInfoMap[line.Id].Name,
			ShopImage: shopInfoMap[line.Id].Image,
			//Driver: line.DriverId,
		}

		result = append(result, re)
	}
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
