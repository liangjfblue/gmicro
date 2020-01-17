package auth

import (
	"time"

	"github.com/gin-gonic/gin"
	identifyv1 "github.com/liangjfblue/gmicro/app/service/identify/proto/v1"
	"github.com/liangjfblue/gmicro/library/common"
	"github.com/liangjfblue/gmicro/library/http/handle"
	"github.com/liangjfblue/gmicro/library/logger"
	"github.com/micro/go-micro/client"
)

type Auth struct {
	Logger *logger.Logger

	identifySrvClient identifyv1.IdentifySrvService
}

func New(logger *logger.Logger) *Auth {
	a := new(Auth)
	a.Logger = logger

	a.identifySrvClient = identifyv1.NewIdentifySrvService(common.IdentifySrvName, client.NewClient(
		client.Retries(0),
		client.DialTimeout(time.Minute*2),
	))

	return a
}

func (m *Auth) AuthMid() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			err    error
			result handle.Result
		)

		token := c.Request.Header.Get("Authorization")

		req := identifyv1.AuthRequest{
			Token: token,
		}

		resp, err := m.identifySrvClient.AuthMid(c, &req)
		if err != nil {
			m.Logger.Error(err.Error())
			result.Failure(c, err)
			c.Abort()
		}

		c.Set("uid", resp.UID)

		c.Next()
	}
}
