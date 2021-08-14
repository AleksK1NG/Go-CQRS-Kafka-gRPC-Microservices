package commands

import (
	"context"
	"github.com/AleksK1NG/cqrs-microservices/pkg/logger"
	"github.com/AleksK1NG/cqrs-microservices/reader_service/config"
	"github.com/AleksK1NG/cqrs-microservices/reader_service/internal/product/repository"
	"github.com/opentracing/opentracing-go"
)

type DeleteProductCmdHandler interface {
	Handle(ctx context.Context, command *DeleteProductCommand) error
}

type deleteProductCmdHandler struct {
	log       logger.Logger
	cfg       *config.Config
	mongoRepo repository.Repository
	redisRepo repository.CacheRepository
}

func NewDeleteProductCmdHandler(log logger.Logger, cfg *config.Config, mongoRepo repository.Repository, redisRepo repository.CacheRepository) *deleteProductCmdHandler {
	return &deleteProductCmdHandler{log: log, cfg: cfg, mongoRepo: mongoRepo, redisRepo: redisRepo}
}

func (c *deleteProductCmdHandler) Handle(ctx context.Context, command *DeleteProductCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "deleteProductCmdHandler.Handle")
	defer span.Finish()

	if err := c.mongoRepo.DeleteProduct(ctx, command.ProductID); err != nil {
		return err
	}

	c.redisRepo.DelProduct(ctx, command.ProductID.String())
	return nil
}
