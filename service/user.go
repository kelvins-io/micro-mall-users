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
	"gitee.com/cristiane/micro-mall-users/pkg/util/email"
	"gitee.com/cristiane/micro-mall-users/proto/micro_mall_pay_proto/pay_business"
	"gitee.com/cristiane/micro-mall-users/proto/micro_mall_search_proto/search_business"
	"gitee.com/cristiane/micro-mall-users/proto/micro_mall_users_proto/users"
	"gitee.com/cristiane/micro-mall-users/repository"
	"gitee.com/cristiane/micro-mall-users/vars"
	"gitee.com/kelvins-io/common/hash"
	"gitee.com/kelvins-io/common/json"
	"gitee.com/kelvins-io/common/password"
	"gitee.com/kelvins-io/kelvins"
	"github.com/google/uuid"
	"strconv"
	"strings"
	"time"
)

func RegisterUser(ctx context.Context, req *users.RegisterRequest) (args.RegisterResult, int) {
	result := args.RegisterResult{}
	isExist, ret := CheckUserExist(ctx, req.CountryCode, req.Phone)
	if ret != code.Success {
		return result, code.ErrorServer
	}
	if isExist {
		return result, code.UserExist
	}

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
		State:        3,
		IdCardNo:     idCardNo,
		Inviter:      int(req.InviterUser),
		InviteCode:   GenInviterCode(),
		CreateTime:   time.Now(),
		UpdateTime:   time.Now(),
	}
	tx := kelvins.XORM_DBEngine.NewSession()
	err := tx.Begin()
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "CreateUser NewSession err: %v", err)
		return result, code.ErrorServer
	}
	err = repository.CreateUser(tx, &user)
	if err != nil {
		tx.Rollback()
		kelvins.ErrLogger.Errorf(ctx, "CreateUser err: %v, user: %v", err, json.MarshalToStringNoError(user))
		if strings.Contains(err.Error(), code.GetMsg(code.DBDuplicateEntry)) {
			return result, code.UserExist
		}
		return result, code.ErrorServer
	}
	result.InviteCode = user.InviteCode
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
		Content: json.MarshalToStringNoError(args.UserRegisterNotice{
			CountryCode: req.CountryCode,
			Phone:       req.Phone,
			Time:        util.ParseTimeOfStr(time.Now().Unix()),
			State:       3,
		}),
	}
	_, ret = pushNoticeService.PushMessage(ctx, businessMsg)
	if ret != code.Success {
		tx.Rollback()
		kelvins.ErrLogger.Errorf(ctx, "PushMessage register req: %v, err: %v", json.MarshalToStringNoError(businessMsg), code.GetMsg(ret))
		return result, code.ErrorServer
	}
	tx.Commit()

	// 通知搜索
	body := args.UserInfoSearch{
		UserName:    req.GetUserName(),
		Phone:       req.GetCountryCode() + ":" + req.GetPhone(),
		Email:       req.GetEmail(),
		IdCardNo:    req.GetIdCardNo(),
		ContactAddr: req.GetContactAddr(),
	}
	userInfoSearchNotice(&body)

	return result, code.Success
}

func userInfoSearchNotice(info *args.UserInfoSearch) {
	kelvins.GPool.SendJob(func() {
		var ctx = context.TODO()
		userInfoMsg := args.CommonBusinessMsg{
			Type:    args.UserInfoSearchNoticeType,
			Tag:     args.GetMsg(args.UserInfoSearchNoticeType),
			UUID:    uuid.New().String(),
			Content: json.MarshalToStringNoError(info),
		}
		vars.QueueServerUserInfoSearchPusher.PushMessage(ctx, &userInfoMsg)
	})
}

const sqlSelectLoginUser = "id,user_name,password,password_salt"

