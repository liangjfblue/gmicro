package article

import (
	"github.com/liangjfblue/gmicro/app/interface/post/service"
	"github.com/liangjfblue/gmicro/library/logger"
)

type ArticleHandle struct {
	Logger *logger.Logger
	Srv    *service.Service
}

func NewArticleHandle(logger *logger.Logger, srv *service.Service) *ArticleHandle {
	return &ArticleHandle{
		Logger: logger,
		Srv:    srv,
	}
}
