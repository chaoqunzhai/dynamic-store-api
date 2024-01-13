/**
@Author: chaoqun
* @Date: 2023/12/27 15:33
*/
package xlsx_export

import (
	"fmt"
	"go-admin/cmd/migrate/migration/models"
	"go-admin/common/redis_db"
	"go-admin/global"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

//更新table中状态

func SaveExportDb(ormId int,cid int,fileName string,successTag bool,msg string,orm *gorm.DB)  error{
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
	orm.Model(&models.CompanyTasks{}).Where("id = ? and c_id = ?",
		ormId,cid).Updates(map[string]interface{}{
		"status":status,
		"path":fileName,
		"msg":msg,
	})
	return nil

}
//如果key下的list位空 ,那就支持清空这个key

func EmptyKey(RedisKey string,keyIndex int) {
	fmt.Println("删除key",RedisKey,keyIndex)
	err :=redis_db.RedisCli.LTrim(global.RedisCtx,RedisKey,1,int64(keyIndex)).Err()
	if err!=nil{
		zap.S().Errorf("清理redis key:%v 数据清理失败:%v",RedisKey,err)
	}else {
		zap.S().Infof("redis key:%v 消费完毕,数据清理成功",RedisKey)
	}
}