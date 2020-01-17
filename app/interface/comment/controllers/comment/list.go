package comment

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/liangjfblue/gmicro/app/interface/comment/models"
	"github.com/liangjfblue/gmicro/library/http/handle"
	"github.com/liangjfblue/gmicro/library/pkg/errno"
)

func (m *CommentHandle) List(c *gin.Context) {
	var (
		err    error
		result handle.Result
		req    models.ListCommentRequest
	)

	uid, ok := c.Get("uid")
	if !ok {
		m.Logger.Error("web comment list err: token no uid")
		result.Failure(c, errors.New("web comment list err: token no uid"))
		return
	}

	if err = c.BindJSON(&req); err != nil {
		m.Logger.Error("web comment list err: %s", err.Error())
		result.Failure(c, errno.ErrBind)
		return
	}

	req.Uid = uid.(string)

	if req.Uid == "" || req.ArticleId <= 0 || req.Page <= 0 || req.Size <= 0 {
		m.Logger.Error("web comment list err: params empty")
		result.Failure(c, errno.ErrParams)
		return
	}

	resp, err := m.Srv.CommentSrv.List(c, &req)
	if err != nil {
		m.Logger.Error("web comment list err: %s", err.Error())
		result.Failure(c, err)
		return
	}

	result.Success(c, resp)
}
