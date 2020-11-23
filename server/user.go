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
	if req.Uid <= 0 {
		return &users.GetUserInfoResponse{
			Common: &users.CommonResponse{
				Code: users.RetCode_USER_NOT_EXIST,
				Msg:  errcode.GetErrMsg(code.UserNotExist),
			},
		}, nil
	}
	userInfo, retCode := service.GetUserInfo(ctx, int(req.Uid))
	if retCode != code.Success {
		return &users.GetUserInfoResponse{
			Common: &users.CommonResponse{
				Code: users.RetCode_ERROR,
				Msg:  errcode.GetErrMsg(code.ErrorServer),
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
	if req.Phone == "" || req.CountryCode == "" {
		return &users.GetUserInfoByPhoneResponse{
			Common: &users.CommonResponse{
				Code: users.RetCode_USER_NOT_EXIST,
				Msg:  errcode.GetErrMsg(code.UserNotExist),
			},
		}, nil
	}
	userInfo, retCode := service.GetUserInfoByPhone(ctx, req.CountryCode, req.Phone)
	if retCode != code.Success {
		return &users.GetUserInfoByPhoneResponse{
			Common: &users.CommonResponse{
				Code: users.RetCode_ERROR,
				Msg:  errcode.GetErrMsg(code.ErrorServer),
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
			Msg:  "",
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
		result.Common.Msg = errcode.GetErrMsg(retCode)
		return result, nil
	}
	result.IsExist = exist
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
		case code.UserExist:
			result.Common.Code = users.RetCode_USER_EXIST
		default:
			result.Common.Code = users.RetCode_ERROR
		}
		result.Common.Msg = errcode.GetErrMsg(retCode)
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
		case code.UserNotExist:
			result.Common.Code = users.RetCode_USER_NOT_EXIST
		case code.UserPwdNotMatch:
			result.Common.Code = users.RetCode_USER_PWD_NOT_MATCH
		default:
			result.Common.Code = users.RetCode_ERROR
		}
		result.Common.Msg = errcode.GetErrMsg(retCode)
		return result, nil
	}
	result.IdentityToken = token
	return result, nil
}

func (u *UsersServer) PasswordReset(ctx context.Context, req *users.PasswordResetRequest) (*users.PasswordResetResponse, error) {
	result := &users.PasswordResetResponse{Common: &users.CommonResponse{
		Code: users.RetCode_SUCCESS,
		Msg:  "",
	}}
	retCode := service.PasswordReset(ctx, req)
	if retCode != code.Success {
		switch retCode {
		case code.UserNotExist:
			result.Common.Code = users.RetCode_USER_NOT_EXIST
		default:
			result.Common.Code = users.RetCode_ERROR
		}
		result.Common.Msg = errcode.GetErrMsg(retCode)
		return result, nil
	}
	return result, nil
}

func (u *UsersServer) CheckUserIdentity(ctx context.Context, req *users.CheckUserIdentityRequest) (*users.CheckUserIdentityResponse, error) {
	result := &users.CheckUserIdentityResponse{Common: &users.CommonResponse{
		Code: users.RetCode_SUCCESS,
		Msg:  "",
	}}
	retCode := service.CheckUserIdentity(ctx, req)
	if retCode != code.Success {
		switch retCode {
		case code.UserNotExist:
			result.Common.Code = users.RetCode_USER_NOT_EXIST
		case code.UserPwdNotMatch:
			result.Common.Code = users.RetCode_USER_PWD_NOT_MATCH
		default:
			result.Common.Code = users.RetCode_ERROR
		}
		result.Common.Msg = errcode.GetErrMsg(retCode)
		return result, nil
	}
	return result, nil
}

func (u *UsersServer) UpdateUserLoginState(ctx context.Context, req *users.UpdateUserLoginStateRequest) (*users.UpdateUserLoginStateResponse, error) {
	result := &users.UpdateUserLoginStateResponse{Common: &users.CommonResponse{
		Code: users.RetCode_SUCCESS,
		Msg:  "",
	}}
	retCode := service.UpdateUserLoginState(ctx, req)
	if retCode != code.Success {
		switch retCode {
		case code.UserNotExist:
			result.Common.Code = users.RetCode_USER_NOT_EXIST
		default:
			result.Common.Code = users.RetCode_ERROR
		}
		result.Common.Msg = errcode.GetErrMsg(retCode)
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
				Msg:  errcode.GetErrMsg(code.ErrorServer),
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
	}, Info: make([]*users.UserDeliveryInfo, 0)}
	list, retCode := service.GetUserDeliveryInfo(ctx, req)
	if retCode != code.Success {
		result.Common.Code = users.RetCode_ERROR
		return result, nil
	}
	result.Info = list
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
