package repository

import (
	"context"
	"github.com/AleksK1NG/cqrs-microservices/pkg/logger"
	"github.com/AleksK1NG/cqrs-microservices/product_writer_service/config"
	"github.com/AleksK1NG/cqrs-microservices/product_writer_service/internal/models"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

type productRepository struct {
	log logger.Logger
	cfg *config.Config
	db  *pgxpool.Pool
}

func NewProductRepository(log logger.Logger, cfg *config.Config, db *pgxpool.Pool) *productRepository {
	return &productRepository{log: log, cfg: cfg, db: db}
}

func (p *productRepository) CreateProduct(ctx context.Context, product *models.Product) (*models.Product, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "productRepository.CreateProduct")
	defer span.Finish()

	var created models.Product
	if err := p.db.QueryRow(ctx, createProductQuery, &product.ProductID, &product.Name, &product.Description, &product.Price).Scan(
		&created.ProductID,
		&created.Name,
		&created.Description,
		&created.Price,
		&created.CreatedAt,
		&created.UpdatedAt,
	); err != nil {
		return nil, errors.Wrap(err, "db.QueryRow")
	}

	p.log.Debugf("created product: %+v", created)
	return &created, nil
}

func (p *productRepository) UpdateProduct(ctx context.Context, product *models.Product) (*models.Product, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "productRepository.UpdateProduct")
	defer span.Finish()

	var prod models.Product
	if err := p.db.QueryRow(
		ctx,
		updateProductQuery,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.ProductID,
	).Scan(&prod.ProductID, &prod.Name, &prod.Description, &prod.Price, &prod.CreatedAt, &prod.UpdatedAt); err != nil {
		return nil, errors.Wrap(err, "Scan")
	}

	p.log.Debugf("updated product: %+v", prod)
	return &prod, nil
}

func (p *productRepository) GetProductById(ctx context.Context, uuid uuid.UUID) (*models.Product, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "productRepository.GetProductById")
	defer span.Finish()

	var product models.Product
	if err := p.db.QueryRow(ctx, getProductByIdQuery, uuid).Scan(
		&product.ProductID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.CreatedAt,
		&product.UpdatedAt,
	); err != nil {
		return nil, errors.Wrap(err, "Scan")
	}

	p.log.Debugf("get product by id: %+v", product)
	return &product, nil
}

func (p *productRepository) DeleteProductByID(ctx context.Context, uuid uuid.UUID) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "productRepository.DeleteProductByID")
	defer span.Finish()

	deleteProductByIdQuery := `DELETE FROM products WHERE product_id = $1`
	result, err := p.db.Exec(ctx, deleteProductByIdQuery, uuid)
	if err != nil {
		return errors.Wrap(err, "Exec")
	}

	p.log.Infof("deleted rows: %v", result.RowsAffected())
	return nil
}
