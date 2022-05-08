package server

import (
	"context"

	"gitee.com/cristiane/micro-mall-users/pkg/code"
	"gitee.com/cristiane/micro-mall-users/pkg/util"
	"gitee.com/cristiane/micro-mall-users/proto/micro_mall_users_proto/users"
	"gitee.com/cristiane/micro-mall-users/service"
	"gitee.com/kelvins-io/common/errcode"
)

type UsersServer struct {
}

func NewUsersServer() users.UsersServiceServer {
	return new(UsersServer)
}

func (u *UsersServer) GetUserInfo(ctx context.Context, req *users.GetUserInfoRequest) (*users.GetUserInfoResponse, error) {
	userInfo, retCode := service.GetUserInfo(ctx, int(req.Uid))
	if retCode != code.Success {
		return &users.GetUserInfoResponse{
			Common: &users.CommonResponse{
				Code: users.RetCode_ERROR,
				Msg:  code.GetMsg(retCode),
			},
		}, nil
	}
	return &users.GetUserInfoResponse{
		Common: &users.CommonResponse{
			Code: users.RetCode_SUCCESS,
			Msg:  errcode.GetErrMsg(code.Success),
		},
		Info: &users.UserInfo{
			Uid:         int64(userInfo.Id),
			AccountId:   userInfo.AccountId,
			UserName:    userInfo.UserName,
			Sex:         int32(userInfo.Sex),
			CountryCode: userInfo.CountryCode,
			Phone:       userInfo.Phone,
			Email:       userInfo.Email,
			State:       int32(userInfo.State),
			IdCardNo:    userInfo.IdCardNo.String,
			Inviter:     int64(userInfo.Inviter),
			InviterCode: userInfo.InviteCode,
			ContactAddr: userInfo.ContactAddr,
			Age:         int32(userInfo.Age),
			CreateTime:  util.ParseTimeOfStr(userInfo.CreateTime.Unix()),
			UpdateTime:  util.ParseTimeOfStr(userInfo.UpdateTime.Unix()),
		},
	}, nil
}

func (u *UsersServer) GetUserInfoByPhone(ctx context.Context, req *users.GetUserInfoByPhoneRequest) (*users.GetUserInfoByPhoneResponse, error) {
	userInfo, retCode := service.GetUserInfoByPhone(ctx, req.CountryCode, req.Phone)
	if retCode != code.Success {
		return &users.GetUserInfoByPhoneResponse{
			Common: &users.CommonResponse{
				Code: users.RetCode_ERROR,
				Msg:  code.GetMsg(retCode),
			},
		}, nil
	}
	return &users.GetUserInfoByPhoneResponse{
		Common: &users.CommonResponse{
			Code: users.RetCode_SUCCESS,
			Msg:  errcode.GetErrMsg(code.Success),
		},
		Info: &users.UserInfo{
			Uid:         int64(userInfo.Id),
			AccountId:   userInfo.AccountId,
			UserName:    userInfo.UserName,
			Sex:         int32(userInfo.Sex),
			CountryCode: userInfo.CountryCode,
			Phone:       userInfo.Phone,
			Email:       userInfo.Email,
			State:       int32(userInfo.State),
			IdCardNo:    userInfo.IdCardNo.String,
			Inviter:     int64(userInfo.Inviter),
			InviterCode: userInfo.InviteCode,
			ContactAddr: userInfo.ContactAddr,
			Age:         int32(userInfo.Age),
			CreateTime:  util.ParseTimeOfStr(userInfo.CreateTime.Unix()),
			UpdateTime:  util.ParseTimeOfStr(userInfo.UpdateTime.Unix()),
		},
	}, nil
}

func (u *UsersServer) CheckUserByPhone(ctx context.Context, req *users.CheckUserByPhoneRequest) (*users.CheckUserByPhoneResponse, error) {
	result := &users.CheckUserByPhoneResponse{
		Common: &users.CommonResponse{
			Code: users.RetCode_SUCCESS,
		},
		IsExist: false,
	}
	exist, retCode := service.CheckUserExist(ctx, req.CountryCode, req.Phone)
	if retCode != code.Success {
		switch retCode {
		case code.UserNotExist:
			result.Common.Code = users.RetCode_USER_NOT_EXIST
		default:
			result.Common.Code = users.RetCode_ERROR
		}
		result.Common.Msg = code.GetMsg(retCode)
		return result, nil
	}
	result.IsExist = exist
	return result, nil
}

