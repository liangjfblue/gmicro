package api

import (
	"time"

	v1 "github.com/liangjfblue/gmicro/app/service/user/proto/v1"

	"github.com/liangjfblue/gmicro/library/common"
	"github.com/micro/go-micro/client"
)

func NewUserSrvClient() v1.UserService {
	c := client.NewClient(
		client.Retries(0),
		client.DialTimeout(time.Minute*2),
	)
	return v1.NewUserService(common.UserSrvName, c)
}
