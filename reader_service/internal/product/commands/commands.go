package commands

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type ProductCommands struct {
	CreateProduct CreateProductCmdHandler
	UpdateProduct UpdateProductCmdHandler
	DeleteProduct DeleteProductCmdHandler
}

func NewProductCommands(
	createProduct CreateProductCmdHandler,
	updateProduct UpdateProductCmdHandler,
	deleteProduct DeleteProductCmdHandler,
) *ProductCommands {
	return &ProductCommands{CreateProduct: createProduct, UpdateProduct: updateProduct, DeleteProduct: deleteProduct}
}

type CreateProductCommand struct {
	ProductID   string    `json:"productId" bson:"_id,omitempty"`
	Name        string    `json:"name,omitempty" bson:"name,omitempty" validate:"required,min=3,max=250"`
	Description string    `json:"description,omitempty" bson:"description,omitempty" validate:"required,min=3,max=500"`
	Price       float64   `json:"price,omitempty" bson:"price,omitempty" validate:"required"`
	CreatedAt   time.Time `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt   time.Time `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
}

func NewCreateProductCommand(productID string, name string, description string, price float64, createdAt time.Time, updatedAt time.Time) *CreateProductCommand {
	return &CreateProductCommand{ProductID: productID, Name: name, Description: description, Price: price, CreatedAt: createdAt, UpdatedAt: updatedAt}
}

type UpdateProductCommand struct {
	ProductID   string    `json:"productId" bson:"_id,omitempty"`
	Name        string    `json:"name,omitempty" bson:"name,omitempty" validate:"required,min=3,max=250"`
	Description string    `json:"description,omitempty" bson:"description,omitempty" validate:"required,min=3,max=500"`
	Price       float64   `json:"price,omitempty" bson:"price,omitempty" validate:"required"`
	UpdatedAt   time.Time `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
}

func NewUpdateProductCommand(productID string, name string, description string, price float64, updatedAt time.Time) *UpdateProductCommand {
	return &UpdateProductCommand{ProductID: productID, Name: name, Description: description, Price: price, UpdatedAt: updatedAt}
}

type DeleteProductCommand struct {
	ProductID uuid.UUID `json:"productId" bson:"_id,omitempty"`
}

func NewDeleteProductCommand(productID uuid.UUID) *DeleteProductCommand {
	return &DeleteProductCommand{ProductID: productID}
}
