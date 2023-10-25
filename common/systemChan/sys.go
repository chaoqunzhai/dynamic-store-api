/**
@Author: chaoqun
* @Date: 2023/10/23 17:27
*/
package systemChan

import (
	"encoding/json"
	"go-admin/cmd/migrate/migration/models"
	"go.uber.org/zap"
)

func saveLoginLog(message *Message)  {
	var l models.SysLoginLog
	var rb []byte
	rb, err := json.Marshal(message.Data)
	if err != nil {
		zap.S().Errorf("用户登录日志记录,Marshal失败 error:%v",err.Error())
		return
	}
	err = json.Unmarshal(rb, &l)
	if err != nil {
		zap.S().Errorf("用户登录日志记录,Unmarshal失败 error:%v",err.Error())
		return
	}
	err = message.Orm.Create(&l).Error
	if err != nil {
		zap.S().Errorf("用户登录日志记录,录入数据失败 error:%v",err.Error())
		return
	}
	return
}