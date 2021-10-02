package server

import (
	"context"
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
		Code: users.RetCode_SUCCESS,
	}
	if req.Info.Uid <= 0 {
		result.Common.Code = users.RetCode_USER_NOT_EXIST
		result.Common.Msg = errcode.GetErrMsg(code.UserNotExist)
		return &result, nil
	}
	merchantId, retCode := service.MerchantsMaterial(ctx, req)
	result.MaterialId = int64(merchantId)
	if retCode != code.Success {
		switch retCode {
		case code.UserNotExist:
			result.Common.Code = users.RetCode_USER_NOT_EXIST
			result.Common.Msg = errcode.GetErrMsg(code.UserNotExist)
		case code.MerchantExist:
			result.Common.Code = users.RetCode_MERCHANT_EXIST
			result.Common.Msg = errcode.GetErrMsg(code.MerchantExist)
		case code.MerchantNotExist:
			result.Common.Code = users.RetCode_MERCHANT_NOT_EXIST
			result.Common.Msg = errcode.GetErrMsg(code.MerchantNotExist)
		default:
			result.Common.Code = users.RetCode_ERROR
			result.Common.Msg = errcode.GetErrMsg(code.ErrorServer)
		}
	}

	return &result, nil
}

func (m *MerchantsServer) MerchantsMaterialAudit(ctx context.Context, req *users.MerchantsMaterialAuditRequest) (*users.MerchantsMaterialAuditResponse, error) {
	return &users.MerchantsMaterialAuditResponse{}, nil
}

func (m *MerchantsServer) GetMerchantsMaterial(ctx context.Context, req *users.GetMerchantsMaterialRequest) (*users.GetMerchantsMaterialResponse, error) {
	var result users.GetMerchantsMaterialResponse
	result.Common = &users.CommonResponse{
		Code: users.RetCode_SUCCESS,
	}
	result.Info = &users.MerchantsMaterialInfo{}
	if req.MaterialId <= 0 {
		result.Common.Code = users.RetCode_MERCHANT_NOT_EXIST
		result.Common.Msg = errcode.GetErrMsg(code.MerchantNotExist)
		return &result, nil
	}
	merchantInfo, retCode := service.GetMerchantsMaterial(ctx, req)
	if retCode != code.Success {
		switch retCode {
		case code.UserNotExist:
			result.Common.Code = users.RetCode_USER_NOT_EXIST
			result.Common.Msg = errcode.GetErrMsg(code.UserNotExist)
		case code.MerchantNotExist:
			result.Common.Code = users.RetCode_MERCHANT_NOT_EXIST
			result.Common.Msg = errcode.GetErrMsg(code.MerchantNotExist)
		default:
			result.Common.Code = users.RetCode_ERROR
			result.Common.Msg = errcode.GetErrMsg(code.ErrorServer)
		}
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

func (u *MerchantsServer) SearchMerchantInfo(ctx context.Context, req *users.SearchMerchantInfoRequest) (*users.SearchMerchantInfoResponse, error) {
	result := &users.SearchMerchantInfoResponse{Common: &users.CommonResponse{
		Code: users.RetCode_SUCCESS,
	}}
	merchantInfoList, retCode := service.SearchMerchantInfo(ctx, req.GetQuery())
	if retCode != code.Success {
		result.Common.Code = users.RetCode_ERROR
		return result, nil
	}
	result.List = merchantInfoList
	return result, nil
}
