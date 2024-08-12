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
	"time"
)
type DataAnalysis struct {
	api.Api
}

type GoodsAnalysisTableRow struct {
	OrderId string `json:"order_id"`
	GoodsName string `json:"goods_name"`
	SpecsName string `json:"specs_name"`
	AllNumber int64 `json:"all_number"`
	AllMoney float64 `json:"all_money"`
	Unit string `json:"unit"`
	RefundCount int64 `json:"refund_number"`
	RefundMoney float64 `json:"refund_money"`
}
type AnalysisQuery struct {
	OrderType    int    `form:"orderType" `
	ClassID    int    `form:"classId" `
	CustomerUser    int    `form:"customerUser" `
	PageIndex    int    `form:"pageIndex"`
	SpecsId    int    `form:"specsId" `
	PageSize    int    `form:"pageSize" `
	BeginTime  string `form:"beginTime"`
	EndTime string `form:"endTime"`
}

type RefundRow struct {
	AllNumber int64 `json:"all_number"`
	AllMoney float64 `json:"all_money"`
}
// 销售统计


func (e DataAnalysis) GoodsCount(c *gin.Context) {
	s := service.Goods{}
	err := e.MakeContext(c).
		MakeOrm().
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	return
}
// 列表

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
	startTime:=nowTime.AddDate(0,0,-15).Format(time.DateTime)
	endTime:=nowTime.Format(time.DateTime)

	//限制最多查询60天的数据
	//开始时间
	if req.BeginTime!=""{
		startTime = req.BeginTime
	}
	if req.EndTime != "" {
		endTime = req.EndTime
	}


	whereRangeTime:=fmt.Sprintf("c_id = %v and created_at >= '%v' and created_at <= '%v' ",
		userDto.CId,startTime,endTime)
	switch req.OrderType {
	case global.ExpressSelf://自提

	case global.ExpressSameCity:
		refundMapCache:=make(map[string]RefundRow,0)
		//考虑 优惠金额
		//考虑退货的商品
		var queryOrderId []string
		var timeRangeOrder []models.Orders
		orderOrm :=e.Orm.Table(splitTableRes.OrderTable).Select("order_id").Where(whereRangeTime)

		if req.CustomerUser > 0{
			orderOrm = orderOrm.Where("shop_id = ?",req.CustomerUser)
		}
		orderOrm.Find(&timeRangeOrder)

		for _,oo:=range timeRangeOrder{
			queryOrderId = append(queryOrderId,oo.OrderId)
		}
		if req.ClassID > 0{//查询商品分类 -> 获取到订单ID

			fmt.Println("查询分类")
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

			var specsBindOrder []models.OrderSpecs
			e.Orm.Table(splitTableRes.OrderSpecs).Where("? and spec_id = ?",
				whereRangeTime,req.SpecsId).Find(&specsBindOrder)
			for _,b:=range specsBindOrder{
				queryOrderId = append(queryOrderId,b.OrderId)
			}

		}

		//获取到了订单ID,去重处理下订单ID
		queryOrderId = utils.RemoveRepeatStr(queryOrderId)
		//开始进行统计
		//因为上面可能从客户角度过滤数据了,需要在过滤一次,必须是已经完成的订单
		var orderList []models.Orders
		orderOrm.Select("order_id").Where("order_id in ? and status =  ?",
			queryOrderId,global.OrderStatusOver).Find(&orderList)

		isGroupBy :=make([]string,0)
		for _,o:=range orderList{ //获取到层层过滤后的订单ID
			isGroupBy = append(isGroupBy,o.OrderId)
		}
		isGroupBy = utils.RemoveRepeatStr(isGroupBy)
		//查询规格订单,把相同的商品数据放一个map中,然后做成list

		var isOkOrderSpecs []models.OrderSpecs
		e.Orm.Table(splitTableRes.OrderSpecs).Where("order_id in ? and status =  ?",isGroupBy,global.OrderStatusOver).Find(&isOkOrderSpecs)

		//查询是否有退货 after_status = 2

		var refundOrderSpecs []models.OrderSpecs
		e.Orm.Table(splitTableRes.OrderSpecs).Select("order_id").Where(
			"order_id in ? and after_status = ?",isGroupBy,global.RefundOk).Find(&refundOrderSpecs)

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
					AllMoney: utils.RoundDecimalFlot64(row.AllMoney),
				}
			}else {
				cacheRow.AllNumber +=int64(row.Number)
				cacheRow.AllMoney +=utils.RoundDecimalFlot64(row.AllMoney)
			}
			cacheRefundRow,refundOk:=refundMapCache[key]
			if refundOk {//有退货
				cacheRow.RefundCount = cacheRefundRow.AllNumber
				cacheRow.RefundMoney = cacheRefundRow.AllMoney
				refundAllCount += cacheRefundRow.AllNumber
				refundAllMoney +=utils.RoundDecimalFlot64(cacheRefundRow.AllMoney)
			}

			queryAllCount +=int64(row.Number)

			queryAllMoney +=utils.RoundDecimalFlot64(row.AllMoney)

			cacheMap[key] = cacheRow
		}
		resultList:=make([]interface{},0)
		for  _,l:=range  cacheMap{
			resultList = append(resultList,l)
		}
		
		queryResult:=map[string]interface{}{
			"calculationCount":map[string]interface{}{
				"queryAllCount":queryAllCount,
				"queryAllMoney":utils.StringDecimal(queryAllMoney),
				"refundAllCount":refundAllCount,
				"refundAllMoney":refundAllMoney,

			},
			"list":resultList,
			"total":len(resultList),
		}
		
		e.OK(queryResult,"successful")
		return



	case global.ExpressEms:

	default:
		e.Error(500, nil, "订单类型不存在")
		return

	}


	fmt.Println("GoodsList",req)
	//商品名称 | 订货数量 | 订货金额 | 优惠金额 | 实际销售金额 | 退货数量 | 退货金额 | 净销售收入
	return
}


// 商品分类
func (e DataAnalysis) GoodsClassList(c *gin.Context) {
	s := service.Goods{}
	err := e.MakeContext(c).
		MakeOrm().
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	//分类名称 | 订货数量 | 订货金额 | 优惠金额 | 实际销售金额 | 退货数量 | 退货金额 | 净销售收入
	return
}

// 商品品牌
func (e DataAnalysis) GoodsBrandList(c *gin.Context) {
	s := service.Goods{}
	err := e.MakeContext(c).
		MakeOrm().
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	//商品品牌 | 订货数量 | 订货金额 | 优惠金额 | 实际销售金额 | 退货数量 | 退货金额 | 净销售收入
	return
}

//毛利统计

func (e DataAnalysis) GrossCount(c *gin.Context) {
	s := service.Goods{}
	err := e.MakeContext(c).
		MakeOrm().
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	return
}

//毛利列表

func (e DataAnalysis) Grosslist(c *gin.Context) {
	s := service.Goods{}
	err := e.MakeContext(c).
		MakeOrm().
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	//商品-规格名称 | 销售收入 | 优惠抵扣 | 	实际销售收入 | 销售成本 | 销售毛利 | 退货金额 | 退货金额 | 销售净收入 | 毛利 | 毛利率
	return
}
