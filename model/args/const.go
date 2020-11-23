package args

type MerchantsMaterialInfo struct {
	Uid          int64
	MaterialId   int64
	RegisterAddr string
	HealthCardNo string
	Identity     int32
	State        int32
	TaxCardNo    string
}

type RegisterResult struct {
	InviteCode string `json:"invite_code"`
}

const (
	TaskNameUserRegisterNotice    = "task_user_register_notice"
	TaskNameUserRegisterNoticeErr = "task_user_register_notice_err"

	TaskNameUserStateNotice    = "task_user_state_notice"
	TaskNameUserStateNoticeErr = "task_user_state_notice_err"
)

const (
	Unknown            = 0
	VerifyCodeRegister = 1
	VerifyCodeLogin    = 2
	VerifyCodePassword = 3
	VerifyCodeTemplate = "【%v】验证码 %v，用于%v，%v分钟内有效，验证码提供给其他人可能导致账号被盗，请勿泄漏，谨防被骗。"
)

var MsgFlags = map[int]string{
	Unknown:                     "未知",
	VerifyCodeRegister:          "注册",
	VerifyCodeLogin:             "登录",
	VerifyCodePassword:          "修改/重置密码",
	UserStateEventTypeRegister:  "注册",
	UserStateEventTypePwdModify: "修改密码",
	UserStateEventTypeLogin:     "登录上线",
	UserStateEventTypeLogout:    "退出登录",
}

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}
	return MsgFlags[Unknown]
}

const (
	CacheKeyUserSate = "user_state_"
)
const (
	UserStateEventTypeRegister  = 10010
	UserStateEventTypeLogin     = 10011
	UserStateEventTypeLogout    = 10012
	UserStateEventTypePwdModify = 10013
)

type CommonBusinessMsg struct {
	Type int    `json:"type"`
	Tag  string `json:"tag"`
	UUID string `json:"uuid"`
	Msg  string `json:"msg"`
}

type UserRegisterNotice struct {
	CountryCode string `json:"country_code"`
	Phone       string `json:"phone"`
	Time        string `json:"time"`
	State       int    `json:"state"`
}
type UserStateNotice struct {
	Uid  int    `json:"uid"`
	Time string `json:"time"`
}

type UserOnlineState struct {
	Uid   int    `json:"uid"`
	State string `json:"state"`
	Time  string `json:"time"`
}

const (
	RpcServiceMicroMallPay = "micro-mall-pay"
)