func LoginUser(ctx context.Context, req *users.LoginUserRequest) (string, int) {
	result := ""
	loginType := ""
	retCode := code.Success
	user := &mysql.User{}
	switch req.GetLoginType() {
	case users.LoginType_VERIFY_CODE:
		loginInfo := req.GetVerifyCode()
		userDB, err := repository.GetUserByPhone(sqlSelectLoginUser, loginInfo.GetPhone().GetCountryCode(), loginInfo.GetPhone().GetPhone())
		if err != nil {
			kelvins.ErrLogger.Errorf(ctx, "GetUserByPhone err: %v, req: %v", err, json.MarshalToStringNoError(req))
			return result, code.ErrorServer
		}
		loginType = "验证码"
		user = userDB
	case users.LoginType_PWD:
		loginInfo := req.GetPwd()
		switch loginInfo.GetLoginKind() {
		case users.LoginPwdKind_MOBILE_PHONE:
			mobile := loginInfo.GetPhone()
			userDB, err := repository.GetUserByPhone(sqlSelectLoginUser, mobile.GetCountryCode(), mobile.GetPhone())
			if err != nil {
				kelvins.ErrLogger.Errorf(ctx, "GetUserByPhone err: %v, req: %v", err, json.MarshalToStringNoError(req))
				return result, code.ErrorServer
			}
			user = userDB
			loginType = "手机号-密码"
		case users.LoginPwdKind_EMAIL:
			userDB, err := repository.GetUserByEmail(sqlSelectLoginUser, loginInfo.GetEmail().GetContent())
			if err != nil {
				kelvins.ErrLogger.Errorf(ctx, "GetUserByPhone err: %v, req: %v", err, json.MarshalToStringNoError(req))
				return result, code.ErrorServer
			}
			user = userDB
			loginType = "邮箱-密码"
		}
	case users.LoginType_TOKEN:
		loginType = "认证token"
	}
	// 检查用户状态
	retCode = checkUserState(ctx, user.Id)
	if retCode != code.Success {
		return "", retCode
	}
	switch req.GetLoginType() {
	case users.LoginType_PWD:
		pwd := password.GeneratePassword(req.GetPwd().GetPwd(), user.PasswordSalt)
		if pwd != user.Password {
			_, retCode := userLoginFailure(ctx, user.Id)
			if retCode != code.Success {
				return result, retCode
			}
			return result, code.UserPwdNotMatch
		}
	default:
	}
	token, err := util.GenerateToken(user.UserName, user.Id)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "GenerateToken err: %v, req: %v", err, json.MarshalToStringNoError(user))
		return token, code.ErrorServer
	}
	result = token

	updateUserState := func() {
		// 更新在线状态
		retCode := updateUserOnlineState(ctx, user.Id, args.UserOnlineStateOnline)
		if retCode != code.Success {
			return
		}
		// 发送登录邮件
		emailNotice := fmt.Sprintf(args.UserLoginTemplate, user.UserName, time.Now().String(), loginType)
		if vars.EmailNoticeSetting != nil && vars.EmailNoticeSetting.Receivers != nil {
			for _, receiver := range vars.EmailNoticeSetting.Receivers {
				err = email.SendEmailNotice(ctx, receiver, kelvins.AppName, emailNotice)
				if err != nil {
					kelvins.ErrLogger.Info(ctx, "SendEmailNotice err, emailNotice: %v", emailNotice)
					return
				}
			}
		}
	}
	kelvins.GPool.SendJob(updateUserState)

	return result, retCode
}

func PasswordReset(ctx context.Context, req *users.PasswordResetRequest) int {
	user, err := repository.GetUserByUid("id,password_salt", int(req.GetUid()))
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
		_md5Pwd := hash.MD5EncodeToString(pwd)
		kelvins.ErrLogger.Errorf(ctx, "UpdateUserInfo err: %v, where: %+v, maps: %q", err, where, _md5Pwd)
		return code.ErrorServer
	}

	// 触发密码变更消息
	userPwdChangeNotify := func() {
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
			Content: json.MarshalToStringNoError(args.UserStateNotice{
				Uid:  user.Id,
				Time: util.ParseTimeOfStr(time.Now().Unix()),
			}),
		}
		_, retCode := pushNoticeService.PushMessage(ctx, businessMsg)
		if retCode != code.Success {
			kelvins.ErrLogger.Errorf(ctx, "Password Reset businessMsg: %v  notice send err: ", json.MarshalToStringNoError(businessMsg), code.GetMsg(retCode))
		}
	}
	kelvins.GPool.SendJob(userPwdChangeNotify)

	return code.Success
}

func UpdateUserLoginState(ctx context.Context, req *users.UpdateUserLoginStateRequest) int {
	if req.Uid <= 0 {
		return code.UserNotExist
	}
	user, err := repository.GetUserByUid("id", int(req.GetUid()))
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "GetUserByUid err: %v, uid: %v", err, req.GetUid())
		return code.ErrorServer
	}
	if user.Id <= 0 {
		return code.UserNotExist
	}

	return updateUserOnlineState(ctx, user.Id, req.GetState().GetContent())
}

