package main

import (
	"context"
	"log"

	"github.com/nishant1479/Microservice-Backend/config"
	jaeger "github.com/nishant1479/Microservice-Backend/pkg/Jaeger"
	"github.com/nishant1479/Microservice-Backend/pkg/logger"
	"github.com/opentracing/opentracing-go"
)

func main() {
	log.Println("Initializing the infra services")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cfg, err := config.ParseConfig()
	if err != nil {
		log.Fatal(err)
	}

	appLogger := logger.NewApiLogger(cfg)
	appLogger.InitLogger()
	appLogger.Info("Starting user server")
	appLogger.Infof(
		"AppVersion: %s, LogLevel: %s,DevelopmentMode: %s",
		cfg.AppVersion,
		cfg.Logger.Level,
		cfg.Server.Development,
	)
	appLogger.Infof("Success parsed config: %#v",cfg.AppVersion)

	tracer, closer,err := jaeger.InitJaeger(cfg)
	if err != nil {
		appLogger.Fatal("cannot create tracer",err)
	}
	appLogger.Info("Jaeger connected")

	
	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()
	appLogger.Info("Opentracing connected")

	mongoDBConn := cfg.MongoDB
	conn ,err := kafka.NewKafkaConn(cfg)
	if err != nil {
		appLogger.Fatal("NewKafkaConn",err)
	}
	defer conn.Close()
	brokers,err := conn.Brokers()
	if err != nil {
		appLogger.Fatal("conn.Brokers",err)
	}
	appLogger.Infof("Kafka connected: %v",brokers)

	redisClient := redis.NewRedisClient(cfg)
	appLogger.Info("Redis connected")

	s:= server.NewServer(appLogger,cfg,tracer,mongoDBConn,redisClient)
	appLogger.Fatal(s.Run())

}