// 通用错误码组件-错误消息相关接口
package errcode

import "log"

var errMsgDict = make(map[int]string)

func RegisterErrMsgDict(dict map[int]string) {
	for code, errMsg := range dict {
		if errMsgDict[code] != "" {
			log.Fatalf("错误码初始化错误,重复定义的code:%d", code)
		}
		if code < 0 || (1 <= code && code <= 9999999) || code > 29999999 {
			log.Fatalf("错误码初始化错误,不符合规范的code:%d", code)
		}
		errMsgDict[code] = errMsg
	}
}

func GetErrMsg(code int) string {
	return errMsgDict[code]
}
