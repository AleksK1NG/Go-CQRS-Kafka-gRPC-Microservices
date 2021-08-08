package queries

import (
	"context"
	"github.com/AleksK1NG/cqrs-microservices/api_gateway_service/config"
	"github.com/AleksK1NG/cqrs-microservices/api_gateway_service/internal/dto"
	"github.com/AleksK1NG/cqrs-microservices/pkg/logger"
	"github.com/AleksK1NG/cqrs-microservices/pkg/tracing"
	readerService "github.com/AleksK1NG/cqrs-microservices/product_reader_service/proto/product_reader"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc/metadata"
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

	textMapCarrier, err := tracing.InjectTextMapCarrier(span.Context())
	if err != nil {
		q.log.WarnMsg("InjectTextMapCarrier", err)
	}

	q.log.Infof("text map carrier: %+v", textMapCarrier)
	md := metadata.New(textMapCarrier)
	ctx = metadata.NewOutgoingContext(ctx, md)

	res, err := q.rsClient.GetProductById(ctx, &readerService.GetProductByIdReq{ProductID: query.ProductID})
	if err != nil {
		return nil, err
	}

	return dto.ProductResponseFromGrpc(res.GetProduct()), err
}
