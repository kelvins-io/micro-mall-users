package code

import "gitee.com/kelvins-io/common/errcode"

const (
	CountryCodeEmpty = 29000001
	MobileEmpty      = 29000002
	AccountIdEmpty   = 29000003
	NameEmpty        = 29000004
	AccountNotExist  = 29000005
)

var ErrMap = make(map[int]string)

func init() {
	dict := map[int]string{
		CountryCodeEmpty: "国际码非法",
		MobileEmpty:      "手机号码非法",
		AccountIdEmpty:   "通行证ID不能为空",
		NameEmpty:        "用户名不能为空",
		AccountNotExist:  "通行证账户不存在",
	}
	errcode.RegisterErrMsgDict(dict)
	for key, _ := range dict {
		ErrMap[key] = dict[key]
	}
}
