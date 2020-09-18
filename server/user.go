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

func (u *UsersServer) CreateUserAccount(ctx context.Context, req *users.CreateUserAccountRequest) (*users.CreateUserAccountResponse, error) {
	var result users.CreateUserAccountResponse
	result.Common = &users.CommonResponse{
		Code: 0,
		Msg:  "",
	}
	if req.Uid <= 0 {
		result.Common.Code = users.RetCode_USER_NOT_EXIST
		result.Common.Msg = errcode.GetErrMsg(code.UserNotExist)
		return &result, nil
	}

	accountCode, retCode := service.CreateUserAccount(ctx, req)
	if retCode != code.Success {
		if retCode == code.UserNotExist {
			result.Common.Code = users.RetCode_USER_NOT_EXIST
			result.Common.Msg = errcode.GetErrMsg(code.UserNotExist)
			return &result, nil
		}
		if retCode == code.AccountExist {
			result.Common.Code = users.RetCode_ACCOUNT_EXIST
			result.Common.Msg = errcode.GetErrMsg(code.AccountExist)
			return &result, nil
		}
		result.Common.Code = users.RetCode_ERROR
		result.Common.Msg = errcode.GetErrMsg(code.ErrorServer)
		return &result, nil
	}
	result.Common.Code = users.RetCode_SUCCESS
	result.Common.Msg = errcode.GetErrMsg(code.Success)
	result.AccountCode = accountCode
	return &result, nil
}
