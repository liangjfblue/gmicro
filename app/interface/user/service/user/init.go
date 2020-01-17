package user

import (
	"github.com/liangjfblue/gmicro/app/interface/user/api"
	userv1 "github.com/liangjfblue/gmicro/app/service/user/proto/v1"

	"github.com/liangjfblue/gmicro/library/logger"
)

type UserSrv struct {
	Logger *logger.Logger

	userSrvClient userv1.UserService
}

func NewUserSrv(logger *logger.Logger) *UserSrv {
	return &UserSrv{
		Logger:        logger,
		userSrvClient: api.NewUserSrvClient(),
	}
}
