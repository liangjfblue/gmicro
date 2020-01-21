package comment

import (
	"context"
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/liangjfblue/gmicro/app/interface/comment/models"
	"github.com/liangjfblue/gmicro/library/config"
	"github.com/liangjfblue/gmicro/library/http/handle"
	"github.com/liangjfblue/gmicro/library/pkg/errno"
	"github.com/liangjfblue/gmicro/library/pkg/tracer"
)

func (m *CommentHandle) Del(c *gin.Context) {
	var (
		err    error
		result handle.Result
		req    models.DelCommentRequest
	)

	//tracer
	cc, exist := c.Get(config.Conf.TraceConf.TraceContext)
	if !exist {
		m.Logger.Error("no TraceContext")
		result.Failure(c, errno.ErrTraceNoContext)
		return
	}
	ctx := cc.(context.Context)
	ctx, span, err := tracer.TraceIntoContext(ctx, "WebCommentDel")
	if err != nil {
		m.Logger.Error("web coin err: %s", err.Error())
		result.Failure(c, errno.ErrTraceIntoContext)
		return
	}
	defer span.Finish()

	uid, ok := c.Get("uid")
	if !ok {
		m.Logger.Error("web comment del err: token no uid")
		result.Failure(c, errors.New("web comment del err: token no uid"))
		return
	}

	req.Uid = uid.(string)

	commentId, _ := strconv.Atoi(c.Param("cid"))
	req.CommentId = int32(commentId)

	if req.CommentId <= 0 {
		m.Logger.Error("web comment get err: wrong CommentId " + strconv.Itoa(int(req.CommentId)))
		result.Failure(c, errno.ErrParams)
		return
	}

	resp, err := m.Srv.CommentSrv.Del(ctx, &req)
	if err != nil {
		m.Logger.Error("web comment del err: %s", err.Error())
		result.Failure(c, err)
		return
	}

	result.Success(c, resp)
}
