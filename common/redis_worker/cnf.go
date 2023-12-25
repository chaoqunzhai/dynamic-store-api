/**
@Author: chaoqun
* @Date: 2023/12/22 11:06
*/
package redis_worker

import "context"

const (
	WorkerOrderStartName = "order" //订单选中导出
	WorkerReportStartName = "report" //配送报表

)
var (
	RedisCtx context.Context
	QueueGroup []string

)
func init()  {
	RedisCtx = context.Background()
	QueueGroup =[]string{
		WorkerOrderStartName,WorkerReportStartName,
	}
	go LoopRedisWorker()
}


type ExportReq struct {
	Queue string `json:"queue"`
	Order []string `json:"order"`
	CId int `json:"c_id"`
	OrmId int `json:"orm_id"`
	Type string `json:"type"` //类型 0:配送订单导出 1:自提订单导出 2:总汇总表导出 3:基于路线导出
}
type GetQueueReq struct {
	CId int `json:"c_id"`
	Name string `json:"name"`
}