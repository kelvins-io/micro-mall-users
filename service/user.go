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
	"gitee.com/cristiane/micro-mall-users/proto/micro_mall_pay_proto/pay_business"
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
	idCardNo := sql.NullString{
		String: req.IdCardNo,
		Valid:  false,
	}
	if req.IdCardNo != "" {
		idCardNo.Valid = true
	}
	user := mysql.User{
		AccountId:    GenAccountId(),
		UserName:     req.UserName,
		Password:     pwd,
		PasswordSalt: salt,
		Sex:          int(req.Sex),
		Age:          int(req.Age),
		ContactAddr:  req.ContactAddr,
		Phone:        req.Phone,
		CountryCode:  req.CountryCode,
		Email:        req.Email,
		State:        0,
		IdCardNo:     idCardNo,
		Inviter:      int(req.InviterUser),
		InviteCode:   GenInviterCode(),
		CreateTime:   time.Now(),
		UpdateTime:   time.Now(),
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

const sqlSelectLoginUser = "id,user_name,password,password_salt"

func LoginUser(ctx context.Context, req *users.LoginUserRequest) (string, int) {
	result := ""
	retCode := code.Success
	user := &mysql.User{}
	switch req.GetLoginType() {
	case users.LoginType_VERIFY_CODE:
		loginInfo := req.GetVerifyCode()
		userDB, err := repository.GetUserByPhone(sqlSelectLoginUser, loginInfo.GetPhone().GetCountryCode(), loginInfo.GetPhone().GetPhone())
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
			userDB, err := repository.GetUserByPhone(sqlSelectLoginUser, mobile.GetCountryCode(), mobile.GetPhone())
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
			userDB, err := repository.GetUserByEmail(sqlSelectLoginUser, loginInfo.GetEmail().GetContent())
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
	userDB, err := repository.GetUserByPhone("password,password_salt", req.GetCountryCode(), req.GetPhone())
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
	user, err := repository.GetUserByPhone("*", countryCode, phone)
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

func ModifyUserDeliveryInfo(ctx context.Context, req *users.ModifyUserDeliveryInfoRequest) int {
	if req.OperationType == users.OperationType_CREATE {
		deliveryInfo := &mysql.UserLogisticsDelivery{
			Uid:          req.Uid,
			DeliveryUser: req.Info.DeliveryUser,
			CountryCode:  "86",
			Phone:        req.Info.MobilePhone,
			Area:         req.Info.Area,
			AreaDetailed: req.Info.DetailedArea,
			Label:        strings.Join(req.Info.Label, "|"),
			IsDefault:    int(req.Info.IsDefault),
			CreateTime:   time.Now(),
			UpdateTime:   time.Now(),
		}
		if req.Info.IsDefault == users.IsDefaultType_DEFAULT_TYPE_TRUE {
			tx := kelvins.XORM_DBEngine.NewSession()
			err := tx.Begin()
			if err != nil {
				kelvins.ErrLogger.Errorf(ctx, "ModifyUserDeliveryInfo create Begin err: %v", err)
				return code.ErrorServer
			}
			where := map[string]interface{}{
				"uid":        req.Uid,
				"is_default": 1,
			}
			maps := map[string]interface{}{
				"is_default":  0,
				"update_time": time.Now(),
			}
			rowAffected, err := repository.UpdateUserLogisticsDeliveryByTx(tx, where, maps)
			if err != nil {
				errCallback := tx.Rollback()
				if errCallback != nil {
					kelvins.ErrLogger.Errorf(ctx, "UpdateUserLogisticsDeliveryByTx Rollback err:%v", errCallback)
				}
				kelvins.ErrLogger.Errorf(ctx, "UpdateUserLogisticsDeliveryByTx err: %v, where: %v", err, where)
				return code.ErrorServer
			}
			if rowAffected <= 0 {
				// 用户第一次添加除外
				//errCallback := tx.Rollback()
				//if errCallback != nil {
				//	kelvins.ErrLogger.Errorf(ctx,"UpdateUserLogisticsDeliveryByTx rowAffected Rollback err:%v",errCallback)
				//}
				//return code.TransactionFailed
			}
			err = repository.CreateUserLogisticsDeliveryByTx(tx, deliveryInfo)
			if err != nil {
				errCallback := tx.Rollback()
				if errCallback != nil {
					kelvins.ErrLogger.Errorf(ctx, "CreateUserLogisticsDeliveryByTx Rollback err:%v", errCallback)
				}
				kelvins.ErrLogger.Errorf(ctx, "CreateUserLogisticsDeliveryByTx err: %v, deliveryInfo: %v", err, deliveryInfo)
				return code.ErrorServer
			}
			err = tx.Commit()
			if err != nil {
				kelvins.ErrLogger.Errorf(ctx, "ModifyUserDeliveryInfo create Commit err: %v", err)
				return code.ErrorServer
			}
			return code.Success
		}
		err := repository.CreateUserLogisticsDelivery(deliveryInfo)
		if err != nil {
			kelvins.ErrLogger.Errorf(ctx, "CreateUserLogisticsDelivery err: %v, deliveryInfo: %v", err, deliveryInfo)
			return code.ErrorServer
		}
		return code.Success
	} else if req.OperationType == users.OperationType_UPDATE {
		if req.Info.Id <= 0 {
			return code.UserDeliveryInfoNotExist
		}
		deliveryInfoDB, err := repository.GetUserLogisticsDelivery("id", req.Info.Id)
		if err != nil {
			kelvins.ErrLogger.Errorf(ctx, "GetUserLogisticsDelivery err: %v, id: %v", err, req.Info.Id)
			return code.ErrorServer
		}
		if deliveryInfoDB.Id <= 0 {
			return code.UserDeliveryInfoNotExist
		}
		deliveryInfo := &mysql.UserLogisticsDelivery{
			DeliveryUser: req.Info.DeliveryUser,
			Phone:        req.Info.MobilePhone,
			Area:         req.Info.Area,
			AreaDetailed: req.Info.DetailedArea,
			Label:        strings.Join(req.Info.Label, "|"),
			IsDefault:    int(req.Info.IsDefault),
			UpdateTime:   time.Now(),
		}
		if req.Info.IsDefault == users.IsDefaultType_DEFAULT_TYPE_TRUE {
			tx := kelvins.XORM_DBEngine.NewSession()
			err := tx.Begin()
			if err != nil {
				kelvins.ErrLogger.Errorf(ctx, "ModifyUserDeliveryInfo create Begin err: %v", err)
				return code.ErrorServer
			}
			where := map[string]interface{}{
				"uid":        req.Uid,
				"is_default": 1,
			}
			maps := map[string]interface{}{
				"is_default":  0,
				"update_time": time.Now(),
			}
			rowAffected, err := repository.UpdateUserLogisticsDeliveryByTx(tx, where, maps)
			if err != nil {
				errCallback := tx.Rollback()
				if errCallback != nil {
					kelvins.ErrLogger.Errorf(ctx, "UpdateUserLogisticsDeliveryByTx Rollback err:%v", errCallback)
				}
				kelvins.ErrLogger.Errorf(ctx, "UpdateUserLogisticsDeliveryByTx err: %v, where: %v", err, where)
				return code.ErrorServer
			}
			if rowAffected <= 0 {
				// 用户第一次添加除外
				//errCallback := tx.Rollback()
				//if errCallback != nil {
				//	kelvins.ErrLogger.Errorf(ctx,"UpdateUserLogisticsDeliveryByTx rowAffected Rollback err:%v",errCallback)
				//}
				//return code.TransactionFailed
			}
			where2 := map[string]interface{}{
				"id": req.Info.Id,
			}
			rowsAffected, err := repository.UpdateUserLogisticsDeliveryByTx(tx, where2, deliveryInfo)
			if err != nil {
				kelvins.ErrLogger.Errorf(ctx, "UpdateUserLogisticsDeliveryByTx err: %v,id: %v, deliveryInfo: %v", err, req.Info.Id, deliveryInfo)
				return code.ErrorServer
			}
			if rowsAffected != 1 {
				errCallback := tx.Rollback()
				if errCallback != nil {
					kelvins.ErrLogger.Errorf(ctx, "UpdateUserLogisticsDeliveryByTx rowAffected Rollback err:%v", errCallback)
				}
				return code.TransactionFailed
			}
			err = tx.Commit()
			if err != nil {
				kelvins.ErrLogger.Errorf(ctx, "ModifyUserDeliveryInfo create Commit err: %v", err)
				return code.ErrorServer
			}
			return code.Success
		}
		where := map[string]interface{}{
			"id": req.Info.Id,
		}
		rowsAffected, err := repository.UpdateUserLogisticsDelivery(where, deliveryInfo)
		if err != nil {
			kelvins.ErrLogger.Errorf(ctx, "UpdateUserLogisticsDelivery err: %v,id: %v, deliveryInfo: %v", err, req.Info.Id, deliveryInfo)
			return code.ErrorServer
		}
		if rowsAffected != 1 {
			return code.TransactionFailed
		}
		return code.Success
	}

	return code.Success
}

const sqlSelectUserDeliveryInfo = "id,delivery_user,country_code,phone,area,area_detailed,is_default,label"

func GetUserDeliveryInfo(ctx context.Context, req *users.GetUserDeliveryInfoRequest) ([]*users.UserDeliveryInfo, int) {
	result := make([]*users.UserDeliveryInfo, 0)
	if req.Uid <= 0 {
		return result, code.UserNotExist
	}
	if req.UserDeliveryId <= 0 {
		list, err := repository.GetUserLogisticsDeliveryList(sqlSelectUserDeliveryInfo, req.Uid)
		if err != nil {
			kelvins.ErrLogger.Errorf(ctx, "GetUserLogisticsDeliveryList err: %v, uid: %v", err, req.Uid)
			return result, code.ErrorServer
		}
		result = make([]*users.UserDeliveryInfo, len(list))
		for i := 0; i < len(list); i++ {
			info := &users.UserDeliveryInfo{
				Id:           list[i].Id,
				DeliveryUser: list[i].DeliveryUser,
				MobilePhone:  fmt.Sprintf("%s-%s", list[i].CountryCode, list[i].Phone),
				Area:         list[i].Area,
				DetailedArea: list[i].AreaDetailed,
				Label:        strings.Split(list[i].Label, "|"),
				IsDefault:    users.IsDefaultType(list[i].IsDefault),
			}
			result[i] = info
		}
	} else {
		infoDB, err := repository.GetUserLogisticsDelivery(sqlSelectUserDeliveryInfo, int64(req.UserDeliveryId))
		if err != nil {
			kelvins.ErrLogger.Errorf(ctx, "GetUserLogisticsDelivery err: %v, uid: %v,id: %v", err, req.Uid, req.UserDeliveryId)
			return result, code.ErrorServer
		}
		info := &users.UserDeliveryInfo{
			Id:           infoDB.Id,
			DeliveryUser: infoDB.DeliveryUser,
			MobilePhone:  fmt.Sprintf("%s-%s", infoDB.CountryCode, infoDB.Phone),
			Area:         infoDB.Area,
			DetailedArea: infoDB.AreaDetailed,
			Label:        strings.Split(infoDB.Label, "|"),
			IsDefault:    users.IsDefaultType(infoDB.IsDefault),
		}
		result = append(result, info)
	}
	return result, code.Success
}

const (
	sqlSelectFindUserInfoMain = "id,user_name,country_code,phone,age,contact_addr"
)

func FindUserInfo(ctx context.Context, req *users.FindUserInfoRequest) (result []*users.UserInfoMain, retCode int) {
	result = make([]*users.UserInfoMain, 0)
	retCode = code.Success
	userInfoList, err := repository.FindUserInfo(sqlSelectFindUserInfoMain, req.GetUidList())
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "FindUserInfo err: %v, uidList: %+v", err, req.GetUidList())
		retCode = code.ErrorServer
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

func UserAccountCharge(ctx context.Context, req *users.UserAccountChargeRequest) (retCode int) {
	retCode = code.Success
	userInfoList, err := repository.FindUserInfo("id,account_id", req.UidList)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "FindUserInfo err: %v, uidList: %+v", err, req.GetUidList())
		retCode = code.ErrorServer
		return
	}
	if len(userInfoList) == 0 {
		retCode = code.UserNotExist
		return
	}
	if len(userInfoList) != len(req.GetUidList()) {
		retCode = code.UserNotExist
		return
	}
	accountIdList := make([]string, 0)
	for i := 0; i < len(userInfoList); i++ {
		accountIdList = append(accountIdList, userInfoList[i].AccountId)
	}
	serverName := args.RpcServiceMicroMallPay
	conn, err := util.GetGrpcClient(serverName)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "GetGrpcClient %v,err: %v", serverName, err)
		return code.ErrorServer
	}
	defer conn.Close()
	payClient := pay_business.NewPayBusinessServiceClient(conn)
	payReq := &pay_business.AccountChargeRequest{
		Owner:       accountIdList,
		AccountType: pay_business.AccountType(req.AccountType),
		CoinType:    pay_business.CoinType(req.CoinType),
		OutTradeNo: req.OutTradeNo,
		Amount:      req.Amount,
		OpMeta: &pay_business.OperationMeta{
			OpUid:      req.OpMeta.OpUid,
			OpIp:       req.OpMeta.OpIp,
			OpPlatform: req.OpMeta.OpPlatform,
			OpDevice:   req.OpMeta.OpDevice,
		},
	}
	payRsp, err := payClient.AccountCharge(ctx, payReq)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "AccountCharge %v,err: %v", serverName, err)
		return code.ErrorServer
	}
	if payRsp.Common.Code != pay_business.RetCode_SUCCESS {
		kelvins.ErrLogger.Errorf(ctx, "AccountCharge  %v,err: %v, req: %+v, rsp: %+v", serverName, err, payReq, payRsp)
		switch payRsp.Common.Code {
		case pay_business.RetCode_USER_ACCOUNT_NOT_EXIST:
			retCode = code.AccountNotExist
		case pay_business.RetCode_TRANSACTION_FAILED:
			retCode = code.TransactionFailed
		case pay_business.RetCode_USER_ACCOUNT_STATE_INVALID:
			retCode = code.AccountStateInvalid
		case pay_business.RetCode_USER_ACCOUNT_STATE_LOCK:
			retCode = code.AccountStateLock
		case pay_business.RetCode_TRADE_PAY_RUN:
			retCode = code.UserChargeRun
		case pay_business.RetCode_TRADE_PAY_SUCCESS:
			retCode = code.UserChargeSuccess
		case pay_business.RetCode_TRADE_UUID_EMPTY:
			retCode = code.UserChargeTradeNoEmpty
		default:
			retCode = code.ErrorServer
		}
		return
	}
	return
}

