package user

import (
	"github.com/gin-gonic/gin"
	"github.com/liangjfblue/gmicro/app/interface/user/models"
	"github.com/liangjfblue/gmicro/library/http/handle"
	"github.com/liangjfblue/gmicro/library/pkg/errno"
)

func (u *UserHandle) Register(c *gin.Context) {
	var (
		err    error
		result handle.Result
		req    models.RegisterRequest
	)

	if err = c.BindJSON(&req); err != nil {
		result.Failure(c, errno.ErrBind)
		return
	}

	resp, err := u.Srv.UserSrv.Register(c, &req)
	if err != nil {
		u.Logger.Error("web user Register err: %s", err.Error())
		result.Failure(c, err)
		return
	}

	result.Success(c, resp)
}
