package router

import (
	"net/http"

	"github.com/liangjfblue/gmicro/app/interface/user/controllers"
	"github.com/liangjfblue/gmicro/app/interface/user/service"
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

	u := r.G.Group("/v1/user")
	u.Use()
	{
		u.POST("/register", handles.UserHandle.Register)
		u.POST("/login", handles.UserHandle.Login)

		u.Use(srv.AuthMid.AuthMid())
		{
			u.GET("/info", handles.UserHandle.Info)
		}
	}

	c := r.G.Group("/v1/coin")
	c.Use(srv.AuthMid.AuthMid())
	{
		c.GET("/get", handles.CoinHandle.Get)
		c.POST("/add", handles.CoinHandle.Add)
	}
}