func (u *UsersServer) GenVerifyCode(ctx context.Context, req *users.GenVerifyCodeRequest) (*users.GenVerifyCodeResponse, error) {
	result := &users.GenVerifyCodeResponse{
		Common: &users.CommonResponse{
			Code: users.RetCode_SUCCESS,
		},
	}
	retCode := service.GenVerifyCode(ctx, req)
	if retCode != code.Success {
		switch retCode {
		case code.ErrorVerifyCodeLimited:
			result.Common.Code = users.RetCode_USER_VERIFY_CODE_LIMITED
		case code.ErrorVerifyCodeInterval:
			result.Common.Code = users.RetCode_USER_VERIFY_CODE_INTERVAL
		default:
			result.Common.Code = users.RetCode_ERROR
		}
	}
	return result, nil
}

func (u *UsersServer) Register(ctx context.Context, req *users.RegisterRequest) (*users.RegisterResponse, error) {
	result := &users.RegisterResponse{
		Common: &users.CommonResponse{
			Code: users.RetCode_SUCCESS,
		},
	}
	reg, retCode := service.RegisterUser(ctx, req)
	if retCode != code.Success {
		switch retCode {
		case code.ErrorVerifyCodeForbidden:
			result.Common.Code = users.RetCode_USER_VERIFY_CODE_FORBIDDEN
		case code.ErrorVerifyCodeExpire:
			result.Common.Code = users.RetCode_USER_VERIFY_CODE_EXPIRE
		case code.ErrorVerifyCodeInvalid:
			result.Common.Code = users.RetCode_USER_VERIFY_CODE_INVALID
		case code.UserExist:
			result.Common.Code = users.RetCode_USER_EXIST
		case code.ErrorInviteCodeInvalid:
			result.Common.Code = users.RetCode_USER_INVITE_CODE_INVALID
		case code.TransactionFailed:
			result.Common.Code = users.RetCode_TRANSACTION_FAILED
		default:
			result.Common.Code = users.RetCode_ERROR
		}
		result.Common.Msg = code.GetMsg(retCode)
		return result, nil
	}
	result.Result = &users.RegisterResult{InviteCode: reg.InviteCode}
	return result, nil
}

func (u *UsersServer) LoginUser(ctx context.Context, req *users.LoginUserRequest) (*users.LoginUserResponse, error) {
	result := &users.LoginUserResponse{
		Common: &users.CommonResponse{
			Code: users.RetCode_SUCCESS,
		},
	}
	token, retCode := service.LoginUser(ctx, req)
	if retCode != code.Success {
		switch retCode {
		case code.ErrorVerifyCodeForbidden:
			result.Common.Code = users.RetCode_USER_VERIFY_CODE_FORBIDDEN
		case code.ErrorVerifyCodeExpire:
			result.Common.Code = users.RetCode_USER_VERIFY_CODE_EXPIRE
		case code.ErrorVerifyCodeInvalid:
			result.Common.Code = users.RetCode_USER_VERIFY_CODE_INVALID
		case code.UserNotExist:
			result.Common.Code = users.RetCode_USER_NOT_EXIST
		case code.UserPwdNotMatch:
			result.Common.Code = users.RetCode_USER_PWD_NOT_MATCH
		case code.UserStateForbiddenLogin:
			result.Common.Code = users.RetCode_USER_STATE_FORBIDDEN_LOGIN
		case code.UserStateNotVerify:
			result.Common.Code = users.RetCode_USER_STATE_NOT_VERIFY
		default:
			result.Common.Code = users.RetCode_ERROR
		}
		result.Common.Msg = code.GetMsg(retCode)
		return result, nil
	}
	result.IdentityToken = token
	return result, nil
}

func (u *UsersServer) PasswordReset(ctx context.Context, req *users.PasswordResetRequest) (*users.PasswordResetResponse, error) {
	result := &users.PasswordResetResponse{Common: &users.CommonResponse{
		Code: users.RetCode_SUCCESS,
	}}
	retCode := service.PasswordReset(ctx, req)
	if retCode != code.Success {
		switch retCode {
		case code.ErrorVerifyCodeForbidden:
			result.Common.Code = users.RetCode_USER_VERIFY_CODE_FORBIDDEN
		case code.ErrorVerifyCodeExpire:
			result.Common.Code = users.RetCode_USER_VERIFY_CODE_EXPIRE
		case code.ErrorVerifyCodeInvalid:
			result.Common.Code = users.RetCode_USER_VERIFY_CODE_INVALID
		case code.UserNotExist:
			result.Common.Code = users.RetCode_USER_NOT_EXIST
		default:
			result.Common.Code = users.RetCode_ERROR
		}
		result.Common.Msg = code.GetMsg(retCode)
		return result, nil
	}
	return result, nil
}

