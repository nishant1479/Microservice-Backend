package kafka

import (
	"context"

	"github.com/nishant1479/Microservice-Backend/config"
	"github.com/segmentio/kafka-go"
)

func NewKafkaConn(cfg *config.Config) (*kafka.Conn,error) {
	return kafka.DialContext(context.Background(),"tcp",cfg.Kafka.Brokers[0])
}