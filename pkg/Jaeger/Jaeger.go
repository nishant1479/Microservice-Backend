package jaeger

import (
	"io"

	"github.com/nishant1479/Microservice-Backend/config"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
)

func InitJaeger(cfg *config.Config) (opentracing.Tracer, io.Closer,error) {
	jaegerCfgInstance := jaegercfg.Configuration{
		ServiceName: cfg.Jaeger.ServiceName,
		Sampler: &jaegercfg.SamplerConfig{
			Type: jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:				cfg.Jaeger.LogSpans,
			LocalAgentHostPort:		cfg.Jaeger.Host,
		},
	}
	return jaegerCfgInstance.NewTracer(
		jaegercfg.Logger(jaegerlog.StdLogger),
		jaegercfg.Metrics(metrics.NullFactory),
	)
}

// reason why jaegerconfig
/* If you’re already importing another config package (for example, your own 	"github.com/nishant1479/Microservice-Backend/config") without aliasing, you’ll get a naming conflict. In that case, having the alias (jaegercfg) is necessary.

If there are no conflicts, it’s fine to use the default name.*/