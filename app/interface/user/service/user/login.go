package user

import (
	"context"

	v1 "github.com/liangjfblue/gmicro/app/service/user/proto/v1"

	"github.com/liangjfblue/gmicro/app/interface/user/models"

	"github.com/liangjfblue/gmicro/library/pkg/errno"
)

func (u *UserSrv) Login(ctx context.Context, req *models.LoginRequest) (*models.LoginRespond, error) {
	result, err := u.userSrvClient.Login(ctx, &v1.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		err = errno.ErrUserLogin
		u.Logger.Error("web user Login err: %s", err.Error())
		return nil, err
	}

	resp := &models.LoginRespond{
		Code:  result.Code,
		Token: result.Token,
	}
	return resp, nil
}
