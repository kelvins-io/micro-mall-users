# micro-mall-users

#### 介绍
微商城-用户系统

#### 软件架构
grpc

#### 框架，库依赖
kelvins框架支持（gRPC，cron，queue，web支持）：https://gitee.com/kelvins-io/kelvins   
g2cache缓存库支持（两级缓存）：https://gitee.com/kelvins-io/g2cache   

#### 安装教程
1.仅构建  sh build.sh   
2 运行  sh build-run.sh   
3 停止 sh stop.sh

#### 使用说明
配置参考
```toml
[kelvins-server]
Environment = "dev"

[kelvins-logger]
RootPath = "./logs"
Level = "debug"

[kelvins-auth]
Token = "c9VW6ForlmzdeDkZE2i8"
TransportSecurity = false

[kelvins-mysql]
Host = "127.0.0.1:3306"
UserName = "root"
Password = "xxxx"
DBName = "micro_mall_user"
Charset = "utf8"
PoolNum =  10
MaxIdleConns = 5
ConnMaxLifeSecond = 3600
MultiStatements = true
ParseTime = true

[kelvins-jwt]
Secret = "&WJof0jaY4ByTHR2"
TokenExpireSecond = 2592000

[kelvins-redis]
Host = "127.0.0.1:6379"
Password = "xxx"
DB = 1
PoolNum = 10

[queue-user-register-notice]
Broker = "amqp://micro-mall:szJ9aePR@localhost:5672/micro-mall"
DefaultQueue = "user_register_notice"
ResultBackend = "redis://xxx@127.0.0.1:6379/8"
ResultsExpireIn = 3600
Exchange = "user_register_notice"
ExchangeType = "direct"
BindingKey = "user_register_notice"
PrefetchCount = 5
TaskRetryCount = 3
TaskRetryTimeout = 3600

[queue-user-state-notice]
Broker = "amqp://micro-mall:szJ9aePR@localhost:5672/micro-mall"
DefaultQueue = "user_state_notice"
ResultBackend = "redis://xxx@127.0.0.1:6379/8"
ResultsExpireIn = 3600
Exchange = "user_state_notice"
ExchangeType = "direct"
BindingKey = "user_state_notice"
PrefetchCount = 5
TaskRetryCount = 3
TaskRetryTimeout = 3600

[email-config]
Enable = false
User = "xxx@qq.com"
Password = "xx"
Host = "smtp.qq.com"
Port = "465"
```

#### 参与贡献

1.  Fork 本仓库
2.  新建 Feat_xxx 分支
3.  提交代码
4.  新建 Pull Request

