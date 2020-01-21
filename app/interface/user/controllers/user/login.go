package user

import (
	"context"
	"fmt"

	"github.com/liangjfblue/gmicro/library/config"
	"github.com/liangjfblue/gmicro/library/pkg/errno"
	"github.com/liangjfblue/gmicro/library/pkg/tracer"

	"github.com/gin-gonic/gin"
	"github.com/liangjfblue/gmicro/app/interface/user/models"
	"github.com/liangjfblue/gmicro/library/http/handle"
)

func (u *UserHandle) Login(c *gin.Context) {
	var (
		err    error
		result handle.Result
		req    models.LoginRequest
	)

	//tracer
	cc, exist := c.Get(config.Conf.TraceConf.TraceContext)
	if !exist {
		u.Logger.Error("no TraceContext")
		result.Failure(c, errno.ErrTraceNoContext)
		return
	}
	ctx := cc.(context.Context)
	ctx, span, err := tracer.TraceIntoContext(ctx, "WebUserLogin")
	if err != nil {
		u.Logger.Error("web user err: %s", err.Error())
		result.Failure(c, errno.ErrTraceIntoContext)
		return
	}
	defer span.Finish()

	if err = c.BindJSON(&req); err != nil {
		u.Logger.Error("web user Login err: %s", err.Error())
		result.Failure(c, errno.ErrBind)
		return
	}

	if req.Username == "" || req.Password == "" {
		u.Logger.Error("web user Login err: %s", fmt.Sprintf("params empty: %s %s", req.Username, req.Password))
		result.Failure(c, errno.ErrParams)
		return
	}

	resp, err := u.Srv.UserSrv.Login(ctx, &req)
	if err != nil {
		u.Logger.Error("web user Login err: %s", err.Error())
		result.Failure(c, err)
		return
	}

	result.Success(c, resp)
}
