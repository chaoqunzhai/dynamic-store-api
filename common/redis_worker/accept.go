/**
@Author: chaoqun
* @Date: 2023/12/22 11:07
*/
package redis_worker

import (
	"encoding/json"
	"fmt"
	"github.com/go-admin-team/go-admin-core/sdk"
	"go-admin/common/redis_db"
	"go-admin/common/xlsx_export"
	"go-admin/global"
	"go.uber.org/zap"
	"math/rand"

	"time"
)

func LoopRedisWorker()  {
	fmt.Println("异步导出任务启动成功！！！")
	for {
		//读取不同的任务Queue
		for _,queueName:=range global.QueueGroup{

			randomSleepTime := time.Duration(rand.Intn(10)+1) * time.Second
			time.Sleep(12 * time.Second + randomSleepTime) //10秒才进行任务处理
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
			for _,key:=range keys{
				data,keyErr:=redis_db.RedisCli.LRange(global.RedisCtx,key,0,-1).Result()

				if keyErr!=nil{
					zap.S().Errorf("读取redis数据失败,key:%v,错误:%v",keys,keyErr)
					continue
				}
				switch queueName {
				case global.WorkerOrderStartName: //订单选中导出
					GetExportQueueInfo(key,data)
				case global.WorkerReportSummaryStartName: //汇总导出

					GetSummaryExportQueueInfo(key,data)
				}

			}
		}


	}
}
//获取到消息了
//开始解析
func GetExportQueueInfo(key string,data []string)   {
	for _,dat:=range data{
		//睡眠500毫秒,缓解压力
		time.Sleep(500*time.Millisecond)
		var err error
		row:=global.ExportRedisInfo{}
		err =json.Unmarshal([]byte(dat),&row)
		if err!=nil{
			continue
		}
		orderExportFunc:=xlsx_export.OrderExportObj{
			RedisKey: key,
			Dat: row,
			Orm:sdk.Runtime.GetDbByKey("*"),
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
		xlsx_export.EmptyKey(key,len(dat))
	}

}

func GetSummaryExportQueueInfo(key string,data []string)   {
	for _,dat:=range data{
		//睡眠500毫秒,缓解压力
		time.Sleep(500*time.Millisecond)
		var err error
		row:=global.ExportRedisInfo{}
		err =json.Unmarshal([]byte(dat),&row)
		if err!=nil{
			continue
		}
		exportFunc:=xlsx_export.SummaryExportObj{
			RedisKey: key,
			Dat: row,
			Orm:sdk.Runtime.GetDbByKey("*"),
			UpCloud: true,
		}
		if exportFunc.Orm == nil{
			zap.S().Errorf("读取redis 导出[汇总表任务数据] Orm对象为空")
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
		xlsx_export.EmptyKey(key,len(dat))
	}


}