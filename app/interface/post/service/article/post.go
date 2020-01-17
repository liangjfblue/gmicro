package article

import (
	"context"

	"github.com/liangjfblue/gmicro/app/interface/post/models"
	postv1 "github.com/liangjfblue/gmicro/app/service/post/proto/v1"
	"github.com/liangjfblue/gmicro/library/pkg/errno"
)

func (u *Srv) Post(ctx context.Context, req *models.PostArticleRequest) (*models.PostArticleRespond, error) {
	result, err := u.articleSrvClient.PostArticle(ctx, &postv1.PostArticleRequest{
		Uid:        req.Uid,
		Title:      req.Title,
		Topic:      req.Topic,
		Author:     req.Author,
		IsOriginal: req.IsOriginal,
		Content:    req.Content,
	})
	if err != nil {
		err = errno.ErrArticlePost
		u.Logger.Error("web post post err: %s", err.Error())
		return nil, err
	}

	resp := &models.PostArticleRespond{
		ArticleId: result.ArticleId,
	}
	return resp, nil
}
