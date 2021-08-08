package queries

import (
	"context"
	"github.com/AleksK1NG/cqrs-microservices/pkg/logger"
	"github.com/AleksK1NG/cqrs-microservices/product_writer_service/config"
	"github.com/AleksK1NG/cqrs-microservices/product_writer_service/internal/models"
	"github.com/AleksK1NG/cqrs-microservices/product_writer_service/internal/product/repository"
)

type GetProductByIdHandler interface {
	Handle(ctx context.Context, query *GetProductByIdQuery) (*models.Product, error)
}

type getProductByIdHandler struct {
	log    logger.Logger
	cfg    *config.Config
	pgRepo repository.Repository
}

func NewGetProductByIdHandler(log logger.Logger, cfg *config.Config, pgRepo repository.Repository) *getProductByIdHandler {
	return &getProductByIdHandler{log: log, cfg: cfg, pgRepo: pgRepo}
}

func (q *getProductByIdHandler) Handle(ctx context.Context, query *GetProductByIdQuery) (*models.Product, error) {
	return q.pgRepo.GetProductById(ctx, query.ProductID)
}
