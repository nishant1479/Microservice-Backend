package repository

import (
	"context"
	"time"
	"github.com/pkg/errors"
productErrors "github.com/nishant1479/Microservice-Backend/pkg/product_errors"
	"github.com/nishant1479/Microservice-Backend/pkg/utlis"
	"github.com/nishant1479/Microservice-Backend/internal/models"
	"github.com/opentracing/opentracing-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	productsDB        = "products"
	productCollection = "products"
)

type productMongoRep struct {
	mongoDB *mongo.Client
}

func NewProductMongoRepo(mongoDB *mongo.Client) *productMongoRep {
	return &productMongoRep{
		mongoDB: mongoDB,
	}
}

func (p *productMongoRep) Create(ctx context.Context, product *models.Product) (*models.Product, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "productMongoRepo.Create")
	defer span.Finish()

	collection := p.mongoDB.Database(productsDB).Collection(productCollection)

	product.CreatedAt = time.Now().UTC()
	product.UpdatedAt = time.Now().UTC()

	result, err := collection.InsertOne(ctx, product, &options.InsertOneOptions{})
	if err != nil {
		return nil, errors.Wrap(err,"InsertOne")
	}

	objectID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok{
		return nil, errors.Wrap(productErrors.ErrObjectIDTypeConversion, "result,InsertedID")
	}

	product.ProductID = objectID
	
	return product, nil
}

func (p *productMongoRep) Update(ctx context.Context, product *models.Product) (*models.Product,error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "productMongoRepo.Update")
	defer span.Finish()

	collection := p.mongoDB.Database(productsDB).Collection(productCollection)

	ops := options.FindOneAndUpdate()
	ops.SetReturnDocument(options.After)
	ops.SetUpsert(true)

	var prod models.Product
	if err := collection.FindOneAndUpdate(ctx, bson.M{"_id": product.ProductID}, bson.M{"$set": product}).Decode(&prod); err != nil{
		return nil ,errors.Wrap(err, "Decode")
	}
	return &prod,nil
}

func (p *productMongoRep) GetByID(ctx context.Context, productID primitive.ObjectID) (*models.Product, error) {
	span,ctx := opentracing.StartSpanFromContext(ctx, "productMongoRepo.GetByID")
	defer span.Finish()

	collection := p.mongoDB.Database(productsDB).Collection(productCollection)

	var prod models.Product
	if err := collection.FindOne(ctx, bson.M{"_id": productID}).Decode(*&prod); err != nil {
		return nil,errors.Wrap(err,"Decode")
	}

	return &prod,nil
}

func (p *productMongoRep) Search(ctx context.Context, search string, pagination *utlis.Pagination) ( *models.ProductsList, error){
	span, ctx := opentracing.StartSpanFromContext(ctx, "productMongoRepo.Search")
	defer span.Finish()

	collection := p.mongoDB.Database(productsDB).Collection(productCollection)

	f := bson.D{
		{Key: "$or", Value: bson.A{
			bson.D{{Key: "name", Value: primitive.Regex{
				Pattern: search,
				Options: "gi",
			}}},
			bson.D{{Key: "description", Value: primitive.Regex{
				Pattern: search,
				Options: "gi",
			}}},
		}},
	}

	count, err := collection.CountDocuments(ctx,f)
	if err != nil {
		return nil, errors.Wrap(err,"CountDocuments")
	}
	if count == 0 {
		return &models.ProductsList{
			TotalCount: 0,
			TotalPages: 0,
			Page: 0,
			Size: 0,
			HasMore: false,
			Products: make([]*models.Product, 0),
		},nil
	}

	limit := int64(pagination.GetLimit())
	skip := int64(pagination.GetOffSet())
	cursor, err := collection.Find(ctx, f, &options.FindOptions{
		Limit: &limit,
		Skip: &skip,
	})
	if err != nil {
		return nil, errors.Wrap(err, "Find")
	}
	defer cursor.Close(ctx)

	products := make([]*models.Product,0,pagination.GetSize())
	for cursor.Next(ctx){
		var prod models.Product
		if err := cursor.Decode(&prod); err != nil{
			return nil, errors.Wrap(err, "Find")
		}
		products = append(products, &prod)
	}

	if err := cursor.Err(); err != nil {
		return nil,errors.Wrap(err, "curson.Err")
	}
	return &models.ProductsList{
		TotalCount: count,
		TotalPages: int64(pagination.GetTotalPages(int(count))),
		Page:       int64(pagination.GetPage()),
		Size:       int64(pagination.GetSize()),
		HasMore:    pagination.GetHasMore(int(count)),
		Products:   products,
	}, nil
}