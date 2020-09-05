package event

import (
	"context"
	"errors"
	"fmt"
	sdk_errors "github.com/aliyun/alibaba-cloud-sdk-go/sdk/errors"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ons"
	"golang.org/x/time/rate"
	"regexp"
	"strings"
	"sync"
)

type aliOnsClient struct {
	*ons.Client

	instanceId string
	topicMap   sync.Map
	groupMap   sync.Map
	waiter     *rate.Limiter

	topicLocker chan struct{}
	groupLocker chan struct{}
}

var (
	onsClient *aliOnsClient
	reg       *regexp.Regexp
)

const EXPR = "^[0-9a-zA-Z-_]{3,64}$"

func initOnsClient(regionId, accessKey, secretKey, instanceId string) error {
	c, err := ons.NewClientWithAccessKey(regionId, accessKey, secretKey)
	if err != nil {
		return err
	}

	onsClient = &aliOnsClient{Client: c, instanceId: instanceId}
	onsClient.waiter = rate.NewLimiter(1, 1) // 阿里云请求限制，1 秒 1 次

	reg, err = regexp.Compile(EXPR)
	if err != nil {
		return err
	}

	onsClient.topicLocker = make(chan struct{}, 1)
	onsClient.groupLocker = make(chan struct{}, 1)

	return nil
}

// 获取 topic
func (c *aliOnsClient) getOnsTopic(businessName, topic string) (string, error) {
	onsClient.topicLocker <- struct{}{}
	defer func() {
		<-onsClient.topicLocker
	}()

	topic = strings.ToLower(topic)
	if !strings.Contains(topic, "-") {
		topic = fmt.Sprintf("%s-%s", businessName, topic)
	}
	if !reg.MatchString(topic) {
		return "", fmt.Errorf("不符合规范 [%s] 的 topic : %s", EXPR, topic)
	}

	if _, ok := c.topicMap.Load(topic); ok {
		return topic, nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c.waiter.Wait(ctx)

	// 查询 topic
	topicStatusReq := ons.CreateOnsTopicStatusRequest()
	topicStatusReq.InstanceId = c.instanceId
	topicStatusReq.Topic = topic
	_, err := c.OnsTopicStatus(topicStatusReq)
	if err != nil {
		return "", errors.New(err.Error() + fmt.Sprintf(", topic: %s", topic))
	}

	c.topicMap.Store(topic, true)

	return topic, nil
}

// 获取 groupId
func (c *aliOnsClient) getGroupId(businessName, consumerName string) (string, error) {
	onsClient.groupLocker <- struct{}{}
	defer func() {
		<-onsClient.groupLocker
	}()

	consumerName = strings.ToLower(consumerName)
	if !strings.Contains(consumerName, "-") {
		consumerName = fmt.Sprintf("%s-%s", businessName, consumerName)
	}

	var groupId = fmt.Sprintf("GID_%s", consumerName)
	if !reg.MatchString(groupId) {
		return "", fmt.Errorf("不符合规范 [%s] 的 groupId：%s", EXPR, groupId)
	}

	if groupId, ok := c.groupMap.Load(consumerName); ok {
		return groupId.(string), nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c.waiter.Wait(ctx)

	// 创建 groupId
	groupCreateReq := ons.CreateOnsGroupCreateRequest()
	groupCreateReq.InstanceId = c.instanceId
	groupCreateReq.GroupId = groupId
	groupCreateReq.Remark = businessName
	_, err := c.OnsGroupCreate(groupCreateReq)

	serverErr, ok := err.(*sdk_errors.ServerError)
	if err != nil && (!ok || (serverErr.ErrorCode() != "BIZ_SUBSCRIPTION_EXISTED")) {
		return "", errors.New(err.Error() + fmt.Sprintf(", %s", groupId))
	}

	c.groupMap.Store(consumerName, groupId)

	return groupId, nil
}
