package repository

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"gitee.com/kelvins-io/common/hash"

	"gitee.com/kelvins-io/kelvins"

	"github.com/gomodule/redigo/redis"
)

//@title Request verification code limit.
//@description Time interval for request verification code and limit on the number of times within the time period .
//Provide redis and local map cache the number of consecutive requests for users in a period of time
//and the expiration flag of the interval between each request.
//key
//@author xhylgogo
//@modify time 2021-04-15

const (
	VerifyCodePeriodLimitCountKeyPrefix          = "micro-mall-users:verify_code_period_limit_count:"
	VerifyCodeIntervalKeyPrefix                  = "micro-mall-users:verify_code_interval:"
	VerifyCodeFailureKey                         = "micro-mall-users:verify_code_failure:"
	VerifyCodeFailureMax                         = 3
	DefaultVerifyCodeSendPeriodLimitCount        = 10
	DefaultVerifyCodeSendPeriodLimitExpireSecond = 3600
	DefaultVerifyCodeSendIntervalExpireSecond    = 60
)

var (
	VerifyFailureForbidden = fmt.Errorf("verify code failure too many")
)

type CheckVerifyCodeLimiter interface {
	//Accumulative number of verification code requests during the acquisition period
	GetVerifyCodePeriodLimitCount(key string) (int, error)
	//The cumulative number of verification code requests within the set time period
	SetVerifyCodePeriodLimitCount(key string, limitCount int, expireTime int64) error
	//The remaining time of the next request for verification code within the time interval
	GetVerifyCodeInterval(key string) (int64, error)
	//The remaining time of the next request for verification code within the set time interval
	SetVerifyCodeInterval(key string, intervalTime int64) error
	// VerifyFailure is verify code failure
	VerifyFailure(key string) error
	// CheckVerifyState is verify code failure max
	CheckVerifyState(key string) error
}

type CheckVerifyCodeRedisLimiter struct {
}

func (c *CheckVerifyCodeRedisLimiter) GetVerifyCodePeriodLimitCount(key string) (int, error) {
	conn := kelvins.RedisConn.Get()
	defer conn.Close()
	key = fmt.Sprintf("%s%s", VerifyCodePeriodLimitCountKeyPrefix, key)
	key = hash.MD5EncodeToString(key)
	count, err := redis.Int(conn.Do("GET", key))
	if err != nil && err != redis.ErrNil {
		return 0, err
	}
	return count, nil
}

func (c *CheckVerifyCodeRedisLimiter) VerifyFailure(key string) error {
	frequency, err := c.getVerifyFailureFrequency(key)
	if err != nil {
		return err
	}
	if frequency >= VerifyCodeFailureMax {
		return VerifyFailureForbidden
	}
	frequency++
	conn := kelvins.RedisConn.Get()
	defer conn.Close()
	key2 := fmt.Sprintf("%v%v", VerifyCodeFailureKey, key)
	key2 = hash.MD5EncodeToString(key2)
	_, err = conn.Do("SET", key2, fmt.Sprintf("%d", frequency))
	if err != nil {
		return err
	}

	expireTime := 24 * 3600
	_, err = conn.Do("EXPIRE", key2, expireTime)
	if err != nil {
		_, err := redis.Bool(conn.Do("DEL", key))
		if err != nil {
			return err
		}
		return err
	}
	return nil
}

func (c *CheckVerifyCodeRedisLimiter) CheckVerifyState(key string) error {
	frequency, err := c.getVerifyFailureFrequency(key)
	if err != nil {
		return err
	}
	if frequency >= VerifyCodeFailureMax {
		return VerifyFailureForbidden
	}
	return nil
}

func (c *CheckVerifyCodeRedisLimiter) getVerifyFailureFrequency(key string) (int, error) {
	var frequency int
	conn := kelvins.RedisConn.Get()
	defer conn.Close()
	key2 := fmt.Sprintf("%v%v", VerifyCodeFailureKey, key)
	key2 = hash.MD5EncodeToString(key2)
	str, err := redis.String(conn.Do("GET", key2))
	if err != nil && err != redis.ErrNil {
		return 0, err
	}
	if str != "" && err != redis.ErrNil {
		frequency, err = strconv.Atoi(str)
		if err != nil {
			return 0, err
		}
	}
	return frequency, nil
}

func (c *CheckVerifyCodeRedisLimiter) SetVerifyCodePeriodLimitCount(key string, limitCount int, expireTime int64) error {
	if expireTime <= 0 {
		expireTime = DefaultVerifyCodeSendPeriodLimitExpireSecond
	}
	conn := kelvins.RedisConn.Get()
	defer conn.Close()
	key = fmt.Sprintf("%s%s", VerifyCodePeriodLimitCountKeyPrefix, key)
	key = hash.MD5EncodeToString(key)
	_, err := conn.Do("SET", key, limitCount)
	if err != nil {
		return err
	}
	_, err = conn.Do("EXPIRE", key, expireTime)
	if err != nil {
		_, err := redis.Bool(conn.Do("DEL", key))
		if err != nil {
			return err
		}
		return err
	}
	return nil
}

