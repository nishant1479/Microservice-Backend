package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (p *productHandlers) MapRoutes() {

	p.group.Get("/test", func(c echo.Context)error{
		return c.String(http.StatusOK,"Hello form test endpoint")
	})
	p.group.POST("",p.CreateProduct())
	p.group.PUT("/:product_id",p.UpdateProduct())
	p.group.GET("/:product_id",p.GetByIDProduct())
	p.group.GET("/search",p.SearchProduct())
}