func (u *UsersServer) UpdateUserLoginState(ctx context.Context, req *users.UpdateUserLoginStateRequest) (*users.UpdateUserLoginStateResponse, error) {
	result := &users.UpdateUserLoginStateResponse{Common: &users.CommonResponse{
		Code: users.RetCode_SUCCESS,
	}}
	retCode := service.UpdateUserLoginState(ctx, req)
	if retCode != code.Success {
		switch retCode {
		case code.UserNotExist:
			result.Common.Code = users.RetCode_USER_NOT_EXIST
		default:
			result.Common.Code = users.RetCode_ERROR
		}
		result.Common.Msg = code.GetMsg(retCode)
		return result, nil
	}
	return result, nil
}

func (u *UsersServer) GetUserInfoByInviteCode(ctx context.Context, req *users.GetUserByInviteCodeRequest) (*users.GetUserByInviteCodeResponse, error) {
	userInfo, retCode := service.GetUserInfoByInviteCode(ctx, req.InviteCode)
	if retCode != code.Success {
		return &users.GetUserByInviteCodeResponse{
			Common: &users.CommonResponse{
				Code: users.RetCode_ERROR,
				Msg:  code.GetMsg(retCode),
			},
		}, nil
	}
	return &users.GetUserByInviteCodeResponse{
		Common: &users.CommonResponse{
			Code: users.RetCode_SUCCESS,
			Msg:  errcode.GetErrMsg(code.Success),
		},
		Info: &users.UserInfo{
			Uid:         int64(userInfo.Id),
			AccountId:   userInfo.AccountId,
			UserName:    userInfo.UserName,
			Sex:         int32(userInfo.Sex),
			CountryCode: userInfo.CountryCode,
			Phone:       userInfo.Phone,
			Email:       userInfo.Email,
			State:       int32(userInfo.State),
			IdCardNo:    userInfo.IdCardNo.String,
			Inviter:     int64(userInfo.Inviter),
			InviterCode: userInfo.InviteCode,
			ContactAddr: userInfo.ContactAddr,
			Age:         int32(userInfo.Age),
			CreateTime:  util.ParseTimeOfStr(userInfo.CreateTime.Unix()),
			UpdateTime:  util.ParseTimeOfStr(userInfo.UpdateTime.Unix()),
		},
	}, nil
}

func (u *UsersServer) ModifyUserDeliveryInfo(ctx context.Context, req *users.ModifyUserDeliveryInfoRequest) (*users.ModifyUserDeliveryInfoResponse, error) {
	result := users.ModifyUserDeliveryInfoResponse{Common: &users.CommonResponse{
		Code: users.RetCode_SUCCESS,
		Msg:  "",
	}}
	retCode := service.ModifyUserDeliveryInfo(ctx, req)
	if retCode != code.Success {
		switch retCode {
		case code.TransactionFailed:
			result.Common.Code = users.RetCode_TRANSACTION_FAILED
		case code.UserDeliveryInfoExist:
			result.Common.Code = users.RetCode_USER_DELIVERY_INFO_EXIST
		case code.UserDeliveryInfoNotExist:
			result.Common.Code = users.RetCode_USER_DELIVERY_INFO_NOT_EXIST
		case code.ErrorServer:
			result.Common.Code = users.RetCode_ERROR
		}
		return &result, nil
	}
	return &result, nil
}

func (u *UsersServer) GetUserDeliveryInfo(ctx context.Context, req *users.GetUserDeliveryInfoRequest) (*users.GetUserDeliveryInfoResponse, error) {
	result := &users.GetUserDeliveryInfoResponse{Common: &users.CommonResponse{
		Code: users.RetCode_SUCCESS,
	}, InfoList: make([]*users.UserDeliveryInfo, 0)}
	list, retCode := service.GetUserDeliveryInfo(ctx, req)
	if retCode != code.Success {
		switch retCode {
		case code.UserNotExist:
			result.Common.Code = users.RetCode_USER_NOT_EXIST
		case code.UserDeliveryInfoNotExist:
			result.Common.Code = users.RetCode_USER_DELIVERY_INFO_NOT_EXIST
		default:
			result.Common.Code = users.RetCode_ERROR
		}
		return result, nil
	}
	result.InfoList = list
	return result, nil
}

func (u *UsersServer) FindUserInfo(ctx context.Context, req *users.FindUserInfoRequest) (*users.FindUserInfoResponse, error) {
	result := &users.FindUserInfoResponse{
		Common: &users.CommonResponse{
			Code: users.RetCode_SUCCESS,
		},
		InfoList: nil,
	}
	userInfoList, retCode := service.FindUserInfo(ctx, req)
	if retCode != code.Success {
		result.Common.Code = users.RetCode_ERROR
		return result, nil
	}
	result.InfoList = userInfoList
	return result, nil
}

