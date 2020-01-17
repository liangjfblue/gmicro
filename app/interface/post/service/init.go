package service

import (
	"github.com/liangjfblue/gmicro/app/interface/post/service/article"
	"github.com/liangjfblue/gmicro/library/http/middleware/auth"
	"github.com/liangjfblue/gmicro/library/logger"
)

type Service struct {
	AuthMid    *auth.Auth
	ArticleSrv *article.Srv
}

func NewService(logger *logger.Logger) *Service {
	return &Service{
		AuthMid:    auth.New(logger),
		ArticleSrv: article.NewArticleSrv(logger),
	}
}
