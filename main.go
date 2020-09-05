package main

import (
	"gitee.com/cristiane/micro-mall-users/vars"
	"gitee.com/kelvins-io/kelvins"
	"gitee.com/cristiane/micro-mall-users/startup"
	"gitee.com/kelvins-io/kelvins/app"
)

const APP_NAME = "micro-mall-users"

func main() {
	application := &kelvins.GRPCApplication{
		Application: &kelvins.Application{
			LoadConfig: startup.LoadConfig,
			SetupVars:  startup.SetupVars,
			Name:       APP_NAME,
		},
		RegisterGRPCServer: startup.RegisterGRPCServer,
		RegisterGateway:    startup.RegisterGateway,
		RegisterHttpRoute:  startup.RegisterHttpRoute,
	}
	app.RunGRPCApplication(application)
	vars.App = application
}
