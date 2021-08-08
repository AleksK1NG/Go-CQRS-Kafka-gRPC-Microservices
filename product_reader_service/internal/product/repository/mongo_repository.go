package repository

import (
	"context"
	"github.com/AleksK1NG/cqrs-microservices/pkg/logger"
	"github.com/AleksK1NG/cqrs-microservices/pkg/utils"
	"github.com/AleksK1NG/cqrs-microservices/product_reader_service/config"
	"github.com/AleksK1NG/cqrs-microservices/product_reader_service/internal/models"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoRepository struct {
	log logger.Logger
	cfg *config.Config
	db  *mongo.Client
}

func NewMongoRepository(log logger.Logger, cfg *config.Config, db *mongo.Client) *mongoRepository {
	return &mongoRepository{log: log, cfg: cfg, db: db}
}

func (p *mongoRepository) CreateProduct(ctx context.Context, product *models.Product) (*models.Product, error) {
	collection := p.db.Database(p.cfg.Mongo.Db).Collection(p.cfg.MongoCollections.Products)

	insertedID, err := collection.InsertOne(ctx, product, &options.InsertOneOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "InsertOne")
	}

	p.log.Infof("INSERTED ID: %+v", insertedID)

	return product, nil
}

func (p *mongoRepository) UpdateProduct(ctx context.Context, product *models.Product) (*models.Product, error) {
	collection := p.db.Database(p.cfg.Mongo.Db).Collection(p.cfg.MongoCollections.Products)

	ops := options.FindOneAndUpdate()
	ops.SetReturnDocument(options.After)
	ops.SetUpsert(true)

	var updated models.Product
	if err := collection.FindOneAndUpdate(ctx, bson.M{"_id": product.ProductID}, bson.M{"$set": product}, ops).Decode(&updated); err != nil {
		return nil, errors.Wrap(err, "Decode")
	}

	p.log.Infof("UPDATED ID: %+v", updated)
	return &updated, nil
}

func (p *mongoRepository) GetProductById(ctx context.Context, uuid uuid.UUID) (*models.Product, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "mongoRepository.GetProductById")
	defer span.Finish()

	collection := p.db.Database(p.cfg.Mongo.Db).Collection(p.cfg.MongoCollections.Products)

	var product models.Product
	if err := collection.FindOne(ctx, bson.M{"_id": uuid.String()}).Decode(&product); err != nil {
		return nil, errors.Wrap(err, "Decode")
	}

	p.log.Infof("GET PRODUCT BY ID ID: %+v", product)
	return &product, nil
}

func (p *mongoRepository) DeleteProduct(ctx context.Context, uuid uuid.UUID) error {
	collection := p.db.Database(p.cfg.Mongo.Db).Collection(p.cfg.MongoCollections.Products)

	return collection.FindOneAndDelete(ctx, bson.M{"_id": uuid}).Err()
}

func (p *mongoRepository) Search(ctx context.Context, search string, pagination *utils.Pagination) (*models.ProductsList, error) {
	collection := p.db.Database(p.cfg.Mongo.Db).Collection(p.cfg.MongoCollections.Products)

	filter := bson.D{
		{Key: "$or", Value: bson.A{
			bson.D{{Key: "name", Value: primitive.Regex{Pattern: search, Options: "gi"}}},
			bson.D{{Key: "description", Value: primitive.Regex{Pattern: search, Options: "gi"}}},
		}},
	}

	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, errors.Wrap(err, "CountDocuments")
	}
	if count == 0 {
		return &models.ProductsList{Products: make([]*models.Product, 0)}, nil
	}

	limit := int64(pagination.GetLimit())
	skip := int64(pagination.GetOffset())
	cursor, err := collection.Find(ctx, filter, &options.FindOptions{
		Limit: &limit,
		Skip:  &skip,
	})
	if err != nil {
		return nil, errors.Wrap(err, "Find")
	}
	defer cursor.Close(ctx)

	products := make([]*models.Product, 0, pagination.GetSize())

	for cursor.Next(ctx) {
		var prod models.Product
		if err := cursor.Decode(&prod); err != nil {
			return nil, errors.Wrap(err, "Find")
		}
		products = append(products, &prod)
	}

	if err := cursor.Err(); err != nil {
		return nil, errors.Wrap(err, "cursor.Err")
	}

	return models.NewProductListWithPagination(products, count, pagination), nil
}
