package user

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/liangjfblue/gmicro/app/interface/user/models"
	"github.com/liangjfblue/gmicro/library/http/handle"
)

func (u *UserHandle) Info(c *gin.Context) {
	var (
		err    error
		result handle.Result
		req    models.InfoRequest
	)

	uid, ok := c.Get("uid")
	if !ok {
		u.Logger.Error("web user err: token no uid")
		result.Failure(c, errors.New("web user err: token no uid"))
		return
	}

	req.Uid = uid.(string)

	resp, err := u.Srv.UserSrv.Info(c, &req)
	if err != nil {
		u.Logger.Error("web user err: %s", err.Error())
		result.Failure(c, err)
		return
	}

	result.Success(c, resp)
}
