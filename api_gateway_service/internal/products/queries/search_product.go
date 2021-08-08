package queries

import (
	"context"
	"github.com/AleksK1NG/cqrs-microservices/api_gateway_service/config"
	"github.com/AleksK1NG/cqrs-microservices/api_gateway_service/internal/dto"
	"github.com/AleksK1NG/cqrs-microservices/pkg/logger"
	"github.com/AleksK1NG/cqrs-microservices/pkg/tracing"
	readerService "github.com/AleksK1NG/cqrs-microservices/product_reader_service/proto/product_reader"
	"github.com/opentracing/opentracing-go"
)

type SearchProductHandler interface {
	Handle(ctx context.Context, query *SearchProductQuery) (*dto.ProductsListResponse, error)
}

type searchProductHandler struct {
	log      logger.Logger
	cfg      *config.Config
	rsClient readerService.ReaderServiceClient
}

func NewSearchProductHandler(log logger.Logger, cfg *config.Config, rsClient readerService.ReaderServiceClient) *searchProductHandler {
	return &searchProductHandler{log: log, cfg: cfg, rsClient: rsClient}
}

func (s *searchProductHandler) Handle(ctx context.Context, query *SearchProductQuery) (*dto.ProductsListResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "searchProductHandler.Handle")
	defer span.Finish()

	ctx = tracing.InjectTextMapCarrierToGrpcMetaData(ctx, span.Context())
	res, err := s.rsClient.SearchProduct(ctx, &readerService.SearchReq{
		Search: query.Text,
		Page:   int64(query.Pagination.GetPage()),
		Size:   int64(query.Pagination.GetSize()),
	})
	if err != nil {
		return nil, err
	}

	return dto.ProductsListResponseFromGrpc(res), nil
}
