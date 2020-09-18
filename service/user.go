package service

import (
	"context"
	"gitee.com/cristiane/micro-mall-users/model/mysql"
	"gitee.com/cristiane/micro-mall-users/pkg/code"
	"gitee.com/cristiane/micro-mall-users/pkg/util"
	"gitee.com/cristiane/micro-mall-users/proto/micro_mall_users_proto/users"
	"gitee.com/cristiane/micro-mall-users/repository"
	"gitee.com/kelvins-io/common/errcode"
	"gitee.com/kelvins-io/kelvins"
	"strings"
	"time"
)

func GetUserInfo(ctx context.Context, uid int) (*mysql.User, int) {
	user, err := repository.GetUserByUid(uid)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "GetUserByUid err: %v, uid: %v", err, uid)
		return user, code.ErrorServer
	}
	return user, code.Success
}

func GetUserInfoByPhone(ctx context.Context, countryCode, phone string) (*mysql.User, int) {
	user, err := repository.GetUserByPhone(countryCode, phone)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "GetUserByPhone err: %v, countryCode: %v, phone:%v", err, countryCode, phone)
		return user, code.ErrorServer
	}
	return user, code.Success
}

func CreateUserAccount(ctx context.Context, req *users.CreateUserAccountRequest) (string, int) {
	userInfo, err := repository.GetUserAccountIdByUid(req.Uid)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "GetUserAccountIdByUid err: %v, uid: %+v", err, req.Uid)
		return "", code.ErrorServer
	}
	if userInfo.AccountId == "" {
		return "", code.UserNotExist
	}
	accountCode := util.GetUUID()
	account := mysql.Account{
		AccountCode: accountCode,
		Owner:       userInfo.AccountId,
		Balance:     req.Balance,
		CoinType:    int(req.CoinType),
		CoinDesc:    req.CoinDesc,
		State:       int(req.AccountState),
		AccountType: int(req.AccountType),
		CreateTime:  time.Now(),
		UpdateTime:  time.Now(),
	}
	err = repository.CreateAccount(&account)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "CreateAccount err: %v, account: %+v", err, account)
		if strings.Contains(err.Error(), errcode.GetErrMsg(code.DBDuplicateEntry)) {
			return "", code.AccountExist
		}
		return "", code.ErrorServer
	}
	return accountCode, code.Success
}
