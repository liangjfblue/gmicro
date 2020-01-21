package service

import (
	"context"

	userv1 "github.com/liangjfblue/gmicro/app/service/user/proto/v1"

	"github.com/liangjfblue/gmicro/app/service/comment/model"
	"github.com/liangjfblue/gmicro/library/pkg/errno"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	commentv1 "github.com/liangjfblue/gmicro/app/service/comment/proto/v1"
)

func (c *Service) AddComment(ctx context.Context, in *commentv1.AddCommentRequest, out *commentv1.AddCommentRespond) error {
	var (
		err error
	)

	if ctx.Err() == context.Canceled {
		return errors.Wrap(status.New(codes.Canceled, "Client cancelled, abandoning").Err(), "service comment")
	}

	defer func() {
		if err != nil {
			err = errno.ErrCommentAdd
		}
	}()

	comment := model.TBComment{
		ArticleId: uint(in.ArticleId),
		Comment:   in.Comment,
		FromId:    in.FromId,
		ToId:      in.ToId,
	}
	if err = comment.Create(); err != nil {
		c.Logger.Error("service comment: %s", err.Error())
		return errors.Wrap(err, " service comment")
	}

	//add coin
	resp2, err := c.userSrvClient.CoinAdd(ctx, &userv1.CoinAddRequest{
		Uid:   in.Uid,
		Value: int32(c.Config.CommentConf.AddCoin),
	})
	if err != nil {
		c.Logger.Error("service comment: %s", err.Error())
		return errors.Wrap(err, " service comment")
	}

	out.Code = resp2.Code

	return nil
}

func (c *Service) DelComment(ctx context.Context, in *commentv1.DelCommentRequest, out *commentv1.DelCommentRespond) error {
	var (
		err error
	)

	if ctx.Err() == context.Canceled {
		return errors.Wrap(status.New(codes.Canceled, "Client cancelled, abandoning").Err(), "service comment")
	}

	defer func() {
		if err != nil {
			err = errno.ErrCommentDel
		}
	}()

	if err = model.DeleteComment(uint(in.CommentId)); err != nil {
		c.Logger.Error("service comment: %s", err.Error())
		return errors.Wrap(err, " service comment")
	}

	out.Code = errno.Success.Code

	return nil
}

func (c *Service) ListComment(ctx context.Context, in *commentv1.ListCommentRequest, out *commentv1.ListCommentRespond) error {
	var (
		err error
	)

	if ctx.Err() == context.Canceled {
		return errors.Wrap(status.New(codes.Canceled, "Client cancelled, abandoning").Err(), "service comment")
	}

	defer func() {
		if err != nil {
			err = errno.ErrCommentList
		}
	}()

	comments, count, err := model.ListComments(uint(in.ArticleId), in.Size*(in.Page-1), in.Size)
	if err != nil {
		c.Logger.Error("service comment: %s", err.Error())
		return errors.Wrap(err, " service comment")
	}

	out.Count = int32(count)
	out.Lists = make(map[int32]*commentv1.Comment, out.Count)

	for _, comment := range comments {
		out.Lists[int32(comment.ID)] = &commentv1.Comment{
			Id:      int32(comment.ID),
			Comment: comment.Comment,
			Time:    comment.CreatedAt.Format("2016-01-02 13:04:05"),
		}
	}

	return nil
}
