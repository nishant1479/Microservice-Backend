package product

import (
	"context"

	"github.com/nishant1479/Microservice-Backend/internal/models"
	"github.com/nishant1479/Microservice-Backend/pkg/utlis"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UseCase Product
type UseCase interface {
	Create(ctx context.Context, product *models.Product) (*models.Product, error)
	Update(ctx context.Context, product *models.Product) (*models.Product, error)
	GetByID(ctx context.Context, productID primitive.ObjectID) (*models.Product, error)
	Search(ctx context.Context, search string, pagination *utlis.Pagination) (*models.ProductsList, error)
	PublishCreate(ctx context.Context, product *models.Product) error
	PublishUpdate(ctx context.Context, product *models.Product) error
}