package service

import (
	"time"

	userv1 "github.com/liangjfblue/gmicro/app/service/user/proto/v1"
	"github.com/liangjfblue/gmicro/library/common"
	"github.com/liangjfblue/gmicro/library/logger"
	"github.com/micro/go-micro/client"
)

type Service struct {
	Logger *logger.Logger

	userSrvClient userv1.UserService
}

func New(logger *logger.Logger) *Service {
	return &Service{
		Logger:        logger,
		userSrvClient: NewUserSrvClient(),
	}
}

func NewUserSrvClient() userv1.UserService {
	c := client.NewClient(
		client.Retries(0),
		client.DialTimeout(time.Minute*2),
	)
	return userv1.NewUserService(common.UserSrvName, c)
}
