package vars

import (
	"gitee.com/kelvins-io/common/log"
	"gitee.com/kelvins-io/common/queue"
	"gitee.com/kelvins-io/kelvins/config/setting"
	"github.com/gomodule/redigo/redis"
	"github.com/jinzhu/gorm"
	"net/http"
	"time"
	"xorm.io/xorm"
)

var (
	DBEngineXORM                       xorm.EngineInterface
	DBEngineGORM                       *gorm.DB
	LoggerSetting                      *setting.LoggerSettingS
	AccessLogger                       log.LoggerContextIface
	ErrorLogger                        log.LoggerContextIface
	BusinessLogger                     log.LoggerContextIface
	ServerSetting                      *setting.ServerSettingS
	JwtSetting                         *setting.JwtSettingS
	MysqlSettingMicroMall              *setting.MysqlSettingS
	RedisSettingMicroMall              *setting.RedisSettingS
	EmailConfigSetting                 *EmailConfigSettingS
	RedisPoolMicroMall                 *redis.Pool
	QueueAMQPSettingUserRegisterNotice *setting.QueueAMQPSettingS
	QueueServerUserRegisterNotice      *queue.MachineryQueue
	QueueAMQPSettingUserStateNotice    *setting.QueueAMQPSettingS
	QueueServerUserStateNotice         *queue.MachineryQueue
	HttpClient                         = &http.Client{Timeout: 30 * time.Second}
)