func GetUserInfo(ctx context.Context, uid int) (*mysql.User, int) {
	if uid <= 0 {
		return nil, code.UserNotExist
	}
	user, err := repository.GetUserByUid("*", uid)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "GetUserByUid err: %v, uid: %v", err, uid)
		return user, code.ErrorServer
	}
	return user, code.Success
}

func GetUserInfoByPhone(ctx context.Context, countryCode, phone string) (*mysql.User, int) {
	if countryCode == "" || phone == "" {
		return nil, code.UserNotExist
	}
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
				kelvins.ErrLogger.Errorf(ctx, "UpdateUserLogisticsDeliveryByTx err: %v, where: %v", err, json.MarshalToStringNoError(where))
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
				kelvins.ErrLogger.Errorf(ctx, "CreateUserLogisticsDeliveryByTx err: %v, deliveryInfo: %v", err, json.MarshalToStringNoError(deliveryInfo))
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
			kelvins.ErrLogger.Errorf(ctx, "CreateUserLogisticsDelivery err: %v, deliveryInfo: %v", err, json.MarshalToStringNoError(deliveryInfo))
			return code.ErrorServer
		}
		return code.Success
	} else if req.OperationType == users.OperationType_UPDATE {
		if req.Info.Id <= 0 {
			return code.UserDeliveryInfoNotExist
		}
		deliveryInfoDB, err := repository.GetUserLogisticsDelivery("id", req.Uid, req.Info.Id)
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
				kelvins.ErrLogger.Errorf(ctx, "UpdateUserLogisticsDeliveryByTx err: %v, where: %v", err, json.MarshalToStringNoError(where))
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
				kelvins.ErrLogger.Errorf(ctx, "UpdateUserLogisticsDeliveryByTx err: %v,id: %v, deliveryInfo: %v", err, req.Info.Id, json.MarshalToStringNoError(deliveryInfo))
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
			kelvins.ErrLogger.Errorf(ctx, "UpdateUserLogisticsDelivery err: %v,id: %v, deliveryInfo: %v", err, req.Info.Id, json.MarshalToStringNoError(deliveryInfo))
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
		if len(list) == 0 {
			return result, code.UserDeliveryInfoNotExist
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
		infoDB, err := repository.GetUserLogisticsDelivery(sqlSelectUserDeliveryInfo, req.Uid, int64(req.UserDeliveryId))
		if err != nil {
			kelvins.ErrLogger.Errorf(ctx, "GetUserLogisticsDelivery err: %v, uid: %v,id: %v", err, req.Uid, req.UserDeliveryId)
			return result, code.ErrorServer
		}
		if infoDB.Id <= 0 {
			return result, code.UserDeliveryInfoNotExist
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

const sqlSelectFindUserInfoByPhone = "*"

func FindUserInfoByPhone(ctx context.Context, countryCode []string, phone []string) (result []*mysql.User, retCode int) {
	var err error
	result = make([]*mysql.User, 0)
	retCode = code.Success
	result, err = repository.FindUserInfoByPhone(sqlSelectFindUserInfoByPhone, countryCode, phone)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "FindUserInfoByPhone err: %v, countryCode: %v, phone: %v", err, countryCode, phone)
		retCode = code.ErrorServer
		return
	}
	return result, retCode
}

func UserAccountCharge(ctx context.Context, req *users.UserAccountChargeRequest) (retCode int) {
	retCode = code.Success
	userInfoList, err := repository.FindUserInfo("id,account_id,user_name", req.UidList)
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
	conn, err := util.GetGrpcClient(ctx, serverName)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "GetGrpcClient %v err: %v", serverName, err)
		return code.ErrorServer
	}
	//defer conn.Close()
	payClient := pay_business.NewPayBusinessServiceClient(conn)
	payReq := &pay_business.AccountChargeRequest{
		Owner:       accountIdList,
		AccountType: pay_business.AccountType(req.AccountType),
		CoinType:    pay_business.CoinType(req.CoinType),
		OutTradeNo:  req.OutTradeNo,
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
		kelvins.ErrLogger.Errorf(ctx, "AccountCharge req: %v, rsp: %v", json.MarshalToStringNoError(payReq), json.MarshalToStringNoError(payRsp))
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

	kelvins.GPool.SendJob(func() {
		// 发送登录邮件
		var un strings.Builder
		for _, v := range userInfoList {
			un.WriteString(v.UserName)
		}
		var coin string
		if req.CoinType == 0 {
			coin = "RMB"
		} else {
			coin = "USD"
		}
		emailNotice := fmt.Sprintf(args.UserAccountChargeTemplate, un.String(), time.Now(), req.Amount, coin)
		if vars.EmailNoticeSetting != nil && vars.EmailNoticeSetting.Receivers != nil {
			for _, receiver := range vars.EmailNoticeSetting.Receivers {
				err = email.SendEmailNotice(ctx, receiver, kelvins.AppName, emailNotice)
				if err != nil {
					kelvins.ErrLogger.Info(ctx, "SendEmailNotice err, emailNotice: %v", emailNotice)
					return
				}
			}
		}
	})

	return
}

