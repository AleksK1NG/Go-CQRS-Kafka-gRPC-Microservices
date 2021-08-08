package v1

import (
	"github.com/AleksK1NG/cqrs-microservices/api_gateway_service/config"
	"github.com/AleksK1NG/cqrs-microservices/api_gateway_service/internal/dto"
	"github.com/AleksK1NG/cqrs-microservices/api_gateway_service/internal/metrics"
	"github.com/AleksK1NG/cqrs-microservices/api_gateway_service/internal/middlewares"
	"github.com/AleksK1NG/cqrs-microservices/api_gateway_service/internal/products/commands"
	"github.com/AleksK1NG/cqrs-microservices/api_gateway_service/internal/products/queries"
	"github.com/AleksK1NG/cqrs-microservices/api_gateway_service/internal/products/service"
	httpErrors "github.com/AleksK1NG/cqrs-microservices/pkg/http_errors"
	"github.com/AleksK1NG/cqrs-microservices/pkg/logger"
	"github.com/AleksK1NG/cqrs-microservices/pkg/tracing"
	"github.com/AleksK1NG/cqrs-microservices/pkg/utils"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

type productsHandlers struct {
	group     *echo.Group
	log       logger.Logger
	mw        middlewares.MiddlewareManager
	cfg       *config.Config
	ps        *service.ProductService
	v         *validator.Validate
	agMetrics *metrics.ApiGatewayMetrics
}

func NewProductsHandlers(
	group *echo.Group,
	log logger.Logger,
	mw middlewares.MiddlewareManager,
	cfg *config.Config,
	ps *service.ProductService,
	v *validator.Validate,
	agMetrics *metrics.ApiGatewayMetrics,
) *productsHandlers {
	return &productsHandlers{group: group, log: log, mw: mw, cfg: cfg, ps: ps, v: v, agMetrics: agMetrics}
}

func (h *productsHandlers) CreateProduct() echo.HandlerFunc {
	return func(c echo.Context) error {
		h.agMetrics.CreateProductHttpRequests.Inc()

		ctx, span := tracing.StartHttpServerTracerSpan(c, "productsHandlers.CreateProduct")
		defer span.Finish()

		createDto := &dto.CreateProductDto{}
		if err := c.Bind(createDto); err != nil {
			h.log.WarnMsg("Bind", err)
			h.agMetrics.ErrorHttpRequests.Inc()
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		createDto.ProductID = uuid.NewV4()
		if err := h.v.StructCtx(ctx, createDto); err != nil {
			h.log.WarnMsg("validate", err)
			h.agMetrics.ErrorHttpRequests.Inc()
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		if err := h.ps.Commands.CreateProduct.Handle(ctx, commands.NewCreateProductCommand(createDto)); err != nil {
			h.log.WarnMsg("CreateProduct", err)
			h.agMetrics.ErrorHttpRequests.Inc()
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		h.agMetrics.SuccessHttpRequests.Inc()
		return c.JSON(http.StatusCreated, dto.CreateProductResponseDto{ProductID: createDto.ProductID})
	}
}

func (h *productsHandlers) GetProductByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		h.agMetrics.GetProductByIdHttpRequests.Inc()

		ctx, span := tracing.StartHttpServerTracerSpan(c, "productsHandlers.GetProductByID")
		defer span.Finish()

		query := queries.NewGetProductByIdQuery(c.Param("id"))
		if err := h.v.StructCtx(ctx, query); err != nil {
			h.log.WarnMsg("validate", err)
			h.agMetrics.ErrorHttpRequests.Inc()
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		response, err := h.ps.Queries.GetProductById.Handle(ctx, query)
		if err != nil {
			h.log.WarnMsg("GetProductById", err)
			h.agMetrics.ErrorHttpRequests.Inc()
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		h.agMetrics.SuccessHttpRequests.Inc()
		return c.JSON(http.StatusOK, response)
	}
}

func (h *productsHandlers) SearchProduct() echo.HandlerFunc {
	return func(c echo.Context) error {
		h.agMetrics.SearchProductHttpRequests.Inc()

		ctx, span := tracing.StartHttpServerTracerSpan(c, "productsHandlers.SearchProduct")
		defer span.Finish()

		pq := utils.NewPaginationFromQueryParams(c.QueryParam("size"), c.QueryParam("page"))

		query := queries.NewSearchProductQuery(c.QueryParam("search"), pq)
		response, err := h.ps.Queries.SearchProduct.Handle(ctx, query)
		if err != nil {
			h.log.WarnMsg("SearchProduct", err)
			h.agMetrics.ErrorHttpRequests.Inc()
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		h.agMetrics.SuccessHttpRequests.Inc()
		return c.JSON(http.StatusOK, response)
	}
}

func (h *productsHandlers) UpdateProduct() echo.HandlerFunc {
	return func(c echo.Context) error {
		h.agMetrics.UpdateProductHttpRequests.Inc()

		ctx, span := tracing.StartHttpServerTracerSpan(c, "productsHandlers.UpdateProduct")
		defer span.Finish()

		productUUID, err := uuid.FromString(c.Param("id"))
		if err != nil {
			h.log.WarnMsg("uuid.FromString", err)
			h.agMetrics.ErrorHttpRequests.Inc()
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		updateDto := &dto.UpdateProductDto{ProductID: productUUID}
		if err := c.Bind(updateDto); err != nil {
			h.log.WarnMsg("Bind", err)
			h.agMetrics.ErrorHttpRequests.Inc()
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		if err := h.v.StructCtx(ctx, updateDto); err != nil {
			h.log.WarnMsg("validate", err)
			h.agMetrics.ErrorHttpRequests.Inc()
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		if err := h.ps.Commands.UpdateProduct.Handle(ctx, commands.NewUpdateProductCommand(updateDto)); err != nil {
			h.log.WarnMsg("UpdateProduct", err)
			h.agMetrics.ErrorHttpRequests.Inc()
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		h.agMetrics.SuccessHttpRequests.Inc()
		return c.JSON(http.StatusOK, updateDto)
	}
}

func (h *productsHandlers) DeleteProduct() echo.HandlerFunc {
	return func(c echo.Context) error {
		h.agMetrics.DeleteProductHttpRequests.Inc()

		ctx, span := tracing.StartHttpServerTracerSpan(c, "productsHandlers.DeleteProduct")
		defer span.Finish()

		productUUID, err := uuid.FromString(c.Param("id"))
		if err != nil {
			h.log.WarnMsg("uuid.FromString", err)
			h.agMetrics.ErrorHttpRequests.Inc()
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		if err := h.ps.Commands.DeleteProduct.Handle(ctx, commands.NewDeleteProductCommand(productUUID)); err != nil {
			h.log.WarnMsg("DeleteProduct", err)
			h.agMetrics.ErrorHttpRequests.Inc()
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		h.agMetrics.SuccessHttpRequests.Inc()
		return c.NoContent(http.StatusOK)
	}
}
