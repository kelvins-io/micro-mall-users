package service

import (
	"context"
	"database/sql"
	"fmt"
	"gitee.com/cristiane/micro-mall-users/model/args"
	"gitee.com/cristiane/micro-mall-users/model/mysql"
	"gitee.com/cristiane/micro-mall-users/pkg/code"
	"gitee.com/cristiane/micro-mall-users/pkg/util"
	"gitee.com/cristiane/micro-mall-users/pkg/util/cache"
	"gitee.com/cristiane/micro-mall-users/proto/micro_mall_users_proto/users"
	"gitee.com/cristiane/micro-mall-users/repository"
	"gitee.com/cristiane/micro-mall-users/vars"
	"gitee.com/kelvins-io/common/json"
	"gitee.com/kelvins-io/common/password"
	"gitee.com/kelvins-io/kelvins"
	"strings"
	"time"
)

func RegisterUser(ctx context.Context, req *users.RegisterRequest) (args.RegisterResult, int) {
	result := args.RegisterResult{}
	retCode := code.Success
	salt := password.GenerateSalt()
	pwd := password.GeneratePassword(req.Password, salt)
	var user = mysql.User{
		AccountId:    GenAccountId(),
		UserName:     req.UserName,
		Password:     pwd,
		PasswordSalt: salt,
		Sex:          int(req.Sex),
		Phone:        req.Phone,
		CountryCode:  req.CountryCode,
		Email:        req.Email,
		State:        0,
		IdCardNo: sql.NullString{
			String: req.IdCardNo,
		},
		Inviter:    int(req.InviterUser),
		InviteCode: GenInviterCode(),
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
	err := repository.CreateUser(&user)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "CreateUser err: %v, user: %+v", err, user)
		if strings.Contains(err.Error(), code.GetMsg(code.DBDuplicateEntry)) {
			return result, code.UserExist
		}
		return result, code.ErrorServer
	}
	result.InviteCode = user.InviteCode
	// 协程触发邮件
	go func() {
		pushNoticeService := NewPushNoticeService(vars.QueueServerUserRegisterNotice, PushMsgTag{
			DeliveryTag:    args.TaskNameUserRegisterNotice,
			DeliveryErrTag: args.TaskNameUserRegisterNoticeErr,
			RetryCount:     vars.QueueAMQPSettingUserRegisterNotice.TaskRetryCount,
			RetryTimeout:   vars.QueueAMQPSettingUserRegisterNotice.TaskRetryTimeout,
		})
		businessMsg := args.CommonBusinessMsg{
			Type: args.UserStateEventTypeRegister,
			Tag:  args.GetMsg(args.UserStateEventTypeRegister),
			UUID: genUUID(),
			Msg: json.MarshalToStringNoError(args.UserRegisterNotice{
				CountryCode: req.CountryCode,
				Phone:       req.Phone,
				Time:        util.ParseTimeOfStr(time.Now().Unix()),
				State:       0,
			}),
		}
		taskUUID, retCode := pushNoticeService.PushMessage(ctx, businessMsg)
		if retCode != code.Success {
			kelvins.ErrLogger.Errorf(ctx, "businessMsg: %+v register notice send err: ", businessMsg, code.GetMsg(retCode))
		}
		kelvins.BusinessLogger.Infof(ctx, "businessMsg: %+v register notice taskUUID :%v", businessMsg, taskUUID)
	}()

	return result, retCode
}

func LoginUser(ctx context.Context, req *users.LoginUserRequest) (string, int) {
	result := ""
	retCode := code.Success
	user := &mysql.User{}
	switch req.GetLoginType() {
	case users.LoginType_VERIFY_CODE:
		loginInfo := req.GetVerifyCode()
		userDB, err := repository.GetUserByPhone(loginInfo.GetPhone().GetCountryCode(), loginInfo.GetPhone().GetPhone())
		if err != nil {
			kelvins.ErrLogger.Errorf(ctx, "GetUserByPhone err: %v, req: %+v", err, req)
			return result, code.ErrorServer
		}
		user = userDB
	case users.LoginType_PWD:
		loginInfo := req.GetPwd()
		switch loginInfo.GetLoginKind() {
		case users.LoginPwdKind_MOBILE_PHONE:
			mobile := loginInfo.GetPhone()
			userDB, err := repository.GetUserByPhone(mobile.GetCountryCode(), mobile.GetPhone())
			if err != nil {
				kelvins.ErrLogger.Errorf(ctx, "GetUserByPhone err: %v, req: %+v", err, req)
				return result, code.ErrorServer
			}
			if userDB.Id <= 0 {
				return "", code.UserNotExist
			}
			pwd := password.GeneratePassword(loginInfo.GetPwd(), userDB.PasswordSalt)
			if pwd != userDB.Password {
				return result, code.UserPwdNotMatch
			}
			user = userDB
		case users.LoginPwdKind_EMAIL:
			userDB, err := repository.GetUserByEmail(loginInfo.GetEmail().GetContent())
			if err != nil {
				kelvins.ErrLogger.Errorf(ctx, "GetUserByPhone err: %v, req: %+v", err, req)
				return result, code.ErrorServer
			}
			user = userDB
		}
	case users.LoginType_TOKEN:
	}
	if user.Id <= 0 {
		return "", code.UserNotExist
	}
	token, err := util.GenerateToken(user.UserName, user.Id)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "GenerateToken err: %v, req: %+v", err, user)
		return token, code.ErrorServer
	}
	result = token
	// 更新用户状态
	go func() {
		state := args.UserOnlineState{
			Uid:   user.Id,
			State: "online",
			Time:  util.ParseTimeOfStr(time.Now().Unix()),
		}
		userLoginKey := fmt.Sprintf("%v%d", args.CacheKeyUserSate, user.Id)
		err := cache.Set(kelvins.RedisConn, userLoginKey, json.MarshalToStringNoError(state), 7200)
		if err != nil {
			kelvins.ErrLogger.Errorf(ctx, "setUserState err: %v, userLoginKey: %+v", err, userLoginKey)
		}
	}()

	return result, retCode
}

