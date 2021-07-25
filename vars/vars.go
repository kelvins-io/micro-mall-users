package vars

import (
	"gitee.com/cristiane/micro-mall-users/pkg/util/goroutine"
	"gitee.com/kelvins-io/common/queue"
	"gitee.com/kelvins-io/kelvins"
	"gitee.com/kelvins-io/kelvins/config/setting"
)

var (
	App                                *kelvins.GRPCApplication
	EmailConfigSetting                 *EmailConfigSettingS
	JwtSetting                         *JwtSettingS
	QueueAMQPSettingUserRegisterNotice *setting.QueueAMQPSettingS
	QueueServerUserRegisterNotice      *queue.MachineryQueue
	QueueAMQPSettingUserStateNotice    *setting.QueueAMQPSettingS
	QueueServerUserStateNotice         *queue.MachineryQueue
	EmailNoticeSetting                 *EmailNoticeSettingS
	GPool                              *goroutine.Pool
)
