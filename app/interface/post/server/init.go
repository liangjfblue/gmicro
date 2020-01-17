package server

import (
	"time"

	"github.com/liangjfblue/gmicro/library/logger"

	"github.com/liangjfblue/gmicro/app/interface/post/router"

	"github.com/micro/go-micro/web"
)

type Server struct {
	serviceName    string
	serviceVersion string

	Logger *logger.Logger

	Service web.Service
	Router  *router.Router
}

func NewServer(serviceName, serviceVersion string) *Server {
	s := new(Server)

	s.serviceName = serviceName
	s.serviceVersion = serviceVersion

	s.Logger = logger.NewLogger(
		logger.LogDirName("../logs"),
		logger.AllowLogLevel(logger.LevelD),
		logger.FlushInterval(time.Duration(2)*time.Second),
	)

	s.Router = router.NewRouter(s.Logger)

	return s
}

func (s *Server) Init() {
	s.Logger.Init()

	//register := etcdv3.NewRegistry(
	//	registry.Addrs("172.16.7.16:9002", "172.16.7.16:9004", "172.16.7.16:9006"),
	//)

	s.Service = web.NewService(
		web.Name(s.serviceName),
		web.Version(s.serviceVersion),
		web.Address("172.16.7.16:7080"),
		web.RegisterTTL(time.Second*30),
		web.RegisterInterval(time.Second*15),
		//web.Registry(register),
	)

	if err := s.Service.Init(); err != nil {
		s.Logger.Error("web user err: %s", err.Error())
		return
	}

	s.Router.Init()

	s.Service.Handle("/", s.Router.G)
}

func (s *Server) Run() {
	s.Logger.Debug("web user server run")
	if err := s.Service.Run(); err != nil {
		s.Logger.Error("web user err: %s", err.Error())
		return
	}
}
