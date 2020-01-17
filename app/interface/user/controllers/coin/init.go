package coin

import (
	"github.com/liangjfblue/gmicro/app/interface/user/service"
	"github.com/liangjfblue/gmicro/library/logger"
)

type CoinHandle struct {
	Logger *logger.Logger
	Srv    *service.Service
}

func NewCoinHandle(logger *logger.Logger, srv *service.Service) *CoinHandle {
	return &CoinHandle{
		Logger: logger,
		Srv:    srv,
	}
}
