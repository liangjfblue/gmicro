package controllers

import (
	"github.com/liangjfblue/gmicro/app/interface/comment/controllers/comment"
	"github.com/liangjfblue/gmicro/app/interface/comment/service"
	"github.com/liangjfblue/gmicro/library/logger"
)

type Handles struct {
	CommentHandle *comment.CommentHandle
}

func NewHandles(logger *logger.Logger, srv *service.Service) *Handles {
	return &Handles{
		CommentHandle: comment.NewCommentHandle(logger, srv),
	}
}
