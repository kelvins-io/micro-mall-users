package vars

import (
	"gitee.com/kelvins-io/common/queue"
	"gitee.com/kelvins-io/kelvins/config/setting"
)

var (
	EmailConfigSetting                 *EmailConfigSettingS
	JwtSetting                         *JwtSettingS
	QueueAMQPSettingUserRegisterNotice *setting.QueueAMQPSettingS
	QueueServerUserRegisterNotice      *queue.MachineryQueue
	QueueAMQPSettingUserStateNotice    *setting.QueueAMQPSettingS
	QueueServerUserStateNotice         *queue.MachineryQueue
	EmailNoticeSetting                 *EmailNoticeSettingS
)
