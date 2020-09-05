package repository

import (
	"gitee.com/cristiane/micro-mall-users/model/mysql"
	"gitee.com/kelvins-io/kelvins"
)

func CreateVerifyCodeRecord(record *mysql.VerifyCodeRecord) (err error) {
	_, err = kelvins.XORM_DBEngine.Table(mysql.TableVerifyCodeRecord).Insert(record)
	return
}

func GetVerifyCode(businessType int, countryCode, phone, verifyCode string) (*mysql.VerifyCodeRecord, error) {
	var result mysql.VerifyCodeRecord
	var err error
	_, err = kelvins.XORM_DBEngine.Table(mysql.TableVerifyCodeRecord).
		Select("id,expire").
		Where("business_type = ? AND country_code = ? AND phone = ? AND verify_code = ?", businessType, countryCode, phone, verifyCode).
		Desc("id").
		Get(&result)
	return &result, err
}