func CheckUserDeliveryInfo(ctx context.Context, req *users.CheckUserDeliveryInfoRequest) (retCode int) {
	retCode = code.Success
	infoList, err := repository.CheckUserLogisticsDelivery(req.Uid, req.DeliveryIds)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "AccountCharge err: %v, req: %v", err, json.MarshalToStringNoError(req))
		retCode = code.ErrorServer
		return
	}
	if len(infoList) != len(req.DeliveryIds) {
		retCode = code.UserDeliveryInfoNotExist
		return
	}
	return
}

func updateUserOnlineState(ctx context.Context, uid int, content string) (retCode int) {
	retCode = code.Success
	if uid <= 0 {
		retCode = code.UserNotExist
		return
	}
	// 更新用户状态
	state := args.UserOnlineState{
		Uid:   uid,
		State: content,
		Time:  util.ParseTimeOfStr(time.Now().Unix()),
	}
	userLoginKey := fmt.Sprintf("%v-%d", args.CacheKeyUserOnlineSate, uid)
	err := cache.Set(kelvins.RedisConn, userLoginKey, json.MarshalToStringNoError(state), 24*3600)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "updateUserOnlineState err: %v, userLoginKey: %v", err, userLoginKey)
		retCode = code.ErrorServer
		return
	}
	return
}

func userLoginFailure(ctx context.Context, uid int) (failureFrequency, retCode int) {
	retCode = code.Success
	userLoginFailureKey := fmt.Sprintf("%v-%d", args.UserLoginFailureFrequency, uid)
	str, err := cache.Get(kelvins.RedisConn, userLoginFailureKey)
	if err != nil && err != cache.CacheNotFound {
		kelvins.ErrLogger.Errorf(ctx, "userLoginFailure err: %v, userLoginFailureKey: %v", err, userLoginFailureKey)
		return
	}
	if str != "" && err != cache.CacheNotFound {
		failureFrequency, err = strconv.Atoi(str)
		if err != nil {
			kelvins.ErrLogger.Errorf(ctx, "userLoginFailure strconv err: %v, str: %v", err, str)
			retCode = code.ErrorServer
			return
		}
	}
	if failureFrequency >= args.UserLoginFailureFrequencyMax {
		retCode = updateUserOnlineState(ctx, uid, args.UserOnlineStateForbiddenLogin)
		if retCode != code.Success {
			return
		}
		retCode = code.UserStateForbiddenLogin
		return
	}
	failureFrequency++
	err = cache.Set(kelvins.RedisConn, userLoginFailureKey, fmt.Sprintf("%d", failureFrequency), 24*3600)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "userLoginFailure err: %v, userLoginFailureKey: %v", err, userLoginFailureKey)
		retCode = code.ErrorServer
		return
	}
	return
}

func getUserOnlineState(ctx context.Context, uid int) (string, error) {
	var content string
	var err error
	userLoginKey := fmt.Sprintf("%v-%d", args.CacheKeyUserOnlineSate, uid)
	str, err := cache.Get(kelvins.RedisConn, userLoginKey)
	if err != nil && err != cache.CacheNotFound {
		kelvins.ErrLogger.Errorf(ctx, "getUserOnlineState err: %v, userLoginKey: %v", err, userLoginKey)
		return content, err
	}
	if err == cache.CacheNotFound {
		return "", nil
	}
	var state args.UserOnlineState
	err = json.Unmarshal(str, &state)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "getUserOnlineState Unmarshal err: %v, str: %v", err, str)
		return content, err
	}
	content = state.State
	return content, nil
}

func checkUserState(ctx context.Context, uid int) (retCode int) {
	retCode = code.Success
	if uid <= 0 {
		retCode = code.UserNotExist
		return
	}
	content, err := getUserOnlineState(ctx, uid)
	if err != nil {
		retCode = code.ErrorServer
		return
	}
	switch content {
	case args.UserOnlineStateOnline:
		return
	case args.UserOnlineStateForbiddenLogin:
		retCode = code.UserStateForbiddenLogin
		return
	default:
		return
	}
}

