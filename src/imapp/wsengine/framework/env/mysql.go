package env

import (
	"fmt"
)

type mysqlDataBase struct {
	Username          string `yaml:"Username"`
	PassWord          string `yaml:"Password"`
	Host              string `yaml:"Host"`
	Port              int    `yaml:"Port"`
	DataBaseName      string `yaml:"DataBaseName"`
	CharSet           string `yaml:"CharSet"`
	LogLevel          string `yaml:"LogLevel"`
	MaxIdleConnection int    `yaml:"MaxIdleConnection"`
	MaxOpenConnection int    `yaml:"MaxOpenConnection"`
	MaxLifeTime       int    `yaml:"MaxLifeTime"`
}

// DSN .
func (m *mysqlDataBase) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		m.Username, m.PassWord, m.Host, m.Port, m.DataBaseName, m.CharSet)
}
