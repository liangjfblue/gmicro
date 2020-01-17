package comment

import (
	"context"
	"sort"

	"github.com/liangjfblue/gmicro/app/interface/comment/models"
	commentv1 "github.com/liangjfblue/gmicro/app/service/comment/proto/v1"
	"github.com/liangjfblue/gmicro/library/pkg/errno"
)

func (u *Srv) List(ctx context.Context, req *models.ListCommentRequest) (*models.ListCommentRespond, error) {
	result, err := u.commentSrvClient.ListComment(ctx, &commentv1.ListCommentRequest{
		Uid:       req.Uid,
		ArticleId: req.ArticleId,
		Page:      req.Page,
		Size:      req.Size,
	})
	if err != nil {
		err = errno.ErrCommentList
		u.Logger.Error("web comment list err: %s", err.Error())
		return nil, err
	}

	resp := &models.ListCommentRespond{
		Count: result.Count,
	}

	for _, comment := range result.Lists {
		resp.Lists = append(resp.Lists, models.Comment{
			Id:      comment.Id,
			Comment: comment.Comment,
			Time:    comment.Time,
		})
	}

	sort.Slice(resp.Lists, func(i, j int) bool {
		if resp.Lists[i].Id > resp.Lists[j].Id {
			return true
		} else {
			return false
		}
	})

	return resp, nil
}
