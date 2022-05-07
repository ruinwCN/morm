package conf

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

var (
	ConfigFile = "mysql.yaml"
)

type BaseConfig struct {
}

type DBConfig struct {
	BaseConfig
	Tag          string `json:"tag"`
	Address      string `json:"address"`
	Database     string `json:"database"`
	Port         int    `json:"port"`
	UserName     string `json:"user_name"`
	Password     string `json:"password"`
	Timeout      int    `json:"timeout"`
	MaxLifeTime  int    `json:"max_life_time"`
	MaxOpenConns int    `json:"max_open_conns"`
	MaxIdleConns int    `json:"max_idle_conns"`
}

func (c DBConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4", c.UserName, c.Password, c.Address, c.Port, c.Database)
}

func GetDBConfig(file string) DBConfig {
	if len(file) != 0 {
		ConfigFile = file
	}
	viper.SetConfigFile(ConfigFile)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			panic(fmt.Errorf("Fatal error ConfigFileNotFoundError "))
		} else {
			panic(fmt.Errorf("Fatal error %s \n", err))
		}
	}
	return DBConfig{
		BaseConfig:   BaseConfig{},
		Address:      viper.GetString("address"),
		Port:         viper.GetInt("port"),
		Timeout:      viper.GetInt("timeout"),
		Database:     viper.GetString("database"),
		UserName:     viper.GetString("user_name"),
		Password:     viper.GetString("password"),
		MaxLifeTime:  viper.GetInt("max_open_conns"),
		MaxOpenConns: viper.GetInt("max_life_time"),
		MaxIdleConns: viper.GetInt("max_idle_conns"),
	}
}
