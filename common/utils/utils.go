/*
*
@Author: chaoqun
* @Date: 2023/5/28 23:46
*/
package utils

import (
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func GenValidateCode(width int) string {
	numeric := [10]byte{0, 1, 2, 3, 5, 6, 7, 8, 9}
	r := len(numeric)
	rand.Seed(time.Now().UnixNano())

	var sb strings.Builder
	for i := 0; i < width; i++ {
		fmt.Fprintf(&sb, "%d", numeric[rand.Intn(r)])
	}
	return sb.String()
}

// 获取当前周几
func HasWeekNumber() int {
	n := time.Now()
	week := 0
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
func IsArray(key string, array []string) bool {
	set := make(map[string]struct{})
	for _, value := range array {
		set[value] = struct{}{}
	}
	_, ok := set[key]
	return ok
}

// 判断当前时间 是否在开始和结束时间区间
// TimeCheckRange("09:00","16:00")
func TimeCheckRange(start, end string) bool {
	now := time.Now()
	yearMD := now.Format("2006-01-02")
	//转换开始时间
	dbStartStr := fmt.Sprintf("%v %v", yearMD, start)
	dbStartTimer, _ := time.Parse("2006-01-02 15:04", dbStartStr)

	//转换结束时间
	dbEndStr := fmt.Sprintf("%v %v", yearMD, end)
	dbEndTimer, _ := time.Parse("2006-01-02 15:04", dbEndStr)

	//转换当前时间
	nowParse, _ := time.Parse("2006-01-02 15:04", time.Now().Format("2006-01-02 15:04"))

	return dbStartTimer.Before(nowParse) && nowParse.Before(dbEndTimer)
}
func ParInt(n float64) float64 {

	value, err := strconv.ParseFloat(fmt.Sprintf("%.2f", n), 64)
	if err != nil {
		return n
	}
	return value
}

// 数组去重
func RemoveRepeatElement(list []string) (result []string) {
	// 创建一个临时map用来存储数组元素
	temp := make(map[string]bool)
	for _, v := range list {
		// 遍历数组元素，判断此元素是否已经存在map中
		_, ok := temp[v]
		if !ok {
			temp[v] = true
			result = append(result, v)
		}
	}
	return result
}

func Avg(a []float64) float64 {
	sum := 0.0

	for i := 0; i < len(a); i++ {
		sum += a[i]
	}
	return ParInt(sum / float64(len(a)))
}
func Min(a []float64) float64 {
	min := a[0]
	for i := 0; i < len(a); i++ {
		if a[i] == 0 {
			continue
		}
		if a[i] < min {
			min = a[i]
		}
	}
	return ParInt(min)
}
func Max(a []float64) float64 {
	max := a[0]
	for i := 0; i < len(a); i++ {
		if a[i] > max {
			max = a[i]
		}
	}
	return ParInt(max)
}

// 处理精度问题
func DecimalMul(n int, k float32) float32 {
	a := decimal.NewFromFloat32(k)
	b := a.Mul(decimal.NewFromInt(int64(n)))

	c, _ := b.Float64()
	return float32(c)

}
func DecimalAdd(n1, n2 float32) float32 {
	a := decimal.NewFromFloat32(n1)
	b := a.Add(decimal.NewFromFloat32(n2))

	c, _ := b.Float64()
	return float32(c)

}

func ReplacePhone(phone string) (err error, phoneText string) {
	//str := `13734351278`ReplacePhone
	pattern := `^(\d{3})(\d{4})(\d{4})$`
	re := regexp.MustCompile(pattern) //确保正则表达式的正确 遇到错误会直接panic
	match := re.MatchString(phone)
	if !match {
		fmt.Println("phone number error")

		return errors.New("非法手机号"), ""
	}
	repStr := re.ReplaceAllString(phone, "$1****$3")

	return nil, repStr
}
func sup(i int64, n int) string {
	m := fmt.Sprintf("%d", i)

	for len(m) < n {
		m = fmt.Sprintf("0%s", m)
	}
	return m
}

// 生成订单ID,1S内支持1万个订单ID,
func GenUUID() int {
	t := time.Now()
	nanoTime := int64(time.Now().Nanosecond())
	rand.Seed(nanoTime)
	p := rand.Intn(10000) % (rand.Intn(100) + 1)
	ps := sup(int64(p), 2)
	s := t.Format("20060102150405")
	m := t.UnixNano()/1e6 - t.UnixNano()/1e9*1e3

	rs := sup(m, 3)

	n := fmt.Sprintf("%v%v%v", s, rs, ps)

	number, _ := strconv.ParseInt(n, 10, 64)
	//fmt.Printf("%v to %v\n",n,number)
	return int(number)
}
