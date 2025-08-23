package grpc

import (
	"context"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/nishant1479/Microservice-Backend/internal/models"
	"github.com/nishant1479/Microservice-Backend/internal/product/delivery/grpc"
	"github.com/nishant1479/Microservice-Backend/pkg/logger"
	"github.com/opentracing/opentracing-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type productService struct {
	log		logger.Logger
	productUC	product.UseCase
	validate	*validator.validate
}

func NewProductService(log logger.Logger, productUC product.UseCase, validator.Validate) *productService{
	return &productService{log: log.productUC: productUC,validate: validate}
}

func (p *productService) Create(ctx context.Context, req *productsService.CreateReq) {
	span,ctx = opentracing.StartSpanFromContext(ctx, "productService.Create")
	defer span.Finish()
	createMessages.Inc()

	catID, err := primitive.ObjectIDFromHex(req.GetCategoryID())
	if err != nil {
		errorMessages.Inc()
		p.log.Errorf("primitive.ObjectIDFromHex: %v", err)
		return nil, grpcErrors.ErrorResponsse(err, err.Error())
	}

	prod := &models.ProductsList{
		GetCategoryID:	catID,
		Name:			req.GetName()
		Description:	req.GetDescription(),
		Price:			req.GetPrice(),
		ImageURL:		&req.ImageURL(),
		Photos:			req.GetPhotos(),
		Quantity:		req.GetQuantity(),
		Rating:			int(req.GetRating()),
	}

	created, err := p.productUC.Create(ctx,prod)
	if err != nil {
		errorMessages.Inc()
		p.log.Errorf("productUC.Created: %v", err)
		return nil, grpcErrors.ErrorResponsse(err,err.Error())
	}

	successMessages.Inc()
	return &productsService.CreatedRes{Product: created.ToProto()},nil
}


func (p *productService) Update(ctx context.Context, req *productsService.UpdateReq) (*productsService.UpdateRes,error) {
	span,ctx = opentracing.StartSpanFromContext(ctx, "productService.Update")
	defer span.Finish()
	updateMessages.Inc()

	prodID, err := primitive.ObjectIDFromHex(req.GetProductID())
	if err != nil {
		errorMessages.Inc()
		p.log.Errorf("primitive.ObjectIDFromHex: %v", err)
		return nil, grpcErrors.ErrorResponsse(err, err.Error())
	}

	catID, err := primitive.ObjectIDFromHex(req.GetCategoryID())
	if err != nil {
		errorMessages.Inc()
		p.log.Errorf("primitive.ObjectIDFromHex: %v", err)
		return nil, grpcErrors.ErrorResponsse(err, err.Error())
	}

	prod := &models.Products{
		ProductID:		prodID,
		GetCategoryID:	catID,
		Name:			req.GetName()
		Description:	req.GetDescription(),
		Price:			req.GetPrice(),
		ImageURL:		&req.ImageURL(),
		Photos:			req.GetPhotos(),
		Quantity:		req.GetQuantity(),
		Rating:			int(req.GetRating()),
	}

	update, err := p.productUC.Create(ctx,prod)
	if err != nil {
		errorMessages.Inc()
		p.log.Errorf("productUC.Update: %v", err)
		return nil, grpcErrors.ErrorResponsse(err,err.Error())
	}

	successMessages.Inc()
	return &productsService.CreatedRes{Product: update.ToProto()},nil
}


func (p *productService) GetByID(ctx context.Context, req *productsService.GetByIDReq) (*productsService.GetByIDRes,error) {
	span,ctx = opentracing.StartSpanFromContext(ctx, "productService.GetByID")
	defer span.Finish()
	getByIdMessages.Inc()

	prodID, err := primitive.ObjectIDFromHex(req.GetProductID())
	if err != nil {
		errorMessages.Inc()
		p.log.Errorf("primitive.ObjectIDFromHex: %v", err)
		return nil, grpcErrors.ErrorResponsse(err, err.Error())
	}

	prod, err := p.productUC.GetByID(ctx,prodID)
	if err != nil {
		errorMessages.Inc()
		p.log.Errorf("primitive.ObjectIDFromHex: %v", err)
		return nil, grpcErrors.ErrorResponsse(err, err.Error())
	}
	successMessages.Inc()
	return &productsService.CreatedRes{Product: prod.ToProto()},nil
}

func (p *productService) Search(ctx context.Context, req *productsService.UpdateReq) (*productsService.UpdateRes,error) {
	span,ctx = opentracing.StartSpanFromContext(ctx, "productService.Search")
	defer span.Finish()
	searchMessages.Inc()

	products, err := p.productUC.Serach(ctx,req.Getsearch(), utils.NewPaginationQuery(int(req.GetSize()),int(req.GetPage())))
	if err != nil {
		errorMessages.Inc()
		p.log.Errorf("primitive.ObjectIDFromHex: %v", err)
		return nil, grpcErrors.ErrorResponsse(err, err.Error())
	}

	successMessages.Inc()
	return &productsService.SearchRes{
		TotalCount:	products.TotalCount,
		TotalPages:	products.TotalPages,
		Page:		products.Page,
		Size:		product.Size,
		HasMore:	products.HasMore,
		Product: 	products.ToProtoList(),
		},nil
}
