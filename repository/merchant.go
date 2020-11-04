package repository

import (
	"gitee.com/cristiane/micro-mall-users/model/mysql"
	"gitee.com/kelvins-io/kelvins"
)

func CreateMerchantsMaterial(model *mysql.Merchant) (err error) {
	_, err = kelvins.XORM_DBEngine.Table(mysql.TableMerchantInfo).Insert(model)
	return
}

func UpdateMerchantsMaterial(query, maps map[string]interface{}) (err error) {
	_, err = kelvins.XORM_DBEngine.Table(mysql.TableMerchantInfo).Where(query).Update(maps)
	return
}

func GetMerchantsMaterialByUid(uid int) (*mysql.Merchant, error) {
	var model mysql.Merchant
	_, err := kelvins.XORM_DBEngine.Table(mysql.TableMerchantInfo).Where("uid = ?", uid).Get(&model)
	return &model, err
}

func GetMerchantsMaterial(merchantId int) (*mysql.Merchant, error) {
	var model mysql.Merchant
	_, err := kelvins.XORM_DBEngine.Table(mysql.TableMerchantInfo).Where("merchant_id = ?", merchantId).Get(&model)
	return &model, err
}
