package email

import (
	"context"
	"strings"
	"sync"
	"time"

	"gitee.com/cristiane/micro-mall-users/vars"
	"gitee.com/kelvins-io/kelvins"
)

var (
	one                sync.Once
	DefaultEmailHelper *Client
)

func initVars() {
	DefaultEmailHelper = NewClient(vars.EmailConfigSetting.User, vars.EmailConfigSetting.Password, vars.EmailConfigSetting.Host, vars.EmailConfigSetting.Port)
}

const maxRetrySendTimes = 3
const retryIdleTime = 500 * time.Millisecond

func SendEmailNotice(ctx context.Context, receivers, subject, msg string) error {
	if vars.EmailConfigSetting == nil || !vars.EmailConfigSetting.Enable {
		return nil
	}
	if receivers == "" {
		return nil
	}
	one.Do(func() {
		initVars()
	})
	var err error
	emailReq := &SendRequest{
		Receivers: strings.Split(receivers, ","),
		Subject:   subject,
		Message:   msg,
	}

	// retry send email
	for retryCount := 0; retryCount < maxRetrySendTimes; retryCount++ {
		err = DefaultEmailHelper.SendEmail(emailReq)
		if err == nil {
			break
		}
		time.Sleep(retryIdleTime)
	}

	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "send email err: %v, req: %v", err, emailReq)
		return err
	}

	return nil
}
