package coin

import (
	"context"

	v1 "github.com/liangjfblue/gmicro/app/service/user/proto/v1"

	"github.com/liangjfblue/gmicro/app/interface/user/models"

	"github.com/liangjfblue/gmicro/library/pkg/errno"
)

func (u *CoinSrv) Get(ctx context.Context, req *models.CoinGetRequest) (*models.CoinGetRespond, error) {
	result, err := u.userSrvClient.CoinGet(ctx, &v1.CoinGetRequest{
		Uid: req.Uid,
	})
	if err != nil {
		err = errno.ErrCoinGet
		u.Logger.Error("web coin get err: %s", err.Error())
		return nil, err
	}

	resp := &models.CoinGetRespond{
		Value: result.Value,
	}

	return resp, nil
}
