package startup

import (
	"gitee.com/cristiane/micro-mall-users/vars"
	"gitee.com/kelvins-io/kelvins/setup"
)

// SetupVars 加载变量
func SetupVars() error {
	var err error
	if vars.QueueAMQPSettingUserRegisterNotice != nil && vars.QueueAMQPSettingUserRegisterNotice.Broker != "" {
		vars.QueueServerUserRegisterNotice, err = setup.NewAMQPQueue(vars.QueueAMQPSettingUserRegisterNotice, nil)
		if err != nil {
			return err
		}
	}

	if vars.QueueAMQPSettingUserStateNotice != nil && vars.QueueAMQPSettingUserStateNotice.Broker != "" {
		vars.QueueServerUserStateNotice, err = setup.NewAMQPQueue(vars.QueueAMQPSettingUserStateNotice, nil)
		if err != nil {
			return err
		}
	}
	return nil
}
