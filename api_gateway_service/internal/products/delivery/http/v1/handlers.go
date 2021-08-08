package v1

import (
	"github.com/AleksK1NG/cqrs-microservices/api_gateway_service/config"
	"github.com/AleksK1NG/cqrs-microservices/api_gateway_service/internal/dto"
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
	"strconv"
)

type productsHandlers struct {
	group *echo.Group
	log   logger.Logger
	mw    middlewares.MiddlewareManager
	cfg   *config.Config
	ps    *service.ProductService
	v     *validator.Validate
}

func NewProductsHandlers(
	group *echo.Group,
	log logger.Logger,
	mw middlewares.MiddlewareManager,
	cfg *config.Config,
	ps *service.ProductService,
	v *validator.Validate,
) *productsHandlers {
	return &productsHandlers{group: group, log: log, mw: mw, cfg: cfg, ps: ps, v: v}
}

func (h *productsHandlers) CreateProduct() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, span := tracing.StartHttpServerTracerSpan(c, "productsHandlers.CreateProduct")
		defer span.Finish()

		createDto := &dto.CreateProductDto{ProductID: uuid.NewV4()}
		if err := c.Bind(createDto); err != nil {
			h.log.WarnMsg("Bind", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		if err := h.v.StructCtx(ctx, createDto); err != nil {
			h.log.WarnMsg("validate", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		if err := h.ps.Commands.CreateProduct.Handle(ctx, commands.NewCreateProductCommand(createDto)); err != nil {
			h.log.WarnMsg("CreateProduct", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		return c.JSON(http.StatusCreated, dto.CreateProductResponseDto{ProductID: createDto.ProductID})
	}
}

func (h *productsHandlers) GetProductByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, span := tracing.StartHttpServerTracerSpan(c, "productsHandlers.GetProductByID")
		defer span.Finish()

		query := queries.NewGetProductByIdQuery(c.Param("id"))
		if err := h.v.StructCtx(ctx, query); err != nil {
			h.log.WarnMsg("validate", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		response, err := h.ps.Queries.GetProductById.Handle(ctx, query)
		if err != nil {
			h.log.WarnMsg("GetProductById", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		return c.JSON(http.StatusOK, response)
	}
}

func (h *productsHandlers) SearchProduct() echo.HandlerFunc {
	return func(c echo.Context) error {

		page, err := strconv.Atoi(c.QueryParam("page"))
		if err != nil {
			h.log.WarnMsg("validate", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}
		size, err := strconv.Atoi(c.QueryParam("size"))
		if err != nil {
			h.log.WarnMsg("validate", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		query := queries.NewSearchProductQuery(c.QueryParam("search"), utils.NewPaginationQuery(size, page))
		response, err := h.ps.Queries.SearchProduct.Handle(c.Request().Context(), query)
		if err != nil {
			h.log.WarnMsg("SearchProduct", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		return c.JSON(http.StatusOK, response)
	}
}

func (h *productsHandlers) UpdateProduct() echo.HandlerFunc {
	return func(c echo.Context) error {

		productUUID, err := uuid.FromString(c.Param("id"))
		if err != nil {
			h.log.WarnMsg("uuid.FromString", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		updateDto := &dto.UpdateProductDto{ProductID: productUUID}
		if err := c.Bind(updateDto); err != nil {
			h.log.WarnMsg("Bind", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		if err := h.v.StructCtx(c.Request().Context(), updateDto); err != nil {
			h.log.WarnMsg("validate", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		if err := h.ps.Commands.UpdateProduct.Handle(c.Request().Context(), commands.NewUpdateProductCommand(updateDto)); err != nil {
			h.log.WarnMsg("UpdateProduct", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		return c.JSON(http.StatusOK, updateDto)
	}
}

func (h *productsHandlers) DeleteProduct() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, "OK")
	}
}
