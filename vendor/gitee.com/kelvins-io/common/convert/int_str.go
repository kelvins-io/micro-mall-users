package convert

import (
	"strconv"
	"strings"
)

func Int64ToStr(num int64) string {
	return strconv.FormatInt(num, 10)
}

func IntToStr(num int) string {
	return strconv.Itoa(num)
}

func StrReplace(str string, old string, news ...string) string {
	strRp := str
	for _, new_ := range news {
		strRp = strings.Replace(strRp, old, new_, 1)
	}
	return strRp
}
