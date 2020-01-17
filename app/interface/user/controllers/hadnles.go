package controllers

import (
	"github.com/liangjfblue/gmicro/app/interface/user/controllers/coin"
	"github.com/liangjfblue/gmicro/app/interface/user/controllers/user"
	"github.com/liangjfblue/gmicro/app/interface/user/service"
	"github.com/liangjfblue/gmicro/library/logger"
)

type Handles struct {
	UserHandle *user.UserHandle
	CoinHandle *coin.CoinHandle
}

func NewHandles(logger *logger.Logger, srv *service.Service) *Handles {
	return &Handles{
		UserHandle: user.NewUserHandle(logger, srv),
		CoinHandle: coin.NewCoinHandle(logger, srv),
	}
}
