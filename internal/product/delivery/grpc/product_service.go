package grpc

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/nishant1479/Microservice-Backend/internal/models"
	"github.com/nishant1479/Microservice-Backend/internal/product"
	"github.com/nishant1479/Microservice-Backend/pkg/logger"
	grpcErrors "github.com/nishant1479/Microservice-Backend/pkg/grpc_errors"
	productsService "github.com/nishant1479/Microservice-Backend/proto/product"
	"github.com/opentracing/opentracing-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type productService struct {
	log		logger.Logger
	productUC	product.UseCase
	validate	*validator.Validate
}

func NewProductService(log logger.Logger, productUC product.UseCase, validate *validator.Validate) *productService{
	return &productService{log: log,productUC: productUC,validate: validate}
}

func (p *productService) Create(ctx context.Context, req *productsService.CreateRequest)  (*productsService.CreateResponse,error) {
	span,ctx := opentracing.StartSpanFromContext(ctx, "productService.Create")
	defer span.Finish()
	createMessages.Inc()

	catID, err := primitive.ObjectIDFromHex(req.GetCategoryID())
	if err != nil {
		errorMessages.Inc()
		p.log.Errorf("primitive.ObjectIDFromHex: %v", err)
		return nil, grpcErrors.ErrorResponse(err, err.Error())
	}

	prod := &models.Product{
		CategoryID:		catID,
		Name:			req.GetName(),
		Description:	req.GetDescription(),
		Price:			req.GetPrice(),
		ImageURL:		&req.ImageURL(),
		Photos:			req.GetPhotos(),
		Quantity:		req.GetQuantity(),
		Rating:			int64(req.GetRating()),
	}

	created, err := p.productUC.Create(ctx,prod)
	if err != nil {
		errorMessages.Inc()
		p.log.Errorf("productUC.Created: %v", err)
		return nil, grpcErrors.ErrorResponse(err,err.Error())
	}

	successMessages.Inc()
	return &productsService.CreateResponse{Product: created.ToProto()},nil
}


func (p *productService) Update(ctx context.Context, req *productsService.UpdateRequest) (*productsService.UpdateResponse,error) {
	span,ctx := opentracing.StartSpanFromContext(ctx, "productService.Update")
	defer span.Finish()
	updateMessages.Inc()

	prodID, err := primitive.ObjectIDFromHex(req.GetProductID())
	if err != nil {
		errorMessages.Inc()
		p.log.Errorf("primitive.ObjectIDFromHex: %v", err)
		return nil, grpcErrors.ErrorResponse(err, err.Error())
	}

	catID, err := primitive.ObjectIDFromHex(req.GetCategoryID())
	if err != nil {
		errorMessages.Inc()
		p.log.Errorf("primitive.ObjectIDFromHex: %v", err)
		return nil, grpcErrors.ErrorResponse(err, err.Error())
	}

	prod := &models.Product{
		ProductID:		prodID,
		CategoryID:		catID,
		Name:			req.GetName(),
		Description:	req.GetDescription(),
		Price:			req.GetPrice(),
		ImageURL:		&req.ImageURL(),
		Photos:			req.GetPhotos(),
		Quantity:		req.GetQuantity(),
		Rating:			int64(req.GetRating()),
	}

	update, err := p.productUC.Create(ctx,prod)
	if err != nil {
		errorMessages.Inc()
		p.log.Errorf("productUC.Update: %v", err)
		return nil, grpcErrors.ErrorResponse(err,err.Error())
	}

	successMessages.Inc()
	return &productsService.UpdateResponse{Product: update.ToProto()},nil
}


func (p *productService) GetByID(ctx context.Context, req *productsService.GetByIDRequest) (*productsService.GetByIDResponse,error) {
	span,ctx := opentracing.StartSpanFromContext(ctx, "productService.GetByID")
	defer span.Finish()
	getByIdMessages.Inc()

	prodID, err := primitive.ObjectIDFromHex(req.GetProductID())
	if err != nil {
		errorMessages.Inc()
		p.log.Errorf("primitive.ObjectIDFromHex: %v", err)
		return nil, grpcErrors.ErrorResponse(err, err.Error())
	}

	prod, err := p.productUC.GetByID(ctx,prodID)
	if err != nil {
		errorMessages.Inc()
		p.log.Errorf("primitive.ObjectIDFromHex: %v", err)
		return nil, grpcErrors.ErrorResponse(err, err.Error())
	}
	successMessages.Inc()
	return &productsService.GetByIDResponse{Product: prod.ToProto()},nil
}

func (p *productService) Search(ctx context.Context, req *productsService.SearchReq) (*productsService.SearchResponse,error) {
	span,ctx := opentracing.StartSpanFromContext(ctx, "productService.Search")
	defer span.Finish()
	searchMessages.Inc()

	products, err := p.productUC.Search(ctx,req.GetSearch(), utils.NewPaginationQuery(int(req.GetSize()),int(req.GetPage())))
	if err != nil {
		errorMessages.Inc()
		p.log.Errorf("primitive.ObjectIDFromHex: %v", err)
		return nil, grpcErrors.ErrorResponse(err, err.Error())
	}

	successMessages.Inc()
	return &productsService.SearchResponse{
		TotalCount:	products.TotalCount,
		TotalPages:	products.TotalPages,
		Page:		products.Page,
		Size:		products.Size,
		HasMore:	products.HasMore,
		Product: 	products.ToProtoList(),
		},nil
}
