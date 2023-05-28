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
