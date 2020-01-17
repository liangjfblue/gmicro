package service

import (
	"fmt"
	"time"

	"github.com/liangjfblue/gmicro/library/config"

	"github.com/liangjfblue/gmicro/library/pkg/token"

	userv1 "github.com/liangjfblue/gmicro/app/service/user/proto/v1"
	"github.com/liangjfblue/gmicro/library/common"
	"github.com/liangjfblue/gmicro/library/logger"
	"github.com/micro/go-micro/client"
)

type Service struct {
	Logger        *logger.Logger
	Token         *token.Token
	userSrvClient userv1.UserService
}

func New(logger *logger.Logger) *Service {
	s := new(Service)
	s.Logger = logger
	s.Token = token.New(config.Conf.JwtConf.JwtKey, config.Conf.JwtConf.JwtTime)

	fmt.Println(config.Conf.JwtConf)

	s.userSrvClient = newUserSrvClient()
	return s
}

func newUserSrvClient() userv1.UserService {
	c := client.NewClient(
		client.Retries(0),
		client.DialTimeout(time.Minute*2),
	)
	return userv1.NewUserService(common.UserSrvName, c)
}
