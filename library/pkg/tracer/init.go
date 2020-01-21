package tracer

import (
	"io"

	"github.com/liangjfblue/gmicro/library/logger"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
)

type Tracer struct {
	log *logger.Logger
	c   io.Closer

	traceAddress string
	serviceName  string
}

func New(log *logger.Logger, traceAddress, serviceName string) *Tracer {
	return &Tracer{
		log:          log,
		traceAddress: traceAddress,
		serviceName:  serviceName,
	}
}

func (t *Tracer) Init() {
	t.c = t.traceingInit(t.traceAddress, t.serviceName)
	t.log.Info("初始化traceing:%+v", opentracing.GlobalTracer())
}

func (t *Tracer) Close() {
	if t.c != nil {
		_ = t.c.Close()
	}
}

func (t *Tracer) traceingInit(address, servicename string) io.Closer {
	cfg := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans: true,
		},
	}

	jLogger := jaegerlog.StdLogger
	jMetricsFactory := metrics.NullFactory

	//metricsFactory := metrics.NewLocalFactory(0)
	_metrics := jaeger.NewMetrics(jMetricsFactory, nil)

	sender, err := jaeger.NewUDPTransport(address, 0)
	if err != nil {
		t.log.Info("could not initialize jaeger sender: " + err.Error())
		return nil
	}

	repoter := jaeger.NewRemoteReporter(sender, jaeger.ReporterOptions.Metrics(_metrics))

	closer, err := cfg.InitGlobalTracer(
		servicename,
		jaegercfg.Logger(jLogger),
		jaegercfg.Metrics(jMetricsFactory),
		jaegercfg.Reporter(repoter),
	)

	if err != nil {
		t.log.Info("could not initialize jaeger tracer: " + err.Error())
		return nil
	}
	return closer
}
