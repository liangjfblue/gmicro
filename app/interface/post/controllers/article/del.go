package article

import (
	"context"
	"errors"
	"strconv"

	"github.com/liangjfblue/gmicro/library/config"
	"github.com/liangjfblue/gmicro/library/pkg/tracer"

	"github.com/liangjfblue/gmicro/library/pkg/errno"

	"github.com/gin-gonic/gin"
	"github.com/liangjfblue/gmicro/app/interface/post/models"
	"github.com/liangjfblue/gmicro/library/http/handle"
)

func (a *ArticleHandle) Del(c *gin.Context) {
	var (
		err    error
		result handle.Result
		req    models.DelArticleRequest
	)

	//tracer
	cc, exist := c.Get(config.Conf.TraceConf.TraceContext)
	if !exist {
		a.Logger.Error("no TraceContext")
		result.Failure(c, errno.ErrTraceNoContext)
		return
	}
	ctx := cc.(context.Context)
	ctx, span, err := tracer.TraceIntoContext(ctx, "WebArticleDel")
	if err != nil {
		a.Logger.Error("web coin err: %s", err.Error())
		result.Failure(c, errno.ErrTraceIntoContext)
		return
	}
	defer span.Finish()

	uid, ok := c.Get("uid")
	if !ok {
		a.Logger.Error("web post del err: token no uid")
		result.Failure(c, errors.New("web post del err: token no uid"))
		return
	}

	articleId, _ := strconv.Atoi(c.Param("aid"))
	req.ArticleId = int32(articleId)

	if req.ArticleId <= 0 {
		a.Logger.Error("web post get err: wrong ArticleId " + strconv.Itoa(int(req.ArticleId)))
		result.Failure(c, errno.ErrParams)
		return
	}

	req.Uid = uid.(string)

	resp, err := a.Srv.ArticleSrv.Del(ctx, &req)
	if err != nil {
		a.Logger.Error("web post del err: %s", err.Error())
		result.Failure(c, err)
		return
	}

	result.Success(c, resp)
}
