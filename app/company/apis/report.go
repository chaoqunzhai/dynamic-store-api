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
	"go-admin/common/actions"
	"go-admin/common/business"
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
	Id     int     `json:"id"`
	LineId int     `json:"line"`
	Name   string  `json:"name"`
	Image  string  `json:"image"`
	Number int     `json:"number"`
	Money  float64 `json:"money"`
}
type GoodsRow struct {
	Number int     `json:"number"`
	Money  float64 `json:"money"`
}

type CycleBaseReq struct {
	Cycle int `json:"cycle" form:"cycle"`
}
type SummaryCnfRow struct {
	GoodsName string `json:"goods_name"`
	GoodsImage string `json:"goods_image"`
	GoodsNumber int `json:"goods_number"`
	GoodsId int `json:"goods_id"`
}
// 汇总指定配送周期下的商品总数

func (e Orders)Summary(c *gin.Context)  {
	req := CycleBaseReq{}
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

	splitTableRes := business.GetTableName(userDto.CId, e.Orm)

	SummaryMap:=make(map[int]*SummaryCnfRow,0)
	var data models2.OrderCycleCnf
	e.Orm.Table(splitTableRes.OrderCycle).Select("uid,id").Scopes(
		actions.PermissionSysUser(splitTableRes.OrderCycle,userDto)).Model(
			&models2.OrderCycleCnf{}).Where("id = ?",req.Cycle).Limit(1).Find(&data)
	if data.Id == 0 {
		e.OK(business.Response{Code: -1,Msg: "暂无周期订单数据"},"")
		return
	}

	orderList:=make([]models2.Orders,0)
	//根据配送UID 统一查一下
	e.Orm.Table(splitTableRes.OrderTable).Select("order_id").Where("uid = ?", data.Uid).Find(&orderList)
	orderIds:=make([]string,0)
	for _,k:=range orderList{
		orderIds = append(orderIds,k.OrderId)
	}
	orderSpecs:=make([]models2.OrderSpecs,0)
	//查下数据 获取规格 在做一次统计
	e.Orm.Table(splitTableRes.OrderSpecs).Select("goods_name,goods_id,number,image").Where("order_id in ?",orderIds).Find(&orderSpecs)

	for _,specs:=range orderSpecs{
		cnf,ok:=SummaryMap[specs.GoodsId]
		if !ok{
			cnf = &SummaryCnfRow{
				GoodsName: specs.GoodsName,
				GoodsImage: func() string {
					if specs.Image == "" {
						return ""
					}
					return business.GetGoodsPathFirst(userDto.CId,specs.Image,global.GoodsPath)
				}(),
				GoodsNumber: specs.Number,
				GoodsId: specs.GoodsId,
			}
		}
		cnf.GoodsNumber +=specs.Number
		SummaryMap[specs.GoodsId] = cnf
		
	}
	e.OK(business.Response{Code: 1,Msg: "successful",Data: SummaryMap},"")
	return
}


// 查询指定配送周期下的路线列表

