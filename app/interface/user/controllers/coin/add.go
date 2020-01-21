package coin

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/liangjfblue/gmicro/app/interface/user/models"
	"github.com/liangjfblue/gmicro/library/config"
	"github.com/liangjfblue/gmicro/library/http/handle"
	"github.com/liangjfblue/gmicro/library/pkg/errno"
	"github.com/liangjfblue/gmicro/library/pkg/tracer"
)

func (u *CoinHandle) Add(c *gin.Context) {
	var (
		err    error
		result handle.Result
		req    models.CoinAddRequest
	)

	//tracer
	cc, exist := c.Get(config.Conf.TraceConf.TraceContext)
	if !exist {
		u.Logger.Error("no TraceContext")
		result.Failure(c, errno.ErrTraceNoContext)
		return
	}
	ctx := cc.(context.Context)
	ctx, span, err := tracer.TraceIntoContext(ctx, "WebCoinAdd")
	if err != nil {
		u.Logger.Error("web coin err: %s", err.Error())
		result.Failure(c, errno.ErrTraceIntoContext)
		return
	}
	defer span.Finish()

	if err = c.BindJSON(&req); err != nil {
		result.Failure(c, errno.ErrBind)
		return
	}

	resp, err := u.Srv.CoinSrv.Add(c, &req)
	if err != nil {
		u.Logger.Error("web coin add err: %s", err.Error())
		result.Failure(c, err)
		return
	}

	result.Success(c, resp)
}
