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
	"github.com/opentracing/opentracing-go"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

type productsHandlers struct {
	group   *echo.Group
	log     logger.Logger
	mw      middlewares.MiddlewareManager
	cfg     *config.Config
	ps      *service.ProductService
	v       *validator.Validate
	metrics *metrics.ApiGatewayMetrics
}

func NewProductsHandlers(
	group *echo.Group,
	log logger.Logger,
	mw middlewares.MiddlewareManager,
	cfg *config.Config,
	ps *service.ProductService,
	v *validator.Validate,
	metrics *metrics.ApiGatewayMetrics,
) *productsHandlers {
	return &productsHandlers{group: group, log: log, mw: mw, cfg: cfg, ps: ps, v: v, metrics: metrics}
}

// CreateProduct
// @Tags Products
// @Summary Create product
// @Description Create new product item
// @Accept json
// @Produce json
// @Success 201 {object} dto.CreateProductResponseDto
// @Router /products [post]
func (h *productsHandlers) CreateProduct() echo.HandlerFunc {
	return func(c echo.Context) error {
		h.metrics.CreateProductHttpRequests.Inc()

		ctx, span := tracing.StartHttpServerTracerSpan(c, "productsHandlers.CreateProduct")
		defer span.Finish()

		createDto := &dto.CreateProductDto{}
		if err := c.Bind(createDto); err != nil {
			h.log.WarnMsg("Bind", err)
			h.traceErr(span, err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		createDto.ProductID = uuid.NewV4()
		if err := h.v.StructCtx(ctx, createDto); err != nil {
			h.log.WarnMsg("validate", err)
			h.traceErr(span, err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		if err := h.ps.Commands.CreateProduct.Handle(ctx, commands.NewCreateProductCommand(createDto)); err != nil {
			h.log.WarnMsg("CreateProduct", err)
			h.metrics.ErrorHttpRequests.Inc()
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		h.metrics.SuccessHttpRequests.Inc()
		return c.JSON(http.StatusCreated, dto.CreateProductResponseDto{ProductID: createDto.ProductID})
	}
}

// GetProductByID
// @Tags Products
// @Summary Get product
// @Description Get product by id
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} dto.ProductResponse
// @Router /products/{id} [get]
func (h *productsHandlers) GetProductByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		h.metrics.GetProductByIdHttpRequests.Inc()

		ctx, span := tracing.StartHttpServerTracerSpan(c, "productsHandlers.GetProductByID")
		defer span.Finish()

		productUUID, err := uuid.FromString(c.Param("id"))
		if err != nil {
			h.log.WarnMsg("uuid.FromString", err)
			h.traceErr(span, err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		query := queries.NewGetProductByIdQuery(productUUID)
		response, err := h.ps.Queries.GetProductById.Handle(ctx, query)
		if err != nil {
			h.log.WarnMsg("GetProductById", err)
			h.metrics.ErrorHttpRequests.Inc()
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		h.metrics.SuccessHttpRequests.Inc()
		return c.JSON(http.StatusOK, response)
	}
}

// SearchProduct
// @Tags Products
// @Summary Search product
// @Description Get product by name with pagination
// @Accept json
// @Produce json
// @Param search query string false "search text"
// @Param page query string false "page number"
// @Param size query string false "number of elements"
// @Success 200 {object} dto.ProductsListResponse
// @Router /products/search [get]
func (h *productsHandlers) SearchProduct() echo.HandlerFunc {
	return func(c echo.Context) error {
		h.metrics.SearchProductHttpRequests.Inc()

		ctx, span := tracing.StartHttpServerTracerSpan(c, "productsHandlers.SearchProduct")
		defer span.Finish()

		pq := utils.NewPaginationFromQueryParams(c.QueryParam("size"), c.QueryParam("page"))

		query := queries.NewSearchProductQuery(c.QueryParam("search"), pq)
		response, err := h.ps.Queries.SearchProduct.Handle(ctx, query)
		if err != nil {
			h.log.WarnMsg("SearchProduct", err)
			h.metrics.ErrorHttpRequests.Inc()
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		h.metrics.SuccessHttpRequests.Inc()
		return c.JSON(http.StatusOK, response)
	}
}

// UpdateProduct
// @Tags Products
// @Summary Update product
// @Description Update existing product
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} dto.UpdateProductDto
// @Router /products/{id} [put]
func (h *productsHandlers) UpdateProduct() echo.HandlerFunc {
	return func(c echo.Context) error {
		h.metrics.UpdateProductHttpRequests.Inc()

		ctx, span := tracing.StartHttpServerTracerSpan(c, "productsHandlers.UpdateProduct")
		defer span.Finish()

		productUUID, err := uuid.FromString(c.Param("id"))
		if err != nil {
			h.log.WarnMsg("uuid.FromString", err)
			h.traceErr(span, err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		updateDto := &dto.UpdateProductDto{ProductID: productUUID}
		if err := c.Bind(updateDto); err != nil {
			h.log.WarnMsg("Bind", err)
			h.traceErr(span, err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		if err := h.v.StructCtx(ctx, updateDto); err != nil {
			h.log.WarnMsg("validate", err)
			h.traceErr(span, err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		if err := h.ps.Commands.UpdateProduct.Handle(ctx, commands.NewUpdateProductCommand(updateDto)); err != nil {
			h.log.WarnMsg("UpdateProduct", err)
			h.metrics.ErrorHttpRequests.Inc()
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		h.metrics.SuccessHttpRequests.Inc()
		return c.JSON(http.StatusOK, updateDto)
	}
}

// DeleteProduct
// @Tags Products
// @Summary Delete product
// @Description Delete existing product
// @Accept json
// @Produce json
// @Success 200 ""
// @Param id path string true "Product ID"
// @Router /products/{id} [delete]
func (h *productsHandlers) DeleteProduct() echo.HandlerFunc {
	return func(c echo.Context) error {
		h.metrics.DeleteProductHttpRequests.Inc()

		ctx, span := tracing.StartHttpServerTracerSpan(c, "productsHandlers.DeleteProduct")
		defer span.Finish()

		productUUID, err := uuid.FromString(c.Param("id"))
		if err != nil {
			h.log.WarnMsg("uuid.FromString", err)
			h.traceErr(span, err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		if err := h.ps.Commands.DeleteProduct.Handle(ctx, commands.NewDeleteProductCommand(productUUID)); err != nil {
			h.log.WarnMsg("DeleteProduct", err)
			h.metrics.ErrorHttpRequests.Inc()
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		h.metrics.SuccessHttpRequests.Inc()
		return c.NoContent(http.StatusOK)
	}
}

func (h *productsHandlers) traceErr(span opentracing.Span, err error) {
	span.SetTag("error", true)
	span.LogKV("error_code", err.Error())
	h.metrics.ErrorHttpRequests.Inc()
}
