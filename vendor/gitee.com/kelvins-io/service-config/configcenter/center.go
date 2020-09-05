package configcenter

import (
	"errors"
	"gitee.com/kelvins-io/common/env"
	"gitee.com/kelvins-io/service-config/configcenter/read"
	"github.com/tidwall/gjson"
	"log"
	"strings"
)

type ConfigCenter struct {
	parse      gjson.Result
	serverName string
	configName string
}

const (
	DefaultGlobalPath     = "/usr/local/etc/global-conf"
	DefaultGlobalFileName = "/usr/local/etc/global-conf/config.json"
	DefaultConfigRootPath = "/usr/local/etc"

	DefaultCertPemPath     = "certs/server.pem"
	DefaultCertKeyPath     = "certs/server-key.pem"
	DefaultCertPemFileName = "server.pem"
	DefaultCertKeyFileName = "server-key.pem"

	ServerName     = "server_name"
	ServerPort     = "server_port"
	ConfigRootPath = "config_root_path"
	CertPemPath    = "cert_pem_path"
	CertKeyPath    = "cert_key_path"
	CertServerName = "cert_server_name"
)

func NewConfigCenter(serverName string) *ConfigCenter {
	fileRead := read.NewFileRead()
	jsonByte, err := fileRead.Read(DefaultGlobalFileName)
	if err != nil {
		log.Fatalf("fileRead.Read err: %v", err)
	}

	env, err := env.GetMode()
	if err != nil {
		log.Fatalf("env.GetMode err: %v", err)
	}

	return &ConfigCenter{
		parse:      gjson.Parse(string(jsonByte[:])),
		serverName: serverName,
		configName: env + "/" + serverName,
	}
}

// 获取服务名称
func (c *ConfigCenter) GetServerName() (string, error) {
	value := c.parse.Get(c.getPath(ServerName))
	if !value.Exists() {
		return ``, errors.New(c.serverName + " server_name not exist")
	}

	return value.String(), nil
}

// 获取服务端口
func (c *ConfigCenter) GetServerPort() (uint64, error) {
	value := c.parse.Get(c.getPath(ServerPort))
	if !value.Exists() {
		return 0, errors.New(c.serverName + " server_port not exist")
	}

	return value.Uint(), nil
}

// 获取服务配置目录
func (c *ConfigCenter) GetServerConfigPath() (string, error) {
	value := c.parse.Get(c.getPath(ConfigRootPath))
	if !value.Exists() {
		return ``, errors.New(c.serverName + " config_root_path not exist")
	}

	return strings.Join([]string{value.String(), c.configName}, "/"), nil
}

func (c *ConfigCenter) MustGetServerConfigPath() string {
	var paths []string
	value := c.parse.Get(c.getPath(ConfigRootPath))

	if !value.Exists() {
		paths = []string{DefaultConfigRootPath, c.configName}
	} else {
		paths = []string{value.String(), c.configName}
	}

	return strings.Join(paths, "/")
}

// 获取证书服务名称
func (c *ConfigCenter) GetCertServerName() (string, error) {
	value := c.parse.Get(c.getPath(CertServerName))
	if !value.Exists() {
		return "", errors.New(c.serverName + " cert_server_name not exist")
	}

	return value.String(), nil
}

// 获取证书 .pem 路径
func (c *ConfigCenter) GetCertPemPath() (string, error) {
	value := c.parse.Get(c.getPath(CertPemPath))
	if !value.Exists() {
		return ``, errors.New(c.serverName + " cert_pem_path not exist")
	}

	return value.String(), nil
}

func (c *ConfigCenter) MustGetCertPemPath() string {
	value := c.parse.Get(c.getPath(CertPemPath))
	if !value.Exists() {
		return DefaultCertPemPath
	}

	return value.String()
}

// 获取证书 .key 路径
func (c *ConfigCenter) GetCertKeyPath() (string, error) {
	value := c.parse.Get(c.getPath(CertKeyPath))
	if !value.Exists() {
		return ``, errors.New(c.serverName + " cert_key_path not exist")
	} else {
		return value.String(), nil
	}
}

func (c *ConfigCenter) MustGetCertKeyPath() string {
	value := c.parse.Get(c.getPath(CertKeyPath))
	if !value.Exists() {
		return DefaultCertKeyPath
	}

	return value.String()
}

func (c *ConfigCenter) getPath(key string) string {
	return strings.Join([]string{c.serverName, key}, ".")
}
