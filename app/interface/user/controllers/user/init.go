package user

import (
	"github.com/liangjfblue/gmicro/app/interface/user/service"
	"github.com/liangjfblue/gmicro/library/logger"
)

type UserHandle struct {
	Logger *logger.Logger

	Srv *service.Service
}

func NewUserHandle(logger *logger.Logger, srv *service.Service) *UserHandle {
	return &UserHandle{
		Logger: logger,
		Srv:    srv,
	}
}
