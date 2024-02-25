/**
@Author: chaoqun
* @Date: 2023/12/22 11:07
*/
package redis_worker

import (
	"encoding/json"
	"fmt"
	"github.com/go-admin-team/go-admin-core/sdk"
	"go-admin/cmd/migrate/migration/models"
	"go-admin/common/redis_db"
	"go-admin/common/xlsx_export"
	"go-admin/global"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"math/rand"
	"time"
)

var (
	WorkerSleep = 10000 * time.Millisecond //默认睡10秒

)

type LoopRedisWorker struct {
	Orm *gorm.DB
}
func (l *LoopRedisWorker)Start()  {
	//任务开启 必须睡眠时间 等待orm加载完毕
	time.Sleep(7*time.Second)
	l.Orm = sdk.Runtime.GetDbByKey("*")
	fmt.Println("异步任务初始化成功！！！")

	for {
		//读取不同的任务Queue
		for _,queueName:=range global.QueueGroup{
			//fmt.Println("开始巡检Queue",queueName)

			//随机睡10S以内的数据
			randomSleepTime := time.Duration(rand.Intn(10000)) * time.Millisecond
			time.Sleep(WorkerSleep + randomSleepTime) //

			redis_db.RedisCli.Do(global.RedisCtx, "select", global.AllQueueChannel)
			//获取所以key
			keys,err:=redis_db.RedisCli.Keys(global.RedisCtx,fmt.Sprintf("%v*",queueName)).Result()
			if err!=nil{
				zap.S().Errorf("读取redis 获取key* 数据失败,key:%v,错误:%v",queueName,err)
				continue
			}
			if len(keys) == 0 {
				continue
			}
			//所有的大B -Key数据
			fmt.Println("queueName",queueName,"keys",keys)
			for _,key:=range keys{
				data,keyErr:=redis_db.RedisCli.LRange(global.RedisCtx,key,0,-1).Result()

				if keyErr!=nil{
					zap.S().Errorf("读取redis数据失败,key:%v,错误:%v",keys,keyErr)
					continue
				}
				switch queueName {
				case global.WorkerOrderStartName: //订单选中导出
					go l.GetExportQueueInfo(key,data)
				case global.WorkerReportSummaryStartName: //汇总导出
					go l.GetSummaryExportQueueInfo(key,data)
				case global.WorkerReportLineStartName: //路线导出
					go l.GetReportLineExportQueueInfo(key,data)

				case global.WorkerReportLineDeliveryStartName: //路线配送表导出

					go l.GetReportLineDeliveryExportQueueInfo(key,data)
				}

			}
		}


	}
}

//查询DB是否已经执行成功,防止未删除导致的异常
func(l *LoopRedisWorker) IsValidWorkTask(cid,ormId int,key string,index int) bool  {
	var task models.CompanyTasks
	if l.Orm == nil{
		l.Orm = sdk.Runtime.GetDbByKey("*")
		return false
	}
	l.Orm.Model(&models.CompanyTasks{}).Select("id,status").Where("id = ? and c_id = ?",
		ormId,cid).Limit(1).Find(&task)
	if task.Id == 0 {
		xlsx_export.EmptyKey(key,index)
		return false
	}
	//只有执行中 才可以启动任务
	if task.Status == 0{
		return true
	}
	xlsx_export.EmptyKey(key,index)
	return false
}
//获取到消息了
//开始解析
func (l *LoopRedisWorker)GetExportQueueInfo(key string,data []string)   {
	for index,dat:=range data{

		time.Sleep(600*time.Millisecond)
		var err error
		row:=global.ExportRedisInfo{}
		err =json.Unmarshal([]byte(dat),&row)
		if err!=nil{
			continue
		}
		//DB中状态为0
		if !l.IsValidWorkTask(row.CId,row.OrmId,key,index){
			continue
		}
		orderExportFunc:=xlsx_export.OrderExportObj{
			RedisKey: key,
			Dat: row,
			Orm:l.Orm,
			UpCloud: true,
		}
		if orderExportFunc.Orm == nil{
			zap.S().Errorf("读取redis 导出[配送订单任务数据] Orm对象为空")
			continue
		}
		successTag:=true
		errorMsg :=""
		sheetData := make(map[int]*xlsx_export.SheetRow,0)
		if sheetData,err =orderExportFunc.ReadOrderDetail();err!=nil{
			successTag =false
			errorMsg = err.Error()
			zap.S().Errorf("读取redis 导出[配送订单任务数据]ReadOrderDetail,错误:%v",err)
			continue

		}
		//保存到云端
		orderExportFunc.SaveExportXlsx(row,sheetData)

		if err = xlsx_export.SaveExportDb(
			orderExportFunc.Dat.OrmId,orderExportFunc.Dat.CId,
			orderExportFunc.FileName,
			successTag,errorMsg,orderExportFunc.Orm);err!=nil{
			zap.S().Errorf("读取redis 导出[配送订单任务数据]SaveExportDb,错误:%v",err)
			continue
		}
		//最后在删除
		xlsx_export.EmptyKey(key,index)
	}

}

