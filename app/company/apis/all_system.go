package apis

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"
	"go-admin/app/company/models"
	"go-admin/common/qiniu"
	"go-admin/common/utils"
	"go-admin/common/xlsx_export"
	"go-admin/config"
	"go-admin/global"
	"path"
	"strings"
	"time"

	models2 "go-admin/cmd/migrate/migration/models"

	"go-admin/common/business"
	"go-admin/common/dto"
	cDto "go-admin/common/dto"
	"go-admin/common/jwt/user"
	"go-admin/common/redis_worker"
)

type Worker struct {
	api.Api
}
type WorkerExportLineSummaryReq struct {
	Cycle int `json:"cycle" form:"cycle"`
	LineId int `json:"line_id" form:"line_id"`
}

type WorkerExportReq struct {
	Order   []string   `json:"order" form:"order"`
	Cycle int `json:"cycle" form:"cycle"`
	Type int `json:"type" form:"type"`
	Detail bool `json:"detail" form:"detail"` //是否单商品导出
	LineId []int `json:"line_id" form:"line_id"`
	LineName []string `json:"line_name" form:"line_name"`
}

type GetPageReq struct {
	dto.Pagination `search:"-"`
	Type string `json:"type" form:"type" search:"-" `
	BeginTime      string `form:"beginTime" search:"type:gte;column:created_at;table:company_tasks" comment:"创建时间"`
	EndTime        string `form:"endTime" search:"type:lte;column:created_at;table:company_tasks" comment:"创建时间"`

}
type OrderMapCnf struct {
	ShopId int `json:"shop_id"`
	ShopName string `json:"shop_name"`
	AllNumber int `json:"all_number"`
	GoodsName string `json:"goods_name"`
	SpecsName string `json:"specs_name"`
}
func (m *GetPageReq) GetNeedSearch() interface{} {
	return *m
}

//获取大B的下载中心任务队列
func (e *Worker)Get(c *gin.Context)  {
	req := GetPageReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	userDto, err := user.GetUserDto(e.Orm, c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}
	var data models2.CompanyTasks
	list:=make([]models2.CompanyTasks,0)
	var count int64
	orm:=e.Orm.Model(&data).Where("c_id = ?",userDto.CId).
		Scopes(
			cDto.MakeCondition(req.GetNeedSearch()),
			cDto.Paginate(req.GetPageSize(), req.GetPageIndex()),
		)

	switch req.Type {
	case "report":
		orm.Where("`type` > ?",global.ExportTypeOrder)
	case "export":
		orm.Where("`type` = ?",global.ExportTypeOrder)
	}
	err = orm.Order("id desc").
		Find(&list).Limit(-1).Offset(-1).
		Count(&count).Error
	//直接读取DB中的数据

	result:=make([]map[string]interface{},0)
	for _,row:=range list{
		result = append(result, map[string]interface{}{
			"create_time":row.CreatedAt.Format("2006-01-02 15:04:05"),
			"path":row.Path,
			"status":row.Status,
			"type":row.Type,
			"id":row.Id,
			"title":row.Title,
			"user":row.UserName,
		})
	}
	e.PageOK(result, int(count), req.GetPageIndex(), req.GetPageSize(), fmt.Sprintf("云端文件最多保留%v天",config.ExtConfig.ExportDay))
}


func (e *Worker)Download(c *gin.Context)  {
	err := e.MakeContext(c).
		MakeOrm().
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	userDto, err := user.GetUserDto(e.Orm, c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}
	uid:=c.Param("uid")
	var data models2.CompanyTasks
	e.Orm.Model(&models2.CompanyTasks{}).Where("c_id = ? and id = ?",userDto.CId,uid).Limit(1).Find(&data)
	if data.Id == 0 {
		e.Error(500, nil,"数据不存在")
		return
	}

	if data.Status == 0{
		e.Error(500, nil,"任务执行中...,请稍后下载")
		return
	}
	downloadUrl :=config.ExtConfig.CloudObsUrl + path.Join(fmt.Sprintf("%v",userDto.CId),global.CloudExportOrderFilePath,data.Path)

	e.OK(downloadUrl,"")
	return

}


