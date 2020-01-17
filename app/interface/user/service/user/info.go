package user

import (
	"context"

	v1 "github.com/liangjfblue/gmicro/app/service/user/proto/v1"

	"github.com/liangjfblue/gmicro/app/interface/user/models"

	"github.com/liangjfblue/gmicro/library/pkg/errno"
)

func (u *UserSrv) Info(ctx context.Context, req *models.InfoRequest) (*models.InfoRespond, error) {
	result, err := u.userSrvClient.Info(ctx, &v1.InfoRequest{
		Uid: req.Uid,
	})
	if err != nil {
		err = errno.ErrUserInfo
		u.Logger.Error("web user Info err: %s", err.Error())
		return nil, err
	}

	resp := &models.InfoRespond{
		Username: result.Username,
		Age:      result.Age,
		Addr:     result.Addr,
	}
	return resp, nil
}
