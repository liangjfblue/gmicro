package article

import (
	"github.com/liangjfblue/gmicro/app/interface/post/api"
	postv1 "github.com/liangjfblue/gmicro/app/service/post/proto/v1"
	v1 "github.com/liangjfblue/gmicro/app/service/user/proto/v1"
	"github.com/liangjfblue/gmicro/library/logger"
)

type Srv struct {
	Logger *logger.Logger

	userSrvClient    v1.UserService
	articleSrvClient postv1.PostArticleSrvService
}

func NewArticleSrv(logger *logger.Logger) *Srv {
	return &Srv{
		Logger:           logger,
		userSrvClient:    api.NewUserSrvClient(),
		articleSrvClient: api.NewPostArticleSrvClient(),
	}
}
