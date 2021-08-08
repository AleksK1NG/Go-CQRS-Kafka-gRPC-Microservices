package commands

import (
	"github.com/AleksK1NG/cqrs-microservices/api_gateway_service/internal/dto"
	uuid "github.com/satori/go.uuid"
)

type ProductCommands struct {
	CreateProduct CreateProductCmdHandler
	UpdateProduct UpdateProductCmdHandler
	DeleteProduct DeleteProductCmdHandler
}

func NewProductCommands(createProduct CreateProductCmdHandler, updateProduct UpdateProductCmdHandler, deleteProduct DeleteProductCmdHandler) *ProductCommands {
	return &ProductCommands{CreateProduct: createProduct, UpdateProduct: updateProduct, DeleteProduct: deleteProduct}
}

type CreateProductCommand struct {
	CreateDto *dto.CreateProductDto
}

func NewCreateProductCommand(createDto *dto.CreateProductDto) *CreateProductCommand {
	return &CreateProductCommand{CreateDto: createDto}
}

type UpdateProductCommand struct {
	UpdateDto *dto.UpdateProductDto
}

func NewUpdateProductCommand(updateDto *dto.UpdateProductDto) *UpdateProductCommand {
	return &UpdateProductCommand{UpdateDto: updateDto}
}

type DeleteProductCommand struct {
	ProductID uuid.UUID `json:"productId" validate:"required"`
}

func NewDeleteProductCommand(productID uuid.UUID) *DeleteProductCommand {
	return &DeleteProductCommand{ProductID: productID}
}
