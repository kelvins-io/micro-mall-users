package code

import "gitee.com/kelvins-io/common/errcode"

const (
	Success      = 29000000
	ErrorServer  = 29000001
	UserNotExist = 29000005
)

var ErrMap = make(map[int]string)

func init() {
	dict := map[int]string{
		Success:      "OK",
		ErrorServer:  "服务器错误",
		UserNotExist: "用户不存在",
	}
	errcode.RegisterErrMsgDict(dict)
	for key, _ := range dict {
		ErrMap[key] = dict[key]
	}
}
