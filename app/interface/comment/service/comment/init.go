package comment

import (
	"github.com/liangjfblue/gmicro/app/interface/comment/api"
	commentv1 "github.com/liangjfblue/gmicro/app/service/comment/proto/v1"
	userv1 "github.com/liangjfblue/gmicro/app/service/user/proto/v1"
	"github.com/liangjfblue/gmicro/library/logger"
)

type Srv struct {
	Logger *logger.Logger

	userSrvClient    userv1.UserService
	commentSrvClient commentv1.PostCommentSrvService
}

func NewCommentSrv(logger *logger.Logger) *Srv {
	return &Srv{
		Logger:           logger,
		userSrvClient:    api.NewUserSrvClient(),
		commentSrvClient: api.NewCommentSrvClient(),
	}
}
