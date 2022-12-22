package env

import (
	"ws/framework/plugin/nacos"
)

// NacosConfig nacos中的配置
var NacosConfig = &nacosConfig{}

type nacosConfig struct {
	MysqlDataBase mysqlDataBase `yaml:"MysqlDataBase"`
	LogLevel      string        `yaml:"LogLevel"`
	HlrLookup     hlrLookup     `yaml:"HlrLookup"`
}

func init() {
	nacos.New(NacosConfig)
}
