package article

import (
	"context"
	"errors"
	"strconv"

	"github.com/liangjfblue/gmicro/library/config"

	"github.com/liangjfblue/gmicro/library/pkg/errno"
	"github.com/liangjfblue/gmicro/library/pkg/tracer"

	"github.com/gin-gonic/gin"
	"github.com/liangjfblue/gmicro/app/interface/post/models"
	"github.com/liangjfblue/gmicro/library/http/handle"
)

func (a *ArticleHandle) Get(c *gin.Context) {
	var (
		err    error
		result handle.Result
		req    models.GetArticleRequest
	)

	//tracer
	cc, exist := c.Get(config.Conf.TraceConf.TraceContext)
	if !exist {
		a.Logger.Error("no TraceContext")
		result.Failure(c, errno.ErrTraceNoContext)
		return
	}
	ctx := cc.(context.Context)
	ctx, span, err := tracer.TraceIntoContext(ctx, "WebArticleGet")
	if err != nil {
		a.Logger.Error("web coin err: %s", err.Error())
		result.Failure(c, errno.ErrTraceIntoContext)
		return
	}
	defer span.Finish()

	uid, ok := c.Get("uid")
	if !ok {
		a.Logger.Error("web post get err: token no uid")
		result.Failure(c, errors.New("web post get err: token no uid"))
		return
	}

	req.Uid = uid.(string)

	articleId, _ := strconv.Atoi(c.Param("aid"))
	req.ArticleId = int32(articleId)

	if req.ArticleId <= 0 {
		a.Logger.Error("web post get err: wrong ArticleId " + strconv.Itoa(int(req.ArticleId)))
		result.Failure(c, errno.ErrParams)
		return
	}

	resp, err := a.Srv.ArticleSrv.Get(ctx, &req)
	if err != nil {
		a.Logger.Error("web post get err: %s", err.Error())
		result.Failure(c, err)
		return
	}

	result.Success(c, resp)

}
