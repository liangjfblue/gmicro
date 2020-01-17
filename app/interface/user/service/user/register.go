package user

import (
	"context"

	v1 "github.com/liangjfblue/gmicro/app/service/user/proto/v1"

	"github.com/liangjfblue/gmicro/library/pkg/errno"

	"github.com/liangjfblue/gmicro/app/interface/user/models"
)

func (u *UserSrv) Register(ctx context.Context, req *models.RegisterRequest) (*models.RegisterRespond, error) {
	result, err := u.userSrvClient.Register(ctx, &v1.RegisterRequest{
		Username: req.Username,
		Password: req.Password,
		Age:      req.Age,
		Addr:     req.Addr,
	})
	if err != nil {
		u.Logger.Error("web user Register err: %s", err.Error())
		err = errno.ErrUserRegister
		return nil, err
	}

	resp := &models.RegisterRespond{
		Code: result.Code,
		Uid:  result.Uid,
	}

	return resp, nil
}
