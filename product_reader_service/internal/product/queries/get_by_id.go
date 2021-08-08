package queries

import (
	"context"
	"github.com/AleksK1NG/cqrs-microservices/pkg/logger"
	"github.com/AleksK1NG/cqrs-microservices/product_reader_service/config"
	"github.com/AleksK1NG/cqrs-microservices/product_reader_service/internal/models"
	"github.com/AleksK1NG/cqrs-microservices/product_reader_service/internal/product/repository"
	"github.com/opentracing/opentracing-go"
)

type GetProductByIdHandler interface {
	Handle(ctx context.Context, query *GetProductByIdQuery) (*models.Product, error)
}

type getProductByIdHandler struct {
	log       logger.Logger
	cfg       *config.Config
	mongoRepo repository.Repository
	redisRepo repository.CacheRepository
}

func NewGetProductByIdHandler(log logger.Logger, cfg *config.Config, mongoRepo repository.Repository, redisRepo repository.CacheRepository) *getProductByIdHandler {
	return &getProductByIdHandler{log: log, cfg: cfg, mongoRepo: mongoRepo, redisRepo: redisRepo}
}

func (q *getProductByIdHandler) Handle(ctx context.Context, query *GetProductByIdQuery) (*models.Product, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "getProductByIdHandler.Handle")
	defer span.Finish()

	if product, err := q.redisRepo.GetProduct(ctx, query.ProductID.String()); err == nil && product != nil {
		return product, nil
	}

	product, err := q.mongoRepo.GetProductById(ctx, query.ProductID)
	if err != nil {
		return nil, err
	}

	q.redisRepo.PutProduct(ctx, product.ProductID, product)
	return product, nil
}
