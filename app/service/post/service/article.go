package service

import (
	"context"

	"github.com/liangjfblue/gmicro/library/pkg/errno"

	"github.com/liangjfblue/gmicro/app/service/post/model"
	v1 "github.com/liangjfblue/gmicro/app/service/post/proto/v1"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c *Service) PostArticle(ctx context.Context, in *v1.PostArticleRequest, out *v1.PostArticleRespond) error {
	var (
		err error
	)

	if ctx.Err() == context.Canceled {
		return errors.Wrap(status.New(codes.Canceled, "Client cancelled, abandoning").Err(), "service coin")
	}

	defer func() {
		if err != nil {
			err = errno.ErrArticlePost
		}
	}()

	tx := model.DB.Begin()

	articleInfo := model.TBArticleInfo{
		Title:      in.Title,
		Topic:      in.Topic,
		Author:     in.Author,
		IsOriginal: int8(in.IsOriginal),
	}
	if err = articleInfo.Create(); err != nil {
		tx.Rollback()
		c.Logger.Error("service post: %s", err.Error())
		return errors.Wrap(err, " service post")
	}

	article := model.TBArticle{
		Uid:           in.Uid,
		Content:       in.Content,
		ArticleInfoId: articleInfo.ID,
	}
	if err = article.Create(); err != nil {
		tx.Rollback()
		c.Logger.Error("service post: %s", err.Error())
		return errors.Wrap(err, " service post")
	}

	tx.Commit()

	out.ArticleId = int32(article.ID)

	return nil
}

func (c *Service) GetArticle(ctx context.Context, in *v1.GetArticleRequest, out *v1.GetArticleRespond) error {
	var (
		err         error
		article     *model.TBArticle
		articleInfo *model.TBArticleInfo
	)

	if ctx.Err() == context.Canceled {
		return errors.Wrap(status.New(codes.Canceled, "Client cancelled, abandoning").Err(), "service article")
	}

	defer func() {
		if err != nil {
			err = errno.ErrArticleGet
		}
	}()

	article, err = model.GetArticle(&model.TBArticle{ID: uint(in.ArticleId), Uid: in.Uid})
	if err != nil {
		c.Logger.Error("service article: %s", err.Error())
		return errors.Wrap(err, " service article")
	}

	articleInfo, err = model.GetArticleInfo(&model.TBArticleInfo{ID: article.ID})
	if err != nil {
		c.Logger.Error("service article: %s", err.Error())
		return errors.Wrap(err, " service article")
	}

	out.Title = articleInfo.Title
	out.Topic = articleInfo.Topic
	out.Author = articleInfo.Author
	out.IsOriginal = int32(articleInfo.IsOriginal)
	out.Content = article.Content
	out.CreateTime = article.CreatedAt.Format("2016-01-02 13:04:05")
	out.LastUpdateTime = article.UpdatedAt.Format("2016-01-02 13:04:05")

	return nil
}

func (c *Service) DelArticle(ctx context.Context, in *v1.DelArticleRequest, out *v1.DelArticleRespond) error {
	var (
		err     error
		article *model.TBArticle
	)

	if ctx.Err() == context.Canceled {
		return errors.Wrap(status.New(codes.Canceled, "Client cancelled, abandoning").Err(), "service article")
	}

	defer func() {
		if err != nil {
			err = errno.ErrArticleDel
		}
	}()

	article, err = model.GetArticle(&model.TBArticle{ID: uint(in.ArticleId), Uid: in.Uid})
	if err != nil {
		c.Logger.Error("service article: %s", err.Error())
		return errors.Wrap(err, " service article")
	}

	ArticleInfoId := article.ArticleInfoId

	tx := model.DB.Begin()

	if err = model.DeleteArticle(uint(in.ArticleId)); err != nil {
		tx.Rollback()
		c.Logger.Error("service article: %s", err.Error())
		return errors.Wrap(err, " service article")
	}

	if err = model.DeleteArticleInfo(ArticleInfoId); err != nil {
		tx.Rollback()
		c.Logger.Error("service article: %s", err.Error())
		return errors.Wrap(err, " service article")
	}

	tx.Commit()

	out.Code = errno.Success.Code

	return nil
}
