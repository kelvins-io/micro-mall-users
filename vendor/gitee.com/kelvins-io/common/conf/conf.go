package conf

import (
	"gitee.com/kelvins-io/service-config/configcenter"
	"log"
	"strings"
)

type ConfService struct {
	configCenter *configcenter.ConfigCenter
}

func NewConfService(serverName string) *ConfService {
	return &ConfService{configCenter: configcenter.NewConfigCenter(serverName)}
}

// 获取当前环境的配置路径
func (c *ConfService) GetConfigModePath(fileName string) string {
	configPath := c.configCenter.MustGetServerConfigPath()
	paths := []string{configPath, fileName}

	return strings.Join(paths, "/")
}

// 获取服务名称
func (c *ConfService) GetServerName() string {
	serverName, err := c.configCenter.GetServerName()
	if err != nil {
		log.Printf("conf.GetServerName err: %v", serverName)
		return ""
	}

	return serverName
}

//  获取服务端口号
func (c *ConfService) GetServerPort() uint64 {
	port, err := c.configCenter.GetServerPort()
	if err != nil {
		log.Printf("conf.GetServerPort err: %v", err)
		return 0
	}

	return port
}

// 获取证书服务名称
func (c *ConfService) GetCertServerName() string {
	serverName, err := c.configCenter.GetCertServerName()
	if err != nil {
		log.Printf("conf.GetServerName err: %v", err)
		return ""
	}

	return serverName
}

// 获取服务证书 PEM 路径
func (c *ConfService) GetCertPemPath() string {
	return c.configCenter.MustGetCertPemPath()
}

// 获取服务证书 PEM 的完整路径
func (c *ConfService) GetCertPemFullPath() string {
	return c.GetConfigModePath(c.GetCertPemPath())
}

// 获取服务证书 KEY 路径
func (y *ConfService) GetCertKeyPath() string {
	return y.configCenter.MustGetCertKeyPath()
}

// 获取服务证书 KEY 的完整路径
func (y *ConfService) GetCertKeyFullPath() string {
	return y.GetConfigModePath(y.GetCertKeyPath())
}
