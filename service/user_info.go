package service

import (
	"context"
	"gitee.com/cristiane/micro-mall-users/pkg/code"
	"gitee.com/cristiane/micro-mall-users/proto/micro_mall_users_proto/users"
	"gitee.com/cristiane/micro-mall-users/repository"
	"gitee.com/kelvins-io/kelvins"
)

const (
	sqlSelectFindUserInfoMain = "id,user_name,country_code,phone,age,contact_addr"
)

func FindUserInfo(ctx context.Context, req *users.FindUserInfoRequest) (result []*users.UserInfoMain, retCode int) {
	result = make([]*users.UserInfoMain, 0)
	retCode = code.Success
	userInfoList, err := repository.FindUserInfo(sqlSelectFindUserInfoMain, req.GetUidList())
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "FindUserInfo err: %v, uidList: %+v", err, req.GetUidList())
		return
	}
	if len(userInfoList) == 0 {
		return
	}
	result = make([]*users.UserInfoMain, len(userInfoList))
	for i := 0; i < len(userInfoList); i++ {
		userInfoMain := &users.UserInfoMain{
			Uid:         int64(userInfoList[i].Id),
			Name:        userInfoList[i].UserName,
			CountryCode: userInfoList[i].CountryCode,
			Phone:       userInfoList[i].Phone,
			Age:         int32(userInfoList[i].Age),
			Address:     userInfoList[i].ContactAddr,
		}
		result[i] = userInfoMain
	}

	return
}
