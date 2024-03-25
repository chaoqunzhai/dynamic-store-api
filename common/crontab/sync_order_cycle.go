package crontab

import (
	"fmt"
	"go-admin/app/company/models"
	models2 "go-admin/cmd/migrate/migration/models"
	"go-admin/global"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"sync"
	"time"
)

//定时获取 大B order_cycle_cnf 周期配置表中
//1.查询delivery_time  = 今天的日期  得到获取到周期的uid
//2.周期ID在订单表中获取 状态为待配送的订单,更改为配送中

var (
	MaxCompany = 50 //同时处理50个大B客户数据
)

type CycleCnf struct {
	CId int `json:"c_id"`
	Uid string `json:"uid"`
}
type ThisDaySyncOrderCycle struct {
	Orm *gorm.DB
	NowDay string
	CompanyMap map[int]*models2.SplitTableMap


}

func (s *ThisDaySyncOrderCycle)WhereStatus(IsOpenApprove bool) string  {
	var WhereOrderStatus string
	if IsOpenApprove{ //,如果开启了 只有审批通过的订单 才会变为配送中
		WhereOrderStatus = fmt.Sprintf("status = %v and approve_status = %v",global.OrderStatusWaitSend,global.OrderApproveOk)
	}else { //没有开启审批,那就只查询 待配送的
		WhereOrderStatus = fmt.Sprintf("status = %v",global.OrderStatusWaitSend)

	}
	return  WhereOrderStatus
}
//获取所有的分表配置

func (s *ThisDaySyncOrderCycle)RunCompanySplitOrderSync()  {
	//获取是否开启了审批权限


	s.CompanyMap = make(map[int]*models2.SplitTableMap,0)//方便取

	s.NowDay = time.Now().Format("2006-01-02") //查询当天

	var SplitTableMap []*models2.SplitTableMap
	s.Orm.Model(&models2.SplitTableMap{}).Find(&SplitTableMap)

	for _,row:=range SplitTableMap{
		//查询大B配置的订单验收天数

		var OrderTrade models2.OrderTrade
		s.Orm.Model(&models2.OrderTrade{}).Select("id,receive_days").Where("c_id = ?",row.CId).Limit(1).Find(&OrderTrade)
		var ReceiveDays int
		if OrderTrade.Id == 0 {
			ReceiveDays = global.OrderReceiveDays //默认天数
		}else {
			ReceiveDays = OrderTrade.ReceiveDays //配置的天数
		}
		//更新配送中订单 设置的数天后自动确认收货, 也就是更新为OrderStatusOver

		row.ReceiveBeforeTime = s.GetReceiveBeforeDaysTime(ReceiveDays)
		s.CompanyMap[row.CId] = row
	}

	//构建周期映射MAP 方便后面订单更新
	//s.MakeOrderCycleCnfTable()

	//订单按量分隔 启动协程 处理
	//s.RunUpdateTableStatus()

	//最后进行一次验收
	s.OverStatus()
}

func (s *ThisDaySyncOrderCycle)GetReceiveBeforeDaysTime(ReceiveDays int) string  {

	//往后推5分钟
	return time.Now().Add(5 * time.Minute).AddDate(0,0,-ReceiveDays).Format("2006-01-02 15:04:05")

}

//只进行完结验证

func (s *ThisDaySyncOrderCycle)OverStatus()  {
	//最后进行一次 验收检测
	for cid,splitCnf:=range s.CompanyMap{
		//查询必须是配送中 +  小于 验收时间的订单
		s.Orm.Table(splitCnf.OrderTable).Where("c_id = ? and `status` = ? and delivery_run_at <= ?",cid,global.OrderWaitConfirm,splitCnf.ReceiveBeforeTime).Updates(map[string]interface{}{
			"status":global.OrderStatusOver,
		})
	}

}
//获取所有周期配置表

func (s *ThisDaySyncOrderCycle)MakeOrderCycleCnfTable()  {


	for cid, splitCnf := range s.CompanyMap { //循环的是每一个大B的周期明细表
		findCycleCnfList := make([]string,0) //
		var orderCycleList []models.OrderCycleCnf
		if !s.Orm.Table(splitCnf.OrderCycle).Migrator().HasTable(splitCnf.OrderCycle){
			fmt.Println("表",splitCnf.CId,splitCnf.OrderCycle,"不存在")
			//如果有不存在表的情况,那就是默认表
			splitCnf = &models2.SplitTableMap{
				CId: splitCnf.CId,
				OrderTable: global.SplitOrderDefaultTableName,
				OrderSpecs: global.SplitOrderDefaultSubTableName,
				OrderCycle: global.SplitOrderCycleSubTableName,
				OrderEdit:global.SplitOrderEdit,
				OrderReturn: global.SplitOrderReturn,
				InventoryRecordLog:global.InventoryRecordLog,
			}
		}

		//查询今天的订单订单
		s.Orm.Table(splitCnf.OrderCycle).Select("uid").Where("delivery_time = ? and c_id = ?",s.NowDay,cid).Find(&orderCycleList)

		var Approve models2.OrderApproveCnf
		s.Orm.Model(&models2.OrderApproveCnf{}).Select("enable").Where("c_id = ?",cid).Limit(1).Find(&Approve)

		splitCnf.IsOpenApprove = Approve.Enable //订单审批赋值

		for _,row:=range orderCycleList{
			findCycleCnfList = append(findCycleCnfList,row.Uid) //给订单用于修改状态
		}
		splitCnf.CycleCnfList = findCycleCnfList //保存起来 订单需要根据UID进行查询
		s.CompanyMap[cid] = splitCnf
	}
}
//获取所有订单表,因为订单表数据量大,是需要启动多协程处理

func (s *ThisDaySyncOrderCycle)RunUpdateTableStatus()  {
	orderTables :=make([]*models2.SplitTableMap,0)

	for _,row:=range s.CompanyMap{
		orderTables = append(orderTables,row)
	}
	wg := &sync.WaitGroup{}
	runStart:=time.Now()
	zap.S().Info("定时同步 修改周期订单配送状态 开始...")
	for i := 0; i < len(orderTables); i += MaxCompany {
		end := i + MaxCompany // 取前MaxCompany个元素，因为数组索引从0开始
		if end >= len(orderTables) {
			end = len(orderTables) // 如果剩余元素不足MaxCompany个，取剩余所有元素
		}

		wg.Add(1) // 增加等待组计数器
		go s.processWaitToConfirm(orderTables[i:end], wg) // 开启协程处理指定范围的数组元素
	}

	wg.Wait() // 等待所有协程完成
	zap.S().Infof("定时同步 修改周期订单配送状态 完成,耗时%v!",time.Since(runStart).Seconds())

}

//批量更新订单 待配送为 配送中
func (s *ThisDaySyncOrderCycle)processWaitToConfirm(tableList []*models2.SplitTableMap, wg *sync.WaitGroup)  {
	defer wg.Done()
	// 处理数组的逻辑


	for _,splitCnf :=range tableList{

		whereStatus:=s.WhereStatus(splitCnf.IsOpenApprove)

		for _,uid:=range splitCnf.CycleCnfList {
			s.Orm.Table(splitCnf.OrderTable).Where(whereStatus).Where("c_id = ? and uid = ? ",splitCnf.CId,uid).Updates(map[string]interface{}{
				"status":global.OrderWaitConfirm, //待配送的订单更新为配送中
				"approve_status":global.OrderApproveOk, //也需要改为审批通过,防止后期开启了审核,订单都为待审核状态
				"delivery_run_at":time.Now(),
			})

		}
	}

}
