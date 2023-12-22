/**
@Author: chaoqun
* @Date: 2023/12/22 11:07
*/
package redis_worker

import (
	"encoding/json"
	"fmt"
	"github.com/go-admin-team/go-admin-core/sdk"
	"github.com/xuri/excelize/v2"
	"go-admin/cmd/migrate/migration/models"
	"go-admin/common/redis_db"
	"go-admin/global"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

type ExportObj struct {
	Orm *gorm.DB
	Dat ExportReq
	RedisKey string
	FileName string
}

func LoopRedisWorker()  {
	fmt.Println("异步导出任务启动成功！！！")
	for {
		time.Sleep(10 * time.Second) //10秒才进行任务处理
		redis_db.RedisCli.Do(RedisCtx, "select", global.AllQueueChannel)
		//获取所以key
		keys,err:=redis_db.RedisCli.Keys(RedisCtx,fmt.Sprintf("%v*",WorkerStartName)).Result()
		if err!=nil{
			zap.S().Errorf("读取redis 获取key* 数据失败,key:%v,错误:%v",WorkerStartName,err)
			continue
		}
		//所有的大B -Key数据
		for _,key:=range keys{
			data,keyErr:=redis_db.RedisCli.LRange(RedisCtx,key,0,-1).Result()

			if keyErr!=nil{
				zap.S().Errorf("读取redis数据失败,key:%v,错误:%v",keys,keyErr)
				continue
			}
			GetExportQueueInfo(key,data)
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
		row:=ExportReq{}
		err =json.Unmarshal([]byte(dat),&row)
		if err!=nil{
			continue
		}
		zipFunc:=ExportObj{
			RedisKey: key,
			Dat: row,
			Orm:sdk.Runtime.GetDbByKey("*"),
		}
		if zipFunc.Orm == nil{
			zap.S().Errorf("读取redis 解析导出任务数据 获取Orm对象为空")
			continue
		}
		successTag:=true
		errorMsg :=""
		//if err =zipFunc.ReadOrderDetail();err!=nil{
		//	successTag =false
		//	errorMsg = err.Error()
		//	zap.S().Errorf("读取redis 解析导出任务数据 ReadOrderDetail,错误:%v",err)
		//
		//}
		if err = zipFunc.SaveExportZIP();err!=nil{
			successTag =false
			errorMsg = err.Error()
			zap.S().Errorf("读取redis 解析导出任务数据 SaveExportZIP,错误:%v",err)
		}

		if err = zipFunc.SaveExportDb(successTag,errorMsg);err!=nil{
			zap.S().Errorf("读取redis 解析导出任务数据 SaveExportDb,错误:%v",err)
			continue
		}
		zipFunc.EmptyKey(len(dat))
	}

}

//开始执行数据导出
//多个订单save为一个excel文件
func (e *ExportObj)ReadOrderDetail() error  {
	file := excelize.NewFile()
	//设置表名
	file.SetSheetName("Sheet1", "订单列表")
	//创建流式写入
	writer, err := file.NewStreamWriter("订单列表")
	//修改列宽
	writer.SetColWidth(1, 20, 18)
	//设置表头
	writer.SetRow("A1", []interface{}{"商品名称", "商品规格", "商品数量", "商品单价", "客户名称", "提交时间","订单状态"})
	if err != nil {
		return err
	}
	for index,order:=range e.Dat.Order{
		fmt.Println("order",order)

		cell, _ := excelize.CoordinatesToCellName(1, index+1)
		//添加的数据
		writer.SetRow(cell, []interface{}{"商品名称", "商品规格", "商品数量", "商品单价", "客户名称", "提交时间","订单状态"})

	}
	//结束流式写入
	writer.Flush()
	xlsxName:=fmt.Sprintf("%v-订单数据.xlsx",time.Now().Format("20060102-150405"))
	file.SaveAs(xlsxName)
	//保存的订单数据,存对象存储中即可
	//保存成功后 上传到对象存储中

	e.FileName = xlsxName
	return nil
}


//导出数据保存在本地zip压缩包
func (e *ExportObj)SaveExportZIP() error {
	fmt.Println("保存到本地zip文件",e.Dat)
	return nil
}

//更新table中状态
func (e *ExportObj)SaveExportDb(successTag bool,msg string)  error{
	var status int
	if !successTag{
		status = 2
	}else {
		status = 1
	}
	if msg !=""{
		if len(msg) > 60{
			msg = msg[:60]
		}
	}
	e.Orm.Model(&models.CompanyTasks{}).Where("id = ? and c_id = ?",
		e.Dat.OrmId,e.Dat.CId).Updates(map[string]interface{}{
		"status":status,
		"path":e.FileName,
		"msg":msg,
	})
	return nil

}
//如果key下的list位空 ,那就支持清空这个key
func  (e *ExportObj)EmptyKey(keyLen int) {
	err :=redis_db.RedisCli.LTrim(RedisCtx,e.RedisKey,1,int64(keyLen)).Err()
	if err!=nil{
		zap.S().Errorf("清理redis key:%v 数据清理失败:%v",e.RedisKey,err)
	}else {
		zap.S().Infof("redis key:%v 消费完毕,数据清理成功",e.RedisKey)
	}
}