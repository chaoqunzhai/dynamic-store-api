/**
@Author: chaoqun
* @Date: 2024/8/4 17:16
*/
package apis

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"go-admin/app/company/models"
	"go-admin/app/company/service"
	models2 "go-admin/cmd/migrate/migration/models"
	"go-admin/common/business"
	customUser "go-admin/common/jwt/user"
	"go-admin/common/utils"
	"go-admin/global"
	"sort"
	"strings"
	"time"
)
type DataAnalysis struct {
	api.Api
}
//商品汇总
type GoodsAnalysisTableRow struct {
	OrderId string `json:"order_id"`
	GoodsName string `json:"goods_name"`
	SpecsName string `json:"specs_name"`
	AllNumber int64 `json:"all_number"`
	AllMoneyValue float64 `json:"-"`
	AllMoney string `json:"all_money"`
	Gross string `json:"gross"`  //毛利=[(实际销售金-实际退货金额)-(销售成本-退货成本)]
	GrossProfit string `json:"gross_profit"` //毛利率=[(实际销售金-实际退货金额)-(销售成本-退货成本)]/实际销售金额
	Cost string `json:"cost"` //销售成本
	SalesGross  string `json:"sales_gross"` //销售毛利
	SalesGrossProfit  string `json:"sales_gross_profit"` //销售毛利率
	Unit string `json:"unit"`
	RefundCount int64 `json:"refund_number"`
	RefundMoneyValue float64 `json:"-"`
	RefundMoney string `json:"refund_money"`
	RefundCost string `json:"refund_cost"`
	Income string `json:"income"`
}
//商品分类汇总
type GoodsClsAnalysisTableRow struct {
	GoodsId int `json:"goods_id"`
	ClsName string `json:"cls_name"`
	ClsId int `json:"cls_id"`
	AllNumber int64 `json:"all_number"`
	AllMoneyValue float64 `json:"-"`
	AllMoney string `json:"all_money"`
	RefundCount int64 `json:"refund_number"`
	RefundMoneyValue float64 `json:"-"`
	RefundMoney string `json:"refund_money"`
	Income string `json:"income"`

}

//商品品牌汇总
type GoodsBrandAnalysisTableRow struct {
	GoodsId int `json:"goods_id"`
	BrandName string `json:"brand_name"`
	BrandId int `json:"brand_id"`
	AllNumber int64 `json:"all_number"`
	AllMoney float64 `json:"all_money"`
	RefundCount int64 `json:"refund_number"`
	RefundMoneyValue float64 `json:"-"`
	RefundMoney string `json:"refund_money"`
	Income string `json:"income"`

}
type ByResultAnalysisMoney []*GoodsAnalysisTableRow

func (a ByResultAnalysisMoney) Len() int           { return len(a) }
func (a ByResultAnalysisMoney) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByResultAnalysisMoney) Less(i, j int) bool { return a[i].AllMoney > a[j].AllMoney } // 注意这里是从大到小


type ByResultClsAnalysisMoney []*GoodsClsAnalysisTableRow

func (a ByResultClsAnalysisMoney) Len() int           { return len(a) }
func (a ByResultClsAnalysisMoney) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByResultClsAnalysisMoney) Less(i, j int) bool { return a[i].AllMoney > a[j].AllMoney } // 注意这里是从大到小



type ByResultBrandAnalysisMoney []*GoodsBrandAnalysisTableRow

func (a ByResultBrandAnalysisMoney) Len() int           { return len(a) }
func (a ByResultBrandAnalysisMoney) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByResultBrandAnalysisMoney) Less(i, j int) bool { return a[i].AllMoney > a[j].AllMoney } // 注意这里是从大到小



type AnalysisQuery struct {
	OrderType    int    `form:"orderType" `
	ClassID    int    `form:"classId" `
	BrandId int `form:"brandId"`
	BandId int `form:"bandId"`
	CustomerUser    int    `form:"customerUser" `
	SpecsId    int    `form:"specsId" `
	BeginTime  string `form:"beginTime"`
	EndTime string `form:"endTime"`
}

type RefundRow struct {
	AllNumber int64 `json:"all_number"`
	AllMoney float64 `json:"all_money"`
}

// 列表
//商品名称 | 订货数量 | 订货金额 | 优惠金额 | 实际销售金额 | 退货数量 | 退货金额 | 净销售收入