func (l *LoopRedisWorker)GetSummaryExportQueueInfo(key string,data []string)   {
	for index,dat:=range data{

		time.Sleep(700*time.Millisecond)
		var err error
		row:=global.ExportRedisInfo{}
		err =json.Unmarshal([]byte(dat),&row)
		if err!=nil{
			continue
		}

		exportFunc:=xlsx_export.SummaryExportObj{
			RedisKey: key,
			Dat: row,
			Orm:l.Orm,
			UpCloud: true,
		}
		if exportFunc.Orm == nil{
			zap.S().Errorf("读取redis 导出[汇总表任务数据] Orm对象为空")
			continue
		}
		//DB中状态不为0
		if !l.IsValidWorkTask(row.CId,row.OrmId,key,index){
			continue
		}
		successTag:=true
		errorMsg :=""
		sheetData := &xlsx_export.SheetRow{}

		if sheetData,err =exportFunc.ReadSummaryDetail();err!=nil{
			successTag =false
			errorMsg = err.Error()
			zap.S().Errorf("读取redis 导出[汇总表任务数据] ReadOrderDetail,错误:%v",err)
			continue
		}
		//保存到云端
		exportFunc.SaveExportXlsx(row,sheetData)

		//数据保存
		if err = xlsx_export.SaveExportDb(
			exportFunc.Dat.OrmId,exportFunc.Dat.CId,
			exportFunc.FileName,
			successTag,errorMsg,exportFunc.Orm);err!=nil{
			zap.S().Errorf("读取redis 导出[配送订单任务数据]SaveExportDb,错误:%v",err)
			continue
		}
		//最后在删除redis

		xlsx_export.EmptyKey(key,index)
	}


}



func (l *LoopRedisWorker)GetReportLineExportQueueInfo(key string,data []string) {
	fmt.Println("GetReportLineExportQueueInfo",key)
	//fmt.Println("开始路线数据导出",key,"DATA",data)

	for index,dat:=range data{
		time.Sleep(800*time.Millisecond)
		var err error
		row:=global.ExportRedisInfo{}
		err =json.Unmarshal([]byte(dat),&row)
		if err!=nil{
			continue
		}
		orderExportFunc:=xlsx_export.ReportLineObj{
			RedisKey: key,
			Dat: row,
			Orm:l.Orm,
			UpCloud: true,
			CycleUid: row.CycleUid,
		}
		if orderExportFunc.Orm == nil{
			zap.S().Errorf("读取redis 导出[路线数据] Orm对象为空")
			continue
		}
		//DB中状态为0
		if !l.IsValidWorkTask(row.CId,row.OrmId,key,index){
			continue
		}
		successTag:=true
		errorMsg :=""
		//兼容多条路线同时存放

		sheetData,detailErr :=orderExportFunc.ReadLineDetail()
		if detailErr!=nil{
			successTag =false
			errorMsg = detailErr.Error()
			zap.S().Errorf("读取redis 导出[路线数据] ReadOrderDetail,错误:%v",detailErr)
			continue

		}

		//保存到云端
		zipFile:=fmt.Sprintf("%v 多路线表导出.zip",row.ExportTime)
		FileName :=xlsx_export.SaveLineExportXlsx("line",zipFile,orderExportFunc.UpCloud,row,sheetData)

		if err = xlsx_export.SaveExportDb(
			orderExportFunc.Dat.OrmId,orderExportFunc.Dat.CId,
			FileName,
			successTag,errorMsg,orderExportFunc.Orm);err!=nil{
			zap.S().Errorf("读取redis 导出[路线数据] SaveExportDb,错误:%v",err)
			continue
		}
		//最后在删除
		xlsx_export.EmptyKey(key,index)
	}

}
func (l *LoopRedisWorker)GetReportLineDeliveryExportQueueInfo(key string,data []string) {

	//fmt.Println("GetReportLineDeliveryExportQueueInfo",key)
	//查询线路下不同的小B列表
	//先查一波路线 -> 路线下的小B
	//不同的是 路线是一个大的点，如果多个路线 那就是多个文件了里面有多个小B


	//同时要支持多个文件导出的逻辑
	//1.如果是单个路线 那就是一个excel
	//2.如果是多个路线 那就是一个zip压缩包
	for index,dat:=range data{
		time.Sleep(900*time.Millisecond)
		var err error
		row:=global.ExportRedisInfo{}
		err =json.Unmarshal([]byte(dat),&row)
		if err!=nil{
			continue
		}
		//DB中状态为0
		if !l.IsValidWorkTask(row.CId,row.OrmId,key,index){
			continue
		}
		orderExportFunc:=xlsx_export.ReportDeliveryLineObj{
			RedisKey: key,
			Dat: row,
			Orm:l.Orm,
			UpCloud: true,
			CycleUid: row.CycleUid,
		}
		if orderExportFunc.Orm == nil{
			zap.S().Errorf("读取redis 导出[路线配送表] Orm对象为空")
			continue
		}
		successTag:=true
		errorMsg :=""


		sheetData,detailErr :=orderExportFunc.ReadLineDeliveryDetail()
		if detailErr!=nil{
			successTag =false
			errorMsg = detailErr.Error()
			zap.S().Errorf("读取redis 导出[路线配送表] ReadOrderDetail,错误:%v",detailErr)
			continue

		}
		//
		//fmt.Println("successTag,",successTag,errorMsg)
		//for _,k:=range sheetData{
		//	for _,v:=range k.DeliveryData{
		//
		//		for _,s:=range v{
		//			fmt.Println("线路",k.LineName,"商家",s.SheetName,"数据", len(s.Table))
		//			for _,b:=range s.Table{
		//				fmt.Printf("--row:%v\n",b.GoodsName)
		//			}
		//		}
		//	}
		//}
		zipFile:=fmt.Sprintf("%v 多路线配送表导出.zip",row.ExportTime)
		FileName :=xlsx_export.SaveLineExportXlsx("delivery",zipFile,orderExportFunc.UpCloud,row,sheetData)

		if err = xlsx_export.SaveExportDb(
			orderExportFunc.Dat.OrmId,orderExportFunc.Dat.CId,
			FileName,
			successTag,errorMsg,orderExportFunc.Orm);err!=nil{
			zap.S().Errorf("读取redis 导出[路线配送表] SaveExportDb,错误:%v",err)
			continue
		}
		//最后在删除
		xlsx_export.EmptyKey(key,index)
	}
}