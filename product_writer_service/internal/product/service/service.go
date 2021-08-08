package service

import (
	kafkaClient "github.com/AleksK1NG/cqrs-microservices/pkg/kafka"
	"github.com/AleksK1NG/cqrs-microservices/pkg/logger"
	"github.com/AleksK1NG/cqrs-microservices/product_writer_service/config"
	"github.com/AleksK1NG/cqrs-microservices/product_writer_service/internal/product/commands"
	"github.com/AleksK1NG/cqrs-microservices/product_writer_service/internal/product/queries"
	"github.com/AleksK1NG/cqrs-microservices/product_writer_service/internal/product/repository"
)

type ProductService struct {
	Commands *commands.ProductCommands
	Queries  *queries.ProductQueries
}

func NewProductService(log logger.Logger, cfg *config.Config, pgRepo repository.Repository, kafkaProducer kafkaClient.Producer) *ProductService {

	updateProductHandler := commands.NewUpdateProductHandler(log, cfg, pgRepo, kafkaProducer)
	createProductHandler := commands.NewCreateProductHandler(log, cfg, pgRepo, kafkaProducer)

	getProductByIdHandler := queries.NewGetProductByIdHandler(log, cfg, pgRepo)

	productCommands := commands.NewProductCommands(createProductHandler, updateProductHandler)
	productQueries := queries.NewProductQueries(getProductByIdHandler)

	return &ProductService{Commands: productCommands, Queries: productQueries}
}
