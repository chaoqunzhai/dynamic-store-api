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
	"go-admin/app/company/service"
	"go-admin/app/company/service/dto"
	"go-admin/app/shop/models"
	"go-admin/common/actions"
	"go-admin/common/business"
	cDto "go-admin/common/dto"
	customUser "go-admin/common/jwt/user"
	"go-admin/common/utils"
	"go-admin/global"
	"sort"
	"time"
)

type IndexReq struct {
	Day string `json:"day" form:"day"`
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

type CycleLineReq struct {
	CycleBaseReq
	LineName string `form:"line_name" json:"line_name"`
}
type SummaryCnfRow struct {
	GoodsName string `json:"goods_name"`
	GoodsImage string `json:"goods_image"`
	GoodsNumber int `json:"goods_number"`
	OrderMoney float64 `json:"order_money"` //订单最终成交价
	GoodsId int `json:"goods_id"`
	Layer int `json:"layer"`
}
type CacheMapping struct {
	LineId int `json:"line_id"`
	OrderMoney float64 `json:"order_money"` //订单最终成交价
}
type TableLineRow struct {
	LineName string  `json:"line_name"`
	LineId int `json:"line_id"`
	LineMoney float64 `json:"line_money"`
	Desc string `json:"desc"`
	DriverName string `json:"driver_name"`
	ExpirationTimeStr string `json:"expiration_time_str" gorm:"-"`
	ExpirationDay int `json:"expiration_day" gorm:"-"`
	RenewalTime    string     `json:"renewal_time" gorm:"type:datetime(3);comment:续费时间"`
	Goods []*SummaryCnfRow `json:"goods"` //路线关联的商品汇总
}
// 汇总指定配送周期下的商品总数,基于goods_id做汇总的

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

	queryStart:=time.Now()
	SummaryMap:=make(map[int]*SummaryCnfRow,0)
	var data models2.OrderCycleCnf
	e.Orm.Table(splitTableRes.OrderCycle).Select("uid,id").Scopes(
		actions.PermissionSysUser(splitTableRes.OrderCycle,userDto)).Where("id = ?",req.Cycle).Limit(1).Find(&data)
	if data.Id == 0 {
		e.OK(business.Response{Code: -1,Msg: "暂无周期订单数据"},"")
		return
	}

	orderList:=make([]models2.Orders,0)

	openApprove,_:=service.IsHasOpenApprove(userDto,e.Orm)


	orm :=e.Orm.Table(splitTableRes.OrderTable).Select("order_id")
	if openApprove{ //开启了审核,那查询状态必须是审核通过的订单

		orm = e.Orm.Table(splitTableRes.OrderTable).Where("approve_status = ?",global.OrderApproveOk)
	}
	//根据配送UID 统一查一下 订单的ID
	orm.Where("uid = ? and c_id = ? and status in ?", data.Uid,userDto.CId,global.OrderEffEct()).Find(&orderList)
	orderIds:=make([]string,0)
	for _,k:=range orderList{
		orderIds = append(orderIds,k.OrderId)
	}
	orderSpecs:=make([]models2.OrderSpecs,0)
	//查下数据 获取规格 在做一次统计
	e.Orm.Table(splitTableRes.OrderSpecs).Select("id,goods_name,goods_id,number,image").Where("order_id in ?",orderIds).Find(&orderSpecs)

	//resultTable:=make([]interface{},0)
	goodsId:=make([]int,0)
	for _,specs:=range orderSpecs{
		goodsId = append(goodsId,specs.GoodsId)
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
				GoodsId: specs.GoodsId,
			}
		}
		cnf.GoodsNumber +=specs.Number

		SummaryMap[specs.GoodsId] = cnf
	}
	goodsId = utils.RemoveRepeatInt(goodsId)
	var goodsList []models2.Goods
	e.Orm.Model(&models2.Goods{}).Select("id,layer").Where("id in ?",goodsId).Order(global.OrderLayerKey).Find(&goodsList)
	sortData:=make([]*SummaryCnfRow,0)
	for _,row:=range goodsList{
		GetData :=SummaryMap[row.Id]
		GetData.Layer = row.Layer
		sortData = append(sortData,GetData)
	}
	sort.Slice(sortData, func(i, j int) bool {
		return sortData[i].Layer > sortData[j].Layer
	})
	e.OK(business.Response{Code: 1,Msg: "操作成功",Data: sortData,Extend: fmt.Sprintf("query run time %v",time.Since(queryStart))},"")
	return
}


// 查询指定配送周期下的路线列表,基于line_id做汇总

