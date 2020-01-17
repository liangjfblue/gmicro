package comment

import (
	"github.com/liangjfblue/gmicro/app/interface/comment/service"
	"github.com/liangjfblue/gmicro/library/logger"
)

type CommentHandle struct {
	Logger *logger.Logger
	Srv    *service.Service
}

func NewCommentHandle(logger *logger.Logger, srv *service.Service) *CommentHandle {
	return &CommentHandle{
		Logger: logger,
		Srv:    srv,
	}
}
