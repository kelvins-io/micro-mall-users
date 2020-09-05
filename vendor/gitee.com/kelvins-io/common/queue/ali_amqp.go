package queue

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"gitee.com/kelvins-io/common/convert"

	"strings"
	"time"
)

func convertAliyunUserName(accessKey string, userId int) string {
	var bf bytes.Buffer
	bf.WriteString("0:")
	bf.WriteString(convert.IntToStr(userId))
	bf.WriteString(":")
	bf.WriteString(accessKey)

	return base64.StdEncoding.EncodeToString(bf.Bytes())
}

func convertAliyunPassword(secretKey string) string {
	var ts = time.Now().UnixNano() / 1e6

	var mac = hmac.New(sha1.New, []byte(convert.Int64ToStr(ts)))
	mac.Write([]byte(secretKey))

	var macStr = strings.ToUpper(hex.EncodeToString(mac.Sum(nil)))

	var bf bytes.Buffer
	bf.WriteString(macStr)
	bf.WriteString(":")
	bf.WriteString(convert.Int64ToStr(ts))

	return base64.StdEncoding.EncodeToString(bf.Bytes())
}
