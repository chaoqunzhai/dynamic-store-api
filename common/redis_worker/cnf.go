/**
@Author: chaoqun
* @Date: 2023/12/22 11:06
*/
package redis_worker

func init()  {
	work :=LoopRedisWorker{
	}

	go work.Start()
}

