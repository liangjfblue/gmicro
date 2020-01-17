package service

import (
	"github.com/liangjfblue/gmicro/app/interface/user/service/coin"
	"github.com/liangjfblue/gmicro/app/interface/user/service/user"
	"github.com/liangjfblue/gmicro/library/http/middleware/auth"
	"github.com/liangjfblue/gmicro/library/logger"
)

type Service struct {
	AuthMid *auth.Auth

	UserSrv *user.UserSrv
	CoinSrv *coin.CoinSrv
}

func NewService(logger *logger.Logger) *Service {
	return &Service{
		AuthMid: auth.New(logger),
		UserSrv: user.NewUserSrv(logger),
		CoinSrv: coin.NewCoinSrv(logger),
	}
}
