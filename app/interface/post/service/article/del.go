package article

import (
	"context"

	"github.com/liangjfblue/gmicro/app/interface/post/models"
	postv1 "github.com/liangjfblue/gmicro/app/service/post/proto/v1"
	"github.com/liangjfblue/gmicro/library/pkg/errno"
)

func (u *Srv) Del(ctx context.Context, req *models.DelArticleRequest) (*models.DelArticleRespond, error) {
	result, err := u.articleSrvClient.DelArticle(ctx, &postv1.DelArticleRequest{
		Uid:       req.Uid,
		ArticleId: req.ArticleId,
	})
	if err != nil {
		err = errno.ErrArticleDel
		u.Logger.Error("web post del err: %s", err.Error())
		return nil, err
	}

	resp := &models.DelArticleRespond{
		Code: result.Code,
	}
	return resp, nil
}
