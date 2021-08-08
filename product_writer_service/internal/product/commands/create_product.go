package commands

import (
	"context"
	kafkaClient "github.com/AleksK1NG/cqrs-microservices/pkg/kafka"
	"github.com/AleksK1NG/cqrs-microservices/pkg/logger"
	"github.com/AleksK1NG/cqrs-microservices/product_writer_service/config"
	"github.com/AleksK1NG/cqrs-microservices/product_writer_service/internal/models"
	"github.com/AleksK1NG/cqrs-microservices/product_writer_service/internal/product/repository"
	"github.com/AleksK1NG/cqrs-microservices/product_writer_service/mappers"
	kafkaMessages "github.com/AleksK1NG/cqrs-microservices/proto/kafka"
	"github.com/golang/protobuf/proto"
	"github.com/segmentio/kafka-go"
	"time"
)

type CreateProductCmdHandler interface {
	Handle(ctx context.Context, command *CreateProductCommand) (*models.Product, error)
}

type createProductHandler struct {
	log           logger.Logger
	cfg           *config.Config
	pgRepo        repository.Repository
	kafkaProducer kafkaClient.Producer
}

func NewCreateProductHandler(log logger.Logger, cfg *config.Config, pgRepo repository.Repository, kafkaProducer kafkaClient.Producer) *createProductHandler {
	return &createProductHandler{log: log, cfg: cfg, pgRepo: pgRepo, kafkaProducer: kafkaProducer}
}

func (c *createProductHandler) Handle(ctx context.Context, command *CreateProductCommand) (*models.Product, error) {
	productDto := &models.Product{
		ProductID:   command.ProductID,
		Name:        command.Name,
		Description: command.Description,
		Price:       command.Price,
	}

	createProduct, err := c.pgRepo.CreateProduct(ctx, productDto)
	if err != nil {
		return nil, err
	}

	msg := &kafkaMessages.ProductCreated{Product: mappers.ProductToGrpcMessage(createProduct)}
	msgBytes, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}

	message := kafka.Message{
		Topic: c.cfg.KafkaTopics.ProductCreated.TopicName,
		Value: msgBytes,
		Time:  time.Now().UTC(),
	}

	c.log.Infof("created product: %+v", createProduct)
	if err := c.kafkaProducer.PublishMessage(ctx, message); err != nil {
		return nil, err
	}

	return createProduct, nil
}
