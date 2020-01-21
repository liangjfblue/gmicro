package server

import (
	"strconv"
	"time"

	config2 "github.com/liangjfblue/gmicro/library/config"

	"github.com/liangjfblue/gmicro/library/pkg/tracer"

	"github.com/liangjfblue/gmicro/library/logger"

	"github.com/liangjfblue/gmicro/app/interface/post/config"
	"github.com/liangjfblue/gmicro/app/interface/post/router"

	"github.com/micro/go-micro/web"
)

type Server struct {
	serviceName    string
	serviceVersion string

	Logger *logger.Logger
	Config *config.Config

	Service web.Service
	Router  *router.Router

	Tracer *tracer.Tracer
}

func NewServer(serviceName, serviceVersion string) *Server {
	s := new(Server)

	s.serviceName = serviceName
	s.serviceVersion = serviceVersion

	s.Logger = logger.NewLogger(
		logger.LogDirName("./logs"),
		logger.AllowLogLevel(logger.LevelD),
		logger.FlushInterval(time.Duration(2)*time.Second),
	)

	s.Config = config.NewConfig()

	s.Router = router.NewRouter(s.Logger)

	s.Tracer = tracer.New(s.Logger, config2.Conf.TraceConf.Addr, s.serviceName)

	return s
}

func (s *Server) Init() {
	s.Logger.Init()

	s.Tracer.Init()

	//register := etcdv3.NewRegistry(
	//	registry.Addrs("172.16.7.16:9002", "172.16.7.16:9004", "172.16.7.16:9006"),
	//)

	s.Service = web.NewService(
		web.Name(s.serviceName),
		web.Version(s.serviceVersion),
		web.Address(":"+strconv.Itoa(s.Config.HttpConf.Port)),
		web.RegisterTTL(time.Second*30),
		web.RegisterInterval(time.Second*15),
		//web.Registry(register),
	)

	if err := s.Service.Init(); err != nil {
		s.Logger.Error("web post err: %s", err.Error())
		return
	}

	s.Router.Init()

	s.Service.Handle("/", s.Router.G)
}

func (s *Server) Run() {
	defer func() {
		s.Logger.Info("web post close, clean and close something")
		s.Tracer.Close()
	}()

	s.Logger.Debug("web post server run")
	if err := s.Service.Run(); err != nil {
		s.Logger.Error("web post err: %s", err.Error())
		return
	}
}
