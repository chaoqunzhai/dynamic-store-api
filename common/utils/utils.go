/**
@Author: chaoqun
* @Date: 2023/5/28 23:46
*/
package utils

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func GenValidateCode(width int) string {
	numeric := [10]byte{0,1,2,3,5,6,7,8,9}
	r := len(numeric)
	rand.Seed(time.Now().UnixNano())

	var sb strings.Builder
	for i := 0; i < width; i++ {
		fmt.Fprintf(&sb, "%d", numeric[ rand.Intn(r) ])
	}
	return sb.String()
}
//获取当前周几
func HasWeekNumber() int {
	n :=time.Now()
	week :=0
	switch n.Weekday().String() {
	case "Sunday":
		week = 0
	case "Monday":
		week = 1
	case "Tuesday":
		week = 2
	case "Wednesday":
		week = 3
	case "Thursday":
		week = 4
	case "Friday":
		week = 5
	case "Saturday":
		week = 6
	}
	return week
}

//判断当前时间 是否在开始和结束时间区间
//TimeCheckRange("09:00","16:00")
func TimeCheckRange(start,end string) bool {
	now :=time.Now()
	yearMD:=now.Format("2006-01-02")
	//转换开始时间
	dbStartStr :=fmt.Sprintf("%v %v",yearMD,start)
	dbStartTimer,_:=time.Parse("2006-01-02 15:04",dbStartStr)

	//转换结束时间
	dbEndStr :=fmt.Sprintf("%v %v",yearMD,end)
	dbEndTimer,_:=time.Parse("2006-01-02 15:04",dbEndStr)

	//转换当前时间
	nowParse,_ := time.Parse("2006-01-02 15:04", time.Now().Format("2006-01-02 15:04"))

	return dbStartTimer.Before(nowParse) && nowParse.Before(dbEndTimer)
}