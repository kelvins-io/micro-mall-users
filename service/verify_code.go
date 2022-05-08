package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"gitee.com/cristiane/micro-mall-users/pkg/util/email"

	"gitee.com/cristiane/micro-mall-users/model/args"
	"gitee.com/cristiane/micro-mall-users/model/mysql"
	"gitee.com/cristiane/micro-mall-users/pkg/code"
	"gitee.com/cristiane/micro-mall-users/proto/micro_mall_users_proto/users"
	"gitee.com/cristiane/micro-mall-users/repository"
	"gitee.com/cristiane/micro-mall-users/vars"
	"gitee.com/kelvins-io/common/json"
	"gitee.com/kelvins-io/common/random"
	"gitee.com/kelvins-io/kelvins"
)

func GenVerifyCode(ctx context.Context, req *users.GenVerifyCodeRequest) int {
	retCode := code.Success
	var (
		err error
		//redis : new(repository.CheckVerifyCodeRedisLimiter)
		limiter = new(repository.CheckVerifyCodeRedisLimiter)
	)

	//Limits on the number of verification code requests and time interval
	limitKey := fmt.Sprintf("%s-%s-%d", req.CountryCode, req.Phone, req.BusinessType)
	limitRetCode, limitCount := checkVerifyCodeLimit(limiter, limitKey, vars.VerifyCodeSetting.SendPeriodLimitCount)
	if limitRetCode != code.Success {
		kelvins.ErrLogger.Infof(ctx, "checkVerifyCodeLimit %v %v is limited", req.CountryCode, req.Phone)
		retCode = limitRetCode
		return retCode
	}

	var uid int
	uid = int(req.Uid)
	userInfo, err := repository.GetUserByPhone("id,user_name,email", req.CountryCode, req.Phone)
	if err != nil {
		retCode = code.ErrorServer
		return retCode
	}
	if req.Uid <= 0 {
		if userInfo != nil {
			uid = userInfo.Id
		}
	}
	verifyCode := random.KrandNum(6)
	verifyCodeExpire := time.Now().Add(time.Duration(vars.VerifyCodeSetting.ExpireMinute) * time.Minute).Unix()
	verifyCodeRecord := mysql.VerifyCodeRecord{
		Uid:          uid,
		BusinessType: int(req.BusinessType),
		VerifyCode:   verifyCode,
		Expire:       int(verifyCodeExpire),
		CountryCode:  req.CountryCode,
		Phone:        req.Phone,
		Email:        req.Receiver,
		CreateTime:   time.Now(),
		UpdateTime:   time.Now(),
	}
	err = repository.CreateVerifyCodeRecord(&verifyCodeRecord)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "CreateVerifyCodeRecord err: %v, req: %v", err, json.MarshalToStringNoError(req))
		retCode = code.ErrorServer
		return retCode
	}

	key := fmt.Sprintf("%s-%s-%d", req.CountryCode, req.Phone, req.BusinessType)
	err = kelvins.G2CacheEngine.Set(verifyCodeCachePrefix+key, &verifyCodeRecord, 60*vars.VerifyCodeSetting.ExpireMinute, false)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "G2CacheEngine Set err: %v, key: %s,val: %v", err, key, json.MarshalToStringNoError(verifyCodeRecord))
		retCode = code.ErrorVerifyCodeInterval
		return retCode
	}

	err = limiter.SetVerifyCodeInterval(limitKey, vars.VerifyCodeSetting.SendIntervalExpireSecond)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "SetVerifyCodeInterval err: %v, req: %v", err, json.MarshalToStringNoError(req))
		retCode = code.ErrorVerifyCodeInterval
		return retCode
	}

	err = limiter.SetVerifyCodePeriodLimitCount(limitKey, limitCount+1, vars.VerifyCodeSetting.SendPeriodLimitExpireSecond)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "SetVerifyCodePeriodLimitCount err: %v, req: %v", err, json.MarshalToStringNoError(req))
		retCode = code.ErrorVerifyCodeInterval
		return retCode
	}

	notice := fmt.Sprintf(args.VerifyCodeTemplate, userInfo.UserName, verifyCode, args.GetMsg(int(req.BusinessType)), vars.VerifyCodeSetting.ExpireMinute)
	if req.Receiver != "" {
		for _, receiver := range strings.Split(req.Receiver, ",") {
			err = email.SendEmailNotice(ctx, receiver, kelvins.AppName, notice)
			if err != nil {
				kelvins.ErrLogger.Errorf(ctx, "SendEmailNotice err %v,receiver:%v, emailNotice: %v", err, receiver, notice)
			}
		}
	} else {
		err = email.SendEmailNotice(ctx, userInfo.Email, kelvins.AppName, notice)
		if err != nil {
			kelvins.ErrLogger.Errorf(ctx, "SendEmailNotice err %v,receiver:%v, emailNotice: %v", err, userInfo.Email, notice)
		}
	}

	return retCode
}

func checkVerifyCode(ctx context.Context, req *checkVerifyCodeArgs) int {
	key := fmt.Sprintf("%s-%s-%d", req.countryCode, req.phone, req.businessType)
	limiter := new(repository.CheckVerifyCodeRedisLimiter)
	err := limiter.CheckVerifyState(key)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "CheckVerifyState err: %v, key: %v", err, key)
		switch err {
		case repository.VerifyFailureForbidden:
			return code.ErrorVerifyCodeForbidden
		default:
			return code.ErrorServer
		}
	}

	var obj mysql.VerifyCodeRecord
	err = kelvins.G2CacheEngine.Get(verifyCodeCachePrefix+key, 60*vars.VerifyCodeSetting.ExpireMinute, &obj, func() (interface{}, error) {
		record, err := repository.GetVerifyCode(req.businessType, req.countryCode, req.phone, req.verifyCode)
		return record, err
	})
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "GetVerifyCode err: %v, req: %v", err, json.MarshalToStringNoError(req))
		return code.ErrorServer
	}

	if obj.VerifyCode != req.verifyCode {
		err := limiter.VerifyFailure(key)
		if err != nil {
			kelvins.ErrLogger.Errorf(ctx, "VerifyFailure err: %v, key: %v", err, key)
			switch err {
			case repository.VerifyFailureForbidden:
				return code.ErrorVerifyCodeForbidden
			default:
				return code.ErrorServer
			}
		}
		return code.ErrorVerifyCodeInvalid
	}

	if int64(obj.Expire) < time.Now().Unix() {
		return code.ErrorVerifyCodeExpire
	}

	return code.Success
}

type checkVerifyCodeArgs struct {
	businessType                   int
	countryCode, phone, verifyCode string
}

const verifyCodeCachePrefix = "micro-mall-users:verify-code:"

func checkVerifyCodeLimit(limiter repository.CheckVerifyCodeLimiter, key string, limitCount int) (int, int) {
	if limitCount <= 0 {
		limitCount = repository.DefaultVerifyCodeSendPeriodLimitCount
	}
	count, err := limiter.GetVerifyCodePeriodLimitCount(key)
	if err != nil {
		return code.ErrorServer, count
	}
	if count >= limitCount {
		return code.ErrorVerifyCodeLimited, count
	}

	intervalTime, err := limiter.GetVerifyCodeInterval(key)
	if err != nil {
		return code.ErrorServer, count
	}
	if intervalTime == 0 {
		return code.Success, count
	}

	return code.ErrorVerifyCodeInterval, count
}
