package model

import (
	"time"

	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var validate = validator.New()

type Invoice struct {
	ID         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	CustomerID primitive.ObjectID `json:"customerId" bson:"customerId"`
	Amount     float64            `json:"amount" bson:"amount"`
	Date       string             `json:"date" bson:"date"`
	Status     string             `json:"status" bson:"status"`
	CreatedAt  time.Time          `bson:"createdAt,omitempty" json:"createdAt"`
	UpdatedAt  time.Time          `bson:"updatedAt,omitempty" json:"updatedAt"`
}

type InvoiceDTO struct {
	ID         primitive.ObjectID `json:"id" bson:"_id"`
	CustomerID primitive.ObjectID `json:"customerId,omitzero" bson:"customerId"`
	Amount     float64            `json:"amount" bson:"amount"`
	Date       string             `json:"date" bson:"date"`
	Status     string             `json:"status" bson:"status"`
	CreatedAt  time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt  time.Time          `json:"updatedAt,omitzero" bson:"updatedAt"`
}

type LatestInvoice struct {
	ID       string  `json:"id" bson:"_id"`
	Name     string  `json:"name" bson:"name"`
	ImageURL string  `json:"imageUrl" bson:"imageUrl"`
	Email    string  `json:"email" bson:"email"`
	Amount   float64 `json:"amount" bson:"amount"`
}

type CreateInvoice struct {
	CustomerID string  `json:"customerId" validate:"required"`
	Amount     float64 `json:"amount" validate:"required"`
	Date       string  `json:"date" validate:"required"`
	Status     string  `json:"status"`
}

type UpdateInvoice struct {
	CustomerID string  `json:"customerId"`
	Amount     float64 `json:"amount"`
	Date       string  `json:"date"`
	Status     string  `json:"status"`
}

type ErrorResponse struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Value string `json:"value,omitempty"`
}

func ValidateStruct[T any](payload T) []ErrorResponse {
	var errors []ErrorResponse
	err := validate.Struct(payload)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.Field = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, element)
		}
	}
	return errors
}
