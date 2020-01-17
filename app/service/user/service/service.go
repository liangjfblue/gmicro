package service

import (
	"github.com/liangjfblue/gmicro/library/logger"
	"github.com/liangjfblue/gmicro/library/pkg/token"
)

type Service struct {
	Logger *logger.Logger
	Token  *token.Token
}

func New(logger *logger.Logger, token *token.Token) *Service {
	return &Service{
		Logger: logger,
		Token:  token,
	}
}