func (e Orders)Line(c *gin.Context){
	req := CycleLineReq{}
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

	queryStart:=time.Now()


	var data models2.OrderCycleCnf
	e.Orm.Table(splitTableRes.OrderCycle).Select("uid,id").Scopes(
		actions.PermissionSysUser(splitTableRes.OrderCycle,userDto)).Where("id = ? ",req.Cycle).Limit(1).Find(&data)
	if data.Id == 0 {
		e.OK(business.Response{Code: -1,Msg: "暂无周期订单数据"},"")
		return
	}

	orderList:=make([]models2.Orders,0)
	//根据配送UID 统一查一下 线路ID和订单ID
	//只查询待配送,配送中,售后中,或者审批通过, 订单完结

	openApprove,_:=service.IsHasOpenApprove(userDto,e.Orm)


	orm :=e.Orm.Table(splitTableRes.OrderTable).Select("line_id,order_id,id")
	if openApprove{ //开启了审核,那查询状态必须是审核通过的订单

		orm = e.Orm.Table(splitTableRes.OrderTable).Where("approve_status = ?",global.OrderApproveOk)
	}

	orderOrm :=orm.Where("uid = ? and c_id = ? and status in ? and order_money > 0", data.Uid,userDto.CId,global.OrderEffEct())

	//进行路线名称查询
	//首先查出的路线 必须在订单中
	if req.LineName != ""{
		searchLine :=make([]models2.Line,0)
		likeVal:=fmt.Sprintf("%%%v%%",req.LineName)
		e.Orm.Model(&models2.Line{}).Select("id").Where("c_id = ? and `name` like ?",userDto.CId,likeVal).Find(&searchLine)
		searchIds:=make([]int,0)
		for _,k:=range searchLine{
			searchIds = append(searchIds,k.Id)
		}
		//在基于路线过滤一次
		orderOrm.Where("line_id in ?",searchIds).Find(&orderList)

	}else {
		orderOrm.Find(&orderList)
	}

	queryOrderTime:=time.Since(queryStart)

	lineIds:=make([]int,0)
	orderIds:=make([]string,0)
	//订单和线路做一个map,方便把订单放到线路里面
	orderLinMapping:=make(map[string]CacheMapping,0)
	for _,k:=range orderList{

		//统一查询路线
		lineIds = append(lineIds,k.LineId)
		//统一查询订单
		orderIds = append(orderIds,k.OrderId)
		//路线和订单做一个映射
		orderLinMapping[k.OrderId] = CacheMapping{
			LineId: k.LineId,
			//OrderMoney: k.OrderMoney, 不进行计算
		}
	}

	//查询线路信息
	lineIds = utils.RemoveRepeatInt(lineIds) //去重
	lineList :=make([]models2.Line,0)
	e.Orm.Model(&models2.Line{}).Where("c_id = ? and id in ?",userDto.CId,lineIds).Order(global.OrderLayerKey).Find(&lineList)

	//存放最终返回的数据
	ResultMap:=make(map[int]*TableLineRow,0)
	for _,l:=range lineList{
		var DriverObj models2.Driver
		e.Orm.Model(&models2.Driver{}).Select("id,name,phone").Where("c_id = ? and id = ?",userDto.CId,l.DriverId).Limit(1).Find(&DriverObj)
		LineRow := &TableLineRow{
			LineName: l.Name,
			Desc: l.Desc,
			RenewalTime:l.RenewalTime.Format("2006-01-02 15:04:05"),
			Goods: make([]*SummaryCnfRow,0),
		}
		if DriverObj.Id > 0 {
			LineRow.DriverName = fmt.Sprintf("%v/%v",DriverObj.Name,DriverObj.Phone)
		}

		if !l.ExpirationTime.Time.IsZero() {
			LineRow.ExpirationDay = int(l.ExpirationTime.Sub(time.Now()).Hours() / 24)
			LineRow.ExpirationTimeStr = l.ExpirationTime.Format("2006-01-02 15:04:05")
		}else {
			LineRow.ExpirationTimeStr = "无期限"
		}
		ResultMap[l.Id] = LineRow
	}
	queryStart2:=time.Now()
	orderSpecs:=make([]models2.OrderSpecs,0)
	//查下数据 获取规格 在做一次统计
	e.Orm.Table(splitTableRes.OrderSpecs).Select("goods_name,goods_id,number,image,order_id").Where("order_id in ?",orderIds).Find(&orderSpecs)

	queryOrderSpecsTime:=time.Since(queryStart2)

	var goodsIdLists []int
	//订单商品放到路线中
	for _,specs:=range orderSpecs{

		//通过订单ID 获取到路线ID
		getLineMapInfo,ok:=orderLinMapping[specs.OrderId]
		if !ok{continue}

		//通过路线ID 获取这个路线的大数据
		lineTableRow,ok:=ResultMap[getLineMapInfo.LineId]
		if !ok{continue}
		//规格的商品信息
		cnf := &SummaryCnfRow{
			GoodsName: specs.GoodsName,
			GoodsImage: func() string {
				if specs.Image == "" {
					return ""
				}
				return business.GetGoodsPathFirst(userDto.CId,specs.Image,global.GoodsPath)
			}(),
			GoodsNumber: specs.Number,
			GoodsId: specs.GoodsId,
			OrderMoney: getLineMapInfo.OrderMoney,
		}
		//直接把商品配置都放路线里面,汇总在统计后 在统一去count
		lineTableRow.Goods = append(lineTableRow.Goods,cnf)

		ResultMap[getLineMapInfo.LineId] = lineTableRow
		goodsIdLists =append(goodsIdLists,specs.GoodsId)
	}
	goodsIdLists = utils.RemoveRepeatInt(goodsIdLists)
	//排序
	var goodsList []models2.Goods
	e.Orm.Model(&models2.Goods{}).Select("id,layer").Where("id in ?",goodsIdLists).Order(global.OrderLayerKey).Find(&goodsList)
	goodsLayerMap:=make(map[int]int,0)
	for _,row:=range goodsList{
		goodsLayerMap[row.Id] = row.Layer
	}
	//对每个路线下的商品数据 在统计count一次
	resultTable:=make([]*TableLineRow,0)

	for lineId:=range ResultMap{
		lineRow:=ResultMap[lineId]
		//对每个路线汇总
		SummaryGoodsMap:=make(map[int]*SummaryCnfRow,0)
		for _,v:=range lineRow.Goods{

			cnf,ok:=SummaryGoodsMap[v.GoodsId]
			if !ok{
				cnf = v
			}else {
				cnf.GoodsNumber +=v.GoodsNumber
			}
			cnf.Layer = goodsLayerMap[v.GoodsId]

			SummaryGoodsMap[v.GoodsId] = cnf

			//lineRow.LineMoney =utils.RoundDecimalFlot64(cnf.OrderMoney) + utils.RoundDecimalFlot64(lineRow.LineMoney)
		}

		newTable:=make([]*SummaryCnfRow,0)
		for goodsId:=range SummaryGoodsMap{
			newTable = append(newTable,SummaryGoodsMap[goodsId])
		}
		//数据还原
		lineRow.Goods = newTable
		lineRow.LineId = lineId

		resultTable = append(resultTable,lineRow)
	}

	for _,linRow:=range resultTable{
		sort.Slice(linRow.Goods, func(i, j int) bool {
			return linRow.Goods[i].Layer > linRow.Goods[j].Layer
		})
	}



	ExtendMap:=map[string]interface{}{
		"run_time":fmt.Sprintf("%v",time.Since(queryStart)),
		"queryOrderTime":fmt.Sprintf("%v",queryOrderTime),
		"queryOrderSpecsTime":fmt.Sprintf("%v",queryOrderSpecsTime),
	}
	e.OK(business.Response{Code: 1,Msg: "操作成功",Data: resultTable,Extend:ExtendMap},"")
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

	e.OK(result, "操作成功")
	return
}
func (e Orders) Detail(c *gin.Context) {
	req := dto.LineDetailReq{}
	err := e.MakeContext(c).
		Bind(&req).
		MakeOrm().
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	result:=make([]map[string]interface{},0)
	LineId:=c.Param("line_id")

	//fmt.Println("uid!",uid,req.Cycle)
	userDto, err := customUser.GetUserDto(e.Orm, c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}

	splitTableRes := business.GetTableName(userDto.CId, e.Orm)
	var data models2.OrderCycleCnf
	e.Orm.Table(splitTableRes.OrderCycle).Select("uid,id").Scopes(
		actions.PermissionSysUser(splitTableRes.OrderCycle,userDto)).Where("id = ? ",req.Cycle).Limit(1).Find(&data)
	if data.Id == 0 {
		e.OK(business.Response{Code: -1,Msg: "暂无周期订单数据"},"")
		return
	}

	orderList:=make([]models2.Orders,0)

	openApprove,_:=service.IsHasOpenApprove(userDto,e.Orm)


	orm :=e.Orm.Table(splitTableRes.OrderTable).Select("order_id,shop_id,order_money,number")
	if openApprove{ //开启了审核,那查询状态必须是审核通过的订单

		orm = e.Orm.Table(splitTableRes.OrderTable).Where("approve_status = ?",global.OrderApproveOk)
	}
	//根据配送UID 统一查一下 订单的ID
	orm.Where("line_id = ? and uid = ? and c_id = ? and status in ?",LineId,data.Uid,userDto.CId,global.OrderEffEct()).Find(&orderList)

	shopIds:=make([]int,0)
	shopGoodsMap:=make(map[int]dto.DetailCount,0)
	for _,k:=range orderList{
		shopIds = append(shopIds,k.ShopId)
		detailRow,ok:=shopGoodsMap[k.ShopId]
		if !ok{
			detailRow = dto.DetailCount{
				Count: k.Number,
				Money: k.OrderMoney,
			}
		}else {
			detailRow.Count +=k.Number
			detailRow.Money +=k.OrderMoney
		}
		shopGoodsMap[k.ShopId] = detailRow
	}
	shopIds = utils.RemoveRepeatInt(shopIds)

	var shopList []models.Shop
	var count int64
	e.Orm.Model(&models.Shop{}).Scopes(
		cDto.Paginate(req.GetPageSize(), req.GetPageIndex()),
		).Order(global.OrderLayerKey).Select("name,id").Where("id in ?",
			shopIds).Find(&shopList).Limit(-1).Offset(-1).
		Count(&count)

	for _,row:=range shopList{
		detailItem :=shopGoodsMap[row.Id]
		item:=map[string]interface{}{
			"name":row.Name,
			"count":detailItem.Count,
			"shop_id":row.Id,
			"all_money":fmt.Sprintf("%v%v",global.SymBol,detailItem.Money),
		}
		result = append(result,item)
	}

	mapData :=map[string]interface{}{
		"list":result,
		"count":count,
		"pageIndex": req.GetPageIndex(),
		"pageSize":req.GetPageSize(),
	}
	e.OK(mapData,"查询成功")
	return
}

