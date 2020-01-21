package comment

import (
	"context"
	"errors"

	"github.com/liangjfblue/gmicro/library/config"
	"github.com/liangjfblue/gmicro/library/pkg/tracer"

	"github.com/gin-gonic/gin"
	"github.com/liangjfblue/gmicro/app/interface/comment/models"
	"github.com/liangjfblue/gmicro/library/http/handle"
	"github.com/liangjfblue/gmicro/library/pkg/errno"
)

func (m *CommentHandle) Add(c *gin.Context) {
	var (
		err    error
		result handle.Result
		req    models.AddCommentRequest
	)

	//tracer
	cc, exist := c.Get(config.Conf.TraceConf.TraceContext)
	if !exist {
		m.Logger.Error("no TraceContext")
		result.Failure(c, errno.ErrTraceNoContext)
		return
	}
	ctx := cc.(context.Context)
	ctx, span, err := tracer.TraceIntoContext(ctx, "WebCommentAdd")
	if err != nil {
		m.Logger.Error("web coin err: %s", err.Error())
		result.Failure(c, errno.ErrTraceIntoContext)
		return
	}
	defer span.Finish()

	uid, ok := c.Get("uid")
	if !ok {
		m.Logger.Error("web comment add err: token no uid")
		result.Failure(c, errors.New("web comment add err: token no uid"))
		return
	}

	if err = c.BindJSON(&req); err != nil {
		m.Logger.Error("web comment add err: %s", err.Error())
		result.Failure(c, errno.ErrBind)
		return
	}

	req.Uid = uid.(string)

	resp, err := m.Srv.CommentSrv.Add(ctx, &req)
	if err != nil {
		m.Logger.Error("web comment add err: %s", err.Error())
		result.Failure(c, err)
		return
	}

	result.Success(c, resp)
}
