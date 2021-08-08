package repository

import (
	"context"
	"github.com/AleksK1NG/cqrs-microservices/product_writer_service/internal/models"
	uuid "github.com/satori/go.uuid"
)

type Repository interface {
	CreateProduct(ctx context.Context, product *models.Product) (*models.Product, error)
	UpdateProduct(ctx context.Context, product *models.Product) (*models.Product, error)
	DeleteProductByID(ctx context.Context, uuid uuid.UUID) error

	GetProductById(ctx context.Context, uuid uuid.UUID) (*models.Product, error)
}
