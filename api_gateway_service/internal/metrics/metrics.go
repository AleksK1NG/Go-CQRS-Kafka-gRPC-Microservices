package metrics

import (
	"fmt"
	"github.com/AleksK1NG/cqrs-microservices/api_gateway_service/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type ApiGatewayMetrics struct {
	SuccessHttpRequests        prometheus.Counter
	ErrorHttpRequests          prometheus.Counter
	CreateProductHttpRequests  prometheus.Counter
	UpdateProductHttpRequests  prometheus.Counter
	DeleteProductHttpRequests  prometheus.Counter
	GetProductByIdHttpRequests prometheus.Counter
	SearchProductHttpRequests  prometheus.Counter
}

func NewApiGatewayMetrics(cfg *config.Config) *ApiGatewayMetrics {
	return &ApiGatewayMetrics{
		SuccessHttpRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_success_http_requests_total", cfg.ServiceName),
			Help: "The total number of success http requests",
		}),
		ErrorHttpRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_error_http_requests_total", cfg.ServiceName),
			Help: "The total number of error http requests",
		}),
		CreateProductHttpRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_create_product_http_requests_total", cfg.ServiceName),
			Help: "The total number of create product http requests",
		}),
		UpdateProductHttpRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_update_product_http_requests_total", cfg.ServiceName),
			Help: "The total number of update product http requests",
		}),
		DeleteProductHttpRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_delete_product_http_requests_total", cfg.ServiceName),
			Help: "The total number of delete product http requests",
		}),
		GetProductByIdHttpRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_get_product_by_id_http_requests_total", cfg.ServiceName),
			Help: "The total number of get product by id http requests",
		}),
		SearchProductHttpRequests: promauto.NewCounter(prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_search_product_http_requests_total", cfg.ServiceName),
			Help: "The total number of search product http requests",
		}),
	}
}
