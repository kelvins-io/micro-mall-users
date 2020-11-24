package mysql

import (
	"database/sql"
	"time"
)

const (
	TableUser                  = "user"
	TableMerchantInfo          = "merchant"
	TableAccount               = "account"
	TableVerifyCodeRecord      = "verify_code_record"
	TableUserLogisticsDelivery = "user_logistics_delivery"
)

type UserLogisticsDelivery struct {
	Id           int64     `xorm:"pk autoincr comment('自增ID') BIGINT"`
	Uid          int64     `xorm:"comment('用户ID') index BIGINT"`
	DeliveryUser string    `xorm:"comment('交付人') VARCHAR(512)"`
	CountryCode  string    `xorm:"not null default '86' comment('区号') VARCHAR(10)"`
	Phone        string    `xorm:"comment('手机号') VARCHAR(255)"`
	Area         string    `xorm:"comment('交付区域') VARCHAR(255)"`
	AreaDetailed string    `xorm:"comment('详细地址') TEXT"`
	Label        string    `xorm:"comment('标签，多个以|分割开') TEXT"`
	IsDefault    int       `xorm:"default 0 comment('是否为默认，1-默认') TINYINT"`
	CreateTime   time.Time `xorm:"not null default CURRENT_TIMESTAMP comment('创建时间') DATETIME"`
	UpdateTime   time.Time `xorm:"not null default CURRENT_TIMESTAMP comment('更新时间') DATETIME"`
}

type Account struct {
	AccountCode string    `xorm:"not null pk comment('账户主键') CHAR(50)"`
	Owner       string    `xorm:"not null comment('账户所有者') unique(account_index) CHAR(36)"`
	Balance     string    `xorm:"comment('账户余额') DECIMAL(32,16)"`
	CoinType    int       `xorm:"not null default 1 comment('币种类型，1-rmb，2-usdt') unique(account_index) TINYINT"`
	CoinDesc    string    `xorm:"comment('币种描述') VARCHAR(64)"`
	State       int       `xorm:"comment('状态，1无效，2锁定，3正常') TINYINT"`
	AccountType int       `xorm:"not null comment('账户类型，1-个人账户，2-公司账户，3-系统账户') unique(account_index) TINYINT"`
	CreateTime  time.Time `xorm:"not null default CURRENT_TIMESTAMP comment('创建时间') index DATETIME"`
	UpdateTime  time.Time `xorm:"not null default CURRENT_TIMESTAMP comment('更新时间') DATETIME"`
}

type VerifyCodeRecord struct {
	Id           int       `xorm:"'id' not null pk autoincr comment('自增id') INT"`
	Uid          int       `xorm:"'uid' not null comment('用户UID') INT"`
	BusinessType int       `xorm:"'business_type' comment('验证类型，1-注册登录，2-购买商品') TINYINT"`
	VerifyCode   string    `xorm:"'verify_code' comment('验证码') index CHAR(6)"`
	Expire       int       `xorm:"'expire' comment('过期时间unix') INT"`
	CountryCode  string    `xorm:"'country_code' comment('验证码下发手机国际码') index(country_code_phone_index) CHAR(5)"`
	Phone        string    `xorm:"'phone' comment('验证码下发手机号') index(country_code_phone_index) CHAR(11)"`
	Email        string    `xorm:"'email' comment('验证码下发邮箱') index VARCHAR(255)"`
	CreateTime   time.Time `xorm:"'create_time' comment('创建时间') DATETIME"`
	UpdateTime   time.Time `xorm:"'update_time' comment('修改时间') DATETIME"`
}
type User struct {
	Id           int            `xorm:"'id' not null pk autoincr comment('用户id') INT"`
	AccountId    string         `xorm:"'account_id' not null comment('账户ID，全局唯一') unique CHAR(36)"`
	UserName     string         `xorm:"'user_name' not null comment('用户名') index VARCHAR(255)"`
	Password     string         `xorm:"'password' not null comment('用户密码md5值') VARCHAR(255)"`
	PasswordSalt string         `xorm:"'password_salt' comment('密码salt值') VARCHAR(255)"`
	Sex          int            `xorm:"'sex' comment('性别，1-男，2-女') TINYINT"`
	Phone        string         `xorm:"'phone' comment('手机号') unique(country_code_phone_index) CHAR(11)"`
	CountryCode  string         `xorm:"'country_code' comment('手机区号') unique(country_code_phone_index) CHAR(5)"`
	Email        string         `xorm:"'email' comment('邮箱') index VARCHAR(255)"`
	State        int            `xorm:"'state' comment('状态，0-未激活，1-审核中，2-审核未通过，3-已审核') TINYINT"`
	IdCardNo     sql.NullString `xorm:"'id_card_no' comment('身份证号') unique CHAR(18)"`
	Inviter      int            `xorm:"'inviter' comment('邀请人uid') INT"`
	InviteCode   string         `xorm:"'invite_code' comment('邀请码') CHAR(20)"`
	ContactAddr  string         `xorm:"'contact_addr' comment('联系地址') TEXT"`
	Age          int            `xorm:"'age' comment('年龄') INT"`
	CreateTime   time.Time      `xorm:"'create_time' not null comment('创建时间') DATETIME"`
	UpdateTime   time.Time      `xorm:"'update_time' not null comment('修改时间') DATETIME"`
}

type Merchant struct {
	MerchantId   int       `xorm:"'merchant_id' not null pk autoincr comment('商户号ID') INT"`
	MerchantCode string    `xorm:"'merchant_code' not null comment('商户唯一code') index CHAR(36)"`
	Uid          int       `xorm:"'uid' not null comment('用户ID') unique INT"`
	RegisterAddr string    `xorm:"'register_addr' not null comment('注册地址') TEXT"`
	HealthCardNo string    `xorm:"'health_card_no' not null comment('健康证号') index CHAR(30)"`
	Identity     int       `xorm:"'identity' comment('身份属性，1-临时店员，2-正式店员，3-经理，4-店长') TINYINT"`
	State        int       `xorm:"'state' comment('状态，0-未审核，1-审核中，2-审核不通过，3-已审核') TINYINT"`
	TaxCardNo    string    `xorm:"'tax_card_no' comment('纳税账户号') index CHAR(30)"`
	CreateTime   time.Time `xorm:"'create_time' default CURRENT_TIMESTAMP comment('创建时间') DATETIME"`
	UpdateTime   time.Time `xorm:"'update_time' default CURRENT_TIMESTAMP comment('修改时间') DATETIME"`
}
