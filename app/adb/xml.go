package adb

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/Eminlin/adbCtrlGo/app/format"
)

func RandInt64(min, max int64) int64 {
	if min >= max || min == 0 || max == 0 {
		return max
	}
	return rand.Int63n(max-min) + min
}

//	[36,1410][301,1485]
func ParseBounds(b string) format.ButtonPoint {
	rtn := format.ButtonPoint{}
	if b == "" {
		return rtn
	}
	b = strings.ReplaceAll(b, "][", ",")
	b = strings.ReplaceAll(b, "[", "")
	b = strings.ReplaceAll(b, "]", "")
	bTemp := strings.Split(b, ",")
	if len(bTemp) != 4 {
		fmt.Println(bTemp)
		return rtn
	}
	temp := []int{}
	for _, v := range bTemp {
		t, _ := strconv.Atoi(v)
		temp = append(temp, t)
	}
	rtn.XRangeL = temp[0]
	rtn.XRangeR = temp[2]
	rtn.YRangeL = temp[1]
	rtn.YRangeR = temp[3]
	return rtn
}
