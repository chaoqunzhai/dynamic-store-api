/**
@Author: chaoqun
* @Date: 2023/10/23 17:28
*/
package systemChan

import (
	"fmt"
	"gorm.io/gorm"
)

const (
	LogIng  = iota

)
type Message struct {
	Orm *gorm.DB
	Table string //表名字
	Data map[string]interface{}

}
var (
	SysChannel chan *Message
)

func init()  {
	SysChannel = make(chan *Message,0)
	go watchSysChannel()
	
}
func SendMessage(m *Message)  {
	SysChannel <- m
}
func watchSysChannel()  {

	for {
		message,ok :=<- SysChannel
		if !ok {
			continue
		}

		fmt.Printf("读取到了消息:%v\n",message)
		switch message.Table {
		case "sys_login_log":
			saveLoginLog(message)
		}
	}
}