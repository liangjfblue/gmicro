package coin

import (
	"context"
	"errors"

	"github.com/liangjfblue/gmicro/app/interface/user/models"

	"github.com/gin-gonic/gin"
	"github.com/liangjfblue/gmicro/library/config"
	"github.com/liangjfblue/gmicro/library/http/handle"
	"github.com/liangjfblue/gmicro/library/pkg/errno"
	"github.com/liangjfblue/gmicro/library/pkg/tracer"
)

func (u *CoinHandle) Get(c *gin.Context) {
	var (
		err    error
		result handle.Result
		req    models.CoinGetRequest
	)

	//tracer
	cc, exist := c.Get(config.Conf.TraceConf.TraceContext)
	if !exist {
		u.Logger.Error("no TraceContext")
		result.Failure(c, errno.ErrTraceNoContext)
		return
	}
	ctx := cc.(context.Context)
	ctx, span, err := tracer.TraceIntoContext(ctx, "WebCoinGet")
	if err != nil {
		u.Logger.Error("web coin err: %s", err.Error())
		result.Failure(c, errno.ErrTraceIntoContext)
		return
	}
	defer span.Finish()

	uid, ok := c.Get("uid")
	if !ok {
		u.Logger.Error("web coin err: token no uid")
		result.Failure(c, errors.New("web coin err: token no uid"))
		return
	}

	req.Uid = uid.(string)

	resp, err := u.Srv.CoinSrv.Get(c, &req)
	if err != nil {
		u.Logger.Error("web coin err: %s", err.Error())
		result.Failure(c, err)
		return
	}

	result.Success(c, resp)
}