func CheckUserState(ctx context.Context, uidList []int64) (retCode int) {
	retCode = code.Success
	infoList, err := repository.FindUserInfo("id,state", uidList)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "AccountCharge err: %v, req: %v", err, json.MarshalToStringNoError(uidList))
		retCode = code.ErrorServer
		return
	}
	if len(infoList) == 0 {
		retCode = code.UserNotExist
		return
	}
	if len(infoList) != len(uidList) {
		retCode = code.UserNotExist
		return
	}
	for i := 0; i < len(infoList); i++ {
		if infoList[i].Id <= 0 {
			retCode = code.UserNotExist
			return
		}
		if infoList[i].State != 3 {
			retCode = code.UserStateNotVerify
			return
		}
		retCode = checkUserState(ctx, infoList[i].Id)
		if retCode != code.Success {
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
		kelvins.ErrLogger.Errorf(ctx, "FindUserInfo err: %v, req: %v", err, json.MarshalToStringNoError(req))
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

func ListUserInfo(ctx context.Context, req *users.ListUserInfoRequest) (result []*users.MobilePhone, retCode int) {
	retCode = code.Success
	result = make([]*users.MobilePhone, 0)
	userInfoList, err := repository.ListUserInfo("country_code,phone", int(req.PageMeta.PageSize), int(req.PageMeta.PageNum))
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "ListUserInfo err: %v, req: %v", err, json.MarshalToStringNoError(req))
		retCode = code.ErrorServer
		return
	}
	result = make([]*users.MobilePhone, len(userInfoList))
	for i := 0; i < len(userInfoList); i++ {
		info := &users.MobilePhone{
			CountryCode: userInfoList[i].CountryCode,
			Phone:       userInfoList[i].Phone,
		}
		result[i] = info
	}
	return
}

func SearchUserInfo(ctx context.Context, query string) (result []*users.SearchUserInfoEntry, retCode int) {
	result = make([]*users.SearchUserInfoEntry, 0)
	retCode = code.Success
	serverName := args.RpcServiceMicroMallSearch
	conn, err := util.GetGrpcClient(ctx, serverName)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "GetGrpcClient %v err: %v", serverName, err)
		return result, code.ErrorServer
	}
	client := search_business.NewSearchBusinessServiceClient(conn)
	rsp, err := client.UserSearch(ctx, &search_business.UserSearchRequest{Query: query})
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "UserSearch err: %v, query: %v", err, query)
		return result, code.ErrorServer
	}
	if rsp.Common.Code != search_business.RetCode_SUCCESS {
		kelvins.ErrLogger.Errorf(ctx, "UserSearch err: %v, query: %v, rsp: %v", err, query, json.MarshalToStringNoError(rsp))
		return result, code.ErrorServer
	}
	if len(rsp.List) == 0 {
		return
	}
	var userCountryCode = make([]string, 0)
	var userPhone = make([]string, 0)
	for i := 0; i < len(rsp.List); i++ {
		v := rsp.List[i]
		vv := strings.SplitN(v.GetPhone(), ":", 2)
		if len(vv) > 0 {
			userCountryCode = append(userCountryCode, vv[0])
		}
		if len(vv) > 1 {
			userPhone = append(userPhone, vv[1])
		}
	}
	userInfoList, retCode := FindUserInfoByPhone(ctx, userCountryCode, userPhone)
	if retCode != code.Success {
		return nil, code.ErrorServer
	}
	if len(userInfoList) == 0 {
		return
	}
	phoneToUserInfo := map[string]*mysql.User{}
	for i := 0; i < len(userInfoList); i++ {
		key := userInfoList[i].CountryCode + ":" + userInfoList[i].Phone
		phoneToUserInfo[key] = userInfoList[i]
	}
	result = make([]*users.SearchUserInfoEntry, 0)
	for i := 0; i < len(rsp.List); i++ {
		if rsp.List[i].GetPhone() == "" {
			continue
		}
		userInfo, ok := phoneToUserInfo[rsp.List[i].GetPhone()]
		if ok {
			entry := &users.SearchUserInfoEntry{
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
					InviterCode: userInfo.InviteCode,
					ContactAddr: userInfo.ContactAddr,
					Age:         int32(userInfo.Age),
					CreateTime:  userInfo.CreateTime.Format(time.RFC3339),
				},
				Score: rsp.List[i].Score,
			}
			result = append(result, entry)
		}
	}
	return
}
