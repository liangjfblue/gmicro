package server

import (
	"time"

	"github.com/liangjfblue/gmicro/app/service/comment/configs"
	"github.com/liangjfblue/gmicro/app/service/comment/model"
	v1 "github.com/liangjfblue/gmicro/app/service/comment/proto/v1"

	commentSrv "github.com/liangjfblue/gmicro/app/service/comment/service"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/server"

	"github.com/liangjfblue/gmicro/library/logger"
)

type Server struct {
	serviceName    string
	serviceVersion string

	Logger *logger.Logger
	Config *configs.Config

	service micro.Service
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

	s.Config = configs.NewConfig()

	return s
}

func (s *Server) Init() {
	s.Logger.Init()

	model.Init(s.Config.MysqlConf)

	//registre := etcdv3.NewRegistry(
	//	registry.Addrs("172.16.7.16:9002", "172.16.7.16:9004", "172.16.7.16:9006"),
	//)

	s.service = micro.NewService(
		micro.Name(s.serviceName),
		micro.Version(s.serviceVersion),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*15),
		//	micro.Registry(registre),
	)

	s.service.Init()

	s.initRegisterHandler()
}

func (s *Server) initRegisterHandler() {
	srv := commentSrv.New(s.Logger)

	if err := v1.RegisterPostCommentSrvHandler(s.service.Server(), srv, server.InternalHandler(true)); err != nil {
		s.Logger.Error("service comment err: %s", err.Error())
		return
	}
}

func (s *Server) Run() {
	s.Logger.Debug("service comment server run")
	if err := s.service.Run(); err != nil {
		s.Logger.Error("service comment err: %s", err.Error())
		return
	}
}
