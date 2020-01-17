package comment

import (
	"context"

	"github.com/liangjfblue/gmicro/app/interface/comment/models"
	commentv1 "github.com/liangjfblue/gmicro/app/service/comment/proto/v1"
	"github.com/liangjfblue/gmicro/library/pkg/errno"
)

func (u *Srv) Del(ctx context.Context, req *models.DelCommentRequest) (*models.DelCommentRespond, error) {
	result, err := u.commentSrvClient.DelComment(ctx, &commentv1.DelCommentRequest{
		Uid:       req.Uid,
		ArticleId: req.ArticleId,
		CommentId: req.CommentId,
	})
	if err != nil {
		err = errno.ErrCommentDel
		u.Logger.Error("web comment del err: %s", err.Error())
		return nil, err
	}

	resp := &models.DelCommentRespond{
		Code: result.Code,
	}
	return resp, nil
}
