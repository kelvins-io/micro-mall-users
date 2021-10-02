package startup

import (
	"gitee.com/cristiane/micro-mall-users/model/args"
	"gitee.com/cristiane/micro-mall-users/vars"
	"gitee.com/kelvins-io/kelvins"
	"gitee.com/kelvins-io/kelvins/setup"
	"gitee.com/kelvins-io/kelvins/util/queue_helper"
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

	if vars.QueueAMQPSettingUserInfoSearchNotice != nil && vars.QueueAMQPSettingUserInfoSearchNotice.Broker != "" {
		vars.QueueServerUserInfoSearch, err = setup.NewAMQPQueue(vars.QueueAMQPSettingUserInfoSearchNotice, nil)
		if err != nil {
			return err
		}
		vars.QueueServerUserInfoSearchPusher, err = queue_helper.NewPublishService(vars.QueueServerUserInfoSearch, &queue_helper.PushMsgTag{
			DeliveryTag:    args.UserInfoSearchNotice,
			DeliveryErrTag: args.UserInfoSearchNoticeErr,
			RetryCount:     vars.QueueAMQPSettingUserInfoSearchNotice.TaskRetryCount,
			RetryTimeout:   vars.QueueAMQPSettingUserInfoSearchNotice.TaskRetryTimeout,
		}, kelvins.BusinessLogger)
		if err != nil {
			return err
		}
	}
	return nil
}
