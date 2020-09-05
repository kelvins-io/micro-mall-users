package event

import (
	"context"
	"encoding/json"

	"gitee.com/kelvins-io/common/log"
	mq_http_sdk "github.com/aliyunmq/mq-http-go-sdk"
)

type EventProducer struct {
	mq_http_sdk.MQProducer
	logger log.LoggerContextIface
}

func (p *EventProducer) publish(tag, key string, body interface{}) (string, error) {
	bs, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	msg := mq_http_sdk.PublishMessageRequest{
		MessageBody: string(bs),          // 消息内容
		MessageTag:  tag,                 // 消息标签
		MessageKey:  key,                 // 消息KEY
		Properties:  map[string]string{}, // 消息属性
	}

	var ret mq_http_sdk.PublishMessageResponse
	for i := 0; i < publishRetryCnt; i++ {
		ret, err = p.PublishMessage(msg)
		if err == nil {
			break
		}

		if i == publishRetryCnt-1 {
			return "", err
		}
	}

	p.logger.Infof(context.Background(), "发布消息成功：%s, %s, %s, %v, 消息 ID：%s", p.TopicName(), tag, key, body, ret.MessageId)

	return ret.MessageId, nil
}
