package config

import (
	"os"
	"strconv"
)

type Config struct {
	JwtConf JwtConfig
}

type JwtConfig struct {
	JwtKey  string
	JwtTime int
}

var (
	Conf *Config
)

func init() {
	Conf = new(Config)

	Conf.JwtConf.JwtKey = os.Getenv("CONF_APPID")
	Conf.JwtConf.JwtTime, _ = strconv.Atoi(os.Getenv("CONF_VERSION"))

	if Conf.JwtConf.JwtKey == "" {
		Conf.JwtConf.JwtKey = "k1jr2907u01h01"
	}
	if Conf.JwtConf.JwtTime <= 0 {
		Conf.JwtConf.JwtTime = 3600
	}
}
