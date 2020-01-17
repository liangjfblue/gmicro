package comment

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/liangjfblue/gmicro/app/interface/comment/models"
	"github.com/liangjfblue/gmicro/library/http/handle"
	"github.com/liangjfblue/gmicro/library/pkg/errno"
)

func (m *CommentHandle) Del(c *gin.Context) {
	var (
		err    error
		result handle.Result
		req    models.DelCommentRequest
	)

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

	resp, err := m.Srv.CommentSrv.Del(c, &req)
	if err != nil {
		m.Logger.Error("web comment del err: %s", err.Error())
		result.Failure(c, err)
		return
	}

	result.Success(c, resp)
}
