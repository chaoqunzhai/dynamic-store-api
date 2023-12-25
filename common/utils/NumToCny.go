/**
@Author: chaoqun
* @Date: 2023/12/25 16:50
*/
package utils

import (
	"go.uber.org/zap"
	"regexp"
	"strconv"
)

func ConvertNumToCny(num float64) string {
	var err error
	strnum := strconv.FormatFloat(num*100, 'f', 0, 64)
	sliceUnit := []string{"仟", "佰", "拾", "亿", "仟", "佰", "拾", "万", "仟", "佰", "拾", "元", "角", "分"}
	// log.Println(sliceUnit[:len(sliceUnit)-2])
	s := sliceUnit[len(sliceUnit)-len(strnum) : len(sliceUnit)]
	upperDigitUnit := map[string]string{"0": "零", "1": "壹", "2": "贰", "3": "叁", "4": "肆", "5": "伍", "6": "陆", "7": "柒", "8": "捌", "9": "玖"}
	str := ""
	for k, v := range strnum[:] {
		str = str + upperDigitUnit[string(v)] + s[k]
	}
	reg, err := regexp.Compile(`零角零分$`)
	str = reg.ReplaceAllString(str, "整")

	reg, err = regexp.Compile(`零角`)
	str = reg.ReplaceAllString(str, "零")

	reg, err = regexp.Compile(`零分$`)
	str = reg.ReplaceAllString(str, "整")

	reg, err = regexp.Compile(`零[仟佰拾]`)
	str = reg.ReplaceAllString(str, "零")

	reg, err = regexp.Compile(`零{2,}`)
	str = reg.ReplaceAllString(str, "零")

	reg, err = regexp.Compile(`零亿`)
	str = reg.ReplaceAllString(str, "亿")

	reg, err = regexp.Compile(`零万`)
	str = reg.ReplaceAllString(str, "万")

	reg, err = regexp.Compile(`零*元`)
	str = reg.ReplaceAllString(str, "元")

	reg, err = regexp.Compile(`亿零{0, 3}万`)
	str = reg.ReplaceAllString(str, "^元")

	reg, err = regexp.Compile(`零元`)
	str = reg.ReplaceAllString(str, "零")

	if err!=nil{
		zap.S().Errorf("转换失败：%v",err)
	}
	return str
}