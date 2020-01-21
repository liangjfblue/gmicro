package user

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/liangjfblue/gmicro/app/interface/user/models"
	"github.com/liangjfblue/gmicro/library/config"
	"github.com/liangjfblue/gmicro/library/http/handle"
	"github.com/liangjfblue/gmicro/library/pkg/errno"
	"github.com/liangjfblue/gmicro/library/pkg/tracer"
)

func (u *UserHandle) Register(c *gin.Context) {
	var (
		err    error
		result handle.Result
		req    models.RegisterRequest
	)

	//tracer
	cc, exist := c.Get(config.Conf.TraceConf.TraceContext)
	if !exist {
		u.Logger.Error("no TraceContext")
		result.Failure(c, errno.ErrTraceNoContext)
		return
	}
	ctx := cc.(context.Context)
	ctx, span, err := tracer.TraceIntoContext(ctx, "WebUserRegister")
	if err != nil {
		u.Logger.Error("web user err: %s", err.Error())
		result.Failure(c, errno.ErrTraceIntoContext)
		return
	}
	defer span.Finish()

	if err = c.BindJSON(&req); err != nil {
		result.Failure(c, errno.ErrBind)
		return
	}

	resp, err := u.Srv.UserSrv.Register(ctx, &req)
	if err != nil {
		u.Logger.Error("web user Register err: %s", err.Error())
		result.Failure(c, err)
		return
	}

	result.Success(c, resp)
}
