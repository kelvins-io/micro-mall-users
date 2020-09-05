package event

import (
	"context"
	"errors"
	"fmt"
	"gitee.com/kelvins-io/common/log"
	mq_http_sdk "github.com/aliyunmq/mq-http-go-sdk"
	"math/rand"
	"reflect"
	"regexp"
	"strings"
	"sync"
	"time"
)

var (
	publishRetryCnt int           = 3
	numOfMessages   int32         = 3
	waitSeconds     int64         = 3
	timeoutSeconds  time.Duration = 35
)

type Config struct {
	BusinessName string
	RegionId     string
	AccessKey    string
	SecretKey    string
	InstanceId   string
	HttpEndpoint string
}

type Message struct {
	// 话题
	Topic string
	// 标签
	Tag string
	// 消息键
	Key string
	// 消息体
	Body interface{}
}

type ProducerIface interface {
	Publish(msg *Message) (string, error)
}

type EventServerIface interface {
	Subscribe(toipc, tag, consumer string, handler interface{}) error
}

type EventServer struct {
	*mq_http_sdk.AliyunMQClient
	businessName string
	instanceId   string
	once         sync.Once
	producerMap  sync.Map
	consumerMap  sync.Map
	logger       log.LoggerContextIface
}

func NewEventServer(config *Config, logger log.LoggerContextIface) (*EventServer, error) {
	if config.AccessKey == "" {
		return nil, errors.New("缺少参数：config.AccessKey")
	}
	if config.SecretKey == "" {
		return nil, errors.New("缺少参数：config.SecretKey")
	}
	if config.HttpEndpoint == "" {
		return nil, errors.New("缺少参数：config.HttpEndpoint")
	}
	if config.InstanceId == "" {
		return nil, errors.New("缺少参数：config.InstanceId")
	}
	if config.BusinessName == "" {
		return nil, errors.New("缺少参数：config.BusinessName")
	}

	expr := "^[a-zA-Z]{1,12}$"
	reg, _ := regexp.Compile(expr)
	if !reg.MatchString(config.BusinessName) {
		return nil, fmt.Errorf("不符合规范 [%s] 的业务名称：%s", expr, config.BusinessName)
	}

	rand.Seed(time.Now().UnixNano())

	// 初始化 ons 客户端
	err := initOnsClient(config.RegionId, config.AccessKey, config.SecretKey, config.InstanceId)
	if err != nil {
		return nil, err
	}

	// 初始化 rocketMQ 客户端
	client := mq_http_sdk.NewAliyunMQClient(config.HttpEndpoint, config.AccessKey, config.SecretKey, "")

	server := &EventServer{
		AliyunMQClient: client.(*mq_http_sdk.AliyunMQClient),
		instanceId:     config.InstanceId,
		businessName:   config.BusinessName,
		logger:         &log.EmptyLoggerContext{},
	}
	if logger != nil {
		server.logger = logger
	}

	return server, nil
}

// 发布消息
func (h *EventServer) Publish(msg *Message) (string, error) {
	if msg.Topic == "" {
		return "", errors.New("topic 不能为空")
	}
	if msg.Tag == "" {
		return "", errors.New("tag 不能为空")
	}
	if msg.Body == nil {
		return "", errors.New("body 不能为空")
	}

	// 获取阿里云 topic
	topic, err := onsClient.getOnsTopic(h.businessName, msg.Topic)
	if err != nil {
		h.logger.Errorf(context.Background(), "获取 topic: %s 失败: %v", msg.Topic, err)

		return "", err
	}
	if !strings.Contains(topic, h.businessName) {
		return "", errors.New("不能发布其它业务线的消息")
	}

	// 创建或获取生产者
	var producer *EventProducer
	if p, ok := h.producerMap.Load(topic); ok {
		producer = p.(*EventProducer)
	} else {
		producer = &EventProducer{
			MQProducer: h.GetProducer(h.instanceId, topic),
			logger:     h.logger,
		}
		h.producerMap.Store(topic, producer)
	}

	return producer.publish(msg.Tag, msg.Key, msg.Body)
}

// 订阅消息
func (h *EventServer) Subscribe(toipc, tag, consumer string, handler interface{}) error {
	if toipc == "" {
		return errors.New("topic 不能为空")
	}
	if tag == "" {
		return errors.New("tag 不能为空")
	}
	if consumer == "" {
		return errors.New("消费者名称不能为空")
	}
	if handler == nil {
		return errors.New("处理方法不能为空")
	}

	// 判断类型
	rHandler := reflect.ValueOf(handler)
	if rHandler.Type().NumIn() != 1 {
		return errors.New("处理方法只允许一个入参")
	}
	if rHandler.Type().NumOut() != 1 || rHandler.Type().Out(0).String() != "error" {
		return errors.New("处理方法返回值必须是一个 error 类型")
	}

	// 获取阿里云 topic
	toipc, err := onsClient.getOnsTopic(h.businessName, toipc)
	if err != nil {
		h.logger.Errorf(context.Background(), "获取 topic 失败：%v", err)
		return err
	}

	// 获取阿里云 group id
	groupId, err := onsClient.getGroupId(h.businessName, consumer)
	if err != nil {
		h.logger.Errorf(context.Background(), "获取 groupId 失败：%v", err)
		return err
	}

	key := fmt.Sprintf("%s_%s", groupId, toipc)

	// 消息者不能重复订阅，否则会出现消息丢失
	if _, exist := h.consumerMap.Load(key); exist {
		return fmt.Errorf("消费者: %s 已订阅消息, 消费者的消费行为必须保持一致！", groupId)
	}

	// 添加消费者
	h.consumerMap.Store(key, &EventConsumer{
		MQConsumer: h.GetConsumer(h.instanceId, toipc, groupId, tag),
		logger:     h.logger,
		handler:    handler,
	})

	h.logger.Infof(context.Background(), "消费者: %s 订阅 topic：%s, tag: %s", groupId, toipc, tag)

	return nil
}

// 开始事件处理
func (h *EventServer) Start() error {
	h.once.Do(func() {
		h.consumerMap.Range(func(key, value interface{}) bool {
			// 错峰请求阿里云
			time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)

			h.logger.Infof(context.Background(), "开始轮询消息：%s", key)

			go func(consumer *EventConsumer) {
				for {
					waitChan := make(chan int)
					errChan := make(chan error)
					respChan := make(chan mq_http_sdk.ConsumeMessageResponse)

					// 处理消息
					go consumer.handleMessages(consumer, respChan, waitChan, errChan)

					// 拉取消息
					go consumer.ConsumeMessage(respChan, errChan, numOfMessages, waitSeconds)

					<-waitChan
				}
			}(value.(*EventConsumer))

			return true
		})
	})

	return nil
}
