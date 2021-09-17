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
	Unknown                        = 0
	VerifyCodeRegister             = 1
	VerifyCodeLogin                = 2
	VerifyCodePassword             = 3
	UserLoginTemplate              = "尊敬的用户【%s】你好，你于：%v 在微商城使用【%s】登录"
	UserAccountChargeTemplate      = "尊敬的用户【%s】你好，你与：%v 在微商城 充值【%s-%s】成功"
	UserApplyMerchantTemplate      = "尊敬的用户【%s】你好。你与：%v 在微商城申请【%v】商户成功"
	UserModifyMerchantInfoTemplate = "尊敬的用户【%s】你好。你与：%v 在微商城变更商户资料成功"
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
	RpcServiceMicroMallPay   = "micro-mall-pay"
	RpcServiceMicroMallUsers = "micro-mall-users"
)
