/**
@Author: chaoqun
* @Date: 2023/12/22 11:06
*/
package redis_worker

import "context"

const (
	WorkerStartName = "task"

)
var (
	RedisCtx context.Context


)
func init()  {
	RedisCtx = context.Background()
	go LoopRedisWorker()
}


type ExportReq struct {
	Order []string `json:"order"`
	CId int `json:"c_id"`
	OrmId int `json:"orm_id"`
	Type string `json:"type"` //类型 0:配送订单导出 1:自提订单导出 2:总汇总表导出 3:基于路线导出
}
type GetQueueReq struct {
	CId int `json:"c_id"`
	Name string `json:"name"`
}