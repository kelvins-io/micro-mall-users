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

	UserInfoSearchNotice    = "user_info_search_notice"
	UserInfoSearchNoticeErr = "user_info_search_notice_err"
)

const (
	Unknown            = 0
	VerifyCodeRegister = 1
	VerifyCodeLogin    = 2
	VerifyCodePassword = 3
)

const (
	DefaultCountryCode = "86"
)

var MsgFlags = map[int]string{
	Unknown:                         "未知",
	VerifyCodeRegister:              "注册",
	VerifyCodeLogin:                 "登录",
	VerifyCodePassword:              "修改/重置密码",
	UserStateEventTypeRegister:      "注册",
	UserStateEventTypePwdModify:     "修改密码",
	UserStateEventTypeLogin:         "登录上线",
	UserStateEventTypeLogout:        "退出登录",
	UserInfoSearchNoticeType:        "用户注册、更新信息",
	MerchantInfoSearchNoticeType:    "商户认证注册、更新信息",
	UserStateEventTypeAccountCharge: "账号充值",
	UserStateEventTypeMerchantInfo:  "商户认证申请、更新信息",
}

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}
	return MsgFlags[Unknown]
}

const (
	CacheKeyUserOnlineSate = "online_state_"
)

const (
	UserStateEventTypeRegister      = 10010
	UserStateEventTypeLogin         = 10011
	UserStateEventTypeLogout        = 10012
	UserStateEventTypePwdModify     = 10013
	UserStateEventTypeAccountCharge = 10014
	UserStateEventTypeMerchantInfo  = 10015
)

const (
	UserInfoSearchNoticeType     = 50001
	MerchantInfoSearchNoticeType = 50002
)

const (
	VerifyCodeTemplate = "尊敬的用户【%v】你好，验证码 %v，用于%v，%v分钟内有效，验证码提供给其他人可能导致账号被盗，请勿泄漏，谨防被骗。"
)

type UserVerifyCode struct {
	VerifyCode string `json:"verify_code"`
	Expire     int64  `json:"expire"`
}

type UserInfoSearch struct {
	UserName    string `json:"user_name"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	IdCardNo    string `json:"id_card_no"`
	ContactAddr string `json:"contact_addr"`
}

type MerchantInfoSearch struct {
	Uid          int64  `json:"uid"`
	UserName     string `json:"user_name"`
	MerchantCode string `json:"merchant_code"`
	RegisterAddr string `json:"register_addr"`
	HealthCardNo string `json:"health_card_no"`
	TaxCardNo    string `json:"tax_card_no"`
}

type CommonBusinessMsg struct {
	Type    int    `json:"type"`
	Tag     string `json:"tag"`
	UUID    string `json:"uuid"`
	Time    string `json:"time"`
	Content string `json:"content"`
}

type UserRegisterNotice struct {
	CountryCode string `json:"country_code"`
	Phone       string `json:"phone"`
	State       int    `json:"state"`
}

type UserStateNotice struct {
	Uid   int               `json:"uid"`
	Extra map[string]string `json:"extra"`
}

type UserOnlineState struct {
	Uid   int    `json:"uid"`
	State string `json:"state"`
	Time  string `json:"time"`
}

const (
	UserOnlineStateOnline         = "online-login"
	UserOnlineStateForbiddenLogin = "forbidden to login"
)

const (
	UserLoginFailureFrequency    = "login-failure-frequency"
	UserLoginFailureFrequencyMax = 3
)

const (
	RpcServiceMicroMallPay    = "micro-mall-pay"
	RpcServiceMicroMallUsers  = "micro-mall-users"
	RpcServiceMicroMallSearch = "micro-mall-search"
)