func (e Orders)DetailShopGoods(c *gin.Context)  {
	req := dto.ShopLineDetailReq{}
	err := e.MakeContext(c).
		Bind(&req).
		MakeOrm().
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	LineId:=c.Param("line_id")
	userDto, err := customUser.GetUserDto(e.Orm, c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}
	//fmt.Println("查询周期: ",req.Cycle,"小B: ",req.ShopId,"LineId: ",LineId)
	splitTableRes := business.GetTableName(userDto.CId, e.Orm)
	//
	var data models2.OrderCycleCnf
	e.Orm.Table(splitTableRes.OrderCycle).Select("uid,id").Scopes(
		actions.PermissionSysUser(splitTableRes.OrderCycle,userDto)).Where("id = ? ",req.Cycle).Limit(1).Find(&data)
	if data.Id == 0 {
		e.OK(business.Response{Code: -1,Msg: "暂无周期订单数据"},"")
		return
	}
	var count int64
	orderList:=make([]models2.Orders,0)

	openApprove,_:=service.IsHasOpenApprove(userDto,e.Orm)

	orm := e.Orm.Table(splitTableRes.OrderTable).Select("order_id").Scopes(
		cDto.Paginate(req.GetPageSize(), req.GetPageIndex()),
	)
	if openApprove{ //开启了审核,那查询状态必须是审核通过的订单

		orm = e.Orm.Table(splitTableRes.OrderTable).Where("approve_status = ?",global.OrderApproveOk)
	}
	//根据配送UID 统一查一下 订单的ID
	orm.Where("shop_id = ? and line_id = ? and uid = ? and c_id = ? and status in ?",
		req.ShopId,LineId,data.Uid,userDto.CId,
		global.OrderEffEct()).Find(&orderList).Limit(-1).Offset(-1).Count(&count)

	orderIds:=make([]string,0)
	for _,row:=range orderList{
		orderIds = append(orderIds,row.OrderId)
	}
	orderIds = utils.RemoveRepeatStr(orderIds)

	result:=make([]interface{},0)


	var orderSpecs []models2.OrderSpecs
	e.Orm.Table(splitTableRes.OrderSpecs).Where("order_id in ?", orderIds).Order("id desc").Find(&orderSpecs)


	mergeMap:=make(map[string]dto.DetailGoodsRow,0)
	for _,row:=range orderSpecs{
		//一样的需要合并
		key :=fmt.Sprintf("%v_%v",row.GoodsId,row.SpecId)
		goodsRow:=dto.DetailGoodsRow{
			Id: row.Id,
			Name: row.SpecsName,
			GoodsName: row.GoodsName,
			Number: row.Number,
			CreatedAt:  row.CreatedAt.Format("2006-01-02 15:04:05"),
			Unit: row.Unit,
			Money: global.SymBol + utils.StringDecimal(row.Money),
			AllMoney: utils.RoundDecimalFlot64(row.Money  * float64(row.Number)),
			Image: func() string {
				if row.Image == "" {
					return ""
				}
				return business.GetGoodsPathFirst(userDto.CId,row.Image,global.GoodsPath)
			}(),
		}
		getData,ok:=mergeMap[key]

		if ok{ //如果有 那就叠加
			getData.Number +=goodsRow.Number
			getData.AllMoney +=goodsRow.AllMoney
			mergeMap[key] = getData
		}else { //如果没有赋值
			mergeMap[key] = goodsRow
		}
	}
	for _,row:=range mergeMap{
		row.AllMoneyValue = global.SymBol + utils.StringDecimal(row.AllMoney)
		result = append(result,row)
	}
	mapData :=map[string]interface{}{
		"list":result,
		"count":count,
		"pageIndex": req.GetPageIndex(),
		"pageSize":req.GetPageSize(),
	}
	e.OK(mapData,"查询成功")
	return
}

