package service

import (
	"github.com/liangjfblue/gmicro/app/interface/comment/service/comment"
	"github.com/liangjfblue/gmicro/library/http/middleware/auth"
	"github.com/liangjfblue/gmicro/library/logger"
)

type Service struct {
	AuthMid    *auth.Auth
	CommentSrv *comment.Srv
}

func NewService(logger *logger.Logger) *Service {
	return &Service{
		AuthMid:    auth.New(logger),
		CommentSrv: comment.NewCommentSrv(logger),
	}
}
