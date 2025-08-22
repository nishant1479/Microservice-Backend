package interceptors

import (
	"context"
	"time"

	"github.com/nishant1479/Microservice-Backend/config"
	"github.com/nishant1479/Microservice-Backend/pkg/logger"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var (
	totalRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name : "products_service_requests_total",
		Help: "The total number of incoming gRPC messages",
	})
)

type InterceptorManger struct{
	logger	logger.Logger
	cfg		*config.Config
}

func NewInterceptorManger(logger logger.Logger, cfg *config.Config) *InterceptorManger {
	return &InterceptorManger{logger: logger,cfg: cfg}
}

func (im *InterceptorManger) Logger (
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
)(resp interface{}, err error) {
	totalRequests.Inc()
	start := time.Now()
	md, _ := metadata.FromIncomingContext(ctx)
	reply, err := handler(ctx,req)
	im.logger.Info("Method: %s, Time: %v, Metadata: %v,Err: %v", info.FullMethod, time.Since(start),md,err)

	return reply,err
}