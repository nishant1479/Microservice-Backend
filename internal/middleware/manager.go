package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/nishant1479/Microservice-Backend/config"
	"github.com/nishant1479/Microservice-Backend/pkg/logger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	httpTotalRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "http_microservices_total_requests",
		Help: "The total number of incoming HTTP requests",
	})
)

type middlewareManager struct{
	log	logger.Logger
	cfg	*config.Config
}

type MiddlewareManager interface{
	Metrics(next echo.HandlerFunc) echo.HandlerFunc
}

func NewMiddlewareManger(log logger.Logger,cfg *config.Config) *middlewareManager{
	return &middlewareManager{log: log,cfg: cfg}
}

func (m *middlewareManager) Metrics(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		httpTotalRequests.Inc()
		return next(c)
	}
}