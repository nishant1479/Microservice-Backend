package server

import (
	"context"
	"net"
	"strings"
	"time"



	
	// gRPC middleware
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"


	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/nishant1479/Microservice-Backend/config"
	"github.com/nishant1479/Microservice-Backend/pkg/logger"
	"github.com/opentracing/opentracing-go"
	"github.com/redis/go-redis"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/protobuf/internal/errors"
)

type server struct {
	log		logger.Logger
	cfg		*config.Config
	tracer	opentracing.Tracer
	mongoDB	*mongo.Client
	echo	*echo.Echo
	redis	*redis.Client
}

func NewServer(
	log		logger.Logger,
	cfg		*config.Config,
	tracer	opentracing.Tracer,
	mongoDB	*mongo.Client,
	redis	*redis.Client,
) *server {
	return &server{
		log: log,
		cfg: cfg,
		tracer: tracer,
		mongoDB: mongoDB,
		echo: echo.New(),
		redis: redis,
	}
}


func (s *server) Run() error{
	ctx,cancel := context.WithCancel(context.Background())
	defer cancel()

	validate := validator.New()
	
	productsProducer := kafka.NewProductsProducer(s.log,s.cfg)
	
	im := interceptors.NewInterceptorManger(s.log,s.cfg)
	mw := middleware.NewMiddlewareManger(s.log,s.cfg)

	grpcAddr := s.cfg.Server.Port
	if !strings.HasPrefix(grpcAddr,":") {	
		grpcAddr = ":" + grpcAddr
	}
	listener, err := net.Listen("tcp",grpcAddr)
	if err != nil {
		return errors.Wrap(err,"net.Listen")
	}

	gppcServer := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: s.cfg.Server.MaxConnectionIdle*time.Minute,
			Timeout: 			s.cfg.Server.Timeout*time.Second,
			MaxConnectionAge: 	s.cfg.Server.MaxConnectionAge*time.Minute,
			Time: 				s.cfg.Server.Timeout*time.Minute,
		}),
		grpc.ChainUnaryInterceptor(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_opentracing.UnaryServerInterceptor(),
			grpc_prometheus.UnaryServerInterceptor(),
			grpcrecovery.UnaryServerInterceptor(),
			im.Logger,
		),
	)
	
	prodSvc := productGrpc.NewProductService(s.log,productUC,validate)
	productsService.RegisterProductsServiceServer(grpcServer,prodSvc)
	grpc_prometheus.Register(grpcServer)

	v1 := s.echo.Group("/api/v1")
	v1.Use(mw.Metrics)

	productHandlers := productsHttpV1.newProductHandlers(s.log,productUC,validate,v1.Group("/products"),mw)
	productHandlers.MapRoutes()

	
	return nil
}