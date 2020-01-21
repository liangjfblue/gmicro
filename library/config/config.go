package config

import (
	"os"
	"strconv"
)

type Config struct {
	JwtConf   JwtConfig
	TraceConf TraceConfig
}

type JwtConfig struct {
	JwtKey  string
	JwtTime int
}

type TraceConfig struct {
	Addr         string
	TraceContext string //tracer
	ReqParam     string //请求参数绑定
}

var (
	Conf Config
)

func init() {
	Conf.JwtConf.JwtKey = os.Getenv("CONF_JWT_JwtKey")
	Conf.JwtConf.JwtTime, _ = strconv.Atoi(os.Getenv("CONF_JWT_JwtTime"))

	if Conf.JwtConf.JwtKey == "" {
		Conf.JwtConf.JwtKey = "k1jr2907u01h01"
	}
	if Conf.JwtConf.JwtTime <= 0 {
		Conf.JwtConf.JwtTime = 3600
	}

	Conf.TraceConf.Addr = os.Getenv("CONFIGOR_TRACE_ADDRESS")
	Conf.TraceConf.TraceContext = os.Getenv("CONFIGOR_TRACE_TRACECONTEXT")
	Conf.TraceConf.ReqParam = os.Getenv("CONFIGOR_TRACE_REQPARAM")
	if Conf.TraceConf.Addr == "" {
		Conf.TraceConf.Addr = "127.0.0.1:6831"
	}
	if Conf.TraceConf.TraceContext == "" {
		Conf.TraceConf.TraceContext = "trace_ctx"
	}
	if Conf.TraceConf.ReqParam == "" {
		Conf.TraceConf.ReqParam = "req_param"
	}
}
