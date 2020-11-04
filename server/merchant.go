package server

import (
	"context"
	"fmt"
	"gitee.com/cristiane/micro-mall-users/pkg/code"
	"gitee.com/cristiane/micro-mall-users/pkg/util"
	"gitee.com/cristiane/micro-mall-users/proto/micro_mall_users_proto/users"
	"gitee.com/cristiane/micro-mall-users/service"
	"gitee.com/kelvins-io/common/errcode"
)

type MerchantsServer struct{}

func NewMerchantsServer() users.MerchantsServiceServer {
	return new(MerchantsServer)
}

func (m *MerchantsServer) MerchantsMaterial(ctx context.Context, req *users.MerchantsMaterialRequest) (*users.MerchantsMaterialResponse, error) {
	var result users.MerchantsMaterialResponse
	result.Common = &users.CommonResponse{
		Code: 0,
		Msg:  "",
	}
	if req.Info.Uid <= 0 {
		result.Common.Code = users.RetCode_USER_NOT_EXIST
		result.Common.Msg = errcode.GetErrMsg(code.UserNotExist)
		return &result, nil
	}
	merchantId, retCode := service.MerchantsMaterial(ctx, req)
	result.MaterialId = int64(merchantId)
	if retCode != code.Success {
		if retCode == code.UserExist {
			result.Common.Code = users.RetCode_USER_EXIST
			result.Common.Msg = errcode.GetErrMsg(code.UserExist)
		} else if retCode == code.UserNotExist {
			result.Common.Code = users.RetCode_USER_NOT_EXIST
			result.Common.Msg = errcode.GetErrMsg(code.UserNotExist)
		} else if retCode == code.MerchantNotExist {
			result.Common.Code = users.RetCode_MERCHANT_NOT_EXIST
			result.Common.Msg = errcode.GetErrMsg(code.MerchantNotExist)
		} else if retCode == code.MerchantExist {
			result.Common.Code = users.RetCode_MERCHANT_EXIST
			result.Common.Msg = errcode.GetErrMsg(code.MerchantExist)
		} else {
			result.Common.Code = users.RetCode_ERROR
			result.Common.Msg = errcode.GetErrMsg(code.ErrorServer)
		}
	} else {
		result.Common.Code = users.RetCode_SUCCESS
		result.Common.Msg = errcode.GetErrMsg(code.Success)
	}

	return &result, nil
}

func (m *MerchantsServer) MerchantsMaterialAudit(ctx context.Context, req *users.MerchantsMaterialAuditRequest) (*users.MerchantsMaterialAuditResponse, error) {
	return &users.MerchantsMaterialAuditResponse{}, nil
}

func (m *MerchantsServer) GetMerchantsMaterial(ctx context.Context, req *users.GetMerchantsMaterialRequest) (*users.GetMerchantsMaterialResponse, error) {
	var result users.GetMerchantsMaterialResponse
	result.Common = &users.CommonResponse{
		Code: 0,
		Msg:  "",
	}
	result.Info = &users.MerchantsMaterialInfo{}
	if req.MaterialId <= 0 {
		result.Common.Code = users.RetCode_MERCHANT_NOT_EXIST
		result.Common.Msg = errcode.GetErrMsg(code.MerchantNotExist)
		return &result, nil
	}
	fmt.Println("开始获取商户认证材料")
	merchantInfo, retCode := service.GetMerchantsMaterial(ctx, req)
	fmt.Printf("商户认证资料 merchantInfo: %+v", merchantInfo)
	if retCode != code.Success {
		if retCode == code.UserExist {
			result.Common.Code = users.RetCode_USER_EXIST
			result.Common.Msg = errcode.GetErrMsg(code.UserExist)
		} else if retCode == code.UserNotExist {
			result.Common.Code = users.RetCode_USER_NOT_EXIST
			result.Common.Msg = errcode.GetErrMsg(code.UserNotExist)
		} else if retCode == code.MerchantNotExist {
			result.Common.Code = users.RetCode_MERCHANT_NOT_EXIST
			result.Common.Msg = errcode.GetErrMsg(code.MerchantNotExist)
		} else if retCode == code.MerchantExist {
			result.Common.Code = users.RetCode_MERCHANT_EXIST
			result.Common.Msg = errcode.GetErrMsg(code.MerchantExist)
		} else {
			result.Common.Code = users.RetCode_ERROR
			result.Common.Msg = errcode.GetErrMsg(code.ErrorServer)
		}
	} else {
		result.Common.Code = users.RetCode_SUCCESS
		result.Common.Msg = errcode.GetErrMsg(code.Success)
	}
	result.Info = &users.MerchantsMaterialInfo{
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
	return &result, nil
}

func (m *MerchantsServer) MerchantsAssociateShop(ctx context.Context, req *users.MerchantsAssociateShopRequest) (*users.MerchantsAssociateShopResponse, error) {
	return &users.MerchantsAssociateShopResponse{}, nil
}
