package convert

import (
	"bytes"
	"strconv"
	"strings"
)

func StringIpToInt(ipstring string) int64 {
	ipSegs := strings.Split(ipstring, ".")
	var ipInt int64 = 0
	var pos uint = 24
	for _, ipSeg := range ipSegs {
		tempInt, _ := strconv.ParseInt(ipSeg, 10, 64)
		tempInt = tempInt << pos
		ipInt = ipInt | tempInt
		pos -= 8
	}
	return ipInt
}

func IpIntToString(ipInt int64) string {
	ipSegs := make([]string, 4)
	var lenght int = len(ipSegs)
	buffer := bytes.NewBufferString("")
	for i := 0; i < lenght; i++ {
		tempInt := ipInt & 0xFF
		ipSegs[lenght-i-1] = strconv.FormatInt(tempInt, 10)
		ipInt = ipInt >> 8
	}
	for i := 0; i < lenght; i++ {
		buffer.WriteString(ipSegs[i])
		if i < lenght-1 {
			buffer.WriteString(".")
		}
	}
	return buffer.String()
}
