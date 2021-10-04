package service

import (
	"context"
	"fmt"
	"gitee.com/cristiane/micro-mall-users/model/args"
	"gitee.com/cristiane/micro-mall-users/model/mysql"
	"gitee.com/cristiane/micro-mall-users/pkg/code"
	"gitee.com/cristiane/micro-mall-users/pkg/util"
	"gitee.com/cristiane/micro-mall-users/pkg/util/email"
	"gitee.com/cristiane/micro-mall-users/proto/micro_mall_search_proto/search_business"
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
		merchantCode := uuid.New().String()
		merchantMaterial := mysql.Merchant{
			Uid:          int(req.Info.Uid),
			MerchantCode: merchantCode,
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

		// 商户申请通知
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

		// 搜索通知
		merchantsMaterialSearchNotice(&args.MerchantInfoSearch{
			Uid:          req.GetInfo().GetUid(),
			MerchantCode: merchantCode,
			RegisterAddr: req.GetInfo().GetRegisterAddr(),
			HealthCardNo: req.GetInfo().GetHealthCardNo(),
			TaxCardNo:    req.GetInfo().GetTaxCardNo(),
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

		// 搜索通知
		merchantsMaterialSearchNotice(&args.MerchantInfoSearch{
			Uid:          req.GetInfo().GetUid(),
			RegisterAddr: req.GetInfo().GetRegisterAddr(),
			HealthCardNo: req.GetInfo().GetHealthCardNo(),
			TaxCardNo:    req.GetInfo().GetTaxCardNo(),
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

func merchantsMaterialSearchNotice(info *args.MerchantInfoSearch) {
	kelvins.GPool.SendJob(func() {
		var ctx = context.TODO()
		var userName string
		user, ret := GetUserInfo(ctx, int(info.Uid))
		if ret == code.Success {
			if user != nil {
				userName = user.UserName
			}
		}
		info.UserName = userName
		userInfoMsg := args.CommonBusinessMsg{
			Type:    args.MerchantsMaterialInfoNoticeType,
			Tag:     args.GetMsg(args.MerchantsMaterialInfoNoticeType),
			UUID:    uuid.New().String(),
			Content: json.MarshalToStringNoError(info),
		}
		vars.QueueServerUserInfoSearchPusher.PushMessage(ctx, &userInfoMsg)
	})
}

func SearchMerchantInfo(ctx context.Context, query string) (result []*users.SearchMerchantsInfoEntry, retCode int) {
	result = make([]*users.SearchMerchantsInfoEntry, 0)
	retCode = code.Success
	serverName := args.RpcServiceMicroMallSearch
	conn, err := util.GetGrpcClient(ctx, serverName)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "GetGrpcClient %v err: %v", serverName, err)
		return result, code.ErrorServer
	}
	client := search_business.NewSearchBusinessServiceClient(conn)
	reqSearch := &search_business.MerchantInfoSearchRequest{Query: query}
	rsp, err := client.MerchantInfoSearch(ctx, reqSearch)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "MerchantInfoSearch err: %v, query: %v", err, query)
		return result, code.ErrorServer
	}
	if rsp.Common.Code != search_business.RetCode_SUCCESS {
		kelvins.ErrLogger.Errorf(ctx, "MerchantInfoSearch err: %v, query: %v, rsp: %v", err, query, json.MarshalToStringNoError(rsp))
		return result, code.ErrorServer
	}
	if len(rsp.List) == 0 {
		return
	}
	merchantCode := make([]string, len(rsp.List))
	for i := 0; i < len(rsp.List); i++ {
		merchantCode[i] = rsp.List[i].GetMerchantCode()
	}
	merchantInfoList, ret := FindMerchantInfo(ctx, merchantCode)
	if ret != code.Success {
		retCode = ret
		return
	}
	if len(merchantInfoList) == 0 {
		return
	}
	merchantCodeToMerchant := map[string]*mysql.Merchant{}
	for i := 0; i < len(merchantInfoList); i++ {
		merchantCodeToMerchant[merchantInfoList[i].MerchantCode] = merchantInfoList[i]
	}
	result = make([]*users.SearchMerchantsInfoEntry, 0)
	for i := 0; i < len(rsp.List); i++ {
		if rsp.List[i].MerchantCode == "" {
			continue
		}
		merchantInfo, ok := merchantCodeToMerchant[rsp.List[i].MerchantCode]
		if ok {
			entry := &users.SearchMerchantsInfoEntry{
				Info: &users.MerchantsMaterialInfo{
					Uid:          int64(merchantInfo.Uid),
					MaterialId:   int64(merchantInfo.MerchantId),
					RegisterAddr: merchantInfo.RegisterAddr,
					HealthCardNo: merchantInfo.HealthCardNo,
					Identity:     int32(merchantInfo.Identity),
					State:        int32(merchantInfo.State),
					TaxCardNo:    merchantInfo.TaxCardNo,
					CreateTime:   merchantInfo.CreateTime.Format(time.RFC3339),
				},
				Score: rsp.List[i].GetScore(),
			}
			result = append(result, entry)
		}
	}
	return
}

const sqlSelectFindMerchantInfo = "*"

func FindMerchantInfo(ctx context.Context, merchantCode []string) ([]*mysql.Merchant, int) {
	retCode := code.Success
	list, err := repository.FindMerchantInfo(sqlSelectFindMerchantInfo, merchantCode)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "FindMerchantInfo err：%v,merchantCode: %v", err, merchantCode)
		return list, code.ErrorServer
	}
	return list, retCode
}
