package comment

import (
	"context"

	"github.com/liangjfblue/gmicro/app/interface/comment/models"
	commentv1 "github.com/liangjfblue/gmicro/app/service/comment/proto/v1"
	"github.com/liangjfblue/gmicro/library/pkg/errno"
)

func (u *Srv) Add(ctx context.Context, req *models.AddCommentRequest) (*models.AddCommentRespond, error) {
	result, err := u.commentSrvClient.AddComment(ctx, &commentv1.AddCommentRequest{
		Uid:       req.Uid,
		ArticleId: req.ArticleId,
		Comment:   req.Comment,
		FromId:    req.FromId,
		ToId:      req.ToId,
	})
	if err != nil {
		err = errno.ErrCommentAdd
		u.Logger.Error("web comment add err: %s", err.Error())
		return nil, err
	}

	resp := &models.AddCommentRespond{
		Code: result.Code,
	}
	return resp, nil
}
