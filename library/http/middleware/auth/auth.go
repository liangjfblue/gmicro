package auth

import (
	"context"
	"time"

	"github.com/liangjfblue/gmicro/library/pkg/errno"

	"github.com/liangjfblue/gmicro/library/config"

	"github.com/gin-gonic/gin"
	identifyv1 "github.com/liangjfblue/gmicro/app/service/identify/proto/v1"
	"github.com/liangjfblue/gmicro/library/common"
	"github.com/liangjfblue/gmicro/library/http/handle"
	"github.com/liangjfblue/gmicro/library/logger"
	"github.com/liangjfblue/gmicro/library/pkg/tracer"
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

		//tracer
		cc, ok := c.Get(config.Conf.TraceConf.TraceContext)
		if !ok {
			m.Logger.Error("no TraceContext")
			result.Failure(c, errno.ErrTraceNoContext)
			c.Abort()
			return
		}

		ctx := cc.(context.Context)
		ctx, span, err := tracer.TraceIntoContext(ctx, "VerifyToken")
		if err != nil {
			m.Logger.Error(err.Error())
			result.Failure(c, errno.ErrTraceIntoContext)
			c.Abort()
			return
		}
		defer span.Finish()

		//jwt
		token := c.Request.Header.Get("Authorization")

		req := identifyv1.AuthRequest{
			Token: token,
		}

		resp, err := m.identifySrvClient.AuthMid(c, &req)
		if err != nil {
			m.Logger.Error(err.Error())
			result.Failure(c, errno.ErrUserAuthMid)
			c.Abort()
			return
		}

		c.Set("uid", resp.UID)

		c.Next()
	}
}
