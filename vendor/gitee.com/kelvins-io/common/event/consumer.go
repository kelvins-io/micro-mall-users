package event

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	"gitee.com/kelvins-io/common/log"
	mq_http_sdk "github.com/aliyunmq/mq-http-go-sdk"
	mq_errors "github.com/gogap/errors"
)

type EventConsumer struct {
	mq_http_sdk.MQConsumer
	logger  log.LoggerContextIface
	handler interface{}
}

func (c *EventConsumer) handleMessages(consumer *EventConsumer, respChan chan mq_http_sdk.ConsumeMessageResponse, waitChan chan int, errChan chan error) {
	defer func() {
		waitChan <- 1
	}()

	select {
	case resp := <-respChan:
		{
			// 处理业务逻辑
			var handles []string
			var wg sync.WaitGroup
			var mutex sync.Mutex
			for _, msg := range resp.Messages {
				wg.Add(1)

				go func(entry mq_http_sdk.ConsumeMessageEntry) {
					defer wg.Done()

					err := c.handleMessage(consumer, entry)
					if err == nil {
						mutex.Lock()
						handles = append(handles, entry.ReceiptHandle)
						mutex.Unlock()
					}
				}(msg)
			}
			wg.Wait()

			if len(handles) == 0 {
				return
			}

			err := consumer.AckMessage(handles)
			if err == nil {
				c.logger.Infof(context.Background(), "%s, 消息回传 Ack 成功：%v", c.getKey(), handles)
			} else {
				errs := err.(mq_errors.ErrCode).Context()["Detail"].([]mq_http_sdk.ErrAckItem)
				for _, errAckItem := range errs {
					c.logger.Errorf(context.Background(), "%s, 消息回传 Ack 失败：%v", c.getKey(), errAckItem)
				}
			}
		}
	case err := <-errChan:
		{
			isNotExist := strings.Contains(err.(mq_errors.ErrCode).Error(), "MessageNotExist")
			if !isNotExist {
				c.logger.Errorf(context.Background(), "%s, 拉取消息失败：%v", c.getKey(), err)
			}

			// 没有新的消息
			time.Sleep(time.Duration(1) * time.Second)
		}
	case <-time.After(timeoutSeconds * time.Second):
		{
			c.logger.Errorf(context.Background(), "%s, 消费消息超时: %d秒", c.getKey(), timeoutSeconds)
		}
	}
}

func (c *EventConsumer) handleMessage(consumer *EventConsumer, msg mq_http_sdk.ConsumeMessageEntry) (res error) {
	defer func() {
		if err := recover(); err != nil {
			msg := fmt.Sprintf("%s recover from panic, msg: %v", c.getKey(), err)
			c.logger.Error(context.Background(), msg)
			res = errors.New(msg)
		}
	}()

	// 消息转换
	rMethod := reflect.ValueOf(consumer.handler)
	bodyType := rMethod.Type().In(0)
	body := reflect.New(bodyType).Interface()
	err := json.Unmarshal([]byte(msg.MessageBody), body)
	if err != nil {
		c.logger.Errorf(context.Background(), "%s, 消息转义错误：%v", c.getKey(), err)
		return err
	}

	// 消费消息
	body = reflect.ValueOf(body).Elem().Interface()
	ret := rMethod.Call([]reflect.Value{reflect.ValueOf(body)})
	if ret[0].Interface() != nil {
		err := ret[0].Interface().(error)
		c.logger.Errorf(context.Background(), "%s, 执行消费消息失败：%v", c.getKey(), err)
		return err
	}

	c.logger.Infof(context.Background(), "%s, 执行消费消息成功：%v", c.getKey(), msg.ReceiptHandle)

	return nil
}

func (c *EventConsumer) getKey() string {
	return fmt.Sprintf("%s_%s_%s", c.TopicName(), c.MessageTag(), c.Consumer())
}
