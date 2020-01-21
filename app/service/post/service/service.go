package service

import (
	"time"

	"github.com/liangjfblue/gmicro/app/service/post/configs"
	userv1 "github.com/liangjfblue/gmicro/app/service/user/proto/v1"
	"github.com/liangjfblue/gmicro/library/common"
	"github.com/liangjfblue/gmicro/library/logger"
	"github.com/micro/go-micro/client"
)

type Service struct {
	Logger        *logger.Logger
	Config        *configs.Config
	userSrvClient userv1.UserService
}

func New(logger *logger.Logger, config *configs.Config) *Service {
	return &Service{
		Logger:        logger,
		Config:        config,
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
