package repository

import (
	"gitee.com/cristiane/micro-mall-users/model/mysql"
	"gitee.com/kelvins-io/kelvins"
)

func CreateMerchantsMaterial(model *mysql.MerchantInfo) (err error) {
	_, err = kelvins.XORM_DBEngine.Table(mysql.TableMerchantInfo).Insert(model)
	return
}

func UpdateMerchantsMaterial(query, maps map[string]interface{}) (err error) {
	_, err = kelvins.XORM_DBEngine.Table(mysql.TableMerchantInfo).Where(query).Update(maps)
	return
}

func GetMerchantsMaterialByUid(uid int) (*mysql.MerchantInfo, error) {
	var model mysql.MerchantInfo
	_, err := kelvins.XORM_DBEngine.Table(mysql.TableMerchantInfo).Where("merchant_id = ?", uid).Get(&model)
	return &model, err
}
