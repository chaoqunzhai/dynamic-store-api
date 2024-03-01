/**
@Author: chaoqun
* @Date: 2023/12/22 11:06
*/
package redis_worker

import (
	"encoding/json"
	"fmt"
	"go-admin/common/redis_db"
	"go-admin/global"
	"go.uber.org/zap"
)

func getExportWorkerName(cid int,q string) string  {

	return 	fmt.Sprintf("%v.export.%v",q,cid)
}
//发送数据到队列中
func SendExportQueue(r global.ExportRedisInfo)  error {
	zap.S().Debugf("%v 报表导出任务,发送至redis:%v 开始",r.CId,global.AllQueueChannel,)
	data,err :=json.Marshal(r)
	if err!=nil{
		zap.S().Errorf("%v 报表导出任务,发送至redis:%v 失败:%v",r.CId,global.AllQueueChannel,err.Error())
		return err
	}

	redis_db.RedisCli.Do(global.RedisCtx, "select", global.AllQueueChannel)
	redis_db.RedisCli.LPush(global.RedisCtx,getExportWorkerName(r.CId,r.Queue),string(data))
	zap.S().Debugf("%v 报表导出任务,发送至redis:%v 成功",r.CId,global.AllQueueChannel)
	return nil
}

func GetExportQueueLength(cid int,q string) int {
	redis_db.RedisCli.Do(global.RedisCtx, "select", global.AllQueueChannel)
	length, err :=redis_db.RedisCli.LLen(global.RedisCtx,getExportWorkerName(cid,q)).Result()
	if err!=nil{
		return 0
	}
	return  int(length)
}