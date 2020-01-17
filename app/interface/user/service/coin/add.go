package coin

import (
	"context"

	v1 "github.com/liangjfblue/gmicro/app/service/user/proto/v1"

	"github.com/liangjfblue/gmicro/app/interface/user/models"
	"github.com/liangjfblue/gmicro/library/pkg/errno"
)

func (u *CoinSrv) Add(ctx context.Context, req *models.CoinAddRequest) (*models.CoinAddRespond, error) {
	result, err := u.userSrvClient.CoinAdd(ctx, &v1.CoinAddRequest{
		Uid:   req.Uid,
		Value: req.Value,
	})
	if err != nil {
		err = errno.ErrCoinAdd
		u.Logger.Error("web coin add err: %s", err.Error())
		return nil, err
	}

	resp := &models.CoinAddRespond{
		Code: result.Code,
	}
	return resp, nil
}