func CheckUserIdentity(ctx context.Context, req *users.CheckUserIdentityRequest) int {
	userDB, err := repository.GetUserByPhone(req.GetCountryCode(), req.GetPhone())
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "GetUserByPhone err: %v, req: %+v", err, req)
		return code.ErrorServer
	}
	if userDB.Id <= 0 {
		return code.UserNotExist
	}
	pwd := password.GeneratePassword(req.GetPwd(), userDB.PasswordSalt)
	if pwd != userDB.Password {
		return code.UserPwdNotMatch
	}
	return code.Success
}

func PasswordReset(ctx context.Context, req *users.PasswordResetRequest) int {
	user, err := repository.GetUserByUid(int(req.GetUid()))
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "GetUserByUid err: %v, uid: %v", err, req.GetUid())
		return code.ErrorServer
	}
	if user.Id <= 0 {
		return code.UserNotExist
	}
	pwd := password.GeneratePassword(req.GetPwd(), user.PasswordSalt)
	where := map[string]interface{}{
		"id": req.GetUid(),
	}
	maps := map[string]interface{}{
		"password": pwd,
	}
	err = repository.UpdateUserInfo(where, maps)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "UpdateUserInfo err: %v, where: %+v, maps: %+v", err, where, maps)
		return code.ErrorServer
	}

	// 触发密码变更消息
	go func() {
		pushNoticeService := NewPushNoticeService(vars.QueueServerUserStateNotice, PushMsgTag{
			DeliveryTag:    args.TaskNameUserStateNotice,
			DeliveryErrTag: args.TaskNameUserStateNoticeErr,
			RetryCount:     vars.QueueAMQPSettingUserStateNotice.TaskRetryCount,
			RetryTimeout:   vars.QueueAMQPSettingUserStateNotice.TaskRetryTimeout,
		})
		businessMsg := args.CommonBusinessMsg{
			Type: args.UserStateEventTypePwdModify,
			Tag:  args.GetMsg(args.UserStateEventTypePwdModify),
			UUID: genUUID(),
			Msg: json.MarshalToStringNoError(args.UserStateNotice{
				Uid:  user.Id,
				Time: util.ParseTimeOfStr(time.Now().Unix()),
			}),
		}
		taskUUID, retCode := pushNoticeService.PushMessage(ctx, businessMsg)
		if retCode != code.Success {
			kelvins.ErrLogger.Errorf(ctx, "Password Reset businessMsg: %+v  notice send err: ", businessMsg, code.GetMsg(retCode))
		}
		kelvins.ErrLogger.Infof(ctx, "Password Reset businessMsg: %+v  taskUUID :%v", businessMsg, taskUUID)
	}()

	return code.Success
}

func UpdateUserLoginState(ctx context.Context, req *users.UpdateUserLoginStateRequest) int {
	user, err := repository.GetUserByUid(int(req.GetUid()))
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "GetUserByUid err: %v, uid: %v", err, req.GetUid())
		return code.ErrorServer
	}
	if user.Id <= 0 {
		return code.UserNotExist
	}
	state := args.UserOnlineState{
		Uid:   int(req.Uid),
		State: "online",
		Time:  util.ParseTimeOfStr(time.Now().Unix()),
	}
	userLoginKey := fmt.Sprintf("%v%d", args.CacheKeyUserSate, req.Uid)
	err = cache.Set(kelvins.RedisConn, userLoginKey, json.MarshalToStringNoError(state), 7200)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "setUserState err: %v, userLoginKey: %+v", err, userLoginKey)
		return code.ErrorServer
	}
	return code.Success
}

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

func CheckUserExist(ctx context.Context, countryCode, phone string) (bool, int) {
	exist, err := repository.CheckUserExistByPhone(countryCode, phone)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "CheckUserExistByPhone err: %v, countryCode: %v, phone:%v", err, countryCode, phone)
		return exist, code.ErrorServer
	}
	return exist, code.Success
}

func GetUserInfoByInviteCode(ctx context.Context, inviteCode string) (*mysql.User, int) {
	user, err := repository.GetUserByInviteCode(inviteCode)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "GetUserByInviteCode err: %v, inviteCode: %v", err, inviteCode)
		return user, code.ErrorServer
	}
	return user, code.Success
}
