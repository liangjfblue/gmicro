package coin

import (
	"errors"

	"github.com/liangjfblue/gmicro/app/interface/user/models"

	"github.com/gin-gonic/gin"
	"github.com/liangjfblue/gmicro/library/http/handle"
)

func (u *CoinHandle) Get(c *gin.Context) {
	var (
		err    error
		result handle.Result
		req    models.CoinGetRequest
	)

	uid, ok := c.Get("uid")
	if !ok {
		u.Logger.Error("web coin err: token no uid")
		result.Failure(c, errors.New("web coin err: token no uid"))
		return
	}

	req.Uid = uid.(string)

	resp, err := u.Srv.CoinSrv.Get(c, &req)
	if err != nil {
		u.Logger.Error("web coin err: %s", err.Error())
		result.Failure(c, err)
		return
	}

	result.Success(c, resp)
}
