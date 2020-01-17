package service

import (
	"context"

	v1 "github.com/liangjfblue/gmicro/app/service/user/proto/v1"

	"github.com/liangjfblue/gmicro/app/service/user/model"
	"github.com/liangjfblue/gmicro/library/pkg/errno"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) CoinAdd(ctx context.Context, in *v1.CoinAddRequest, out *v1.CoinAddResponse) error {
	var (
		err   error
		money *model.TBMoney
	)

	if ctx.Err() == context.Canceled {
		return errors.Wrap(status.New(codes.Canceled, "Client cancelled, abandoning").Err(), "service coin")
	}

	money, err = model.GetMoney(&model.TBMoney{Uid: in.Uid})
	if err != nil {
		s.Logger.Error("service coin: %s", err.Error())
		return errors.Wrap(err, "service coin")
	}

	money.Coin += in.Value
	if money.Coin < 0 {
		money.Coin = 0
	}

	if err = money.Update(); err != nil {
		s.Logger.Error("service coin: %s", err.Error())
		return errors.Wrap(err, "service coin")
	}

	out.Code = errno.Success.Code

	return nil

}
func (s *Service) CoinGet(ctx context.Context, in *v1.CoinGetRequest, out *v1.CoinGetResponse) error {
	var (
		err   error
		money *model.TBMoney
	)

	if ctx.Err() == context.Canceled {
		return errors.Wrap(status.New(codes.Canceled, "Client cancelled, abandoning").Err(), "service coin")
	}

	money, err = model.GetMoney(&model.TBMoney{Uid: in.Uid})
	if err != nil {
		s.Logger.Error("service coin: %s", err.Error())
		return errors.Wrap(err, "service coin")
	}

	out.Value = money.Coin

	return nil
}
