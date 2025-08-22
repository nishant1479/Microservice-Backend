package server

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/swaggo/swag/example/basic/docs"
)

func (s *server) runHTTPServer() error {
	s.echo.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK,"Server running on HTTP")
	})

	s.echo.GET("/test",func(c echo.Context) error {
		return c.String(http.StatusOK,"Test route working!")
	})

	s.echo.GET("/health",func(c echo.Context) error {
		return c.String(http.StatusOK,"Ok")

	})

	s.echo.GET("/metrics".echo.WrapHandler(promhttp.Handler()))
	
	docs.SwaggerInfo.Title = "Products microservice"
	docs.SwaggerInfo.Description = "Products REST API microservice."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Host = "localhost:5007"
	docs.SwaggerInfo.Schemes = []string{"http"}

	s.echo.GET("/swagger/*", echoSwagger.WrapHandler)
	s.echo.GET("/swagger", func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})

	s.echo.Use(middleware.Logger())
	s.echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderXRequestID,
			csrfTokenHeader,
		},
	}))
	s.echo.Use(middleware.Recover())
	s.echo.Use(middleware.RequestID())
	s.echo.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: gzipLevel,
		Skipper: func(c echo.Context) bool {
			return strings.Contains(c.Request().URL.Path, "swagger")
		},
	}))
	s.echo.Use(middleware.BodyLimit(bodyLimit))

	addr := s.cfg.Http.Port
	if !strings.HasPrefix(addr, ":") {
		addr = ":" + addr
	}

	s.echo.Server.ReadTimeout = time.Second * s.cfg.Http.ReadTimeout
	s.echo.Server.WriteTimeout = time.Second * s.cfg.Http.WriteTimeout
	s.echo.Server.MaxHeaderBytes = maxHeaderBytes

	s.log.Infof("Starting HTTP server on %s", addr)
	return s.echo.Start(addr)

}
