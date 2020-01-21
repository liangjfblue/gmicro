package article

import (
	"context"
	"errors"

	"github.com/liangjfblue/gmicro/library/config"
	"github.com/liangjfblue/gmicro/library/pkg/tracer"

	"github.com/gin-gonic/gin"
	"github.com/liangjfblue/gmicro/app/interface/post/models"
	"github.com/liangjfblue/gmicro/library/http/handle"
	"github.com/liangjfblue/gmicro/library/pkg/errno"
)

func (a *ArticleHandle) Post(c *gin.Context) {
	var (
		err    error
		result handle.Result
		req    models.PostArticleRequest
	)

	//tracer
	cc, exist := c.Get(config.Conf.TraceConf.TraceContext)
	if !exist {
		a.Logger.Error("no TraceContext")
		result.Failure(c, errno.ErrTraceNoContext)
		return
	}
	ctx := cc.(context.Context)
	ctx, span, err := tracer.TraceIntoContext(ctx, "WebArticlePost")
	if err != nil {
		a.Logger.Error("web coin err: %s", err.Error())
		result.Failure(c, errno.ErrTraceIntoContext)
		return
	}
	defer span.Finish()

	uid, ok := c.Get("uid")
	if !ok {
		a.Logger.Error("web post post err: token no uid")
		result.Failure(c, errors.New("web post post err: token no uid"))
		return
	}

	if err = c.BindJSON(&req); err != nil {
		a.Logger.Error("web post post err: %s", err.Error())
		result.Failure(c, errno.ErrBind)
		return
	}

	req.Uid = uid.(string)

	resp, err := a.Srv.ArticleSrv.Post(ctx, &req)
	if err != nil {
		a.Logger.Error("web post post err: %s", err.Error())
		result.Failure(c, err)
		return
	}

	result.Success(c, resp)
}
