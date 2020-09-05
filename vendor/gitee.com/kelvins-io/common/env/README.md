# 配置模块

该包用于接管全局公共配置类的一些相关操作，统一管理

## 环境变量
```
export GO_ENV=dev
```

## 示例

- 获取当前环境的运行模式：`env.GetMode()`
- 判断是否为开发模式：`env.IsDevMode()`
- 判断是否为测试模式：`env.IsTestMode()`
- 判断是否预发布模式：`env.IsReleaseMode()`
- 判断是否生产环境：`env.IsProdMode()`