func CheckUserDeliveryInfo(ctx context.Context, req *users.CheckUserDeliveryInfoRequest) (retCode int) {
	retCode = code.Success
	infoList, err := repository.CheckUserLogisticsDelivery(req.Uid, req.DeliveryIds)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "AccountCharge err: %v, req: %+v", err, req)
		retCode = code.ErrorServer
		return
	}
	if len(infoList) != len(req.DeliveryIds) {
		retCode = code.UserDeliveryInfoNotExist
		return
	}
	return
}

func CheckUserState(ctx context.Context, req *users.CheckUserStateRequest) (retCode int) {
	retCode = code.Success
	infoList, err := repository.FindUserInfo("id,state", req.GetUidList())
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "AccountCharge err: %v, req: %+v", err, req)
		retCode = code.ErrorServer
		return
	}
	if len(infoList) == 0 {
		retCode = code.UserNotExist
		return
	}
	if len(infoList) != len(req.GetUidList()) {
		retCode = code.UserNotExist
		return
	}
	for i := 0; i < len(infoList); i++ {
		if infoList[i].Id <= 0 {
			retCode = code.UserNotExist
			return
		}
		if infoList[i].State != 3 {
			retCode = code.UserNotExist
			return
		}
	}
	return
}

func GetUserAccountId(ctx context.Context, req *users.GetUserAccountIdRequest) (result []*users.UserAccountInfo, retCode int) {
	retCode = code.Success
	userList, err := repository.FindUserInfo("id,account_id", req.GetUidList())
	result = make([]*users.UserAccountInfo, len(userList))
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "FindUserInfo err: %v, req: %+v", err, req)
		retCode = code.ErrorServer
		return
	}
	if len(userList) == 0 {
		retCode = code.UserNotExist
		return
	}
	if len(userList) != len(req.GetUidList()) {
		retCode = code.UserNotExist
		return
	}
	for i := 0; i < len(userList); i++ {
		if userList[i].Id < 0 || userList[i].AccountId == "" {
			retCode = code.UserNotExist
			return
		}
		accountInfo := &users.UserAccountInfo{
			Uid:       int64(userList[i].Id),
			AccountId: userList[i].AccountId,
		}
		result[i] = accountInfo
	}
	return
}
