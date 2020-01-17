package coin

import (
	"github.com/liangjfblue/gmicro/app/interface/user/api"
	v1 "github.com/liangjfblue/gmicro/app/service/user/proto/v1"
	"github.com/liangjfblue/gmicro/library/logger"
)

type CoinSrv struct {
	Logger *logger.Logger

	userSrvClient v1.UserService
}

func NewCoinSrv(logger *logger.Logger) *CoinSrv {
	return &CoinSrv{
		Logger:        logger,
		userSrvClient: api.NewUserSrvClient(),
	}
}
