package api

import (
	"time"

	userv1 "github.com/liangjfblue/gmicro/app/service/user/proto/v1"

	commentv1 "github.com/liangjfblue/gmicro/app/service/comment/proto/v1"
	"github.com/liangjfblue/gmicro/library/common"
	"github.com/micro/go-micro/client"
)

func NewUserSrvClient() userv1.UserService {
	c := client.NewClient(
		client.Retries(0),
		client.DialTimeout(time.Minute*2),
	)
	return userv1.NewUserService(common.UserSrvName, c)
}

func NewCommentSrvClient() commentv1.PostCommentSrvService {
	c := client.NewClient(
		client.Retries(0),
		client.DialTimeout(time.Minute*2),
	)
	return commentv1.NewPostCommentSrvService(common.CommentSrvName, c)
}