func (e *Worker)CustomerBindUser(c *gin.Context) {
	req := WorkerExportLineSummaryReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	userDto, err := user.GetUserDto(e.Orm, c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}

	//获取选择周期的 所以客户和商品的对应表
	splitTableRes := business.GetTableName(userDto.CId, e.Orm)
	//先查询配送数据
	//在查询配送周期下的订单 + 路线
	var CycleCnfObj models.OrderCycleCnf
	e.Orm.Table(splitTableRes.OrderCycle).Model(
		&models.OrderCycleCnf{}).Where("id = ?",req.Cycle).Limit(1).Find(&CycleCnfObj)
	if CycleCnfObj.Id == 0 {

		e.Error(500, nil, "暂无周期")
		return
	}
	CycleUid  := CycleCnfObj.Uid

	var list []models2.Orders
	e.Orm.Table(splitTableRes.OrderTable).Select("shop_id,number,id,order_id").Where(
		"c_id = ?  and uid = ? and status in ? ",
		userDto.CId, CycleUid,global.OrderEffEct(),
	).Find(&list)



	var CustomerList []models2.Shop
	shopMap:=make(map[int]string,0)
	e.Orm.Model(&models2.Shop{}).Where("c_id = ?",userDto.CId).Select("name,id").Find(&CustomerList).Order("layer desc")

	for _,row:=range CustomerList{
		shopMap[row.Id] = row.Name
	}

	//orderGoodsMap:=make(map[int]int,0)
	orderLists:=make([]string,0)
	orderMapCnf :=make(map[string]OrderMapCnf,0)
	for _,orders:=range list{
		orderLists = append(orderLists,orders.OrderId)
		orderMapCnf[orders.OrderId] = OrderMapCnf{
			ShopId: orders.ShopId,
			AllNumber: 0,
			ShopName: shopMap[orders.ShopId],
		}
	}

	var orderSpecs []models.OrderSpecs
	e.Orm.Table(splitTableRes.OrderSpecs).Where("order_id in ? and c_id = ?",orderLists,userDto.CId).Find(&orderSpecs)

	//是每一个客户 和商品的对等关系
	allGoodsMap := make(map[string]OrderMapCnf,0)

	//商品和数量的map
	for _,orders:=range orderSpecs{


		orderShopMapCnf,ok :=orderMapCnf[orders.OrderId] //orderMapCnf 是这个商品的客户信息 还有数量
		if !ok{continue} //必须是在大的父ID 订单中

	//订单关联 用商品的规则ID做一个唯一的
		goodsKey:=fmt.Sprintf("%v/%vSRE+%v",orders.GoodsName,orders.SpecsName,orderShopMapCnf.ShopId)

		cacheDat,cacheOk:=allGoodsMap[goodsKey]
		if !cacheOk{
			cacheDat = OrderMapCnf{
				ShopId: orderShopMapCnf.ShopId,
				ShopName: orderShopMapCnf.ShopName,
				GoodsName: orders.GoodsName,
				SpecsName: orders.SpecsName,
			}
		}
		cacheDat.AllNumber += orders.Number
		allGoodsMap[goodsKey] = cacheDat

	}
	newMap:=make(map[string][]string,0)
	for _,shop :=range CustomerList{
		mathKey:=fmt.Sprintf("SRE+%v",shop.Id)
		cacheDat := make([]string,0)
		for or,orDat:=range allGoodsMap {
			var val string
			if strings.HasSuffix(or,mathKey){
				val =fmt.Sprintf("%vDEVOPS%v",orDat.ShopName,orDat.AllNumber)
			}else {
				val =fmt.Sprintf("%vDEVOPS%v",orDat.ShopName,0)
			}
		}

	}
	//
	//fmt.Println("allGoodsMap",allGoodsMap)
	//export :=xlsx_export.XlsxBaseExport{}
	////查询 所有商品 + 所有客户
	//DeliveryTime:=CycleCnfObj.DeliveryTime.Format(time.DateOnly)
	//xlsxPath,_ := export.CustomerBindUser(userDto.CId,DeliveryTime,allGoodsMap)
	//
	//downloadUrl :=path.Join("/company/api/v1/report/file",xlsxPath)
	//e.OK(downloadUrl,"")
	e.OK("","")
	return

}

