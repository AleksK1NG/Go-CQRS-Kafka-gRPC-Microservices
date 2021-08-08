package grpc

import (
	"context"
	"github.com/AleksK1NG/cqrs-microservices/pkg/logger"
	"github.com/AleksK1NG/cqrs-microservices/pkg/tracing"
	"github.com/AleksK1NG/cqrs-microservices/pkg/utils"
	"github.com/AleksK1NG/cqrs-microservices/product_reader_service/config"
	"github.com/AleksK1NG/cqrs-microservices/product_reader_service/internal/models"
	"github.com/AleksK1NG/cqrs-microservices/product_reader_service/internal/product/commands"
	"github.com/AleksK1NG/cqrs-microservices/product_reader_service/internal/product/queries"
	"github.com/AleksK1NG/cqrs-microservices/product_reader_service/internal/product/service"
	"github.com/AleksK1NG/cqrs-microservices/product_reader_service/proto/product_reader"
	"github.com/go-playground/validator"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type grpcService struct {
	log logger.Logger
	cfg *config.Config
	v   *validator.Validate
	ps  *service.ProductService
}

func NewReaderGrpcService(log logger.Logger, cfg *config.Config, v *validator.Validate, ps *service.ProductService) *grpcService {
	return &grpcService{log: log, cfg: cfg, v: v, ps: ps}
}

func (s *grpcService) CreateProduct(ctx context.Context, req *readerService.CreateProductReq) (*readerService.CreateProductRes, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "grpcService.CreateProduct")
	defer span.Finish()

	command := commands.NewCreateProductCommand(req.GetProductID(), req.GetName(), req.GetDescription(), req.GetPrice(), time.Now(), time.Now())
	if err := s.v.StructCtx(ctx, command); err != nil {
		s.log.WarnMsg("validate", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}

	if err := s.ps.Commands.CreateProduct.Handle(ctx, command); err != nil {
		s.log.WarnMsg("CreateProduct", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}

	return &readerService.CreateProductRes{ProductID: req.GetProductID()}, nil
}

func (s *grpcService) UpdateProduct(ctx context.Context, req *readerService.UpdateProductReq) (*readerService.UpdateProductRes, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "grpcService.UpdateProduct")
	defer span.Finish()

	command := commands.NewUpdateProductCommand(req.GetProductID(), req.GetName(), req.GetDescription(), req.GetPrice(), time.Now())
	if err := s.v.StructCtx(ctx, command); err != nil {
		s.log.WarnMsg("validate", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}

	if err := s.ps.Commands.UpdateProduct.Handle(ctx, command); err != nil {
		s.log.WarnMsg("UpdateProduct", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}

	return &readerService.UpdateProductRes{ProductID: req.GetProductID()}, nil
}

func (s *grpcService) GetProductById(ctx context.Context, req *readerService.GetProductByIdReq) (*readerService.GetProductByIdRes, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "grpcService.GetProductById")
	defer span.Finish()

	productUUID, err := uuid.FromString(req.GetProductID())
	if err != nil {
		s.log.WarnMsg("uuid.FromString", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}

	query := queries.NewGetProductByIdQuery(productUUID)
	if err := s.v.StructCtx(ctx, query); err != nil {
		s.log.WarnMsg("validate", err)
		return nil, s.errResponse(codes.InvalidArgument, err)
	}

	product, err := s.ps.Queries.GetProductById.Handle(ctx, query)
	if err != nil {
		s.log.WarnMsg("GetProductById", err)
		return nil, s.errResponse(codes.Internal, err)
	}

	return &readerService.GetProductByIdRes{Product: models.ProductToGrpcMessage(product)}, nil
}

func (s *grpcService) SearchProduct(ctx context.Context, req *readerService.SearchReq) (*readerService.SearchRes, error) {
	ctx, span := tracing.StartGrpcServerTracerSpan(ctx, "grpcService.SearchProduct")
	defer span.Finish()

	pq := utils.NewPaginationQuery(int(req.GetSize()), int(req.GetPage()))

	query := queries.NewSearchProductQuery(req.GetSearch(), pq)
	productsList, err := s.ps.Queries.SearchProduct.Handle(ctx, query)
	if err != nil {
		s.log.WarnMsg("SearchProduct", err)
		return nil, s.errResponse(codes.Internal, err)
	}

	return models.ProductListToGrpc(productsList), nil
}

func (s *grpcService) errResponse(c codes.Code, err error) error {
	return status.Error(c, err.Error())
}
