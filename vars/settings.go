package vars

type EmailConfigSettingS struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     string `json:"port"`
}

type QueueRedisSettingS struct {
	Broker          string
	DefaultQueue    string
	ResultBackend   string
	ResultsExpireIn int
}

type QueueAMQPSettingS struct {
	Broker           string
	DefaultQueue     string
	ResultBackend    string
	ResultsExpireIn  int
	Exchange         string
	ExchangeType     string
	BindingKey       string
	PrefetchCount    int
	TaskRetryCount   int
	TaskRetryTimeout int
}

// QueueAliAMQPSettingS defines for aliyun AMQP queue
type QueueAliAMQPSettingS struct {
	AccessKey       string
	SecretKey       string
	AliUid          int
	EndPoint        string
	VHost           string
	DefaultQueue    string
	ResultBackend   string
	ResultsExpireIn int
	Exchange        string
	ExchangeType    string
	BindingKey      string
	PrefetchCount   int
}
