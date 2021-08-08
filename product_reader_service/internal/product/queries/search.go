package queries

import (
	"context"
	"github.com/AleksK1NG/cqrs-microservices/pkg/logger"
	"github.com/AleksK1NG/cqrs-microservices/product_reader_service/config"
	"github.com/AleksK1NG/cqrs-microservices/product_reader_service/internal/models"
	"github.com/AleksK1NG/cqrs-microservices/product_reader_service/internal/product/repository"
)

type SearchProductHandler interface {
	Handle(ctx context.Context, query *SearchProductQuery) (*models.ProductsList, error)
}

type searchProductHandler struct {
	log       logger.Logger
	cfg       *config.Config
	mongoRepo repository.Repository
	redisRepo repository.CacheRepository
}

func NewSearchProductHandler(log logger.Logger, cfg *config.Config, mongoRepo repository.Repository, redisRepo repository.CacheRepository) *searchProductHandler {
	return &searchProductHandler{log: log, cfg: cfg, mongoRepo: mongoRepo, redisRepo: redisRepo}
}

func (s *searchProductHandler) Handle(ctx context.Context, query *SearchProductQuery) (*models.ProductsList, error) {
	return s.mongoRepo.Search(ctx, query.Text, query.Pagination)
}