func (e *Worker)LineSummary(c *gin.Context) {
	req := WorkerExportLineSummaryReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	userDto, err := user.GetUserDto(e.Orm, c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}

	splitTableRes := business.GetTableName(userDto.CId, e.Orm)
	//先查询配送数据
	//在查询配送周期下的订单 + 路线
	var CycleCnfObj models.OrderCycleCnf
	e.Orm.Table(splitTableRes.OrderCycle).Model(
		&models.OrderCycleCnf{}).Where("id = ?",req.Cycle).Limit(1).Find(&CycleCnfObj)
	if CycleCnfObj.Id == 0 {

		e.Error(500, nil, "暂无周期")
		return
	}
	CycleUid  := CycleCnfObj.Uid

	var list []models2.Orders
	e.Orm.Table(splitTableRes.OrderTable).Select("shop_id,number,order_money").Where(
		"c_id = ? and line_id = ? and uid = ? and status in ? ",
		userDto.CId, req.LineId,CycleUid,global.OrderEffEct(),
		).Find(&list)

	var lineObject models2.Line
	e.Orm.Model(&lineObject).Where("c_id = ? and id = ?",userDto.CId,req.LineId).Limit(1).Find(&lineObject)
	cacheShopMap:=make(map[int]xlsx_export.LineSummaryRow,0)
	for _,row:=range list{
		shopData,ok:=cacheShopMap[row.ShopId]
		if !ok{
			var shopObject models2.Shop
			e.Orm.Model(&shopObject).Where("c_id = ? and id = ?",userDto.CId,row.ShopId).Limit(1).Find(&shopObject)
			if shopObject.Id == 0 { //没有客户就退出
				continue
			}

			shopData = xlsx_export.LineSummaryRow{
				Layer: shopObject.Layer,
				ShopId: row.ShopId,
				ShopName: shopObject.Name,
				ShopAddress: shopObject.Address,
				ShopPhone:shopObject.Phone,
			}
		}
		shopData.OrderCount +=row.Number
		shopData.OrderMoney +=utils.RoundDecimalFlot64(row.OrderMoney)
		cacheShopMap[row.ShopId] = shopData
	}
	fmt.Println("cacheShopMap",cacheShopMap)
	//sort.Slice(cacheShopMap, func(i, j int) bool {
	//	fmt.Println("cacheShopMap[i]",cacheShopMap[i])
	//	return cacheShopMap[i].Layer > cacheShopMap[j].Layer
	//})
	export :=xlsx_export.XlsxBaseExport{}
	xlsxPath,_ := export.ExportLineSummary(userDto.CId,lineObject.Name,cacheShopMap)

	downloadUrl :=path.Join("/company/api/v1/report/file",xlsxPath)
	e.OK(downloadUrl,"")
	return

}
func (e *Worker)Create(c *gin.Context)  {
	req:=WorkerExportReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	if req.Type == 0 {
		if len(req.Order) == 0 {
			e.Error(500, nil, "请选择订单ID")
			return
		}
	}
	userDto, err := user.GetUserDto(e.Orm, c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}

	splitTableRes := business.GetTableName(userDto.CId, e.Orm)
	var CycleUid string //配送周期的UID
	if req.Type != 0  {
		var CycleCnfObj models.OrderCycleCnf
		e.Orm.Table(splitTableRes.OrderCycle).Select("uid,id").Model(
			&models.OrderCycleCnf{}).Where("id = ?",req.Cycle).Limit(1).Find(&CycleCnfObj)
		if CycleCnfObj.Id == 0 {

			e.Error(500, nil, "暂无周期")
			return
		}
		CycleUid  = CycleCnfObj.Uid
	}

	//查询队列限制
	CompanyCnf := business.GetCompanyCnf(userDto.CId, "export_worker", e.Orm)
	MaxNumber := CompanyCnf["export_worker"]
	//获取订单的导出队列
	thisWorker:=redis_worker.GetExportQueueLength(userDto.CId,global.WorkerOrderStartName)
	if thisWorker >= MaxNumber {
		e.Error(500, nil,fmt.Sprintf("最多同时支持%v个任务执行,请稍后重试",MaxNumber))
		return
	}
	//用一个key来做防止大量的重复提交,给DB带来非必要的查询
	//如果是订单 那就查询 当前时间戳  精确到分 ,查询到这个key 已经存在 + 处理中   = 返回请一分钟后重试

	//如果是汇总 查询cid:周期ID
	var Queue string
	var title string
	var LineExport int
	//限制一次
	mathKey := time.Now().Add(-10 * time.Second).Format("200601021504")

	switch req.Type {

	case global.ExportTypeOrder:
		if req.Detail {
			title = "配送订单导出"
			mathKey = fmt.Sprintf("%v_detail_order",mathKey)

		}else {
			title = "批量导出配送订单"
			mathKey = fmt.Sprintf("%v_order",mathKey)
		}
		Queue = global.WorkerOrderStartName
	case global.ExportTypeSummary:
		title = "导出配送汇总表"
		Queue = global.WorkerReportSummaryStartName
		mathKey = fmt.Sprintf("%v_summary",mathKey)
	case global.ExportTypeLine:
		if len(req.LineName) > 1{
			title = fmt.Sprintf("批量导出【%v】条路线汇总表",len(req.LineName))
		}else {
			title = fmt.Sprintf("导出【%v】路线汇总表",req.LineName[0])
		}
		LineExport = 0
		Queue = global.WorkerReportLineStartName
		mathKey = fmt.Sprintf("%v_line",mathKey)
	case global.ExportTypeLineShopDelivery:
		if len(req.LineName) > 1{
			title = fmt.Sprintf("批量导出【%v】条路线明细表",len(req.LineName))
		}else {
			title = fmt.Sprintf("导出【%v】路线明细表",req.LineName[0])
		}
		LineExport = 1
		Queue = global.WorkerReportLineDeliveryStartName
		mathKey = fmt.Sprintf("%v_delivery",mathKey)
	}
	var count int64

	e.Orm.Model(&models2.CompanyTasks{}).Where("`key` = ? and c_id = ? and type = ?",mathKey,userDto.CId,req.Type).Count(&count)

	if count > 0 {
		e.Error(500, nil,"请勿在10秒内重复提交相同任务")
		return
	}

	taskTable:=models2.CompanyTasks{
		UserName: userDto.Username,
		Title: title,
		CreateBy: userDto.UserId,
		CId:      userDto.CId,
		Type:     req.Type,
		Key:      mathKey,
		Status:   0,
	}
	//记录在DB中
	e.Orm.Create(&taskTable)
	exportReq :=global.ExportRedisInfo{
		CId: userDto.CId,
		Order: req.Order,
		Cycle: req.Cycle,
		CycleUid: CycleUid,
		OrmId: taskTable.Id,
		LineId: req.LineId,
		ExportUser: userDto.Username,
		ExportTime: time.Now().Format("2006-01-02 15:04:05"),
		Queue: Queue,
		LineExport: LineExport,
	}

	//先发送到redis中
	_=redis_worker.SendExportQueue(exportReq)
	e.OK("","任务创建成功.请前往任务中心查看 ！")
	return
}


func (e *Worker)Remove(c *gin.Context) {
	err := e.MakeContext(c).
		MakeOrm().
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	userDto, err := user.GetUserDto(e.Orm, c)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	}
	uid:=c.Param("uid")
	var data models2.CompanyTasks
	e.Orm.Model(&models2.CompanyTasks{}).Where("c_id = ? and id = ?",userDto.CId,uid).Limit(1).Find(&data)
	if data.Id == 0 {
		e.Error(500, nil,"数据不存在")
		return
	}

	e.Orm.Model(&data).Where("id = ?",data.Id).Delete(&data)
	//
	obsUrl :=path.Join(fmt.Sprintf("%v",userDto.CId),global.CloudExportOrderFilePath,data.Path)

	//数据删除
	buckClient:=qiniu.QinUi{
		CId: userDto.CId,
	}
	buckClient.InitClient()

	buckClient.RemoveFile(obsUrl)
	e.OK("","删除成功")
	return

}