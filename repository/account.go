package repository

import (
	"gitee.com/cristiane/micro-mall-users/model/mysql"
	"gitee.com/kelvins-io/kelvins"
)

func CreateAccount(model *mysql.Account) (err error) {
	_, err = kelvins.XORM_DBEngine.Table(mysql.TableAccount).Insert(model)
	return
}