func (e Orders)Line(c *gin.Context){
	req := CycleBaseReq{}
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


	return
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
//	  "shop_name":"生鲜超市",
//	  "shop_image":"",
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
	//orderTableName := business.GetTableName(userDto.CId, e.Orm)
	splitTableRes := business.GetTableName(userDto.CId, e.Orm)
	//获取指定天数的订单的商家列表
	//大B + 选择天数 + 待送 + 有用的单子 +
	//只聚合查询出，哪些客户=哪些路线  哪些商品=商品的配送
	whereSql := fmt.Sprintf("select shop_id,good_id,line_id from orders where  enable = %v and delivery_time = '%v' and status ='%v' GROUP BY shop_id,good_id,line_id",
		true, req.Day, global.OrderStatusWaitSend)
	orderResult := make([]OrderShopResult, 0)
	e.Orm.Table(splitTableRes.OrderTable).Scopes(actions.PermissionSysUser(splitTableRes.OrderTable,userDto)).Raw(whereSql).Scan(&orderResult)

	//todo:统一聚合查询,统一查询资源
	shopList := make([]int, 0)
	lineList := make([]int, 0)
	goodsList := make([]int, 0)
	for _, row := range orderResult {
		goodsList = append(goodsList, row.GoodId)
		shopList = append(shopList, row.ShopId)
		lineList = append(lineList, row.LineId)
	}
	//todo:统一查询,做map的key

	//todo:商品信息,要把查询到对的商品放到指定的路线下
	goodsModelLists := make([]models2.Goods, 0)
	var goods models2.Goods
	e.Orm.Model(&models2.Goods{}).Scopes(actions.PermissionSysUser(goods.TableName(),userDto)).Select("name,image,id").Where("enable = ?  and id in ?", true, goodsList).Find(&goodsModelLists)
	goodsMapData := make(map[int]models2.Goods, 0)
	for _, g := range goodsModelLists {
		goodsMapData[g.Id] = g
	}
	//todo:商家的信息
	shopModelLists := make([]models.Shop, 0)
	var shopDto models.Shop
	e.Orm.Model(&models.Shop{}).Scopes(actions.PermissionSysUser(shopDto.TableName(),userDto)).Select("name,image,line_id,id").Where("enable = ?  and id in ?", true,  shopList).Find(&shopModelLists)
	shopInfoMap := make(map[int]models.Shop)
	for _, s := range shopModelLists {
		shopInfoMap[s.LineId] = s
	}
	//todo:路线信息
	lineModelLists := make([]models2.Line, 0)
	var lineDto models2.Line
	e.Orm.Model(&models2.Line{}).Scopes(actions.PermissionSysUser(lineDto.TableName(),userDto)).Select("name,driver_id,id").Where("enable = ?  and id in ?", true, lineList).Find(&lineModelLists)

	reportCache := make(map[int]ReportResult, 0)
	//todo:线路数据汇总
	for _, line := range lineModelLists {
		if _, ok := shopInfoMap[line.Id]; !ok {
			fmt.Println("路线和商家数据严重不符合！！！")
			continue
		}
		var DriverObject models2.Driver
		e.Orm.Model(&DriverObject).Where("id = ?", line.DriverId).Limit(1).Find(&DriverObject)
		re := ReportResult{
			Line:      line.Name,
			Id:        line.Id,
			ShopName:  shopInfoMap[line.Id].Name,
			ShopImage: shopInfoMap[line.Id].Image,
			Driver:    DriverObject.Name,
		}
		reportCache[line.Id] = re
	}

	var list []models2.Orders

	e.Orm.Table(splitTableRes.OrderTable).Scopes(actions.PermissionSysUser(splitTableRes.OrderTable,userDto)).Select("number,good_id,line_id,money").Where("enable = ? and delivery_time = ? and status =? ",  true, req.Day, global.OrderStatusWaitSend).Find(&list)

	//todo:商品聚合计算
	//cacheGoods := make(map[int]GoodsRow, 0)
	cacheReportGoods := make(map[int]reportGoods, 0)
	//for _, row := range list {
		//fmt.Println("商品ID", row.GoodId, "路线ID", row.LineId, "商品ID", row.GoodId)
		//goodsRow, ok := goodsMapData[row.GoodsId]
		//if !ok {
		//	fmt.Println("订单中的商品不在统一数据中！")
		//	continue
		//}
		////todo:一样的商品做一个数量和价格的叠加
		//cacheGood, validOk := cacheGoods[row.GoodsId]
		//if validOk {
		//	cacheGood.Number += row.Number
		//	cacheGood.Money += row.Money
		//	cacheGoods[row.GoodsId] = cacheGood
		//} else {
		//	cacheGoods[row.GoodsId] = GoodsRow{
		//		Number: row.Number,
		//		Money:  row.Money,
		//	}
		//}
		//newCacheGoods, _ := cacheGoods[row.GoodsId]
		//report := reportGoods{
		//	Id:     goodsRow.Id,
		//	Name:   goodsRow.Name,
		//	Image:  goodsRow.Image,
		//	Number: newCacheGoods.Number,
		//	Money:  newCacheGoods.Money,
		//	LineId: row.LineId,
		//}
		//cacheReportGoods[row.GoodsId] = report
	//}
	//fmt.Println("cache2", cacheReportGoods)

	result := make([]ReportResult, 0)
	//todo:聚合的商品放到自己的路线下
	for _, goodsData := range cacheReportGoods {
		reportRow, okReport := reportCache[goodsData.LineId]
		if !okReport {
			fmt.Println("订单中的路线不在统一数据中！")
			continue
		}
		reportRow.Goods = append(reportRow.Goods, goodsData)
		reportCache[goodsData.LineId] = reportRow
	}
	for _, row := range reportCache {

		result = append(result, row)
	}
	//fmt.Println("reportRow", result)
	e.OK(result, "successful")
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
