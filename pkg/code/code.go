package code

import "gitee.com/kelvins-io/common/errcode"

const (
	Success                  = 29000000
	ErrorServer              = 29000001
	UserNotExist             = 29000005
	UserExist                = 29000006
	DBDuplicateEntry         = 29000007
	MerchantExist            = 29000008
	MerchantNotExist         = 29000009
	AccountExist             = 29000010
	AccountNotExist          = 29000011
	UserPwdNotMatch          = 29000012
	UserDeliveryInfoExist    = 29000013
	UserDeliveryInfoNotExist = 29000014
	TransactionFailed        = 29000015
)

var ErrMap = make(map[int]string)

func init() {
	dict := map[int]string{
		Success:                  "OK",
		ErrorServer:              "服务器错误",
		UserNotExist:             "用户不存在",
		DBDuplicateEntry:         "Duplicate entry",
		UserExist:                "已存在用户记录，请勿重复创建",
		MerchantExist:            "商户认证材料已存在",
		MerchantNotExist:         "商户未提交材料",
		AccountExist:             "账户已存在",
		AccountNotExist:          "账户不存在",
		UserPwdNotMatch:          "用户密码不匹配",
		UserDeliveryInfoExist:    "用户物流交付信息存在",
		UserDeliveryInfoNotExist: "用户物流交付信息不存在",
		TransactionFailed:        "事务执行失败",
	}
	errcode.RegisterErrMsgDict(dict)
	for key, _ := range dict {
		ErrMap[key] = dict[key]
	}
}

func GetMsg(code int) string {
	v, ok := ErrMap[code]
	if !ok {
		return ""
	}
	return v
}
