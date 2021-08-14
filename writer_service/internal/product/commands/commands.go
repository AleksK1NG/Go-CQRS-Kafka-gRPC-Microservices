package commands

import uuid "github.com/satori/go.uuid"

type ProductCommands struct {
	CreateProduct CreateProductCmdHandler
	UpdateProduct UpdateProductCmdHandler
	DeleteProduct DeleteProductCmdHandler
}

func NewProductCommands(createProduct CreateProductCmdHandler, updateProduct UpdateProductCmdHandler, deleteProduct DeleteProductCmdHandler) *ProductCommands {
	return &ProductCommands{CreateProduct: createProduct, UpdateProduct: updateProduct, DeleteProduct: deleteProduct}
}

type CreateProductCommand struct {
	ProductID   uuid.UUID `json:"productId" validate:"required"`
	Name        string    `json:"name" validate:"required,gte=0,lte=255"`
	Description string    `json:"description" validate:"required,gte=0,lte=5000"`
	Price       float64   `json:"price" validate:"required,gte=0"`
}

func NewCreateProductCommand(productID uuid.UUID, name string, description string, price float64) *CreateProductCommand {
	return &CreateProductCommand{ProductID: productID, Name: name, Description: description, Price: price}
}

type UpdateProductCommand struct {
	ProductID   uuid.UUID `json:"productId" validate:"required,gte=0,lte=255"`
	Name        string    `json:"name" validate:"required,gte=0,lte=255"`
	Description string    `json:"description" validate:"required,gte=0,lte=5000"`
	Price       float64   `json:"price" validate:"required,gte=0"`
}

func NewUpdateProductCommand(productID uuid.UUID, name string, description string, price float64) *UpdateProductCommand {
	return &UpdateProductCommand{ProductID: productID, Name: name, Description: description, Price: price}
}

type DeleteProductCommand struct {
	ProductID uuid.UUID `json:"productId" validate:"required"`
}

func NewDeleteProductCommand(productID uuid.UUID) *DeleteProductCommand {
	return &DeleteProductCommand{ProductID: productID}
}
