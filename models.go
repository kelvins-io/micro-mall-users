package models

import (
	"time"
)

type Merchant struct {
	MerchantId   int64     `xorm:"not null pk autoincr comment('商户号ID') BIGINT"`
	MerchantCode string    `xorm:"not null comment('商户唯一code') index CHAR(36)"`
	Uid          int64     `xorm:"not null comment('用户ID') unique BIGINT"`
	RegisterAddr string    `xorm:"not null comment('注册地址') TEXT"`
	HealthCardNo string    `xorm:"not null comment('健康证号') CHAR(30)"`
	Identity     int       `xorm:"comment('身份属性，1-临时店员，2-正式店员，3-经理，4-店长') TINYINT"`
	State        int       `xorm:"comment('状态，0-未审核，1-审核中，2-审核不通过，3-已审核') TINYINT"`
	TaxCardNo    string    `xorm:"comment('纳税账户号') CHAR(30)"`
	CreateTime   time.Time `xorm:"not null default CURRENT_TIMESTAMP comment('创建时间') DATETIME"`
	UpdateTime   time.Time `xorm:"not null default CURRENT_TIMESTAMP comment('修改时间') DATETIME"`
}

type User struct {
	Id           int64     `xorm:"pk autoincr comment('自增ID') BIGINT"`
	AccountId    string    `xorm:"not null comment('账户ID，全局唯一') unique CHAR(36)"`
	UserName     string    `xorm:"not null comment('用户名') index VARCHAR(255)"`
	Password     string    `xorm:"not null comment('用户密码md5值') VARCHAR(255)"`
	PasswordSalt string    `xorm:"comment('密码salt值') VARCHAR(255)"`
	Sex          int       `xorm:"comment('性别，1-男，2-女') TINYINT(1)"`
	Phone        string    `xorm:"comment('手机号') unique(country_code_phone_index) CHAR(11)"`
	CountryCode  string    `xorm:"comment('手机区号') unique(country_code_phone_index) CHAR(5)"`
	Email        string    `xorm:"comment('邮箱') index VARCHAR(255)"`
	State        int       `xorm:"comment('状态，0-未激活，1-审核中，2-审核未通过，3-已审核') TINYINT(1)"`
	IdCardNo     string    `xorm:"comment('身份证号') unique CHAR(18)"`
	Inviter      int64     `xorm:"comment('邀请人uid') BIGINT"`
	InviteCode   string    `xorm:"comment('邀请码') CHAR(20)"`
	CreateTime   time.Time `xorm:"not null default CURRENT_TIMESTAMP comment('创建时间') DATETIME"`
	UpdateTime   time.Time `xorm:"not null default CURRENT_TIMESTAMP comment('修改时间') DATETIME"`
	ContactAddr  string    `xorm:"comment('联系地址') TEXT"`
	Age          int       `xorm:"comment('年龄') INT"`
}

type UserLogisticsDelivery struct {
	Id           int64     `xorm:"pk autoincr comment('自增ID') BIGINT"`
	Uid          int64     `xorm:"comment('用户ID') index BIGINT"`
	DeliveryUser string    `xorm:"comment('交付人') VARCHAR(512)"`
	CountryCode  string    `xorm:"default '86' comment('区号') VARCHAR(10)"`
	Phone        string    `xorm:"comment('手机号') VARCHAR(255)"`
	Area         string    `xorm:"comment('交付区域') VARCHAR(255)"`
	AreaDetailed string    `xorm:"comment('详细地址') TEXT"`
	Label        string    `xorm:"comment('标签，多个以|分割开') TEXT"`
	IsDefault    int       `xorm:"default 0 comment('是否为默认，1-默认') TINYINT"`
	CreateTime   time.Time `xorm:"not null default CURRENT_TIMESTAMP comment('创建时间') DATETIME"`
	UpdateTime   time.Time `xorm:"not null default CURRENT_TIMESTAMP comment('更新时间') DATETIME"`
}
