package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"gitee.com/cristiane/micro-mall-users/model/args"
	"gitee.com/cristiane/micro-mall-users/model/mysql"
	"gitee.com/cristiane/micro-mall-users/pkg/code"
	"gitee.com/cristiane/micro-mall-users/pkg/util"
	"gitee.com/cristiane/micro-mall-users/proto/micro_mall_search_proto/search_business"
	"gitee.com/cristiane/micro-mall-users/proto/micro_mall_users_proto/users"
	"gitee.com/cristiane/micro-mall-users/repository"
	"gitee.com/cristiane/micro-mall-users/vars"
	"gitee.com/kelvins-io/common/errcode"
	"gitee.com/kelvins-io/common/json"
	"gitee.com/kelvins-io/kelvins"
	"github.com/google/uuid"
)

func MerchantsMaterial(ctx context.Context, req *users.MerchantsMaterialRequest) (merchantId int, retCode int) {
	retCode = code.Success
	exist, err := repository.CheckUserExistById(int(req.Info.Uid))
	if err != nil {
		retCode = code.ErrorServer
		kelvins.ErrLogger.Errorf(ctx, "CheckUserExistById err: %v,req : %v", err, json.MarshalToStringNoError(req))
		return
	}
	if !exist {
		retCode = code.UserNotExist
		return
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
			retCode = code.ErrorServer
			if strings.Contains(err.Error(), errcode.GetErrMsg(code.DBDuplicateEntry)) {
				retCode = code.MerchantExist
				return
			}
			kelvins.ErrLogger.Errorf(ctx, "CreateMerchantsMaterial err: %v,merchantMaterial:%v", err, json.MarshalToStringNoError(merchantMaterial))
			return
		}
		record, err := repository.GetMerchantIdByUid(int(req.Info.Uid))
		if err != nil {
			retCode = code.ErrorServer
			kelvins.ErrLogger.Errorf(ctx, "GetMerchantsMaterialByUid err: %v,uid : %v", err, req.Info.Uid)
			return
		}
		merchantId = record.MerchantId
		// 事件通知
		merchantsMaterialEventNotice(&args.MerchantInfoSearch{
			Uid:          req.GetInfo().GetUid(),
			MerchantCode: merchantCode,
			RegisterAddr: req.GetInfo().GetRegisterAddr(),
			HealthCardNo: req.GetInfo().GetHealthCardNo(),
			TaxCardNo:    req.GetInfo().GetTaxCardNo(),
		}, "create")

		return
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
		err = repository.UpdateMerchantsMaterial(query, maps)
		if err != nil {
			retCode = code.ErrorServer
			kelvins.ErrLogger.Errorf(ctx, "UpdateMerchantsMaterial err: %v,query : %+v, maps: %+v", err, query, maps)
			return
		}
		// 事件通知
		merchantsMaterialEventNotice(&args.MerchantInfoSearch{
			Uid:          req.GetInfo().GetUid(),
			RegisterAddr: req.GetInfo().GetRegisterAddr(),
			HealthCardNo: req.GetInfo().GetHealthCardNo(),
			TaxCardNo:    req.GetInfo().GetTaxCardNo(),
		}, "update")

		return
	}
	return
}

func GetMerchantsMaterial(ctx context.Context, req *users.GetMerchantsMaterialRequest) (*users.MerchantsMaterialInfo, int) {
	result := &users.MerchantsMaterialInfo{}
	merchantInfo, err := repository.GetMerchantsMaterial(int(req.MaterialId))
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "GetMerchantsMaterialByUid err: %v,MaterialId : %v", err, req.GetMaterialId())
		return result, code.ErrorServer
	}

	result = &users.MerchantsMaterialInfo{
		Uid:          int64(merchantInfo.Uid),
		MaterialId:   int64(merchantInfo.MerchantId),
		RegisterAddr: merchantInfo.RegisterAddr,
		HealthCardNo: merchantInfo.HealthCardNo,
		Identity:     int32(merchantInfo.Identity),
		State:        int32(merchantInfo.State),
		TaxCardNo:    merchantInfo.TaxCardNo,
		CreateTime:   util.ParseTimeOfStr(merchantInfo.CreateTime.Unix()),
		UpdateTime:   util.ParseTimeOfStr(merchantInfo.UpdateTime.Unix()),
	}

	if merchantInfo.Uid > 0 {
		userByUid, err := repository.GetUserByUid("id,user_name,email", merchantInfo.Uid)
		if err != nil {
			kelvins.ErrLogger.Errorf(ctx, "GetUserByUid err: %v,uid : %v", err, merchantInfo.Uid)
			return result, code.ErrorServer
		}
		result.MerchantName = userByUid.UserName
		result.MerchantEmail = userByUid.Email
	}

	return result, code.Success
}

func merchantsMaterialEventNotice(info *args.MerchantInfoSearch, operationType string) {
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
		if info.MerchantCode == "" {
			merchantIdByUid, err := repository.GetMerchantIdByUid(int(info.Uid))
			if err != nil {
				kelvins.ErrLogger.Errorf(ctx, "GetMerchantIdByUid err: %v,uid : %v", err, info.Uid)
				return
			}
			info.MerchantCode = merchantIdByUid.MerchantCode
		}

		// 1 搜索事件
		userInfoMsg := args.CommonBusinessMsg{
			Type:    args.MerchantInfoSearchNoticeType,
			Tag:     args.GetMsg(args.MerchantInfoSearchNoticeType),
			UUID:    uuid.New().String(),
			Content: json.MarshalToStringNoError(info),
		}
		vars.QueueServerUserInfoSearchPusher.PushMessage(ctx, &userInfoMsg)

		// 2 用户状态事件
		businessMsg := args.CommonBusinessMsg{
			Type: args.UserStateEventTypeMerchantInfo,
			Tag:  args.GetMsg(args.UserStateEventTypeMerchantInfo),
			UUID: genUUID(),
			Time: util.ParseTimeOfStr(time.Now().Unix()),
			Content: json.MarshalToStringNoError(args.UserStateNotice{
				Uid: int(info.Uid),
				Extra: map[string]string{
					"operation_type": operationType,
					"merchant_code":  info.MerchantCode,
				},
			}),
		}
		pushUserStateNoticeService.PushMessage(ctx, businessMsg)
	})
}

func SearchMerchantInfo(ctx context.Context, query string) (result []*users.SearchMerchantsInfoEntry, retCode int) {
	result = make([]*users.SearchMerchantsInfoEntry, 0)
	retCode = code.Success
	searchKey := "micro-mall-users:search-merchant:" + query
	err := kelvins.G2CacheEngine.Get(searchKey, 120, &result, func() (interface{}, error) {
		ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
		defer cancel()
		list, ret := searchMerchantInfo(ctx, query)
		if ret != code.Success {
			return &list, fmt.Errorf("searchMerchantInfo ret %v", ret)
		}
		return &list, nil
	})
	if err != nil {
		retCode = code.ErrorServer
		return
	}
	return
}

func searchMerchantInfo(ctx context.Context, query string) (result []*users.SearchMerchantsInfoEntry, retCode int) {
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
