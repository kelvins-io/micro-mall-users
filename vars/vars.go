package vars

import (
	"gitee.com/kelvins-io/common/queue"
	"gitee.com/kelvins-io/kelvins/config/setting"
	"gitee.com/kelvins-io/kelvins/util/queue_helper"
)

var (
	EmailConfigSetting                   *EmailConfigSettingS
	QueueAMQPSettingUserRegisterNotice   *setting.QueueAMQPSettingS
	QueueServerUserRegisterNotice        *queue.MachineryQueue
	QueueAMQPSettingUserStateNotice      *setting.QueueAMQPSettingS
	QueueServerUserStateNotice           *queue.MachineryQueue
	QueueAMQPSettingUserInfoSearchNotice *setting.QueueAMQPSettingS
	QueueServerUserInfoSearch            *queue.MachineryQueue
	QueueServerUserInfoSearchPusher      *queue_helper.PublishService
	VerifyCodeSetting                    *VerifyCodeSettingS
)