func (u *UsersServer) UserAccountCharge(ctx context.Context, req *users.UserAccountChargeRequest) (*users.UserAccountChargeResponse, error) {
	result := &users.UserAccountChargeResponse{Common: &users.CommonResponse{
		Code: users.RetCode_SUCCESS,
	}}
	retCode := service.UserAccountCharge(ctx, req)
	if retCode != code.Success {
		switch retCode {
		case code.UserNotExist:
			result.Common.Code = users.RetCode_USER_NOT_EXIST
		case code.AccountNotExist:
			result.Common.Code = users.RetCode_ACCOUNT_NOT_EXIST
		case code.TransactionFailed:
			result.Common.Code = users.RetCode_TRANSACTION_FAILED
		case code.AccountStateLock:
			result.Common.Code = users.RetCode_ACCOUNT_LOCK
		case code.AccountStateInvalid:
			result.Common.Code = users.RetCode_ACCOUNT_INVALID
		case code.UserChargeRun:
			result.Common.Code = users.RetCode_USER_CHARGE_RUN
		case code.UserChargeSuccess:
			result.Common.Code = users.RetCode_USER_CHARGE_SUCCESS
		case code.UserChargeTradeNoEmpty:
			result.Common.Code = users.RetCode_USER_CHARGE_TRADE_NO_EMPTY
		default:
			result.Common.Code = users.RetCode_ERROR
		}
		return result, nil
	}
	return result, nil
}

func (u *UsersServer) CheckUserDeliveryInfo(ctx context.Context, req *users.CheckUserDeliveryInfoRequest) (*users.CheckUserDeliveryInfoResponse, error) {
	result := &users.CheckUserDeliveryInfoResponse{Common: &users.CommonResponse{
		Code: users.RetCode_SUCCESS,
	}}
	retCode := service.CheckUserDeliveryInfo(ctx, req)
	if retCode != code.Success {
		switch retCode {
		case code.UserDeliveryInfoNotExist:
			result.Common.Code = users.RetCode_USER_DELIVERY_INFO_NOT_EXIST
		default:
			result.Common.Code = users.RetCode_ERROR
		}
		return result, nil
	}
	return result, nil
}

func (u *UsersServer) CheckUserState(ctx context.Context, req *users.CheckUserStateRequest) (*users.CheckUserStateResponse, error) {
	result := &users.CheckUserStateResponse{Common: &users.CommonResponse{
		Code: users.RetCode_SUCCESS,
	}}
	retCode := service.CheckUserState(ctx, req.GetUidList())
	if retCode != code.Success {
		switch retCode {
		case code.UserNotExist:
			result.Common.Code = users.RetCode_USER_NOT_EXIST
		case code.UserStateNotVerify:
			result.Common.Code = users.RetCode_USER_STATE_NOT_VERIFY
		case code.UserStateForbiddenLogin:
			result.Common.Code = users.RetCode_USER_STATE_FORBIDDEN_LOGIN
		default:
			result.Common.Code = users.RetCode_ERROR
		}
	}
	return result, nil
}

func (u *UsersServer) GetUserAccountId(ctx context.Context, req *users.GetUserAccountIdRequest) (*users.GetUserAccountIdResponse, error) {
	result := &users.GetUserAccountIdResponse{
		Common: &users.CommonResponse{
			Code: users.RetCode_SUCCESS,
		},
		InfoList: nil,
	}
	accountInfoList, retCode := service.GetUserAccountId(ctx, req)
	if retCode != code.Success {
		switch retCode {
		case code.UserNotExist:
			result.Common.Code = users.RetCode_USER_NOT_EXIST
		default:
			result.Common.Code = users.RetCode_ERROR
		}
		return result, nil
	}
	result.InfoList = accountInfoList
	return result, nil
}

func (u *UsersServer) ListUserInfo(ctx context.Context, req *users.ListUserInfoRequest) (*users.ListUserInfoResponse, error) {
	result := &users.ListUserInfoResponse{
		Common: &users.CommonResponse{
			Code: users.RetCode_SUCCESS,
		},
		UserInfoList: nil,
	}
	userInfoList, retCode := service.ListUserInfo(ctx, req)
	result.UserInfoList = userInfoList
	if retCode != code.Success {
		result.Common.Code = users.RetCode_ERROR
		return result, nil
	}
	return result, nil
}

func (u *UsersServer) SearchUserInfo(ctx context.Context, req *users.SearchUserInfoRequest) (*users.SearchUserInfoResponse, error) {
	result := &users.SearchUserInfoResponse{Common: &users.CommonResponse{
		Code: users.RetCode_SUCCESS,
	}}
	userInfoList, retCode := service.SearchUserInfo(ctx, req.GetQuery())
	if retCode != code.Success {
		result.Common.Code = users.RetCode_ERROR
		return result, nil
	}
	result.List = userInfoList
	return result, nil
}
