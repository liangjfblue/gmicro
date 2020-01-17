package router

import (
	"net/http"

	"github.com/liangjfblue/gmicro/app/interface/comment/controllers"
	"github.com/liangjfblue/gmicro/app/interface/comment/service"
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

	c := r.G.Group("/v1/comment")
	c.Use(srv.AuthMid.AuthMid())
	{
		c.POST("/add", handles.CommentHandle.Add)
		c.DELETE("/:cid", handles.CommentHandle.Del)
		c.GET("/list", handles.CommentHandle.List)
	}
}