func (c *CheckVerifyCodeRedisLimiter) GetVerifyCodeInterval(key string) (int64, error) {
	conn := kelvins.RedisConn.Get()
	defer conn.Close()
	key = fmt.Sprintf("%s%s", VerifyCodeIntervalKeyPrefix, key)
	key = hash.MD5EncodeToString(key)
	expireTime, err := redis.Int64(conn.Do("TTL", key))
	if err != nil || expireTime <= 0 {
		return 0, err
	}
	return expireTime, nil
}

func (c *CheckVerifyCodeRedisLimiter) SetVerifyCodeInterval(key string, intervalTime int64) error {
	if intervalTime <= 0 {
		intervalTime = DefaultVerifyCodeSendIntervalExpireSecond
	}
	conn := kelvins.RedisConn.Get()
	defer conn.Close()
	key = fmt.Sprintf("%s%s", VerifyCodeIntervalKeyPrefix, key)
	key = hash.MD5EncodeToString(key)
	endTime := time.Now().Add(time.Duration(intervalTime) * time.Second).Unix()
	_, err := conn.Do("SET", key, endTime)
	if err != nil {
		return err
	}
	_, err = conn.Do("EXPIRE", key, intervalTime)
	if err != nil {
		_, err := redis.Bool(conn.Do("DEL", key))
		if err != nil {
			return err
		}
		return err
	}
	return nil
}

var (
	verifyCodeLimitedMapCache  = new(sync.Map)
	verifyCodeIntervalMapCache = new(sync.Map)
	verifyCodeFailureMapCache  = new(sync.Map)
)

type CheckVerifyCodeMapCacheLimiter struct {
}

type limitCacheModel struct {
	LimitCount int
	ExpireTime int64
}

func (c CheckVerifyCodeMapCacheLimiter) GetVerifyCodePeriodLimitCount(key string) (int, error) {
	limitCountInterface, ok := verifyCodeLimitedMapCache.Load(key)
	if !ok {
		return 0, nil
	}
	limitCount, ok := limitCountInterface.(limitCacheModel)
	if limitCount.ExpireTime <= time.Now().Unix() {
		verifyCodeLimitedMapCache.Delete(key)
		return 0, nil
	}
	return limitCount.LimitCount, nil
}

func (c CheckVerifyCodeMapCacheLimiter) SetVerifyCodePeriodLimitCount(key string, limitCount int, expireTime int64) error {
	if expireTime <= 0 {
		expireTime = DefaultVerifyCodeSendPeriodLimitExpireSecond
	}
	verifyCodeLimitedMapCache.Store(key, limitCacheModel{
		LimitCount: limitCount,
		ExpireTime: time.Now().Add(time.Duration(expireTime) * time.Second).Unix(),
	})
	return nil
}

func (c CheckVerifyCodeMapCacheLimiter) GetVerifyCodeInterval(key string) (int64, error) {
	endTimeInterface, ok := verifyCodeIntervalMapCache.Load(key)
	if !ok {
		return 0, nil
	}
	nowTime := time.Now().Unix()
	expireTime, ok := endTimeInterface.(int64)
	if !ok {
		verifyCodeIntervalMapCache.Delete(key)
		return 0, errors.New("GetVerifyCodeInterval error : endTime assert failed")
	}
	if expireTime <= nowTime {
		verifyCodeIntervalMapCache.Delete(key)
		return 0, nil
	}
	intervalTime := expireTime - nowTime
	return intervalTime, nil
}

func (c CheckVerifyCodeMapCacheLimiter) SetVerifyCodeInterval(key string, intervalTime int64) error {
	if intervalTime <= 0 {
		intervalTime = DefaultVerifyCodeSendIntervalExpireSecond
	}
	expireTime := time.Now().Add(time.Duration(intervalTime) * time.Second).Unix()
	verifyCodeIntervalMapCache.Store(key, expireTime)
	return nil
}

func (c CheckVerifyCodeMapCacheLimiter) VerifyFailure(key string) error {
	var frequency int
	v, ok := verifyCodeFailureMapCache.Load(key)
	if ok {
		frequency = v.(int)
	}
	if frequency >= VerifyCodeFailureMax {
		return VerifyFailureForbidden
	}
	frequency++
	verifyCodeFailureMapCache.Store(key, frequency)
	return nil
}

func (c CheckVerifyCodeMapCacheLimiter) CheckVerifyState(key string) error {
	var frequency int
	v, ok := verifyCodeFailureMapCache.Load(key)
	if ok {
		frequency = v.(int)
	}
	if frequency >= VerifyCodeFailureMax {
		return VerifyFailureForbidden
	}
	return nil
}
