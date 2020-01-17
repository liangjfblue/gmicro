package comment

import (
	"errors"

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

	_, ok := c.Get("uid")
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

	resp, err := m.Srv.CommentSrv.Add(c, &req)
	if err != nil {
		m.Logger.Error("web comment add err: %s", err.Error())
		result.Failure(c, err)
		return
	}

	result.Success(c, resp)
}
