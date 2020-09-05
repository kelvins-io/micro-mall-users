package env

import (
	"errors"
	"os"
	"strings"
)

const (
	EnvName     = "GO_ENV"  //标识名
	DevMode     = "dev"     //开发环境
	TestMode    = "test"    //测试环境
	ReleaseMode = "release" //预发布环境
	ProdMode    = "prod"    //生产环境
)

// 获取当前环境的运行模式
func GetMode() (env string, err error) {
	env = strings.ToLower(os.Getenv(EnvName))
	if env == "" {
		err = errors.New("Can not find ENV '" + EnvName + "'")
	}

	return env, err
}

// 判断是否为开发模式
func IsDevMode() bool {
	return checkMode(DevMode)
}

// 判断是否测试模式
func IsTestMode() bool {
	return checkMode(TestMode)
}

// 判断是否预发布模式
func IsReleaseMode() bool {
	return checkMode(ReleaseMode)
}

// 判断是否生产环境
func IsProdMode() bool {
	return checkMode(ProdMode)
}

// 检查当前运行模式
func checkMode(mode string) bool {
	env := strings.ToLower(os.Getenv(EnvName))
	return env == mode
}
