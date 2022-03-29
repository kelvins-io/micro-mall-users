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
	AccountStateLock         = 29000016
	AccountStateInvalid      = 29000017
	UserChargeRun            = 29000018
	UserChargeSuccess        = 29000019
	UserChargeTradeNoEmpty   = 29000020
	UserStateNotVerify       = 29000021
	UserStateForbiddenLogin  = 29000022
	ErrorVerifyCodeEmpty     = 29000023
	ErrorVerifyCodeInvalid   = 29000024
	ErrorVerifyCodeExpire    = 29000025
	ErrorVerifyCodeInterval  = 29000027
	ErrorVerifyCodeLimited   = 29000028
	ErrorVerifyCodeForbidden = 29000029
	ErrorInviteCodeInvalid   = 29000030
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
		AccountStateLock:         "用户账户锁定中",
		AccountStateInvalid:      "用户账户无效",
		UserChargeRun:            "本次充值交易正在进行中",
		UserChargeSuccess:        "本次充值交易已成功",
		UserChargeTradeNoEmpty:   "本次充值交易号为空",
		UserStateNotVerify:       "用户身份未验证或审核或被锁定",
		UserStateForbiddenLogin:  "用户被禁止登录",
		ErrorVerifyCodeEmpty:     "验证码为空",
		ErrorVerifyCodeInvalid:   "验证码无效",
		ErrorVerifyCodeExpire:    "验证码过期",
		ErrorVerifyCodeInterval:  "验证码出错",
		ErrorVerifyCodeForbidden: "验证码权限不够",
		ErrorInviteCodeInvalid:   "无效的邀请码",
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
