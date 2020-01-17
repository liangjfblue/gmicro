package service

import (
	"context"
	"time"

	v1 "github.com/liangjfblue/gmicro/app/service/user/proto/v1"

	"github.com/liangjfblue/gmicro/library/pkg/auth"
	"github.com/liangjfblue/gmicro/library/pkg/token"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/liangjfblue/gmicro/library/pkg/uuid"

	"github.com/jinzhu/gorm"
	"github.com/liangjfblue/gmicro/app/service/user/model"
	"github.com/pkg/errors"

	"github.com/liangjfblue/gmicro/library/pkg/errno"
)

func (s *Service) Register(ctx context.Context, in *v1.RegisterRequest, out *v1.RegisterRespond) error {
	if ctx.Err() == context.Canceled {
		return ctx.Err()
	}

	if _, err := model.GetUser(&model.TBUser{Username: in.Username}); err != nil && !gorm.IsRecordNotFoundError(err) {
		s.Logger.Error("service user: %s", err.Error())
		return errors.Wrap(err, " service user")
	}

	user := model.TBUser{
		Uid:         uuid.UUID(),
		Username:    in.Username,
		Password:    in.Password,
		Age:         in.Age,
		Address:     in.Addr,
		IsAvailable: 1,
		LastLogin:   time.Now(),
	}

	if err := user.Validate(); err != nil {
		s.Logger.Error("service user: %s", err.Error())
		return errors.Wrap(err, " service user")
	}

	if err := user.Encrypt(); err != nil {
		s.Logger.Error("service user: %s", err.Error())
		return errors.Wrap(err, " service user")
	}

	if err := user.Create(); err != nil {
		s.Logger.Error("service user: %s", err.Error())
		return errors.Wrap(err, " service user")
	}

	//coin new record
	money := model.TBMoney{
		Uid:  user.Uid,
		Coin: 0,
	}

	if err := money.Create(); err != nil {
		s.Logger.Error("service user: %s", err.Error())
		return errors.Wrap(err, " service user")
	}

	//update user coin id
	user.MId = money.ID
	if err := user.Update(); err != nil {
		s.Logger.Error("service user: %s", err.Error())
		return errors.Wrap(err, " service user")
	}

	out.Code = errno.Success.Code
	out.Uid = user.Uid

	return nil
}

func (s *Service) Login(ctx context.Context, in *v1.LoginRequest, out *v1.LoginRespond) error {
	if ctx.Err() == context.Canceled {
		return errors.Wrap(status.New(codes.Canceled, "Client cancelled, abandoning.").Err(), "service account")
	}

	var (
		err      error
		user     *model.TBUser
		tokenStr string
	)

	user, err = model.GetUser(&model.TBUser{Username: in.Username})
	if err != nil {
		s.Logger.Error("service user: %s", err.Error())
		return errors.Wrap(err, "service account")
	}

	if err = auth.Compare(user.Password, in.Password); err != nil {
		s.Logger.Error("service user: %s", err.Error())
		return errors.Wrap(err, "service account")
	}

	if user.IsAvailable != 1 {
		s.Logger.Error("account unavailable")
		return errors.Wrap(errors.New("account unavailable"), "service account")
	}

	user.LastLogin = time.Now()
	if err = user.Update(); err != nil {
		s.Logger.Error("service user: %s", err.Error())
		return errors.Wrap(err, "service account")
	}

	tokenStr, err = s.Token.SignToken(token.Context{Uid: user.Uid})
	if err != nil {
		s.Logger.Error("service user: %s", err.Error())
		return errors.Wrap(err, "service account")
	}

	out.Code = errno.Success.Code
	out.Token = tokenStr

	return nil
}

func (s *Service) Info(ctx context.Context, in *v1.InfoRequest, out *v1.InfoRespond) error {
	var (
		err  error
		user *model.TBUser
	)

	if ctx.Err() == context.Canceled {
		return errors.Wrap(status.New(codes.Canceled, "Client cancelled, abandoning").Err(), "service user")
	}

	user, err = model.GetUser(&model.TBUser{Uid: in.Uid})
	if err != nil {
		s.Logger.Error("service user: %s", err.Error())
		return errors.Wrap(err, "service user")
	}

	out.Username = user.Username
	out.Age = user.Age
	out.Addr = user.Address

	return nil
}
