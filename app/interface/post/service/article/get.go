package article

import (
	"context"

	"github.com/liangjfblue/gmicro/app/interface/post/models"
	postv1 "github.com/liangjfblue/gmicro/app/service/post/proto/v1"
	"github.com/liangjfblue/gmicro/library/pkg/errno"
)

func (u *Srv) Get(ctx context.Context, req *models.GetArticleRequest) (*models.GetArticleRespond, error) {
	result, err := u.articleSrvClient.GetArticle(ctx, &postv1.GetArticleRequest{
		Uid:       req.Uid,
		ArticleId: req.ArticleId,
	})
	if err != nil {
		err = errno.ErrArticleGet
		u.Logger.Error("web post get err: %s", err.Error())
		return nil, err
	}

	resp := &models.GetArticleRespond{
		Title:          result.Title,
		Topic:          result.Topic,
		Author:         result.Author,
		IsOriginal:     result.IsOriginal,
		Content:        result.Content,
		CreateTime:     result.CreateTime,
		LastUpdateTime: result.LastUpdateTime,
	}
	return resp, nil
}
