package commands

import "github.com/AleksK1NG/cqrs-microservices/api_gateway_service/internal/dto"

type ProductCommands struct {
	CreateProduct CreateProductCmdHandler
	UpdateProduct UpdateProductCmdHandler
}

func NewProductCommands(createProduct CreateProductCmdHandler, updateProduct UpdateProductCmdHandler) *ProductCommands {
	return &ProductCommands{CreateProduct: createProduct, UpdateProduct: updateProduct}
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
