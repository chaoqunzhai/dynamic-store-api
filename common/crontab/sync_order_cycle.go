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
//获取所有的分表配置

func (s *ThisDaySyncOrderCycle)GetCompanySplit()  {

	s.CompanyMap = make(map[int]*models2.SplitTableMap,0)//方便取

	s.NowDay = time.Now().Format("2006-01-02") //查询当天

	var SplitTableMap []*models2.SplitTableMap
	s.Orm.Model(&models2.SplitTableMap{}).Find(&SplitTableMap)

	for _,row:=range SplitTableMap{

		s.CompanyMap[row.Id] = row
	}

	//构建周期映射MAP 方便后面订单更新
	s.MakeOrderCycleCnfTable()

	//订单按量分隔 启动协程 处理
	s.RunUpdateTableStatus()
}
//获取所有周期配置表

func (s *ThisDaySyncOrderCycle)MakeOrderCycleCnfTable()  {


	for cid, splitCnf := range s.CompanyMap { //循环的是每一个大B的周期配送表
		findCycleCnfList := make([]string,0) //
		var orderCycleList []models.OrderCycleCnf
		s.Orm.Table(splitCnf.OrderCycle).Select("uid").Where("delivery_time = ? and c_id = ?",s.NowDay,cid).Find(&orderCycleList)

		for _,row:=range orderCycleList{
			findCycleCnfList = append(findCycleCnfList,row.Uid) //给订单用于修改状态
		}
		splitCnf.CycleCnfList = findCycleCnfList //保存起来 订单需要用
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
			end = len(orderTables) - 1 // 如果剩余元素不足MaxCompany个，取剩余所有元素
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
	fmt.Println("Processing array:", tableList)

	for _,splitCnf :=range tableList{

		for _,uid:=range splitCnf.CycleCnfList {
			s.Orm.Table(splitCnf.OrderTable).Where("c_id = ? and status = ? and uid = ?",splitCnf.CId,global.OrderStatusWaitSend,uid).Updates(map[string]interface{}{
				"status":global.OrderWaitConfirm,
			})

		}
	}

}
