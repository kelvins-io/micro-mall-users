package configcenter

import (
	"fmt"
	"strings"
)

type ConfigCenterV2 struct {
	serverName string
}

func NewConfigCenterV2(serverName string) *ConfigCenterV2 {
	return &ConfigCenterV2{
		serverName: serverName,
	}
}

// 获取证书 .pem 路径
func (c *ConfigCenterV2) GetCertPemPath() (string, error) {
	return c.getPath(DefaultCertPemFileName)
}

// 获取证书 .key 路径
func (c *ConfigCenterV2) GetCertKeyPath() (string, error) {
	return c.getPath(DefaultCertKeyFileName)
}

func (c *ConfigCenterV2) getPath(key string) (string, error) {
	names := strings.Split(c.serverName, "-")
	if len(names) < 1 {
		return "", fmt.Errorf("ConfigCenterV2.getPath is empty.")
	}
	serverTag := names[0] + "-cert"

	return strings.Join([]string{DefaultGlobalPath, "certs", serverTag, key}, "/"), nil
}
