package queries

import (
	"context"
	"github.com/AleksK1NG/cqrs-microservices/api_gateway_service/config"
	"github.com/AleksK1NG/cqrs-microservices/api_gateway_service/internal/dto"
	"github.com/AleksK1NG/cqrs-microservices/pkg/logger"
	readerService "github.com/AleksK1NG/cqrs-microservices/product_reader_service/proto/product_reader"
)

type GetProductByIdHandler interface {
	Handle(ctx context.Context, query *GetProductByIdQuery) (*dto.ProductResponse, error)
}

type getProductByIdHandler struct {
	log      logger.Logger
	cfg      *config.Config
	rsClient readerService.ReaderServiceClient
}

func NewGetProductByIdHandler(log logger.Logger, cfg *config.Config, rsClient readerService.ReaderServiceClient) *getProductByIdHandler {
	return &getProductByIdHandler{log: log, cfg: cfg, rsClient: rsClient}
}

func (q *getProductByIdHandler) Handle(ctx context.Context, query *GetProductByIdQuery) (*dto.ProductResponse, error) {
	res, err := q.rsClient.GetProductById(ctx, &readerService.GetProductByIdReq{ProductID: query.ProductID})
	if err != nil {
		return nil, err
	}

	return dto.ProductResponseFromGrpc(res.GetProduct()), err
}
