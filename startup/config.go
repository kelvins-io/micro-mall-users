package startup

import (
	"gitee.com/cristiane/micro-mall-users/vars"
	"gitee.com/kelvins-io/kelvins/config"
	"gitee.com/kelvins-io/kelvins/config/setting"
	"log"
)

const (
	SectionEmailConfig             = "email-config"
	SectionQueueUserRegisterNotice = "queue-user-register-notice"
	SectionQueueUserStateNotice    = "queue-user-state-notice"
)

// LoadConfig 加载配置对象映射
func LoadConfig() error {
	// 加载email数据源
	log.Printf("[info] Load default config %s", SectionEmailConfig)
	vars.EmailConfigSetting = new(vars.EmailConfigSettingS)
	config.MapConfig(SectionEmailConfig, vars.EmailConfigSetting)
	// 用户注册通知
	log.Printf("[info] Load default config %s", SectionQueueUserRegisterNotice)
	vars.QueueAMQPSettingUserRegisterNotice = new(setting.QueueAMQPSettingS)
	config.MapConfig(SectionQueueUserRegisterNotice, vars.QueueAMQPSettingUserRegisterNotice)
	// 用户事件通知
	log.Printf("[info] Load default config %s", SectionQueueUserStateNotice)
	vars.QueueAMQPSettingUserStateNotice = new(setting.QueueAMQPSettingS)
	config.MapConfig(SectionQueueUserStateNotice, vars.QueueAMQPSettingUserStateNotice)
	return nil
}
