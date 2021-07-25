package startup

import (
	"gitee.com/cristiane/micro-mall-users/vars"
	"gitee.com/kelvins-io/kelvins/config"
	"gitee.com/kelvins-io/kelvins/config/setting"
)

const (
	SectionEmailConfig             = "email-config"
	SectionQueueUserRegisterNotice = "queue-user-register-notice"
	SectionQueueUserStateNotice    = "queue-user-state-notice"
	SectionJwt                     = "kelvins-jwt"
	EmailNotice                    = "email-notice"
)

// LoadConfig 加载配置对象映射
func LoadConfig() error {
	// 加载email数据源
	vars.EmailConfigSetting = new(vars.EmailConfigSettingS)
	config.MapConfig(SectionEmailConfig, vars.EmailConfigSetting)
	// 邮件通知
	vars.EmailNoticeSetting = new(vars.EmailNoticeSettingS)
	config.MapConfig(EmailNotice, vars.EmailNoticeSetting)
	// 用户注册通知
	vars.QueueAMQPSettingUserRegisterNotice = new(setting.QueueAMQPSettingS)
	config.MapConfig(SectionQueueUserRegisterNotice, vars.QueueAMQPSettingUserRegisterNotice)
	// 用户事件通知
	vars.QueueAMQPSettingUserStateNotice = new(setting.QueueAMQPSettingS)
	config.MapConfig(SectionQueueUserStateNotice, vars.QueueAMQPSettingUserStateNotice)
	// 用户认证token
	vars.JwtSetting = new(vars.JwtSettingS)
	config.MapConfig(SectionJwt, vars.JwtSetting)
	return nil
}
