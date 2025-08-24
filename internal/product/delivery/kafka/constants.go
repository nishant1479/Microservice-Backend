package kafka

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	incomingMessages = promauto.NewCounter(prometheus.CounterOpts{

	})

	successMessages = promauto.NewCounter(prometheus.CounterOpts{

	})

	errorMessages = promauto.NewCounter(prometheus.CounterOpts{

	})
)

const (
	minBytes               = 10e3 // 10KB
	maxBytes               = 10e6 // 10MB
	queueCapacity          = 100
	heartbeatInterval      = 3 * time.Second
	commitInterval         = 0
	partitionWatchInterval = 5 * time.Second
	maxAttempts            = 3
	dialTimeout            = 3 * time.Minute

	writerReadTimeout  = 10 * time.Second
	writerWriteTimeout = 10 * time.Second
	writerRequiredAcks = -1
	writerMaxAttempts  = 3

	createProductTopic   = "create-product"
	createProductWorkers = 3
	updateProductTopic   = "update-product"
	updateProductWorkers = 3

	deadLetterQueueTopic = "dead-letter-queue"

	productsGroupID = "products_group"
)