package queries

import "github.com/AleksK1NG/cqrs-microservices/pkg/utils"

type ProductQueries struct {
	GetProductById GetProductByIdHandler
	SearchProduct  SearchProductHandler
}

func NewProductQueries(getProductById GetProductByIdHandler, searchProduct SearchProductHandler) *ProductQueries {
	return &ProductQueries{GetProductById: getProductById, SearchProduct: searchProduct}
}

type GetProductByIdQuery struct {
	ProductID string `json:"productId" validate:"required,gte=0,lte=255"`
}

func NewGetProductByIdQuery(productID string) *GetProductByIdQuery {
	return &GetProductByIdQuery{ProductID: productID}
}

type SearchProductQuery struct {
	Text       string            `json:"text"`
	Pagination *utils.Pagination `json:"pagination"`
}

func NewSearchProductQuery(text string, pagination *utils.Pagination) *SearchProductQuery {
	return &SearchProductQuery{Text: text, Pagination: pagination}
}
