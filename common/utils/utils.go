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
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func DirNotCreate(dir string) {

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// 目录不存在，创建它
		err = os.Mkdir(dir, 0755) // 0755是权限设置，你可以根据需要修改
		if err != nil {
			fmt.Printf("创建目录 %s 失败: %s\n", dir, err)
			return
		}
		fmt.Printf("目录 %s 已创建\n", dir)
	} else if err != nil {
		// 出现了其他错误
		fmt.Printf("获取目录 %s 信息时出错: %s\n", dir, err)
		return
	} else {
		// 目录已经存在
		//fmt.Printf("目录 %s 已存在\n", dir)
	}
	return

}
// GetWeekdayTimestamps 获取指定星期几的开始和结束时间戳
func GetWeekdayTimestamps(weekdayNumber int) (weekTime time.Time, err error) {
	// 获取当前时间

	now := time.Now()
	var weekday time.Weekday
	switch weekdayNumber {
	case 0:
		weekday = time.Sunday

		return now.AddDate(0, 0, int(now.Weekday()) + 1),nil
	case 1:
		weekday = time.Monday
	case 2:
		weekday = time.Tuesday
	case 3:
		weekday = time.Wednesday
	case 4:
		weekday = time.Thursday
	case 5:
		weekday = time.Friday
	case 6:
		weekday = time.Saturday

	default:

		return  time.Time{},errors.New("非法日期")
	}
	// 获取当前周的第一天的时间
	startOfWeek := now.AddDate(0, 0, int(-now.Weekday()+weekday))


	return startOfWeek,nil
}
func StructToMap(obj interface{}) map[string]interface{} {

	value := reflect.ValueOf(obj)
	structType := value.Type()
	structMap := make(map[string]interface{})

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		valueField := value.Field(i)
		jsonTag :=field.Tag.Get("json")
		key:=""
		if jsonTag != ""{
			key = jsonTag
		}else {
			key = field.Name
		}
		structMap[key] = valueField.Interface()
	}

	return structMap
}

//对数值进行补0 或者 带有小数点的数字 只保留2位
func StringDecimal(value interface{}) string {
	amount,err :=decimal.NewFromString(fmt.Sprintf("%v",value))
	if err!=nil{
		return fmt.Sprintf("%v",value)
	}
	return  amount.StringFixed(2)

}
func StringToInt(v interface{}) int {
	n, _ := strconv.Atoi(fmt.Sprintf("%v", v))
	return n
}
func StringToFloat64(v interface{}) float64 {
	n, _ := strconv.ParseFloat(fmt.Sprintf("%v", v), 64)
	return n
}
//
func RoundDecimalFlot64(value interface{}) float64 {
	toStr := fmt.Sprintf("%v", value)
	amount3, _ := decimal.NewFromString(toStr)
	f,_:=amount3.Round(2).Float64()

	return f
}
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
func MinAndMax(values []float64) (float64, float64) {
	min1 := values[0] //assign the first element equal to min
	max1 := values[0] //assign the first element equal to max
	for _, number := range values {
		if number < min1 {
			min1 = number
		}
		if number > max1 {
			max1 = number
		}
	}
	return min1, max1
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
func IsArrayInt(key int, array []int) bool {
	set := make(map[int]struct{})
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
func RemoveRepeatStr(list []string) (result []string) {
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
// 数值去重
func RemoveRepeatInt(list []int) (result []int) {
	// 创建一个临时map用来存储数组元素
	temp := make(map[int]bool)
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

func RoundDecimal(value interface{}) decimal.Decimal {
	toStr := fmt.Sprintf("%v", value)
	amount3, _ := decimal.NewFromString(toStr) //0.8
	return amount3.Round(2)  //0.80
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
func CreateCode() string {
	return fmt.Sprintf("%06v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000)) //这里面前面的04v是和后面的1000相对应的
}

// 求并集
func Union(slice1, slice2 []string) []string {
	m := make(map[string]int)
	for _, v := range slice1 {
		m[v]++
	}

	for _, v := range slice2 {
		times, _ := m[v]
		if times == 0 {
			slice1 = append(slice1, v)
		}
	}
	return slice1
}
func IntersectInt(slice1, slice2 []int) []int {
	m := make(map[int]int)
	nn := make([]int, 0)
	for _, v := range slice1 {
		m[v]++
	}

	for _, v := range slice2 {
		times, _ := m[v]
		if times == 1 {
			nn = append(nn, v)
		}
	}
	return nn
}

// 求交集
func Intersect(slice1, slice2 []string) []string {
	m := make(map[string]int)
	nn := make([]string, 0)
	for _, v := range slice1 {
		m[v]++
	}

	for _, v := range slice2 {
		times, _ := m[v]
		if times == 1 {
			nn = append(nn, v)
		}
	}
	return nn
}

// 求差集 slice1-并集
func Difference(slice1, slice2 []string) []string {
	m := make(map[string]int)
	nn := make([]string, 0)
	inter := Intersect(slice1, slice2)
	for _, v := range inter {
		m[v]++
	}

	for _, value := range slice1 {
		times, _ := m[value]
		if times == 0 {
			nn = append(nn, value)
		}
	}
	return nn
}
func DifferenceInt(slice1, slice2 []int) []int {
	m := make(map[int]int)
	nn := make([]int, 0)
	inter := IntersectInt(slice1, slice2)
	for _, v := range inter {
		m[v]++
	}

	for _, value := range slice1 {
		times, _ := m[value]
		if times == 0 {
			nn = append(nn, value)
		}
	}
	return nn
}


//随机字符串
func GetRandStr(n int)  string {
	rand.Seed(time.Now().UnixNano())

	// 生成4位随机字母
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	result := make([]byte, n)
	for i := 0; i < n; i++ {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}

func RemoveDirectory(dir string) error {
	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			err := os.RemoveAll(path) // 递归删除目录及其内容
			fmt.Println("删除path",path)
			if err != nil {
				return err
			}
		} else {
			err := os.Remove(path) // 删除文件
			fmt.Println("删除文件",path)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return os.Remove(dir) // 删除空目录本身
}
// isTimeOverlap 检查两个时间范围是否有重叠
func IsTimeOverlap(start1, end1, start2, end2 time.Time) bool {
	// 如果第一个时间范围的结束时间在第二个时间范围的开始时间之前，
	// 或者第二个时间范围的结束时间在第一个时间范围的开始时间之前，
	// 则两个时间范围不重叠。
	if end1.Before(start2) || end2.Before(start1) {
		return false
	}
	// 否则，两个时间范围有重叠。
	return true
}
