package server

import (
	"time"

	ot "github.com/micro/go-plugins/wrapper/trace/opentracing"
	"github.com/opentracing/opentracing-go"

	userv1 "github.com/liangjfblue/gmicro/app/service/user/proto/v1"
	"github.com/liangjfblue/gmicro/library/config"
	config2 "github.com/liangjfblue/gmicro/library/config"

	"github.com/liangjfblue/gmicro/app/service/user/model"

	"github.com/liangjfblue/gmicro/app/service/user/configs"

	"github.com/liangjfblue/gmicro/library/pkg/token"
	"github.com/liangjfblue/gmicro/library/pkg/tracer"

	service2 "github.com/liangjfblue/gmicro/app/service/user/service"
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

	s.Config = configs.NewConfig()

	s.Tracer = tracer.New(s.Logger, config2.Conf.TraceConf.Addr, s.serviceName)

	return s
}

func (s *Server) Init() {
	s.Logger.Init()

	model.Init(s.Config.MysqlConf)

	s.Tracer.Init()

	//registre := etcdv3.NewRegistry(
	//	registry.Addrs("172.16.7.16:9002", "172.16.7.16:9004", "172.16.7.16:9006"),
	//)

	s.service = micro.NewService(
		micro.Name(s.serviceName),
		micro.Version(s.serviceVersion),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*15),
		micro.WrapClient(ot.NewClientWrapper(opentracing.GlobalTracer())),
		micro.WrapHandler(ot.NewHandlerWrapper(opentracing.GlobalTracer())),
		//	micro.Registry(registre),
	)

	s.service.Init()

	s.initRegisterHandler()
}

func (s *Server) initRegisterHandler() {
	jwt := token.New(config.Conf.JwtConf.JwtKey, config.Conf.JwtConf.JwtTime)

	userSrv := &service2.Service{
		Logger: s.Logger,
		Token:  jwt,
	}
	if err := userv1.RegisterUserHandler(s.service.Server(), userSrv, server.InternalHandler(true)); err != nil {
		s.Logger.Error("service user err: %s", err.Error())
		return
	}
}

func (s *Server) Run() {
	defer func() {
		s.Logger.Info("srv user close, clean and close something")
		s.Tracer.Close()
	}()

	s.Logger.Debug("service user server run")
	if err := s.service.Run(); err != nil {
		s.Logger.Error("service user err: %s", err.Error())
		return
	}
}
