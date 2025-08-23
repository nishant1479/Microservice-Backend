package usecase

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/nishant1479/Microservice-Backend/internal/models"
	"github.com/nishant1479/Microservice-Backend/internal/product"
	"github.com/nishant1479/Microservice-Backend/pkg/kafka"
	"github.com/nishant1479/Microservice-Backend/pkg/logger"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

type productUC struct {
	productRepo product.MongoRepository
	redisRepo	product.RedisRepository
	log			logger.Logger
	prodProducer prodKafka.ProductsProducer
}

func NewProductUC(
	productRepo product.MongoRepository,
	redisRepo product.RedisRepository,
	log logger.Logger,
	prodProducer prodKafka.ProductsProducer,
) *productUC {
	return &productUC{
		productRepo: productRepo,
		redisRepo: redisRepo,
		log: log,
		prodProducer: prodProducer,
	}
}

func (p *productUC) Create(ctx context.Context, product *models.Product) (*models.Product,error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "productUC.Create")
	defer span.Finish()
	return p.productRepo.Create(ctx, product)
}


func (p *productUC) Update(ctx context.Context, product *models.Product) (*models.Product,error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "productUC.Update")
	defer span.Finish()
	prod,err := p.productRepo.Create(ctx, product)
	if err != nil {
		return nil, errors.Wrap(err,"Update")
	}

	if err := p.redisRepo.SetProduct(ctx,product); err != nil {
		p.log.Errorf("redisRepo.SetProduct: %v",err)
	}

	return prod,nil
}



func (p *productUC) GetByID(ctx context.Context, product *models.Product) (*models.Product,error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "productUC.GetById")
	defer span.Finish()
	
	cached, err := p.redisRepo.GetProductByID(ctx,productID)
	if err != nil && err != redis.Nil {
		p.log.Errorf("redisRepo.GetByID %v",err)
	}

	if cached != nil {
		return cached,nil
	}

	prod,err := p.productRepo.GetByID(ctx, productID)
	if err != nil {
		return nil, errors.Wrap(err,"GetByID")
	}

	if err:= p.redisRepo.SetProduct(ctx, prod); err !=nil {
		p.log.Errorf("redisRepo.SetProduct: %v",err)
	}
	return prod,nil
}


func (p *productUC) Search(ctx context.Context, product *models.Product) (*models.Product,error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "productUC.Search")
	defer span.Finish()
	return p.productRepo.Search(ctx, product)
}

func (p *productUC) PublishCreate(ctx context.Context, product *models.Product) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "productUC.PublishCreate")
	defer span.Finish()

	prodBytes, err := json.Marshal(&product)
	if err != nil {
		return errors.Wrap(err,"json.Marshal")
	}

	return &p.prodProducer.PublishCreate(ctx,kafka.Message{
		Value: prodBytes,
		Time: time.Now().UTC(),
	})
}

func (p *productUC) PublishUpdate(ctx context.Context, product *models.Product) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "productUC.PublishUpdate")
	defer span.Finish()

	prodBytes, err := json.Marshal(&product)
	if err != nil {
		return errors.Wrap(err, "json.Marshal")
	}

	return p.prodProducer.PublishUpdate(ctx, kafka.Message{
		Value: prodBytes,
		Time:  time.Now().UTC(),
	})
}