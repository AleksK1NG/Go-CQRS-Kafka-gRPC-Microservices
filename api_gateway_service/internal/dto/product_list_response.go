package dto

import readerService "github.com/AleksK1NG/cqrs-microservices/reader_service/proto/product_reader"

type ProductsListResponse struct {
	TotalCount int64              `json:"totalCount" bson:"totalCount"`
	TotalPages int64              `json:"totalPages" bson:"totalPages"`
	Page       int64              `json:"page" bson:"page"`
	Size       int64              `json:"size" bson:"size"`
	HasMore    bool               `json:"hasMore" bson:"hasMore"`
	Products   []*ProductResponse `json:"products" bson:"products"`
}

func ProductsListResponseFromGrpc(listResponse *readerService.SearchRes) *ProductsListResponse {
	list := make([]*ProductResponse, 0, len(listResponse.GetProducts()))
	for _, product := range listResponse.GetProducts() {
		list = append(list, ProductResponseFromGrpc(product))
	}

	return &ProductsListResponse{
		TotalCount: listResponse.GetTotalCount(),
		TotalPages: listResponse.GetTotalPages(),
		Page:       listResponse.GetPage(),
		Size:       listResponse.GetSize(),
		HasMore:    listResponse.GetHasMore(),
		Products:   list,
	}
}
