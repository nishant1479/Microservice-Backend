package server

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/nishant1479/Microservice-Backend/config"
	"github.com/nishant1479/Microservice-Backend/pkg/logger"
	"github.com/opentracing/opentracing-go"
	"github.com/redis/go-redis"
	"go.mongodb.org/mongo-driver/mongo"
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
	
	return nil
}