package grpc

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	successMessages = promauto.NewCounter(prometheus.CounterOpts{
		Name: "products_update_incoming_grpc_requests_total",
		Help: "The total number of success incoming success gRPC messages",
	})

	errorMessages = promauto.NewCounter(prometheus.CounterOpts{
		Name: "products_update_incoming_grpc_requests_total",
		Help: "The total number of error incoming succes gRPC messages",

	})

	createMessages = promauto.NewCounter(prometheus.CounterOpts{
		Name: "products_update_incoming_grpc_requests_total",
		Help: "The total number of incoming create products gRPC messages",

	})

	updateMessages = promauto.NewCounter(prometheus.CounterOpts{
		Name: "products_update_incoming_grpc_requests_total",
		Help: "The total number of incoming update products gRPC messages",

	})
	getByIdMessages = promauto.NewCounter(prometheus.CounterOpts{
		Name: "products_update_incoming_grpc_requests_total",
		Help: "The total number of incoming get by products gRPC messages",

	})
	
	searchMessages = promauto.NewCounter(prometheus.CounterOpts{
		Name: "products_update_incoming_grpc_requests_total",
		Help: "The total number of incoming search products gRPC messages",
	})
)