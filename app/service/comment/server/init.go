package server

import (
	"time"

	config2 "github.com/liangjfblue/gmicro/library/config"

	"github.com/liangjfblue/gmicro/library/pkg/tracer"

	"github.com/opentracing/opentracing-go"

	"github.com/liangjfblue/gmicro/app/service/comment/configs"
	"github.com/liangjfblue/gmicro/app/service/comment/model"
	v1 "github.com/liangjfblue/gmicro/app/service/comment/proto/v1"
	commentSrv "github.com/liangjfblue/gmicro/app/service/comment/service"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/server"
	ot "github.com/micro/go-plugins/wrapper/trace/opentracing"

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
	srv := commentSrv.New(s.Logger, s.Config)

	if err := v1.RegisterPostCommentSrvHandler(s.service.Server(), srv, server.InternalHandler(true)); err != nil {
		s.Logger.Error("service comment err: %s", err.Error())
		return
	}
}

func (s *Server) Run() {
	defer func() {
		s.Logger.Info("srv comment close, clean and close something")
		s.Tracer.Close()
	}()

	s.Logger.Debug("service comment server run")
	if err := s.service.Run(); err != nil {
		s.Logger.Error("service comment err: %s", err.Error())
		return
	}
}
