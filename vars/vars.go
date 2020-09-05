package vars

import (
	"gitee.com/kelvins-io/common/queue"
	"gitee.com/kelvins-io/kelvins"
	"gitee.com/kelvins-io/kelvins/config/setting"
	"net/http"
	"time"
)

var (
	App                                *kelvins.GRPCApplication
	EmailConfigSetting                 *EmailConfigSettingS
	QueueAMQPSettingUserRegisterNotice *setting.QueueAMQPSettingS
	QueueServerUserRegisterNotice      *queue.MachineryQueue
	QueueAMQPSettingUserStateNotice    *setting.QueueAMQPSettingS
	QueueServerUserStateNotice         *queue.MachineryQueue
	HttpClient                         = &http.Client{Timeout: 30 * time.Second}
)
