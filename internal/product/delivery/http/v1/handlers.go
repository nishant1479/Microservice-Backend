package v1

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/nishant1479/Microservice-Backend/internal/middleware"
	"github.com/nishant1479/Microservice-Backend/internal/product"
	"github.com/nishant1479/Microservice-Backend/pkg/logger"
)

type productHandlers struct {
	log 	logger.Logger
	productUC	product.UseCase
	validate	*validator.Validate
	group		*echo.Group
	mw			middleware.MiddlewareManager
}

func newProductHandlers(
	log logger.Logger,
	productUC product.UseCase,
	validate validator.Validate,
	group *echo.Group,
	mw middleware.MiddlewareManager,
) *productHandlers {
	return &productHandlers{log: log,productUC: productUC,validate: &validate,group: group,mw: mw}
}

CreateProduct()

UpdateProduct()

GetByIDProduct()

SearchProduct()