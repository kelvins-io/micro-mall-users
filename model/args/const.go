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
	Unknown:                         "未知",
	VerifyCodeRegister:              "注册",
	VerifyCodeLogin:                 "登录",
	VerifyCodePassword:              "修改/重置密码",
	UserStateEventTypeRegister:      "注册",
	UserStateEventTypePwdModify:     "修改密码",
	UserStateEventTypeLogin:         "登录上线",
	UserStateEventTypeLogout:        "退出登录",
	UserInfoSearchNoticeType:        "用户注册/更新信息",
	MerchantsMaterialInfoNoticeType: "商户申请通知",
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
	UserInfoSearchNoticeType        = 10014
	MerchantsMaterialInfoNoticeType = 10015
)

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
	Content string `json:"content"`
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
	UserOnlineStateOnline = "online-login"
	UserOnlineStateForbiddenLogin = "forbidden to login"
)

const (
	UserLoginFailureFrequency = "login-failure-frequency"
	UserLoginFailureFrequencyMax = 3
)

const (
	RpcServiceMicroMallPay    = "micro-mall-pay"
	RpcServiceMicroMallUsers  = "micro-mall-users"
	RpcServiceMicroMallSearch = "micro-mall-search"
)
