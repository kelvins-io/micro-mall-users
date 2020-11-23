package repository

import (
	"gitee.com/cristiane/micro-mall-users/model/mysql"
	"gitee.com/kelvins-io/kelvins"
)

func CreateUserLogisticsDelivery(model *mysql.UserLogisticsDelivery) error {
	_, err := kelvins.XORM_DBEngine.Table(mysql.TableUserLogisticsDelivery).Insert(model)
	return err
}

func GetUserLogisticsDelivery(sqlSelect string, uid, id int64) (*mysql.UserLogisticsDelivery, error) {
	var model mysql.UserLogisticsDelivery
	_, err := kelvins.XORM_DBEngine.Table(mysql.TableUserLogisticsDelivery).Select(sqlSelect).Where("id = ? AND uid = ?", id, uid).Get(&model)
	return &model, err
}

func GetUserLogisticsDeliveryList(sqlSelect string, uid int64) ([]mysql.UserLogisticsDelivery, error) {
	var result = make([]mysql.UserLogisticsDelivery, 0)
	session := kelvins.XORM_DBEngine.Table(mysql.TableUserLogisticsDelivery).Select(sqlSelect).Where("uid = ?", uid)
	err := session.Find(&result)
	return result, err
}

func UpdateUserLogisticsDelivery(where, maps interface{}) (int64, error) {
	return kelvins.XORM_DBEngine.Table(mysql.TableUserLogisticsDelivery).Where(where).Update(maps)
}
