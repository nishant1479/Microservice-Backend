package kafka

import (
	"context"

	"github.com/nishant1479/Microservice-Backend/config"
	"github.com/nishant1479/Microservice-Backend/pkg/logger"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/compress"
)

type ProductsProducer interface {
	PublishCreate(ctx context.Context, msgs ...kafka.Message) error
	PublishUpdate(ctx context.Context, msgs ...kafka.Message) error
	Close()
	Run()
	GetNewKafkaWriter(topic string) *kafka.Writer
}

type productsProducer struct {
	log          logger.Logger
	cfg          *config.Config
	createWriter *kafka.Writer
	updateWriter *kafka.Writer
}

func NewProductsProducer(log logger.Logger, cfg *config.Config) *productsProducer {
	return &productsProducer{
		log: log,
		cfg: cfg,
	}
}

func (p *productsProducer) GetNewKafkaWriter(topic string) *kafka.Writer  {
	w := &kafka.Writer{
		Addr:			kafka.TCP(p.cfg.Kafka.Brokers...),
		Topic:			topic,
		Balancer:		&kafka.LeastBytes{},
		RequiredAcks:	writerRequiredAcks,
		MaxAttempts:	writerMaxAttempts,
		Logger:			kafka.LoggerFunc(p.log.Debugf),
		ErrorLogger:	kafka.LoggerFunc(p.log.Errorf),
		Compression:	compress.Snappy,
		ReadTimeout:	writerReadTimeout,
		WriteTimeout:	writerWriteTimeout,
	}
	return w
}

func (p productsProducer) Close() {
	p.createWriter.Close()
	p.updateWriter.Close()
}

func (p *productsProducer) PublishCreate(ctx context.Context, msgs ...kafka.Message) error {
	return p.createWriter.WriteMessages(ctx, msgs...)
}


func (p *productsProducer) PublishUpdate(ctx context.Context, msgs ...kafka.Message) error {
	return p.updateWriter.WriteMessages(ctx, msgs...)
}