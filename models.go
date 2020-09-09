package models

import (
	"time"
)

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

type MerchantInfo struct {
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

type ShopBusinessInfo struct {
	ShopId           int64     `xorm:"not null pk autoincr comment('店铺ID') BIGINT"`
	NickName         string    `xorm:"not null comment('简称') unique(legal_person_nick_name_index) VARCHAR(512)"`
	ShopCode         string    `xorm:"not null comment('店铺唯一code') unique CHAR(36)"`
	FullName         string    `xorm:"not null comment('店铺全称') TEXT"`
	RegisterAddr     string    `xorm:"not null comment('注册地址') TEXT"`
	BusinessAddr     string    `xorm:"not null comment('实际经营地址') TEXT"`
	LegalPerson      int64     `xorm:"not null comment('店铺法人') index unique(legal_person_nick_name_index) BIGINT"`
	BusinessLicense  string    `xorm:"not null comment('经营许可证') CHAR(36)"`
	TaxCardNo        string    `xorm:"not null comment('纳税号') CHAR(36)"`
	BusinessDesc     string    `xorm:"not null comment('经营描述') TEXT"`
	SocialCreditCode string    `xorm:"not null comment('统一社会信用代码') CHAR(36)"`
	OrganizationCode string    `xorm:"not null comment('组织机构代码') CHAR(36)"`
	State            int       `xorm:"not null default 0 comment('状态，0-未审核，1-审核不通过，2-审核通过') TINYINT"`
	CreateTime       time.Time `xorm:"not null default CURRENT_TIMESTAMP comment('创建时间') DATETIME"`
	UpdateTime       time.Time `xorm:"not null default CURRENT_TIMESTAMP comment('修改时间') DATETIME"`
}

type SkuInventory struct {
	Id         int64     `xorm:"pk autoincr comment('商品库存ID') BIGINT"`
	SkuCode    string    `xorm:"not null comment('商品编码') unique CHAR(64)"`
	Amount     int64     `xorm:"comment('库存数量') BIGINT"`
	Price      string    `xorm:"comment('入库单价') DECIMAL(32,16)"`
	ShopId     int64     `xorm:"not null comment('所属店铺ID') index BIGINT"`
	OpUid      int64     `xorm:"not null comment('操作用户UID') BIGINT"`
	OpIp       string    `xorm:"comment('操作的IP') CHAR(16)"`
	CreateTime time.Time `xorm:"not null default CURRENT_TIMESTAMP comment('创建时间') DATETIME"`
	UpdateTime time.Time `xorm:"not null default CURRENT_TIMESTAMP comment('修改时间') DATETIME"`
}

type SkuProperty struct {
	Id           int64  `xorm:"pk autoincr comment('ID') BIGINT"`
	SkuCode      string `xorm:"not null comment('商品唯一编号') index CHAR(64)"`
	SkuType1th   int    `xorm:"comment('商品分类，一级分类') INT"`
	SkuType2th   int    `xorm:"comment('商品分类，二级分类') INT"`
	SkuType3th   int    `xorm:"comment('商品分类，三级分类') INT"`
	StoragePrice string `xorm:"comment('商品入库价格') DECIMAL(10,2)"`
	Price        string `xorm:"comment('商品当前价格') DECIMAL(10,2)"`
	PricePrev    string `xorm:"comment('商品之前价格') DECIMAL(10)"`
	SkuName      string `xorm:"comment('商品名称') index VARCHAR(255)"`
	SkuDesc      string `xorm:"comment('商品描述') TEXT"`
	Production   string `xorm:"comment('生产企业') VARCHAR(1024)"`
	Supplier     string `xorm:"comment('供应商') VARCHAR(1024)"`
}

type Transaction struct {
	Id              int64     `xorm:"pk comment('交易ID') BIGINT"`
	FromAccountCode string    `xorm:"not null default '0' comment('转出账户ID') CHAR(36)"`
	FromBalance     string    `xorm:"default 0.0000000000000000 comment('转出后账户余额') DECIMAL(32,16)"`
	ToAccountCode   string    `xorm:"not null default '0' comment('转入账户ID') CHAR(36)"`
	ToBalance       string    `xorm:"comment('转入后账户余额') DECIMAL(32,16)"`
	Amount          string    `xorm:"comment('交易金额') DECIMAL(32,16)"`
	Meta            string    `xorm:"comment('转账说明') VARCHAR(255)"`
	Scene           string    `xorm:"comment('支付场景') VARCHAR(64)"`
	OpUid           int64     `xorm:"not null comment('操作用户UID') BIGINT"`
	OpIp            string    `xorm:"comment('操作的IP') VARCHAR(16)"`
	TxId            string    `xorm:"comment('对应交易号') CHAR(36)"`
	Fingerprint     string    `xorm:"not null comment('防篡改指纹') VARCHAR(32)"`
	PayType         int       `xorm:"default 0 comment('支付方式，0系统操作，1-银行卡，2-信用卡,3-支付宝,4-微信支付,5-京东支付') TINYINT"`
	PayDesc         string    `xorm:"comment('支付方式描述') VARCHAR(36)"`
	CreateTime      time.Time `xorm:"not null default CURRENT_TIMESTAMP comment('创建时间') DATETIME"`
	UpdateTime      time.Time `xorm:"not null default CURRENT_TIMESTAMP comment('修改时间') DATETIME"`
}

type UserInfo struct {
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

type VerifyCodeRecord struct {
	Id           int64     `xorm:"pk autoincr comment('自增id') BIGINT"`
	Uid          int64     `xorm:"not null comment('用户UID') BIGINT"`
	BusinessType int       `xorm:"comment('验证类型，1-注册登录，2-购买商品') TINYINT"`
	VerifyCode   string    `xorm:"comment('验证码') index CHAR(6)"`
	Expire       int       `xorm:"comment('过期时间unix') INT"`
	CountryCode  string    `xorm:"comment('验证码下发手机国际码') index(country_code_phone_index) CHAR(5)"`
	Phone        string    `xorm:"comment('验证码下发手机号') index(country_code_phone_index) CHAR(11)"`
	Email        string    `xorm:"comment('验证码下发邮箱') index VARCHAR(255)"`
	CreateTime   time.Time `xorm:"not null default CURRENT_TIMESTAMP comment('创建时间') DATETIME"`
	UpdateTime   time.Time `xorm:"not null default CURRENT_TIMESTAMP comment('修改时间') DATETIME"`
}