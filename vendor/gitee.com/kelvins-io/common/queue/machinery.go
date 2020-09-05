package queue

import (
	"bytes"
	"context"
	"errors"

	"github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/backends/result"
	"github.com/RichardKnop/machinery/v1/config"
	"github.com/RichardKnop/machinery/v1/tasks"

	backendsiface "github.com/RichardKnop/machinery/v1/backends/iface"
	brokersiface "github.com/RichardKnop/machinery/v1/brokers/iface"
)

const (
	QueueRabbitmqExchangetypeDirect = "direct"
	QueueRabbitmqExchangetypeFanout = "fanout"
	QueueRabbitmqExchangetypeTopic  = "topic"
)

// mockgen -destination=./machinery_mock.go -package=queue gitee.com/kelvins-io/common/queue TaskServerIface
type TaskServerIface interface {
	NewWorker(consumerTag string, concurrency int) *machinery.Worker
	NewCustomQueueWorker(consumerTag string, concurrency int, queue string) *machinery.Worker
	GetBroker() brokersiface.Broker
	SetBroker(broker brokersiface.Broker)
	GetBackend() backendsiface.Backend
	SetBackend(backend backendsiface.Backend)
	GetConfig() *config.Config
	SetConfig(cnf *config.Config)
	RegisterTasks(namedTaskFuncs map[string]interface{}) error
	RegisterTask(name string, taskFunc interface{}) error
	IsTaskRegistered(name string) bool
	GetRegisteredTask(name string) (interface{}, error)
	SendTaskWithContext(ctx context.Context, signature *tasks.Signature) (*result.AsyncResult, error)
	SendTask(signature *tasks.Signature) (*result.AsyncResult, error)
	SendChainWithContext(ctx context.Context, chain *tasks.Chain) (*result.ChainAsyncResult, error)
	SendChain(chain *tasks.Chain) (*result.ChainAsyncResult, error)
	SendGroupWithContext(ctx context.Context, group *tasks.Group, sendConcurrency int) ([]*result.AsyncResult, error)
	SendGroup(group *tasks.Group, sendConcurrency int) ([]*result.AsyncResult, error)
	SendChordWithContext(ctx context.Context, chord *tasks.Chord, sendConcurrency int) (*result.ChordAsyncResult, error)
	SendChord(chord *tasks.Chord, sendConcurrency int) (*result.ChordAsyncResult, error)
	GetRegisteredTaskNames() []string
}

type MachineryQueue struct {
	TaskServer TaskServerIface
}

func newMachineryQueue(broker string, defaultQueue string, resultBackend string, resultsExpireIn int, amqp *config.AMQPConfig, namedTaskFuncs map[string]interface{}) (*MachineryQueue, error) {
	config2 := &config.Config{
		Broker:          broker,
		DefaultQueue:    defaultQueue,
		ResultBackend:   resultBackend,
		ResultsExpireIn: resultsExpireIn,
		AMQP:            amqp,
	}

	server, err := machinery.NewServer(config2)
	if err != nil {
		return nil, err
	}

	err = server.RegisterTasks(namedTaskFuncs)
	if err != nil {
		return nil, err
	}

	return &MachineryQueue{TaskServer: server}, nil
}

func NewRabbitMqQueue(broker string, defaultQueue string, resultBackend string, resultsExpireIn int, exchange string, exchangeType string,
	bindingKey string, prefetchCount int, queueBindingArgs map[string]interface{}, namedTaskFuncs map[string]interface{}) (*MachineryQueue, error) {

	allowExchangeType := map[string]bool{
		QueueRabbitmqExchangetypeDirect: true,
		QueueRabbitmqExchangetypeFanout: true,
		QueueRabbitmqExchangetypeTopic:  true,
	}

	if _, ok := allowExchangeType[exchangeType]; !ok {
		return nil, errors.New("error exchange type for rabbitMq")
	}

	amqp := &config.AMQPConfig{
		Exchange:         exchange,
		ExchangeType:     exchangeType,
		BindingKey:       bindingKey,
		PrefetchCount:    prefetchCount,
		QueueBindingArgs: queueBindingArgs,
	}

	return newMachineryQueue(
		broker,
		defaultQueue,
		resultBackend,
		resultsExpireIn,
		amqp,
		namedTaskFuncs,
	)
}

func NewRedisQueue(broker string, defaultQueue string, resultBackend string, resultsExpireIn int, namedTaskFuncs map[string]interface{}) (*MachineryQueue, error) {
	return newMachineryQueue(
		broker,
		defaultQueue,
		resultBackend,
		resultsExpireIn,
		nil,
		namedTaskFuncs,
	)
}

type AliAMQPConfig struct {
	AccessKey        string // 阿里云accesskey
	SecretKey        string // 阿里云secretKey
	AliUid           int    // 阿里云资源owner账户 ID 信息，点击在控制台右上角客户头像进入账号管理查看
	EndPoint         string // 阿里云amqp接入点
	VHost            string // vhost
	DefaultQueue     string // 默认队列
	ResultBackend    string // 任务状态存储后台
	ResultsExpireIn  int    // 任务状态存储时间
	Exchange         string
	ExchangeType     string
	BindingKey       string
	PrefetchCount    int
	QueueBindingArgs map[string]interface{}
	NamedTaskFuncs   map[string]interface{}
}

// 阿里云AMQP消息队列
func NewAliAMQPMqQueue(c *AliAMQPConfig) (*MachineryQueue, error) {
	var broker = buildAliAMQPBroker(
		c.AccessKey,
		c.SecretKey,
		c.AliUid,
		c.EndPoint,
		c.VHost,
	)

	return NewRabbitMqQueue(
		broker,
		c.DefaultQueue,
		c.ResultBackend,
		c.ResultsExpireIn,
		c.Exchange,
		c.ExchangeType,
		c.BindingKey,
		c.PrefetchCount,
		c.QueueBindingArgs,
		c.NamedTaskFuncs,
	)
}

func buildAliAMQPBroker(accessKey string, secretKey string, aliUid int, endPoint string, vHost string) string {
	var bf bytes.Buffer
	bf.WriteString("amqp://")
	bf.WriteString(convertAliyunUserName(accessKey, aliUid))
	bf.WriteString(":")
	bf.WriteString(convertAliyunPassword(secretKey))
	bf.WriteString("@")
	bf.WriteString(endPoint)
	bf.WriteString("/")
	bf.WriteString(vHost)

	return bf.String()
}