func (e Orders)LineGoodsDetail(c *gin.Context)  {
	req := dto.LineGoodsDetail{}
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
	//fmt.Println("查询周期: ",req.Cycle,"小B: ",req.ShopId,"LineId: ",LineId)
	splitTableRes := business.GetTableName(userDto.CId, e.Orm)
	//
	var data models2.OrderCycleCnf
	e.Orm.Table(splitTableRes.OrderCycle).Select("uid,id").Scopes(
		actions.PermissionSysUser(splitTableRes.OrderCycle,userDto)).Where("id = ? ",req.CycleId).Limit(1).Find(&data)
	if data.Id == 0 {
		e.OK(business.Response{Code: -1,Msg: "暂无周期订单数据"},"")
		return
	}
	orderList:=make([]models2.Orders,0)

	openApprove,_:=service.IsHasOpenApprove(userDto,e.Orm)

	orm := e.Orm.Table(splitTableRes.OrderTable).Select("order_id")
	if openApprove{ //开启了审核,那查询状态必须是审核通过的订单

		orm = e.Orm.Table(splitTableRes.OrderTable).Where("approve_status = ?",global.OrderApproveOk)
	}
	//根据配送UID 统一查一下 订单的ID
	orm.Where("line_id = ? and uid = ? and c_id = ? and status in ?",
		req.LineId,data.Uid,userDto.CId,
		global.OrderEffEct()).Find(&orderList)

	orderIds:=make([]string,0)
	for _,row:=range orderList{
		orderIds = append(orderIds,row.OrderId)
	}
	orderIds = utils.RemoveRepeatStr(orderIds)

	result:=make([]interface{},0)


	//必须查商品,因为是从商品里面查出来的
	var count int64
	var orderSpecs []models2.OrderSpecs
	e.Orm.Table(splitTableRes.OrderSpecs).Where("goods_id = ? and order_id in ?",
		req.GoodsId,orderIds).Scopes(
		cDto.Paginate(req.GetPageSize(), req.GetPageIndex()),
	).Order("id desc").Find(&orderSpecs).Limit(-1).Offset(-1).Count(&count)

	//查询商品的照片
	for _,row:=range orderSpecs{
		//一样的需要合并
		goodsRow:=dto.DetailGoodsRow{
			Id: row.Id,
			Name: row.SpecsName,
			GoodsName: row.GoodsName,
			Number: row.Number,
			CreatedAt:  row.CreatedAt.Format("2006-01-02 15:04:05"),
			Unit: row.Unit,
			Money: global.SymBol+ utils.StringDecimal(row.Money),
			Image: func() string {
				if row.Image == "" {
					return ""
				}
				return business.GetGoodsPathFirst(userDto.CId,row.Image,global.GoodsPath)
			}(),
			AllMoney: utils.RoundDecimalFlot64(row.Money  * float64(row.Number)),
		}
		goodsRow.AllMoneyValue = global.SymBol + utils.StringDecimal(row.AllMoney)
		result = append(result,goodsRow)
	}
	mapData :=map[string]interface{}{
		"list":result,
		"count":count,
		"pageIndex": req.GetPageIndex(),
		"pageSize":req.GetPageSize(),
	}
	e.OK(mapData,"查询成功")
	return


}