package controllers

import (
	"github.com/liangjfblue/gmicro/app/interface/post/controllers/article"
	"github.com/liangjfblue/gmicro/app/interface/post/service"
	"github.com/liangjfblue/gmicro/library/logger"
)

type Handles struct {
	ArticleHandle *article.ArticleHandle
}

func NewHandles(logger *logger.Logger, srv *service.Service) *Handles {
	return &Handles{
		ArticleHandle: article.NewArticleHandle(logger, srv),
	}
}
