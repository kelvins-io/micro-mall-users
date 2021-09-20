package service

import (
	"context"
	"fmt"
	"gitee.com/cristiane/micro-mall-users/model/args"
	"gitee.com/cristiane/micro-mall-users/model/mysql"
	"gitee.com/cristiane/micro-mall-users/pkg/code"
	"gitee.com/cristiane/micro-mall-users/pkg/util/email"
	"gitee.com/cristiane/micro-mall-users/proto/micro_mall_users_proto/users"
	"gitee.com/cristiane/micro-mall-users/repository"
	"gitee.com/cristiane/micro-mall-users/vars"
	"gitee.com/kelvins-io/common/errcode"
	"gitee.com/kelvins-io/common/json"
	"gitee.com/kelvins-io/kelvins"
	"github.com/google/uuid"
	"strings"
	"time"
)

func MerchantsMaterial(ctx context.Context, req *users.MerchantsMaterialRequest) (int, int) {
	var merchantId int
	exist, err := repository.CheckUserExistById(int(req.Info.Uid))
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "CheckUserExistById err: %v,req : %v", err, json.MarshalToStringNoError(req))
		return merchantId, code.ErrorServer
	}
	if !exist {
		return merchantId, code.UserNotExist
	}

	if req.OperationType == users.OperationType_CREATE {
		merchantMaterial := mysql.Merchant{
			Uid:          int(req.Info.Uid),
			MerchantCode: uuid.New().String(),
			RegisterAddr: req.Info.RegisterAddr,
			HealthCardNo: req.Info.HealthCardNo,
			Identity:     int(req.Info.Identity),
			State:        int(req.Info.State),
			TaxCardNo:    req.Info.TaxCardNo,
			CreateTime:   time.Now(),
			UpdateTime:   time.Now(),
		}
		err := repository.CreateMerchantsMaterial(&merchantMaterial)
		if err != nil {
			if strings.Contains(err.Error(), errcode.GetErrMsg(code.DBDuplicateEntry)) {
				return merchantId, code.MerchantExist
			}
			kelvins.ErrLogger.Errorf(ctx, "CreateMerchantsMaterial err: %v,merchantMaterial:%v", err, json.MarshalToStringNoError(merchantMaterial))
			return merchantId, code.ErrorServer
		}
		record, err := repository.GetMerchantIdByUid(int(req.Info.Uid))
		if err != nil {
			kelvins.ErrLogger.Errorf(ctx, "GetMerchantsMaterialByUid err: %v,uid : %v", err, req.Info.Uid)
			return merchantId, code.ErrorServer
		}
		merchantId = record.MerchantId

		kelvins.GPool.SendJob(func() {
			u, ret := GetUserInfo(ctx, int(req.Info.Uid))
			if ret == code.Success {
				emailNotice := fmt.Sprintf(args.UserApplyMerchantTemplate, u.UserName, time.Now(), req.Info.RegisterAddr)
				if vars.EmailNoticeSetting != nil && vars.EmailNoticeSetting.Receivers != nil {
					for _, receiver := range vars.EmailNoticeSetting.Receivers {
						err = email.SendEmailNotice(ctx, receiver, kelvins.AppName, emailNotice)
						if err != nil {
							kelvins.ErrLogger.Info(ctx, "SendEmailNotice err, emailNotice: %v", emailNotice)
							return
						}
					}
				}
			}
		})

		return merchantId, code.Success
	} else if req.OperationType == users.OperationType_UPDATE {
		query := map[string]interface{}{
			"uid": req.Info.Uid,
		}
		maps := map[string]interface{}{
			"register_addr":  req.Info.RegisterAddr,
			"health_card_no": req.Info.HealthCardNo,
			"identity":       req.Info.Identity,
			"state":          req.Info.State,
			"tax_card_no":    req.Info.TaxCardNo,
			"update_time":    time.Now(),
		}
		err := repository.UpdateMerchantsMaterial(query, maps)
		if err != nil {
			kelvins.ErrLogger.Errorf(ctx, "UpdateMerchantsMaterial err: %v,query : %+v, maps: %+v", err, query, maps)
			return merchantId, code.ErrorServer
		}

		kelvins.GPool.SendJob(func() {
			u, ret := GetUserInfo(ctx, int(req.Info.Uid))
			if ret == code.Success {
				emailNotice := fmt.Sprintf(args.UserModifyMerchantInfoTemplate, u.UserName, time.Now())
				if vars.EmailNoticeSetting != nil && vars.EmailNoticeSetting.Receivers != nil {
					for _, receiver := range vars.EmailNoticeSetting.Receivers {
						err = email.SendEmailNotice(ctx, receiver, kelvins.AppName, emailNotice)
						if err != nil {
							kelvins.ErrLogger.Info(ctx, "SendEmailNotice err, emailNotice: %v", emailNotice)
							return
						}
					}
				}
			}
		})
		return merchantId, code.Success
	}
	return merchantId, code.Success
}

func GetMerchantsMaterial(ctx context.Context, req *users.GetMerchantsMaterialRequest) (*mysql.Merchant, int) {
	merchantInfo, err := repository.GetMerchantsMaterial(int(req.MaterialId))
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "GetMerchantsMaterialByUid err: %v,MaterialId : %v", err, req.GetMaterialId())
		return merchantInfo, code.ErrorServer
	}
	return merchantInfo, code.Success
}
