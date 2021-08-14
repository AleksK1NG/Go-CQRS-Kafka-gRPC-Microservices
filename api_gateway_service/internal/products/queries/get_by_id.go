package queries

import (
	"context"
	"github.com/AleksK1NG/cqrs-microservices/api_gateway_service/config"
	"github.com/AleksK1NG/cqrs-microservices/api_gateway_service/internal/dto"
	"github.com/AleksK1NG/cqrs-microservices/pkg/logger"
	"github.com/AleksK1NG/cqrs-microservices/pkg/tracing"
	readerService "github.com/AleksK1NG/cqrs-microservices/reader_service/proto/product_reader"
	"github.com/opentracing/opentracing-go"
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "getProductByIdHandler.Handle")
	defer span.Finish()

	ctx = tracing.InjectTextMapCarrierToGrpcMetaData(ctx, span.Context())
	res, err := q.rsClient.GetProductById(ctx, &readerService.GetProductByIdReq{ProductID: query.ProductID.String()})
	if err != nil {
		return nil, err
	}

	return dto.ProductResponseFromGrpc(res.GetProduct()), nil
}
