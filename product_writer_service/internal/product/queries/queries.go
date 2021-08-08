package queries

import uuid "github.com/satori/go.uuid"

type ProductQueries struct {
	GetProductById GetProductByIdHandler
}

func NewProductQueries(getProductById GetProductByIdHandler) *ProductQueries {
	return &ProductQueries{GetProductById: getProductById}
}

type GetProductByIdQuery struct {
	ProductID uuid.UUID `json:"productId" validate:"required,gte=0,lte=255"`
}

func NewGetProductByIdQuery(productID uuid.UUID) *GetProductByIdQuery {
	return &GetProductByIdQuery{ProductID: productID}
}
