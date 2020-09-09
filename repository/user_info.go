package repository

import (
	"gitee.com/cristiane/micro-mall-users/model/mysql"
	"gitee.com/kelvins-io/kelvins"
)

func CreateUser(user *mysql.UserInfo) (err error) {
	_, err = kelvins.XORM_DBEngine.Table(mysql.TableUser).Insert(user)
	return
}

func GetUserByUserName(username string) (*mysql.UserInfo, error) {
	var user mysql.UserInfo
	var err error
	_, err = kelvins.XORM_DBEngine.Table(mysql.TableUser).Where("user_name = ?", username).Get(&user)
	return &user, err
}

func GetUserByUid(uid int) (*mysql.UserInfo, error) {
	var user mysql.UserInfo
	var err error
	_, err = kelvins.XORM_DBEngine.Table(mysql.TableUser).Where("id = ?", uid).Get(&user)
	return &user, err
}

func UpdateUserInfo(query, maps map[string]interface{}) (err error) {
	_, err = kelvins.XORM_DBEngine.Table(mysql.TableUser).Where(query).Update(maps)
	return
}

func GetUserByPhone(countryCode, phone string) (*mysql.UserInfo, error) {
	var user mysql.UserInfo
	var err error
	_, err = kelvins.XORM_DBEngine.Table(mysql.TableUser).Where("country_code = ? and phone = ?", countryCode, phone).Get(&user)
	return &user, err
}

func CheckUserExistById(id int) (exist bool, err error) {
	var user mysql.UserInfo
	exist, err = kelvins.XORM_DBEngine.Table(mysql.TableUser).
		Select("id").
		Where("id = ?", id).Get(&user)
	if err != nil {
		return false, err
	}
	if user.Id != 0 {
		return true, nil
	}

	return false, nil
}

func CheckUserExistByPhone(countryCode, phone string) (exist bool, err error) {
	var user mysql.UserInfo
	exist, err = kelvins.XORM_DBEngine.Table(mysql.TableUser).
		Select("id").
		Where("country_code = ? and phone = ?", countryCode, phone).Get(&user)
	if err != nil {
		return false, err
	}
	if user.Id != 0 {
		return true, nil
	}

	return false, nil
}
