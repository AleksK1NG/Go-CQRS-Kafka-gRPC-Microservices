package products

import "github.com/labstack/echo/v4"

type HttpDelivery interface {
	CreateProduct() echo.HandlerFunc
	UpdateProduct() echo.HandlerFunc
	DeleteProduct() echo.HandlerFunc

	GetProductByID() echo.HandlerFunc
	SearchProduct() echo.HandlerFunc
}