func (e DataAnalysis) GoodsList(c *gin.Context) {
	req:=AnalysisQuery{}
	s := service.Goods{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
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
	nowTime:=time.Now()

	//查询的总订单数
	var queryAllCount int64
	//查询的总金额
	var queryAllMoney float64

	//退货总数量
	var refundAllCount int64
	var refundAllMoney float64
	//获取客户的分表
	splitTableRes := business.GetTableName(userDto.CId, e.Orm)

	//默认查询15天的时间
	//限制最多查询60天的数据
	//开始时间
	var startTime string
	var endTime string
	if req.BeginTime!=""{
		startTime = req.BeginTime
	}else {
		startTime = nowTime.AddDate(0,0,-15).Format(time.DateTime)

	}
	if req.EndTime != "" {
		endTime = req.EndTime
	}else {
		endTime =nowTime.Format(time.DateTime)

	}


	whereRangeTime:=fmt.Sprintf("c_id = %v and created_at >= '%v' and created_at <= '%v' ",
		userDto.CId,startTime,endTime)

	//不同的订单完结状态不一样
	orderTypes:=make([]int,0)
	switch req.OrderType {
	case global.ExpressSelf:
		orderTypes = append(orderTypes,global.ExpressSelf)
	case global.ExpressSameCity:
		orderTypes = append(orderTypes,global.ExpressSameCity)
	case global.ExpressEms:
		orderTypes = append(orderTypes,global.ExpressEms)
	default:
		orderTypes = append(orderTypes,global.ExpressEms)
		orderTypes = append(orderTypes,global.ExpressSelf)
		orderTypes = append(orderTypes,global.ExpressSameCity)
	}
	refundMapCache:=make(map[string]RefundRow,0)
	//考虑 优惠金额
	//考虑退货的商品
	var queryOrderId []string
	var timeRangeOrder []models.Orders
	orderOrm :=e.Orm.Table(splitTableRes.OrderTable).Select("order_id").Where(whereRangeTime).Where("delivery_type in ?",
		orderTypes).Order(global.OrderTimeKey)

	if req.CustomerUser > 0{
		orderOrm = orderOrm.Where("shop_id = ?",req.CustomerUser)
	}

	orderOrm.Find(&timeRangeOrder)
	for _,oo:=range timeRangeOrder{
		queryOrderId = append(queryOrderId,oo.OrderId)
	}
	if req.ClassID > 0{//查询商品分类 -> 获取到订单ID
		queryOrderId = []string{} //重置查询的ID
		var bindGoodsId []int
		e.Orm.Raw(fmt.Sprintf("select goods_id from goods_mark_class where class_id in (%v)",req.ClassID)).Scan(&bindGoodsId)

		var bindGoodsSpecs []models.OrderSpecs
		e.Orm.Table(splitTableRes.OrderSpecs).Select("order_id").Where(whereRangeTime).Where("goods_id in ?",bindGoodsId).Find(&bindGoodsSpecs)

		if len(bindGoodsSpecs) == 0 {
			//查询分类没有那就返回
			queryResult :=map[string]interface{}{
				"calculationCount":map[string]interface{}{
					"queryAllCount":0,
					"queryAllMoney":"0.0",
					"refundAllCount":0,
					"refundAllMoney":"0.0",
				},
				"list":make([]string,0),
				"total":0,
			}
			e.OK(queryResult,"")
			return
		}
		for _,b:=range bindGoodsSpecs{
			queryOrderId = append(queryOrderId,b.OrderId)
		}

	}
	if req.SpecsId > 0 {//规格查询 -> 获取到订单ID
		queryOrderId = []string{} //重置查询的ID
		var specsBindOrder []models.OrderSpecs
		e.Orm.Table(splitTableRes.OrderSpecs).Where(whereRangeTime).Where("spec_id = ?",
			req.SpecsId).Find(&specsBindOrder)
		for _,b:=range specsBindOrder{
			queryOrderId = append(queryOrderId,b.OrderId)
		}

	}

	//获取到了订单ID,去重处理下订单ID
	queryOrderId = utils.RemoveRepeatStr(queryOrderId)
	//开始进行统计
	//因为上面可能从客户角度过滤数据了,需要在过滤一次,必须是已经完成的订单
	var orderList []models.Orders
	orderOrm.Select("order_id").Where("order_id in ? and status = ? and delivery_type in ? ",
		queryOrderId,global.OrderStatusOver,orderTypes).Find(&orderList)

	isOkOrderIds :=make([]string,0)
	for _,o:=range orderList{ //获取到层层过滤后的订单ID
		isOkOrderIds = append(isOkOrderIds,o.OrderId)
	}
	isOkOrderIds = utils.RemoveRepeatStr(isOkOrderIds)
	//查询规格订单,把相同的商品数据放一个map中,然后做成list

	var isOkOrderSpecs []models.OrderSpecs
	e.Orm.Table(splitTableRes.OrderSpecs).Where("order_id in ? and status =  ?",
		isOkOrderIds,global.OrderStatusOver).Order(global.OrderTimeKey).Find(&isOkOrderSpecs)

	//查询是否有退货 after_status = 2

	var refundOrderSpecs []models.OrderSpecs
	e.Orm.Table(splitTableRes.OrderSpecs).Select("order_id").Where(
		"order_id in ? and after_status = ?",isOkOrderIds,global.RefundOk).Find(&refundOrderSpecs)

	//sum优惠金额
	var totalCouponMoney string
	SumCouponSql:=fmt.Sprintf("SELECT SUM(coupon_money) AS total_coupon_money  FROM %v WHERE c_id = %v and order_id IN (%v);",
		splitTableRes.OrderTable,userDto.CId,strings.Join(isOkOrderIds,","))
	e.Orm.Table(splitTableRes.OrderTable).Raw(SumCouponSql).Scan(&totalCouponMoney)

	refundOrderIds:=make([]string,0)
	//退货的orderId 在退货表 order_return中查询
	for _,refundRow:=range refundOrderSpecs{
		refundOrderIds = append(refundOrderIds,refundRow.OrderId)
	}
	if len(refundOrderIds) > 0 {
		refundOrderIds = utils.RemoveRepeatStr(refundOrderIds)
		var OrderReturnList []models2.OrderReturn
		e.Orm.Table(splitTableRes.OrderReturn).Select("goods_id,spec_id,number,refund_apply_money,order_id").Where(
			"order_id in ? and status =  ?",refundOrderIds,global.RefundOk).Find(&OrderReturnList)

		//做一个map,统计同样商品退货数据,退货金额

		for _,row:=range OrderReturnList{

			key := utils.Md5(fmt.Sprintf("%v-%v",row.GoodsId,row.SpecId))
			cacheRefundRow,ok:=refundMapCache[key]
			if ok{
				cacheRefundRow = RefundRow{}
			}
			cacheRefundRow.AllNumber += int64(row.Number)

			cacheRefundRow.AllMoney +=row.RefundApplyMoney

			refundMapCache[key] = cacheRefundRow
		}
	}

	//相同规格
	cacheMap:=make(map[string]*GoodsAnalysisTableRow,0)
	for _,row:=range  isOkOrderSpecs{
		key := utils.Md5(fmt.Sprintf("%v-%v",row.GoodsId,row.SpecId))
		//fmt.Println("返回key!!",key,row.GoodsId,row.SpecId)
		cacheRow,ok:=cacheMap[key]
		if !ok {
			cacheRow = &GoodsAnalysisTableRow{
				OrderId: row.OrderId,
				GoodsName: row.GoodsName,
				SpecsName: row.SpecsName,
				Unit: row.Unit,
				AllNumber: int64(row.Number),
				AllMoneyValue: utils.RoundDecimalFlot64(row.AllMoney),
				AllMoney: utils.StringDecimal(row.AllMoney),
			}
		}else {
			cacheRow.AllNumber +=int64(row.Number)
			cacheRow.AllMoneyValue +=utils.RoundDecimalFlot64(row.AllMoney)
			cacheRow.AllMoney  = utils.StringDecimal(cacheRow.AllMoneyValue)
		}
		cacheRefundRow,refundOk:=refundMapCache[key]
		cacheRow.RefundMoney = "0.00"
		if refundOk {//有退货
			cacheRow.RefundCount = cacheRefundRow.AllNumber
			cacheRow.RefundMoneyValue = cacheRefundRow.AllMoney
			cacheRow.RefundMoney = utils.StringDecimal(cacheRow.RefundMoneyValue )
			refundAllCount += cacheRefundRow.AllNumber
			refundAllMoney +=utils.RoundDecimalFlot64(cacheRefundRow.AllMoney)
		}
		cacheRow.Income = utils.StringDecimal(cacheRow.AllMoneyValue - cacheRow.RefundMoneyValue)
		queryAllCount +=int64(row.Number)

		queryAllMoney +=utils.RoundDecimalFlot64(row.AllMoney)

		cacheMap[key] = cacheRow
	}
	resultList:=make([]*GoodsAnalysisTableRow,0)
	for  _,l:=range  cacheMap{
		resultList = append(resultList,l)
	}

	sort.Sort(ByResultAnalysisMoney(resultList))
	queryResult:=map[string]interface{}{
		"calculationCount":map[string]interface{}{
			"queryAllCount":queryAllCount,
			"queryAllMoney":utils.StringDecimal(queryAllMoney),
			"refundAllCount":refundAllCount,
			"refundAllMoney":utils.StringDecimal(refundAllMoney),
			"totalCouponMoney":totalCouponMoney,

		},
		"list":resultList,
		"total":len(resultList),
	}

	e.OK(queryResult,"successful")
	return



}


// 商品分类
func (e DataAnalysis) GoodsClassList(c *gin.Context) {
	req:=AnalysisQuery{}
	s := service.Goods{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
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
	//查询的总订单数
	var queryAllCount int64
	//查询的总金额
	var queryAllMoney float64

	//退货总数量
	var refundAllCount int64
	var refundAllMoney float64
	var queryOrderId []string
	var startTime string
	var endTime string

	nowTime:=time.Now()
	if req.BeginTime!=""{
		startTime = req.BeginTime
	}else {
		startTime = nowTime.AddDate(0,0,-15).Format(time.DateTime)

	}
	if req.EndTime != "" {
		endTime = req.EndTime
	}else {
		endTime =nowTime.Format(time.DateTime)

	}
	whereRangeTime:=fmt.Sprintf("c_id = %v and created_at >= '%v' and created_at <= '%v' ",
		userDto.CId,startTime,endTime)

	//先获取所有的分类->用分类下的商品+规格查询订单（其实就是直接进行了一次商品的查询）

	orderTypes:=make([]int,0)
	switch req.OrderType {
	case global.ExpressSelf:
		orderTypes = append(orderTypes,global.ExpressSelf)
	case global.ExpressSameCity:
		orderTypes = append(orderTypes,global.ExpressSameCity)
	case global.ExpressEms:
		orderTypes = append(orderTypes,global.ExpressEms)
	default:
		orderTypes = append(orderTypes,global.ExpressEms)
		orderTypes = append(orderTypes,global.ExpressSelf)
		orderTypes = append(orderTypes,global.ExpressSameCity)
	}

	refundMapCache:=make(map[int]RefundRow,0)
	//分类聚合的话  只需要做 分类_商品即可。因为商品是直接隶属于分类的
	//获取客户的分表
	splitTableRes := business.GetTableName(userDto.CId, e.Orm)

	var timeRangeOrder []models.Orders
	orderOrm :=e.Orm.Table(splitTableRes.OrderTable).Select("order_id").Where(whereRangeTime).Where("delivery_type in ?",
		orderTypes).Order(global.OrderTimeKey)

	if req.CustomerUser > 0{
		orderOrm = orderOrm.Where("shop_id = ?",req.CustomerUser)

	}
	orderOrm.Find(&timeRangeOrder)
	for _,oo:=range timeRangeOrder{
		queryOrderId = append(queryOrderId,oo.OrderId)
	}
	if req.ClassID > 0{//查询商品分类 -> 获取到订单ID
		queryOrderId = []string{} //重置查询的ID
		var bindGoodsId []int
		e.Orm.Raw(fmt.Sprintf("select goods_id from goods_mark_class where class_id in (%v)",req.ClassID)).Scan(&bindGoodsId)

		var bindGoodsSpecs []models.OrderSpecs
		e.Orm.Table(splitTableRes.OrderSpecs).Select("order_id").Where(whereRangeTime).Where("goods_id in ?",bindGoodsId).Find(&bindGoodsSpecs)

		if len(bindGoodsSpecs) == 0 {
			//查询分类没有那就返回
			queryResult :=map[string]interface{}{
				"calculationCount":map[string]interface{}{
					"queryAllCount":0,
					"queryAllMoney":"0.0",
					"refundAllCount":0,
					"refundAllMoney":"0.0",
				},
				"list":make([]string,0),
				"total":0,
			}
			e.OK(queryResult,"")
			return
		}
		for _,b:=range bindGoodsSpecs{
			queryOrderId = append(queryOrderId,b.OrderId)
		}
	}

	//去重
	queryOrderId = utils.RemoveRepeatStr(queryOrderId)

	//开始进行统计
	//因为上面可能从客户角度过滤数据了,需要在过滤一次,必须是已经完成的订单
	var orderList []models.Orders
	orderOrm.Select("order_id").Where("order_id in ? and status = ? and delivery_type in ? ",
		queryOrderId,global.OrderStatusOver,orderTypes).Find(&orderList)

	isOkOrderIds :=make([]string,0)
	for _,o:=range orderList{ //获取到层层过滤后的订单ID
		isOkOrderIds = append(isOkOrderIds,o.OrderId)
	}
	isOkOrderIds = utils.RemoveRepeatStr(isOkOrderIds)

	//查询规格订单,把相同的商品数据放一个map中,然后做成list

	var isOkOrderSpecs []models.OrderSpecs
	e.Orm.Table(splitTableRes.OrderSpecs).Where("order_id in ? and status =  ?",
		isOkOrderIds,global.OrderStatusOver).Order(global.OrderTimeKey).Find(&isOkOrderSpecs)

	//查询是否有退货 after_status = 2

	var refundOrderSpecs []models.OrderSpecs
	e.Orm.Table(splitTableRes.OrderSpecs).Select("order_id").Where(
		"order_id in ? and after_status = ?",isOkOrderIds,global.RefundOk).Find(&refundOrderSpecs)

	refundOrderIds:=make([]string,0)
	//退货的orderId 在退货表 order_return中查询
	for _,refundRow:=range refundOrderSpecs{
		refundOrderIds = append(refundOrderIds,refundRow.OrderId)
	}
	if len(refundOrderIds) > 0 {
		refundOrderIds = utils.RemoveRepeatStr(refundOrderIds)
		var OrderReturnList []models2.OrderReturn
		e.Orm.Table(splitTableRes.OrderReturn).Select("goods_id,spec_id,number,refund_apply_money,order_id").Where(
			"order_id in ? and status =  ?",refundOrderIds,global.RefundOk).Find(&OrderReturnList)

		//做一个map,统计同样商品退货数据,退货金额
		for _,row:=range OrderReturnList{

			cacheRefundRow,ok:=refundMapCache[row.GoodsId]
			if ok{
				cacheRefundRow = RefundRow{}
			}
			cacheRefundRow.AllNumber += int64(row.Number)

			cacheRefundRow.AllMoney +=row.RefundApplyMoney
			refundMapCache[row.GoodsId] = cacheRefundRow
		}
	}

	//sum优惠金额
	var totalCouponMoney string
	SumCouponSql:=fmt.Sprintf("SELECT SUM(coupon_money) AS total_coupon_money  FROM %v WHERE c_id = %v and order_id IN (%v);",
		splitTableRes.OrderTable,userDto.CId,strings.Join(isOkOrderIds,","))
	e.Orm.Table(splitTableRes.OrderTable).Raw(SumCouponSql).Scan(&totalCouponMoney)

	//存放到分类中,分类是直接和商品关联的,不是和分类关联
	cacheMap:=make(map[int]*GoodsClsAnalysisTableRow,0)
	//cacheMap是把同一个商品的所以规格统计放一起了
	isOkGoodsId:=make([]string,0)
	for _,row:=range  isOkOrderSpecs{
		cacheRow,ok:=cacheMap[row.GoodsId]
		if !ok {
			cacheRow = &GoodsClsAnalysisTableRow{
				GoodsId: row.GoodsId,
				AllNumber: int64(row.Number),
				AllMoneyValue: utils.RoundDecimalFlot64(row.AllMoney),
				AllMoney: utils.StringDecimal(row.AllMoney),
			}
		}else {
			cacheRow.AllNumber +=int64(row.Number)
			cacheRow.AllMoneyValue +=utils.RoundDecimalFlot64(row.AllMoney)
			cacheRow.AllMoney = utils.StringDecimal(cacheRow.AllMoneyValue)
		}
		cacheRefundRow,refundOk:=refundMapCache[row.GoodsId]
		if refundOk {//有退货
			cacheRow.RefundCount = cacheRefundRow.AllNumber
			cacheRow.RefundMoneyValue = cacheRefundRow.AllMoney
			cacheRow.RefundMoney = utils.StringDecimal(cacheRefundRow.AllMoney)
			refundAllCount += cacheRefundRow.AllNumber
			refundAllMoney +=utils.RoundDecimalFlot64(cacheRefundRow.AllMoney)
		}
		cacheRow.Income = utils.StringDecimal(cacheRow.AllMoneyValue - cacheRow.RefundMoneyValue)
		queryAllCount +=int64(row.Number)

		queryAllMoney +=utils.RoundDecimalFlot64(row.AllMoney)

		cacheMap[row.GoodsId] = cacheRow

		isOkGoodsId = append(isOkGoodsId,fmt.Sprintf("%v",row.GoodsId))
	}
	isOkGoodsId = utils.RemoveRepeatStr(isOkGoodsId)

	//通过商品ID获取分类
	var classList []models.GoodsClass
	e.Orm.Model(&models.GoodsClass{}).Where("c_id = ? and id in ?",userDto.CId,isOkGoodsId).Find(&classList)

	var bindClsGoodsId []struct{
		GoodsId int `json:"goods_id"`
		ClassId int `json:"class_id"`
	}
	e.Orm.Raw(fmt.Sprintf("select * from goods_mark_class where goods_id in (%v)",
		strings.Join(isOkGoodsId,","))).Scan(&bindClsGoodsId)
	clsIds:=make([]int,0)
	goodsBindCls:=make(map[int]int,0)//商品:分类ID
	for _,cls:=range bindClsGoodsId{
		clsIds = append(clsIds,cls.ClassId)
		goodsBindCls[cls.GoodsId] = cls.ClassId
	}
	//构造一个分类的map信息
	var goodsClsList []models.GoodsClass
	e.Orm.Model(&models.GoodsClass{}).Where("c_id = ? and id in ?",userDto.CId,clsIds).Find(&goodsClsList)
	clsInfoMap:=make(map[int]models.GoodsClass,0)
	for _,cls:=range goodsClsList{
		clsInfoMap[cls.Id] = cls
	}

	//创建一个无分类的默认缓存内容
	noBindCache:=models.GoodsClass{
		Name: "暂无分类!!!",
	}
	noBindCache.Id = -1
	clsInfoMap[-1] = noBindCache

	mergeMap:=make(map[int]*GoodsClsAnalysisTableRow,0)//把相同分类的商品合并在一起

	//把同一个分类下的商品归纳到一起
	for  _,l:=range cacheMap{ //对有效的订单进行一次 商品分类查询
		clsId,ok:=goodsBindCls[l.GoodsId]
		if !ok{
			clsId = -1
		}
		clsModule,clsOk:=clsInfoMap[clsId]
		if !clsOk{
			//fmt.Println("这个商品的分类 的model数据不存在",l.GoodsId,clsId)
			clsId = -1

			clsModule = clsInfoMap[-1]
		}
		mergeRow,mergeOk:=mergeMap[clsId]
		if !mergeOk{
			mergeRow = &GoodsClsAnalysisTableRow{}
		}
		mergeRow.ClsId = clsModule.Id
		mergeRow.AllMoneyValue +=l.AllMoneyValue
		mergeRow.AllMoney = utils.StringDecimal(mergeRow.AllMoneyValue)
		mergeRow.ClsName = clsModule.Name
		mergeRow.AllNumber +=l.AllNumber
		mergeRow.RefundMoneyValue +=l.RefundMoneyValue
		mergeRow.RefundMoney = utils.StringDecimal(mergeRow.RefundMoneyValue)
		mergeRow.RefundCount +=l.RefundCount
		mergeRow.Income = utils.StringDecimal(mergeRow.AllMoneyValue - mergeRow.RefundMoneyValue)
		mergeMap[clsId] = mergeRow
	}

	resultList:=make([]*GoodsClsAnalysisTableRow,0)
	for  _,l:=range  mergeMap{
		resultList = append(resultList,l)
	}

	sort.Sort(ByResultClsAnalysisMoney(resultList))
	queryResult:=map[string]interface{}{
		"calculationCount":map[string]interface{}{
			"queryAllCount":queryAllCount,
			"queryAllMoney":utils.StringDecimal(queryAllMoney),
			"refundAllCount":refundAllCount,
			"refundAllMoney":refundAllMoney,
			"totalCouponMoney":totalCouponMoney,
		},
		"list":resultList,
		"total":len(resultList),
	}

	e.OK(queryResult,"successful")


	//分类名称 | 订货数量 | 订货金额 | 优惠金额 | 实际销售金额 | 退货数量 | 退货金额 | 净销售收入
	return
}



// 商品品牌

func (e DataAnalysis) GoodsBrandList(c *gin.Context) {
	req:=AnalysisQuery{}
	s := service.Goods{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
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
	//查询的总订单数
	var queryAllCount int64
	//查询的总金额
	var queryAllMoney float64

	//退货总数量
	var refundAllCount int64
	var refundAllMoney float64
	var queryOrderId []string
	var startTime string
	var endTime string

	nowTime:=time.Now()
	if req.BeginTime!=""{
		startTime = req.BeginTime
	}else {
		startTime = nowTime.AddDate(0,0,-15).Format(time.DateTime)

	}
	if req.EndTime != "" {
		endTime = req.EndTime
	}else {
		endTime =nowTime.Format(time.DateTime)

	}
	whereRangeTime:=fmt.Sprintf("c_id = %v and created_at >= '%v' and created_at <= '%v' ",
		userDto.CId,startTime,endTime)

	//先获取所有的分类->用分类下的商品+规格查询订单（其实就是直接进行了一次商品的查询）

	orderTypes:=make([]int,0)
	switch req.OrderType {
	case global.ExpressSelf:
		orderTypes = append(orderTypes,global.ExpressSelf)
	case global.ExpressSameCity:
		orderTypes = append(orderTypes,global.ExpressSameCity)
	case global.ExpressEms:
		orderTypes = append(orderTypes,global.ExpressEms)
	default:
		orderTypes = append(orderTypes,global.ExpressEms)
		orderTypes = append(orderTypes,global.ExpressSelf)
		orderTypes = append(orderTypes,global.ExpressSameCity)
	}

	refundMapCache:=make(map[int]RefundRow,0)
	//分类聚合的话  只需要做 分类_商品即可。因为商品是直接隶属于分类的
	//获取客户的分表
	splitTableRes := business.GetTableName(userDto.CId, e.Orm)

	var timeRangeOrder []models.Orders
	orderOrm :=e.Orm.Table(splitTableRes.OrderTable).Select("order_id").Where(whereRangeTime).Where("delivery_type in ?",
		orderTypes).Order(global.OrderTimeKey)

	if req.CustomerUser > 0{
		orderOrm = orderOrm.Where("shop_id = ?",req.CustomerUser)

	}
	orderOrm.Find(&timeRangeOrder)
	for _,oo:=range timeRangeOrder{
		queryOrderId = append(queryOrderId,oo.OrderId)
	}
	if req.BrandId > 0{//查询商品分类 -> 获取到订单ID
		queryOrderId = []string{} //重置查询的ID
		var bindGoodsId []int
		e.Orm.Raw(fmt.Sprintf("select goods_id from goods_mark_brand where class_id in (%v)",req.BrandId)).Scan(&bindGoodsId)

		var bindGoodsSpecs []models.OrderSpecs
		e.Orm.Table(splitTableRes.OrderSpecs).Select("order_id").Where(whereRangeTime).Where("goods_id in ?",bindGoodsId).Find(&bindGoodsSpecs)

		if len(bindGoodsSpecs) == 0 {
			//查询分类没有那就返回
			queryResult :=map[string]interface{}{
				"calculationCount":map[string]interface{}{
					"queryAllCount":0,
					"queryAllMoney":"0.0",
					"refundAllCount":0,
					"refundAllMoney":"0.0",
				},
				"list":make([]string,0),
				"total":0,
			}
			e.OK(queryResult,"")
			return
		}
		for _,b:=range bindGoodsSpecs{
			queryOrderId = append(queryOrderId,b.OrderId)
		}
	}

	//去重
	queryOrderId = utils.RemoveRepeatStr(queryOrderId)

	//开始进行统计
	//因为上面可能从客户角度过滤数据了,需要在过滤一次,必须是已经完成的订单
	var orderList []models.Orders
	orderOrm.Select("order_id").Where("order_id in ? and status = ? and delivery_type in ? ",
		queryOrderId,global.OrderStatusOver,orderTypes).Find(&orderList)

	isOkOrderIds :=make([]string,0)
	for _,o:=range orderList{ //获取到层层过滤后的订单ID
		isOkOrderIds = append(isOkOrderIds,o.OrderId)
	}
	isOkOrderIds = utils.RemoveRepeatStr(isOkOrderIds)

	//查询规格订单,把相同的商品数据放一个map中,然后做成list

	var isOkOrderSpecs []models.OrderSpecs
	e.Orm.Table(splitTableRes.OrderSpecs).Where("order_id in ? and status =  ?",
		isOkOrderIds,global.OrderStatusOver).Order(global.OrderTimeKey).Find(&isOkOrderSpecs)

	//查询是否有退货 after_status = 2

	var refundOrderSpecs []models.OrderSpecs
	e.Orm.Table(splitTableRes.OrderSpecs).Select("order_id").Where(
		"order_id in ? and after_status = ?",isOkOrderIds,global.RefundOk).Find(&refundOrderSpecs)

	refundOrderIds:=make([]string,0)
	//退货的orderId 在退货表 order_return中查询
	for _,refundRow:=range refundOrderSpecs{
		refundOrderIds = append(refundOrderIds,refundRow.OrderId)
	}
	if len(refundOrderIds) > 0 {
		refundOrderIds = utils.RemoveRepeatStr(refundOrderIds)
		var OrderReturnList []models2.OrderReturn
		e.Orm.Table(splitTableRes.OrderReturn).Select("goods_id,spec_id,number,refund_apply_money,order_id").Where(
			"order_id in ? and status =  ?",refundOrderIds,global.RefundOk).Find(&OrderReturnList)

		//做一个map,统计同样商品退货数据,退货金额
		for _,row:=range OrderReturnList{

			cacheRefundRow,ok:=refundMapCache[row.GoodsId]
			if ok{
				cacheRefundRow = RefundRow{}
			}
			cacheRefundRow.AllNumber += int64(row.Number)

			cacheRefundRow.AllMoney +=row.RefundApplyMoney
			refundMapCache[row.GoodsId] = cacheRefundRow
		}
	}

	//sum优惠金额
	var totalCouponMoney string
	SumCouponSql:=fmt.Sprintf("SELECT SUM(coupon_money) AS total_coupon_money  FROM %v WHERE c_id = %v and order_id IN (%v);",
		splitTableRes.OrderTable,userDto.CId,strings.Join(isOkOrderIds,","))
	e.Orm.Table(splitTableRes.OrderTable).Raw(SumCouponSql).Scan(&totalCouponMoney)

	//存放到分类中,分类是直接和商品关联的,不是和分类关联
	cacheMap:=make(map[int]*GoodsClsAnalysisTableRow,0)
	//cacheMap是把同一个商品的所以规格统计放一起了
	isOkGoodsId:=make([]string,0)
	for _,row:=range  isOkOrderSpecs{
		cacheRow,ok:=cacheMap[row.GoodsId]
		if !ok {
			cacheRow = &GoodsClsAnalysisTableRow{
				GoodsId: row.GoodsId,
				AllNumber: int64(row.Number),
				AllMoneyValue: utils.RoundDecimalFlot64(row.AllMoney),
				AllMoney: utils.StringDecimal(row.AllMoney),
			}
		}else {
			cacheRow.AllNumber +=int64(row.Number)
			cacheRow.AllMoneyValue +=utils.RoundDecimalFlot64(row.AllMoney)
			cacheRow.AllMoney = utils.StringDecimal(cacheRow.AllMoneyValue)

		}
		cacheRefundRow,refundOk:=refundMapCache[row.GoodsId]
		if refundOk {//有退货
			cacheRow.RefundCount = cacheRefundRow.AllNumber
			cacheRow.RefundMoneyValue = cacheRefundRow.AllMoney
			cacheRow.RefundMoney = utils.StringDecimal(cacheRefundRow.AllMoney)
			refundAllCount += cacheRefundRow.AllNumber
			refundAllMoney +=utils.RoundDecimalFlot64(cacheRefundRow.AllMoney)
		}
		cacheRow.Income = utils.StringDecimal(cacheRow.AllMoneyValue - cacheRow.RefundMoneyValue)
		queryAllCount +=int64(row.Number)

		queryAllMoney +=utils.RoundDecimalFlot64(row.AllMoney)

		cacheMap[row.GoodsId] = cacheRow

		isOkGoodsId = append(isOkGoodsId,fmt.Sprintf("%v",row.GoodsId))
	}
	isOkGoodsId = utils.RemoveRepeatStr(isOkGoodsId)

	//通过商品ID获取分类
	var brandList []models.GoodsBrand
	e.Orm.Model(&models.GoodsBrand{}).Where("c_id = ? and id in ?",userDto.CId,isOkGoodsId).Find(&brandList)

	var bindBrandGoodsId []struct{
		GoodsId int `json:"goods_id"`
		BrandId int `json:"brand_id"`
	}
	e.Orm.Raw(fmt.Sprintf("select * from goods_mark_brand where goods_id in (%v)",
		strings.Join(isOkGoodsId,","))).Scan(&bindBrandGoodsId)
	brandIds:=make([]int,0)
	goodsBindBrand:=make(map[int]int,0)//商品:品牌ID
	for _,brand:=range bindBrandGoodsId{
		brandIds = append(brandIds,brand.BrandId)
		goodsBindBrand[brand.GoodsId] = brand.BrandId
	}
	//构造一个分类的map信息
	var goodsBrandList []models.GoodsBrand
	e.Orm.Model(&models.GoodsBrand{}).Where("c_id = ? and id in ?",userDto.CId,brandIds).Find(&goodsBrandList)
	brandInfoMap:=make(map[int]models.GoodsBrand,0)
	for _,brand:=range goodsBrandList{
		brandInfoMap[brand.Id] = brand
	}

	//创建一个无分类的默认缓存内容
	noBindCache:=models.GoodsBrand{
		Name: "暂无品牌 !",
	}
	noBindCache.Id = -1
	brandInfoMap[-1] = noBindCache

	mergeMap:=make(map[int]*GoodsBrandAnalysisTableRow,0)//把相同分类的商品合并在一起

	//把同一个分类下的商品归纳到一起
	for  _,l:=range cacheMap{ //对有效的订单进行一次 商品分类查询 ,品牌是一个非必填的字段
		brandId,ok:=goodsBindBrand[l.GoodsId]
		if !ok{
			brandId = -1
		}
		Module,clsOk:=brandInfoMap[brandId]
		if !clsOk{
			//fmt.Println("这个商品的分类 的model数据不存在",l.GoodsId,clsId)
			brandId = -1

			Module = brandInfoMap[-1]
		}
		mergeRow,mergeOk:=mergeMap[brandId]
		if !mergeOk{
			mergeRow = &GoodsBrandAnalysisTableRow{}
		}
		mergeRow.BrandId = Module.Id
		mergeRow.AllMoney +=l.AllMoneyValue
		mergeRow.BrandName = Module.Name
		mergeRow.AllNumber +=l.AllNumber
		mergeRow.RefundMoneyValue +=l.RefundMoneyValue
		mergeRow.RefundMoney = utils.StringDecimal(mergeRow.RefundMoneyValue )
		mergeRow.RefundCount +=l.RefundCount
		mergeRow.Income = utils.StringDecimal(mergeRow.AllMoney - mergeRow.RefundMoneyValue)
		mergeMap[brandId] = mergeRow
	}

	resultList:=make([]*GoodsBrandAnalysisTableRow,0)
	for  _,l:=range  mergeMap{
		resultList = append(resultList,l)
	}

	sort.Sort(ByResultBrandAnalysisMoney(resultList))
	queryResult:=map[string]interface{}{
		"calculationCount":map[string]interface{}{
			"queryAllCount":queryAllCount,
			"queryAllMoney":utils.StringDecimal(queryAllMoney),
			"refundAllCount":refundAllCount,
			"refundAllMoney":utils.StringDecimal(refundAllMoney),
			"totalCouponMoney":utils.StringDecimal(totalCouponMoney),
		},
		"list":resultList,
		"total":len(resultList),
	}

	e.OK(queryResult,"successful")

	return
}

//毛利统计
//商品-规格名称 | 销售收入 | 优惠抵扣 | 	实际销售收入 | 销售成本 | 销售毛利 | 退货金额 | 退货金额 | 销售净收入 | 毛利 | 毛利率

func (e DataAnalysis) Grosslist(c *gin.Context) {
	req:=AnalysisQuery{}
	s := service.Goods{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
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

	//查询的总订单数
	var queryAllCount int64
	//查询的总金额
	var queryAllMoney float64

	//退货总数量
	var refundAllCount int64
	var refundAllMoney float64
	//获取客户的分表
	splitTableRes := business.GetTableName(userDto.CId, e.Orm)
	var startTime string
	var endTime string
	nowTime:=time.Now()
	if req.BeginTime!=""{
		startTime = req.BeginTime
	}else {
		startTime = nowTime.AddDate(0,0,-15).Format(time.DateTime)

	}
	if req.EndTime != "" {
		endTime = req.EndTime
	}else {
		endTime =nowTime.Format(time.DateTime)

	}


	whereRangeTime:=fmt.Sprintf("c_id = %v and created_at >= '%v' and created_at <= '%v' ",
		userDto.CId,startTime,endTime)

	//不同的订单完结状态不一样
	orderTypes:=make([]int,0)
	switch req.OrderType {
	case global.ExpressSelf:
		orderTypes = append(orderTypes,global.ExpressSelf)
	case global.ExpressSameCity:
		orderTypes = append(orderTypes,global.ExpressSameCity)
	case global.ExpressEms:
		orderTypes = append(orderTypes,global.ExpressEms)
	default:
		orderTypes = append(orderTypes,global.ExpressEms)
		orderTypes = append(orderTypes,global.ExpressSelf)
		orderTypes = append(orderTypes,global.ExpressSameCity)
	}
	refundMapCache:=make(map[string]RefundRow,0)
	//考虑 优惠金额
	//考虑退货的商品
	var queryOrderId []string
	var timeRangeOrder []models.Orders
	orderOrm :=e.Orm.Table(splitTableRes.OrderTable).Select("order_id").Where(whereRangeTime).Where("delivery_type in ?",
		orderTypes).Order(global.OrderTimeKey)

	if req.CustomerUser > 0{
		orderOrm = orderOrm.Where("shop_id = ?",req.CustomerUser)
	}

	orderOrm.Find(&timeRangeOrder)
	for _,oo:=range timeRangeOrder{
		queryOrderId = append(queryOrderId,oo.OrderId)
	}
	if req.ClassID > 0{//查询商品分类 -> 获取到订单ID
		queryOrderId = []string{} //重置查询的ID
		var bindGoodsId []int
		e.Orm.Raw(fmt.Sprintf("select goods_id from goods_mark_class where class_id in (%v)",req.ClassID)).Scan(&bindGoodsId)

		var bindGoodsSpecs []models.OrderSpecs
		e.Orm.Table(splitTableRes.OrderSpecs).Select("order_id").Where(whereRangeTime).Where("goods_id in ?",bindGoodsId).Find(&bindGoodsSpecs)

		if len(bindGoodsSpecs) == 0 {
			//查询分类没有那就返回
			queryResult :=map[string]interface{}{
				"calculationCount":map[string]interface{}{
					"queryAllCount":0,
					"queryAllMoney":"0.0",
					"refundAllCount":0,
					"refundAllMoney":"0.0",
				},
				"list":make([]string,0),
				"total":0,
			}
			e.OK(queryResult,"")
			return
		}
		for _,b:=range bindGoodsSpecs{
			queryOrderId = append(queryOrderId,b.OrderId)
		}

	}
	if req.SpecsId > 0 {//规格查询 -> 获取到订单ID
		queryOrderId = []string{} //重置查询的ID
		var specsBindOrder []models.OrderSpecs
		e.Orm.Table(splitTableRes.OrderSpecs).Where(whereRangeTime).Where("spec_id = ?",
			req.SpecsId).Find(&specsBindOrder)
		for _,b:=range specsBindOrder{
			queryOrderId = append(queryOrderId,b.OrderId)
		}

	}

	//获取到了订单ID,去重处理下订单ID
	queryOrderId = utils.RemoveRepeatStr(queryOrderId)
	//开始进行统计
	//因为上面可能从客户角度过滤数据了,需要在过滤一次,必须是已经完成的订单
	var orderList []models.Orders
	orderOrm.Select("order_id").Where("order_id in ? and status = ? and delivery_type in ? ",
		queryOrderId,global.OrderStatusOver,orderTypes).Find(&orderList)

	isOkOrderIds :=make([]string,0)
	for _,o:=range orderList{ //获取到层层过滤后的订单ID
		isOkOrderIds = append(isOkOrderIds,o.OrderId)
	}
	isOkOrderIds = utils.RemoveRepeatStr(isOkOrderIds)
	//查询规格订单,把相同的商品数据放一个map中,然后做成list

	var isOkOrderSpecs []models.OrderSpecs
	e.Orm.Table(splitTableRes.OrderSpecs).Where("order_id in ? and status =  ?",
		isOkOrderIds,global.OrderStatusOver).Order(global.OrderTimeKey).Find(&isOkOrderSpecs)

	//查询是否有退货 after_status = 2
	var refundOrderSpecs []models.OrderSpecs
	e.Orm.Table(splitTableRes.OrderSpecs).Select("order_id").Where(
		"order_id in ? and after_status = ?",isOkOrderIds,global.RefundOk).Find(&refundOrderSpecs)
	refundOrderIds:=make([]string,0)
	//退货的orderId 在退货表 order_return中查询
	for _,refundRow:=range refundOrderSpecs{
		refundOrderIds = append(refundOrderIds,refundRow.OrderId)
	}
	if len(refundOrderIds) > 0 {
		refundOrderIds = utils.RemoveRepeatStr(refundOrderIds)
		var OrderReturnList []models2.OrderReturn
		e.Orm.Table(splitTableRes.OrderReturn).Select("goods_id,spec_id,number,refund_apply_money,order_id").Where(
			"order_id in ? and status =  ?",refundOrderIds,global.RefundOk).Find(&OrderReturnList)

		//做一个map,统计同样商品退货数据,退货金额

		for _,row:=range OrderReturnList{

			key := utils.Md5(fmt.Sprintf("%v-%v",row.GoodsId,row.SpecId))
			cacheRefundRow,ok:=refundMapCache[key]
			if ok{
				cacheRefundRow = RefundRow{}
			}
			cacheRefundRow.AllNumber += int64(row.Number)

			cacheRefundRow.AllMoney +=row.RefundApplyMoney

			refundMapCache[key] = cacheRefundRow
		}
	}


	//sum优惠金额
	var totalCouponMoney string
	SumCouponSql:=fmt.Sprintf("SELECT SUM(coupon_money) AS total_coupon_money  FROM %v WHERE c_id = %v and order_id IN (%v);",
		splitTableRes.OrderTable,userDto.CId,strings.Join(isOkOrderIds,","))
	e.Orm.Table(splitTableRes.OrderTable).Raw(SumCouponSql).Scan(&totalCouponMoney)


	//查询商品的成本价
	//如果开启了库存,应该是去库存中获取成本价
	originalPriceMap:=make(map[string]float64,0)
	isOpenInventory := service.IsOpenInventory(userDto.CId,e.Orm)
	var searchGoodsSql []string
	for _,row:=range  isOkOrderSpecs{
		var queryStr string
		if isOpenInventory{
			queryStr =fmt.Sprintf("(goods_id = %v AND spec_id = %v)",row.GoodsId,row.SpecId)
		}else {
			queryStr =fmt.Sprintf("(goods_id = %v AND id = %v)",row.GoodsId,row.SpecId)
		}

		searchGoodsSql = append(searchGoodsSql,queryStr)
	}

	if isOpenInventory{
		var Inventorying []models2.Inventory
		e.Orm.Model(&models2.Inventory{}).Select("original_price,goods_id,spec_id").Where(
			strings.Join(searchGoodsSql,"OR")).Find(&Inventorying)
		for _,row:=range Inventorying{
			originalPriceMap[fmt.Sprintf("%v-%v",row.GoodsId,row.SpecId)] = row.OriginalPrice
		}
	}else {
		var GoodsSpecsList []models.GoodsSpecs
		e.Orm.Model(&models.GoodsSpecs{}).Select("original,goods_id,id").Where(
			strings.Join(searchGoodsSql,"OR")).Find(&GoodsSpecsList)
		for _,row:=range GoodsSpecsList{
			originalPriceMap[fmt.Sprintf("%v-%v",row.GoodsId,row.Id)] = row.Original
		}
	}

	//相同规格
	cacheMap:=make(map[string]*GoodsAnalysisTableRow,0)
	for _,row:=range  isOkOrderSpecs{

		key := utils.Md5(fmt.Sprintf("%v-%v",row.GoodsId,row.SpecId))
		//fmt.Println("返回key!!",key,row.GoodsId,row.SpecId)
		cacheRow,ok:=cacheMap[key]
		if !ok {
			cacheRow = &GoodsAnalysisTableRow{
				OrderId: row.OrderId,
				GoodsName: row.GoodsName,
				SpecsName: row.SpecsName,
				Unit: row.Unit,
				AllNumber: int64(row.Number),
				AllMoneyValue: utils.RoundDecimalFlot64(row.AllMoney),
				AllMoney: utils.StringDecimal(row.AllMoney),
			}
		}else {
			cacheRow.AllNumber +=int64(row.Number)
			cacheRow.AllMoneyValue +=utils.RoundDecimalFlot64(row.AllMoney)
			cacheRow.AllMoney  = utils.StringDecimal(cacheRow.AllMoneyValue)
		}
		cacheRefundRow,refundOk:=refundMapCache[key]
		cacheRow.RefundMoney = "0.00"
		if refundOk {//有退货
			cacheRow.RefundCount = cacheRefundRow.AllNumber
			cacheRow.RefundMoneyValue = cacheRefundRow.AllMoney
			cacheRow.RefundMoney = utils.StringDecimal(cacheRow.RefundMoneyValue )
			refundAllCount += cacheRefundRow.AllNumber
			refundAllMoney +=utils.RoundDecimalFlot64(cacheRefundRow.AllMoney)
		}
		originalPrice:=originalPriceMap[fmt.Sprintf("%v-%v",row.GoodsId,row.SpecId)]

		if originalPrice > 0 { //只有成本价才值得计算
			//销售成本  = 销售数量 * 成本
			SalesCost :=utils.RoundDecimalFlot64(float64(cacheRow.AllNumber) * originalPrice)
			cacheRow.Cost =  utils.StringDecimal(SalesCost)


			SalesGross :=utils.RoundDecimalFlot64(cacheRow.AllMoneyValue  - SalesCost)
			//销售毛利 = （实际销售收入-实际销售成本），不包含退货数据。
			cacheRow.SalesGross  = utils.StringDecimal(SalesGross)
			//销售毛利率=销售毛利/实际销售收入
			cacheRow.SalesGrossProfit  = fmt.Sprintf("%.2f%%",(SalesGross / cacheRow.AllMoneyValue)  * 100)
			// 退回成本=货物单价×退回的货物数量

			RefundCost:= utils.RoundDecimalFlot64(originalPrice * float64(cacheRow.RefundCount))
			cacheRow.RefundCost = utils.StringDecimal(RefundCost)

			//整个毛利计算

			//毛利=[(实际销售金-实际退货金额)-(销售成本-退货成本)]
			Gross := utils.RoundDecimalFlot64(
				(cacheRow.AllMoneyValue - cacheRow.RefundMoneyValue) - (SalesCost -  RefundCost))
			cacheRow.Gross = utils.StringDecimal(Gross)

			//毛利率=[(实际销售金-实际退货金额)-(销售成本-退货成本)]/实际销售金额

			cacheRow.GrossProfit = fmt.Sprintf("%.2f%%",(Gross / cacheRow.AllMoneyValue)  * 100)
		}



		cacheRow.Income = utils.StringDecimal(cacheRow.AllMoneyValue - cacheRow.RefundMoneyValue)
		queryAllCount +=int64(row.Number)

		queryAllMoney +=utils.RoundDecimalFlot64(row.AllMoney)

		cacheMap[key] = cacheRow
	}
	resultList:=make([]*GoodsAnalysisTableRow,0)
	for  _,l:=range  cacheMap{
		resultList = append(resultList,l)
	}

	sort.Sort(ByResultAnalysisMoney(resultList))
	queryResult:=map[string]interface{}{
		"calculationCount":map[string]interface{}{
			"queryAllCount":queryAllCount,
			"queryAllMoney":utils.StringDecimal(queryAllMoney),
			"refundAllCount":refundAllCount,
			"refundAllMoney":utils.StringDecimal(refundAllMoney),
			"totalCouponMoney":utils.StringDecimal(totalCouponMoney),

		},
		"list":resultList,
		"total":len(resultList),
	}

	e.OK(queryResult,"successful")
	return



}



