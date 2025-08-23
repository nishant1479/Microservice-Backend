package v1

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	succesRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "https_products_success_incoming_messages_total",
		Help: "The total number of succes incoming succes HTTP requests",
	})
	errorRequests = promauto.NewCounter(prometheus.CounterOpts{
	Name: "https_products_success_incoming_messages_total",
		Help: "The total number of succes incoming succes HTTP requests",
	})
	createRequests = promauto.NewCounter(prometheus.CounterOpts{
Name: "https_products_success_incoming_messages_total",
		Help: "The total number of succes incoming succes HTTP requests",
	})
	updateRequests = promauto.NewCounter(prometheus.CounterOpts{
Name: "https_products_success_incoming_messages_total",
		Help: "The total number of succes incoming succes HTTP requests",
	})
	getByIdRequests = promauto.NewCounter(prometheus.CounterOpts{
Name: "https_products_success_incoming_messages_total",
		Help: "The total number of succes incoming succes HTTP requests",
	})
	SearchRequests = promauto.NewCounter(prometheus.CounterOpts{
Name: "https_products_success_incoming_messages_total",
		Help: "The total number of succes incoming succes HTTP requests",
	})
)