package startup

import (
	"gitee.com/cristiane/micro-mall-users/vars"
	"gitee.com/kelvins-io/kelvins/setup"
)

// SetupVars 加载变量
func SetupVars() error {
	if vars.QueueAMQPSettingUserRegisterNotice != nil && vars.QueueAMQPSettingUserRegisterNotice.Broker != "" {
		vars.QueueServerUserRegisterNotice = setup.NewAMQPQueue(vars.QueueAMQPSettingUserRegisterNotice, nil)
	}

	if vars.QueueAMQPSettingUserStateNotice != nil && vars.QueueAMQPSettingUserStateNotice.Broker != "" {
		vars.QueueServerUserStateNotice = setup.NewAMQPQueue(vars.QueueAMQPSettingUserStateNotice, nil)
	}
	return nil
}
