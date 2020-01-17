package router

import (
	"net/http"

	"github.com/liangjfblue/gmicro/app/interface/post/controllers"
	"github.com/liangjfblue/gmicro/app/interface/post/service"
	"github.com/liangjfblue/gmicro/library/logger"

	"github.com/gin-gonic/gin"
)

type Router struct {
	G      *gin.Engine
	Logger *logger.Logger
}

func NewRouter(logger *logger.Logger) *Router {
	return &Router{
		G:      gin.Default(),
		Logger: logger,
	}
}

func (r *Router) Init() {
	r.G.Use(gin.Recovery())
	r.G.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "The incorrect API route")
	})

	r.initRouter()
}

func (r *Router) initRouter() {
	srv := service.NewService(r.Logger)

	handles := controllers.NewHandles(r.Logger, srv)

	u := r.G.Group("/v1/post/article")
	u.Use(srv.AuthMid.AuthMid())
	{
		u.POST("/post", handles.ArticleHandle.Post)
		u.GET("/:aid", handles.ArticleHandle.Get)
		u.DELETE("/:aid", handles.ArticleHandle.Del)
	}
}
