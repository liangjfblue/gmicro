package service

import (
	"context"

	"github.com/liangjfblue/gmicro/library/pkg/token"

	"github.com/liangjfblue/gmicro/app/service/identify/model"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	v1 "github.com/liangjfblue/gmicro/app/service/identify/proto/v1"
)

func (u *Service) AuthMid(ctx context.Context, in *v1.AuthRequest, out *v1.AuthResponse) error {
	var (
		err  error
		t    *token.Context
		user *model.TBUser
	)

	if ctx.Err() == context.Canceled {
		u.Logger.Error("service identify: %s", err.Error())
		return errors.Wrap(status.New(codes.Canceled, "Client cancelled, abandoning.").Err(), "service identify")
	}

	if t, err = u.Token.ParseRequest(in.Token); err != nil {
		u.Logger.Error("service identify: %s", err.Error())
		return errors.Wrap(err, "service identify")
	}

	if t.Uid == "" {
		u.Logger.Error("service identify: uid empty")
		return errors.Wrap(errors.New("token uid is empty"), "service identify")
	}

	user, err = model.GetUser(&model.TBUser{Uid: t.Uid})
	if err != nil {
		u.Logger.Error("service identify: %s", err.Error())
		return errors.Wrap(err, "service identify")
	}

	out.UID = user.Uid

	return nil
}
