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
)

func getExportWorkerName(cid int) string  {

	return 	fmt.Sprintf("%v.export.%v",WorkerStartName,cid)
}
//发送数据到队列中
func SendExportQueue(r ExportReq)  {

	data,err :=json.Marshal(r)
	if err!=nil{
		return
	}

	redis_db.RedisCli.Do(RedisCtx, "select", global.AllQueueChannel)
	redis_db.RedisCli.LPush(RedisCtx,getExportWorkerName(r.CId),string(data))
}

func GetExportQueueLength(cid int) int {
	redis_db.RedisCli.Do(RedisCtx, "select", global.AllQueueChannel)
	length, err :=redis_db.RedisCli.LLen(RedisCtx,getExportWorkerName(cid)).Result()
	if err!=nil{
		return 0
	}
	return  int(length)
}